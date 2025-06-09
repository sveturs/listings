package handler

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type DocsHandler struct {
	rootPath string
}

func NewDocsHandler(rootPath string) *DocsHandler {
	return &DocsHandler{
		rootPath: rootPath,
	}
}

type DocFile struct {
	Name     string     `json:"name"`
	Path     string     `json:"path"`
	Type     string     `json:"type"` // "file" or "directory"
	Children []DocFile  `json:"children,omitempty"`
}

// GetDocFiles returns the list of markdown files
// @Summary Get documentation files
// @Description Returns the hierarchical list of markdown files
// @Tags docs
// @Accept json
// @Produce json
// @Success 200 {object} map[string][]DocFile
// @Router /api/v1/docs/files [get]
func (h *DocsHandler) GetDocFiles(c *fiber.Ctx) error {
	files, err := h.scanDirectory("")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to scan documentation files",
			"details": err.Error(),
			"rootPath": h.rootPath,
		})
	}

	return c.JSON(fiber.Map{
		"files": files,
		"rootPath": h.rootPath,
	})
}

// GetFileContent returns the content of a specific markdown file
// @Summary Get file content
// @Description Returns the content of a specific markdown file
// @Tags docs
// @Accept json
// @Produce json
// @Param path query string true "File path"
// @Success 200 {object} map[string]string
// @Router /api/v1/docs/content [get]
func (h *DocsHandler) GetFileContent(c *fiber.Ctx) error {
	filePath := c.Query("path")
	if filePath == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "File path is required",
		})
	}

	// Security: prevent directory traversal
	cleanPath := filepath.Clean(filePath)
	if strings.Contains(cleanPath, "..") {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid file path",
		})
	}

	fullPath := filepath.Join(h.rootPath, cleanPath)
	
	// Check if file exists and is a markdown file
	info, err := os.Stat(fullPath)
	if err != nil || info.IsDir() {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "File not found",
		})
	}

	if !strings.HasSuffix(strings.ToLower(fullPath), ".md") {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Only markdown files are allowed",
		})
	}

	content, err := os.ReadFile(fullPath)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to read file",
		})
	}

	return c.JSON(fiber.Map{
		"content": string(content),
		"path":    filePath,
	})
}

func (h *DocsHandler) scanDirectory(relativePath string) ([]DocFile, error) {
	fullPath := filepath.Join(h.rootPath, relativePath)
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

// RegisterRoutes registers all documentation routes
func (h *DocsHandler) RegisterRoutes(router fiber.Router) {
	docs := router.Group("/docs")
	// Docs routes are public - no auth required
	docs.Get("/files", h.GetDocFiles)
	docs.Get("/content", h.GetFileContent)
}