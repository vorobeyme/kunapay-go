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
	"time"

	"github.com/vorobeyme/kunapay-go"
)

func main() {
	ctx := context.Background()
	client, err := kunapay.New(os.Getenv("KUNAPAY_PUBLIC_KEY"), os.Getenv("KUNAPAY_PRIVATE_KEY"))
	if err != nil {
		log.Fatal(err)
	}

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

			var amount string
			fmt.Print("Enter amount: ")
			fmt.Scanf("%s", &amount)

			var currency string
			fmt.Print("Enter currency (UAH, EUR, USDT): ")
			fmt.Scanf("%s", &currency)

			i, _, err := client.Invoice.Create(ctx, &kunapay.CreateInvoiceRequest{
				Amount:             amount,
				Asset:              currency,
				ExternalOrderID:    fmt.Sprintf("test-%d", time.Now().Unix()),
				ProductDescription: "Test invoice",
				ProductCategory:    "Test",
				CallbackURL:        "https://example.com/callback",
			})
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Invoice %s created. Payment link: %s\n", i.ID, i.PaymentLink)

		case "list":
			fmt.Println("Getting invoices...")
			i, _, err := client.Invoice.List(ctx, nil)
			if err != nil {
				log.Fatal(err)
			}
			for _, invoice := range i {
				fmt.Printf("ID: %s\nAddresID: %s\nExternalOrderID: %s\nPaymentAmount: %s\nInvoiceAmount: %s\nInvoiceAssetCode: %s"+
					"\nPaymentAssetCode: %s\nExpireAt: %s\nCompletedAt: %s\nCreatedAt: %s\n\n",
					invoice.ID, invoice.AddressID, invoice.ExternalOrderID, invoice.PaymentAmount,
					invoice.InvoiceAmount, invoice.InvoiceAssetCode, invoice.PaymentAssetCode,
					invoice.ExpireAt, invoice.CompletedAt, invoice.CreatedAt,
				)
			}

		case "get":
			var ID string
			fmt.Print("Enter invoice ID: ")
			fmt.Scanf("%s", &ID)
			fmt.Printf("Getting invoice %s...\n\n", ID)
			t, _, err := client.Invoice.Get(ctx, ID)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("ID: %s\nStatus: %s\nExternalOrderID: %s\nAddressID: %s\nCreatorID: %s\nInvoiceAmount: %s\nPaymentAmount: %s"+
				"\nInvoiceAssetCode: %s\nPaymentAssetCode: %s\nProductCategory: %s\nProductDescription: %s"+
				"\nIsCreatedByAPI: %t\nExpireAt: %s\nCompletedAt: %s\nCreatedAt: %s\nUpdatedAt: %s\nTransactions: %v\n",
				t.ID, t.Status, t.ExternalOrderID, t.AddressID, t.CreatorID, t.InvoiceAmount,
				t.PaymentAmount, t.InvoiceAssetCode, t.PaymentAssetCode, t.ProductCategory,
				t.ProductDescription, t.IsCreatedByAPI, t.ExpireAt, t.CompletedAt, t.CreatedAt, t.UpdateAt, t.Transactions,
			)

		case "currencies":
			fmt.Println("Getting currencies...")
			c, _, err := client.Invoice.GetCurrencies(ctx, nil)
			if err != nil {
				log.Fatal(err)
			}
			for _, currency := range c {
				fmt.Printf("Code: %s\nName: %s\nDescription: %s\nPosition: %d\nPrecision: %d\nType: %s\nSVG icon: %s\nPNG icon: %s\n\n",
					currency.Code, currency.Name, currency.Description, currency.Position, currency.Precision, currency.Type,
					currency.Icons.SVG, currency.Icons.PNG)
			}

		default:
			fmt.Println("Unknown method")
		}
	}
}
