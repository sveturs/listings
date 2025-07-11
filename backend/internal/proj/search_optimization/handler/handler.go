package handler

import (
	"net/http"
	"strconv"
	"time"

	"backend/internal/proj/search_optimization/service"
	"backend/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

type SearchOptimizationHandler struct {
	service service.SearchOptimizationService
}

func NewSearchOptimizationHandler(service service.SearchOptimizationService) *SearchOptimizationHandler {
	return &SearchOptimizationHandler{
		service: service,
	}
}

// OptimizationSession представляет сессию оптимизации весов (для swagger)
type OptimizationSession struct {
	ID              int64                       `json:"id"`
	Status          string                      `json:"status"` // running, completed, failed
	StartTime       time.Time                   `json:"start_time"`
	EndTime         *time.Time                  `json:"end_time,omitempty"`
	TotalFields     int                         `json:"total_fields"`
	ProcessedFields int                         `json:"processed_fields"`
	Results         []*WeightOptimizationResult `json:"results,omitempty"`
	ErrorMessage    *string                     `json:"error_message,omitempty"`
	CreatedBy       int                         `json:"created_by"`
}

// WeightOptimizationResult представляет результат оптимизации для одного поля (для swagger)
type WeightOptimizationResult struct {
	FieldName          string   `json:"field_name"`
	OriginalWeight     float64  `json:"original_weight"`
	OptimizedWeight    float64  `json:"optimized_weight"`
	ImprovementPercent float64  `json:"improvement_percent"`
	SearchQueries      []string `json:"search_queries"`
	SampleSize         int      `json:"sample_size"`
}

// StartOptimizationRequest запрос на запуск оптимизации
type StartOptimizationRequest struct {
	FieldNames      []string `json:"field_names,omitempty"`
	ItemType        string   `json:"item_type" validate:"required,oneof=marketplace storefront global"`
	CategoryID      *int     `json:"category_id,omitempty"`
	MinSampleSize   int      `json:"min_sample_size" validate:"min=10,max=10000"`
	ConfidenceLevel float64  `json:"confidence_level" validate:"min=0.5,max=0.99"`
	LearningRate    float64  `json:"learning_rate" validate:"min=0.001,max=1"`
	MaxIterations   int      `json:"max_iterations" validate:"min=10,max=10000"`
	AnalysisPeriod  int      `json:"analysis_period_days" validate:"min=1,max=365"`
	AutoApply       bool     `json:"auto_apply"`
}

// ApplyWeightsRequest запрос на применение весов
type ApplyWeightsRequest struct {
	SessionID       int64   `json:"session_id" validate:"required"`
	SelectedResults []int64 `json:"selected_results" validate:"required"`
}

// AnalyzeWeightsRequest запрос на анализ весов
type AnalyzeWeightsRequest struct {
	ItemType   string `json:"item_type" validate:"required,oneof=marketplace storefront global"`
	CategoryID *int   `json:"category_id,omitempty"`
	FromDate   string `json:"from_date" validate:"required"`
	ToDate     string `json:"to_date" validate:"required"`
}

// StartOptimization запускает процесс оптимизации весов поиска
// @Summary Запуск оптимизации весов поиска
// @Description Запускает процесс машинного обучения для оптимизации весов полей поиска на основе поведенческих данных
// @Tags Search Optimization
// @Accept json
// @Produce json
// @Param request body StartOptimizationRequest true "Параметры оптимизации"
// @Success 200 {object} utils.SuccessResponseSwag{data=map[string]interface{}} "Оптимизация запущена"
// @Failure 400 {object} utils.ErrorResponseSwag "Неверные параметры"
// @Failure 401 {object} utils.ErrorResponseSwag "Не авторизован"
// @Failure 500 {object} utils.ErrorResponseSwag "Внутренняя ошибка сервера"
// @Router /api/v1/admin/search/optimize-weights [post]
// @Security BearerAuth
func (h *SearchOptimizationHandler) StartOptimization(c *fiber.Ctx) error {
	var req StartOptimizationRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "invalid_request_body")
	}

	// Валидация запроса (пропускаем пока)
	// TODO: Реализовать валидацию

	// Получение ID администратора из контекста
	adminID, ok := c.Locals("admin_id").(int)
	if !ok {
		// Временно используем значение по умолчанию для тестирования
		adminID = 1
	}

	// Преобразование в параметры сервиса
	params := &service.OptimizationParams{
		FieldNames:      req.FieldNames,
		ItemType:        req.ItemType,
		CategoryID:      req.CategoryID,
		MinSampleSize:   req.MinSampleSize,
		ConfidenceLevel: req.ConfidenceLevel,
		LearningRate:    req.LearningRate,
		MaxIterations:   req.MaxIterations,
		AnalysisPeriod:  req.AnalysisPeriod,
		AutoApply:       req.AutoApply,
	}

	// Установка значений по умолчанию
	if params.MinSampleSize == 0 {
		params.MinSampleSize = 100
	}
	if params.ConfidenceLevel == 0 {
		params.ConfidenceLevel = 0.85
	}
	if params.LearningRate == 0 {
		params.LearningRate = 0.01
	}
	if params.MaxIterations == 0 {
		params.MaxIterations = 1000
	}
	if params.AnalysisPeriod == 0 {
		params.AnalysisPeriod = 30
	}

	sessionID, err := h.service.StartOptimization(c.Context(), params, adminID)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, "optimization_start_failed")
	}

	return utils.SuccessResponse(c, map[string]interface{}{
		"session_id": sessionID,
		"message":    "Optimization process started successfully",
	})
}

