package main

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/ONSdigital/dis-search-upstream-stub/config"
	"github.com/ONSdigital/dis-search-upstream-stub/models"
	kafka "github.com/ONSdigital/dp-kafka/v4"
	"github.com/ONSdigital/log.go/v2/log"
)

const (
	serviceName = "dp-search-upstream-stub"
	dataDir     = "data" // Directory with JSON files
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

	// Process each JSON file in the data directory
	files, err := os.ReadDir(dataDir)
	if err != nil {
		log.Error(ctx, "error reading data directory", err)
		os.Exit(1)
	}

	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".json" {
			continue
		}

		filePath := filepath.Join(dataDir, file.Name())
		log.Info(ctx, "processing file", log.Data{"file": filePath})

		// Read file content
		content, err := os.ReadFile(filePath)
		if err != nil {
			log.Error(ctx, "error reading file", err, log.Data{"file": filePath})
			continue
		}

		// Unmarshal into Resources struct
		var resources models.Resources
		if err := json.Unmarshal(content, &resources); err != nil {
			log.Error(ctx, "error unmarshalling JSON", err, log.Data{"file": filePath})
			continue
		}

		log.Info(ctx, "successfully parsed resources", log.Data{"resources": resources})

		// Marshal each resource to Kafka message format and send
		for _, item := range resources.Items {
			var messageBytes []byte
			messageBytes, err = json.Marshal(item)
			if err != nil {
				log.Error(ctx, "error marshalling resource to JSON", err, log.Data{"item": item})
				continue
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

	log.Info(ctx, "completed processing all JSON files")
	time.Sleep(500 * time.Millisecond) // Ensure all messages are sent
}
