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

const (
	serviceName  = "dp-search-upstream-stub"
	resourceType = "SearchContentUpdatedResource"
)

func main() {
	log.Namespace = serviceName
	ctx := context.Background()

	// Get Config
	cfg, err := config.Get()
	if err != nil {
		log.Error(ctx, "error getting config", err)
		os.Exit(1)
	}

	// Create Kafka Producer
	pConfig := &kafka.ProducerConfig{
		BrokerAddrs:     cfg.Kafka.Addr,
		Topic:           cfg.Kafka.SearchContentUpdatedTopic,
		KafkaVersion:    &cfg.Kafka.Version,
		MaxMessageBytes: &cfg.Kafka.MaxBytes,
	}
	if cfg.Kafka.SecProtocol == config.KafkaTLSProtocol {
		pConfig.SecurityConfig = kafka.GetSecurityConfig(
			cfg.Kafka.SecCACerts,
			cfg.Kafka.SecClientCert,
			cfg.Kafka.SecClientKey,
			cfg.Kafka.SecSkipVerify,
		)
	}
	kafkaProducer, err := kafka.NewProducer(ctx, pConfig)
	if err != nil {
		log.Error(ctx, "fatal error trying to create kafka producer", err, log.Data{"topic": cfg.Kafka.ContentUpdatedTopic})
		os.Exit(1)
	}

	// kafka error logging go-routines
	kafkaProducer.LogErrors(ctx)

	time.Sleep(500 * time.Millisecond)

	// Initialize ResourceStore
	resourceStore := &data.ResourceStore{}

	// Define Options for fetching resources
	options := data.Options{
		Offset: 0,
		Limit:  100,
	}

	// Call GetResources to fetch resources
	resources, err := resourceStore.GetResourcesWithType(ctx, resourceType, options)
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
		_, err := fmt.Scanln(&selection)
		if err != nil || selection < 1 || selection > len(resources.Items) { // Adjust range for 1-based indexing
			fmt.Println("Invalid selection. Please try again.")
			continue
		}
		break
	}

	// Get the selected item (adjust for 1-based indexing)
	selectedItem := &resources.Items[selection-1]

	// Assert the type of the selected item
	var messageBytes []byte
	switch r := (*selectedItem).(type) {
	case models.ContentUpdatedResource:
		// Marshal the ContentUpdatedResource
		messageBytes, err = schema.ContentPublishedEvent.Marshal(r)
		if err != nil {
			log.Error(ctx, "content-update event error", err)
			os.Exit(1)
		}

	case models.SearchContentUpdatedResource:
		// Marshal the SearchContentUpdatedResource
		messageBytes, err = schema.SearchContentUpdateEvent.Marshal(r)
		if err != nil {
			log.Error(ctx, "search-content-update event error", err)
			os.Exit(1)
		}

	default:
		log.Error(ctx, "unsupported resource type", nil)
		os.Exit(1)
	}

	// Create a Kafka BytesMessage from the byte slice
	kafkaMessage := kafka.BytesMessage{
		Value: messageBytes,
	}

	// Send message to Kafka
	kafkaProducer.Channels().Output <- kafkaMessage
	// Log the actual message being sent
	log.Info(context.Background(), "message sent to Kafka", log.Data{"message": string(messageBytes)})

	log.Info(context.Background(), "resource sent to Kafka", log.Data{"item": selectedItem})
}
