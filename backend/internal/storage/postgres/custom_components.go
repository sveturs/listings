package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"backend/internal/domain/models"
)

// CustomComponentStorage определяет методы для работы с кастомными компонентами
type CustomComponentStorage interface {
	CreateComponent(ctx context.Context, component *models.CustomUIComponent) (int, error)
	GetComponent(ctx context.Context, id int) (*models.CustomUIComponent, error)
	GetComponentByName(ctx context.Context, name string) (*models.CustomUIComponent, error)
	ListComponents(ctx context.Context, filters map[string]interface{}) ([]*models.CustomUIComponent, error)
	UpdateComponent(ctx context.Context, id int, updates map[string]interface{}) error
	DeleteComponent(ctx context.Context, id int) error
	
	// Использование компонентов
	AddComponentUsage(ctx context.Context, usage *models.CustomUIComponentUsage) (int, error)
	RemoveComponentUsage(ctx context.Context, id int) error
	GetCategoryComponents(ctx context.Context, categoryID int, context string) ([]*models.CustomUIComponentUsage, error)
	CreateComponentUsage(ctx context.Context, usage *models.CustomUIComponentUsage) (int, error)
	GetComponentUsage(ctx context.Context, entityType string, entityID int) (*models.CustomUIComponentUsage, error)
	ListComponentUsages(ctx context.Context, componentID int) ([]*models.CustomUIComponentUsage, error)
	UpdateComponentUsage(ctx context.Context, id int, updates map[string]interface{}) error
	DeleteComponentUsage(ctx context.Context, id int) error
	
	// Шаблоны
	CreateTemplate(ctx context.Context, template *models.CustomUITemplate) (int, error)
	ListTemplates(ctx context.Context) ([]*models.CustomUITemplate, error)
	GetTemplate(ctx context.Context, id int) (*models.CustomUITemplate, error)
}

// customComponentStorage реализация интерфейса
type customComponentStorage struct {
	pool *pgxpool.Pool
}

// NewCustomComponentStorage создает новое хранилище для кастомных компонентов
func NewCustomComponentStorage(pool *pgxpool.Pool) CustomComponentStorage {
	return &customComponentStorage{pool: pool}
}

// CreateComponent создает новый компонент
func (s *customComponentStorage) CreateComponent(ctx context.Context, component *models.CustomUIComponent) (int, error) {
	query := `
		INSERT INTO custom_ui_components (
			name, display_name, description, component_type, component_code,
			configuration, dependencies, is_active, created_by, updated_by
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id`

	var id int
	err := s.pool.QueryRow(ctx, query,
		component.Name,
		component.DisplayName,
		component.Description,
		component.ComponentType,
		component.ComponentCode,
		component.Configuration,
		component.Dependencies,
		component.IsActive,
		component.CreatedBy,
		component.UpdatedBy,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("ошибка создания компонента: %w", err)
	}

	return id, nil
}

// GetComponent получает компонент по ID
func (s *customComponentStorage) GetComponent(ctx context.Context, id int) (*models.CustomUIComponent, error) {
	query := `
		SELECT id, name, display_name, description, component_type, component_code,
			   configuration, dependencies, is_active, created_at, updated_at,
			   created_by, updated_by, compiled_code, compilation_errors, last_compiled_at
		FROM custom_ui_components
		WHERE id = $1`

	component := &models.CustomUIComponent{}
	err := s.pool.QueryRow(ctx, query, id).Scan(
		&component.ID,
		&component.Name,
		&component.DisplayName,
		&component.Description,
		&component.ComponentType,
		&component.ComponentCode,
		&component.Configuration,
		&component.Dependencies,
		&component.IsActive,
		&component.CreatedAt,
		&component.UpdatedAt,
		&component.CreatedBy,
		&component.UpdatedBy,
		&component.CompiledCode,
		&component.CompilationErrors,
		&component.LastCompiledAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("компонент не найден")
	}
	if err != nil {
		return nil, fmt.Errorf("ошибка получения компонента: %w", err)
	}

	return component, nil
}

