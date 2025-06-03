# Load Test Script

This module contains the `Makefile` and a `mass-produce` script to load test our search pipeline and examine its
performance including kafka.

## Default Variables

The following default variables are used in the `Makefile`:

## Environment Variables

| Variable               | Description                                                            | Default Value                |
|------------------------|------------------------------------------------------------------------|------------------------------|
| `APP`                  | Name of the application.                                               | `dis-search-upstream-stub`   |
| `ENV`                  | The environment to deploy to.                                          | `sandbox`                    |
| `SUBNET`               | The subnet to deploy to.                                               | `publishing`                 |
| `TIMEOUT`              | Timeout value for operations.                                          | `30s`                        |
| `host_num`             | Host number for remote deployment.                                     | `publishing 3`               |
| `host_bin`             | The binary directory on the remote host (constructed as `bin-$(APP)`). | `bin-$(APP)`                 |
| `GOOS`                 | Target operating system.                                               | `linux`                      |
| `GOARCH`               | Target architecture.                                                   | `amd64`                      |
| `BUILD`                | The build directory.                                                   | `build`                      |
| `BUILD_ARCH`           | Architecture-specific build directory (`$(BUILD)/$(GOOS)-$(GOARCH)`).  | `$(BUILD)/$(GOOS)-$(GOARCH)` |
| `DP_CONFIGS`           | Path to `dp-configs`.                                                  | `../../../dp-configs`        |
| `SECRETS_APP`          | Application name for secrets fetching.                                 | `dis-search-upstream-stub`   |
| `MESSAGE_COUNT_LEGACY` | Number of legacy messages.                                             | `6000`                       |
| `MESSAGE_COUNT_NEW`    | Number of new messages.                                                | `100`                        |

## How to Run

### To Build, Deploy, and Clean:

Simply run:

```bash
make
```

This will clean, build, deploy, and clean up.

To build the application:

```bash
make build
```

To deploy the application:

```bash
make deploy
```

To clean build and deployment artifacts:

```bash
make clean
```
