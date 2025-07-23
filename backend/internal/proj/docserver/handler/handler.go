// Package handler
// backend/internal/proj/docserver/handler/handler.go
package handler

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"

	"backend/internal/config"
	"backend/internal/logger"
	"backend/pkg/utils"
)

// NewHandler создает новый экземпляр Handler
func NewHandler(cfg config.DocsConfig) *Handler {
	return &Handler{
		cfg: cfg,
	}
}

type Handler struct {
	cfg config.DocsConfig
}

// DocFile определен в responses.go

// GetDocFiles returns the list of markdown files
// @Summary Get documentation files
// @Description Returns the hierarchical list of markdown files in the documentation directory
// @Tags docs
// @Accept json
// @Produce json
// @Success 200 {object} utils.SuccessResponseSwag{data=DocFilesResponse} "List of documentation files"
// @Failure 500 {object} utils.ErrorResponseSwag "docs.scanError"
// @Router /api/v1/docs/files [get]
func (h *Handler) GetDocFiles(c *fiber.Ctx) error {
	files, err := h.scanDirectory("")
	if err != nil {
		logger.Error().Err(err).Str("rootPath", h.cfg.RootPath).Msg("Failed to scan documentation files")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "docs.scanError")
	}

	response := DocFilesResponse{
		Files:    files,
		RootPath: h.cfg.RootPath,
	}

	return utils.SuccessResponse(c, response)
}

// GetFileContent returns the content of a specific markdown file
// @Summary Get file content
// @Description Returns the content of a specific markdown file
// @Tags docs
// @Accept json
// @Produce json
// @Param path query string true "File path relative to docs root"
// @Success 200 {object} utils.SuccessResponseSwag{data=DocContentResponse} "File content"
// @Failure 400 {object} utils.ErrorResponseSwag "docs.pathRequired or docs.invalidPath or docs.onlyMarkdownAllowed"
// @Failure 404 {object} utils.ErrorResponseSwag "docs.fileNotFound"
// @Failure 500 {object} utils.ErrorResponseSwag "docs.readError"
// @Router /api/v1/docs/content [get]
func (h *Handler) GetFileContent(c *fiber.Ctx) error {
	filePath := c.Query("path")
	if filePath == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "docs.pathRequired")
	}

	// Security: prevent directory traversal
	cleanPath := filepath.Clean(filePath)
	if strings.Contains(cleanPath, "..") {
		logger.Warn().Str("path", filePath).Msg("Attempt to access file with directory traversal")
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "docs.invalidPath")
	}

	fullPath := filepath.Join(h.cfg.RootPath, cleanPath)

	// Additional security: ensure the resolved path is within the root directory
	absRootPath, _ := filepath.Abs(h.cfg.RootPath)
	absFullPath, _ := filepath.Abs(fullPath)
	if !strings.HasPrefix(absFullPath, absRootPath) {
		logger.Warn().Str("path", filePath).Msg("Attempt to access file outside root directory")
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "docs.invalidPath")
	}

	// Check if file exists and is a markdown file
	info, err := os.Stat(fullPath)
	if err != nil || info.IsDir() {
		if err != nil {
			logger.Debug().Err(err).Str("path", cleanPath).Msg("File not found")
		}
		return utils.ErrorResponse(c, fiber.StatusNotFound, "docs.fileNotFound")
	}

	if !strings.HasSuffix(strings.ToLower(fullPath), ".md") {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "docs.onlyMarkdownAllowed")
	}

	content, err := os.ReadFile(fullPath) // #nosec G304 -- path validated above
	if err != nil {
		logger.Error().Err(err).Str("path", cleanPath).Msg("Failed to read file")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "docs.readError")
	}

	response := DocContentResponse{
		Content: string(content),
		Path:    filePath,
	}

	return utils.SuccessResponse(c, response)
}

func (h *Handler) scanDirectory(relativePath string) ([]DocFile, error) {
	fullPath := filepath.Join(h.cfg.RootPath, relativePath)
	entries, err := os.ReadDir(fullPath)
	if err != nil {
		return nil, err
	}

	var files []DocFile
	for _, entry := range entries {
		// Skip hidden files and directories
		if strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		// Skip non-markdown files except directories
		if !entry.IsDir() && !strings.HasSuffix(strings.ToLower(entry.Name()), ".md") {
			continue
		}

		entryPath := filepath.Join(relativePath, entry.Name())

		if entry.IsDir() {
			// Skip certain directories
			skipDirs := []string{
				"node_modules", "vendor", "dist", ".git", ".next",
				"build", "coverage", ".turbo", "out", ".vercel",
				"data", "uploads", "bin", "tmp", "temp", ".idea",
				".vscode", "logs", "__pycache__", ".cache",
			}
			skip := false
			for _, skipDir := range skipDirs {
				if entry.Name() == skipDir {
					skip = true
					break
				}
			}
			if skip {
				continue
			}

			children, err := h.scanDirectory(entryPath)
			if err != nil {
				continue
			}

			// Only include directory if it has markdown files
			if len(children) > 0 {
				files = append(files, DocFile{
					Name:     entry.Name(),
					Path:     entryPath,
					Type:     "directory",
					Children: children,
				})
			}
		} else {
			files = append(files, DocFile{
				Name: entry.Name(),
				Path: entryPath,
				Type: "file",
			})
		}
	}

	return files, nil
}
