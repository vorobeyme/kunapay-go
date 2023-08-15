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
	ctx := context.Background()
	client, err := kunapay.New(os.Getenv("KUNAPAY_PUBLIC_KEY"), os.Getenv("KUNAPAY_PRIVATE_KEY"))
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
				var fields []string
				for _, field := range method.Fields {
					fields = append(fields, fmt.Sprintf("\n\tName: %s \n\tLabel: %s \n\tDescription: %s \n\tPosition: %d \n\tType: %s"+
						"\n\tIsRequired: %t \n\tIsMasked: %t \n\tIsResultField: %t\n",
						field.Name, field.Label, field.Description, field.Position, field.Type, field.IsRequired, field.IsMasked, field.IsResultField))
				}
				fmt.Printf("Code: %s \nAsset: %s \nNetwork: %s \nPosition: %d, \nName: %s \nIcon: %s \nDescription: %s \nCustomTitle: %s \nFields: %+v\n\n",
					method.Code, method.Asset, method.Network, method.Position, method.Name, method.Icon, method.Description, method.CustomTitle, fields)
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
