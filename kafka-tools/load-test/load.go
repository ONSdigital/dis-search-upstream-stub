package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/ONSdigital/dis-search-upstream-stub/config"
	"github.com/ONSdigital/dis-search-upstream-stub/data"
	"github.com/ONSdigital/dis-search-upstream-stub/models"
	"github.com/ONSdigital/dis-search-upstream-stub/schema"
	kafka "github.com/ONSdigital/dp-kafka/v4"
	"github.com/ONSdigital/log.go/v2/log"
)

const (
	serviceName = "dis-search-upstream-stub"
	// Default values for message counts
	defaultLegacyContentUpdatedMessages = 1
	defaultSearchContentUpdatedMessages = 1
	defaultSearchContentDeletedMessages = 1
)

type producers struct {
	ContentUpdatedProducer       *kafka.Producer
	SearchContentUpdatedProducer *kafka.Producer
	SearchContentDeletedProducer *kafka.Producer
}

type resources struct {
	contentUpdatedResources       *models.Resources
	searchContentUpdatedResources *models.Resources
	searchContentDeletedResources *models.Resources
}

func makeStubTraceID(loopIndex int) string {
	return fmt.Sprintf("stub-%d-%d", time.Now().UnixMilli(), loopIndex)
}

func sendMessageToKafka(producer *kafka.Producer, item models.Resource, loopIndex int, wg *sync.WaitGroup) {
	defer wg.Done()

	var (
		messageBytes []byte
		err          error
		eventType    string
		traceID      string
	)

	// Marshal the resource to Kafka message format based on the topic type
	switch r := item.(type) {
	case models.ContentUpdatedResource:
		// Marshal for the legacy topic (ContentPublishedEvent)
		if r.TraceID == "" {
			r.TraceID = makeStubTraceID(loopIndex)
		}
		traceID = r.TraceID
		messageBytes, err = schema.ContentPublishedEvent.Marshal(r)
		eventType = "ContentPublishedEvent(AVRO)"
	case models.SearchContentUpdatedResource:
		// Marshal for the new topic (SearchContentUpdatedEvent)
		if r.TraceID == "" {
			r.TraceID = makeStubTraceID(loopIndex)
		}
		traceID = r.TraceID
		messageBytes, err = json.Marshal(r)
		eventType = "SearchContentUpdatedEvent(JSON)"
	case models.SearchContentDeletedResource:
		// Marshal for the new topic (SearchContentDeletedEvent)
		if r.TraceID == "" {
			r.TraceID = makeStubTraceID(loopIndex)
		}
		traceID = r.TraceID
		messageBytes, err = json.Marshal(r)
		eventType = "SearchContentDeleteEvent(JSON)"
	default:
		log.Error(context.Background(), "unsupported resource type", fmt.Errorf("resource type: %v", item))
		return
	}

	if err != nil {
		log.Error(context.Background(), "event marshal error", err, log.Data{"event_type": eventType})
		return
	}

	// Send message to Kafka
	producer.Channels().Output <- kafka.BytesMessage{Value: messageBytes}
	log.Info(context.Background(), "resource sent to Kafka", log.Data{
		"event_type": eventType,
		"trace_id":   traceID,
		"item":       item,
	})
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

func buildAndInitProducers(ctx context.Context, cfg *config.Config) (producers, error) {
	// inner helper: create + log errors + initialise
	mk := func(topic string) (*kafka.Producer, error) {
		p, err := createKafkaProducer(ctx, cfg, topic)
		if err != nil {
			return nil, err
		}
		p.LogErrors(ctx)
		if err := p.Initialise(ctx); err != nil {
			_ = p.Close(ctx) // best-effort cleanup on init failure
			return nil, err
		}
		return p, nil
	}

	contentUpdated, err := mk(cfg.Kafka.ContentUpdatedTopic)
	if err != nil {
		return producers{}, fmt.Errorf("legacy content-updated producer: %w", err)
	}
	searchContentUpdated, err := mk(cfg.Kafka.SearchContentUpdatedTopic)
	if err != nil {
		return producers{}, fmt.Errorf("search-content-updated producer: %w", err)
	}
	searchContentDeleted, err := mk(cfg.Kafka.SearchContentDeletedTopic)
	if err != nil {
		return producers{}, fmt.Errorf("search-content-deleted producer: %w", err)
	}

	return producers{ContentUpdatedProducer: contentUpdated, SearchContentUpdatedProducer: searchContentUpdated, SearchContentDeletedProducer: searchContentDeleted}, nil
}

func fetchAllResources(ctx context.Context) (resources, error) {
	rs := &data.ResourceStore{}
	opts := data.Options{Offset: 0, Limit: 100}

	contentUpdated, err := rs.GetResourcesWithType(ctx, "ContentUpdatedResource", opts)
	if err != nil {
		return resources{}, fmt.Errorf("failed to retrieve legacy content-updated resources: %w", err)
	}

	searchContentUpdated, err := rs.GetResourcesWithType(ctx, "SearchContentUpdatedResource", opts)
	if err != nil {
		return resources{}, fmt.Errorf("failed to retrieve search-content-updated resources: %w", err)
	}

	searchContentDeleted, err := rs.GetResourcesWithType(ctx, "SearchContentDeletedResource", opts)
	if err != nil {
		return resources{}, fmt.Errorf("failed to retrieve search-content-deleted resources: %w", err)
	}

	return resources{
		contentUpdatedResources:       contentUpdated,
		searchContentUpdatedResources: searchContentUpdated,
		searchContentDeletedResources: searchContentDeleted,
	}, nil
}

func main() {
	log.Namespace = serviceName
	ctx := context.Background()

	// Define flags for how many messages to send per topic
	var numContentUpdated, numSearchContentUpdated, numSearchContentDeleted int
	flag.IntVar(&numContentUpdated, "num-content-updated", defaultLegacyContentUpdatedMessages,
		"Number of legacy 'content-updated' messages to send (Avro)")
	flag.IntVar(&numSearchContentUpdated, "num-search-content-updated", defaultSearchContentUpdatedMessages,
		"Number of 'search-content-updated' messages to send (JSON)")
	flag.IntVar(&numSearchContentDeleted, "num-search-content-deleted", defaultSearchContentDeletedMessages,
		"Number of 'search-content-deleted' messages to send (JSON)")
	flag.Parse()

	// Get Config
	cfg, err := config.Get()
	if err != nil {
		log.Error(ctx, "error getting config", err)
		os.Exit(1)
	}

	log.Info(ctx, "Script config", log.Data{"cfg": cfg})

	prods, err := buildAndInitProducers(ctx, cfg)
	if err != nil {
		log.Error(ctx, "failed to create/initialise producers", err)
		os.Exit(1)
	}
	log.Info(ctx, "kafka producers initialised", log.Data{})

	res, err := fetchAllResources(ctx)
	if err != nil {
		log.Error(ctx, "failed to fetch resources", err)
		os.Exit(1)
	}

	// Validate message numbers
	if numContentUpdated < 0 || numSearchContentUpdated < 0 || numSearchContentDeleted < 0 {
		log.Error(ctx, "invalid message numbers for topics", nil)
		os.Exit(1)
	}

	// Guard: avoid modulo by zero if caller asked for >0, but we have no resources
	if numContentUpdated > 0 && len(res.contentUpdatedResources.Items) == 0 {
		log.Error(ctx, "no content-updated resources available but requested legacy count > 0", nil)
		numContentUpdated = 0
	}
	if numSearchContentUpdated > 0 && len(res.searchContentUpdatedResources.Items) == 0 {
		log.Error(ctx, "no search-content-updated resources available but requested new count > 0", nil)
		numSearchContentUpdated = 0
	}
	if numSearchContentDeleted > 0 && len(res.searchContentDeletedResources.Items) == 0 {
		log.Error(ctx, "no search-content-deleted resources available but requested deleted count > 0", nil)
		numSearchContentDeleted = 0
	}

	// Create a WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup

	// Send messages to legacy topic concurrently
	for i := 0; i < numContentUpdated; i++ {
		item := &res.contentUpdatedResources.Items[i%len(res.contentUpdatedResources.Items)]
		wg.Add(1)
		go sendMessageToKafka(prods.ContentUpdatedProducer, *item, i, &wg)
	}

	// Send messages to search-content-updated concurrently
	for i := 0; i < numSearchContentUpdated; i++ {
		item := &res.searchContentUpdatedResources.Items[i%len(res.searchContentUpdatedResources.Items)]
		wg.Add(1)
		go sendMessageToKafka(prods.SearchContentUpdatedProducer, *item, i, &wg)
	}

	// Send messages to search-content-deleted concurrently
	for i := 0; i < numSearchContentDeleted; i++ {
		item := &res.searchContentDeletedResources.Items[i%len(res.searchContentDeletedResources.Items)]
		wg.Add(1)
		go sendMessageToKafka(prods.SearchContentDeletedProducer, *item, i, &wg)
	}

	// Wait for all messages to be processed
	wg.Wait()

	// Close producers to flush and release resources
	if err := prods.ContentUpdatedProducer.Close(ctx); err != nil {
		log.Error(ctx, "error closing legacy kafka producer", err)
	}
	if err := prods.SearchContentUpdatedProducer.Close(ctx); err != nil {
		log.Error(ctx, "error closing updated kafka producer", err)
	}
	if err := prods.SearchContentDeletedProducer.Close(ctx); err != nil {
		log.Error(ctx, "error closing deleted kafka producer", err)
	}

	log.Info(ctx, "done sending messages", log.Data{})
}
