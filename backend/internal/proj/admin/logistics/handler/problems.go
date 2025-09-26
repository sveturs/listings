package handler

import (
	"fmt"
	"strconv"
	"time"

	"backend/internal/domain/logistics"
	"backend/internal/proj/admin/logistics/service"
	"backend/pkg/logger"
	"backend/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

const (
	problemNotFoundErr = "problem not found"
)

// ProblemsHandler обработчик для управления проблемными отправлениями
type ProblemsHandler struct {
	problemService *service.ProblemService
	logger         *logger.Logger
}

// NewProblemsHandler создает новый обработчик проблем
func NewProblemsHandler(problemService *service.ProblemService, logger *logger.Logger) *ProblemsHandler {
	return &ProblemsHandler{
		problemService: problemService,
		logger:         logger,
	}
}

// GetProblems godoc
// @Summary Получить список проблемных отправлений
// @Description Возвращает список проблем с отправлениями с возможностью фильтрации
// @Tags admin-logistics
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param status query string false "Фильтр по статусу проблемы"
// @Param severity query string false "Фильтр по критичности (low, medium, high, critical)"
// @Param problem_type query string false "Фильтр по типу проблемы"
// @Param assigned_to query int false "Фильтр по назначенному пользователю"
// @Param page query int false "Номер страницы" default(1)
// @Param limit query int false "Количество элементов на странице" default(20)
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=map[string]interface{}} "List of problems"
// @Failure 401 {object} backend_pkg_utils.ErrorResponseSwag "Unauthorized"
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/admin/logistics/problems [get]
func (h *ProblemsHandler) GetProblems(c *fiber.Ctx) error {
	// Проверка прав доступа
	userID := c.Locals("user_id")
	if userID == nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "auth.unauthorized")
	}

	// Парсим фильтры
	filter := service.ProblemsFilter{
		Page:  c.QueryInt("page", 1),
		Limit: c.QueryInt("limit", 20),
	}

	if status := c.Query("status"); status != "" {
		filter.Status = &status
	}

	if severity := c.Query("severity"); severity != "" {
		filter.Severity = &severity
	}

	if problemType := c.Query("problem_type"); problemType != "" {
		filter.ProblemType = &problemType
	}

	if assignedTo := c.QueryInt("assigned_to", 0); assignedTo > 0 {
		filter.AssignedTo = &assignedTo
	}

	// Получаем проблемы
	problems, total, err := h.problemService.GetProblems(c.Context(), filter)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "logistics.problems.error")
	}

	// Формируем ответ
	response := map[string]interface{}{
		"problems": problems,
		"total":    total,
		"page":     filter.Page,
		"limit":    filter.Limit,
		"pages":    (total + filter.Limit - 1) / filter.Limit,
	}

	return utils.SuccessResponse(c, response)
}

// CreateProblem godoc
// @Summary Создать новую проблему
// @Description Создает новую запись о проблеме с отправлением
// @Tags admin-logistics
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param body body logistics.ProblemShipment true "Данные проблемы"
// @Success 201 {object} backend_pkg_utils.SuccessResponseSwag{data=logistics.ProblemShipment} "Created problem"
// @Failure 401 {object} backend_pkg_utils.ErrorResponseSwag "Unauthorized"
// @Failure 400 {object} backend_pkg_utils.ErrorResponseSwag "Bad request - invalid data"
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/admin/logistics/problems [post]
func (h *ProblemsHandler) CreateProblem(c *fiber.Ctx) error {
	// Проверка прав доступа
	userID := c.Locals("user_id")
	if userID == nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "auth.unauthorized")
	}

	// Парсим тело запроса
	var problem logistics.ProblemShipment
	if err := c.BodyParser(&problem); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "common.invalid_request_body")
	}

	// Валидация обязательных полей
	if problem.ShipmentID == 0 || problem.ShipmentType == "" || problem.ProblemType == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "logistics.problem.missing_fields")
	}

	// Устанавливаем значения по умолчанию
	if problem.Severity == "" {
		problem.Severity = logistics.SeverityMedium
	}
	if problem.Status == "" {
		problem.Status = logistics.StatusOpen
	}

	// Создаем проблему
	createdProblem, err := h.problemService.CreateProblem(c.Context(), &problem)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "logistics.problem.create_error")
	}

	c.Status(fiber.StatusCreated)
	return utils.SuccessResponse(c, createdProblem)
}

