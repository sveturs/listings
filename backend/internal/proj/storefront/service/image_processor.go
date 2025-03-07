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

 	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
	
 	//"github.com/disintegration/imaging"
 	//"github.com/golang/freetype/truetype"
	"github.com/fogleman/gg"
    "github.com/disintegration/imaging"


)

// ProcessImportImages обрабатывает изображения для импортируемого объявления
func (s *StorefrontService) ProcessImportImages(
	ctx context.Context, 
	listingID int, 
	imagesStr string,
	zipReader *zip.Reader, // Опционально: если есть zip-архив с изображениями
) error {
	// Предположим, что изображения разделены запятыми
	imagesList := strings.Split(imagesStr, ",")
	
	// Директория для сохранения обработанных изображений
	uploadDir := "./uploads"
	
	// Убедимся, что директория существует
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return fmt.Errorf("error creating upload directory: %w", err)
	}
	
	// Обработка каждого изображения
	for i, imagePath := range imagesList {
		imagePath = strings.TrimSpace(imagePath)
		if imagePath == "" {
			continue
		}
		
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
				continue
			}
		} else if zipReader != nil {
			// 2. Ищем изображение в ZIP-архиве
			imgData, contentType, err = s.extractImageFromZip(zipReader, imagePath)
			if err != nil {
				log.Printf("Ошибка при извлечении изображения %s из архива: %v", imagePath, err)
				continue
			}
		} else {
			// 3. Пробуем интерпретировать как локальный путь (только для внутреннего использования)
			imgData, contentType, err = s.readLocalImage(imagePath)
			if err != nil {
				log.Printf("Ошибка при чтении локального изображения %s: %v", imagePath, err)
				continue
			}
		}
		
		if imgData == nil || len(imgData) == 0 {
			log.Printf("Не удалось получить данные для изображения %s", imagePath)
			continue
		}
		
		// Обрабатываем изображение (уменьшаем и добавляем водяной знак)
		processedData, err := s.processImage(imgData, contentType)
		if err != nil {
			log.Printf("Ошибка при обработке изображения %s: %v", imagePath, err)
			continue
		}
		
		// Генерируем уникальное имя файла
		fileName := fmt.Sprintf("%d_%d%s", listingID, time.Now().UnixNano(), s.getExtensionFromContentType(contentType))
		filePath := filepath.Join(uploadDir, fileName)
		
		// Сохраняем файл
		if err := ioutil.WriteFile(filePath, processedData, 0644); err != nil {
			log.Printf("Ошибка при сохранении изображения %s: %v", filePath, err)
			continue
		}
		
		// Создаем запись об изображении в базе данных
		image := &models.MarketplaceImage{
			ListingID:   listingID,
			FilePath:    fileName,
			FileName:    filepath.Base(imagePath),
			FileSize:    len(processedData),
			ContentType: contentType,
			IsMain:      i == 0, // Первое изображение - основное
		}
		
		_, err = s.storage.AddListingImage(ctx, image)
		if err != nil {
			log.Printf("Ошибка при добавлении информации об изображении в БД: %v", err)
			continue
		}
		
		log.Printf("Изображение %s успешно обработано и добавлено для объявления %d", fileName, listingID)
	}
	
	return nil
}

// downloadImage загружает изображение по URL
func (s *StorefrontService) downloadImage(url string) ([]byte, string, error) {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	
	resp, err := client.Get(url)
	if err != nil {
		return nil, "", fmt.Errorf("ошибка HTTP-запроса: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("некорректный статус ответа: %d", resp.StatusCode)
	}
	
	contentType := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		return nil, "", fmt.Errorf("некорректный тип содержимого: %s", contentType)
	}
	
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, "", fmt.Errorf("ошибка чтения содержимого: %w", err)
	}
	
	return data, contentType, nil
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
    textWidth := fontSize * float64(len(watermarkText)) * 0.55  // Уменьшено с 0.6 до 0.55
    textHeight := fontSize
    
    // Меньшие отступы для более компактного вида
    padding := fontSize * 0.2  // Уменьшено с 0.3 до 0.2
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
    
    if strings.Contains(contentType, "jpeg") || strings.Contains(contentType, "jpg") {
        img, err = jpeg.Decode(reader)
        log.Printf("Декодирование JPEG изображения")
    } else if strings.Contains(contentType, "png") {
        img, err = png.Decode(reader)
        log.Printf("Декодирование PNG изображения")
    } else {
        img, _, err = image.Decode(reader)
        log.Printf("Декодирование изображения неизвестного формата")
    }
    
    if err != nil {
        log.Printf("Ошибка декодирования изображения: %v", err)
        return nil, fmt.Errorf("ошибка декодирования изображения: %w", err)
    }
    
    // Логирование размеров изображения
    width := img.Bounds().Dx()
    height := img.Bounds().Dy()
    log.Printf("Изображение успешно декодировано. Размеры: %dx%d", width, height)
    
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
    
    if strings.Contains(contentType, "jpeg") || strings.Contains(contentType, "jpg") {
        err = jpeg.Encode(&buf, imgWithWatermark, &jpeg.Options{Quality: 85})
        log.Printf("Кодирование в JPEG")
    } else if strings.Contains(contentType, "png") {
        err = png.Encode(&buf, imgWithWatermark)
        log.Printf("Кодирование в PNG")
    } else {
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
	switch contentType {
	case "image/jpeg", "image/jpg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/gif":
		return ".gif"
	case "image/webp":
		return ".webp"
	default:
		return ".jpg" // По умолчанию
	}
}