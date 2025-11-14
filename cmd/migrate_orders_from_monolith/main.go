package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

// MigrationStats tracks migration progress
type MigrationStats struct {
	ReservationsTotal    int
	ReservationsMigrated int
	ReservationsSkipped  int
	ReservationsFailed   int
	StartTime            time.Time
	EndTime              time.Time
}

// Config holds database connection strings
type Config struct {
	MonolithDSN    string
	MicroserviceDSN string
	DryRun         bool
	Verbose        bool
}

// InventoryReservation represents a reservation record
type InventoryReservation struct {
	ID        int64
	ProductID int64  // source: monolith uses product_id
	VariantID *int64
	OrderID   *int64 // nullable: order may not exist yet
	Quantity  int
	Status    string
	ExpiresAt time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

func main() {
	var cfg Config

	// Parse command-line flags
	flag.StringVar(&cfg.MonolithDSN, "monolith-dsn",
		"postgres://postgres:mX3g1XGhMRUZEX3l@localhost:5433/svetubd?sslmode=disable",
		"Monolith database connection string")
	flag.StringVar(&cfg.MicroserviceDSN, "microservice-dsn",
		"postgres://listings_user:listings_secret@localhost:35434/listings_dev_db?sslmode=disable",
		"Microservice database connection string")
	flag.BoolVar(&cfg.DryRun, "dry-run", true, "Dry run mode (no actual writes)")
	flag.BoolVar(&cfg.Verbose, "verbose", false, "Verbose logging")

	flag.Parse()

	if !cfg.DryRun {
		fmt.Println("⚠️  WARNING: Running in EXECUTE mode. Data will be written to microservice database.")
		fmt.Println("Press Ctrl+C within 5 seconds to cancel...")
		time.Sleep(5 * time.Second)
	}

	// Initialize logger
	logger := log.New(os.Stdout, "", log.LstdFlags)

	// Run migration
	stats, err := runMigration(cfg, logger)
	if err != nil {
		logger.Fatalf("❌ Migration failed: %v", err)
	}

	// Print summary
	printSummary(stats, cfg.DryRun, logger)
}

func runMigration(cfg Config, logger *log.Logger) (*MigrationStats, error) {
	stats := &MigrationStats{
		StartTime: time.Now(),
	}

	logger.Println("🚀 Starting Orders data migration from monolith to microservice")
	logger.Printf("   Mode: %s\n", mode(cfg.DryRun))
	logger.Println()

	// Connect to databases
	logger.Println("📡 Connecting to databases...")
	monolithDB, err := sql.Open("postgres", cfg.MonolithDSN)
	if err != nil {
		return stats, fmt.Errorf("failed to connect to monolith: %w", err)
	}
	defer monolithDB.Close()

	microserviceDB, err := sql.Open("postgres", cfg.MicroserviceDSN)
	if err != nil {
		return stats, fmt.Errorf("failed to connect to microservice: %w", err)
	}
	defer microserviceDB.Close()

	// Ping databases
	if err := monolithDB.Ping(); err != nil {
		return stats, fmt.Errorf("monolith database unreachable: %w", err)
	}
	if err := microserviceDB.Ping(); err != nil {
		return stats, fmt.Errorf("microservice database unreachable: %w", err)
	}

	logger.Println("✅ Connected to both databases")
	logger.Println()

	// Migrate inventory_reservations
	logger.Println("📦 Migrating inventory_reservations...")
	err = migrateInventoryReservations(monolithDB, microserviceDB, stats, cfg, logger)
	if err != nil {
		return stats, fmt.Errorf("failed to migrate inventory_reservations: %w", err)
	}

	stats.EndTime = time.Now()
	return stats, nil
}

func migrateInventoryReservations(monolithDB, microserviceDB *sql.DB, stats *MigrationStats, cfg Config, logger *log.Logger) error {
	// Query active reservations from monolith
	query := `
		SELECT
			id, product_id, variant_id, order_id, quantity,
			status::text, expires_at, created_at, updated_at
		FROM inventory_reservations
		WHERE expires_at > NOW() AND status = 'active'
		ORDER BY id
	`

	rows, err := monolithDB.Query(query)
	if err != nil {
		return fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	// Count total first
	countQuery := `SELECT COUNT(*) FROM inventory_reservations WHERE expires_at > NOW() AND status = 'active'`
	err = monolithDB.QueryRow(countQuery).Scan(&stats.ReservationsTotal)
	if err != nil {
		return fmt.Errorf("count query failed: %w", err)
	}

	logger.Printf("   Found %d active reservations to migrate\n", stats.ReservationsTotal)

	if stats.ReservationsTotal == 0 {
		logger.Println("   ℹ️  No active reservations found")
		return nil
	}

	// Start transaction for microservice writes
	var tx *sql.Tx
	if !cfg.DryRun {
		tx, err = microserviceDB.Begin()
		if err != nil {
			return fmt.Errorf("failed to start transaction: %w", err)
		}
		defer tx.Rollback()
	}

	// Process each reservation
	for rows.Next() {
		var r InventoryReservation

		err := rows.Scan(
			&r.ID, &r.ProductID, &r.VariantID, &r.OrderID, &r.Quantity,
			&r.Status, &r.ExpiresAt, &r.CreatedAt, &r.UpdatedAt,
		)
		if err != nil {
			stats.ReservationsFailed++
			logger.Printf("   ⚠️  Failed to scan reservation: %v\n", err)
			continue
		}

		// Check if listing exists in microservice (product_id -> listing_id mapping)
		var listingExists bool
		checkQuery := `SELECT EXISTS(SELECT 1 FROM listings WHERE id = $1)`
		err = microserviceDB.QueryRow(checkQuery, r.ProductID).Scan(&listingExists)
		if err != nil {
			stats.ReservationsFailed++
			logger.Printf("   ⚠️  Failed to check listing %d: %v\n", r.ProductID, err)
			continue
		}

		if !listingExists {
			stats.ReservationsSkipped++
			if cfg.Verbose {
				logger.Printf("   ⏭️  Skipped reservation %d: listing %d not found\n", r.ID, r.ProductID)
			}
			continue
		}

		// Check if reservation already exists
		var existingID int64
		checkExistingQuery := `SELECT id FROM inventory_reservations WHERE id = $1`
		err = microserviceDB.QueryRow(checkExistingQuery, r.ID).Scan(&existingID)
		if err == nil {
			stats.ReservationsSkipped++
			if cfg.Verbose {
				logger.Printf("   ⏭️  Skipped reservation %d: already exists\n", r.ID)
			}
			continue
		} else if err != sql.ErrNoRows {
			stats.ReservationsFailed++
			logger.Printf("   ⚠️  Failed to check existing reservation %d: %v\n", r.ID, err)
			continue
		}

		// Check if order_id exists (if not NULL, verify FK)
		if r.OrderID != nil && *r.OrderID > 0 {
			var orderExists bool
			checkOrderQuery := `SELECT EXISTS(SELECT 1 FROM orders WHERE id = $1)`
			err = microserviceDB.QueryRow(checkOrderQuery, *r.OrderID).Scan(&orderExists)
			if err != nil {
				stats.ReservationsFailed++
				logger.Printf("   ⚠️  Failed to check order %d: %v\n", *r.OrderID, err)
				continue
			}

			if !orderExists {
				// Set order_id to NULL if order doesn't exist
				if cfg.Verbose {
					logger.Printf("   ⚠️  Order %d not found, setting order_id to NULL for reservation %d\n", *r.OrderID, r.ID)
				}
				r.OrderID = nil
			}
		}

		// Insert into microservice (dry-run or execute)
		if cfg.DryRun {
			if cfg.Verbose {
				logger.Printf("   [DRY-RUN] Would insert reservation %d (listing=%d, quantity=%d, expires=%s)\n",
					r.ID, r.ProductID, r.Quantity, r.ExpiresAt.Format(time.RFC3339))
			}
			stats.ReservationsMigrated++
		} else {
			insertQuery := `
				INSERT INTO inventory_reservations
					(id, listing_id, variant_id, order_id, quantity, status, expires_at, created_at, updated_at)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
			`
			_, err = tx.Exec(insertQuery,
				r.ID, r.ProductID, r.VariantID, r.OrderID, r.Quantity,
				r.Status, r.ExpiresAt, r.CreatedAt, r.UpdatedAt,
			)
			if err != nil {
				stats.ReservationsFailed++
				logger.Printf("   ⚠️  Failed to insert reservation %d: %v\n", r.ID, err)
				continue
			}

			if cfg.Verbose {
				logger.Printf("   ✅ Inserted reservation %d (listing=%d, quantity=%d)\n",
					r.ID, r.ProductID, r.Quantity)
			}
			stats.ReservationsMigrated++
		}
	}

	if err = rows.Err(); err != nil {
		return fmt.Errorf("rows iteration failed: %w", err)
	}

	// Commit transaction
	if !cfg.DryRun && tx != nil {
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit transaction: %w", err)
		}
		logger.Println("   ✅ Transaction committed successfully")
	}

	// Verify FK constraints
	if !cfg.DryRun {
		logger.Println()
		logger.Println("🔍 Verifying foreign key constraints...")

		// Check orphaned reservations (listing not found)
		var orphanedCount int
		orphanQuery := `
			SELECT COUNT(*)
			FROM inventory_reservations ir
			LEFT JOIN listings l ON ir.listing_id = l.id
			WHERE l.id IS NULL
		`
		err = microserviceDB.QueryRow(orphanQuery).Scan(&orphanedCount)
		if err != nil {
			logger.Printf("   ⚠️  Failed to verify FK constraints: %v\n", err)
		} else if orphanedCount > 0 {
			logger.Printf("   ⚠️  WARNING: Found %d orphaned reservations (listing not found)\n", orphanedCount)
		} else {
			logger.Println("   ✅ All FK constraints valid")
		}
	}

	return nil
}

func printSummary(stats *MigrationStats, dryRun bool, logger *log.Logger) {
	logger.Println()
	logger.Println("═══════════════════════════════════════════════════════════")
	logger.Println("📊 MIGRATION SUMMARY")
	logger.Println("═══════════════════════════════════════════════════════════")
	logger.Printf("Mode:                    %s\n", mode(dryRun))
	logger.Printf("Duration:                %s\n", stats.EndTime.Sub(stats.StartTime).Round(time.Millisecond))
	logger.Println()
	logger.Println("Inventory Reservations:")
	logger.Printf("  Total found:           %d\n", stats.ReservationsTotal)
	logger.Printf("  Successfully migrated: %d\n", stats.ReservationsMigrated)
	logger.Printf("  Skipped:               %d\n", stats.ReservationsSkipped)
	logger.Printf("  Failed:                %d\n", stats.ReservationsFailed)
	logger.Println("═══════════════════════════════════════════════════════════")

	if dryRun {
		logger.Println()
		logger.Println("ℹ️  This was a DRY RUN. No data was written.")
		logger.Println("To execute migration, run with --dry-run=false")
	} else {
		logger.Println()
		logger.Println("✅ Migration completed successfully!")
	}
}

func mode(dryRun bool) string {
	if dryRun {
		return "DRY-RUN (no writes)"
	}
	return "EXECUTE (writes enabled)"
}
