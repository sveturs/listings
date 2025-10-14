package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"backend/internal/domain"
	"backend/internal/domain/logistics"
	"backend/pkg/logger"
)

// ProblemService сервис для управления проблемными отправлениями
type ProblemService struct {
	db     *sql.DB
	logger *logger.Logger
}

// ProblemsFilter фильтры для получения проблем
type ProblemsFilter struct {
	Status      *string
	Severity    *string
	ProblemType *string
	AssignedTo  *int
	Page        int
	Limit       int
}

// NewProblemService создает новый сервис проблем
func NewProblemService(db *sql.DB) *ProblemService {
	return &ProblemService{
		db:     db,
		logger: logger.GetLogger(),
	}
}

// GetProblems получает список проблем с фильтрами
func (s *ProblemService) GetProblems(ctx context.Context, filter ProblemsFilter) ([]logistics.ProblemShipment, int, error) {
	// Базовый запрос
	baseQuery := `
		SELECT 
			p.id, p.shipment_id, p.shipment_type, p.tracking_number,
			p.problem_type, p.severity, p.description, p.status,
			p.assigned_to, p.resolution, p.order_id, p.user_id,
			p.metadata, p.created_at, p.updated_at, p.resolved_at,
			u.name, u.email,
			cu.name, cu.email
		FROM problem_shipments p
		LEFT JOIN users u ON u.id = p.assigned_to
		LEFT JOIN users cu ON cu.id = p.user_id
		WHERE 1=1
	`

	var conditions []string
	var args []interface{}
	argCount := 0

	// Применяем фильтры
	if filter.Status != nil {
		argCount++
		conditions = append(conditions, fmt.Sprintf("p.status = $%d", argCount))
		args = append(args, *filter.Status)
	}

	if filter.Severity != nil {
		argCount++
		conditions = append(conditions, fmt.Sprintf("p.severity = $%d", argCount))
		args = append(args, *filter.Severity)
	}

	if filter.ProblemType != nil {
		argCount++
		conditions = append(conditions, fmt.Sprintf("p.problem_type = $%d", argCount))
		args = append(args, *filter.ProblemType)
	}

	if filter.AssignedTo != nil {
		argCount++
		conditions = append(conditions, fmt.Sprintf("p.assigned_to = $%d", argCount))
		args = append(args, *filter.AssignedTo)
	}

	// Добавляем условия к запросу
	for _, cond := range conditions {
		baseQuery += " AND " + cond
	}

	// Добавляем сортировку
	baseQuery += " ORDER BY p.created_at DESC"

	// Получаем общее количество
	countQuery := "SELECT COUNT(*) FROM problem_shipments p WHERE 1=1"
	for _, cond := range conditions {
		countQuery += " AND " + cond
	}

	var total int
	err := s.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get total count: %w", err)
	}

	// Добавляем пагинацию
	if filter.Limit <= 0 {
		filter.Limit = 20
	}
	if filter.Page <= 0 {
		filter.Page = 1
	}
	offset := (filter.Page - 1) * filter.Limit

	argCount++
	baseQuery += fmt.Sprintf(" LIMIT $%d", argCount)
	args = append(args, filter.Limit)

	argCount++
	baseQuery += fmt.Sprintf(" OFFSET $%d", argCount)
	args = append(args, offset)

	// Выполняем запрос
	rows, err := s.db.QueryContext(ctx, baseQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get problems: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var problems []logistics.ProblemShipment
	for rows.Next() {
		var p logistics.ProblemShipment
		var assignedName, assignedEmail sql.NullString
		var userName, userEmail sql.NullString
		var metadataBytes []byte

		err := rows.Scan(
			&p.ID,
			&p.ShipmentID,
			&p.ShipmentType,
			&p.TrackingNumber,
			&p.ProblemType,
			&p.Severity,
			&p.Description,
			&p.Status,
			&p.AssignedTo,
			&p.Resolution,
			&p.OrderID,
			&p.UserID,
			&metadataBytes,
			&p.CreatedAt,
			&p.UpdatedAt,
			&p.ResolvedAt,
			&assignedName,
			&assignedEmail,
			&userName,
			&userEmail,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan problem: %w", err)
		}

		// Парсим metadata
		if len(metadataBytes) > 0 {
			err = json.Unmarshal(metadataBytes, &p.Metadata)
			if err != nil {
				p.Metadata = make(logistics.JSONB)
			}
		}

		// Добавляем информацию о назначенном пользователе
		if assignedName.Valid && p.AssignedTo != nil {
			p.AssignedUser = &logistics.User{
				ID:    *p.AssignedTo,
				Name:  assignedName.String,
				Email: assignedEmail.String,
			}
		}

		// Добавляем информацию о пользователе
		if userName.Valid && p.UserID != nil {
			p.User = &logistics.User{
				ID:    *p.UserID,
				Name:  userName.String,
				Email: userEmail.String,
			}
		}

		problems = append(problems, p)
	}

	return problems, total, rows.Err()
}

