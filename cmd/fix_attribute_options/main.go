package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	// Target database (microservice)
	targetDSN = "postgres://listings_user:listings_secret@localhost:35434/listings_dev_db?sslmode=disable"
)

// AttributeOptionOld represents the old format with string label
type AttributeOptionOld struct {
	Value string `json:"value"`
	Label string `json:"label"` // Old format: plain string
}

// AttributeOptionNew represents the new format with i18n label
type AttributeOptionNew struct {
	Value string            `json:"value"`
	Label map[string]string `json:"label"` // New format: i18n object
}

// Stats tracks processing statistics
type Stats struct {
	TotalAttributes    int
	ProcessedOptions   int
	SkippedAttributes  int
	SkippedOptions     int
	Errors             int
	ConvertedOptions   int
}

func main() {
	// Flags
	dryRun := flag.Bool("dry-run", true, "Perform dry run without actual changes (default: true)")
	verbose := flag.Bool("v", false, "Verbose output")
	flag.Parse()

	// Setup logger
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	})

	ctx := context.Background()
	stats := &Stats{}

	log.Info().Msg("=== Attribute Options Label Format Fixer ===")
	log.Info().Bool("dry_run", *dryRun).Msg("Mode")
	if *dryRun {
		log.Warn().Msg("Running in DRY-RUN mode. No changes will be saved.")
	} else {
		log.Info().Msg("Running in LIVE mode. Changes will be committed!")
	}
	log.Info().Msg("")

	// Connect to database
	log.Info().Msg("Connecting to database...")
	pool, err := pgxpool.New(ctx, targetDSN)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}
	defer pool.Close()

	// Test connection
	if err := pool.Ping(ctx); err != nil {
		log.Fatal().Err(err).Msg("Failed to ping database")
	}
	log.Info().Msg("✓ Database connection established")
	log.Info().Msg("")

	// Fetch attributes with options
	log.Info().Msg("Fetching attributes with options...")
	attributes, err := fetchAttributesWithOptions(ctx, pool)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to fetch attributes")
	}
	stats.TotalAttributes = len(attributes)
	log.Info().Int("count", stats.TotalAttributes).Msg("✓ Attributes fetched")
	log.Info().Msg("")

	if stats.TotalAttributes == 0 {
		log.Info().Msg("No attributes with options found. Nothing to do.")
		return
	}

	// Process each attribute
	log.Info().Msg("Processing attributes...")
	updates := make([]AttributeUpdate, 0)

	for _, attr := range attributes {
		needsUpdate, newOptions, convertedCount, err := processAttributeOptions(attr, *verbose)
		if err != nil {
			log.Error().
				Err(err).
				Int32("id", attr.ID).
				Str("code", attr.Code).
				Msg("Failed to process attribute")
			stats.Errors++
			continue
		}

		if needsUpdate {
			updates = append(updates, AttributeUpdate{
				ID:      attr.ID,
				Code:    attr.Code,
				Options: newOptions,
			})
			stats.ProcessedOptions++
			stats.ConvertedOptions += convertedCount

			if *verbose {
				log.Info().
					Int32("id", attr.ID).
					Str("code", attr.Code).
					Int("converted", convertedCount).
					Msg("Attribute needs update")
			}
		} else {
			stats.SkippedAttributes++
			if *verbose {
				log.Debug().
					Int32("id", attr.ID).
					Str("code", attr.Code).
					Msg("Attribute skipped (already correct format)")
			}
		}
	}

	log.Info().Msg("")
	log.Info().Msg("Processing Summary:")
	log.Info().Int("total_attributes", stats.TotalAttributes).Msg("  Total attributes")
	log.Info().Int("processed", stats.ProcessedOptions).Msg("  Attributes to update")
	log.Info().Int("skipped", stats.SkippedAttributes).Msg("  Attributes skipped")
	log.Info().Int("converted_options", stats.ConvertedOptions).Msg("  Options converted")
	log.Info().Int("errors", stats.Errors).Msg("  Errors")
	log.Info().Msg("")

	if len(updates) == 0 {
		log.Info().Msg("✓ All attributes already have correct format!")
		return
	}

	if *dryRun {
		log.Info().Msg("DRY-RUN: Showing sample updates (max 5):")
		for i, update := range updates {
			if i >= 5 {
				break
			}
			log.Info().
				Int32("id", update.ID).
				Str("code", update.Code).
				Interface("options", update.Options).
				Msg("Would update")
		}
		log.Info().Msg("")
		log.Info().Msg("DRY-RUN complete. Run with --dry-run=false to apply changes.")
		return
	}

	// Apply updates
	log.Info().Msg("Applying updates...")
	updated, err := applyUpdates(ctx, pool, updates, *verbose)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to apply updates")
	}

	log.Info().Msg("")
	log.Info().Int("updated", updated).Msg("✓ Updates applied successfully!")

	// Validate
	log.Info().Msg("")
	log.Info().Msg("Validating updates...")
	if err := validateUpdates(ctx, pool, updates); err != nil {
		log.Fatal().Err(err).Msg("Validation failed")
	}

	log.Info().Msg("✓ All validations passed!")
	log.Info().Msg("")
	log.Info().Msg("=== Completed successfully ===")
}

