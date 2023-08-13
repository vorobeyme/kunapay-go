// Withdraw example
//
// It's runnable with the following command:
// export KUNAPAY_PUBLIC_KEY=your_public_key
// export KUNAPAY_PRIVATE_KEY=your_private_key
// go run examples/withdraw/main.go
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
		fmt.Print("Enter method [methods, create] or `q` to quit: ")
		fmt.Scanf("%s", &method)
		if strings.ToLower(method) == "q" {
			return
		}

		switch method {
		case "methods":
			fmt.Println("Enter asset code (3 letters, e.g. BTC): ")
			var assetCode string
			fmt.Scanf("%s", &assetCode)
			fmt.Printf("Getting withdraw methods for %s...\n", assetCode)
			w, _, err := client.Withdraw.GetMethods(ctx, assetCode)
			if err != nil {
				log.Fatal(err)
			}
			for _, method := range w {
				fmt.Printf("%+v\n", method)
			}

		case "create":
			fmt.Println("Creating withdraw...")
			w, _, err := client.Withdraw.Create(ctx, &kunapay.CreateWithdrawRequest{
				Asset:         "USDT",
				Amount:        "1",
				PaymentMethod: "USDT",
			})
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("%+v\n", w)

		default:
			fmt.Println("Unknown method")
		}
	}
}