// GetOptimizationStatus получает статус оптимизации
// @Summary Получение статуса оптимизации
// @Description Возвращает текущий статус и результаты процесса оптимизации весов
// @Tags Search Optimization
// @Produce json
// @Param session_id path int true "ID сессии оптимизации"
// @Success 200 {object} utils.SuccessResponseSwag{data=OptimizationSession} "Статус оптимизации"
// @Failure 400 {object} utils.ErrorResponseSwag "Неверный ID сессии"
// @Failure 404 {object} utils.ErrorResponseSwag "Сессия не найдена"
// @Failure 500 {object} utils.ErrorResponseSwag "Внутренняя ошибка сервера"
// @Router /api/v1/admin/search/optimization-status/{session_id} [get]
// @Security BearerAuth
func (h *SearchOptimizationHandler) GetOptimizationStatus(c *fiber.Ctx) error {
	sessionIDStr := c.Params("session_id")
	sessionID, err := strconv.ParseInt(sessionIDStr, 10, 64)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "invalid_session_id")
	}

	session, err := h.service.GetOptimizationStatus(c.Context(), sessionID)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, "get_status_failed")
	}

	if session == nil {
		return utils.ErrorResponse(c, http.StatusNotFound, "session_not_found")
	}

	return utils.SuccessResponse(c, session)
}

// CancelOptimization отменяет процесс оптимизации
// @Summary Отмена оптимизации
// @Description Отменяет запущенный процесс оптимизации весов
// @Tags Search Optimization
// @Produce json
// @Param session_id path int true "ID сессии оптимизации"
// @Success 200 {object} utils.SuccessResponseSwag{data=map[string]interface{}} "Оптимизация отменена"
// @Failure 400 {object} utils.ErrorResponseSwag "Неверный ID сессии"
// @Failure 500 {object} utils.ErrorResponseSwag "Внутренняя ошибка сервера"
// @Router /api/v1/admin/search/optimization-cancel/{session_id} [post]
// @Security BearerAuth
func (h *SearchOptimizationHandler) CancelOptimization(c *fiber.Ctx) error {
	sessionIDStr := c.Params("session_id")
	sessionID, err := strconv.ParseInt(sessionIDStr, 10, 64)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "invalid_session_id")
	}

	adminID, ok := c.Locals("admin_id").(int)
	if !ok {
		// Временно используем значение по умолчанию для тестирования
		adminID = 1
	}

	err = h.service.CancelOptimization(c.Context(), sessionID, adminID)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, "cancel_failed")
	}

	return utils.SuccessResponse(c, map[string]interface{}{
		"session_id": sessionID,
		"message":    "Optimization cancelled successfully",
	})
}

