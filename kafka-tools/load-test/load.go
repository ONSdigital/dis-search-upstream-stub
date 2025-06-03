package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"sync"

	"github.com/ONSdigital/dis-search-upstream-stub/config"
	"github.com/ONSdigital/dis-search-upstream-stub/data"
	"github.com/ONSdigital/dis-search-upstream-stub/models"
	"github.com/ONSdigital/dis-search-upstream-stub/schema"
	kafka "github.com/ONSdigital/dp-kafka/v4"
	"github.com/ONSdigital/log.go/v2/log"
)

const (
	serviceName          = "dis-search-upstream-stub"
	KafkaTLSProtocolFlag = "TLS"
	// Default values for message counts
	defaultLegacyMessages = 6000
	defaultNewMessages    = 100
)

func sendMessageToKafka(producer *kafka.Producer, item models.Resource, wg *sync.WaitGroup) {
	defer wg.Done()

	var messageBytes []byte
	var err error
	var eventType string

	// Marshal the resource to Kafka message format based on the topic type
	switch r := item.(type) {
	case models.ContentUpdatedResource:
		// Marshal for the legacy topic (ContentPublishedEvent)
		messageBytes, err = schema.ContentPublishedEvent.Marshal(r)
		eventType = "ContentPublishedEvent"
	case models.SearchContentUpdatedResource:
		// Marshal for the new topic (SearchContentUpdateEvent)
		messageBytes, err = schema.SearchContentUpdateEvent.Marshal(r)
		eventType = "SearchContentUpdateEvent"
	default:
		log.Error(context.Background(), "unsupported resource type", fmt.Errorf("resource type: %v", item))
		return
	}

	if err != nil {
		log.Error(context.Background(), "update event error", err)
		return
	}

	// Create a Kafka BytesMessage
	kafkaMessage := kafka.BytesMessage{
		Value: messageBytes,
	}

	// Send message to Kafka
	producer.Channels().Output <- kafkaMessage
	log.Info(context.Background(), "resource sent to Kafka", log.Data{
		"item":       item,
		"event_type": eventType,
	})
}

func main() {
	log.Namespace = serviceName
	ctx := context.Background()

	// Define flags for MESSAGE_COUNT_LEGACY and MESSAGE_COUNT_NEW
	var messageCountLegacy, messageCountNew int
	flag.IntVar(&messageCountLegacy, "legacy", defaultLegacyMessages, "Number of messages for the legacy topic (default 10)")
	flag.IntVar(&messageCountNew, "new", defaultNewMessages, "Number of messages for the new topic (default 10)")
	flag.Parse()

	// Get Config
	cfg, err := config.Get()
	if err != nil {
		log.Error(ctx, "error getting config", err)
		os.Exit(1)
	}

	log.Info(ctx, "Script config", log.Data{"cfg": cfg})

	// Create Kafka Producers for both topics
	legacyKafkaProducer, err := createKafkaProducer(ctx, cfg, cfg.Kafka.ContentUpdatedTopic)
	if err != nil {
		log.Error(ctx, "failed to create kafka producer for legacy topic", err)
		os.Exit(1)
	}

	newKafkaProducer, err := createKafkaProducer(ctx, cfg, cfg.Kafka.SearchContentUpdatedTopic)
	if err != nil {
		log.Error(ctx, "failed to create kafka producer for new topic", err)
		os.Exit(1)
	}

	// kafka error logging go-routines
	legacyKafkaProducer.LogErrors(ctx)
	newKafkaProducer.LogErrors(ctx)

	// Initialize ResourceStore
	resourceStore := &data.ResourceStore{}

	// Define Options for fetching resources
	options := data.Options{
		Offset: 0,
		Limit:  100,
	}

	// Call GetResources to fetch legacy resources
	legacyResources, err := resourceStore.GetResourcesWithType(ctx, "ContentUpdatedResource", options)
	if err != nil {
		log.Error(ctx, "failed to retrieve resources", err)
		os.Exit(1)
	}

	// Call GetResources to fetch new resources
	newResources, err := resourceStore.GetResourcesWithType(ctx, "SearchContentUpdatedResource", options)
	if err != nil {
		log.Error(ctx, "failed to retrieve resources", err)
		os.Exit(1)
	}

	// Validate passed message count values
	if messageCountLegacy <= 0 || messageCountNew <= 0 {
		log.Error(ctx, "invalid message counts for topics", nil)
		os.Exit(1)
	}

	// Create a WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup

	// Send messages to legacy topic concurrently
	for i := 0; i < messageCountLegacy; i++ {
		// Use resources in a cyclic manner
		item := &legacyResources.Items[i%len(legacyResources.Items)] // Reusing items for simplicity

		// Send message to legacy topic concurrently
		wg.Add(1)
		go sendMessageToKafka(legacyKafkaProducer, *item, &wg)
	}

	// Send messages to new topic concurrently
	for i := 0; i < messageCountNew; i++ {
		// Use resources in a cyclic manner
		item := &newResources.Items[i%len(newResources.Items)] // Reusing items for simplicity

		// Send message to new topic concurrently
		wg.Add(1)
		go sendMessageToKafka(newKafkaProducer, *item, &wg)
	}

	// Wait for all messages to be processed
	wg.Wait()
}

func createKafkaProducer(ctx context.Context, cfg *config.Config, topic string) (*kafka.Producer, error) {
	// Create Kafka producer configuration
	producerConfig := &kafka.ProducerConfig{
		BrokerAddrs:     cfg.Kafka.Addr,
		Topic:           topic,
		KafkaVersion:    &cfg.Kafka.Version,
		MaxMessageBytes: &cfg.Kafka.MaxBytes,
	}

	// Set security config if needed
	if cfg.Kafka.SecProtocol == config.KafkaTLSProtocol {
		producerConfig.SecurityConfig = kafka.GetSecurityConfig(
			cfg.Kafka.SecCACerts,
			cfg.Kafka.SecClientCert,
			cfg.Kafka.SecClientKey,
			cfg.Kafka.SecSkipVerify,
		)
	}

	// Create and return the Kafka producer
	producer, err := kafka.NewProducer(ctx, producerConfig)
	if err != nil {
		return nil, fmt.Errorf("fatal error trying to create kafka producer for topic %s: %w", topic, err)
	}

	return producer, nil
}