// GetComponentByName получает компонент по имени
func (s *customComponentStorage) GetComponentByName(ctx context.Context, name string) (*models.CustomUIComponent, error) {
	query := `
		SELECT id, name, display_name, description, component_type, component_code,
			   configuration, dependencies, is_active, created_at, updated_at,
			   created_by, updated_by, compiled_code, compilation_errors, last_compiled_at
		FROM custom_ui_components
		WHERE name = $1`

	component := &models.CustomUIComponent{}
	err := s.pool.QueryRow(ctx, query, name).Scan(
		&component.ID,
		&component.Name,
		&component.DisplayName,
		&component.Description,
		&component.ComponentType,
		&component.ComponentCode,
		&component.Configuration,
		&component.Dependencies,
		&component.IsActive,
		&component.CreatedAt,
		&component.UpdatedAt,
		&component.CreatedBy,
		&component.UpdatedBy,
		&component.CompiledCode,
		&component.CompilationErrors,
		&component.LastCompiledAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("компонент не найден")
	}
	if err != nil {
		return nil, fmt.Errorf("ошибка получения компонента: %w", err)
	}

	return component, nil
}

// ListComponents получает список компонентов с фильтрами
func (s *customComponentStorage) ListComponents(ctx context.Context, filters map[string]interface{}) ([]*models.CustomUIComponent, error) {
	query := `
		SELECT id, name, display_name, description, component_type, component_code,
			   configuration, dependencies, is_active, created_at, updated_at,
			   created_by, updated_by
		FROM custom_ui_components
		WHERE 1=1`

	args := []interface{}{}
	argPos := 1

	if componentType, ok := filters["component_type"].(string); ok && componentType != "" {
		query += fmt.Sprintf(" AND component_type = $%d", argPos)
		args = append(args, componentType)
		argPos++
	}

	if isActive, ok := filters["is_active"].(bool); ok {
		query += fmt.Sprintf(" AND is_active = $%d", argPos)
		args = append(args, isActive)
		argPos++
	}

	query += " ORDER BY created_at DESC"

	rows, err := s.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения списка компонентов: %w", err)
	}
	defer rows.Close()

	components := []*models.CustomUIComponent{}
	for rows.Next() {
		component := &models.CustomUIComponent{}
		err := rows.Scan(
			&component.ID,
			&component.Name,
			&component.DisplayName,
			&component.Description,
			&component.ComponentType,
			&component.ComponentCode,
			&component.Configuration,
			&component.Dependencies,
			&component.IsActive,
			&component.CreatedAt,
			&component.UpdatedAt,
			&component.CreatedBy,
			&component.UpdatedBy,
		)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования компонента: %w", err)
		}
		components = append(components, component)
	}

	return components, nil
}

// UpdateComponent обновляет компонент
func (s *customComponentStorage) UpdateComponent(ctx context.Context, id int, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return nil
	}

	setClause := ""
	args := []interface{}{}
	argPos := 1

	for key, value := range updates {
		if setClause != "" {
			setClause += ", "
		}
		setClause += fmt.Sprintf("%s = $%d", key, argPos)
		args = append(args, value)
		argPos++
	}

	args = append(args, id)
	query := fmt.Sprintf(`UPDATE custom_ui_components SET %s WHERE id = $%d`, setClause, argPos)

	_, err := s.pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("ошибка обновления компонента: %w", err)
	}

	return nil
}

// DeleteComponent удаляет компонент
func (s *customComponentStorage) DeleteComponent(ctx context.Context, id int) error {
	query := `DELETE FROM custom_ui_components WHERE id = $1`
	
	result, err := s.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("ошибка удаления компонента: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("компонент не найден")
	}

	return nil
}

// CreateComponentUsage создает запись использования компонента
func (s *customComponentStorage) CreateComponentUsage(ctx context.Context, usage *models.CustomUIComponentUsage) (int, error) {
	query := `
		INSERT INTO custom_ui_component_usage (
			component_id, category_id, usage_context, placement, priority,
			configuration, conditions_logic, is_active, created_by, updated_by
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id`

	var id int
	err := s.pool.QueryRow(ctx, query,
		usage.ComponentID,
		usage.CategoryID,
		usage.UsageContext,
		usage.Placement,
		usage.Priority,
		usage.Configuration,
		usage.ConditionsLogic,
		usage.IsActive,
		usage.CreatedBy,
		usage.UpdatedBy,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("ошибка создания использования компонента: %w", err)
	}

	return id, nil
}

