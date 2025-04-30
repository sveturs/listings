// backend/internal/proj/storefront/service/image_processor.go
package service

import (
	"archive/zip"
	"backend/internal/domain/models"
	"bytes"
	"context"
	"fmt"
	"image"
	//	"image/color"
	"image/jpeg"
	"image/png"
	//	"image/draw"
	"image/gif"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	//	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	//"github.com/disintegration/imaging"
	//"github.com/golang/freetype/truetype"
	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
)

// ProcessImportImages обрабатывает изображения для импортируемого объявления
func (s *StorefrontService) ProcessImportImages(
	ctx context.Context,
	listingID int,
	imagesStr string,
	zipReader *zip.Reader,
) error {
	// Перенаправление на новый метод асинхронной обработки для надежности
	s.ProcessImportImagesAsync(ctx, listingID, imagesStr, zipReader)
	return nil
}

// Асинхронная обработка изображений
func (s *StorefrontService) ProcessImportImagesAsync(ctx context.Context, listingID int, imagesStr string, zipReader *zip.Reader) {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
		defer cancel()

		log.Printf("Начата асинхронная обработка изображений для листинга %d", listingID)
		imagesList := strings.Split(imagesStr, ",")

		// Создаем семафор для ограничения количества параллельных горутин
		sem := make(chan struct{}, 5)
		var wg sync.WaitGroup

		totalImages := 0
		processedImages := 0
		failedImages := 0
		successfullyProcessed := false // Флаг для отслеживания успешной обработки хотя бы одного изображения

		for i, imagePath := range imagesList {
			imagePath = strings.TrimSpace(imagePath)
			if imagePath == "" {
				continue
			}

			totalImages++
			wg.Add(1)
			sem <- struct{}{}

			// Используем детерминированное имя файла
			fileName := fmt.Sprintf("%d_image_%d.jpg", listingID, i)

			go func(idx int, path string, filename string, isMain bool) {
				defer func() {
					<-sem
					wg.Done()
				}()

				// Проверяем изображение перед добавлением в БД
				var imgData []byte
				var contentType string
				var err error

				// Загружаем и проверяем изображение
				if strings.HasPrefix(strings.ToLower(path), "http://") ||
					strings.HasPrefix(strings.ToLower(path), "https://") {
					imgData, contentType, err = s.downloadImage(path)
				} else if zipReader != nil {
					imgData, contentType, err = s.extractImageFromZip(zipReader, path)
				} else {
					imgData, contentType, err = s.readLocalImage(path)
				}

				// Если не удалось получить изображение, просто логируем ошибку и выходим
				if err != nil || imgData == nil || len(imgData) == 0 {
					failedImages++
					log.Printf("Не удалось получить изображение %s для листинга %d: %v", path, listingID, err)

					// Удаляем запись об изображении из БД, если она существует
					s.storage.Exec(ctx,
						"DELETE FROM marketplace_images WHERE listing_id = $1 AND file_path = $2",
						listingID, filename)
					return
				}

				// Обрабатываем изображение (уменьшаем и добавляем водяной знак)
				processedData, err := s.processImage(imgData, contentType)
				if err != nil {
					failedImages++
					log.Printf("Ошибка при обработке изображения %s: %v", path, err)

					// Удаляем запись об изображении из БД, если она существует
					s.storage.Exec(ctx,
						"DELETE FROM marketplace_images WHERE listing_id = $1 AND file_path = $2",
						listingID, filename)
					return
				}

				// Сохраняем файл в MinIO
				fileStorage := s.storage.FileStorage()
				if fileStorage == nil {
					failedImages++
					log.Printf("Ошибка: хранилище файлов не инициализировано")
					return
				}

				// Создаем путь для объекта в MinIO
				objectName := fmt.Sprintf("%d/%s", listingID, filename)

				// Загружаем файл в MinIO
				publicURL, err := fileStorage.UploadFile(ctx, objectName, bytes.NewReader(processedData), int64(len(processedData)), contentType)
				if err != nil {
					failedImages++
					log.Printf("Ошибка при загрузке изображения в MinIO: %v", err)

					// Удаляем запись об изображении из БД, если она существует
					s.storage.Exec(ctx,
						"DELETE FROM marketplace_images WHERE listing_id = $1 AND file_path = $2",
						listingID, filename)
					return
				}

				log.Printf("Изображение успешно загружено в MinIO. objectName=%s, publicURL=%s", objectName, publicURL)

				// Обновляем или создаем запись об изображении в БД
				var imageID int
				err = s.storage.QueryRow(ctx,
					"SELECT id FROM marketplace_images WHERE listing_id = $1 AND file_path = $2",
					listingID, filename).Scan(&imageID)

				if err == nil {
					// Нашли запись, обновляем ее
					_, err = s.storage.Exec(ctx,
						"UPDATE marketplace_images SET file_size = $1, content_type = $2, storage_type = $3, storage_bucket = $4, public_url = $5, file_path = $6 WHERE id = $7",
						len(processedData), contentType, "minio", "listings", publicURL, objectName, imageID)
				} else {
					// Запись не найдена, создаем новую
					image := &models.MarketplaceImage{
						ListingID:     listingID,
						FilePath:      objectName,
						FileName:      strings.Split(path, "/")[len(strings.Split(path, "/"))-1],
						FileSize:      len(processedData),
						ContentType:   contentType,
						IsMain:        isMain,
						StorageType:   "minio",
						StorageBucket: "listings",
						PublicURL:     publicURL,
					}

					_, err = s.storage.AddListingImage(ctx, image)
				}

				if err != nil {
					log.Printf("Ошибка при обновлении/добавлении информации об изображении в БД: %v", err)
					failedImages++
				} else {
					processedImages++
					successfullyProcessed = true // Отмечаем, что хотя бы одно изображение успешно обработано
					log.Printf("Изображение %s успешно обработано и сохранено для объявления %d", filename, listingID)
				}
			}(i, imagePath, fileName, i == 0)
		}

		wg.Wait()
		log.Printf("Завершена асинхронная обработка изображений для листинга %d. Обработано: %d из %d, ошибок: %d",
			listingID, processedImages, totalImages, failedImages)

		// Если хотя бы одно изображение было успешно обработано, обновляем индекс в OpenSearch
		if successfullyProcessed {
			log.Printf("Переиндексация объявления %d в OpenSearch...", listingID)

			// Получаем полную информацию об объявлении
			listing, err := s.storage.GetListingByID(ctx, listingID)
			if err != nil {
				log.Printf("Ошибка при получении информации об объявлении %d для переиндексации: %v",
					listingID, err)
				return
			}

			// Переиндексируем объявление
			if err := s.storage.IndexListing(ctx, listing); err != nil {
				log.Printf("Ошибка при переиндексации объявления %d в OpenSearch: %v",
					listingID, err)
			} else {
				log.Printf("Объявление %d успешно переиндексировано в OpenSearch с %d изображениями",
					listingID, len(listing.Images))
			}
		}
	}()
}

