package main

import (
	"context"
	"os"
	"time"

	"github.com/ONSdigital/dis-search-upstream-stub/config"
	"github.com/ONSdigital/dis-search-upstream-stub/data"
	"github.com/ONSdigital/dis-search-upstream-stub/schema"
	kafka "github.com/ONSdigital/dp-kafka/v4"
	"github.com/ONSdigital/log.go/v2/log"
)

const (
	serviceName = "dp-search-upstream-stub"
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
		Topic:           cfg.Kafka.ContentUpdatedTopic,
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
	resources, err := resourceStore.GetResources(ctx, options)
	if err != nil {
		log.Error(ctx, "failed to retrieve resources", err)
		os.Exit(1)
	}

	// Marshal each resource to Kafka message format and send
	for i := range resources.Items {
		item := &resources.Items[i]
		messageBytes, err := schema.SearchContentUpdateEvent.Marshal(item)
		if err != nil {
			log.Error(ctx, "content-update event error", err)
			os.Exit(1)
		}

		// Create a Kafka BytesMessage from the byte slice
		kafkaMessage := kafka.BytesMessage{
			Value: messageBytes,
		}

		// Send the BytesMessage to Kafka
		if err := kafkaProducer.Initialise(ctx); err != nil {
			log.Warn(ctx, "failed to initialise kafka producer")
			return
		}
		kafkaProducer.Channels().Output <- kafkaMessage
		log.Info(ctx, "resource sent to Kafka", log.Data{"item": item})
	}
}
