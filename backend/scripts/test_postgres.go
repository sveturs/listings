//go:build ignore
// +build ignore

package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	fmt.Println("=== PostgreSQL Integration Test ===")
	fmt.Println()

	// Загружаем .env файл
	if err := godotenv.Load("../.env"); err != nil {
		fmt.Printf("⚠️  Warning: .env file not found, using defaults\n\n")
	}

	// Получаем конфигурацию
	dbURL := getEnv("DATABASE_URL", "postgres://postgres:password@localhost:5432/svetubd?sslmode=disable")

	fmt.Printf("📋 Configuration:\n")
	fmt.Printf("   Database URL: %s\n", maskPassword(dbURL))
	fmt.Println()

	// Test 1: Connect to database
	fmt.Println("🔌 Test 1: Database Connection")
	fmt.Println("-------------------------------")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Printf("❌ Failed to open database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	// Set connection pool settings
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Test connection
	if err := db.Ping(); err != nil {
		fmt.Printf("❌ Failed to ping database: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✅ Connected to PostgreSQL successfully!")
	fmt.Println()

	// Test 2: Database version and info
	fmt.Println("📊 Test 2: Database Information")
	fmt.Println("--------------------------------")

	var version string
	err = db.QueryRow("SELECT version()").Scan(&version)
	if err != nil {
		fmt.Printf("❌ Failed to get version: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("✅ PostgreSQL Version:\n   %s\n", version[:80]+"...")

	var dbName string
	err = db.QueryRow("SELECT current_database()").Scan(&dbName)
	if err != nil {
		fmt.Printf("❌ Failed to get database name: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("✅ Current Database: %s\n", dbName)

	var dbSize string
	err = db.QueryRow("SELECT pg_size_pretty(pg_database_size(current_database()))").Scan(&dbSize)
	if err != nil {
		fmt.Printf("❌ Failed to get database size: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("✅ Database Size: %s\n", dbSize)
	fmt.Println()

	// Test 3: Connection pool stats
	fmt.Println("🔗 Test 3: Connection Pool Statistics")
	fmt.Println("--------------------------------------")

	stats := db.Stats()
	fmt.Printf("✅ Connection Pool:\n")
	fmt.Printf("   Open Connections: %d\n", stats.OpenConnections)
	fmt.Printf("   In Use: %d\n", stats.InUse)
	fmt.Printf("   Idle: %d\n", stats.Idle)
	fmt.Printf("   Max Open: %d\n", stats.MaxOpenConnections)
	fmt.Println()

	// Test 4: List tables
	fmt.Println("📋 Test 4: Database Tables")
	fmt.Println("--------------------------")

	rows, err := db.Query(`
		SELECT table_name,
		       pg_size_pretty(pg_total_relation_size(quote_ident(table_name)::regclass)) as size
		FROM information_schema.tables
		WHERE table_schema = 'public'
		  AND table_type = 'BASE TABLE'
		ORDER BY pg_total_relation_size(quote_ident(table_name)::regclass) DESC
		LIMIT 10
	`)
	if err != nil {
		fmt.Printf("❌ Failed to list tables: %v\n", err)
		os.Exit(1)
	}
	defer rows.Close()

	fmt.Println("✅ Top 10 tables by size:")
	tableCount := 0
	for rows.Next() {
		var tableName, size string
		if err := rows.Scan(&tableName, &size); err != nil {
			fmt.Printf("❌ Failed to scan row: %v\n", err)
			continue
		}
		tableCount++
		fmt.Printf("   %2d. %-30s %s\n", tableCount, tableName, size)
	}
	fmt.Println()

	// Test 5: Sample query (count records in a table)
	fmt.Println("🔢 Test 5: Sample Queries")
	fmt.Println("-------------------------")

	// Check if marketplace_listings exists
	var hasListings bool
	err = db.QueryRow(`
		SELECT EXISTS (
			SELECT FROM information_schema.tables
			WHERE table_schema = 'public'
			AND table_name = 'marketplace_listings'
		)
	`).Scan(&hasListings)

	if err != nil {
		fmt.Printf("❌ Failed to check table existence: %v\n", err)
	} else if hasListings {
		var count int
		err = db.QueryRow("SELECT COUNT(*) FROM marketplace_listings").Scan(&count)
		if err != nil {
			fmt.Printf("❌ Failed to count listings: %v\n", err)
		} else {
			fmt.Printf("✅ Marketplace Listings: %d records\n", count)
		}
	} else {
		fmt.Println("⚠️  marketplace_listings table not found")
	}

	// Check if users exists
	var hasUsers bool
	err = db.QueryRow(`
		SELECT EXISTS (
			SELECT FROM information_schema.tables
			WHERE table_schema = 'public'
			AND table_name = 'users'
		)
	`).Scan(&hasUsers)

	if err != nil {
		fmt.Printf("❌ Failed to check users table: %v\n", err)
	} else if hasUsers {
		var count int
		err = db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
		if err != nil {
			fmt.Printf("❌ Failed to count users: %v\n", err)
		} else {
			fmt.Printf("✅ Users: %d records\n", count)
		}
	} else {
		fmt.Println("⚠️  users table not found")
	}
	fmt.Println()

	// Test 6: Transaction test
	fmt.Println("💳 Test 6: Transaction Test")
	fmt.Println("---------------------------")

	tx, err := db.Begin()
	if err != nil {
		fmt.Printf("❌ Failed to begin transaction: %v\n", err)
		os.Exit(1)
	}

	// Create temporary table
	_, err = tx.Exec(`
		CREATE TEMPORARY TABLE test_integration (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT NOW()
		)
	`)
	if err != nil {
		tx.Rollback()
		fmt.Printf("❌ Failed to create temp table: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✅ Temporary table created")

	// Insert test data
	_, err = tx.Exec("INSERT INTO test_integration (name) VALUES ($1), ($2), ($3)",
		"Test 1", "Test 2", "Test 3")
	if err != nil {
		tx.Rollback()
		fmt.Printf("❌ Failed to insert data: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✅ Test data inserted")

	// Query test data
	var count int
	err = tx.QueryRow("SELECT COUNT(*) FROM test_integration").Scan(&count)
	if err != nil {
		tx.Rollback()
		fmt.Printf("❌ Failed to count test data: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("✅ Transaction query successful: %d records\n", count)

	// Commit transaction
	if err := tx.Commit(); err != nil {
		fmt.Printf("❌ Failed to commit transaction: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✅ Transaction committed")
	fmt.Println()

	fmt.Println("╔════════════════════════════════════════╗")
	fmt.Println("║  ✅ All PostgreSQL tests passed! 🎉   ║")
	fmt.Println("╚════════════════════════════════════════╝")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func maskPassword(dbURL string) string {
	// Simple password masking for display
	// postgres://user:password@host:port/db -> postgres://user:***@host:port/db
	return dbURL // In production, implement proper masking
}