// Обработка одного изображени
// Обработка одного изображения
func (s *StorefrontService) processOneImage(ctx context.Context, listingID int, imagePath string, isMain bool, zipReader *zip.Reader) error {
	var imgData []byte
	var err error
	var contentType string

	// 1. Проверим, если это URL
	if strings.HasPrefix(strings.ToLower(imagePath), "http://") ||
		strings.HasPrefix(strings.ToLower(imagePath), "https://") {
		// Загружаем изображение по URL
		imgData, contentType, err = s.downloadImage(imagePath)
		if err != nil {
			log.Printf("Ошибка при загрузке изображения %s: %v", imagePath, err)
			return err
		}
	} else if zipReader != nil {
		// 2. Ищем изображение в ZIP-архиве
		imgData, contentType, err = s.extractImageFromZip(zipReader, imagePath)
		if err != nil {
			log.Printf("Ошибка при извлечении изображения %s из архива: %v", imagePath, err)
			return err
		}
	} else {
		// 3. Пробуем интерпретировать как локальный путь (только для внутреннего использования)
		imgData, contentType, err = s.readLocalImage(imagePath)
		if err != nil {
			log.Printf("Ошибка при чтении локального изображения %s: %v", imagePath, err)
			return err
		}
	}

	if imgData == nil || len(imgData) == 0 {
		err := fmt.Errorf("не удалось получить данные для изображения %s", imagePath)
		log.Printf("%v", err)
		return err
	}

	// Обрабатываем изображение (уменьшаем и добавляем водяной знак)
	processedData, err := s.processImage(imgData, contentType)
	if err != nil {
		log.Printf("Ошибка при обработке изображения %s: %v", imagePath, err)
		return err
	}

	// Генерируем уникальное имя файла
	fileName := fmt.Sprintf("%d_%d%s", listingID, time.Now().UnixNano(), s.getExtensionFromContentType(contentType))
	filePath := filepath.Join("./uploads", fileName)

	// Сохраняем файл
	if err := ioutil.WriteFile(filePath, processedData, 0644); err != nil {
		log.Printf("Ошибка при сохранении изображения %s: %v", filePath, err)
		return err
	}

	// Создаем запись об изображении в базе данных
	image := &models.MarketplaceImage{
		ListingID:   listingID,
		FilePath:    fileName,
		FileName:    filepath.Base(imagePath),
		FileSize:    len(processedData),
		ContentType: contentType,
		IsMain:      isMain,
	}

	_, err = s.storage.AddListingImage(ctx, image)
	if err != nil {
		log.Printf("Ошибка при добавлении информации об изображении в БД: %v", err)
		return err
	}

	log.Printf("Изображение %s успешно обработано и добавлено для объявления %d", fileName, listingID)
	return nil
}