// AttributeRecord represents an attribute record from database
type AttributeRecord struct {
	ID      int32           `db:"id"`
	Code    string          `db:"code"`
	Options json.RawMessage `db:"options"`
}

// AttributeUpdate represents an attribute to be updated
type AttributeUpdate struct {
	ID      int32
	Code    string
	Options []AttributeOptionNew
}

func fetchAttributesWithOptions(ctx context.Context, pool *pgxpool.Pool) ([]AttributeRecord, error) {
	query := `
		SELECT id, code, options
		FROM attributes
		WHERE options IS NOT NULL
		  AND jsonb_typeof(options) = 'array'
		  AND jsonb_array_length(options) > 0
		ORDER BY id
	`

	rows, err := pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	attributes := make([]AttributeRecord, 0)
	for rows.Next() {
		var attr AttributeRecord
		if err := rows.Scan(&attr.ID, &attr.Code, &attr.Options); err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}
		attributes = append(attributes, attr)
	}

	return attributes, rows.Err()
}

func processAttributeOptions(attr AttributeRecord, verbose bool) (bool, []AttributeOptionNew, int, error) {
	// Try to unmarshal as array of interface{}
	var rawOptions []interface{}
	if err := json.Unmarshal(attr.Options, &rawOptions); err != nil {
		return false, nil, 0, fmt.Errorf("failed to unmarshal options: %w", err)
	}

	if len(rawOptions) == 0 {
		return false, nil, 0, nil
	}

	needsUpdate := false
	convertedCount := 0
	newOptions := make([]AttributeOptionNew, 0, len(rawOptions))

	for i, rawOpt := range rawOptions {
		// Check if it's an object
		optMap, ok := rawOpt.(map[string]interface{})
		if !ok {
			// It's a plain string or other type - skip this format
			if verbose {
				log.Debug().
					Int32("id", attr.ID).
					Str("code", attr.Code).
					Int("index", i).
					Msg("Skipping non-object option")
			}
			continue
		}

		// Check if it has "value" and "label" fields
		value, hasValue := optMap["value"].(string)
		label := optMap["label"]

		if !hasValue {
			if verbose {
				log.Debug().
					Int32("id", attr.ID).
					Str("code", attr.Code).
					Int("index", i).
					Msg("Skipping option without value field")
			}
			continue
		}

		var newOption AttributeOptionNew
		newOption.Value = value

		// Check label type
		switch labelVal := label.(type) {
		case string:
			// Old format: label is a string - convert it!
			needsUpdate = true
			convertedCount++
			newOption.Label = map[string]string{
				"en": labelVal,
				"ru": labelVal,
				"sr": labelVal,
			}
			if verbose {
				log.Debug().
					Int32("id", attr.ID).
					Str("code", attr.Code).
					Str("value", value).
					Str("old_label", labelVal).
					Msg("Converting string label to i18n")
			}

		case map[string]interface{}:
			// New format: label is already an object
			// Verify it has the expected structure
			i18nLabel := make(map[string]string)
			for k, v := range labelVal {
				if strVal, ok := v.(string); ok {
					i18nLabel[k] = strVal
				}
			}

			// Check if it has at least en/ru/sr keys
			if _, hasEn := i18nLabel["en"]; hasEn {
				newOption.Label = i18nLabel
				if verbose {
					log.Debug().
						Int32("id", attr.ID).
						Str("code", attr.Code).
						Str("value", value).
						Msg("Label already in i18n format")
				}
			} else {
				// Object but missing expected keys - treat as error
				return false, nil, 0, fmt.Errorf("label object missing 'en' key for option value=%s", value)
			}

		default:
			// Unknown label type
			return false, nil, 0, fmt.Errorf("unexpected label type %T for option value=%s", label, value)
		}

		newOptions = append(newOptions, newOption)
	}

	return needsUpdate, newOptions, convertedCount, nil
}

