package service

import (
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"backend/internal/config"
	"backend/internal/domain/models"
	"backend/internal/logger"
	"backend/internal/storage"
	"backend/internal/storage/filestorage"
	"backend/pkg/utils"
)

// ChatAttachmentServiceInterface определяет методы для работы с вложениями чата
type ChatAttachmentServiceInterface interface {
	UploadAttachments(ctx context.Context, messageID int, files []*multipart.FileHeader) ([]*models.ChatAttachment, error)
	GetAttachment(ctx context.Context, attachmentID int) (*models.ChatAttachment, error)
	GetMessageAttachments(ctx context.Context, messageID int) ([]*models.ChatAttachment, error)
	DeleteAttachment(ctx context.Context, attachmentID int, userID int) error
	ValidateFile(file *multipart.FileHeader, fileUploadConfig config.FileUploadConfig) error
}

// ChatAttachmentService реализует сервис для работы с вложениями
type ChatAttachmentService struct {
	storage          storage.Storage
	fileStorage      filestorage.FileStorageInterface
	fileUploadConfig config.FileUploadConfig
}

// NewChatAttachmentService создает новый экземпляр сервиса
func NewChatAttachmentService(storage storage.Storage, fileStorage filestorage.FileStorageInterface, config config.FileUploadConfig) *ChatAttachmentService {
	if fileStorage == nil {
		log.Printf("WARNING: fileStorage is nil in NewChatAttachmentService")
	}

	// Создаем отдельное хранилище для файлов чата
	chatFileStorage := createChatFileStorage(fileStorage)

	return &ChatAttachmentService{
		storage:          storage,
		fileStorage:      chatFileStorage,
		fileUploadConfig: config,
	}
}

// createChatFileStorage создает файловое хранилище для чата
func createChatFileStorage(defaultStorage filestorage.FileStorageInterface) filestorage.FileStorageInterface {
	// Пока используем то же хранилище, но с модифицированными путями
	// В будущем можно создать отдельный bucket
	return &chatFileStorageWrapper{
		baseStorage: defaultStorage,
		bucketName:  bucketChatFiles,
	}
}

const (
	// File extensions for attachment validation
	extensionJPG  = ".jpg"
	extensionJPEG = ".jpeg"
	extensionPNG  = ".png"
	extensionGIF  = ".gif"
	extensionWEBP = ".webp"
	extensionMP4  = ".mp4"
	extensionWEBM = ".webm"
	extensionPDF  = ".pdf"
	extensionDOC  = ".doc"
	extensionDOCX = ".docx"
	extensionTXT  = ".txt"

	// URL paths
	listingsPath  = "/listings/"
	chatFilesPath = "/chat-files/"

	// Storage configuration
	storageMinio    = "minio"
	bucketChatFiles = "chat-files"
)

// chatFileStorageWrapper обертка для работы с файлами чата
type chatFileStorageWrapper struct {
	baseStorage filestorage.FileStorageInterface
	bucketName  string
}

func (w *chatFileStorageWrapper) UploadFile(ctx context.Context, objectName string, reader io.Reader, size int64, contentType string) (string, error) {
	if w.baseStorage == nil {
		return "", fmt.Errorf("base storage is not initialized")
	}

	// Загружаем файл через базовое хранилище
	publicURL, err := w.baseStorage.UploadFile(ctx, objectName, reader, size, contentType)
	if err != nil {
		return "", err
	}

	// Заменяем путь на правильный для chat-files
	publicURL = strings.Replace(publicURL, listingsPath, chatFilesPath, 1)
	return publicURL, nil
}

func (w *chatFileStorageWrapper) DeleteFile(ctx context.Context, objectName string) error {
	return w.baseStorage.DeleteFile(ctx, objectName)
}

func (w *chatFileStorageWrapper) GetURL(ctx context.Context, objectName string) (string, error) {
	url, err := w.baseStorage.GetURL(ctx, objectName)
	if err != nil {
		return "", err
	}
	// Заменяем путь на правильный для chat-files
	url = strings.Replace(url, listingsPath, chatFilesPath, 1)
	return url, nil
}

