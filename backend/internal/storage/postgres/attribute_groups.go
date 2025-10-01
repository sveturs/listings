package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"backend/internal/domain/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// AttributeGroupStorage интерфейс для работы с группами атрибутов
type AttributeGroupStorage interface {
	// Группы атрибутов
	CreateAttributeGroup(ctx context.Context, group *models.AttributeGroup) (int, error)
	GetAttributeGroup(ctx context.Context, id int) (*models.AttributeGroup, error)
	GetAttributeGroupByName(ctx context.Context, name string) (*models.AttributeGroup, error)
	ListAttributeGroups(ctx context.Context) ([]*models.AttributeGroup, error)
	UpdateAttributeGroup(ctx context.Context, id int, updates map[string]interface{}) error
	DeleteAttributeGroup(ctx context.Context, id int) error

	// Элементы групп
	AddItemToGroup(ctx context.Context, groupID int, item *models.AttributeGroupItem) (int, error)
	RemoveItemFromGroup(ctx context.Context, groupID, attributeID int) error
	GetGroupItems(ctx context.Context, groupID int) ([]*models.AttributeGroupItem, error)
	UpdateGroupItem(ctx context.Context, id int, updates map[string]interface{}) error
	GetAttributeGroupWithItems(ctx context.Context, id int) (*models.AttributeGroupWithItems, error)

	// Привязка к категориям
	AttachGroupToCategory(ctx context.Context, categoryID int, group *models.CategoryAttributeGroup) (int, error)
	DetachGroupFromCategory(ctx context.Context, categoryID, groupID int) error
	GetCategoryGroups(ctx context.Context, categoryID int) ([]*models.CategoryAttributeGroup, error)
	UpdateCategoryGroup(ctx context.Context, id int, updates map[string]interface{}) error
}

type attributeGroupStorage struct {
	pool *pgxpool.Pool
}

// NewAttributeGroupStorage создает новое хранилище для групп атрибутов
func NewAttributeGroupStorage(pool *pgxpool.Pool) AttributeGroupStorage {
	return &attributeGroupStorage{pool: pool}
}

// CreateAttributeGroup создает новую группу атрибутов
func (s *attributeGroupStorage) CreateAttributeGroup(ctx context.Context, group *models.AttributeGroup) (int, error) {
	query := `
		INSERT INTO attribute_groups (code, name, display_name, description, icon, sort_order, is_active, is_system)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id`

	var id int
	err := s.pool.QueryRow(ctx, query,
		group.Code,
		group.Name,
		group.DisplayName,
		group.Description,
		group.Icon,
		group.SortOrder,
		group.IsActive,
		group.IsSystem,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("ошибка создания группы атрибутов: %w", err)
	}

	return id, nil
}

// GetAttributeGroup получает группу атрибутов по ID
func (s *attributeGroupStorage) GetAttributeGroup(ctx context.Context, id int) (*models.AttributeGroup, error) {
	query := `
		SELECT id, code, name, display_name, description, icon, sort_order, is_active, is_system, created_at, updated_at
		FROM attribute_groups
		WHERE id = $1`

	group := &models.AttributeGroup{}
	err := s.pool.QueryRow(ctx, query, id).Scan(
		&group.ID,
		&group.Code,
		&group.Name,
		&group.DisplayName,
		&group.Description,
		&group.Icon,
		&group.SortOrder,
		&group.IsActive,
		&group.IsSystem,
		&group.CreatedAt,
		&group.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("группа атрибутов не найдена")
	}
	if err != nil {
		return nil, fmt.Errorf("ошибка получения группы атрибутов: %w", err)
	}

	return group, nil
}

// GetAttributeGroupByName получает группу атрибутов по имени
func (s *attributeGroupStorage) GetAttributeGroupByName(ctx context.Context, name string) (*models.AttributeGroup, error) {
	query := `
		SELECT id, code, name, display_name, description, icon, sort_order, is_active, is_system, created_at, updated_at
		FROM attribute_groups
		WHERE name = $1`

	group := &models.AttributeGroup{}
	err := s.pool.QueryRow(ctx, query, name).Scan(
		&group.ID,
		&group.Code,
		&group.Name,
		&group.DisplayName,
		&group.Description,
		&group.Icon,
		&group.SortOrder,
		&group.IsActive,
		&group.IsSystem,
		&group.CreatedAt,
		&group.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("группа атрибутов не найдена")
	}
	if err != nil {
		return nil, fmt.Errorf("ошибка получения группы атрибутов: %w", err)
	}

	return group, nil
}

