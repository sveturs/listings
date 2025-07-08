package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

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

	// Check search_logs table
	fmt.Println("=== Checking search_logs table ===")

	var count int
	err = pool.QueryRow(ctx, "SELECT COUNT(*) FROM search_logs").Scan(&count)
	if err != nil {
		log.Fatalf("Failed to count search_logs: %v", err)
	}
	fmt.Printf("Total records in search_logs: %d\n", count)

	// Get last 10 searches
	if count > 0 {
		fmt.Println("\n=== Last 10 searches ===")
		rows, err := pool.Query(ctx, `
			SELECT query_text, results_count, response_time_ms, device_type, created_at 
			FROM search_logs 
			ORDER BY created_at DESC 
			LIMIT 10
		`)
		if err != nil {
			log.Printf("Failed to get last searches: %v", err)
		} else {
			defer rows.Close()
			for rows.Next() {
				var query string
				var results int
				var responseTime int64
				var deviceType string
				var createdAt time.Time

				if err := rows.Scan(&query, &results, &responseTime, &deviceType, &createdAt); err != nil {
					log.Printf("Failed to scan row: %v", err)
					continue
				}
				fmt.Printf("- %s: '%s' (%d results, %dms, %s)\n",
					createdAt.Format("2006-01-02 15:04:05"),
					query, results, responseTime, deviceType)
			}
		}
	}

	// Check search_analytics table
	fmt.Println("\n=== Checking search_analytics table ===")
	err = pool.QueryRow(ctx, "SELECT COUNT(*) FROM search_analytics").Scan(&count)
	if err != nil {
		log.Printf("Failed to count search_analytics: %v (table might not exist yet)", err)
	} else {
		fmt.Printf("Total records in search_analytics: %d\n", count)
	}

	// Check search_queries table (old name)
	fmt.Println("\n=== Checking search_queries table ===")
	err = pool.QueryRow(ctx, "SELECT COUNT(*) FROM search_queries").Scan(&count)
	if err != nil {
		log.Printf("Failed to count search_queries: %v (table might not exist)", err)
	} else {
		fmt.Printf("Total records in search_queries: %d\n", count)
	}
}
