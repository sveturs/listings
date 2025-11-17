package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	// Source database (monolith)
	sourceDSN = "postgres://postgres:mX3g1XGhMRUZEX3l@localhost:5433/svetubd?sslmode=disable"

	// Target database (microservice)
	targetDSN = "postgres://listings_user:listings_secret@localhost:35434/listings_dev_db?sslmode=disable"
)

// SourceAttribute represents the monolith's unified_attributes structure
type SourceAttribute struct {
	ID                              int
	Code                            string
	Name                            string // VARCHAR - will convert to JSONB
	DisplayName                     string // VARCHAR - will convert to JSONB
	AttributeType                   string
	Purpose                         string
	Options                         []byte // JSONB
	ValidationRules                 []byte // JSONB
	UISettings                      []byte // JSONB
	IsSearchable                    bool
	IsFilterable                    bool
	IsRequired                      bool
	AffectsStock                    bool
	AffectsPrice                    bool
	SortOrder                       int
	IsActive                        bool
	CreatedAt                       time.Time
	UpdatedAt                       time.Time
	LegacyCategoryAttributeID       *int
	LegacyProductVariantAttributeID *int
	IsVariantCompatible             bool
	Icon                            string
	ShowInCard                      bool
}

// I18nField represents a multilingual JSONB field
type I18nField struct {
	En string `json:"en"`
	Ru string `json:"ru"`
	Sr string `json:"sr"`
}

func main() {
	dryRun := flag.Bool("dry-run", false, "Perform dry run without actual migration")
	verbose := flag.Bool("v", false, "Verbose output")
	flag.Parse()

	ctx := context.Background()

	log.Println("=== Attributes Migration Tool ===")
	log.Printf("Mode: %s", func() string {
		if *dryRun {
			return "DRY RUN"
		}
		return "LIVE MIGRATION"
	}())
	log.Println()

	// Connect to source database
	log.Println("Connecting to source database (monolith)...")
	sourcePool, err := pgxpool.New(ctx, sourceDSN)
	if err != nil {
		log.Fatalf("Failed to connect to source DB: %v", err)
	}
	defer sourcePool.Close()

	// Connect to target database
	log.Println("Connecting to target database (microservice)...")
	targetPool, err := pgxpool.New(ctx, targetDSN)
	if err != nil {
		log.Fatalf("Failed to connect to target DB: %v", err)
	}
	defer targetPool.Close()

	// Fetch source data
	log.Println("Fetching attributes from source database...")
	attributes, err := fetchSourceAttributes(ctx, sourcePool)
	if err != nil {
		log.Fatalf("Failed to fetch source attributes: %v", err)
	}
	log.Printf("✓ Found %d attributes in source database\n\n", len(attributes))

	// Check existing attributes in target
	log.Println("Checking existing attributes in target database...")
	existingCodes, err := fetchExistingCodes(ctx, targetPool)
	if err != nil {
		log.Fatalf("Failed to fetch existing codes: %v", err)
	}
	log.Printf("✓ Found %d existing attributes in target database\n\n", len(existingCodes))

	// Filter attributes to migrate
	toMigrate := make([]*SourceAttribute, 0)
	skipped := 0
	for _, attr := range attributes {
		if existingCodes[attr.Code] {
			skipped++
			if *verbose {
				log.Printf("  [SKIP] %s (already exists)\n", attr.Code)
			}
		} else {
			toMigrate = append(toMigrate, attr)
		}
	}

	log.Printf("Migration plan:\n")
	log.Printf("  - Total attributes: %d\n", len(attributes))
	log.Printf("  - Already migrated: %d\n", skipped)
	log.Printf("  - To migrate: %d\n\n", len(toMigrate))

	if len(toMigrate) == 0 {
		log.Println("✓ Nothing to migrate. All attributes already exist.")
		return
	}

	if *dryRun {
		log.Println("DRY RUN - Showing first 5 attributes to be migrated:")
		for i, attr := range toMigrate {
			if i >= 5 {
				break
			}
			log.Printf("  [%d] %s: %s → %s\n", attr.ID, attr.Code, attr.Name, attr.DisplayName)
		}
		log.Println("\nDRY RUN complete. Use without --dry-run to perform actual migration.")
		return
	}

	// Perform migration
	log.Println("Starting migration...")
	migrated, err := migrateAttributes(ctx, targetPool, toMigrate, *verbose)
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	log.Printf("\n✓ Migration complete!\n")
	log.Printf("  - Successfully migrated: %d attributes\n", migrated)

	// Validate
	log.Println("\nValidating migration...")
	if err := validateMigration(ctx, targetPool, toMigrate); err != nil {
		log.Fatalf("Validation failed: %v", err)
	}

	log.Println("✓ All validations passed!")
}