// CreateProblem создает новую проблему
func (s *ProblemService) CreateProblem(ctx context.Context, problem *logistics.ProblemShipment) (*logistics.ProblemShipment, error) {
	// Подготавливаем metadata
	metadataBytes, err := json.Marshal(problem.Metadata)
	if err != nil {
		metadataBytes = []byte("{}")
	}

	// Вставляем запись
	query := `
		INSERT INTO problem_shipments (
			shipment_id, shipment_type, tracking_number,
			problem_type, severity, description, status,
			assigned_to, order_id, user_id, metadata
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, created_at, updated_at
	`

	err = s.db.QueryRowContext(ctx, query,
		problem.ShipmentID,
		problem.ShipmentType,
		problem.TrackingNumber,
		problem.ProblemType,
		problem.Severity,
		problem.Description,
		problem.Status,
		problem.AssignedTo,
		problem.OrderID,
		problem.UserID,
		metadataBytes,
	).Scan(&problem.ID, &problem.CreatedAt, &problem.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create problem: %w", err)
	}

	// Загружаем полную информацию о созданной проблеме
	return s.GetProblemByID(ctx, problem.ID)
}

// UpdateProblem обновляет проблему
func (s *ProblemService) UpdateProblem(ctx context.Context, problemID int, updates map[string]interface{}) (*logistics.ProblemShipment, error) {
	// Проверяем существование проблемы
	exists, err := s.problemExists(ctx, problemID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, domain.ErrProblemNotFound
	}

	// Строим динамический запрос обновления
	setClauses := []string{}
	args := []interface{}{}
	argCount := 0

	for key, value := range updates {
		argCount++
		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", key, argCount))
		args = append(args, value)
	}

	if len(setClauses) == 0 {
		return s.GetProblemByID(ctx, problemID)
	}

	argCount++
	args = append(args, problemID)

	//nolint:gosec // setClauses are controlled, not user input
	query := fmt.Sprintf(`
		UPDATE problem_shipments 
		SET %s, updated_at = NOW()
		WHERE id = $%d
	`, joinStrings(setClauses, ", "), argCount)

	_, err = s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to update problem: %w", err)
	}

	return s.GetProblemByID(ctx, problemID)
}

// ResolveProblem решает проблему
func (s *ProblemService) ResolveProblem(ctx context.Context, problemID int, resolution string, resolvedBy int) error {
	// Проверяем существование проблемы
	exists, err := s.problemExists(ctx, problemID)
	if err != nil {
		return err
	}
	if !exists {
		return domain.ErrProblemNotFound
	}

	query := `
		UPDATE problem_shipments 
		SET status = 'resolved', 
		    resolution = $1, 
		    resolved_at = NOW(),
		    updated_at = NOW()
		WHERE id = $2
	`

	_, err = s.db.ExecContext(ctx, query, resolution, problemID)
	if err != nil {
		return fmt.Errorf("failed to resolve problem: %w", err)
	}

	// Логируем действие
	s.logAdminAction(ctx, resolvedBy, "problem", &problemID, "resolve", map[string]interface{}{
		"resolution": resolution,
	})

	return nil
}

// AssignProblem назначает проблему пользователю
func (s *ProblemService) AssignProblem(ctx context.Context, problemID int, assignTo int) error {
	// Проверяем существование проблемы
	exists, err := s.problemExists(ctx, problemID)
	if err != nil {
		return err
	}
	if !exists {
		return domain.ErrProblemNotFound
	}

	query := `
		UPDATE problem_shipments 
		SET assigned_to = $1,
		    status = CASE WHEN status = 'open' THEN 'investigating' ELSE status END,
		    updated_at = NOW()
		WHERE id = $2
	`

	_, err = s.db.ExecContext(ctx, query, assignTo, problemID)
	if err != nil {
		return fmt.Errorf("failed to assign problem: %w", err)
	}

	return nil
}

