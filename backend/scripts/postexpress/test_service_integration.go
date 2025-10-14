// Test script for Post Express B2B Manifest API integration using updated service
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"backend/internal/proj/postexpress"
	peservice "backend/internal/proj/postexpress/service"
	"backend/pkg/logger"
)

func main() {
	fmt.Println("=== Post Express B2B Manifest Integration Test ===")
	fmt.Println()

	// Создаём WSP клиента с правильными настройками
	config := &peservice.WSPConfig{
		Endpoint:        "http://212.62.32.201/WspWebApi/transakcija",
		Username:        "b2b@svetu.rs",
		Password:        "Sv5et@U!",
		Language:        "sr-Latn-RS",
		DeviceType:      "2", // ВАЖНО: строка "2" для веб-приложения
		Timeout:         30 * time.Second,
		MaxRetries:      2,
		RetryDelay:      2 * time.Second,
		TestMode:        true,
		DeviceName:      "SVETU-Backend",
		ApplicationName: "SVETU-Platform",
		Version:         "0.2.1",
		PartnerID:       10109, // Partner ID для svetu.rs
	}

	log := logger.New()
	client := peservice.NewWSPClient(config, *log)

	ctx := context.Background()

	// Тест 1: Standard shipment
	fmt.Println("📦 Test 1: Creating Standard Shipment via new B2B API")
	fmt.Println(strings.Repeat("-", 60))

	standardShipment := &peservice.WSPShipmentRequest{
		SenderName:          "SVETU Platforma d.o.o.",
		SenderAddress:       "Bulevar kralja Aleksandra 73",
		SenderCity:          "Beograd",
		SenderPostalCode:    "11000",
		SenderPhone:         "+381112345678",
		RecipientName:       "Marko Marković",
		RecipientAddress:    "Takovska 2",
		RecipientCity:       "Beograd",
		RecipientPostalCode: "11000",
		RecipientPhone:      "+381691234567",
		Weight:              0.75,  // kg
		CODAmount:           0,     // no COD
		InsuranceAmount:     0,     // no insurance
		ServiceType:         "PE_Danas_za_sutra_12",
		Content:             "Test sadržaj - knjige",
		Note:                "Test integration - Standard",
	}

	resp1, err := client.CreateShipmentViaManifest(ctx, standardShipment)
	if err != nil {
		fmt.Printf("❌ ERROR: %v\n", err)
		os.Exit(1)
	}

	printManifestResponse(resp1, "Standard Shipment")

	// Тест 2: COD shipment
	fmt.Println()
	fmt.Println("💰 Test 2: Creating COD Shipment via new B2B API")
	fmt.Println(strings.Repeat("-", 60))

	codShipment := &peservice.WSPShipmentRequest{
		SenderName:          "SVETU Platforma d.o.o.",
		SenderAddress:       "Bulevar kralja Aleksandra 73",
		SenderCity:          "Beograd",
		SenderPostalCode:    "11000",
		SenderPhone:         "+381112345678",
		RecipientName:       "Ana Anić",
		RecipientAddress:    "Kneza Miloša 10",
		RecipientCity:       "Beograd",
		RecipientPostalCode: "11000",
		RecipientPhone:      "+381691234568",
		Weight:              0.5,     // kg
		CODAmount:           5000.00, // 5000 RSD COD
		InsuranceAmount:     5000.00, // 5000 RSD insurance
		ServiceType:         "PE_Danas_za_sutra_12",
		Content:             "Test sadržaj COD - elektronika",
		Note:                "Test integration - COD",
	}

	resp2, err := client.CreateShipmentViaManifest(ctx, codShipment)
	if err != nil {
		fmt.Printf("❌ ERROR: %v\n", err)
		os.Exit(1)
	}

	printManifestResponse(resp2, "COD Shipment")

	fmt.Println()
	fmt.Println("✅ All tests completed!")
}

func printManifestResponse(resp *postexpress.ManifestResponse, testName string) {
	fmt.Printf("\n🎯 %s Result:\n", testName)
	fmt.Println(strings.Repeat("=", 60))

	// Pretty print JSON
	jsonBytes, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		fmt.Printf("❌ Failed to marshal response: %v\n", err)
		return
	}
	fmt.Println(string(jsonBytes))
	fmt.Println()

	// Анализ результата
	if resp.Rezultat == 0 {
		fmt.Println("✅ Manifest created successfully!")
		fmt.Printf("   Manifest ID: %d\n", resp.IDManifesta)
		fmt.Printf("   External ID: %s\n", resp.ExtIDManifest)

		if len(resp.Porudzbine) > 0 {
			fmt.Printf("   Orders: %d\n", len(resp.Porudzbine))

			for i, order := range resp.Porudzbine {
				fmt.Printf("\n   Order #%d: %s\n", i+1, order.BrojPorudzbine)
				fmt.Printf("   Shipments: %d\n", len(order.Posiljke))

				for j, shipment := range order.Posiljke {
					fmt.Printf("\n     Shipment #%d:\n", j+1)
					fmt.Printf("       Broj Posiljke: %s\n", shipment.BrojPosiljke)
					fmt.Printf("       ID Posiljke: %d\n", shipment.IDPosiljke)
					fmt.Printf("       Tracking: %s\n", shipment.TrackingNumber)
					fmt.Printf("       Status: %s\n", shipment.Status)
					if shipment.Rezultat == 0 {
						fmt.Printf("       ✅ Success\n")
					} else {
						fmt.Printf("       ⚠️  Warning/Error: %s\n", shipment.Poruka)
					}
				}
			}
		}
	} else {
		fmt.Printf("❌ Manifest creation failed!\n")
		fmt.Printf("   Result code: %d\n", resp.Rezultat)
		fmt.Printf("   Error: %s\n", resp.Poruka)
	}

	// Validation errors
	if len(resp.GreskeValidaci) > 0 {
		fmt.Printf("\n⚠️  Validation warnings (%d):\n", len(resp.GreskeValidaci))
		for i, err := range resp.GreskeValidaci {
			fmt.Printf("   %d. Field: %s, Value: %s\n", i+1, err.Polje, err.Vrednost)
			fmt.Printf("      Message: %s\n", err.Poruka)
		}
	}
}

// Simple logger removed - using backend/pkg/logger instead