func (w *chatFileStorageWrapper) GetPresignedURL(ctx context.Context, objectName string, expiry time.Duration) (string, error) {
	return w.baseStorage.GetPresignedURL(ctx, objectName, expiry)
}

func (w *chatFileStorageWrapper) GetFile(ctx context.Context, objectName string) (io.ReadCloser, error) {
	return w.baseStorage.GetFile(ctx, objectName)
}

// UploadAttachments загружает вложения для сообщения
func (s *ChatAttachmentService) UploadAttachments(ctx context.Context, messageID int, files []*multipart.FileHeader) ([]*models.ChatAttachment, error) {
	logger.Debug().
		Int("messageID", messageID).
		Int("filesCount", len(files)).
		Msg("ChatAttachmentService.UploadAttachments")
	var attachments []*models.ChatAttachment

	for _, fileHeader := range files {
		// Валидация файла
		if err := s.ValidateFile(fileHeader, s.fileUploadConfig); err != nil {
			return nil, fmt.Errorf("file validation error: %w", err)
		}

		// Открываем файл
		file, err := fileHeader.Open()
		if err != nil {
			return nil, fmt.Errorf("error opening file: %w", err)
		}
		defer func() {
			if err := file.Close(); err != nil {
				logger.Warn().Err(err).Msg("Failed to close file")
			}
		}()

		// Определяем тип файла
		fileType := s.getFileType(fileHeader.Header.Get("Content-Type"))

		// Генерируем путь для сохранения
		now := time.Now()
		fileName := fmt.Sprintf("%d_%d_%s", messageID, now.Unix(), fileHeader.Filename)
		objectPath := fmt.Sprintf("%s/%d/%02d/%02d/%s",
			fileType, now.Year(), now.Month(), now.Day(), fileName)

		// Загружаем файл в хранилище
		publicURL, err := s.fileStorage.UploadFile(ctx, objectPath, file, fileHeader.Size, fileHeader.Header.Get("Content-Type"))
		if err != nil {
			return nil, fmt.Errorf("error uploading file: %w", err)
		}

		// Создаем запись в БД
		attachment := &models.ChatAttachment{
			MessageID:     messageID,
			FileType:      fileType,
			FilePath:      objectPath,
			FileName:      fileHeader.Filename,
			FileSize:      fileHeader.Size,
			ContentType:   fileHeader.Header.Get("Content-Type"),
			StorageType:   storageMinio,
			StorageBucket: bucketChatFiles,
			PublicURL:     publicURL,
		}

		// Сохраняем в БД
		if err := s.storage.CreateChatAttachment(ctx, attachment); err != nil {
			// Если не удалось сохранить в БД, удаляем файл из хранилища
			if err := s.fileStorage.DeleteFile(ctx, objectPath); err != nil {
				logger.Warn().Err(err).Str("objectPath", objectPath).Msg("Failed to delete file after DB error")
			}
			return nil, fmt.Errorf("error saving attachment to database: %w", err)
		}

		attachments = append(attachments, attachment)
	}

	// Обновляем счетчик вложений в сообщении
	if err := s.storage.UpdateMessageAttachmentsCount(ctx, messageID, len(attachments)); err != nil {
		// Логируем ошибку, но не возвращаем её, так как файлы уже загружены
		logger.Error().
			Err(err).
			Int("messageID", messageID).
			Int("attachmentsCount", len(attachments)).
			Msg("Error updating message attachments count")
	}

	return attachments, nil
}

// GetAttachment получает информацию о вложении
func (s *ChatAttachmentService) GetAttachment(ctx context.Context, attachmentID int) (*models.ChatAttachment, error) {
	return s.storage.GetChatAttachment(ctx, attachmentID)
}

// GetMessageAttachments получает все вложения сообщения
func (s *ChatAttachmentService) GetMessageAttachments(ctx context.Context, messageID int) ([]*models.ChatAttachment, error) {
	return s.storage.GetMessageAttachments(ctx, messageID)
}