// UpdateProblem godoc
// @Summary Обновить проблему
// @Description Обновляет информацию о проблеме с отправлением
// @Tags admin-logistics
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "ID проблемы"
// @Param body body map[string]interface{} true "Обновляемые поля"
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=logistics.ProblemShipment} "Updated problem"
// @Failure 401 {object} backend_pkg_utils.ErrorResponseSwag "Unauthorized"
// @Failure 400 {object} backend_pkg_utils.ErrorResponseSwag "Bad request - invalid data"
// @Failure 404 {object} backend_pkg_utils.ErrorResponseSwag "Problem not found"
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/admin/logistics/problems/{id} [put]
func (h *ProblemsHandler) UpdateProblem(c *fiber.Ctx) error {
	// Проверка прав доступа
	userID := c.Locals("user_id")
	if userID == nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "auth.unauthorized")
	}

	// Парсим ID проблемы
	problemID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "logistics.invalid_problem_id")
	}

	// Парсим тело запроса
	var updates map[string]interface{}
	if err := c.BodyParser(&updates); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "common.invalid_request_body")
	}

	// Обновляем проблему
	updatedProblem, err := h.problemService.UpdateProblem(c.Context(), problemID, updates)
	if err != nil {
		if err.Error() == problemNotFoundErr {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "logistics.problem_not_found")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "logistics.problem.update_error")
	}

	return utils.SuccessResponse(c, updatedProblem)
}

// ResolveProblem godoc
// @Summary Решить проблему
// @Description Отмечает проблему как решенную с указанием резолюции
// @Tags admin-logistics
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "ID проблемы"
// @Param body body map[string]string true "Резолюция проблемы"
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag "Problem resolved"
// @Failure 401 {object} backend_pkg_utils.ErrorResponseSwag "Unauthorized"
// @Failure 400 {object} backend_pkg_utils.ErrorResponseSwag "Bad request - invalid data"
// @Failure 404 {object} backend_pkg_utils.ErrorResponseSwag "Problem not found"
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/admin/logistics/problems/{id}/resolve [post]
func (h *ProblemsHandler) ResolveProblem(c *fiber.Ctx) error {
	// Проверка прав доступа
	userID := c.Locals("user_id")
	if userID == nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "auth.unauthorized")
	}

	userIDInt := userID.(int)

	// Парсим ID проблемы
	problemID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "logistics.invalid_problem_id")
	}

	// Парсим тело запроса
	var request struct {
		Resolution string `json:"resolution"`
	}

	if err := c.BodyParser(&request); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "common.invalid_request_body")
	}

	if request.Resolution == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "logistics.resolution_required")
	}

	// Решаем проблему
	err = h.problemService.ResolveProblem(c.Context(), problemID, request.Resolution, userIDInt)
	if err != nil {
		if err.Error() == problemNotFoundErr {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "logistics.problem_not_found")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "logistics.problem.resolve_error")
	}

	// Добавляем комментарий о решении проблемы
	_, err = h.problemService.AddProblemComment(c.Context(), problemID, userIDInt,
		"Проблема решена: "+request.Resolution,
		"resolution",
		map[string]interface{}{
			"resolution":  request.Resolution,
			"resolved_by": userIDInt,
		})
	if err != nil {
		// Логируем ошибку, но не прерываем выполнение
		// так как основное действие уже выполнено
		h.logger.Error("Failed to add problem history for resolution: %v", err)
	}

	return utils.SuccessResponse(c, map[string]interface{}{
		"problem_id":  problemID,
		"resolved_at": time.Now(),
	})
}

// AssignProblem godoc
// @Summary Назначить проблему пользователю
// @Description Назначает проблему на конкретного пользователя для решения
// @Tags admin-logistics
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "ID проблемы"
// @Param body body map[string]int true "ID пользователя для назначения"
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag "Problem assigned"
// @Failure 401 {object} backend_pkg_utils.ErrorResponseSwag "Unauthorized"
// @Failure 400 {object} backend_pkg_utils.ErrorResponseSwag "Bad request - invalid data"
// @Failure 404 {object} backend_pkg_utils.ErrorResponseSwag "Problem not found"
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/admin/logistics/problems/{id}/assign [post]
func (h *ProblemsHandler) AssignProblem(c *fiber.Ctx) error {
	// Проверка прав доступа
	userID := c.Locals("user_id")
	if userID == nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "auth.unauthorized")
	}

	// Парсим ID проблемы
	problemID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "logistics.invalid_problem_id")
	}

	// Парсим тело запроса
	var request struct {
		AssignTo int `json:"assign_to"`
	}

	if err := c.BodyParser(&request); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "common.invalid_request_body")
	}

	if request.AssignTo == 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "logistics.assign_user_required")
	}

	// Назначаем проблему
	err = h.problemService.AssignProblem(c.Context(), problemID, request.AssignTo)
	if err != nil {
		if err.Error() == problemNotFoundErr {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "logistics.problem_not_found")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "logistics.problem.assign_error")
	}

	// Добавляем комментарий о назначении проблемы
	_, err = h.problemService.AddProblemComment(c.Context(), problemID, userID.(int),
		fmt.Sprintf("Проблема назначена администратору с ID %d", request.AssignTo),
		"assignment",
		map[string]interface{}{
			"assigned_to": request.AssignTo,
			"assigned_by": userID.(int),
		})
	if err != nil {
		// Логируем ошибку, но не прерываем выполнение
		h.logger.Error("Failed to add problem history for assignment: %v", err)
	}

	return utils.SuccessResponse(c, map[string]interface{}{
		"problem_id":  problemID,
		"assigned_to": request.AssignTo,
	})
}

