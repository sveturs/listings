// Simple test to verify Variant Detector works with Digital Vision data
package main

import (
	"fmt"

	"backend/internal/proj/storefronts/service"
)

func main() {
	// Create variant detector
	detector := service.NewVariantDetector()

	// Sample Digital Vision products (simplified from actual data)
	products := []*service.ProductVariant{
		{Name: "Tastatura Gembird KB-UM-104 crna USB", SKU: "TASTATURAGEMBIRDKB-UM-104-CRNA", Price: 800},
		{Name: "Tastatura Gembird KB-UM-104 bela USB", SKU: "TASTATURAGEMBIRDKB-UM-104-BELA", Price: 800},
		{Name: "Tastatura Gembird KB-UM-104 siva USB", SKU: "TASTATURAGEMBIRDKB-UM-104-SIVA", Price: 800},
		{Name: "Miš Genius DX-110 crni USB", SKU: "MISGENIUSDX-110-CRNI", Price: 500},
		{Name: "Miš Genius DX-110 beli USB", SKU: "MISGENIUSDX-110-BELI", Price: 500},
		{Name: "Monitor Samsung 24\" S24C310", SKU: "MONITORS24C310", Price: 15000},
	}

	fmt.Println("=== Testing Variant Detector on Digital Vision Data ===")

	// Test base name extraction
	fmt.Println("1. Base Name Extraction:")
	for _, p := range products {
		baseName := detector.ExtractBaseName(p.Name)
		fmt.Printf("   %s -> %s\n", p.Name, baseName)
	}

	// Test variant attribute extraction
	fmt.Println("\n2. Variant Attributes Extraction:")
	for _, p := range products {
		attrs := detector.ExtractVariantAttributes(p.Name)
		if len(attrs) > 0 {
			fmt.Printf("   %s:", p.Name)
			for key, val := range attrs {
				fmt.Printf(" %s=%s", key, val)
			}
			fmt.Println()
		}
	}

	// Test grouping
	fmt.Println("\n3. Product Grouping:")
	groups := detector.GroupProducts(products)
	for i, group := range groups {
		fmt.Printf("   Group %d: %s\n", i+1, group.BaseName)
		fmt.Printf("      Variants: %d (confidence: %.2f)\n", group.VariantCount, group.Confidence)
		fmt.Printf("      Is Grouped: %v\n", group.IsGrouped)
		if len(group.VariantAttributes) > 0 {
			fmt.Printf("      Variant Attributes: %v\n", group.VariantAttributes)
		}
		fmt.Println("      Products:")
		for j, variant := range group.Variants {
			attrs := ""
			for k, v := range variant.VariantAttributes {
				attrs += fmt.Sprintf(" %s=%s", k, v)
			}
			fmt.Printf("         %d. %s (SKU: %s)%s\n", j+1, variant.Name, variant.SKU, attrs)
		}
		fmt.Println()
	}

	// Validate groups
	fmt.Println("4. Validation:")
	for i, group := range groups {
		if group.VariantCount > 1 {
			warnings := detector.ValidateVariantGroup(group)
			if len(warnings) > 0 {
				fmt.Printf("   Group %d (%s) has warnings:\n", i+1, group.BaseName)
				for _, warning := range warnings {
					fmt.Printf("      - %s\n", warning)
				}
			} else {
				fmt.Printf("   Group %d (%s) ✓ Valid\n", i+1, group.BaseName)
			}
		}
	}

	fmt.Println("\n=== Test Completed ===")
}