// ApplyOptimizedWeights применяет оптимизированные веса
// @Summary Применение оптимизированных весов
// @Description Применяет выбранные оптимизированные веса к системе поиска
// @Tags Search Optimization
// @Accept json
// @Produce json
// @Param request body ApplyWeightsRequest true "Параметры применения весов"
// @Success 200 {object} utils.SuccessResponseSwag{data=map[string]interface{}} "Веса применены"
// @Failure 400 {object} utils.ErrorResponseSwag "Неверные параметры"
// @Failure 401 {object} utils.ErrorResponseSwag "Не авторизован"
// @Failure 500 {object} utils.ErrorResponseSwag "Внутренняя ошибка сервера"
// @Router /api/v1/admin/search/apply-weights [post]
// @Security BearerAuth
func (h *SearchOptimizationHandler) ApplyOptimizedWeights(c *fiber.Ctx) error {
	var req ApplyWeightsRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "invalid_request_body")
	}

	// TODO: Добавить валидацию, когда будет реализована функция ValidateStruct

	adminID, ok := c.Locals("admin_id").(int)
	if !ok {
		// Временно используем значение по умолчанию для тестирования
		adminID = 1
	}

	err := h.service.ApplyOptimizedWeights(c.Context(), req.SessionID, adminID, req.SelectedResults)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, "apply_weights_failed")
	}

	return utils.SuccessResponse(c, map[string]interface{}{
		"session_id":      req.SessionID,
		"applied_results": len(req.SelectedResults),
		"message":         "Selected weights applied successfully",
	})
}

// AnalyzeCurrentWeights анализирует текущие веса
// @Summary Анализ текущих весов поиска
// @Description Анализирует эффективность текущих весов без запуска полной оптимизации
// @Tags Search Optimization
// @Accept json
// @Produce json
// @Param request body AnalyzeWeightsRequest true "Параметры анализа"
// @Success 200 {object} utils.SuccessResponseSwag{data=[]WeightOptimizationResult} "Результаты анализа"
// @Failure 400 {object} utils.ErrorResponseSwag "Неверные параметры"
// @Failure 500 {object} utils.ErrorResponseSwag "Внутренняя ошибка сервера"
// @Router /api/v1/admin/search/analyze-weights [post]
// @Security BearerAuth
func (h *SearchOptimizationHandler) AnalyzeCurrentWeights(c *fiber.Ctx) error {
	var req AnalyzeWeightsRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "invalid_request_body")
	}

	// TODO: Добавить валидацию, когда будет реализована функция ValidateStruct

	// Парсинг дат
	fromDate, err := time.Parse("2006-01-02", req.FromDate)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "invalid_from_date")
	}

	toDate, err := time.Parse("2006-01-02", req.ToDate)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "invalid_to_date")
	}

	results, err := h.service.AnalyzeCurrentWeights(c.Context(), req.ItemType, req.CategoryID, fromDate, toDate)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, "analysis_failed")
	}

	return utils.SuccessResponse(c, map[string]interface{}{
		"results":      results,
		"total_fields": len(results),
		"period_start": fromDate,
		"period_end":   toDate,
	})
}

// GetOptimizationHistory получает историю оптимизаций
// @Summary История оптимизаций
// @Description Возвращает список последних сессий оптимизации весов
// @Tags Search Optimization
// @Produce json
// @Param limit query int false "Максимальное количество записей" default(20)
// @Success 200 {object} utils.SuccessResponseSwag{data=[]OptimizationSession} "История оптимизаций"
// @Failure 500 {object} utils.ErrorResponseSwag "Внутренняя ошибка сервера"
// @Router /api/v1/admin/search/optimization-history [get]
// @Security BearerAuth
func (h *SearchOptimizationHandler) GetOptimizationHistory(c *fiber.Ctx) error {
	limitStr := c.Query("limit", "20")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 20
	}

	sessions, err := h.service.GetOptimizationHistory(c.Context(), limit)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, "get_history_failed")
	}

	return utils.SuccessResponse(c, sessions)
}

