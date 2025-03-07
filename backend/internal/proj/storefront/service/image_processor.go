// backend/internal/proj/storefront/service/image_processor.go
package service

import (
	"archive/zip"
	"backend/internal/domain/models"
	"bytes"
	"context"
 	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"image/draw"

 	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
	
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


// Обновленная функция processImage
// Обновленная функция processImage
func (s *StorefrontService) processImage(imgData []byte, contentType string) ([]byte, error) {
	// Декодируем изображение
	reader := bytes.NewReader(imgData)
	var img image.Image
	var err error
	
	if strings.Contains(contentType, "jpeg") || strings.Contains(contentType, "jpg") {
		img, err = jpeg.Decode(reader)
	} else if strings.Contains(contentType, "png") {
		img, err = png.Decode(reader)
	} else {
		img, _, err = image.Decode(reader)
	}
	
	if err != nil {
		return nil, fmt.Errorf("ошибка декодирования изображения: %w", err)
	}
	
	// Изменяем размер, если изображение слишком большое
	if img.Bounds().Dx() > 1920 || img.Bounds().Dy() > 1080 {
		img = imaging.Resize(img, 1920, 1080, imaging.Lanczos)
	}
	
	// Добавляем водяной знак
	imgWithWatermark, err := s.addWatermark(img)
	if err != nil {
		log.Printf("Ошибка при добавлении водяного знака: %v, пропускаем", err)
		// В случае ошибки используем изображение без водяного знака
		imgWithWatermark = img
	}
	
	// Кодируем изображение обратно в байты
	var buf bytes.Buffer
	
	if strings.Contains(contentType, "jpeg") || strings.Contains(contentType, "jpg") {
		err = jpeg.Encode(&buf, imgWithWatermark, &jpeg.Options{Quality: 85})
	} else if strings.Contains(contentType, "png") {
		err = png.Encode(&buf, imgWithWatermark)
	} else {
		err = jpeg.Encode(&buf, imgWithWatermark, &jpeg.Options{Quality: 85})
	}
	
	if err != nil {
		return nil, fmt.Errorf("ошибка кодирования изображения: %w", err)
	}
	
	// Выводим лог для подтверждения
	log.Printf("Водяной знак успешно добавлен")
	
	return buf.Bytes(), nil
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
// addWatermark добавляет водяной знак "SveTu.rs" к изображению
func (s *StorefrontService) addWatermark(img image.Image) (image.Image, error) {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	
	// Создаем новое RGBA изображение на основе исходного
	dst := imaging.Clone(img)
	
	// Настройки водяного знака
	watermarkText := "SveTu.rs"
	
	// Уменьшаем размер точки для более компактного размещения
	dotSize := int(math.Max(2, float64(height)*0.008))
	
	// Массив с представлением символов (1 = белая точка, 0 = прозрачность)
	// Каждый символ имеет размер 5x5 точек для компактности
	letters := map[rune][][]int{
		'S': {
			{1, 1, 1, 1, 1},
			{1, 0, 0, 0, 0},
			{1, 1, 1, 1, 1},
			{0, 0, 0, 0, 1},
			{1, 1, 1, 1, 1},
		},
		'v': {
			{1, 0, 0, 0, 1},
			{1, 0, 0, 0, 1},
			{0, 1, 0, 1, 0},
			{0, 1, 0, 1, 0},
			{0, 0, 1, 0, 0},
		},
		'e': {
			{1, 1, 1, 1, 1},
			{1, 0, 0, 0, 0},
			{1, 1, 1, 1, 1},
			{1, 0, 0, 0, 0},
			{1, 1, 1, 1, 1},
		},
		'T': {
			{1, 1, 1, 1, 1},
			{0, 0, 1, 0, 0},
			{0, 0, 1, 0, 0},
			{0, 0, 1, 0, 0},
			{0, 0, 1, 0, 0},
		},
		'u': {
			{1, 0, 0, 0, 1},
			{1, 0, 0, 0, 1},
			{1, 0, 0, 0, 1},
			{1, 0, 0, 0, 1},
			{1, 1, 1, 1, 1},
		},
		'.': {
			{0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0},
			{0, 0, 1, 0, 0},
		},
		'r': {
			{1, 1, 1, 1, 0},
			{1, 0, 0, 0, 1},
			{1, 1, 1, 1, 0},
			{1, 0, 1, 0, 0},
			{1, 0, 0, 1, 0},
		},
		's': {
			{0, 1, 1, 1, 0},
			{1, 0, 0, 0, 0},
			{0, 1, 1, 1, 0},
			{0, 0, 0, 0, 1},
			{0, 1, 1, 1, 0},
		},
	}
	
	// Рассчитываем общую ширину водяного знака
	letterWidth := 5 * dotSize
	letterSpacing := dotSize
	totalWidth := (letterWidth + letterSpacing) * len(watermarkText) - letterSpacing
	totalHeight := 5 * dotSize
	
	// Добавляем отступы
	padding := dotSize * 2
	bgWidth := totalWidth + padding*2
	bgHeight := totalHeight + padding*2
	
	// Положение водяного знака (правый нижний угол с отступом)
	bgX := width - bgWidth - 10
	bgY := height - bgHeight - 10
	
	// Защита от выхода за границы
	if bgX < 0 {
		bgX = 5
	}
	if bgY < 0 {
		bgY = 5
	}
	
	// Создаем черный фон с полупрозрачностью для водяного знака
	bgRect := image.Rect(bgX, bgY, bgX+bgWidth, bgY+bgHeight)
	draw.Draw(dst, bgRect, image.NewUniform(color.RGBA{0, 0, 0, 180}), image.Point{}, draw.Over)
	
	// Рисуем каждую букву
	xOffset := bgX + padding
	yOffset := bgY + padding
	
	for _, char := range watermarkText {
		charPattern, exists := letters[char]
		if !exists {
			// Если символа нет в нашей коллекции, пропускаем
			xOffset += dotSize * 3
			continue
		}
		
		// Рисуем символ
		for y := 0; y < len(charPattern); y++ {
			for x := 0; x < len(charPattern[y]); x++ {
				if charPattern[y][x] == 1 {
					// Рисуем белую точку
					pixelX := xOffset + x*dotSize
					pixelY := yOffset + y*dotSize
					
					// Рисуем квадратную точку размером dotSize
					for dx := 0; dx < dotSize; dx++ {
						for dy := 0; dy < dotSize; dy++ {
							if pixelX+dx < width && pixelY+dy < height {
								dst.Set(pixelX+dx, pixelY+dy, color.RGBA{255, 255, 255, 255})
							}
						}
					}
				}
			}
		}
		
		// Переходим к следующему символу
		xOffset += letterWidth + letterSpacing
	}
	
	// Выводим лог для подтверждения
	log.Printf("Добавлен полный водяной знак '%s', размер: %dx%d", watermarkText, bgWidth, bgHeight)
	
	return dst, nil
}


// drawBoldText рисует жирный текст на изображении с использованием простых линий
func drawBoldText(img *image.NRGBA, text string, x, y, thickness int, col color.RGBA) {
	// Расстояние между буквами
	letterSpacing := thickness * 4
	
	// Координаты для рисования
	currentX := x
	
	// Рисуем каждую букву
	for _, char := range text {
		// Увеличиваем X после каждой буквы
		currentX = drawBoldLetter(img, char, currentX, y, thickness, col)
		currentX += letterSpacing
	}
}

// drawBoldLetter рисует одну букву как набор жирных линий
func drawBoldLetter(img *image.NRGBA, char rune, x, y, thickness int, col color.RGBA) int {
	// Базовая ширина символа
	charWidth := thickness * 5
	
	// Рисуем линии с заданной толщиной
	drawLine := func(x1, y1, x2, y2 int) {
		// Рисуем жирную линию
		dx := x2 - x1
		dy := y2 - y1
		length := int(math.Sqrt(float64(dx*dx + dy*dy)))
		if length == 0 {
			return
		}
		
		// Нормализуем вектор направления
		nx := float64(dx) / float64(length)
		ny := float64(dy) / float64(length)
		
		// Рисуем основную линию
		for i := 0; i <= length; i++ {
			px := int(float64(x1) + float64(i)*nx)
			py := int(float64(y1) + float64(i)*ny)
			
			// Рисуем точку с заданной толщиной
			for tx := -thickness; tx <= thickness; tx++ {
				for ty := -thickness; ty <= thickness; ty++ {
					if tx*tx+ty*ty <= thickness*thickness {
						img.Set(px+tx, py+ty, col)
					}
				}
			}
		}
	}
	
	// Высота символа
	height := thickness * 10
	
	// Рисуем символы
	switch char {
	case 'S':
		// Горизонтальные линии
		drawLine(x, y-height/3, x+charWidth, y-height/3)
		drawLine(x, y, x+charWidth, y)
		drawLine(x, y+height/3, x+charWidth, y+height/3)
		// Вертикальные линии
		drawLine(x, y-height/3, x, y)
		drawLine(x+charWidth, y, x+charWidth, y+height/3)
		return x + charWidth
		
	case 'v':
		// Диагональные линии
		drawLine(x, y-height/3, x+charWidth/2, y+height/3)
		drawLine(x+charWidth/2, y+height/3, x+charWidth, y-height/3)
		return x + charWidth
		
	case 'e':
		// Вертикальная линия
		drawLine(x, y-height/3, x, y+height/3)
		// Горизонтальные линии
		drawLine(x, y-height/3, x+charWidth, y-height/3)
		drawLine(x, y, x+charWidth, y)
		drawLine(x, y+height/3, x+charWidth, y+height/3)
		return x + charWidth
		
	case 'T':
		// Горизонтальная линия
		drawLine(x, y-height/3, x+charWidth, y-height/3)
		// Вертикальная линия
		drawLine(x+charWidth/2, y-height/3, x+charWidth/2, y+height/3)
		return x + charWidth
		
	case 'u':
		// Вертикальные линии
		drawLine(x, y-height/3, x, y+height/3)
		drawLine(x+charWidth, y-height/3, x+charWidth, y+height/3)
		// Горизонтальная линия внизу
		drawLine(x, y+height/3, x+charWidth, y+height/3)
		return x + charWidth
		
	case '.':
		// Точка
		for tx := -thickness; tx <= thickness; tx++ {
			for ty := -thickness; ty <= thickness; ty++ {
				if tx*tx+ty*ty <= thickness*thickness {
					img.Set(x+tx, y+height/3+ty, col)
				}
			}
		}
		return x + thickness*2
		
	case 'r':
		// Вертикальная линия
		drawLine(x, y-height/3, x, y+height/3)
		// Верхняя горизонтальная линия
		drawLine(x, y-height/3, x+charWidth, y-height/3)
		// Диагональная линия
		drawLine(x, y, x+charWidth, y+height/3)
		return x + charWidth
		
	case 's':
		// Меньшая версия S
		drawLine(x, y-height/4, x+charWidth*3/4, y-height/4)
		drawLine(x, y, x+charWidth*3/4, y)
		drawLine(x, y+height/4, x+charWidth*3/4, y+height/4)
		drawLine(x, y-height/4, x, y)
		drawLine(x+charWidth*3/4, y, x+charWidth*3/4, y+height/4)
		return x + charWidth*3/4
		
	default:
		// Любой другой символ - просто вертикальная линия
		drawLine(x, y-height/3, x, y+height/3)
		return x + thickness*2
	}
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