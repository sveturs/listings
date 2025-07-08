package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(".env"); err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	// Get database URL
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	// Connect to database
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	// Check all search-related tables
	fmt.Println("=== Search-related tables in database ===")
	rows, err := pool.Query(ctx, `
		SELECT table_name 
		FROM information_schema.tables 
		WHERE table_schema = 'public' 
		AND (table_name LIKE '%search%' OR table_name LIKE '%analytic%')
		ORDER BY table_name
	`)
	if err != nil {
		log.Fatalf("Failed to query tables: %v", err)
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			log.Printf("Failed to scan row: %v", err)
			continue
		}
		tables = append(tables, tableName)
		fmt.Printf("- %s\n", tableName)
	}

	// Check counts for each table
	fmt.Println("\n=== Table record counts ===")
	for _, table := range tables {
		var count int
		err := pool.QueryRow(ctx, fmt.Sprintf("SELECT COUNT(*) FROM %s", table)).Scan(&count)
		if err != nil {
			fmt.Printf("- %s: Error counting: %v\n", table, err)
		} else {
			fmt.Printf("- %s: %d records\n", table, count)
		}
	}

	// Check search_statistics table structure
	fmt.Println("\n=== search_statistics table structure ===")
	rows, err = pool.Query(ctx, `
		SELECT column_name, data_type 
		FROM information_schema.columns 
		WHERE table_name = 'search_statistics' 
		ORDER BY ordinal_position
	`)
	if err != nil {
		fmt.Println("Table search_statistics not found or error:", err)
	} else {
		defer rows.Close()
		for rows.Next() {
			var colName, dataType string
			rows.Scan(&colName, &dataType)
			fmt.Printf("- %s: %s\n", colName, dataType)
		}
	}
}