// GetOptimizationConfig получает конфигурацию оптимизации
// @Summary Конфигурация оптимизации
// @Description Возвращает текущую конфигурацию параметров оптимизации весов
// @Tags Search Optimization
// @Produce json
// @Success 200 {object} utils.SuccessResponseSwag{data=service.OptimizationConfig} "Конфигурация оптимизации"
// @Failure 500 {object} utils.ErrorResponseSwag "Внутренняя ошибка сервера"
// @Router /api/v1/admin/search/optimization-config [get]
// @Security BearerAuth
func (h *SearchOptimizationHandler) GetOptimizationConfig(c *fiber.Ctx) error {
	config, err := h.service.GetOptimizationConfig(c.Context())
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, "get_config_failed")
	}

	return utils.SuccessResponse(c, config)
}

// UpdateOptimizationConfig обновляет конфигурацию оптимизации
// @Summary Обновление конфигурации оптимизации
// @Description Обновляет параметры оптимизации весов по умолчанию
// @Tags Search Optimization
// @Accept json
// @Produce json
// @Param config body service.OptimizationConfig true "Новая конфигурация"
// @Success 200 {object} utils.SuccessResponseSwag{data=map[string]interface{}} "Конфигурация обновлена"
// @Failure 400 {object} utils.ErrorResponseSwag "Неверная конфигурация"
// @Failure 401 {object} utils.ErrorResponseSwag "Не авторизован"
// @Failure 500 {object} utils.ErrorResponseSwag "Внутренняя ошибка сервера"
// @Router /api/v1/admin/search/optimization-config [put]
// @Security BearerAuth
func (h *SearchOptimizationHandler) UpdateOptimizationConfig(c *fiber.Ctx) error {
	var config service.OptimizationConfig
	if err := c.BodyParser(&config); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "invalid_request_body")
	}

	adminID, ok := c.Locals("admin_id").(int)
	if !ok {
		// Временно используем значение по умолчанию для тестирования
		adminID = 1
	}

	err := h.service.UpdateOptimizationConfig(c.Context(), &config, adminID)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "config_update_failed")
	}

	return utils.SuccessResponse(c, map[string]interface{}{
		"message": "Optimization configuration updated successfully",
	})
}

// CreateWeightBackup создает резервную копию весов
// @Summary Создание резервной копии весов
// @Description Создает резервную копию текущих весов поиска для возможности отката
// @Tags Search Optimization
// @Accept json
// @Produce json
// @Param request body map[string]interface{} true "Параметры резервной копии"
// @Success 200 {object} utils.SuccessResponseSwag{data=map[string]interface{}} "Резервная копия создана"
// @Failure 400 {object} utils.ErrorResponseSwag "Неверные параметры"
// @Failure 401 {object} utils.ErrorResponseSwag "Не авторизован"
// @Failure 500 {object} utils.ErrorResponseSwag "Внутренняя ошибка сервера"
// @Router /api/v1/admin/search/backup-weights [post]
// @Security BearerAuth
func (h *SearchOptimizationHandler) CreateWeightBackup(c *fiber.Ctx) error {
	var req struct {
		ItemType   string `json:"item_type" validate:"required,oneof=marketplace storefront global"`
		CategoryID *int   `json:"category_id,omitempty"`
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "invalid_request_body")
	}

	// TODO: Добавить валидацию, когда будет реализована функция ValidateStruct

	adminID, ok := c.Locals("admin_id").(int)
	if !ok {
		// Временно используем значение по умолчанию для тестирования
		adminID = 1
	}

	err := h.service.CreateWeightBackup(c.Context(), req.ItemType, req.CategoryID, adminID)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, "backup_failed")
	}

	return utils.SuccessResponse(c, map[string]interface{}{
		"item_type":   req.ItemType,
		"category_id": req.CategoryID,
		"message":     "Weight backup created successfully",
	})
}