// ListAttributeGroups получает список всех групп атрибутов
func (s *attributeGroupStorage) ListAttributeGroups(ctx context.Context) ([]*models.AttributeGroup, error) {
	query := `
		SELECT id, code, name, display_name, description, icon, sort_order, is_active, is_system, created_at, updated_at
		FROM attribute_groups
		ORDER BY sort_order, name`

	rows, err := s.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения списка групп атрибутов: %w", err)
	}
	defer rows.Close()

	var groups []*models.AttributeGroup
	for rows.Next() {
		group := &models.AttributeGroup{}
		err := rows.Scan(
			&group.ID,
			&group.Code,
			&group.Name,
			&group.DisplayName,
			&group.Description,
			&group.Icon,
			&group.SortOrder,
			&group.IsActive,
			&group.IsSystem,
			&group.CreatedAt,
			&group.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования группы атрибутов: %w", err)
		}
		groups = append(groups, group)
	}

	return groups, nil
}

// UpdateAttributeGroup обновляет группу атрибутов
func (s *attributeGroupStorage) UpdateAttributeGroup(ctx context.Context, id int, updates map[string]interface{}) error {
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
	query := fmt.Sprintf(`UPDATE attribute_groups SET %s WHERE id = $%d`, setClause, argPos)

	_, err := s.pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("ошибка обновления группы атрибутов: %w", err)
	}

	return nil
}

// DeleteAttributeGroup удаляет группу атрибутов
func (s *attributeGroupStorage) DeleteAttributeGroup(ctx context.Context, id int) error {
	// Проверяем, что группа не является системной
	var isSystem bool
	err := s.pool.QueryRow(ctx, "SELECT is_system FROM attribute_groups WHERE id = $1", id).Scan(&isSystem)
	if err != nil {
		return fmt.Errorf("ошибка проверки группы: %w", err)
	}

	if isSystem {
		return fmt.Errorf("нельзя удалить системную группу")
	}

	query := `DELETE FROM attribute_groups WHERE id = $1`

	result, err := s.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("ошибка удаления группы атрибутов: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("группа атрибутов не найдена")
	}

	return nil
}

// AddItemToGroup добавляет атрибут в группу
func (s *attributeGroupStorage) AddItemToGroup(ctx context.Context, groupID int, item *models.AttributeGroupItem) (int, error) {
	query := `
		INSERT INTO attribute_group_items (group_id, attribute_id, icon, sort_order, custom_display_name, visibility_condition)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (group_id, attribute_id) DO UPDATE SET
			icon = EXCLUDED.icon,
			sort_order = EXCLUDED.sort_order,
			custom_display_name = EXCLUDED.custom_display_name,
			visibility_condition = EXCLUDED.visibility_condition
		RETURNING id`

	var id int
	err := s.pool.QueryRow(ctx, query,
		groupID,
		item.AttributeID,
		item.Icon,
		item.SortOrder,
		item.CustomDisplayName,
		item.VisibilityCondition,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("ошибка добавления атрибута в группу: %w", err)
	}

	return id, nil
}

// RemoveItemFromGroup удаляет атрибут из группы
func (s *attributeGroupStorage) RemoveItemFromGroup(ctx context.Context, groupID, attributeID int) error {
	query := `
		DELETE FROM attribute_group_items 
		WHERE group_id = $1 AND attribute_id = $2`

	result, err := s.pool.Exec(ctx, query, groupID, attributeID)
	if err != nil {
		return fmt.Errorf("ошибка удаления атрибута из группы: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("атрибут не найден в группе")
	}

	return nil
}

// GetGroupItems получает все атрибуты группы
func (s *attributeGroupStorage) GetGroupItems(ctx context.Context, groupID int) ([]*models.AttributeGroupItem, error) {
	query := `
		SELECT agi.id, agi.group_id, agi.attribute_id, agi.icon, agi.sort_order, 
		       agi.custom_display_name, agi.visibility_condition, agi.created_at,
		       ca.id, ca.name, ca.display_name, ca.attribute_type, ca.options, 
		       ca.validation_rules, ca.is_searchable, ca.is_filterable, ca.is_required,
		       ca.sort_order, ca.created_at
		FROM attribute_group_items agi
		LEFT JOIN category_attributes ca ON agi.attribute_id = ca.id
		WHERE agi.group_id = $1
		ORDER BY agi.sort_order, ca.display_name`

	rows, err := s.pool.Query(ctx, query, groupID)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения атрибутов группы: %w", err)
	}
	defer rows.Close()

	var items []*models.AttributeGroupItem
	for rows.Next() {
		item := &models.AttributeGroupItem{
			Attribute: &models.CategoryAttribute{},
		}
		var validRules sql.NullString
		err := rows.Scan(
			&item.ID,
			&item.GroupID,
			&item.AttributeID,
			&item.Icon,
			&item.SortOrder,
			&item.CustomDisplayName,
			&item.VisibilityCondition,
			&item.CreatedAt,
			&item.Attribute.ID,
			&item.Attribute.Name,
			&item.Attribute.DisplayName,
			&item.Attribute.AttributeType,
			&item.Attribute.Options,
			&validRules,
			&item.Attribute.IsSearchable,
			&item.Attribute.IsFilterable,
			&item.Attribute.IsRequired,
			&item.Attribute.SortOrder,
			&item.Attribute.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования атрибута группы: %w", err)
		}
		if validRules.Valid {
			item.Attribute.ValidRules = json.RawMessage(validRules.String)
		}
		items = append(items, item)
	}

	return items, nil
}

