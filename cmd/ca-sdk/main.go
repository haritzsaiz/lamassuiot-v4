package main

import (
	"context"

	"github.com/lamassuiot/lamassuiot/v4/pkg/ca"
	"github.com/lamassuiot/lamassuiot/v4/pkg/models"
)

func main() {
	// Initialize the CASdkService
	caSdkService := ca.NewCASdkService()

	// Example usage of CreateCA
	err := caSdkService.CreateCA(context.Background(), ca.CreateCAInput{Name: "MyCA"})
	if err != nil {
		panic(err)
	}

	// Example usage of GetCAs
	cas, err := caSdkService.GetCAs(context.Background(), ca.GetCAsInput{
		ApplyFunc: func(ca models.CACertificate) {
			// Process each CA certificate
			println("CA Name:", ca.ID)
		},
	})
	if err != nil {
		panic(err)
	}

	println("Retrieved CAs:", cas)
}