func fetchSourceAttributes(ctx context.Context, pool *pgxpool.Pool) ([]*SourceAttribute, error) {
	query := `
		SELECT
			id, code, name, display_name, attribute_type, purpose,
			COALESCE(options, '{}'::jsonb),
			COALESCE(validation_rules, '{}'::jsonb),
			COALESCE(ui_settings, '{}'::jsonb),
			COALESCE(is_searchable, false),
			COALESCE(is_filterable, false),
			COALESCE(is_required, false),
			COALESCE(affects_stock, false),
			COALESCE(affects_price, false),
			COALESCE(sort_order, 0),
			COALESCE(is_active, true),
			COALESCE(created_at, CURRENT_TIMESTAMP),
			COALESCE(updated_at, CURRENT_TIMESTAMP),
			legacy_category_attribute_id,
			legacy_product_variant_attribute_id,
			COALESCE(is_variant_compatible, false),
			COALESCE(icon, ''),
			COALESCE(show_in_card, false)
		FROM unified_attributes
		ORDER BY id
	`

	rows, err := pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	attributes := make([]*SourceAttribute, 0)
	for rows.Next() {
		var attr SourceAttribute
		err := rows.Scan(
			&attr.ID, &attr.Code, &attr.Name, &attr.DisplayName,
			&attr.AttributeType, &attr.Purpose,
			&attr.Options, &attr.ValidationRules, &attr.UISettings,
			&attr.IsSearchable, &attr.IsFilterable, &attr.IsRequired,
			&attr.AffectsStock, &attr.AffectsPrice, &attr.SortOrder,
			&attr.IsActive, &attr.CreatedAt, &attr.UpdatedAt,
			&attr.LegacyCategoryAttributeID, &attr.LegacyProductVariantAttributeID,
			&attr.IsVariantCompatible, &attr.Icon, &attr.ShowInCard,
		)
		if err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}
		attributes = append(attributes, &attr)
	}

	return attributes, rows.Err()
}

func fetchExistingCodes(ctx context.Context, pool *pgxpool.Pool) (map[string]bool, error) {
	query := `SELECT code FROM attributes`

	rows, err := pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	codes := make(map[string]bool)
	for rows.Next() {
		var code string
		if err := rows.Scan(&code); err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}
		codes[code] = true
	}

	return codes, rows.Err()
}

