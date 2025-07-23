package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run query_db.go <SQL query>")
	}

	query := os.Args[1]

	// Connection string from .env
	connStr := "postgres://postgres:password@localhost:5432/hostel_db?sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Failed to close database connection: %v", err)
		}
	}()

	rows, err := db.Query(query)
	if err != nil {
		log.Printf("Query failed: %v", err)
		return
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("Failed to close rows: %v", err)
		}
	}()

	cols, err := rows.Columns()
	if err != nil {
		log.Fatal("Failed to get columns:", err)
	}

	// Print column headers
	for i, col := range cols {
		if i > 0 {
			fmt.Print(" | ")
		}
		fmt.Print(col)
	}
	fmt.Println()
	fmt.Println("---")

	// Print rows
	values := make([]interface{}, len(cols))
	scanArgs := make([]interface{}, len(cols))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			log.Fatal("Failed to scan row:", err)
		}

		for i, value := range values {
			if i > 0 {
				fmt.Print(" | ")
			}
			switch v := value.(type) {
			case nil:
				fmt.Print("NULL")
			case []byte:
				fmt.Print(string(v))
			default:
				fmt.Print(v)
			}
		}
		fmt.Println()
	}

	// Check for iteration errors
	if err := rows.Err(); err != nil {
		log.Printf("Error iterating over rows: %v", err)
	}
}