// GetComponentUsage получает использование компонента для сущности
func (s *customComponentStorage) GetComponentUsage(ctx context.Context, entityType string, entityID int) (*models.CustomUIComponentUsage, error) {
	// Этот метод устарел, но оставлен для совместимости
	return nil, fmt.Errorf("метод устарел, используйте GetCategoryComponents")
}

// GetCategoryComponents получает компоненты для категории
func (s *customComponentStorage) GetCategoryComponents(ctx context.Context, categoryID int, context string) ([]*models.CustomUIComponentUsage, error) {
	query := `
		SELECT u.id, u.component_id, u.category_id, u.usage_context, u.placement,
			   u.priority, u.configuration, u.conditions_logic, u.is_active,
			   u.created_at, u.updated_at, u.created_by, u.updated_by,
			   c.id, c.name, c.display_name, c.description, c.component_type,
			   c.component_code, c.configuration, c.dependencies, c.is_active
		FROM custom_ui_component_usage u
		JOIN custom_ui_components c ON u.component_id = c.id
		WHERE u.category_id = $1`
	args := []interface{}{categoryID}

	if context != "" {
		query += " AND u.usage_context = $2"
		args = append(args, context)
	}

	query += " ORDER BY u.priority ASC"

	rows, err := s.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения компонентов категории: %w", err)
	}
	defer rows.Close()

	var usages []*models.CustomUIComponentUsage
	for rows.Next() {
		usage := &models.CustomUIComponentUsage{Component: &models.CustomUIComponent{}}
		err := rows.Scan(
			&usage.ID,
			&usage.ComponentID,
			&usage.CategoryID,
			&usage.UsageContext,
			&usage.Placement,
			&usage.Priority,
			&usage.Configuration,
			&usage.ConditionsLogic,
			&usage.IsActive,
			&usage.CreatedAt,
			&usage.UpdatedAt,
			&usage.CreatedBy,
			&usage.UpdatedBy,
			&usage.Component.ID,
			&usage.Component.Name,
			&usage.Component.DisplayName,
			&usage.Component.Description,
			&usage.Component.ComponentType,
			&usage.Component.ComponentCode,
			&usage.Component.Configuration,
			&usage.Component.Dependencies,
			&usage.Component.IsActive,
		)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования использования компонента: %w", err)
		}
		usages = append(usages, usage)
	}

	return usages, nil
}

// ListComponentUsages получает список использований компонента
func (s *customComponentStorage) ListComponentUsages(ctx context.Context, componentID int) ([]*models.CustomUIComponentUsage, error) {
	query := `
		SELECT id, component_id, category_id, usage_context, placement,
			   priority, configuration, conditions_logic, is_active,
			   created_at, updated_at, created_by, updated_by
		FROM custom_ui_component_usage
		WHERE component_id = $1
		ORDER BY priority ASC, created_at DESC`

	rows, err := s.pool.Query(ctx, query, componentID)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения списка использований: %w", err)
	}
	defer rows.Close()

	usages := []*models.CustomUIComponentUsage{}
	for rows.Next() {
		usage := &models.CustomUIComponentUsage{}
		err := rows.Scan(
			&usage.ID,
			&usage.ComponentID,
			&usage.CategoryID,
			&usage.UsageContext,
			&usage.Placement,
			&usage.Priority,
			&usage.Configuration,
			&usage.ConditionsLogic,
			&usage.IsActive,
			&usage.CreatedAt,
			&usage.UpdatedAt,
			&usage.CreatedBy,
			&usage.UpdatedBy,
		)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования использования: %w", err)
		}
		usages = append(usages, usage)
	}

	return usages, nil
}

