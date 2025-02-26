package postgres

import (
	"backend/internal/domain/models"
	"context"
	"time"
)

// CreateStorefront создает новую витрину
func (db *Database) CreateStorefront(ctx context.Context, storefront *models.Storefront) (int, error) {
	var id int
	err := db.pool.QueryRow(ctx, `
		INSERT INTO user_storefronts (
			user_id, name, description, logo_path, slug,
			status, creation_transaction_id, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id
	`,
		storefront.UserID, storefront.Name, storefront.Description, storefront.LogoPath,
		storefront.Slug, storefront.Status, storefront.CreationTransactionID,
		storefront.CreatedAt, storefront.UpdatedAt,
	).Scan(&id)

	return id, err
}

// GetUserStorefronts возвращает все витрины пользователя
func (db *Database) GetUserStorefronts(ctx context.Context, userID int) ([]models.Storefront, error) {
	rows, err := db.pool.Query(ctx, `
		SELECT id, user_id, name, description, logo_path, slug,
		       status, creation_transaction_id, created_at, updated_at
		FROM user_storefronts
		WHERE user_id = $1
		ORDER BY created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var storefronts []models.Storefront
	for rows.Next() {
		var s models.Storefront
		err := rows.Scan(
			&s.ID, &s.UserID, &s.Name, &s.Description, &s.LogoPath,
			&s.Slug, &s.Status, &s.CreationTransactionID, &s.CreatedAt, &s.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		storefronts = append(storefronts, s)
	}

	return storefronts, rows.Err()
}

// GetStorefrontByID возвращает витрину по ID
func (db *Database) GetStorefrontByID(ctx context.Context, id int) (*models.Storefront, error) {
	var s models.Storefront
	err := db.pool.QueryRow(ctx, `
		SELECT id, user_id, name, description, logo_path, slug,
		       status, creation_transaction_id, created_at, updated_at
		FROM user_storefronts
		WHERE id = $1
	`, id).Scan(
		&s.ID, &s.UserID, &s.Name, &s.Description, &s.LogoPath,
		&s.Slug, &s.Status, &s.CreationTransactionID, &s.CreatedAt, &s.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &s, nil
}

// UpdateStorefront обновляет информацию о витрине
func (db *Database) UpdateStorefront(ctx context.Context, storefront *models.Storefront) error {
	_, err := db.pool.Exec(ctx, `
		UPDATE user_storefronts
		SET name = $1, description = $2, logo_path = $3, slug = $4,
		    status = $5, updated_at = $6
		WHERE id = $7
	`,
		storefront.Name, storefront.Description, storefront.LogoPath,
		storefront.Slug, storefront.Status, time.Now(), storefront.ID,
	)
	return err
}

// DeleteStorefront удаляет витрину
func (db *Database) DeleteStorefront(ctx context.Context, id int) error {
	_, err := db.pool.Exec(ctx, `
		DELETE FROM user_storefronts WHERE id = $1
	`, id)
	return err
}

// CreateImportSource создает новый источник импорта
func (db *Database) CreateImportSource(ctx context.Context, source *models.ImportSource) (int, error) {
	var id int
	err := db.pool.QueryRow(ctx, `
		INSERT INTO import_sources (
			storefront_id, type, url, auth_data, schedule, mapping,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`,
		source.StorefrontID, source.Type, source.URL, source.AuthData,
		source.Schedule, source.Mapping, source.CreatedAt, source.UpdatedAt,
	).Scan(&id)

	return id, err
}

// GetImportSourceByID возвращает источник импорта по ID
func (db *Database) GetImportSourceByID(ctx context.Context, id int) (*models.ImportSource, error) {
	var s models.ImportSource
	err := db.pool.QueryRow(ctx, `
		SELECT id, storefront_id, type, url, auth_data, schedule, mapping,
		       last_import_at, last_import_status, last_import_log,
		       created_at, updated_at
		FROM import_sources
		WHERE id = $1
	`, id).Scan(
		&s.ID, &s.StorefrontID, &s.Type, &s.URL, &s.AuthData,
		&s.Schedule, &s.Mapping, &s.LastImportAt, &s.LastImportStatus,
		&s.LastImportLog, &s.CreatedAt, &s.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &s, nil
}

// GetImportSources возвращает источники импорта для витрины
func (db *Database) GetImportSources(ctx context.Context, storefrontID int) ([]models.ImportSource, error) {
	rows, err := db.pool.Query(ctx, `
		SELECT id, storefront_id, type, url, auth_data, schedule, mapping,
		       last_import_at, last_import_status, last_import_log,
		       created_at, updated_at
		FROM import_sources
		WHERE storefront_id = $1
		ORDER BY created_at DESC
	`, storefrontID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sources []models.ImportSource
	for rows.Next() {
		var s models.ImportSource
		err := rows.Scan(
			&s.ID, &s.StorefrontID, &s.Type, &s.URL, &s.AuthData,
			&s.Schedule, &s.Mapping, &s.LastImportAt, &s.LastImportStatus,
			&s.LastImportLog, &s.CreatedAt, &s.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		sources = append(sources, s)
	}

	return sources, rows.Err()
}

// UpdateImportSource обновляет источник импорта
func (db *Database) UpdateImportSource(ctx context.Context, source *models.ImportSource) error {
	_, err := db.pool.Exec(ctx, `
		UPDATE import_sources
		SET type = $1, url = $2, auth_data = $3, schedule = $4, mapping = $5,
		    last_import_at = $6, last_import_status = $7, last_import_log = $8,
		    updated_at = $9
		WHERE id = $10
	`,
		source.Type, source.URL, source.AuthData, source.Schedule, source.Mapping,
		source.LastImportAt, source.LastImportStatus, source.LastImportLog,
		time.Now(), source.ID,
	)
	return err
}

// DeleteImportSource удаляет источник импорта
func (db *Database) DeleteImportSource(ctx context.Context, id int) error {
	_, err := db.pool.Exec(ctx, `
		DELETE FROM import_sources WHERE id = $1
	`, id)
	return err
}

// CreateImportHistory создает новую запись в истории импорта
func (db *Database) CreateImportHistory(ctx context.Context, history *models.ImportHistory) (int, error) {
	var id int
	err := db.pool.QueryRow(ctx, `
		INSERT INTO import_history (
			source_id, status, items_total, items_imported, items_failed,
			log, started_at, finished_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`,
		history.SourceID, history.Status, history.ItemsTotal, history.ItemsImported,
		history.ItemsFailed, history.Log, history.StartedAt, history.FinishedAt,
	).Scan(&id)

	return id, err
}

// GetImportHistory возвращает историю импорта
func (db *Database) GetImportHistory(ctx context.Context, sourceID int, limit, offset int) ([]models.ImportHistory, error) {
	rows, err := db.pool.Query(ctx, `
		SELECT id, source_id, status, items_total, items_imported, items_failed,
		       log, started_at, finished_at
		FROM import_history
		WHERE source_id = $1
		ORDER BY started_at DESC
		LIMIT $2 OFFSET $3
	`, sourceID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []models.ImportHistory
	for rows.Next() {
		var h models.ImportHistory
		err := rows.Scan(
			&h.ID, &h.SourceID, &h.Status, &h.ItemsTotal, &h.ItemsImported,
			&h.ItemsFailed, &h.Log, &h.StartedAt, &h.FinishedAt,
		)
		if err != nil {
			return nil, err
		}
		history = append(history, h)
	}

	return history, rows.Err()
}

// UpdateImportHistory обновляет запись в истории импорта
func (db *Database) UpdateImportHistory(ctx context.Context, history *models.ImportHistory) error {
	_, err := db.pool.Exec(ctx, `
		UPDATE import_history
		SET status = $1, items_total = $2, items_imported = $3, items_failed = $4,
		    log = $5, finished_at = $6
		WHERE id = $7
	`,
		history.Status, history.ItemsTotal, history.ItemsImported, history.ItemsFailed,
		history.Log, history.FinishedAt, history.ID,
	)
	return err
}