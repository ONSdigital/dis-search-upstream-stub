# dis-search-upstream-stub SDK

## Overview

The `dis-search-upstream-stub` contains a Go client for interacting a generic upstream API for search purposes. The client contains a methods for each API endpoint
so that any Go application wanting to interact with an upstream API can do so. Please refer to the [swagger specification](../swagger.yaml)
as the source of truth of how each endpoint works.

## Example use of the API SDK

Initialise new Search API client

```go
package main

import (
    "context"
    "github.com/ONSdigital/dis-search-upstream-stub/sdk"
)

func main() {
    ...
    upstreamAPIClient := sdk.NewClient("http://localhost:29600", "/resources")
    ...
}
```

### Getting Resources

Use the GetResources method to send a request to get Resources. Authorisation header is needed if hitting private instance of applications.

```go
...
    // Set query parameters - no limit to which keys and values you set - please refer to swagger spec for list of available parameters
    query := url.Values{}
    query.Add("offset", 10)

    resp, err := upstreamAPIClient.GetResources(ctx, sdk.Options{sdk.Limit: query})
    if err != nil {
        // handle error
    }
...
```

### Handling errors

The error returned from the method contains status code that can be accessed via `Status()` method and similar to extracting the error message using `Error()` method; see snippet below:

```go
...
    _, err := upstreamAPIClient.GetResources(ctx, Options{})
    if err != nil {
        // Retrieve status code from error
        statusCode := err.Status()
        // Retrieve error message from error
        errorMessage := err.Error()

        // log message, below uses "github.com/ONSdigital/log.go/v2/log" package
        log.Error(ctx, "failed to retrieve resources", err, log.Data{"code": statusCode})

        return err
    }
...
```
