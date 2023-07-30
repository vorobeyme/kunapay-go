# kunapay-go

[![GoDoc](https://godoc.org/github.com/vorobeyme/kunapay-go?status.svg)](https://godoc.org/github.com/vorobeyme/kunapay-go)
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

For example, to list all transactions:
```go
package main

import (
    "context"
    "fmt"

    "github.com/vorobeyme/kunapay-go"
)

func main() {
    client := kunapay.New("public_key", "private_key", nil)
    transactions, _, err := client.Transaction.List(context.Background(), &kunapay.TransactionListOpts{})

    fmt.Println(transactions, err)
}
```

## Examples

To find code examples that demonstrate how to call the KunaPay API client, see the [examples](/examples/) folder.


## License

[MIT License](./LICENSE)