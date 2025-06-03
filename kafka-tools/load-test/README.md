# Load Test Script

This module contains the `Makefile` and a `mass-produce` script to load test our search pipeline and examine its performance in including kafka.

## Default Variables

The following default variables are used in the `Makefile`:

- `APP`: Name of the application. Default: `dis-search-upstream-stub`.
- `ENV`: The environment to deploy to. Default: `sandbox`.
- `SUBNET`: The subnet to deploy to. Default: `publishing`.
- `TIMEOUT`: Timeout value for operations. Default: `30s`.
- `host_num`: Host number for remote deployment. Default: `publishing 3`.
- `host_bin`: The binary directory on the remote host (constructed as `bin-$(APP)`).
- `GOOS`: Target operating system. Default: `linux`.
- `GOARCH`: Target architecture. Default: `amd64`.
- `BUILD`: The build directory. Default: `build`.
- `BUILD_ARCH`: Architecture-specific build directory (`$(BUILD)/$(GOOS)-$(GOARCH)`).
- `DP_CONFIGS`: Path to `dp-configs`. Default: `../../../dp-configs`.
- `SECRETS_APP`: Application name for secrets fetching. Default: `dis-search-upstream-stub`.
- `MESSAGE_COUNT_LEGACY`: Number of legacy messages. Default: `6000`.
- `MESSAGE_COUNT_NEW`: Number of new messages. Default: `100`.

## How to Run

### To Build, Deploy, and Clean:

Simply run:

```bash
make
```

This will clean, build, deploy, and clean up.

To Build the Application:
```bash
make build
```

To Deploy the Application:
```bash
make deploy
```

To Clean Build and Deployment Artifacts:
```bash
make clean
```


