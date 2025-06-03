# dis-search-upstream-stub

Fake upstream search stub to mimic upstream services using a generic search contract to integrate with our search stack

## Getting started

* Run `make debug` to run application on http://localhost:29600
* Run `make help` to see full list of make targets

### Dependencies

* golang 1.23.x
* No further dependencies other than those defined in `go.mod`

To run `make validate-specification` you require Node v20.x and to install @redocly/cli:

```sh
   npm install -g @redocly/cli
```

### Configuration

| Environment variable         | Default                  | Description                                                                                                        |
|------------------------------|--------------------------|--------------------------------------------------------------------------------------------------------------------|
| BIND_ADDR                    | :29600                   | The host and port to bind to                                                                                       |
| DEFAULT_LIMIT                | 20                       | The default number of items to be returned from a list endpoint                                                    |
| DEFAULT_MAXIMUM_LIMIT        | 1000                     | The maximum number of items to be returned in any list endpoint (to prevent performance issues)                    |
| DEFAULT_OFFSET               | 0                        | The number of items into the full list (i.e. the 0-based index) that a particular response is starting at          |
| GRACEFUL_SHUTDOWN_TIMEOUT    | 5s                       | The graceful shutdown timeout in seconds (`time.Duration` format)                                                  |
| HEALTHCHECK_INTERVAL         | 30s                      | Time between self-healthchecks (`time.Duration` format)                                                            |
| HEALTHCHECK_CRITICAL_TIMEOUT | 90s                      | Time to wait until an unhealthy dependent propagates its state to make this app unhealthy (`time.Duration` format) |
| OTEL_EXPORTER_OTLP_ENDPOINT  | localhost:4317           | Endpoint for OpenTelemetry service                                                                                 |
| OTEL_SERVICE_NAME            | dis-search-upstream-stub | Label of service for OpenTelemetry service                                                                         |
| OTEL_BATCH_TIMEOUT           | 5s                       | Timeout for OpenTelemetry                                                                                          |
| OTEL_ENABLED                 | false                    | Feature flag to enable OpenTelemetry                                                                               |


### Note:
The `type` parameter in the resource API is optional for the upstream service and is intended for internal team use. It allows specifying the resource type as either "old" - `content-updated` or "new" - `search-content-updated` By default, it returns "new" if not specified.

## Contributing

See [CONTRIBUTING](CONTRIBUTING.md) for details.

## License

Copyright © 2025, Office for National Statistics (https://www.ons.gov.uk)

Released under MIT license, see [LICENSE](LICENSE.md) for details.
