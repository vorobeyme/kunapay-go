// Example of getting asset balance
//
// It's runnable with the following command:
// export KUNAPAY_PUBLIC_KEY=your_public_key
// export KUNAPAY_PRIVATE_KEY=your_private_key
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
	client, err := kunapay.New(os.Getenv("KUNAPAY_PUBLIC_KEY"), os.Getenv("KUNAPAY_PRIVATE_KEY"))
	if err != nil {
		log.Fatal(err)
	}

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
			b, _, err := client.Asset.GetBalance(context.Background())
			if err != nil {
				log.Fatal(err)
			}
			for _, asset := range b {
				fmt.Printf("Balance: %s\nLockBalance: %s\nCode: %s\nName: %s\nSVG icon: %s\nPNG icon: %s\n\n",
					asset.Balance, asset.LockBalance, asset.Code, asset.Name, asset.Icons.SVG, asset.Icons.PNG)
			}
		default:
			fmt.Println("Unknown method")
		}
	}
}
