package main

import (
	"context"
	"encoding/json"
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
	// keep main tiny to satisfy gocyclo
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	log.Namespace = serviceName
	ctx := context.Background()

	// Get Config
	cfg, err := config.Get()
	if err != nil {
		log.Error(ctx, "error getting config", err)
		return err
	}

	// Prompt and map to resource type + topic + label
	topic, resourceType, eventLabel, err := promptTopicAndMap(cfg)
	if err != nil {
		return err
	}
	log.Info(ctx, "Selected topic", log.Data{
		"topic":         topic,
		"resource_type": resourceType,
		"event_type":    eventLabel,
	})

	// Fetch resources
	resourceStore := &data.ResourceStore{}
	resources, err := resourceStore.GetResourcesWithType(ctx, resourceType, data.Options{Offset: 0, Limit: 100})
	if err != nil {
		log.Error(ctx, "failed to retrieve resources", err)
		return err
	}
	if len(resources.Items) == 0 {
		fmt.Println("No resources available for this type.")
		return nil
	}

	// Show and select one
	printResources(resources.Items)
	idx, err := promptSelection(len(resources.Items))
	if err != nil {
		return err
	}
	selected := resources.Items[idx]

	// Marshal payload (legacy Avro, new JSON)
	payload, eventType, err := marshalPayload(selected)
	if err != nil {
		log.Error(ctx, "failed to marshal event", err, log.Data{"event_type": eventType})
		return err
	}

	// Producer setup
	producer, err := newProducerForTopic(ctx, cfg.Kafka, topic)
	if err != nil {
		log.Error(ctx, "fatal error creating kafka producer", err, log.Data{"topic": topic})
		return err
	}
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if cerr := producer.Close(shutdownCtx); cerr != nil {
			log.Error(ctx, "failed to close kafka producer", cerr)
		}
	}()

	producer.LogErrors(ctx)
	if err := producer.Initialise(ctx); err != nil {
		log.Error(ctx, "failed to initialise kafka producer", err)
		return err
	}

	// Send
	producer.Channels().Output <- kafka.BytesMessage{Value: payload}
	log.Info(ctx, "resource sent to Kafka", log.Data{
		"event_type": eventType,
		"topic":      topic,
		"resource":   fmt.Sprintf("%T", selected),
	})

	return nil
}

func promptTopicAndMap(cfg *config.Config) (topic, resourceType, eventLabel string, err error) {
	fmt.Println("Select the Kafka topic to send messages to:")
	fmt.Println("1) content-updated (legacy)")
	fmt.Println("2) search-content-updated (new)")
	fmt.Println("3) search-content-deleted (new)")

	var choice int
	for {
		fmt.Print("Enter choice (1-3): ")
		if _, err = fmt.Scanln(&choice); err != nil || choice < 1 || choice > 3 {
			fmt.Println("Invalid selection. Please try again.")
			continue
		}
		break
	}

	switch choice {
	case 1:
		return cfg.Kafka.ContentUpdatedTopic, "ContentUpdatedResource", "ContentUpdated(AVRO)", nil
	case 2:
		return cfg.Kafka.SearchContentUpdatedTopic, "SearchContentUpdatedResource", "SearchContentUpdated(JSON)", nil
	default:
		return cfg.Kafka.SearchContentDeletedTopic, "SearchContentDeletedResource", "SearchContentDeleted(JSON)", nil
	}
}

func printResources(items []models.Resource) {
	fmt.Println("Available resources:")
	for i, item := range items {
		switch r := item.(type) {
		case models.ContentUpdatedResource:
			fmt.Printf("[%d] URL: %s, CollectionID: %s, Data Type: %s\n", i+1, r.URI, r.CollectionID, r.DataType)
		case models.SearchContentUpdatedResource:
			fmt.Printf("[%d] Title: %s, URL: %s, Content type: %s\n", i+1, r.Title, r.URI, r.ContentType)
		case models.SearchContentDeletedResource:
			fmt.Printf("[%d] URL: %s, CollectionID: %s\n", i+1, r.URI, r.CollectionID)
		default:
			fmt.Printf("[%d] Unknown resource type\n", i+1)
		}
	}
}

func promptSelection(count int) (int, error) {
	var selection int
	for {
		fmt.Print("Enter the number of the resource to send: ")
		if _, err := fmt.Scanln(&selection); err != nil || selection < 1 || selection > count {
			fmt.Println("Invalid selection. Please try again.")
			continue
		}
		return selection - 1, nil
	}
}

func marshalPayload(item models.Resource) (payload []byte, eventType string, err error) {
	switch r := item.(type) {
	case models.ContentUpdatedResource:
		b, err := schema.ContentPublishedEvent.Marshal(r) // legacy Avro
		return b, "ContentPublishedEvent", err
	case models.SearchContentUpdatedResource:
		b, err := json.Marshal(r) // JSON
		return b, "SearchContentUpdatedEvent", err
	case models.SearchContentDeletedResource:
		b, err := json.Marshal(r) // JSON
		return b, "SearchContentDeletedEvent", err
	default:
		return nil, "UnknownEvent", fmt.Errorf("unsupported resource type: %T", item)
	}
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