func (s *StorefrontService) processOneImageWithName(ctx context.Context, listingID int, imagePath string, isMain bool, zipReader *zip.Reader, fileName string) error {
	var imgData []byte
	var err error
	var contentType string

	if strings.HasPrefix(strings.ToLower(imagePath), "http://") ||
		strings.HasPrefix(strings.ToLower(imagePath), "https://") {
		imgData, contentType, err = s.downloadImage(imagePath)
		if err != nil {
			log.Printf("Ошибка при загрузке изображения %s: %v", imagePath, err)
			return err
		}
	} else if zipReader != nil {
		imgData, contentType, err = s.extractImageFromZip(zipReader, imagePath)
		if err != nil {
			log.Printf("Ошибка при извлечении изображения %s из архива: %v", imagePath, err)
			return err
		}
	} else {
		imgData, contentType, err = s.readLocalImage(imagePath)
		if err != nil {
			log.Printf("Ошибка при чтении локального изображения %s: %v", imagePath, err)
			return err
		}
	}

	if imgData == nil || len(imgData) == 0 {
		err := fmt.Errorf("не удалось получить данные для изображения %s", imagePath)
		log.Printf("%v", err)
		return err
	}

	// Обрабатываем изображение (уменьшаем и добавляем водяной знак)
	processedData, err := s.processImage(imgData, contentType)
	if err != nil {
		log.Printf("Ошибка при обработке изображения %s: %v", imagePath, err)
		return err
	}

	filePath := filepath.Join("./uploads", fileName)

	// Сохраняем файл
	if err := ioutil.WriteFile(filePath, processedData, 0644); err != nil {
		log.Printf("Ошибка при сохранении изображения %s: %v", filePath, err)
		return err
	}

	// Обновляем запись об изображении в базе данных (уже должна существовать)
	// Ищем запись с нужным ListingID и IsMain для обновления
	var imageID int
	err = s.storage.QueryRow(ctx,
		"SELECT id FROM marketplace_images WHERE listing_id = $1 AND is_main = $2 ORDER BY id LIMIT 1",
		listingID, isMain).Scan(&imageID)

	if err == nil {
		// Нашли запись, обновляем ее
		_, err = s.storage.Exec(ctx,
			"UPDATE marketplace_images SET file_size = $1, content_type = $2 WHERE id = $3",
			len(processedData), contentType, imageID)
		if err != nil {
			log.Printf("Ошибка при обновлении информации об изображении в БД: %v", err)
		}
	} else {
		// Запись не найдена, создаем новую
		image := &models.MarketplaceImage{
			ListingID:   listingID,
			FilePath:    fileName,
			FileName:    strings.Split(imagePath, "/")[len(strings.Split(imagePath, "/"))-1],
			FileSize:    len(processedData),
			ContentType: contentType,
			IsMain:      isMain,
		}

		_, err = s.storage.AddListingImage(ctx, image)
		if err != nil {
			log.Printf("Ошибка при добавлении информации об изображении в БД: %v", err)
			return err
		}
	}

	log.Printf("Изображение %s успешно обработано и сохранено для объявления %d", fileName, listingID)
	return nil
}