// UpdateGroupItem обновляет атрибут в группе
func (s *attributeGroupStorage) UpdateGroupItem(ctx context.Context, id int, updates map[string]interface{}) error {
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
	query := fmt.Sprintf(`UPDATE attribute_group_items SET %s WHERE id = $%d`, setClause, argPos)

	_, err := s.pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("ошибка обновления атрибута в группе: %w", err)
	}

	return nil
}

// GetAttributeGroupWithItems получает группу с атрибутами
func (s *attributeGroupStorage) GetAttributeGroupWithItems(ctx context.Context, id int) (*models.AttributeGroupWithItems, error) {
	group, err := s.GetAttributeGroup(ctx, id)
	if err != nil {
		return nil, err
	}

	items, err := s.GetGroupItems(ctx, id)
	if err != nil {
		return nil, err
	}

	return &models.AttributeGroupWithItems{
		AttributeGroup: group,
		Items:          items,
	}, nil
}

// AttachGroupToCategory привязывает группу к категории
func (s *attributeGroupStorage) AttachGroupToCategory(ctx context.Context, categoryID int, group *models.CategoryAttributeGroup) (int, error) {
	query := `
		INSERT INTO category_attribute_groups (
			category_id, group_id, component_id, sort_order, is_active, 
			display_mode, collapsed_by_default, configuration
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (category_id, group_id) DO UPDATE SET
			component_id = EXCLUDED.component_id,
			sort_order = EXCLUDED.sort_order,
			is_active = EXCLUDED.is_active,
			display_mode = EXCLUDED.display_mode,
			collapsed_by_default = EXCLUDED.collapsed_by_default,
			configuration = EXCLUDED.configuration
		RETURNING id`

	var id int
	err := s.pool.QueryRow(ctx, query,
		categoryID,
		group.GroupID,
		group.ComponentID,
		group.SortOrder,
		group.IsActive,
		group.DisplayMode,
		group.CollapsedByDefault,
		group.Configuration,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("ошибка привязки группы к категории: %w", err)
	}

	return id, nil
}

// DetachGroupFromCategory отвязывает группу от категории
func (s *attributeGroupStorage) DetachGroupFromCategory(ctx context.Context, categoryID, groupID int) error {
	query := `
		DELETE FROM category_attribute_groups 
		WHERE category_id = $1 AND group_id = $2`

	result, err := s.pool.Exec(ctx, query, categoryID, groupID)
	if err != nil {
		return fmt.Errorf("ошибка отвязки группы от категории: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("связь группы с категорией не найдена")
	}

	return nil
}

// GetCategoryGroups получает все группы категории
func (s *attributeGroupStorage) GetCategoryGroups(ctx context.Context, categoryID int) ([]*models.CategoryAttributeGroup, error) {
	query := `
		SELECT cag.id, cag.category_id, cag.group_id, cag.component_id, cag.sort_order,
		       cag.is_active, cag.display_mode, cag.collapsed_by_default, cag.configuration, cag.created_at,
		       ag.id, ag.name, ag.display_name, ag.description, ag.icon, ag.sort_order,
		       ag.is_active, ag.is_system, ag.created_at, ag.updated_at
		FROM category_attribute_groups cag
		LEFT JOIN attribute_groups ag ON cag.group_id = ag.id
		WHERE cag.category_id = $1
		ORDER BY cag.sort_order, ag.sort_order`

	rows, err := s.pool.Query(ctx, query, categoryID)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения групп категории: %w", err)
	}
	defer rows.Close()

	var groups []*models.CategoryAttributeGroup
	for rows.Next() {
		group := &models.CategoryAttributeGroup{
			Group: &models.AttributeGroup{},
		}
		err := rows.Scan(
			&group.ID,
			&group.CategoryID,
			&group.GroupID,
			&group.ComponentID,
			&group.SortOrder,
			&group.IsActive,
			&group.DisplayMode,
			&group.CollapsedByDefault,
			&group.Configuration,
			&group.CreatedAt,
			&group.Group.ID,
			&group.Group.Name,
			&group.Group.DisplayName,
			&group.Group.Description,
			&group.Group.Icon,
			&group.Group.SortOrder,
			&group.Group.IsActive,
			&group.Group.IsSystem,
			&group.Group.CreatedAt,
			&group.Group.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования группы категории: %w", err)
		}
		groups = append(groups, group)
	}

	return groups, nil
}

// UpdateCategoryGroup обновляет связь группы с категорией
func (s *attributeGroupStorage) UpdateCategoryGroup(ctx context.Context, id int, updates map[string]interface{}) error {
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
	query := fmt.Sprintf(`UPDATE category_attribute_groups SET %s WHERE id = $%d`, setClause, argPos)

	_, err := s.pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("ошибка обновления связи группы с категорией: %w", err)
	}

	return nil
}
