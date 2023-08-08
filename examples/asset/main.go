// Example of getting asset balance
//
// It's runnable with the following command:
// export KUNAPAY_PUBLIC_KEY=bwumw6Q8nsvvwBSoIUWMHbaLEjGXKiToNtXXU2jQv6Q=
// export KUNAPAY_PRIVATE_KEY=DNYKfHBWxNbfAa2fTl4+d6rfY04Y+W9Dd3KqEyVM2Eo=
// go run examples/asset/main.go
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/vorobeyme/kunapay-go"
)

func main() {
	pubKey := os.Getenv("KUNAPAY_PUBLIC_KEY")
	if pubKey == "" {
		log.Fatal("Public key is not set")
	}
	privKey := os.Getenv("KUNAPAY_PRIVATE_KEY")
	if privKey == "" {
		log.Fatal("Private key is not set")
	}

	client := kunapay.New(pubKey, privKey, nil)

	for {
		var method string
		fmt.Print("Enter method [balance] or `q` to quit: ")
		fmt.Scanf("%s", &method)
		if strings.ToLower(method) == "q" {
			return
		}

		switch method {
		case "balance":
			fmt.Println("Getting balance...")
			b, _, err := client.Asset.GetBalance(context.Background(), nil)
			if err != nil {
				log.Fatal(err)
			}
			for _, asset := range b {
				fmt.Printf("%s: %s\n", asset.Code, asset.Balance)
			}
		default:
			fmt.Println("Unknown method")
		}
	}
}
