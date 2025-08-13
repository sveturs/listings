package translation_admin

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"backend/internal/domain/models"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/rs/zerolog"
)

// Repository handles database operations for translation admin
type Repository struct {
	db     *sqlx.DB
	logger zerolog.Logger
}

// NewRepository creates a new translation admin repository
func NewRepository(db *sqlx.DB, logger zerolog.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// GetTranslations retrieves translations with filters
func (r *Repository) GetTranslations(ctx context.Context, filters map[string]interface{}) ([]models.Translation, error) {
	query := `SELECT id, entity_type, entity_id, language, field_name, 
	          translated_text, is_machine_translated, is_verified, 
	          created_at, updated_at, metadata, COALESCE(version, 1) as version
	          FROM translations WHERE 1=1`

	args := []interface{}{}
	argCount := 0

	// Build dynamic query based on filters
	if entityType, ok := filters["entity_type"].(string); ok && entityType != "" {
		argCount++
		query += fmt.Sprintf(" AND entity_type = $%d", argCount)
		args = append(args, entityType)
	}

	if entityID, ok := filters["entity_id"].(int); ok && entityID > 0 {
		argCount++
		query += fmt.Sprintf(" AND entity_id = $%d", argCount)
		args = append(args, entityID)
	}

	if language, ok := filters["language"].(string); ok && language != "" {
		argCount++
		query += fmt.Sprintf(" AND language = $%d", argCount)
		args = append(args, language)
	}

	if fieldName, ok := filters["field_name"].(string); ok && fieldName != "" {
		argCount++
		query += fmt.Sprintf(" AND field_name = $%d", argCount)
		args = append(args, fieldName)
	}

	if isVerified, ok := filters["is_verified"].(bool); ok {
		argCount++
		query += fmt.Sprintf(" AND is_verified = $%d", argCount)
		args = append(args, isVerified)
	}

	// Add ordering and limit
	query += " ORDER BY created_at DESC"

	if limit, ok := filters["limit"].(int); ok && limit > 0 {
		argCount++
		query += fmt.Sprintf(" LIMIT $%d", argCount)
		args = append(args, limit)
	}

	if offset, ok := filters["offset"].(int); ok && offset > 0 {
		argCount++
		query += fmt.Sprintf(" OFFSET $%d", argCount)
		args = append(args, offset)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query translations: %w", err)
	}
	defer rows.Close()

	var translations []models.Translation
	for rows.Next() {
		var t models.Translation
		var metadata sql.NullString
		var version sql.NullInt64

		err := rows.Scan(
			&t.ID, &t.EntityType, &t.EntityID, &t.Language, &t.FieldName,
			&t.TranslatedText, &t.IsMachineTranslated, &t.IsVerified,
			&t.CreatedAt, &t.UpdatedAt, &metadata, &version,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan translation: %w", err)
		}

		if metadata.Valid {
			if err := json.Unmarshal([]byte(metadata.String), &t.Metadata); err != nil {
				r.logger.Error().Err(err).Msg("Failed to unmarshal metadata")
			}
		}

		if version.Valid {
			t.Version = int(version.Int64)
		} else {
			t.Version = 1
		}

		translations = append(translations, t)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return translations, nil
}

// CreateTranslation creates a new translation
func (r *Repository) CreateTranslation(ctx context.Context, translation *models.Translation) error {
	metadataJSON, err := json.Marshal(translation.Metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	query := `INSERT INTO translations (entity_type, entity_id, language, field_name, 
	          translated_text, is_machine_translated, is_verified, metadata)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	          RETURNING id, created_at, updated_at`

	err = r.db.QueryRowContext(
		ctx, query,
		translation.EntityType, translation.EntityID, translation.Language, translation.FieldName,
		translation.TranslatedText, translation.IsMachineTranslated, translation.IsVerified, metadataJSON,
	).Scan(&translation.ID, &translation.CreatedAt, &translation.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create translation: %w", err)
	}

	return nil
}

// UpdateTranslation updates an existing translation
func (r *Repository) UpdateTranslation(ctx context.Context, translation *models.Translation) error {
	metadataJSON, err := json.Marshal(translation.Metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	query := `UPDATE translations SET 
	          translated_text = $1, is_machine_translated = $2, is_verified = $3, 
	          metadata = $4, updated_at = CURRENT_TIMESTAMP
	          WHERE id = $5`

	result, err := r.db.ExecContext(
		ctx, query,
		translation.TranslatedText, translation.IsMachineTranslated, translation.IsVerified,
		metadataJSON, translation.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update translation: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("translation not found")
	}

	return nil
}

// DeleteTranslation deletes a translation
func (r *Repository) DeleteTranslation(ctx context.Context, id int) error {
	query := `DELETE FROM translations WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete translation: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("translation not found")
	}

	return nil
}

// GetTranslationVersions retrieves version history for a translation
func (r *Repository) GetTranslationVersions(ctx context.Context, translationID int) ([]models.TranslationVersion, error) {
	query := `SELECT id, translation_id, entity_type, entity_id, field_name,
	          language, translated_text, previous_text, version, 
	          change_type, changed_by, changed_at, change_reason, metadata
	          FROM translation_versions 
	          WHERE translation_id = $1
	          ORDER BY version DESC`

	rows, err := r.db.QueryContext(ctx, query, translationID)
	if err != nil {
		return nil, fmt.Errorf("failed to query versions: %w", err)
	}
	defer rows.Close()

	var versions []models.TranslationVersion
	for rows.Next() {
		var v models.TranslationVersion
		var previousText, changeReason sql.NullString
		var changedBy sql.NullInt64
		var metadata sql.NullString

		err := rows.Scan(
			&v.ID, &v.TranslationID, &v.EntityType, &v.EntityID, &v.FieldName,
			&v.Language, &v.TranslatedText, &previousText, &v.Version,
			&v.ChangeType, &changedBy, &v.ChangedAt, &changeReason, &metadata,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan version: %w", err)
		}

		if previousText.Valid {
			v.PreviousText = &previousText.String
		}

		if changedBy.Valid {
			v.ChangedBy = &changedBy.Int64
		}

		if changeReason.Valid {
			v.ChangeReason = &changeReason.String
		}

		if metadata.Valid {
			if err := json.Unmarshal([]byte(metadata.String), &v.Metadata); err != nil {
				r.logger.Error().Err(err).Msg("Failed to unmarshal metadata")
			}
		}

		versions = append(versions, v)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return versions, nil
}

// GetVersionsByEntity retrieves all versions for a specific entity
func (r *Repository) GetVersionsByEntity(ctx context.Context, entityType string, entityID int) ([]models.TranslationVersion, error) {
	query := `SELECT id, translation_id, entity_type, entity_id, field_name,
	          language, translated_text, previous_text, version, 
	          change_type, changed_by, changed_at, change_reason, metadata
	          FROM translation_versions 
	          WHERE entity_type = $1 AND entity_id = $2
	          ORDER BY changed_at DESC`

	rows, err := r.db.QueryContext(ctx, query, entityType, entityID)
	if err != nil {
		return nil, fmt.Errorf("failed to query versions by entity: %w", err)
	}
	defer rows.Close()

	var versions []models.TranslationVersion
	for rows.Next() {
		var v models.TranslationVersion
		var previousText, changeReason sql.NullString
		var changedBy sql.NullInt64
		var metadata sql.NullString

		err := rows.Scan(
			&v.ID, &v.TranslationID, &v.EntityType, &v.EntityID, &v.FieldName,
			&v.Language, &v.TranslatedText, &previousText, &v.Version,
			&v.ChangeType, &changedBy, &v.ChangedAt, &changeReason, &metadata,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan version: %w", err)
		}

		if previousText.Valid {
			v.PreviousText = &previousText.String
		}

		if changedBy.Valid {
			v.ChangedBy = &changedBy.Int64
		}

		if changeReason.Valid {
			v.ChangeReason = &changeReason.String
		}

		if metadata.Valid {
			if err := json.Unmarshal([]byte(metadata.String), &v.Metadata); err != nil {
				r.logger.Error().Err(err).Msg("Failed to unmarshal metadata")
			}
		}

		versions = append(versions, v)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return versions, nil
}

// GetVersionDiff gets the difference between two versions
func (r *Repository) GetVersionDiff(ctx context.Context, versionID1, versionID2 int) (*models.VersionDiff, error) {
	query := `SELECT id, translation_id, entity_type, entity_id, field_name,
	          language, translated_text, previous_text, version, 
	          change_type, changed_by, changed_at, change_reason
	          FROM translation_versions 
	          WHERE id IN ($1, $2)
	          ORDER BY version`

	rows, err := r.db.QueryContext(ctx, query, versionID1, versionID2)
	if err != nil {
		return nil, fmt.Errorf("failed to query versions for diff: %w", err)
	}
	defer rows.Close()

	var versions []models.TranslationVersion
	for rows.Next() {
		var v models.TranslationVersion
		var previousText, changeReason sql.NullString
		var changedBy sql.NullInt64

		err := rows.Scan(
			&v.ID, &v.TranslationID, &v.EntityType, &v.EntityID, &v.FieldName,
			&v.Language, &v.TranslatedText, &previousText, &v.Version,
			&v.ChangeType, &changedBy, &v.ChangedAt, &changeReason,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan version: %w", err)
		}

		if previousText.Valid {
			v.PreviousText = &previousText.String
		}
		if changedBy.Valid {
			v.ChangedBy = &changedBy.Int64
		}
		if changeReason.Valid {
			v.ChangeReason = &changeReason.String
		}

		versions = append(versions, v)
	}

	if len(versions) != 2 {
		return nil, fmt.Errorf("expected 2 versions, got %d", len(versions))
	}

	diff := &models.VersionDiff{
		Version1: versions[0],
		Version2: versions[1],
	}

	return diff, nil
}

// RollbackToVersion rolls back a translation to a specific version
func (r *Repository) RollbackToVersion(ctx context.Context, translationID int, versionID int, userID int) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Get the version to rollback to
	var version models.TranslationVersion
	query := `SELECT translated_text FROM translation_versions WHERE id = $1 AND translation_id = $2`
	err = tx.QueryRowContext(ctx, query, versionID, translationID).Scan(&version.TranslatedText)
	if err != nil {
		return fmt.Errorf("failed to get version: %w", err)
	}

	// Update the translation with the old version's text
	updateQuery := `UPDATE translations 
	                SET translated_text = $1, 
	                    last_modified_by = $2,
	                    last_modified_at = CURRENT_TIMESTAMP,
	                    current_version = current_version + 1
	                WHERE id = $3`

	_, err = tx.ExecContext(ctx, updateQuery, version.TranslatedText, userID, translationID)
	if err != nil {
		return fmt.Errorf("failed to update translation: %w", err)
	}

	// The trigger will automatically create a new version entry

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// CreateSyncConflict creates a new sync conflict record
func (r *Repository) CreateSyncConflict(ctx context.Context, conflict *models.TranslationSyncConflict) error {
	query := `INSERT INTO translation_sync_conflicts 
	          (source_type, target_type, entity_identifier, source_value, target_value, conflict_type)
	          VALUES ($1, $2, $3, $4, $5, $6)
	          RETURNING id, created_at`

	err := r.db.QueryRowContext(
		ctx, query,
		conflict.SourceType, conflict.TargetType, conflict.EntityIdentifier,
		conflict.SourceValue, conflict.TargetValue, conflict.ConflictType,
	).Scan(&conflict.ID, &conflict.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create sync conflict: %w", err)
	}

	return nil
}

// ResolveSyncConflict resolves a sync conflict
func (r *Repository) ResolveSyncConflict(ctx context.Context, id int, resolution string, userID int) error {
	query := `UPDATE translation_sync_conflicts 
	          SET resolved = true, resolved_by = $1, resolved_at = CURRENT_TIMESTAMP, resolution_type = $2
	          WHERE id = $3`

	result, err := r.db.ExecContext(ctx, query, userID, resolution, id)
	if err != nil {
		return fmt.Errorf("failed to resolve conflict: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("conflict not found")
	}

	return nil
}

// GetUnresolvedConflicts retrieves all unresolved conflicts
func (r *Repository) GetUnresolvedConflicts(ctx context.Context) ([]models.TranslationSyncConflict, error) {
	query := `SELECT id, source_type, target_type, entity_identifier, source_value, 
	          target_value, conflict_type, resolved, resolved_by, resolved_at, 
	          resolution_type, created_at
	          FROM translation_sync_conflicts
	          WHERE resolved = false
	          ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query conflicts: %w", err)
	}
	defer rows.Close()

	var conflicts []models.TranslationSyncConflict
	for rows.Next() {
		var c models.TranslationSyncConflict
		var sourceValue, targetValue sql.NullString
		var resolvedBy sql.NullInt64
		var resolvedAt sql.NullTime
		var resolutionType sql.NullString

		err := rows.Scan(
			&c.ID, &c.SourceType, &c.TargetType, &c.EntityIdentifier, &sourceValue,
			&targetValue, &c.ConflictType, &c.Resolved, &resolvedBy, &resolvedAt,
			&resolutionType, &c.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan conflict: %w", err)
		}

		if sourceValue.Valid {
			c.SourceValue = &sourceValue.String
		}
		if targetValue.Valid {
			c.TargetValue = &targetValue.String
		}
		if resolvedBy.Valid {
			val := int(resolvedBy.Int64)
			c.ResolvedBy = &val
		}
		if resolvedAt.Valid {
			c.ResolvedAt = &resolvedAt.Time
		}
		if resolutionType.Valid {
			c.ResolutionType = &resolutionType.String
		}

		conflicts = append(conflicts, c)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return conflicts, nil
}

// GetProviders retrieves all translation providers
func (r *Repository) GetProviders(ctx context.Context) ([]models.TranslationProvider, error) {
	query := `SELECT id, name, provider_type, settings, usage_limit, usage_current, 
	          is_active, priority, created_at, updated_at
	          FROM translation_providers
	          ORDER BY priority ASC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query providers: %w", err)
	}
	defer rows.Close()

	var providers []models.TranslationProvider
	for rows.Next() {
		var p models.TranslationProvider
		var settings sql.NullString
		var usageLimit sql.NullInt64

		err := rows.Scan(
			&p.ID, &p.Name, &p.ProviderType, &settings, &usageLimit, &p.UsageCurrent,
			&p.IsActive, &p.Priority, &p.CreatedAt, &p.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan provider: %w", err)
		}

		if settings.Valid {
			if err := json.Unmarshal([]byte(settings.String), &p.Settings); err != nil {
				r.logger.Error().Err(err).Msg("Failed to unmarshal settings")
			}
		}

		if usageLimit.Valid {
			val := int(usageLimit.Int64)
			p.UsageLimit = &val
		}

		providers = append(providers, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return providers, nil
}

// UpdateProvider updates a translation provider
func (r *Repository) UpdateProvider(ctx context.Context, provider *models.TranslationProvider) error {
	settingsJSON, err := json.Marshal(provider.Settings)
	if err != nil {
		return fmt.Errorf("failed to marshal settings: %w", err)
	}

	query := `UPDATE translation_providers SET 
	          name = $1, provider_type = $2, settings = $3, usage_limit = $4, 
	          usage_current = $5, is_active = $6, priority = $7
	          WHERE id = $8`

	result, err := r.db.ExecContext(
		ctx, query,
		provider.Name, provider.ProviderType, settingsJSON, provider.UsageLimit,
		provider.UsageCurrent, provider.IsActive, provider.Priority, provider.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update provider: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("provider not found")
	}

	return nil
}

// CreateTask creates a new translation task
func (r *Repository) CreateTask(ctx context.Context, task *models.TranslationTask) error {
	entityRefsJSON, err := json.Marshal(task.EntityReferences)
	if err != nil {
		return fmt.Errorf("failed to marshal entity references: %w", err)
	}

	metadataJSON, err := json.Marshal(task.Metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	query := `INSERT INTO translation_tasks 
	          (task_type, status, source_language, target_languages, entity_references, 
	           provider_id, created_by, assigned_to, metadata)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	          RETURNING id, created_at`

	err = r.db.QueryRowContext(
		ctx, query,
		task.TaskType, task.Status, task.SourceLanguage, pq.Array(task.TargetLanguages),
		entityRefsJSON, task.ProviderID, task.CreatedBy, task.AssignedTo, metadataJSON,
	).Scan(&task.ID, &task.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create task: %w", err)
	}

	return nil
}

// UpdateTask updates a translation task
func (r *Repository) UpdateTask(ctx context.Context, task *models.TranslationTask) error {
	metadataJSON, err := json.Marshal(task.Metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	query := `UPDATE translation_tasks SET 
	          status = $1, assigned_to = $2, started_at = $3, completed_at = $4, 
	          error_message = $5, metadata = $6
	          WHERE id = $7`

	result, err := r.db.ExecContext(
		ctx, query,
		task.Status, task.AssignedTo, task.StartedAt, task.CompletedAt,
		task.ErrorMessage, metadataJSON, task.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("task not found")
	}

	return nil
}

// GetQualityMetrics retrieves quality metrics for a translation
func (r *Repository) GetQualityMetrics(ctx context.Context, translationID int) (*models.TranslationQualityMetrics, error) {
	query := `SELECT id, translation_id, quality_score, character_count, word_count, 
	          has_placeholders, has_html_tags, checked_at, checked_by, issues
	          FROM translation_quality_metrics
	          WHERE translation_id = $1
	          ORDER BY checked_at DESC
	          LIMIT 1`

	var m models.TranslationQualityMetrics
	var qualityScore sql.NullFloat64
	var characterCount, wordCount sql.NullInt64
	var issues sql.NullString

	err := r.db.QueryRowContext(ctx, query, translationID).Scan(
		&m.ID, &m.TranslationID, &qualityScore, &characterCount, &wordCount,
		&m.HasPlaceholders, &m.HasHTMLTags, &m.CheckedAt, &m.CheckedBy, &issues,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("quality metrics not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query quality metrics: %w", err)
	}

	if qualityScore.Valid {
		m.QualityScore = &qualityScore.Float64
	}
	if characterCount.Valid {
		val := int(characterCount.Int64)
		m.CharacterCount = &val
	}
	if wordCount.Valid {
		val := int(wordCount.Int64)
		m.WordCount = &val
	}
	if issues.Valid {
		if err := json.Unmarshal([]byte(issues.String), &m.Issues); err != nil {
			r.logger.Error().Err(err).Msg("Failed to unmarshal issues")
		}
	}

	return &m, nil
}

// SaveQualityMetrics saves quality metrics for a translation
func (r *Repository) SaveQualityMetrics(ctx context.Context, metrics *models.TranslationQualityMetrics) error {
	issuesJSON, err := json.Marshal(metrics.Issues)
	if err != nil {
		return fmt.Errorf("failed to marshal issues: %w", err)
	}

	query := `INSERT INTO translation_quality_metrics 
	          (translation_id, quality_score, character_count, word_count, 
	           has_placeholders, has_html_tags, checked_by, issues)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	          ON CONFLICT (translation_id) DO UPDATE SET
	          quality_score = EXCLUDED.quality_score,
	          character_count = EXCLUDED.character_count,
	          word_count = EXCLUDED.word_count,
	          has_placeholders = EXCLUDED.has_placeholders,
	          has_html_tags = EXCLUDED.has_html_tags,
	          checked_at = CURRENT_TIMESTAMP,
	          checked_by = EXCLUDED.checked_by,
	          issues = EXCLUDED.issues
	          RETURNING id, checked_at`

	err = r.db.QueryRowContext(
		ctx, query,
		metrics.TranslationID, metrics.QualityScore, metrics.CharacterCount, metrics.WordCount,
		metrics.HasPlaceholders, metrics.HasHTMLTags, metrics.CheckedBy, issuesJSON,
	).Scan(&metrics.ID, &metrics.CheckedAt)
	if err != nil {
		return fmt.Errorf("failed to save quality metrics: %w", err)
	}

	return nil
}

// LogAction logs an audit action
func (r *Repository) LogAction(ctx context.Context, log *models.TranslationAuditLog) error {
	query := `INSERT INTO translation_audit_log 
	          (user_id, action, entity_type, entity_id, old_value, new_value, ip_address, user_agent)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	          RETURNING id, created_at`

	err := r.db.QueryRowContext(
		ctx, query,
		log.UserID, log.Action, log.EntityType, log.EntityID,
		log.OldValue, log.NewValue, log.IPAddress, log.UserAgent,
	).Scan(&log.ID, &log.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to log action: %w", err)
	}

	return nil
}

// GetRecentLogs retrieves recent audit logs
func (r *Repository) GetRecentLogs(ctx context.Context, limit int) ([]models.TranslationAuditLog, error) {
	query := `SELECT id, user_id, action, entity_type, entity_id, old_value, new_value, 
	          ip_address, user_agent, created_at
	          FROM translation_audit_log
	          ORDER BY created_at DESC
	          LIMIT $1`

	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query audit logs: %w", err)
	}
	defer rows.Close()

	var logs []models.TranslationAuditLog
	for rows.Next() {
		var l models.TranslationAuditLog
		var userID sql.NullInt64
		var entityType, oldValue, newValue, ipAddress, userAgent sql.NullString
		var entityID sql.NullInt64

		err := rows.Scan(
			&l.ID, &userID, &l.Action, &entityType, &entityID,
			&oldValue, &newValue, &ipAddress, &userAgent, &l.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan audit log: %w", err)
		}

		if userID.Valid {
			val := int(userID.Int64)
			l.UserID = &val
		}
		if entityType.Valid {
			l.EntityType = &entityType.String
		}
		if entityID.Valid {
			val := int(entityID.Int64)
			l.EntityID = &val
		}
		if oldValue.Valid {
			l.OldValue = &oldValue.String
		}
		if newValue.Valid {
			l.NewValue = &newValue.String
		}
		if ipAddress.Valid {
			l.IPAddress = &ipAddress.String
		}
		if userAgent.Valid {
			l.UserAgent = &userAgent.String
		}

		logs = append(logs, l)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return logs, nil
}

// GetTranslationByID is a helper method to get a single translation by ID
func (r *Repository) GetTranslationByID(ctx context.Context, id int) (*models.Translation, error) {
	query := `SELECT id, entity_type, entity_id, language, field_name, 
	          translated_text, is_machine_translated, is_verified, 
	          created_at, updated_at, metadata, COALESCE(version, 1) as version
	          FROM translations WHERE id = $1`

	var t models.Translation
	var metadata sql.NullString
	var version sql.NullInt64

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&t.ID, &t.EntityType, &t.EntityID, &t.Language, &t.FieldName,
		&t.TranslatedText, &t.IsMachineTranslated, &t.IsVerified,
		&t.CreatedAt, &t.UpdatedAt, &metadata, &version,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("translation not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get translation: %w", err)
	}

	if metadata.Valid {
		if err := json.Unmarshal([]byte(metadata.String), &t.Metadata); err != nil {
			r.logger.Error().Err(err).Msg("Failed to unmarshal metadata")
		}
	}

	if version.Valid {
		t.Version = int(version.Int64)
	} else {
		t.Version = 1
	}

	return &t, nil
}
