package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"

	"backend/internal/domain/models"
)

// CreateComponent создает новый компонент
func (db *Database) CreateComponent(ctx context.Context, component *models.CustomUIComponent) (int, error) {
	var id int
	err := db.pool.QueryRow(ctx, `
		INSERT INTO custom_ui_components (
			name, component_type, description, 
			template_code, styles, props_schema,
			is_active, created_by, updated_by
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id`,
		component.Name, component.ComponentType, component.Description,
		component.TemplateCode, component.Styles, component.PropsSchema,
		component.IsActive, component.CreatedBy, component.UpdatedBy,
	).Scan(&id)

	return id, err
}

// GetComponent получает компонент по ID
func (db *Database) GetComponent(ctx context.Context, id int) (*models.CustomUIComponent, error) {
	component := &models.CustomUIComponent{}
	err := db.pool.QueryRow(ctx, `
		SELECT id, name, component_type, description, 
		       template_code, styles, props_schema, is_active,
		       created_at, updated_at, created_by, updated_by
		FROM custom_ui_components WHERE id = $1`, id,
	).Scan(
		&component.ID, &component.Name, &component.ComponentType,
		&component.Description, &component.TemplateCode, &component.Styles,
		&component.PropsSchema, &component.IsActive,
		&component.CreatedAt, &component.UpdatedAt,
		&component.CreatedBy, &component.UpdatedBy,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("component not found")
	}

	return component, err
}

// UpdateComponent обновляет компонент
func (db *Database) UpdateComponent(ctx context.Context, id int, updates map[string]interface{}) error {
	var setClauses []string
	var args []interface{}
	argCounter := 1

	for key, value := range updates {
		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", key, argCounter))
		args = append(args, value)
		argCounter++
	}

	setClauses = append(setClauses, "updated_at = CURRENT_TIMESTAMP")

	query := fmt.Sprintf(`
		UPDATE custom_ui_components 
		SET %s
		WHERE id = $%d`,
		strings.Join(setClauses, ", "),
		argCounter,
	)

	args = append(args, id)

	_, err := db.pool.Exec(ctx, query, args...)
	return err
}

// DeleteComponent удаляет компонент
func (db *Database) DeleteComponent(ctx context.Context, id int) error {
	_, err := db.pool.Exec(ctx,
		"DELETE FROM custom_ui_components WHERE id = $1", id)
	return err
}

// ListComponents получает список компонентов с фильтрацией
func (db *Database) ListComponents(ctx context.Context, filters map[string]interface{}) ([]*models.CustomUIComponent, error) {
	log.Println("ListComponents called with filters:", filters)

	query := `
		SELECT id, name, component_type, description, 
		       template_code, styles, props_schema, is_active,
		       created_at, updated_at, created_by, updated_by
		FROM custom_ui_components WHERE 1=1`

	args := []interface{}{}
	argNum := 1

	if componentType, ok := filters["component_type"].(string); ok && componentType != "" {
		query += fmt.Sprintf(" AND component_type = $%d", argNum)
		args = append(args, componentType)
		argNum++
	}

	if active, ok := filters["active"].(string); ok && active != "" {
		if active == "true" {
			query += fmt.Sprintf(" AND is_active = $%d", argNum)
			args = append(args, true)
			argNum++
		} else if active == "false" {
			query += fmt.Sprintf(" AND is_active = $%d", argNum)
			args = append(args, false)
			argNum++
		}
	}

	query += " ORDER BY created_at DESC"

	log.Printf("Executing query: %s with args: %v", query, args)
	rows, err := db.pool.Query(ctx, query, args...)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return nil, err
	}
	defer rows.Close()

	components := []*models.CustomUIComponent{}
	for rows.Next() {
		component := &models.CustomUIComponent{}
		err := rows.Scan(
			&component.ID, &component.Name, &component.ComponentType,
			&component.Description, &component.TemplateCode, &component.Styles,
			&component.PropsSchema, &component.IsActive,
			&component.CreatedAt, &component.UpdatedAt,
			&component.CreatedBy, &component.UpdatedBy,
		)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			return nil, err
		}
		components = append(components, component)
	}

	log.Printf("Found %d components", len(components))
	return components, nil
}

// AddComponentUsage добавляет использование компонента для категории
func (db *Database) AddComponentUsage(ctx context.Context, usage *models.CustomUIComponentUsage) (int, error) {
	var id int
	err := db.pool.QueryRow(ctx, `
		INSERT INTO custom_ui_component_usage (
			component_id, category_id, usage_context, placement,
			priority, configuration, conditions_logic, is_active,
			created_by, updated_by
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id`,
		usage.ComponentID, usage.CategoryID, usage.UsageContext,
		usage.Placement, usage.Priority, usage.Configuration,
		usage.ConditionsLogic, usage.IsActive,
		usage.CreatedBy, usage.UpdatedBy,
	).Scan(&id)

	return id, err
}