// DeleteAttachment удаляет вложение
func (s *ChatAttachmentService) DeleteAttachment(ctx context.Context, attachmentID int, userID int) error {
	// Получаем информацию о вложении
	attachment, err := s.storage.GetChatAttachment(ctx, attachmentID)
	if err != nil {
		return fmt.Errorf("error getting attachment: %w", err)
	}

	// Проверяем права доступа
	message, err := s.storage.GetMessageByID(ctx, attachment.MessageID)
	if err != nil {
		return fmt.Errorf("error getting message: %w", err)
	}

	if message.SenderID != userID {
		return fmt.Errorf("permission denied")
	}

	// Удаляем файл из хранилища
	if err := s.fileStorage.DeleteFile(ctx, attachment.FilePath); err != nil {
		// Логируем ошибку, но продолжаем удаление из БД
		logger.Error().
			Err(err).
			Str("filePath", attachment.FilePath).
			Int("attachmentID", attachmentID).
			Msg("Error deleting file from storage")
	}

	// Удаляем запись из БД
	if err := s.storage.DeleteChatAttachment(ctx, attachmentID); err != nil {
		return fmt.Errorf("error deleting attachment from database: %w", err)
	}

	// Обновляем счетчик вложений
	attachments, err := s.storage.GetMessageAttachments(ctx, attachment.MessageID)
	if err == nil {
		if err := s.storage.UpdateMessageAttachmentsCount(ctx, attachment.MessageID, len(attachments)); err != nil {
			logger.Warn().Err(err).Int("messageID", attachment.MessageID).Msg("Failed to update message attachments count")
		}
	}

	return nil
}

// ValidateFile проверяет файл на соответствие ограничениям
func (s *ChatAttachmentService) ValidateFile(file *multipart.FileHeader, config config.FileUploadConfig) error {
	// Санитизация и валидация имени файла
	sanitizedName, err := utils.ValidateFileName(file.Filename)
	if err != nil {
		return fmt.Errorf("invalid filename: %w", err)
	}
	file.Filename = sanitizedName

	contentType := file.Header.Get("Content-Type")

	// Если Content-Type пустой, пытаемся определить по расширению
	if contentType == "" {
		ext := strings.ToLower(filepath.Ext(file.Filename))
		switch ext {
		case extensionJPG, extensionJPEG:
			contentType = "image/jpeg"
		case extensionPNG:
			contentType = "image/png"
		case extensionGIF:
			contentType = "image/gif"
		case extensionWEBP:
			contentType = "image/webp"
		case extensionMP4:
			contentType = "video/mp4"
		case extensionWEBM:
			contentType = "video/webm"
		case extensionPDF:
			contentType = "application/pdf"
		case extensionDOC:
			contentType = "application/msword"
		case extensionDOCX:
			contentType = "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
		case extensionTXT:
			contentType = "text/plain"
		default:
			return fmt.Errorf("unable to determine file type for extension: %s", ext)
		}
	}

	// Используем централизованную валидацию
	fileType, err := utils.ValidateFileType(contentType)
	if err != nil {
		return err
	}

	// Валидация размера файла
	if err := utils.ValidateFileSize(fileType, file.Size); err != nil {
		return err
	}

	return nil
}

// getFileType определяет тип файла по MIME типу
func (s *ChatAttachmentService) getFileType(contentType string) string {
	if strings.HasPrefix(contentType, "image/") {
		return models.FileTypeImage
	}
	if strings.HasPrefix(contentType, "video/") {
		return models.FileTypeVideo
	}
	// Все остальное считаем документом
	return models.FileTypeDocument
}


// GenerateVideoThumbnail генерирует превью для видео (заглушка)
func (s *ChatAttachmentService) GenerateVideoThumbnail(videoPath string) (string, error) {
	// TODO: Реализовать генерацию превью для видео
	return "", nil
}

// ExtractDocumentMetadata извлекает метаданные документа (заглушка)
func (s *ChatAttachmentService) ExtractDocumentMetadata(documentPath string) (map[string]interface{}, error) {
	// TODO: Реализовать извлечение метаданных
	return nil, nil
}
