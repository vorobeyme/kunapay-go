# kunapay-go
KunaPay API client for Go.

[![GoDoc](https://godoc.org/github.com/vorobeyme/kunapay-go?status.svg)](https://godoc.org/github.com/vorobeyme/kunapay-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/vorobeyme/kunapay-go)](https://goreportcard.com/report/github.com/vorobeyme/kunapay-go)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](./LICENSE)

The public API documentation is available at [https://docs-pay.kuna.io](https://docs-pay.kuna.io/reference).

## Installation
```bash
go get github.com/vorobeyme/kunapay-go
```

## Usage

```go
import "github.com/vorobeyme/kunapay-go"
```

Create a new KunaPay client, then use the exposed services to access different parts of the KunaPay API.

```go
package main

import (
    "log"

    "github.com/vorobeyme/kunapay-go"
)

func main() {

}
```

## Examples

To find code examples that demonstrate how to call the KunaPay API client, see the [examples](/examples/) folder.


## License

[MIT License](./LICENSE)