// Example of getting transactions
//
// It's runnable with the following command:
// export KUNAPAY_PUBLIC_KEY=your_public_key
// export KUNAPAY_PRIVATE_KEY=your_private_key
// go run examples/transaction/main.go
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

	ctx := context.Background()
	client, err := kunapay.New(pubKey, privKey)
	if err != nil {
		log.Fatal(err)
	}

	for {
		var method string
		fmt.Print("Enter method [list, get] or `q` to quit: ")
		fmt.Scanf("%s", &method)
		if strings.ToLower(method) == "q" {
			return
		}

		switch method {
		case "list":
			fmt.Println("Getting transactions...")
			t, _, err := client.Transaction.List(ctx, nil)
			if err != nil {
				log.Fatal(err)
			}
			for _, tx := range t {
				fmt.Printf("%s: %s\n", tx.ID, tx.Amount)
			}
		case "get":
			var ID string
			fmt.Print("Enter transaction ID: ")
			fmt.Scanf("%s", &ID)
			fmt.Printf("Getting transaction %s...\n", ID)
			t, _, err := client.Transaction.Get(ctx, ID)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("%+v\n", t)
		default:
			fmt.Println("Unknown method")
		}
	}
}
