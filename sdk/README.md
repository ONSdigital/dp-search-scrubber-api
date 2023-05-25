dp-search-scrubber-api SDK
======================

## Overview

The scrubber API contains a Go client for interacting with the API. The client contains a methods for each API endpoint
so that any Go application wanting to interact with the search api can do so. Please refer to the [swagger specification](../swagger.yaml)
as the source of truth of how each endpoint works.

## Example use of the API SDK

Initialise new Search API client

```go
package main

import (
	"context"
	"github.com/ONSdigital/dp-search-scrubber-api/sdk"
)

func main() {
    ...
	searchAPIClient := sdk.NewClient("http://localhost:28700")
    ...
}
```

### Get Search Results

Use the GetSearch method to send a request to find search results based on query parameters.

```go
...
    // Set query parameters - no limit to which keys and values you set - please refer to swagger spec for list of available parameters
    query := url.Values{}
    query.Add("q", "E00000013,01220")

    resp, err := searchAPIClient.GetSearch(ctx, sdk.Options{sdk.Query: query})
    if err != nil {
        // handle error
    }
...
```

### Handling errors

The error returned from the method contains status code that can be accessed via `Status()` method and similar to extracting the error message using `Error()` method; see snippet below:

```go
...
    _, err := searchAPIClient.GetSearch(ctx, Options{})
    if err != nil {
        // Retrieve status code from error
        statusCode := err.Status()
        // Retrieve error message from error
        errorMessage := err.Error()

        // log message, below uses "github.com/ONSdigital/log.go/v2/log" package
        log.Error(ctx, "failed to retrieve search results", err, log.Data{"code": statusCode})

        return err
    }
...
```