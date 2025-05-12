package config

import (
	"os"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestConfig(t *testing.T) {
	Convey("Given an environment with no environment variables set", t, func() {
		os.Clearenv()
		cfg, err := Get()

		Convey("When the config values are retrieved", func() {
			Convey("Then there should be no error returned", func() {
				So(err, ShouldBeNil)
			})

			Convey("Then the values should be set to the expected defaults", func() {
				So(cfg.BindAddr, ShouldEqual, ":29600")
				So(cfg.DefaultLimit, ShouldEqual, 20)
				So(cfg.DefaultMaxLimit, ShouldEqual, 1000)
				So(cfg.DefaultOffset, ShouldEqual, 0)
				So(cfg.GracefulShutdownTimeout, ShouldEqual, 5*time.Second)
				So(cfg.HealthCheckInterval, ShouldEqual, 30*time.Second)
				So(cfg.HealthCheckCriticalTimeout, ShouldEqual, 90*time.Second)
				So(cfg.OTBatchTimeout, ShouldEqual, 5*time.Second)
				So(cfg.OTExporterOTLPEndpoint, ShouldEqual, "localhost:4317")
				So(cfg.OTServiceName, ShouldEqual, "dis-search-upstream-stub")
				So(cfg.OtelEnabled, ShouldBeFalse)
				So(cfg.Kafka.ContentUpdatedGroup, ShouldEqual, "dis-search-upstream-stub")
				So(cfg.Kafka.ContentUpdatedTopic, ShouldEqual, "content-updated")
				So(cfg.Kafka.SearchContentUpdatedTopic, ShouldEqual, "search-content-updated")
				So(cfg.Kafka.Addr, ShouldResemble, []string{"localhost:9092", "localhost:9093", "localhost:9094"})
				So(cfg.Kafka.Version, ShouldEqual, "1.0.2")
				So(cfg.Kafka.OffsetOldest, ShouldBeTrue)
				So(cfg.Kafka.NumWorkers, ShouldEqual, 1)
				So(cfg.Kafka.SecProtocol, ShouldEqual, "")
				So(cfg.Kafka.SecCACerts, ShouldEqual, "")
				So(cfg.Kafka.SecClientCert, ShouldEqual, "")
				So(cfg.Kafka.SecClientKey, ShouldEqual, "")
				So(cfg.Kafka.SecSkipVerify, ShouldBeFalse)
				So(cfg.Kafka.MaxBytes, ShouldEqual, 2000000)
				So(cfg.Kafka.ConsumerMinBrokersHealthy, ShouldEqual, 1)
				So(cfg.Kafka.ProducerMinBrokersHealthy, ShouldEqual, 1)
			})
		})
	})
}
