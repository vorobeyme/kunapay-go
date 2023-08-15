# kunapay-go

[![Go Reference](https://pkg.go.dev/badge/github.com/vorobeyme/kunapay-go.svg)](https://pkg.go.dev/github.com/vorobeyme/kunapay-go)
[![Go](https://github.com/vorobeyme/kunapay-go/actions/workflows/go.yml/badge.svg)](https://github.com/vorobeyme/kunapay-go/actions/workflows/go.yml)
[![codecov](https://codecov.io/gh/vorobeyme/kunapay-go/branch/main/graph/badge.svg?token=HV37K62JA3)](https://codecov.io/gh/vorobeyme/kunapay-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/vorobeyme/kunapay-go)](https://goreportcard.com/report/github.com/vorobeyme/kunapay-go)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](./LICENSE)

KunaPay API client for Go is the HTTP library implementation with [net/http](https://pkg.go.dev/net/http).

The public API documentation is available at [https://docs-pay.kuna.io](https://docs-pay.kuna.io/reference).

## Installation
```bash
go get github.com/vorobeyme/kunapay-go
```

## Usage

```go
import "github.com/vorobeyme/kunapay-go"
```

To begin, create a new KunaPay client, then use the available services to interact with various sections of the KunaPay API.

For example, to get the balance of the assets:
```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"

    "github.com/vorobeyme/kunapay-go"
)

func main() {
    // Create a new API client using signature authentication with your public and private keys
    client, err := kunapay.New(os.Getenv("KUNAPAY_PUBLIC_KEY"), os.Getenv("KUNAPAY_PRIVATE_KEY"))
    if err != nil {
        log.Fatal(err)
    }

    // Alternatively, you can use an API key
    // client, err := kunapay.NewWithAPIKey(os.Getenv("KUNAPAY_API_KEY"))

    // Create a new API object with a custom HTTP client and/or user agent
    // client, err := kunapay.NewWithAPIKey(
    //    os.Getenv("KUNAPAY_API_KEY"),
    //    kunapay.WithHTTPClient(&http.Client{}),
    //    kunapay.SetUserAgent("MyApp/1.0.0"),
    // )

    // API calls require a context
    ctx := context.Background()

    // Get the balance of the specified assets
    balance, _, err := client.Assets.GetBalance(ctx, "btc", "uah")
    if err != nil {
        log.Fatal(err)
    }

    // Print the balance
    fmt.Println(balance)
}
```

## Examples

To find code examples that demonstrate how to call the KunaPay API client, see the [examples](/examples/) folder.


## License

[MIT License](./LICENSE)