// UpdateComponentUsage обновляет использование компонента
func (s *customComponentStorage) UpdateComponentUsage(ctx context.Context, id int, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return nil
	}

	setClause := ""
	args := []interface{}{}
	argPos := 1

	for key, value := range updates {
		if setClause != "" {
			setClause += ", "
		}
		setClause += fmt.Sprintf("%s = $%d", key, argPos)
		args = append(args, value)
		argPos++
	}

	args = append(args, id)
	query := fmt.Sprintf(`UPDATE custom_ui_component_usage SET %s WHERE id = $%d`, setClause, argPos)

	_, err := s.pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("ошибка обновления использования: %w", err)
	}

	return nil
}

// DeleteComponentUsage удаляет использование компонента
func (s *customComponentStorage) DeleteComponentUsage(ctx context.Context, id int) error {
	query := `DELETE FROM custom_ui_component_usage WHERE id = $1`
	
	result, err := s.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("ошибка удаления использования: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("использование не найдено")
	}

	return nil
}

// ListTemplates получает список шаблонов
func (s *customComponentStorage) ListTemplates(ctx context.Context) ([]*models.CustomUITemplate, error) {
	query := `
		SELECT id, name, display_name, description, template_code,
			   COALESCE(variables, '{}'), is_shared, created_at, updated_at,
			   created_by, updated_by
		FROM custom_ui_templates
		ORDER BY display_name`

	rows, err := s.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения списка шаблонов: %w", err)
	}
	defer rows.Close()

	templates := []*models.CustomUITemplate{}
	for rows.Next() {
		template := &models.CustomUITemplate{}
		err := rows.Scan(
			&template.ID,
			&template.Name,
			&template.DisplayName,
			&template.Description,
			&template.TemplateCode,
			&template.Variables,
			&template.IsShared,
			&template.CreatedAt,
			&template.UpdatedAt,
			&template.CreatedBy,
			&template.UpdatedBy,
		)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования шаблона: %w", err)
		}
		templates = append(templates, template)
	}

	return templates, nil
}

// GetTemplate получает шаблон по ID
func (s *customComponentStorage) GetTemplate(ctx context.Context, id int) (*models.CustomUITemplate, error) {
	query := `
		SELECT id, name, display_name, description, template_code,
			   variables, is_shared, created_at, updated_at,
			   created_by, updated_by
		FROM custom_ui_templates
		WHERE id = $1`

	template := &models.CustomUITemplate{}
	err := s.pool.QueryRow(ctx, query, id).Scan(
		&template.ID,
		&template.Name,
		&template.DisplayName,
		&template.Description,
		&template.TemplateCode,
		&template.Variables,
		&template.IsShared,
		&template.CreatedAt,
		&template.UpdatedAt,
		&template.CreatedBy,
		&template.UpdatedBy,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("шаблон не найден")
	}
	if err != nil {
		return nil, fmt.Errorf("ошибка получения шаблона: %w", err)
	}

	return template, nil
}

// CreateTemplate создает новый шаблон
func (s *customComponentStorage) CreateTemplate(ctx context.Context, template *models.CustomUITemplate) (int, error) {
	query := `
		INSERT INTO custom_ui_templates (
			name, display_name, description, template_code,
			variables, is_shared, created_by, updated_by
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id`

	var id int
	err := s.pool.QueryRow(ctx, query,
		template.Name,
		template.DisplayName,
		template.Description,
		template.TemplateCode,
		template.Variables,
		template.IsShared,
		template.CreatedBy,
		template.UpdatedBy,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("ошибка создания шаблона: %w", err)
	}

	return id, nil
}

// AddComponentUsage - alias для CreateComponentUsage для совместимости с интерфейсом
func (s *customComponentStorage) AddComponentUsage(ctx context.Context, usage *models.CustomUIComponentUsage) (int, error) {
	return s.CreateComponentUsage(ctx, usage)
}

// RemoveComponentUsage удаляет использование компонента
func (s *customComponentStorage) RemoveComponentUsage(ctx context.Context, id int) error {
	return s.DeleteComponentUsage(ctx, id)
}