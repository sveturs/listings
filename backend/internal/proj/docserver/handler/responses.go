// Package handler
// backend/internal/proj/docserver/handler/responses.go
package handler

// DocFile представляет файл или директорию документации
type DocFile struct {
	Name     string    `json:"name" example:"README.md"`
	Path     string    `json:"path" example:"docs/README.md"`
	Type     string    `json:"type" example:"file" enums:"file,directory"`
	Children []DocFile `json:"children,omitempty"`
}

// DocFilesResponse представляет ответ со списком файлов документации
type DocFilesResponse struct {
	Files    []DocFile `json:"files"`
	RootPath string    `json:"rootPath" example:"./docs"`
}

// DocContentResponse представляет ответ с содержимым файла
type DocContentResponse struct {
	Content string `json:"content" example:"# Documentation\n\nThis is the content of the markdown file."`
	Path    string `json:"path" example:"docs/README.md"`
}
