package handler

import (
	"context"

	"github.com/gofiber/fiber/v2"

	"backend/pkg/utils"
)

// GetIndexInfo возвращает информацию об индексе
// @Summary      Get search index information
// @Description  Returns detailed information about the search index
// @Tags         admin-search
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} backend_pkg_utils.SuccessResponseSwag{data=backend_internal_proj_search_admin_service.IndexInfo}
// @Failure      401 {object} backend_pkg_utils.ErrorResponseSwag
// @Failure      403 {object} backend_pkg_utils.ErrorResponseSwag
// @Failure      500 {object} backend_pkg_utils.ErrorResponseSwag
// @Router       /api/v1/admin/search/index/info [get]
func (h *Handler) GetIndexInfo(c *fiber.Ctx) error {
	info, err := h.service.GetIndexInfo(c.Context())
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "admin.search.indexInfoError")
	}

	return utils.SuccessResponse(c, info)
}

// GetIndexStatistics возвращает статистику индекса
// @Summary      Get search index statistics
// @Description  Returns statistics about indexed documents
// @Tags         admin-search
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} backend_pkg_utils.SuccessResponseSwag{data=backend_internal_proj_search_admin_service.IndexStatistics}
// @Failure      401 {object} backend_pkg_utils.ErrorResponseSwag
// @Failure      403 {object} backend_pkg_utils.ErrorResponseSwag
// @Failure      500 {object} backend_pkg_utils.ErrorResponseSwag
// @Router       /api/v1/admin/search/index/statistics [get]
func (h *Handler) GetIndexStatistics(c *fiber.Ctx) error {
	stats, err := h.service.GetIndexStatistics(c.Context())
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "admin.search.statisticsError")
	}

	return utils.SuccessResponse(c, stats)
}

// SearchIndexedDocuments поиск проиндексированных документов
// @Summary      Search indexed documents
// @Description  Search and filter indexed documents
// @Tags         admin-search
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        query        query    string false "Search query"
// @Param        type         query    string false "Document type (listing/product)"
// @Param        category_id  query    int    false "Category ID"
// @Param        page         query    int    false "Page number" default(1)
// @Param        limit        query    int    false "Items per page" default(20)
// @Success      200 {object} backend_pkg_utils.SuccessResponseSwag{data=[]backend_internal_proj_search_admin_service.IndexedDocument}
// @Failure      401 {object} backend_pkg_utils.ErrorResponseSwag
// @Failure      403 {object} backend_pkg_utils.ErrorResponseSwag
// @Failure      500 {object} backend_pkg_utils.ErrorResponseSwag
// @Router       /api/v1/admin/search/index/documents [get]
func (h *Handler) SearchIndexedDocuments(c *fiber.Ctx) error {
	// Параметры запроса
	query := c.Query("query", "")
	docType := c.Query("type", "")
	limit := c.QueryInt("limit", 20)

	result, err := h.service.SearchIndexedDocuments(c.Context(), query, docType, limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "admin.search.documentsError")
	}

	return utils.SuccessResponse(c, result)
}

// ReindexDocuments запускает переиндексацию документов
// @Summary      Reindex search documents
// @Description  Trigger reindexing of all marketplace documents
// @Tags         admin-search
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} backend_pkg_utils.SuccessResponseSwag
// @Failure      401 {object} backend_pkg_utils.ErrorResponseSwag
// @Failure      403 {object} backend_pkg_utils.ErrorResponseSwag
// @Failure      500 {object} backend_pkg_utils.ErrorResponseSwag
// @Router       /api/v1/admin/search/index/reindex [post]
func (h *Handler) ReindexDocuments(c *fiber.Ctx) error {
	h.logger.Info("=== REINDEX HANDLER CALLED ===")

	// Получаем тип документов для переиндексации
	var req struct {
		Type string `json:"type"` // "listing", "product" или пустая строка для всех
	}

	if err := c.BodyParser(&req); err != nil {
		// Если не удалось распарсить, переиндексируем все
		req.Type = ""
	}

	h.logger.Info("Request type: %s", req.Type)

	// ВРЕМЕННО: Запускаем синхронно для отладки
	// TODO: Вернуть асинхронное выполнение после fix логирования
	h.logger.Info("Starting reindexing synchronously for debugging")
	if err := h.service.ReindexDocuments(context.Background(), req.Type); err != nil {
		h.logger.Error("Failed to reindex documents: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "search.reindexError")
	}
	h.logger.Info("Reindexing completed successfully")

	return utils.SuccessResponse(c, map[string]string{
		"message": "Reindexing completed",
	})
}
