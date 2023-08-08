// Example of using kunapay-go library to create, list, get invoices
// and get invoice currencies.
//
// It's runnable with the following command:
// export KUNAPAY_PUBLIC_KEY=your_public_key
// export KUNAPAY_PRIVATE_KEY=your_private_key
// go run examples/invoice/main.go
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
	client := kunapay.New(pubKey, privKey, nil)

	for {
		var method string
		fmt.Print("Enter method [create, list, get, currencies] or `q` to quit: ")
		fmt.Scanf("%s", &method)
		if strings.ToLower(method) == "q" {
			return
		}

		switch method {
		case "create":
			fmt.Println("Creating invoice...")
			i, _, err := client.Invoice.Create(ctx, &kunapay.CreateInvoiceRequest{
				Amount:             "1.00",
				Asset:              "USDT",
				ExternalOrderID:    "1234567890",
				ProductDescription: "Test invoice",
				ProductCategory:    "Test",
				CallbackUrl:        "https://example.com/callback",
			})
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("%+v\n", i)

		case "list":
			fmt.Println("Getting invoices...")
			i, _, err := client.Invoice.List(ctx, nil)
			if err != nil {
				log.Fatal(err)
			}
			for _, invoice := range i {
				fmt.Printf(`ID: %s\n AddresID: %s\n ExternalOrderID: %s\n PaymentAmount: %s\n 
					InvoiceAmount: %s\n InvoiceAssetCode: %s\n PaymentAssetCode: %s\n
					ExpireAt: %s\n CompletedAt: %s\n CreatedAt: %s\n`,
					invoice.ID, invoice.AddressID, invoice.ExternalOrderID, invoice.PaymentAmount,
					invoice.InvoiceAmount, invoice.InvoiceAssetCode, invoice.PaymentAssetCode,
					invoice.ExpireAt, invoice.CompletedAt, invoice.CreatedAt,
				)
			}

		case "get":
			var ID string
			fmt.Print("Enter invoice ID: ")
			fmt.Scanf("%s", &ID)
			fmt.Printf("Getting invoice %s...\n", ID)
			t, _, err := client.Invoice.Get(ctx, ID)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("%+v\n", t)

		case "currencies":
			fmt.Println("Getting currencies...")
			c, _, err := client.Invoice.GetCurrencies(ctx, nil)
			if err != nil {
				log.Fatal(err)
			}
			for _, currency := range c {
				fmt.Printf("%s: %s\n", currency.Code, currency.Name)
			}

		default:
			fmt.Println("Unknown method")
		}
	}
}