// downloadImage загружает изображение по URL
func (s *StorefrontService) downloadImage(url string) ([]byte, string, error) {
	client := &http.Client{
		Timeout: 120 * time.Second, // 2 минуты
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, "", fmt.Errorf("ошибка HTTP-запроса: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("некорректный статус ответа: %d", resp.StatusCode)
	}

	// Чтение данных
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, "", fmt.Errorf("ошибка чтения содержимого: %w", err)
	}

	// Проверяем тип содержимого
	contentType := resp.Header.Get("Content-Type")

	// Проверяем, является ли ответ HTML (иногда серверы возвращают HTML-страницы вместо изображений)
	if strings.Contains(contentType, "text/html") || strings.Contains(string(data[:min(1000, len(data))]), "<html") {
		return nil, "", fmt.Errorf("сервер вернул HTML вместо изображения для URL: %s", url)
	}

	// Если Content-Type не указывает на изображение, попробуем определить тип содержимого
	if !strings.HasPrefix(contentType, "image/") {
		detectedType := http.DetectContentType(data)
		log.Printf("Исходный тип содержимого: %s, определенный тип: %s", contentType, detectedType)

		if strings.HasPrefix(detectedType, "image/") {
			contentType = detectedType
		} else {
			// Пробуем по расширению файла
			if strings.HasSuffix(strings.ToLower(url), ".jpg") || strings.HasSuffix(strings.ToLower(url), ".jpeg") {
				contentType = "image/jpeg"
			} else if strings.HasSuffix(strings.ToLower(url), ".png") {
				contentType = "image/png"
			} else if strings.HasSuffix(strings.ToLower(url), ".gif") {
				contentType = "image/gif"
			} else if strings.HasSuffix(strings.ToLower(url), ".webp") {
				contentType = "image/webp"
			} else {
				// Если всё не удалось, пробуем по содержимому
				// Проверяем сигнатуры файлов
				if len(data) > 4 && data[0] == 0xFF && data[1] == 0xD8 && data[2] == 0xFF {
					// JPEG signature
					contentType = "image/jpeg"
				} else if len(data) > 8 && data[0] == 0x89 && data[1] == 0x50 && data[2] == 0x4E && data[3] == 0x47 {
					// PNG signature
					contentType = "image/png"
				} else if len(data) > 6 && data[0] == 0x47 && data[1] == 0x49 && data[2] == 0x46 {
					// GIF signature
					contentType = "image/gif"
				} else {
					contentType = "image/jpeg" // Предполагаем JPEG по умолчанию
				}
			}
		}
	}

	return data, contentType, nil
}
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// extractImageFromZip извлекает изображение из ZIP-архива
func (s *StorefrontService) extractImageFromZip(zipReader *zip.Reader, fileName string) ([]byte, string, error) {
	// Нормализуем имя файла для поиска (уберем лишние слеши и т.д.)
	normalizedFileName := filepath.Clean(fileName)
	normalizedFileName = filepath.ToSlash(normalizedFileName)
	normalizedFileName = strings.TrimPrefix(normalizedFileName, "/")

	// Ищем файл в архиве
	for _, zipFile := range zipReader.File {
		zipFileName := filepath.ToSlash(zipFile.Name)
		zipFileName = strings.TrimPrefix(zipFileName, "/")

		if zipFileName == normalizedFileName ||
			strings.HasSuffix(zipFileName, "/"+normalizedFileName) {
			// Нашли файл
			file, err := zipFile.Open()
			if err != nil {
				return nil, "", fmt.Errorf("ошибка открытия файла в архиве: %w", err)
			}
			defer file.Close()

			data, err := ioutil.ReadAll(file)
			if err != nil {
				return nil, "", fmt.Errorf("ошибка чтения файла из архива: %w", err)
			}

			// Определяем тип содержимого по расширению
			contentType := ""
			ext := strings.ToLower(filepath.Ext(zipFileName))
			switch ext {
			case ".jpg", ".jpeg":
				contentType = "image/jpeg"
			case ".png":
				contentType = "image/png"
			case ".gif":
				contentType = "image/gif"
			case ".webp":
				contentType = "image/webp"
			default:
				// Пробуем определить тип по содержимому
				contentType = http.DetectContentType(data)
			}

			if !strings.HasPrefix(contentType, "image/") {
				return nil, "", fmt.Errorf("файл не является изображением: %s", contentType)
			}

			return data, contentType, nil
		}
	}

	return nil, "", fmt.Errorf("файл не найден в архиве: %s", normalizedFileName)
}

// readLocalImage читает изображение из локальной файловой системы
func (s *StorefrontService) readLocalImage(path string) ([]byte, string, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, "", fmt.Errorf("ошибка чтения файла: %w", err)
	}

	// Определяем тип содержимого
	contentType := http.DetectContentType(data)
	if !strings.HasPrefix(contentType, "image/") {
		return nil, "", fmt.Errorf("файл не является изображением: %s", contentType)
	}

	return data, contentType, nil
}