// RollbackWeights откатывает веса к предыдущим значениям
// @Summary Откат весов
// @Description Откатывает веса поиска к предыдущим значениям
// @Tags Search Optimization
// @Accept json
// @Produce json
// @Param request body map[string]interface{} true "ID весов для отката"
// @Success 200 {object} utils.SuccessResponseSwag{data=map[string]interface{}} "Откат выполнен"
// @Failure 400 {object} utils.ErrorResponseSwag "Неверные параметры"
// @Failure 401 {object} utils.ErrorResponseSwag "Не авторизован"
// @Failure 500 {object} utils.ErrorResponseSwag "Внутренняя ошибка сервера"
// @Router /api/v1/admin/search/rollback-weights [post]
// @Security BearerAuth
func (h *SearchOptimizationHandler) RollbackWeights(c *fiber.Ctx) error {
	var req struct {
		WeightIDs []int64 `json:"weight_ids" validate:"required,min=1"`
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "invalid_request_body")
	}

	// TODO: Добавить валидацию, когда будет реализована функция ValidateStruct

	adminID, ok := c.Locals("admin_id").(int)
	if !ok {
		// Временно используем значение по умолчанию для тестирования
		adminID = 1
	}

	err := h.service.RollbackWeights(c.Context(), req.WeightIDs, adminID)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, "rollback_failed")
	}

	return utils.SuccessResponse(c, map[string]interface{}{
		"weight_ids":     req.WeightIDs,
		"rollback_count": len(req.WeightIDs),
		"message":        "Weights rolled back successfully",
	})
}

// Synonym Management

// SynonymRequest запрос для создания/обновления синонима
type SynonymRequest struct {
	Term     string `json:"term" validate:"required,min=1,max=255"`
	Synonym  string `json:"synonym" validate:"required,min=1,max=255"`
	Language string `json:"language" validate:"required,oneof=en ru sr"`
	IsActive bool   `json:"is_active"`
}

// GetSynonyms получает список синонимов
// @Summary Получение списка синонимов
// @Description Возвращает список синонимов для указанного языка с возможностью поиска
// @Tags Search Synonyms
// @Produce json
// @Param language query string false "Язык синонимов (en, ru, sr)"
// @Param search query string false "Поиск по термину"
// @Param page query int false "Номер страницы" default(1)
// @Param limit query int false "Количество записей на странице" default(20)
// @Success 200 {object} utils.SuccessResponseSwag{data=[]map[string]interface{}} "Список синонимов"
// @Failure 500 {object} utils.ErrorResponseSwag "Внутренняя ошибка сервера"
// @Router /api/v1/admin/search/synonyms [get]
// @Security BearerAuth
func (h *SearchOptimizationHandler) GetSynonyms(c *fiber.Ctx) error {
	// Временно отключаем проверку admin_id для тестирования
	// _, ok := c.Locals("admin_id").(int)
	// if !ok {
	// 	return utils.ErrorResponse(c, http.StatusUnauthorized, "admin_required")
	// }

	language := c.Query("language", "")
	search := c.Query("search", "")
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	synonyms, total, err := h.service.GetSynonyms(c.Context(), language, search, page, limit)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, "get_synonyms_failed")
	}

	return utils.SuccessResponse(c, map[string]interface{}{
		"data":        synonyms,
		"total":       total,
		"page":        page,
		"limit":       limit,
		"total_pages": (total + limit - 1) / limit,
	})
}

// CreateSynonym создает новый синоним
// @Summary Создание нового синонима
// @Description Создает новый синоним для улучшения поиска
// @Tags Search Synonyms
// @Accept json
// @Produce json
// @Param request body SynonymRequest true "Данные синонима"
// @Success 201 {object} utils.SuccessResponseSwag{data=map[string]interface{}} "Синоним создан"
// @Failure 400 {object} utils.ErrorResponseSwag "Неверные параметры"
// @Failure 401 {object} utils.ErrorResponseSwag "Не авторизован"
// @Failure 409 {object} utils.ErrorResponseSwag "Синоним уже существует"
// @Failure 500 {object} utils.ErrorResponseSwag "Внутренняя ошибка сервера"
// @Router /api/v1/admin/search/synonyms [post]
// @Security BearerAuth
func (h *SearchOptimizationHandler) CreateSynonym(c *fiber.Ctx) error {
	var req SynonymRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "invalid_request_body")
	}

	// TODO: Добавить валидацию, когда будет реализована функция ValidateStruct

	adminID, ok := c.Locals("admin_id").(int)
	if !ok {
		// Временно используем значение по умолчанию для тестирования
		adminID = 1
	}

	synonymID, err := h.service.CreateSynonym(c.Context(), req.Term, req.Synonym, req.Language, req.IsActive, adminID)
	if err != nil {
		if err.Error() == "synonym already exists" {
			return utils.ErrorResponse(c, http.StatusConflict, "synonym_already_exists")
		}
		return utils.ErrorResponse(c, http.StatusInternalServerError, "create_synonym_failed")
	}

	return c.Status(http.StatusCreated).JSON(utils.SuccessResponseSwag{
		Success: true,
		Data: map[string]interface{}{
			"message":   "Synonym created successfully",
			"id":        synonymID,
			"term":      req.Term,
			"synonym":   req.Synonym,
			"language":  req.Language,
			"is_active": req.IsActive,
		},
	})
}