func applyUpdates(ctx context.Context, pool *pgxpool.Pool, updates []AttributeUpdate, verbose bool) (int, error) {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	updateQuery := `
		UPDATE attributes
		SET options = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`

	updated := 0
	for _, update := range updates {
		// Marshal new options to JSON
		optionsJSON, err := json.Marshal(update.Options)
		if err != nil {
			return updated, fmt.Errorf("failed to marshal options for %s: %w", update.Code, err)
		}

		// Validate JSON before update
		var testUnmarshal []AttributeOptionNew
		if err := json.Unmarshal(optionsJSON, &testUnmarshal); err != nil {
			return updated, fmt.Errorf("validation failed for %s: marshaled JSON is invalid: %w", update.Code, err)
		}

		// Execute update
		result, err := tx.Exec(ctx, updateQuery, optionsJSON, update.ID)
		if err != nil {
			return updated, fmt.Errorf("failed to update attribute %s: %w", update.Code, err)
		}

		rowsAffected := result.RowsAffected()
		if rowsAffected != 1 {
			return updated, fmt.Errorf("expected to update 1 row for %s, but updated %d", update.Code, rowsAffected)
		}

		updated++
		if verbose && updated%10 == 0 {
			log.Info().Int("updated", updated).Int("total", len(updates)).Msg("Progress")
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return updated, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return updated, nil
}

func validateUpdates(ctx context.Context, pool *pgxpool.Pool, updates []AttributeUpdate) error {
	// Collect IDs
	ids := make([]int32, len(updates))
	for i, update := range updates {
		ids[i] = update.ID
	}

	query := `
		SELECT id, code, options
		FROM attributes
		WHERE id = ANY($1)
		ORDER BY id
	`

	rows, err := pool.Query(ctx, query, ids)
	if err != nil {
		return fmt.Errorf("validation query failed: %w", err)
	}
	defer rows.Close()

	validated := 0
	for rows.Next() {
		var attr AttributeRecord
		if err := rows.Scan(&attr.ID, &attr.Code, &attr.Options); err != nil {
			return fmt.Errorf("validation scan failed: %w", err)
		}

		// Unmarshal and validate structure
		var options []AttributeOptionNew
		if err := json.Unmarshal(attr.Options, &options); err != nil {
			return fmt.Errorf("validation failed for %s: %w", attr.Code, err)
		}

		// Check each option has i18n label
		for _, opt := range options {
			if opt.Label == nil {
				return fmt.Errorf("validation failed for %s: option %s has nil label", attr.Code, opt.Value)
			}

			// Check for required locales
			if _, hasEn := opt.Label["en"]; !hasEn {
				return fmt.Errorf("validation failed for %s: option %s missing 'en' locale", attr.Code, opt.Value)
			}
		}

		validated++
	}

	if validated != len(updates) {
		return fmt.Errorf("validation count mismatch: expected %d, validated %d", len(updates), validated)
	}

	log.Info().Int("validated", validated).Msg("  Records validated")

	// Show sample
	log.Info().Msg("  Sample validated records (max 3):")
	for i := 0; i < len(updates) && i < 3; i++ {
		log.Info().
			Int32("id", updates[i].ID).
			Str("code", updates[i].Code).
			Int("options_count", len(updates[i].Options)).
			Msg("    ✓")
	}

	return rows.Err()
}
