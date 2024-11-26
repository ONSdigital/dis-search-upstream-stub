package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

// KafkaTLSProtocol informs service to use TLS protocol for kafka
const KafkaTLSProtocol = "TLS"

// Config represents service configuration for dis-search-upstream-stub
type Config struct {
	BindAddr                   string        `envconfig:"BIND_ADDR"`
	DefaultLimit               int           `envconfig:"DEFAULT_LIMIT"`
	DefaultMaxLimit            int           `envconfig:"DEFAULT_MAXIMUM_LIMIT"`
	DefaultOffset              int           `envconfig:"DEFAULT_OFFSET"`
	GracefulShutdownTimeout    time.Duration `envconfig:"GRACEFUL_SHUTDOWN_TIMEOUT"`
	HealthCheckInterval        time.Duration `envconfig:"HEALTHCHECK_INTERVAL"`
	HealthCheckCriticalTimeout time.Duration `envconfig:"HEALTHCHECK_CRITICAL_TIMEOUT"`
	OTBatchTimeout             time.Duration `encconfig:"OTEL_BATCH_TIMEOUT"`
	OTExporterOTLPEndpoint     string        `envconfig:"OTEL_EXPORTER_OTLP_ENDPOINT"`
	OTServiceName              string        `envconfig:"OTEL_SERVICE_NAME"`
	OtelEnabled                bool          `envconfig:"OTEL_ENABLED"`
	Kafka                      *Kafka
}

// Kafka contains the config required to connect to Kafka
type Kafka struct {
	ContentUpdatedGroup       string   `envconfig:"KAFKA_CONTENT_UPDATED_GROUP"`
	ContentUpdatedTopic       string   `envconfig:"KAFKA_CONTENT_UPDATED_TOPIC"`
	Addr                      []string `envconfig:"KAFKA_ADDR"`
	Version                   string   `envconfig:"KAFKA_VERSION"`
	OffsetOldest              bool     `envconfig:"KAFKA_OFFSET_OLDEST"`
	NumWorkers                int      `envconfig:"KAFKA_NUM_WORKERS"`
	SecProtocol               string   `envconfig:"KAFKA_SEC_PROTO"`
	SecCACerts                string   `envconfig:"KAFKA_SEC_CA_CERTS"            json:"-"`
	SecClientCert             string   `envconfig:"KAFKA_SEC_CLIENT_CERT"         json:"-"`
	SecClientKey              string   `envconfig:"KAFKA_SEC_CLIENT_KEY"          json:"-"`
	SecSkipVerify             bool     `envconfig:"KAFKA_SEC_SKIP_VERIFY"`
	MaxBytes                  int      `envconfig:"KAFKA_MAX_BYTES"`
	ConsumerMinBrokersHealthy int      `envconfig:"KAFKA_CONSUMER_MIN_BROKERS_HEALTHY"`
	ProducerMinBrokersHealthy int      `envconfig:"KAFKA_PRODUCER_MIN_BROKERS_HEALTHY"`
}

var cfg *Config

// Get returns the default config with any modifications through environment
// variables
func Get() (*Config, error) {
	if cfg != nil {
		return cfg, nil
	}

	cfg = &Config{
		BindAddr:                   ":29600",
		DefaultLimit:               20,
		DefaultMaxLimit:            1000,
		DefaultOffset:              0,
		GracefulShutdownTimeout:    5 * time.Second,
		HealthCheckInterval:        30 * time.Second,
		HealthCheckCriticalTimeout: 90 * time.Second,
		OTBatchTimeout:             5 * time.Second,
		OTExporterOTLPEndpoint:     "localhost:4317",
		OTServiceName:              "dis-search-upstream-stub",
		OtelEnabled:                false,
		Kafka: &Kafka{
			ContentUpdatedGroup:       "dis-search-upstream-stub",
			ContentUpdatedTopic:       "search-content-updated",
			Addr:                      []string{"localhost:9092", "localhost:9093", "localhost:9094"},
			Version:                   "1.0.0",
			OffsetOldest:              true,
			NumWorkers:                1,
			SecProtocol:               "",
			SecCACerts:                "",
			SecClientCert:             "",
			SecClientKey:              "",
			SecSkipVerify:             false,
			MaxBytes:                  2000000,
			ConsumerMinBrokersHealthy: 1,
			ProducerMinBrokersHealthy: 1,
		},
	}

	return cfg, envconfig.Process("", cfg)
}