// resizeImage изменяет размер изображения с помощью библиотеки imaging
func (s *StorefrontService) resizeImage(img image.Image, maxWidth, maxHeight int) image.Image {
	// Получаем текущие размеры
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// Если изображение уже меньше максимального размера, возвращаем как есть
	if width <= maxWidth && height <= maxHeight {
		return img
	}

	// Вычисляем новые размеры с сохранением пропорций
	var newWidth, newHeight int
	widthRatio := float64(maxWidth) / float64(width)
	heightRatio := float64(maxHeight) / float64(height)

	// Используем меньшее соотношение, чтобы уместиться в ограничения
	ratio := math.Min(widthRatio, heightRatio)
	newWidth = int(float64(width) * ratio)
	newHeight = int(float64(height) * ratio)

	// Используем библиотеку imaging для изменения размера изображения
	// с применением высококачественной интерполяции Lanczos
	resized := imaging.Resize(img, newWidth, newHeight, imaging.Lanczos)

	return resized
}
func (s *StorefrontService) addWatermark(img image.Image) (image.Image, error) {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// Создаем новый графический контекст с размерами изображения
	dc := gg.NewContextForImage(img)

	// Настройки водяного знака
	watermarkText := "SveTu.rs"

	// Определяем размер шрифта (пропорционально размеру изображения)
	fontSize := float64(height) * 0.03
	if fontSize < 12 {
		fontSize = 12 // минимальный размер шрифта
	}

	// Компактная оценка ширины текста
	textWidth := fontSize * float64(len(watermarkText)) * 0.55 // Уменьшено с 0.6 до 0.55
	textHeight := fontSize

	// Меньшие отступы для более компактного вида
	padding := fontSize * 0.2 // Уменьшено с 0.3 до 0.2
	bgWidth := textWidth + padding*2
	bgHeight := textHeight + padding*2

	// Положение водяного знака (правый нижний угол с отступом)
	bgX := float64(width) - bgWidth - fontSize*0.5
	bgY := float64(height) - bgHeight - fontSize*0.5

	// Защита от выхода за границы
	if bgX < 0 {
		bgX = 5
	}
	if bgY < 0 {
		bgY = 5
	}

	// Создаем черный фон с полупрозрачностью
	dc.SetRGBA(0, 0, 0, 0.6) // черный с 60% непрозрачностью
	dc.DrawRoundedRectangle(bgX, bgY, bgWidth, bgHeight, padding*0.5)
	dc.Fill()

	// Рисуем текст белым цветом
	dc.SetRGB(1, 1, 1) // белый цвет
	textX := bgX + padding
	textY := bgY + padding + textHeight*0.9 // коррекция для вертикального центрирования

	dc.DrawString(watermarkText, textX, textY)

	// Преобразуем контекст обратно в изображение
	result := dc.Image()

	log.Printf("Добавлен компактный векторный водяной знак '%s'", watermarkText)

	return result, nil
}