// GetProblemByID получает проблему по ID
func (s *ProblemService) GetProblemByID(ctx context.Context, problemID int) (*logistics.ProblemShipment, error) {
	query := `
		SELECT 
			p.id, p.shipment_id, p.shipment_type, p.tracking_number,
			p.problem_type, p.severity, p.description, p.status,
			p.assigned_to, p.resolution, p.order_id, p.user_id,
			p.metadata, p.created_at, p.updated_at, p.resolved_at,
			u.name, u.email
		FROM problem_shipments p
		LEFT JOIN users u ON u.id = p.assigned_to
		WHERE p.id = $1
	`

	var p logistics.ProblemShipment
	var assignedName, assignedEmail sql.NullString
	var metadataBytes []byte

	err := s.db.QueryRowContext(ctx, query, problemID).Scan(
		&p.ID,
		&p.ShipmentID,
		&p.ShipmentType,
		&p.TrackingNumber,
		&p.ProblemType,
		&p.Severity,
		&p.Description,
		&p.Status,
		&p.AssignedTo,
		&p.Resolution,
		&p.OrderID,
		&p.UserID,
		&metadataBytes,
		&p.CreatedAt,
		&p.UpdatedAt,
		&p.ResolvedAt,
		&assignedName,
		&assignedEmail,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrProblemNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get problem: %w", err)
	}

	// Парсим metadata
	if len(metadataBytes) > 0 {
		err = json.Unmarshal(metadataBytes, &p.Metadata)
		if err != nil {
			p.Metadata = make(logistics.JSONB)
		}
	}

	// Добавляем информацию о назначенном пользователе
	if assignedName.Valid && p.AssignedTo != nil {
		p.AssignedUser = &logistics.User{
			ID:    *p.AssignedTo,
			Name:  assignedName.String,
			Email: assignedEmail.String,
		}
	}

	return &p, nil
}

