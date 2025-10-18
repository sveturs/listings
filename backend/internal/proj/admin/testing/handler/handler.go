// Package handler implements HTTP handlers for testing module
// backend/internal/proj/admin/testing/handler/handler.go
package handler

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	authMiddleware "github.com/sveturs/auth/pkg/http/fiber/middleware"

	"backend/internal/proj/admin/testing/domain"
	"backend/internal/proj/admin/testing/service"
	"backend/pkg/utils"
)

// Handler handles testing module HTTP requests
type Handler struct {
	testRunner  *service.TestRunner
	logger      zerolog.Logger
	jwtParserMW fiber.Handler
}

// NewHandler creates new testing handler instance
func NewHandler(
	testRunner *service.TestRunner,
	jwtParserMW fiber.Handler,
	logger zerolog.Logger,
) *Handler {
	return &Handler{
		testRunner:  testRunner,
		jwtParserMW: jwtParserMW,
		logger:      logger.With().Str("component", "testing_handler").Logger(),
	}
}

// RegisterRoutes registers all testing routes
func (h *Handler) RegisterRoutes(app *fiber.App) {
	// All testing endpoints require admin role
	tests := app.Group("/api/v1/admin/tests", h.jwtParserMW, authMiddleware.RequireAuthString("admin"))

	tests.Post("/run", h.RunTest)
	tests.Get("/suites", h.GetTestSuites)
	tests.Get("/runs", h.ListTestRuns)
	tests.Get("/runs/latest", h.GetLatestTestRun)
	tests.Get("/runs/:id", h.GetTestRunDetail)
	tests.Get("/runs/:id/status", h.GetTestRunStatus)
	tests.Delete("/runs/:id", h.CancelTestRun)
}

// GetPrefix returns routing prefix for this handler
func (h *Handler) GetPrefix() string {
	return "/api/v1/admin/tests"
}

// RunTest godoc
// @Summary Run test suite
// @Description Initiates execution of specified test suite
// @Tags testing
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body domain.RunTestRequest true "Test run request"
// @Success 200 {object} domain.RunTestResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/v1/admin/tests/run [post]
func (h *Handler) RunTest(c *fiber.Ctx) error {
	var req domain.RunTestRequest
	if err := c.BodyParser(&req); err != nil {
		h.logger.Error().Err(err).Msg("Failed to parse request body")
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	// Validate request
	if req.TestSuite == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Test suite is required")
	}

	// Get user ID from context
	userID, ok := authMiddleware.GetUserID(c)
	if !ok {
		h.logger.Error().Msg("Failed to get user ID from context")
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "User not authenticated")
	}

	// Run test suite (or specific test if test_name provided)
	testRun, err := h.testRunner.RunTestSuite(c.Context(), req.TestSuite, req.TestName, userID, req.Parallel)
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("suite", req.TestSuite).
			Str("test_name", req.TestName).
			Msg("Failed to run test suite")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to initiate test run")
	}

	message := "Test suite execution initiated"
	if req.TestName != "" {
		message = fmt.Sprintf("Test '%s' execution initiated", req.TestName)
	}

	return c.Status(fiber.StatusOK).JSON(domain.RunTestResponse{
		TestRunID: testRun.ID,
		RunUUID:   testRun.RunUUID,
		Status:    string(testRun.Status),
		Message:   message,
	})
}

// GetTestSuites godoc
// @Summary Get available test suites
// @Description Returns list of all available test suites
// @Tags testing
// @Produce json
// @Security BearerAuth
// @Success 200 {array} domain.TestSuite
// @Failure 401 {object} utils.ErrorResponse
// @Router /api/v1/admin/tests/suites [get]
func (h *Handler) GetTestSuites(c *fiber.Ctx) error {
	suites := h.testRunner.GetAvailableTestSuites()
	return c.Status(fiber.StatusOK).JSON(suites)
}

// GetLatestTestRun godoc
// @Summary Get latest test run
// @Description Returns the most recent test run with details
// @Tags testing
// @Produce json
// @Security BearerAuth
// @Success 200 {object} domain.TestRunDetail
// @Failure 401 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/v1/admin/tests/runs/latest [get]
func (h *Handler) GetLatestTestRun(c *fiber.Ctx) error {
	testRuns, err := h.testRunner.ListTestRuns(c.Context(), 1, 0)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get latest test run")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve latest test run")
	}

	if len(testRuns) == 0 {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "No test runs found")
	}

	// Get full details of the latest run
	detail, err := h.testRunner.GetTestRunDetail(c.Context(), testRuns[0].ID)
	if err != nil {
		h.logger.Error().Err(err).Int64("id", testRuns[0].ID).Msg("Failed to get test run detail")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve test run details")
	}

	return c.Status(fiber.StatusOK).JSON(detail)
}

// ListTestRuns godoc
// @Summary List test runs
// @Description Returns paginated list of test runs
// @Tags testing
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Limit" default(20)
// @Param offset query int false "Offset" default(0)
// @Success 200 {array} domain.TestRun
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/v1/admin/tests/runs [get]
func (h *Handler) ListTestRuns(c *fiber.Ctx) error {
	limit := c.QueryInt("limit", 20)
	offset := c.QueryInt("offset", 0)

	if limit < 1 || limit > 100 {
		limit = 20
	}

	testRuns, err := h.testRunner.ListTestRuns(c.Context(), limit, offset)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to list test runs")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve test runs")
	}

	if testRuns == nil {
		testRuns = []*domain.TestRun{}
	}

	return c.Status(fiber.StatusOK).JSON(testRuns)
}

// GetTestRunDetail godoc
// @Summary Get test run details
// @Description Returns detailed information about test run including results and logs
// @Tags testing
// @Produce json
// @Security BearerAuth
// @Param id path int true "Test Run ID"
// @Success 200 {object} domain.TestRunDetail
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/v1/admin/tests/runs/{id} [get]
func (h *Handler) GetTestRunDetail(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid test run ID")
	}

	detail, err := h.testRunner.GetTestRunDetail(c.Context(), id)
	if err != nil {
		h.logger.Error().Err(err).Int64("id", id).Msg("Failed to get test run detail")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve test run details")
	}

	if detail == nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Test run not found")
	}

	return c.Status(fiber.StatusOK).JSON(detail)
}

// GetTestRunStatus godoc
// @Summary Get test run status
// @Description Returns current status of test run (for polling)
// @Tags testing
// @Produce json
// @Security BearerAuth
// @Param id path int true "Test Run ID"
// @Success 200 {object} domain.TestRun
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/v1/admin/tests/runs/{id}/status [get]
func (h *Handler) GetTestRunStatus(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid test run ID")
	}

	testRun, err := h.testRunner.GetTestRunStatus(c.Context(), id)
	if err != nil {
		h.logger.Error().Err(err).Int64("id", id).Msg("Failed to get test run status")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve test run status")
	}

	if testRun == nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Test run not found")
	}

	return c.Status(fiber.StatusOK).JSON(testRun)
}

// CancelTestRun godoc
// @Summary Cancel test run
// @Description Cancels running test execution
// @Tags testing
// @Produce json
// @Security BearerAuth
// @Param id path int true "Test Run ID"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/v1/admin/tests/runs/{id} [delete]
func (h *Handler) CancelTestRun(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid test run ID")
	}

	if err := h.testRunner.CancelTestRun(id); err != nil {
		h.logger.Error().Err(err).Int64("id", id).Msg("Failed to cancel test run")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to cancel test run")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Test run cancelled successfully",
	})
}