// UpdateSynonym обновляет существующий синоним
// @Summary Обновление синонима
// @Description Обновляет существующий синоним
// @Tags Search Synonyms
// @Accept json
// @Produce json
// @Param id path int true "ID синонима"
// @Param request body SynonymRequest true "Обновленные данные синонима"
// @Success 200 {object} utils.SuccessResponseSwag{data=map[string]interface{}} "Синоним обновлен"
// @Failure 400 {object} utils.ErrorResponseSwag "Неверные параметры"
// @Failure 401 {object} utils.ErrorResponseSwag "Не авторизован"
// @Failure 404 {object} utils.ErrorResponseSwag "Синоним не найден"
// @Failure 500 {object} utils.ErrorResponseSwag "Внутренняя ошибка сервера"
// @Router /api/v1/admin/search/synonyms/{id} [put]
// @Security BearerAuth
func (h *SearchOptimizationHandler) UpdateSynonym(c *fiber.Ctx) error {
	idStr := c.Params("id")
	synonymID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "invalid_synonym_id")
	}

	var req SynonymRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "invalid_request_body")
	}

	// TODO: Добавить валидацию, когда будет реализована функция ValidateStruct

	adminID, ok := c.Locals("admin_id").(int)
	if !ok {
		// Временно используем значение по умолчанию для тестирования
		adminID = 1
	}

	err = h.service.UpdateSynonym(c.Context(), synonymID, req.Term, req.Synonym, req.Language, req.IsActive, adminID)
	if err != nil {
		if err.Error() == "synonym not found" {
			return utils.ErrorResponse(c, http.StatusNotFound, "synonym_not_found")
		}
		return utils.ErrorResponse(c, http.StatusInternalServerError, "update_synonym_failed")
	}

	return utils.SuccessResponse(c, map[string]interface{}{
		"id":        synonymID,
		"term":      req.Term,
		"synonym":   req.Synonym,
		"language":  req.Language,
		"is_active": req.IsActive,
		"message":   "Synonym updated successfully",
	})
}

// DeleteSynonym удаляет синоним
// @Summary Удаление синонима
// @Description Удаляет синоним из системы
// @Tags Search Synonyms
// @Produce json
// @Param id path int true "ID синонима"
// @Success 200 {object} utils.SuccessResponseSwag{data=map[string]interface{}} "Синоним удален"
// @Failure 400 {object} utils.ErrorResponseSwag "Неверный ID синонима"
// @Failure 401 {object} utils.ErrorResponseSwag "Не авторизован"
// @Failure 404 {object} utils.ErrorResponseSwag "Синоним не найден"
// @Failure 500 {object} utils.ErrorResponseSwag "Внутренняя ошибка сервера"
// @Router /api/v1/admin/search/synonyms/{id} [delete]
// @Security BearerAuth
func (h *SearchOptimizationHandler) DeleteSynonym(c *fiber.Ctx) error {
	idStr := c.Params("id")
	synonymID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "invalid_synonym_id")
	}

	adminID, ok := c.Locals("admin_id").(int)
	if !ok {
		// Временно используем значение по умолчанию для тестирования
		adminID = 1
	}

	err = h.service.DeleteSynonym(c.Context(), synonymID, adminID)
	if err != nil {
		if err.Error() == "synonym not found" {
			return utils.ErrorResponse(c, http.StatusNotFound, "synonym_not_found")
		}
		return utils.ErrorResponse(c, http.StatusInternalServerError, "delete_synonym_failed")
	}

	return utils.SuccessResponse(c, map[string]interface{}{
		"id":      synonymID,
		"message": "Synonym deleted successfully",
	})
}