// problemExists проверяет существование проблемы
func (s *ProblemService) problemExists(ctx context.Context, problemID int) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM problem_shipments WHERE id = $1)`
	err := s.db.QueryRowContext(ctx, query, problemID).Scan(&exists)
	return exists, err
}

// logAdminAction логирует действие администратора
func (s *ProblemService) logAdminAction(ctx context.Context, adminID int, entityType string, entityID *int, action string, details map[string]interface{}) {
	detailsBytes, _ := json.Marshal(details)

	query := `
		INSERT INTO logistics_admin_logs (
			admin_id, admin_email, entity_type, entity_id, action, details
		) VALUES (
			$1, 
			(SELECT email FROM users WHERE id = $1),
			$2, $3, $4, $5
		)
	`

	_, err := s.db.ExecContext(ctx, query, adminID, entityType, entityID, action, detailsBytes)
	if err != nil {
		s.logger.Error("Failed to log admin action: %v", err)
	}
}

// AutoCreateProblemsForDelays автоматически создает проблемы для задержанных отправлений
func (s *ProblemService) AutoCreateProblemsForDelays(ctx context.Context) error {
	// Получаем настройки
	var delayThreshold int
	err := s.db.QueryRowContext(ctx, `
		SELECT delay_threshold_hours 
		FROM logistics_monitoring_settings 
		WHERE id = 1
	`).Scan(&delayThreshold)
	if err != nil {
		delayThreshold = 72 // По умолчанию 3 дня
	}

	// Находим задержанные отправления без проблем
	query := `
		WITH delayed_shipments AS (
			SELECT 
				id, 'bex' as type, tracking_number, marketplace_order_id as order_id, buyer_id as user_id
			FROM bex_shipments
			WHERE status IN ('in_transit', 'picked_up') 
				AND created_at < NOW() - INTERVAL '%d hours'
				AND NOT EXISTS (
					SELECT 1 FROM problem_shipments 
					WHERE shipment_id = bex_shipments.id 
					AND shipment_type = 'bex'
					AND status NOT IN ('resolved', 'closed')
				)
			UNION ALL
			SELECT 
				id, 'postexpress' as type, tracking_number, order_id, NULL as user_id
			FROM post_express_shipments
			WHERE status IN ('in_transit', 'processing') 
				AND created_at < NOW() - INTERVAL '%d hours'
				AND NOT EXISTS (
					SELECT 1 FROM problem_shipments 
					WHERE shipment_id = post_express_shipments.id 
					AND shipment_type = 'postexpress'
					AND status NOT IN ('resolved', 'closed')
				)
		)
		INSERT INTO problem_shipments (
			shipment_id, shipment_type, tracking_number, problem_type, 
			severity, description, status, order_id, user_id
		)
		SELECT 
			id, type, tracking_number, 'delayed',
			CASE 
				WHEN EXTRACT(EPOCH FROM (NOW() - created_at))/3600 > 168 THEN 'critical'
				WHEN EXTRACT(EPOCH FROM (NOW() - created_at))/3600 > 120 THEN 'high'
				ELSE 'medium'
			END,
			'Shipment delayed for more than ' || %d || ' hours',
			'open', order_id, user_id
		FROM delayed_shipments
	`

	query = fmt.Sprintf(query, delayThreshold, delayThreshold, delayThreshold)

	result, err := s.db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to create problems for delays: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected > 0 {
		s.logger.Info("Created problems for delayed shipments: count=%d", rowsAffected)
	}

	return nil
}

// GetProblemComments получает комментарии к проблеме
func (s *ProblemService) GetProblemComments(ctx context.Context, problemID int) ([]logistics.ProblemComment, error) {
	query := `
		SELECT 
			c.id, c.problem_id, c.admin_id, c.comment, c.comment_type,
			c.metadata, c.created_at,
			u.name as admin_name, u.email as admin_email
		FROM problem_comments c
		LEFT JOIN users u ON u.id = c.admin_id
		WHERE c.problem_id = $1
		ORDER BY c.created_at ASC
	`

	rows, err := s.db.QueryContext(ctx, query, problemID)
	if err != nil {
		return nil, fmt.Errorf("failed to get problem comments: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var comments []logistics.ProblemComment
	for rows.Next() {
		var comment logistics.ProblemComment
		var adminName, adminEmail sql.NullString
		var metadataBytes []byte

		err := rows.Scan(
			&comment.ID,
			&comment.ProblemID,
			&comment.AdminID,
			&comment.Comment,
			&comment.CommentType,
			&metadataBytes,
			&comment.CreatedAt,
			&adminName,
			&adminEmail,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan comment: %w", err)
		}

		// Парсим metadata
		if len(metadataBytes) > 0 {
			err = json.Unmarshal(metadataBytes, &comment.Metadata)
			if err != nil {
				comment.Metadata = make(logistics.JSONB)
			}
		}

		// Добавляем информацию об админе
		if adminName.Valid {
			comment.Admin = &logistics.User{
				ID:    comment.AdminID,
				Name:  adminName.String,
				Email: adminEmail.String,
			}
		}

		comments = append(comments, comment)
	}

	return comments, rows.Err()
}

// AddProblemComment добавляет комментарий к проблеме
func (s *ProblemService) AddProblemComment(ctx context.Context, problemID, adminID int, commentText, commentType string, metadata map[string]interface{}) (*logistics.ProblemComment, error) {
	// Подготавливаем metadata
	metadataBytes, err := json.Marshal(metadata)
	if err != nil {
		metadataBytes = []byte("{}")
	}

	// Вставляем комментарий
	query := `
		INSERT INTO problem_comments (problem_id, admin_id, comment, comment_type, metadata)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`

	var comment logistics.ProblemComment
	err = s.db.QueryRowContext(ctx, query, problemID, adminID, commentText, commentType, metadataBytes).Scan(
		&comment.ID,
		&comment.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to add comment: %w", err)
	}

	// Заполняем остальные поля
	comment.ProblemID = problemID
	comment.AdminID = adminID
	comment.Comment = commentText
	comment.CommentType = commentType
	if len(metadata) > 0 {
		comment.Metadata = metadata
	}

	// Получаем информацию об админе
	adminQuery := `SELECT name, email FROM users WHERE id = $1`
	var adminName, adminEmail string
	err = s.db.QueryRowContext(ctx, adminQuery, adminID).Scan(&adminName, &adminEmail)
	if err == nil {
		comment.Admin = &logistics.User{
			ID:    adminID,
			Name:  adminName,
			Email: adminEmail,
		}
	}

	return &comment, nil
}

// GetProblemHistory получает историю изменений статуса проблемы
func (s *ProblemService) GetProblemHistory(ctx context.Context, problemID int) ([]logistics.ProblemStatusHistory, error) {
	query := `
		SELECT 
			h.id, h.problem_id, h.admin_id, h.old_status, h.new_status,
			h.old_assigned_to, h.new_assigned_to, h.comment, h.metadata, h.created_at,
			u.name as admin_name, u.email as admin_email,
			u1.name as old_assigned_name, u1.email as old_assigned_email,
			u2.name as new_assigned_name, u2.email as new_assigned_email
		FROM problem_status_history h
		LEFT JOIN users u ON u.id = h.admin_id
		LEFT JOIN users u1 ON u1.id = h.old_assigned_to
		LEFT JOIN users u2 ON u2.id = h.new_assigned_to
		WHERE h.problem_id = $1
		ORDER BY h.created_at ASC
	`

	rows, err := s.db.QueryContext(ctx, query, problemID)
	if err != nil {
		return nil, fmt.Errorf("failed to get problem history: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var history []logistics.ProblemStatusHistory
	for rows.Next() {
		var h logistics.ProblemStatusHistory
		var adminName, adminEmail sql.NullString
		var oldAssignedName, oldAssignedEmail sql.NullString
		var newAssignedName, newAssignedEmail sql.NullString
		var metadataBytes []byte

		err := rows.Scan(
			&h.ID,
			&h.ProblemID,
			&h.AdminID,
			&h.OldStatus,
			&h.NewStatus,
			&h.OldAssignedTo,
			&h.NewAssignedTo,
			&h.Comment,
			&metadataBytes,
			&h.CreatedAt,
			&adminName,
			&adminEmail,
			&oldAssignedName,
			&oldAssignedEmail,
			&newAssignedName,
			&newAssignedEmail,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan history: %w", err)
		}

		// Парсим metadata
		if len(metadataBytes) > 0 {
			err = json.Unmarshal(metadataBytes, &h.Metadata)
			if err != nil {
				h.Metadata = make(logistics.JSONB)
			}
		}

		// Добавляем информацию об админе
		if adminName.Valid && h.AdminID != nil {
			h.Admin = &logistics.User{
				ID:    *h.AdminID,
				Name:  adminName.String,
				Email: adminEmail.String,
			}
		}

		// Добавляем информацию о старом назначенном
		if oldAssignedName.Valid && h.OldAssignedTo != nil {
			h.OldAssigned = &logistics.User{
				ID:    *h.OldAssignedTo,
				Name:  oldAssignedName.String,
				Email: oldAssignedEmail.String,
			}
		}

		// Добавляем информацию о новом назначенном
		if newAssignedName.Valid && h.NewAssignedTo != nil {
			h.NewAssigned = &logistics.User{
				ID:    *h.NewAssignedTo,
				Name:  newAssignedName.String,
				Email: newAssignedEmail.String,
			}
		}

		history = append(history, h)
	}

	return history, rows.Err()
}

// GetProblemWithDetails получает проблему со всеми деталями (комментарии, история)
func (s *ProblemService) GetProblemWithDetails(ctx context.Context, problemID int) (*logistics.ProblemShipment, error) {
	// Получаем основную информацию о проблеме
	problem, err := s.GetProblemByID(ctx, problemID)
	if err != nil {
		return nil, err
	}

	// Получаем комментарии
	comments, err := s.GetProblemComments(ctx, problemID)
	if err != nil {
		s.logger.Error("Failed to get problem comments: error=%v problem_id=%d", err, problemID)
		// Не возвращаем ошибку, просто логируем
		comments = []logistics.ProblemComment{}
	}

	// Получаем историю
	history, err := s.GetProblemHistory(ctx, problemID)
	if err != nil {
		s.logger.Error("Failed to get problem history: error=%v problem_id=%d", err, problemID)
		// Не возвращаем ошибку, просто логируем
		history = []logistics.ProblemStatusHistory{}
	}

	problem.Comments = comments
	problem.History = history

	return problem, nil
}

// joinStrings объединяет строки через разделитель
func joinStrings(strs []string, sep string) string {
	result := ""
	for i, s := range strs {
		if i > 0 {
			result += sep
		}
		result += s
	}
	return result
}