// GetComponentUsages получает все использования компонентов с фильтрацией
func (db *Database) GetComponentUsages(ctx context.Context, componentID, categoryID *int) ([]*models.CustomUIComponentUsage, error) {
	query := `
		SELECT 
			ucu.id, ucu.component_id, ucu.category_id, ucu.usage_context,
			ucu.placement, ucu.priority, ucu.configuration, ucu.conditions_logic,
			ucu.is_active, ucu.created_at, ucu.updated_at, ucu.created_by, ucu.updated_by,
			cuc.name as component_name,
			mc.name as category_name
		FROM custom_ui_component_usage ucu
		JOIN custom_ui_components cuc ON ucu.component_id = cuc.id
		JOIN marketplace_categories mc ON ucu.category_id = mc.id
		WHERE 1=1`

	args := []interface{}{}
	argNum := 1

	if componentID != nil {
		query += fmt.Sprintf(" AND ucu.component_id = $%d", argNum)
		args = append(args, *componentID)
		argNum++
	}

	if categoryID != nil {
		query += fmt.Sprintf(" AND ucu.category_id = $%d", argNum)
		args = append(args, *categoryID)
		argNum++
	}

	query += " ORDER BY ucu.created_at DESC"

	rows, err := db.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	usages := []*models.CustomUIComponentUsage{}
	for rows.Next() {
		usage := &models.CustomUIComponentUsage{}
		err := rows.Scan(
			&usage.ID, &usage.ComponentID, &usage.CategoryID, &usage.UsageContext,
			&usage.Placement, &usage.Priority, &usage.Configuration, &usage.ConditionsLogic,
			&usage.IsActive, &usage.CreatedAt, &usage.UpdatedAt, &usage.CreatedBy, &usage.UpdatedBy,
			&usage.ComponentName, &usage.CategoryName,
		)
		if err != nil {
			return nil, err
		}
		usages = append(usages, usage)
	}

	return usages, nil
}

// RemoveComponentUsage удаляет использование компонента
func (db *Database) RemoveComponentUsage(ctx context.Context, id int) error {
	_, err := db.pool.Exec(ctx,
		"DELETE FROM custom_ui_component_usage WHERE id = $1", id)
	return err
}

// GetCategoryComponents получает компоненты для категории
func (db *Database) GetCategoryComponents(ctx context.Context, categoryID int, usageContext string) ([]*models.CustomUIComponentUsage, error) {
	query := `
		SELECT 
			ucu.id, ucu.component_id, ucu.category_id, ucu.usage_context,
			ucu.placement, ucu.priority, ucu.configuration, ucu.conditions_logic,
			ucu.is_active, ucu.created_at, ucu.updated_at,
			cuc.name, cuc.component_type, cuc.template_code, cuc.styles, cuc.props_schema
		FROM custom_ui_component_usage ucu
		JOIN custom_ui_components cuc ON ucu.component_id = cuc.id
		WHERE ucu.category_id = $1 AND ucu.is_active = true`

	args := []interface{}{categoryID}

	if usageContext != "" {
		query += " AND ucu.usage_context = $2"
		args = append(args, usageContext)
	}

	query += " ORDER BY ucu.priority ASC"

	rows, err := db.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	components := []*models.CustomUIComponentUsage{}
	for rows.Next() {
		usage := &models.CustomUIComponentUsage{}
		component := &models.CustomUIComponent{}

		err := rows.Scan(
			&usage.ID, &usage.ComponentID, &usage.CategoryID, &usage.UsageContext,
			&usage.Placement, &usage.Priority, &usage.Configuration, &usage.ConditionsLogic,
			&usage.IsActive, &usage.CreatedAt, &usage.UpdatedAt,
			&component.Name, &component.ComponentType, &component.TemplateCode,
			&component.Styles, &component.PropsSchema,
		)
		if err != nil {
			return nil, err
		}

		usage.Component = component
		components = append(components, usage)
	}

	return components, nil
}

// CreateTemplate создает шаблон компонента
func (db *Database) CreateTemplate(ctx context.Context, template *models.ComponentTemplate) (int, error) {
	var id int
	err := db.pool.QueryRow(ctx, `
		INSERT INTO component_templates (
			component_id, name, description, template_config,
			preview_image, category_id, created_by
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id`,
		template.ComponentID, template.Name, template.Description,
		template.TemplateConfig, template.PreviewImage, template.CategoryID,
		template.CreatedBy,
	).Scan(&id)

	return id, err
}

// ListTemplates получает список шаблонов компонента
func (db *Database) ListTemplates(ctx context.Context, componentID int) ([]*models.ComponentTemplate, error) {
	query := `
		SELECT id, component_id, name, description, template_config,
		       preview_image, category_id, created_at, created_by
		FROM component_templates`

	args := []interface{}{}
	if componentID > 0 {
		query += " WHERE component_id = $1"
		args = append(args, componentID)
	}

	query += " ORDER BY created_at DESC"

	rows, err := db.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	templates := []*models.ComponentTemplate{}
	for rows.Next() {
		template := &models.ComponentTemplate{}
		err := rows.Scan(
			&template.ID, &template.ComponentID, &template.Name,
			&template.Description, &template.TemplateConfig,
			&template.PreviewImage, &template.CategoryID,
			&template.CreatedAt, &template.CreatedBy,
		)
		if err != nil {
			return nil, err
		}
		templates = append(templates, template)
	}

	return templates, nil
}