func migrateAttributes(ctx context.Context, pool *pgxpool.Pool, attributes []*SourceAttribute, verbose bool) (int, error) {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	insertQuery := `
		INSERT INTO attributes (
			id, code, name, display_name, attribute_type, purpose,
			options, validation_rules, ui_settings,
			is_searchable, is_filterable, is_required,
			affects_stock, affects_price, sort_order, is_active,
			created_at, updated_at,
			legacy_category_attribute_id, legacy_product_variant_attribute_id,
			is_variant_compatible, icon, show_in_card
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12,
			$13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23
		)
	`

	batch := &pgx.Batch{}
	for _, attr := range attributes {
		// Convert VARCHAR to JSONB i18n format
		nameJSON, err := json.Marshal(I18nField{
			En: attr.Name,
			Ru: attr.Name,
			Sr: attr.Name,
		})
		if err != nil {
			return 0, fmt.Errorf("failed to marshal name for %s: %w", attr.Code, err)
		}

		displayNameJSON, err := json.Marshal(I18nField{
			En: attr.DisplayName,
			Ru: attr.DisplayName,
			Sr: attr.DisplayName,
		})
		if err != nil {
			return 0, fmt.Errorf("failed to marshal display_name for %s: %w", attr.Code, err)
		}

		batch.Queue(insertQuery,
			attr.ID, attr.Code, nameJSON, displayNameJSON,
			attr.AttributeType, attr.Purpose,
			attr.Options, attr.ValidationRules, attr.UISettings,
			attr.IsSearchable, attr.IsFilterable, attr.IsRequired,
			attr.AffectsStock, attr.AffectsPrice, attr.SortOrder, attr.IsActive,
			attr.CreatedAt, attr.UpdatedAt,
			attr.LegacyCategoryAttributeID, attr.LegacyProductVariantAttributeID,
			attr.IsVariantCompatible, attr.Icon, attr.ShowInCard,
		)

		if verbose && len(batch.QueuedQueries)%50 == 0 {
			log.Printf("  Queued %d attributes...\n", len(batch.QueuedQueries))
		}
	}

	log.Printf("Executing batch insert of %d attributes...\n", len(batch.QueuedQueries))
	results := tx.SendBatch(ctx, batch)

	migrated := 0
	for i := 0; i < len(attributes); i++ {
		_, err := results.Exec()
		if err != nil {
			results.Close()
			return migrated, fmt.Errorf("failed to insert attribute %s: %w", attributes[i].Code, err)
		}
		migrated++

		if verbose && migrated%50 == 0 {
			log.Printf("  Inserted %d/%d attributes...\n", migrated, len(attributes))
		}
	}

	// Close batch results before running sequence update
	results.Close()

	// Update sequence
	_, err = tx.Exec(ctx, "SELECT setval('attributes_id_seq', (SELECT COALESCE(MAX(id), 0) FROM attributes))")
	if err != nil {
		return migrated, fmt.Errorf("failed to update sequence: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return migrated, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return migrated, nil
}

func validateMigration(ctx context.Context, pool *pgxpool.Pool, migrated []*SourceAttribute) error {
	// Check counts
	var count int
	err := pool.QueryRow(ctx, "SELECT COUNT(*) FROM attributes").Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to count attributes: %w", err)
	}

	log.Printf("  - Total attributes in target: %d\n", count)

	// Validate JSONB structure for first few records
	query := `
		SELECT id, code, name, display_name
		FROM attributes
		WHERE id = ANY($1)
		ORDER BY id
		LIMIT 5
	`

	ids := make([]int, 0, 5)
	for i := 0; i < len(migrated) && i < 5; i++ {
		ids = append(ids, migrated[i].ID)
	}

	rows, err := pool.Query(ctx, query, ids)
	if err != nil {
		return fmt.Errorf("validation query failed: %w", err)
	}
	defer rows.Close()

	log.Println("  - Sample migrated attributes:")
	for rows.Next() {
		var id int
		var code string
		var name, displayName []byte

		if err := rows.Scan(&id, &code, &name, &displayName); err != nil {
			return fmt.Errorf("validation scan failed: %w", err)
		}

		var nameI18n, displayNameI18n I18nField
		if err := json.Unmarshal(name, &nameI18n); err != nil {
			return fmt.Errorf("invalid name JSONB for %s: %w", code, err)
		}
		if err := json.Unmarshal(displayName, &displayNameI18n); err != nil {
			return fmt.Errorf("invalid display_name JSONB for %s: %w", code, err)
		}

		log.Printf("    [%d] %s: name=%s, display_name=%s\n",
			id, code, nameI18n.En, displayNameI18n.En)
	}

	return rows.Err()
}
