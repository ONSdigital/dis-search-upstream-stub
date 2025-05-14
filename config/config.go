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
	SearchContentUpdatedTopic string   `envconfig:"KAFKA_SEARCH_CONTENT_UPDATED_TOPIC"`
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
			ContentUpdatedTopic:       "content-updated",
			SearchContentUpdatedTopic: "search-content-updated",
			Addr:                      []string{"localhost:9092", "localhost:9093", "localhost:9094"},
			Version:                   "2.6.1",
			OffsetOldest:              true,
			NumWorkers:                1,
			SecProtocol:               "TLS",
			SecCACerts:                "",
			SecClientCert:             "-----BEGIN CERTIFICATE-----\nMIIEcjCCAlqgAwIBAgIQQJqjYdD7p5kredGcGeKdzTANBgkqhkiG9w0BAQsFADAh\nMR8wHQYDVQQDDBZzYW5kYm94LmRwLm9ucy5wcml2YXRlMB4XDTI1MDUxMjEyNTc1\nMFoXDTI3MDUxMjEzNTc1MFowSzFJMEcGA1UEAxNAZGlzLXNlYXJjaC11cHN0cmVh\nbS1zdHViLnB1Ymxpc2hpbmcua2Fma2Euc2FuZGJveC5kcC5vbnMucHJpdmF0ZTCC\nASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAIfn2LUWWYdpqmdzgXYKhmqN\n5PHRZyeUVXRAGlrBOTouJk5tPw4Hi9btbtJF+KdV+0jEVFrh8xobe7gKjJV6qb1A\nFVCdpz2KM9Am7m9c+Fd/Lzo0zPqYCnuL3Q+0U+iIqnVJ+MsYI6SPZsUNMNCtwzO6\neP3lJWXPzzD7wVr3+3hczfDTOS/dzkHvehT1pZVzIwbCKXPYFieqzH6izbeezfVk\ngBgiNbAjtTJx5Kk2zM5F5gNrYr+I/E3BO61GkJtCxn5ES923KIlKx1iE3ZPwmDIJ\nWUYf5T1i5e2jMUNbtqC1bLMXOCkp6cKOdcbzu7y3K6D1GOx7sqBmEqPMh4YmCgEC\nAwEAAaN8MHowCQYDVR0TBAIwADAfBgNVHSMEGDAWgBSu8uIXk1nVCiwa/aK7lLLP\njruRxjAdBgNVHQ4EFgQUexvMFZ8QUKzS/vZZ0WlWk8UYVbIwDgYDVR0PAQH/BAQD\nAgWgMB0GA1UdJQQWMBQGCCsGAQUFBwMBBggrBgEFBQcDAjANBgkqhkiG9w0BAQsF\nAAOCAgEAoLW9ITGWuYddTLiKN+56lytn92dNNxSnnF+vmr0AKsH06FLWUqoV0Pbs\ntiIOt+FeMSyrjlWJt0ngWUD6cG3//UzgvjC3fKHOcIvZl38CEkGorW4VHhmFaTNz\n1iA9zfHbL9lPYTYPHRDT3zpUrDo3Er2dFCSn/7/HVAXbyyo5AVVHb5vUooiR1RRL\nTTU+p3OBDKpX6u1cq4GL/DhooVSs2eA1H/hMeb5SCxdrnhsjb2IetaskNfta2iEb\niQl13bkt86y+0XintQY+YCOsTtzICHuaSDwK5ElTSfUSIetTcCysBXSHfenbJz6+\nO6E+x2rw8R28uYLVpHMc7XRrgpyK7iv58IREfhN+BFIaoYf0UpfvBLKz0QsFc55W\nu64a6cJttKLRJVMEzKEY0GRPeQGjAeJa2m/HXjtR7rLCg/jCnPdqffb4ICRFdR2J\nTi0NRt4w5aI3Y/JDLLPBVgVAg/LR06dj6prqgbvwNMcLlfF31UESNbkw0ob6Gaqv\njT/BPmjimgmrAO7j22jTcsZZ+TGQZ2WT9GcmgpAe0r/++I3OZ1e7uObQfQHLRxFe\n+Ka1sibDLQoFtsVSUWk7Ufh0YwwekKFRip2FXyEYuwo6qwRhoWGkTa1EkKJsYKQa\nLgDbh/9fzU4zoXkmab9loZUy23QrlAqQuYlQbjTCJDXrvqfwHrY=\n-----END CERTIFICATE-----",
			SecClientKey:              "-----BEGIN PRIVATE KEY-----\nMIIEvAIBADANBgkqhkiG9w0BAQEFAASCBKYwggSiAgEAAoIBAQCH59i1FlmHaapn\nc4F2CoZqjeTx0WcnlFV0QBpawTk6LiZObT8OB4vW7W7SRfinVftIxFRa4fMaG3u4\nCoyVeqm9QBVQnac9ijPQJu5vXPhXfy86NMz6mAp7i90PtFPoiKp1SfjLGCOkj2bF\nDTDQrcMzunj95SVlz88w+8Fa9/t4XM3w0zkv3c5B73oU9aWVcyMGwilz2BYnqsx+\nos23ns31ZIAYIjWwI7UyceSpNszOReYDa2K/iPxNwTutRpCbQsZ+REvdtyiJSsdY\nhN2T8JgyCVlGH+U9YuXtozFDW7agtWyzFzgpKenCjnXG87u8tyug9Rjse7KgZhKj\nzIeGJgoBAgMBAAECggEASbiMJv7upbO9ycwwJ7Xh4B8EV+A9Uoz2Jc83/I7f2h16\njgRYteWGB5pYCRwHm83aN0i4cWwrkLfjzpt6UwNs28nNRiOeupUjMgBMSoBl/iBx\nn+NQZYbf+NCPo5swAO1RebociR6ZBwT6vF1BY5E+V+sJAsCwHqAxReLqqcvmzwzN\nddeolKlg9xxRTnXhstKLYvIlOczL3uzZ+4QNkC1FL5an/wU1KEWIjXj0CJtnWvvQ\nzoD9y4h4HbFa2iz3x+iJaixkuAOdiSPP8POByHPt0HsU8EKYYT7bqTvPM2IpHkyr\nDPlHmF33DQrG0oH6nN1uqRzkzLCtotiS0LoisZhi+QKBgQDJuzJjYKFVVKL/jKuj\n4dfxjz5cHrzc2leE24dmfwnrPPtWZdraRRT4lwFgh10kTCHW1+mNvI94XrdQJzYl\ndPSu3P2jO20VBT4DAIU4P525XSqTumPme9NaQLVgb9su6ey/GrYHyYmlCk9s5hh9\n2iQ7CJoeLXMyHctFw8gid/Ir5wKBgQCsd2GBZUd2xwmBUxuAKaUzh7EjneF5IoRB\nXphkf1HGlinTJCYOzU4DqSAN6EYqQWTxYkBmhzbXuZEnOXIJjb1edIPkAvDBuFIB\nVO9WUaWftBk/vUCtr92oKZlByfGOLWNR2wCIBza7ocBLvR5ggII8k9YXLeQem2sS\n3FLGpaYd1wKBgGgmygEc5q7Tn8QosIVQGNmShzOwevnbkMv7O5Djjg9x0KHuvGts\nt0MRU5iuypvu4pm1p9ORwtD2tdYgKIh2Nc4CMsGP8OWlazrJjf5Yeeo1+8GBvgpF\na/1w4zQDDDrQc3bHJ6wllXcsN42VzpdLhOEls8xY0tzRHR3L0wxYuSOJAoGAYHHu\nI/coLKMHjLuV8GjZimSCScGbeis0PH4SyHhumZgV0Y4wfiyPSPrGAyD2Q+EH+viP\nvQY2RBLwujekrvUFhhGwQ8zlJ9/UdAw0P1gvP4zuZbeGuNpVIRoKK0EsBO8a0Iag\n2HD4SZsdtv0ORLb4nbmqipHONNOC4Cw3WgD+UUcCgYBU1kJzSG9MwDWuQyPXmnbp\nCRtD3zVNr8+GBzHXnh97cQ8jqztpQHtq2oBS4W9km+OOwpIWSF4IgznsYI3i1VoF\nMVcMyrDlr60n4RRThn/og4CZCaN3N3z639DRVMnVeJr0SFanNqpQ8T++JFPsjj59\nrfXuwaJF+cQiMJcLjInhCQ==\n-----END PRIVATE KEY-----",
			SecSkipVerify:             false,
			MaxBytes:                  2000000,
			ConsumerMinBrokersHealthy: 1,
			ProducerMinBrokersHealthy: 1,
		},
	}

	return cfg, envconfig.Process("", cfg)
}