// AddProblemComment godoc
// @Summary Добавить комментарий к проблеме
// @Description Добавляет комментарий к существующей проблеме
// @Tags admin-logistics
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "ID проблемы"
// @Param body body map[string]string true "Текст комментария"
// @Success 201 {object} backend_pkg_utils.SuccessResponseSwag{data=logistics.ProblemComment} "Added comment"
// @Failure 401 {object} backend_pkg_utils.ErrorResponseSwag "Unauthorized"
// @Failure 400 {object} backend_pkg_utils.ErrorResponseSwag "Bad request - invalid data"
// @Failure 404 {object} backend_pkg_utils.ErrorResponseSwag "Problem not found"
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/admin/logistics/problems/{id}/comments [post]
func (h *ProblemsHandler) AddProblemComment(c *fiber.Ctx) error {
	// Проверка прав доступа
	userID := c.Locals("user_id")
	if userID == nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "auth.unauthorized")
	}

	userIDInt := userID.(int)

	// Парсим ID проблемы
	problemID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "logistics.invalid_problem_id")
	}

	// Парсим тело запроса
	var request struct {
		Comment string `json:"comment"`
	}

	if err := c.BodyParser(&request); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "common.invalid_request_body")
	}

	if request.Comment == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "logistics.comment_required")
	}

	// Добавляем комментарий
	comment, err := h.problemService.AddProblemComment(c.Context(), problemID, userIDInt, request.Comment, "comment", map[string]interface{}{})
	if err != nil {
		if err.Error() == problemNotFoundErr {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "logistics.problem_not_found")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "logistics.comment.add_error")
	}

	c.Status(fiber.StatusCreated)
	return utils.SuccessResponse(c, comment)
}

// GetProblemComments godoc
// @Summary Получить комментарии к проблеме
// @Description Возвращает список всех комментариев к проблеме
// @Tags admin-logistics
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "ID проблемы"
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=[]logistics.ProblemComment} "List of comments"
// @Failure 401 {object} backend_pkg_utils.ErrorResponseSwag "Unauthorized"
// @Failure 404 {object} backend_pkg_utils.ErrorResponseSwag "Problem not found"
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/admin/logistics/problems/{id}/comments [get]
func (h *ProblemsHandler) GetProblemComments(c *fiber.Ctx) error {
	// Проверка прав доступа
	userID := c.Locals("user_id")
	if userID == nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "auth.unauthorized")
	}

	// Парсим ID проблемы
	problemID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "logistics.invalid_problem_id")
	}

	// Получаем комментарии
	comments, err := h.problemService.GetProblemComments(c.Context(), problemID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "logistics.comments.get_error")
	}

	return utils.SuccessResponse(c, comments)
}

// GetProblemHistory godoc
// @Summary Получить историю изменений проблемы
// @Description Возвращает историю всех изменений статуса и назначения проблемы
// @Tags admin-logistics
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "ID проблемы"
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=[]logistics.ProblemStatusHistory} "History of changes"
// @Failure 401 {object} backend_pkg_utils.ErrorResponseSwag "Unauthorized"
// @Failure 404 {object} backend_pkg_utils.ErrorResponseSwag "Problem not found"
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/admin/logistics/problems/{id}/history [get]
func (h *ProblemsHandler) GetProblemHistory(c *fiber.Ctx) error {
	// Проверка прав доступа
	userID := c.Locals("user_id")
	if userID == nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "auth.unauthorized")
	}

	// Парсим ID проблемы
	problemID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "logistics.invalid_problem_id")
	}

	// Получаем историю
	history, err := h.problemService.GetProblemHistory(c.Context(), problemID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "logistics.history.get_error")
	}

	return utils.SuccessResponse(c, history)
}

// GetProblemDetails godoc
// @Summary Получить детальную информацию о проблеме
// @Description Возвращает полную информацию о проблеме включая комментарии и историю
// @Tags admin-logistics
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "ID проблемы"
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=logistics.ProblemShipment} "Problem details with comments and history"
// @Failure 401 {object} backend_pkg_utils.ErrorResponseSwag "Unauthorized"
// @Failure 404 {object} backend_pkg_utils.ErrorResponseSwag "Problem not found"
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/admin/logistics/problems/{id}/details [get]
func (h *ProblemsHandler) GetProblemDetails(c *fiber.Ctx) error {
	// Проверка прав доступа
	userID := c.Locals("user_id")
	if userID == nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "auth.unauthorized")
	}

	// Парсим ID проблемы
	problemID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "logistics.invalid_problem_id")
	}

	// Получаем детали проблемы
	problem, err := h.problemService.GetProblemWithDetails(c.Context(), problemID)
	if err != nil {
		if err.Error() == problemNotFoundErr {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "logistics.problem_not_found")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "logistics.problem_details.get_error")
	}

	return utils.SuccessResponse(c, problem)
}