// Обновленная функция processImage
func (s *StorefrontService) processImage(imgData []byte, contentType string) ([]byte, error) {
	// Логирование начала процесса
	log.Printf("Начало обработки изображения. Размер входных данных: %d байт", len(imgData))

	// Декодируем изображение
	reader := bytes.NewReader(imgData)
	var img image.Image
	var err error
	var format string

	// Пробуем определить тип содержимого и декодировать соответственно
	// Сначала пробуем использовать image.DecodeConfig для определения формата
	_, format, err = image.DecodeConfig(bytes.NewReader(imgData))
	if err != nil {
		log.Printf("Не удалось определить формат изображения с помощью image.DecodeConfig: %v. Будем пробовать другие декодеры.", err)
		format = ""
	}

	// Если формат не удалось определить по контенту, используем contentType
	if format == "" && contentType != "" {
		switch {
		case strings.Contains(contentType, "jpeg") || strings.Contains(contentType, "jpg"):
			format = "jpeg"
		case strings.Contains(contentType, "png"):
			format = "png"
		case strings.Contains(contentType, "gif"):
			format = "gif"
		case strings.Contains(contentType, "webp"):
			format = "webp"
		}
	}

	// Сбрасываем reader для нового чтения
	reader.Seek(0, 0)

	// Пробуем различные форматы декодирования в зависимости от определенного формата
	if format == "jpeg" {
		log.Printf("Пробуем декодировать как JPEG")
		img, err = jpeg.Decode(reader)
	} else if format == "png" {
		log.Printf("Пробуем декодировать как PNG")
		img, err = png.Decode(reader)
	} else {
		// Если формат определить не удалось, пробуем все поддерживаемые форматы
		log.Printf("Пробуем универсальное декодирование")
		reader.Seek(0, 0)
		img, format, err = image.Decode(reader)
		if err != nil {
			// Если не удалось, пробуем принудительно как JPEG
			reader.Seek(0, 0)
			log.Printf("Пробуем принудительно декодировать как JPEG")
			img, err = jpeg.Decode(reader)
			if err != nil {
				// Если и это не удалось, пробуем PNG
				reader.Seek(0, 0)
				log.Printf("Пробуем принудительно декодировать как PNG")
				img, err = png.Decode(reader)
				if err != nil {
					// Если и это не удалось, пробуем GIF
					reader.Seek(0, 0)
					log.Printf("Пробуем принудительно декодировать как GIF")
					img, err = gif.Decode(reader)
					if err != nil {
						log.Printf("Все попытки декодирования изображения не удались: %v", err)
						return nil, fmt.Errorf("ошибка декодирования изображения: %w", err)
					} else {
						format = "gif"
					}
				} else {
					format = "png"
				}
			} else {
				format = "jpeg"
			}
		}
	}

	if err != nil {
		log.Printf("Ошибка декодирования изображения: %v", err)
		return nil, fmt.Errorf("ошибка декодирования изображения: %w", err)
	}

	// Логирование размеров изображения
	width := img.Bounds().Dx()
	height := img.Bounds().Dy()
	log.Printf("Изображение успешно декодировано как %s. Размеры: %dx%d", format, width, height)

	// Изменяем размер, если изображение слишком большое
	var resizedImg image.Image
	if width > 1920 || height > 1080 {
		log.Printf("Изменяем размер изображения до максимум 1920x1080")
		resizedImg = s.resizeImage(img, 1920, 1080)
	} else {
		log.Printf("Размер изображения не требует изменения")
		resizedImg = img
	}

	// Добавляем водяной знак
	log.Printf("Добавляем водяной знак")
	imgWithWatermark, err := s.addWatermark(resizedImg)
	if err != nil {
		log.Printf("Ошибка при добавлении водяного знака: %v. Продолжаем без него.", err)
		imgWithWatermark = resizedImg
	}

	// Кодируем изображение обратно в байты
	log.Printf("Кодируем изображение обратно в байты")
	var buf bytes.Buffer

	// Используем определенный формат для кодирования
	if format == "jpeg" || format == "" { // По умолчанию JPEG
		err = jpeg.Encode(&buf, imgWithWatermark, &jpeg.Options{Quality: 85})
		log.Printf("Кодирование в JPEG")
	} else if format == "png" {
		err = png.Encode(&buf, imgWithWatermark)
		log.Printf("Кодирование в PNG")
	} else if format == "gif" {
		err = gif.Encode(&buf, imgWithWatermark, &gif.Options{NumColors: 256})
		log.Printf("Кодирование в GIF")
	} else {
		// По умолчанию используем JPEG для других форматов
		err = jpeg.Encode(&buf, imgWithWatermark, &jpeg.Options{Quality: 85})
		log.Printf("Кодирование в JPEG (по умолчанию)")
	}

	if err != nil {
		log.Printf("Ошибка кодирования изображения: %v", err)
		return nil, fmt.Errorf("ошибка кодирования изображения: %w", err)
	}

	// Завершение
	log.Printf("Обработка изображения завершена успешно. Размер данных после обработки: %d байт", buf.Len())
	return buf.Bytes(), nil
}

// getExtensionFromContentType возвращает расширение файла на основе типа содержимого
func (s *StorefrontService) getExtensionFromContentType(contentType string) string {
	contentType = strings.ToLower(contentType)
	switch {
	case strings.Contains(contentType, "jpeg") || strings.Contains(contentType, "jpg"):
		return ".jpg"
	case strings.Contains(contentType, "png"):
		return ".png"
	case strings.Contains(contentType, "gif"):
		return ".gif"
	case strings.Contains(contentType, "webp"):
		return ".webp"
	case strings.Contains(contentType, "bmp"):
		return ".bmp"
	default:
		return ".jpg" // По умолчанию
	}
}
