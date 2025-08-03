package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/ONSdigital/dis-search-upstream-stub/config"
	"github.com/ONSdigital/dis-search-upstream-stub/data"
	"github.com/ONSdigital/dis-search-upstream-stub/models"
	"github.com/ONSdigital/dis-search-upstream-stub/schema"
	kafka "github.com/ONSdigital/dp-kafka/v4"
	"github.com/ONSdigital/log.go/v2/log"
)

const serviceName = "dis-search-upstream-stub"

func main() {
	log.Namespace = serviceName
	ctx := context.Background()

	// Get Config
	cfg, err := config.Get()
	if err != nil {
		log.Error(ctx, "error getting config", err)
		os.Exit(1)
	}

	// Call GetResources to fetch resources
	resourceStore := &data.ResourceStore{}
	resources, err := resourceStore.GetResources(ctx, "", data.Options{Offset: 0, Limit: 100})
	if err != nil {
		log.Error(ctx, "failed to retrieve resources", err)
		os.Exit(1)
	}

	// Display the list of resources with title and URI
	fmt.Println("Available resources:")
	for i, item := range resources.Items {
		switch r := item.(type) {
		case models.ContentUpdatedResource:
			// Display ContentUpdatedResource fields
			fmt.Printf("[%d] URL: %s, CollectionID: %s, Data Type: %s\n", i+1, r.URI, r.CollectionID, r.DataType)
		case models.SearchContentUpdatedResource:
			// Display SearchContentUpdatedResource fields
			fmt.Printf("[%d] Title: %s, URL: %s, Content type: %s\n", i+1, r.Title, r.URI, r.ContentType)
		default:
			// Default case in case of unknown resource type
			fmt.Printf("[%d] Unknown resource type\n", i+1)
		}
	}

	// Ask the user to select a resource
	var selection int
	for {
		fmt.Print("Enter the number of the resource to send: ")
		if _, err := fmt.Scanln(&selection); err != nil || selection < 1 || selection > len(resources.Items) {
			fmt.Println("Invalid selection. Please try again.")
			continue
		}
		break
	}
	selectedItem := resources.Items[selection-1]

	// Decide topic + marshal payload
	var (
		topic        string
		messageBytes []byte
		eventType    string
	)

	// Marshal the resource to Kafka message format based on the topic type
	switch r := selectedItem.(type) {
	case models.ContentUpdatedResource:
		// Marshal for the legacy topic (ContentPublishedEvent)
		topic = cfg.Kafka.ContentUpdatedTopic
		messageBytes, err = schema.ContentPublishedEvent.Marshal(r)
		eventType = "ContentPublishedEvent"
	case models.SearchContentUpdatedResource:
		// Marshal for the new topic (SearchContentUpdateEvent)
		topic = cfg.Kafka.SearchContentUpdatedTopic
		messageBytes, err = schema.SearchContentUpdateEvent.Marshal(r)
		eventType = "SearchContentUpdateEvent"
	default:
		log.Error(context.Background(), "unsupported resource type", fmt.Errorf("resource type: %v", selectedItem))
		return
	}

	if err != nil {
		log.Error(ctx, "failed to marshal event", err, log.Data{"event_type": eventType})
		os.Exit(1)
	}

	// Build a producer for the chosen topic
	producer, err := newProducerForTopic(ctx, cfg.Kafka, topic)
	if err != nil {
		log.Error(ctx, "fatal error creating kafka producer", err, log.Data{"topic": topic})
		os.Exit(1)
	}
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if cerr := producer.Close(shutdownCtx); cerr != nil {
			log.Error(ctx, "failed to close kafka producer", cerr)
		}
	}()

	producer.LogErrors(ctx)

	// Initialise and send
	if err := producer.Initialise(ctx); err != nil {
		log.Error(ctx, "failed to initialise kafka producer", err)
	}

	producer.Channels().Output <- kafka.BytesMessage{Value: messageBytes}

	log.Info(ctx, "resource sent to Kafka", log.Data{
		"event_type": eventType,
		"topic":      topic,
	})
}

func newProducerForTopic(ctx context.Context, kcfg *config.Kafka, topic string) (kafka.IProducer, error) {
	pcfg := &kafka.ProducerConfig{
		BrokerAddrs:     kcfg.Addr,
		Topic:           topic,
		KafkaVersion:    &kcfg.Version,
		MaxMessageBytes: &kcfg.MaxBytes,
	}
	if kcfg.SecProtocol == config.KafkaTLSProtocol {
		pcfg.SecurityConfig = kafka.GetSecurityConfig(
			kcfg.SecCACerts,
			kcfg.SecClientCert,
			kcfg.SecClientKey,
			kcfg.SecSkipVerify,
		)
	}
	return kafka.NewProducer(ctx, pcfg)
}
