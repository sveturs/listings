// backend/internal/proj/marketplace/service/cloudinary.go
package service

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"

	//	"strings"
	vision "cloud.google.com/go/vision/apiv1"
	visionpb "cloud.google.com/go/vision/v2/apiv1/visionpb"
)

type CloudinaryService struct {
	cld *cloudinary.Cloudinary
}

func NewCloudinaryService() (*CloudinaryService, error) {
	cloudinaryURL := os.Getenv("CLOUDINARY_URL")
	log.Printf("NewCloudinaryService: используем CLOUDINARY_URL: [скрыто для безопасности]")

	if cloudinaryURL == "" {
		log.Printf("NewCloudinaryService: CLOUDINARY_URL не установлен в переменных окружения")
		return nil, fmt.Errorf("CLOUDINARY_URL not set")
	}

	log.Printf("NewCloudinaryService: создание клиента Cloudinary")
	cld, err := cloudinary.NewFromURL(cloudinaryURL)
	if err != nil {
		log.Printf("NewCloudinaryService: ошибка создания клиента Cloudinary: %v", err)
		return nil, err
	}

	log.Printf("NewCloudinaryService: клиент Cloudinary успешно создан")
	return &CloudinaryService{cld: cld}, nil
}

func (s *CloudinaryService) MakeUglyPhotoBeautiful(ctx context.Context, imagePath string) (map[string]interface{}, error) {
	log.Printf("MakeUglyPhotoBeautiful: начинаем обработку изображения: %s", imagePath)

	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		log.Printf("MakeUglyPhotoBeautiful: файл не существует: %s", imagePath)
		return nil, fmt.Errorf("file does not exist: %s", imagePath)
	}

	transformation := "e_background_removal,e_improve,w_1000,h_1000,c_fill,b_white,e_brightness:30,e_contrast:30,e_sharpen,q_auto"
	log.Printf("MakeUglyPhotoBeautiful: применяем трансформации: %s", transformation)

	uploadParams := uploader.UploadParams{
		Transformation: transformation,
		Tags:           []string{"ugly_to_beautiful", "auto_enhanced"},
	}

	resp, err := s.cld.Upload.Upload(ctx, imagePath, uploadParams)
	if err != nil {
		log.Printf("MakeUglyPhotoBeautiful: ошибка загрузки в Cloudinary: %v", err)
		return nil, err
	}

	log.Printf("MakeUglyPhotoBeautiful: успешно, URL: %s", resp.SecureURL)
	return map[string]interface{}{
		"url":       resp.SecureURL,
		"public_id": resp.PublicID,
		"width":     resp.Width,
		"height":    resp.Height,
	}, nil
}

func (s *CloudinaryService) ModerateImage(ctx context.Context, imagePath string) (map[string]interface{}, error) {
	log.Printf("ModerateImage: старт модерации файла %s", imagePath)

	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		log.Printf("ModerateImage: файл %s не существует", imagePath)
		return nil, fmt.Errorf("file does not exist: %s", imagePath)
	}

	result := map[string]interface{}{
		"safe":      true,
		"issues":    []string{},
		"url":       "",
		"public_id": "",
		"reason":    "",
	}

	// 1. Проверка Google Vision SafeSearch
	log.Printf("ModerateImage: запускаем Google Vision SafeSearch анализ")
	visionClient, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		log.Printf("ModerateImage: ошибка создания Vision клиента: %v", err)
		return result, nil
	}
	defer func() {
		if err := visionClient.Close(); err != nil {
			log.Printf("ModerateImage: failed to close vision client: %v", err)
		}
	}()

	file, err := os.Open(imagePath)
	if err != nil {
		log.Printf("ModerateImage: ошибка открытия файла Vision: %v", err)
		return result, nil
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("ModerateImage: failed to close file: %v", err)
		}
	}()

	img, err := vision.NewImageFromReader(file)
	if err != nil {
		log.Printf("ModerateImage: ошибка создания изображения Vision: %v", err)
		return result, nil
	}

	safe, err := visionClient.DetectSafeSearch(ctx, img, nil)
	if err != nil {
		log.Printf("ModerateImage: ошибка SafeSearch анализа: %v", err)
		return result, nil
	}

	// для модерации
	type unsafeCategory struct {
		Label      string
		Readable   string
		Likelihood visionpb.Likelihood
		Threshold  visionpb.Likelihood
	}

	categories := []unsafeCategory{
		{"adult", "эротический контент", safe.Adult, visionpb.Likelihood_LIKELY},
		{"violence", "насильственный контент", safe.Violence, visionpb.Likelihood_POSSIBLE},
		{"racy", "провокационный контент", safe.Racy, visionpb.Likelihood_VERY_LIKELY},
		{"medical", "медицинский (шокирующий) контент", safe.Medical, visionpb.Likelihood_VERY_LIKELY},
		{"spoof", "фальшивый или пародийный контент", safe.Spoof, visionpb.Likelihood_LIKELY},
	}

	for _, cat := range categories {
		if cat.Likelihood >= cat.Threshold {
			result["safe"] = false
			issueText := fmt.Sprintf("На изображении вероятно обнаружен %s (%s).",
				cat.Readable, translateLikelihood(cat.Likelihood))
			result["reason"] = issueText
			result["issues"] = append(result["issues"].([]string), issueText)
			log.Printf("ModerateImage: отклонено Google Vision: %s => %s (порог: %s)",
				cat.Label, cat.Likelihood.String(), cat.Threshold.String())
			log.Printf("ModerateImage: итоговый результат: %+v", result)
			return result, nil
		}
	}

	log.Printf("ModerateImage: Google Vision проверка пройдена успешно, продолжаем Cloudinary модерацию")

	// 2. Загрузка в Cloudinary и проверка через webpurify
	uploadParams := uploader.UploadParams{
		Moderation: "webpurify",
	}

	resp, err := s.cld.Upload.Upload(ctx, imagePath, uploadParams)
	if err != nil {
		log.Printf("ModerateImage: ошибка загрузки в Cloudinary: %v", err)
		return nil, err
	}

	log.Printf("ModerateImage: файл успешно загружен в Cloudinary, публичный ID: %s", resp.PublicID)

	result["url"] = resp.SecureURL
	result["public_id"] = resp.PublicID

	if len(resp.Moderation) > 0 {
		for _, mod := range resp.Moderation {
			if mod.Status == "rejected" {
				result["safe"] = false
				msg := fmt.Sprintf("Отклонено Cloudinary: %s", mod.Kind)
				result["issues"] = append(result["issues"].([]string), msg)
				result["reason"] = msg
				log.Printf("ModerateImage: изображение отклонено Cloudinary: %s", msg)
			}
		}
	} else {
		log.Printf("ModerateImage: Cloudinary модерация успешно пройдена")
	}

	log.Printf("ModerateImage: итоговый результат: %+v", result)
	return result, nil
}

func translateLikelihood(likelihood visionpb.Likelihood) string {
	switch likelihood {
	case visionpb.Likelihood_VERY_LIKELY:
		return "очень вероятно"
	case visionpb.Likelihood_LIKELY:
		return "вероятно"
	case visionpb.Likelihood_POSSIBLE:
		return "возможно"
	default:
		return "неопределено"
	}
}

// EnhanceImage улучшает изображение для товарной карточки
func (s *CloudinaryService) EnhanceImage(ctx context.Context, imagePath string) (map[string]interface{}, error) {
	log.Printf("EnhanceImage: начинаем улучшение товарного фото: %s", imagePath)

	// Проверяем существование файла
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		log.Printf("EnhanceImage: файл %s не существует", imagePath)
		return nil, fmt.Errorf("file does not exist: %s", imagePath)
	}

	// Усиленная трансформация
	transformationStr := "e_improve:outdoor:70,e_brightness:25,e_contrast:30,c_pad,w_1000,h_1000,b_white,q_auto:best"
	log.Printf("EnhanceImage: используем следующую трансформацию: %s", transformationStr)

	// Создаем параметры загрузки
	uploadParams := uploader.UploadParams{
		Transformation: transformationStr,
		Tags:           []string{"enhanced", "product"},
	}

	// Пробуем включить удаление фона, если оно доступно
	// Оставим закомментированным до проверки доступности
	// uploadParams.BackgroundRemoval = "cloudinary_ai"

	resp, err := s.cld.Upload.Upload(ctx, imagePath, uploadParams)
	if err != nil {
		log.Printf("EnhanceImage: ошибка улучшения изображения: %v", err)
		return nil, err
	}

	// Проверяем и логируем ответ
	if resp.SecureURL == "" {
		log.Printf("EnhanceImage: предупреждение - пустой URL в ответе Cloudinary")
	} else {
		log.Printf("EnhanceImage: изображение успешно улучшено, URL: %s", resp.SecureURL)
	}

	return map[string]interface{}{
		"url":       resp.SecureURL,
		"public_id": resp.PublicID,
		"width":     resp.Width,
		"height":    resp.Height,
	}, nil
}

// TestAvailableTransformations проверяет доступные трансформации
func (s *CloudinaryService) TestAvailableTransformations(ctx context.Context, imagePath string) {
	// Тест базовой трансформации
	log.Printf("Тестирование базовых трансформаций...")
	_, err1 := s.cld.Upload.Upload(ctx, imagePath, uploader.UploadParams{
		Transformation: "e_background_removal,e_improve,w_1000,h_1000,c_fill,b_white,e_brightness:30,e_contrast:30,e_sharpen,q_auto",
		Tags:           []string{"test_basic"},
	})
	log.Printf("Результат базовых трансформаций: %v", err1 == nil)

	// Тест удаления фона методом 1
	log.Printf("Тестирование удаления фона (метод 1)...")
	_, err2 := s.cld.Upload.Upload(ctx, imagePath, uploader.UploadParams{
		Transformation: "e_background_removal",
		Tags:           []string{"test_bgremoval_1"},
	})
	log.Printf("Результат удаления фона (метод 1): %v", err2 == nil)

	// Тест удаления фона методом 2
	log.Printf("Тестирование удаления фона (метод 2)...")
	_, err3 := s.cld.Upload.Upload(ctx, imagePath, uploader.UploadParams{
		BackgroundRemoval: "cloudinary_ai",
		Tags:              []string{"test_bgremoval_2"},
	})
	log.Printf("Результат удаления фона (метод 2): %v", err3 == nil)
}

func (s *CloudinaryService) TestBackgroundRemoval(ctx context.Context, imagePath string) (bool, error) {
	_, err := s.cld.Upload.Upload(ctx, imagePath, uploader.UploadParams{
		Transformation: "e_background_removal",
		Tags:           []string{"test", "background_removal"},
	})
	if err != nil {
		log.Printf("Тест удаления фона: ошибка: %v", err)
		return false, err
	}

	log.Printf("Тест удаления фона: успешно")
	return true, nil
}

// EnhancePreview создаёт предпросмотр улучшенного изображения
func (s *CloudinaryService) EnhancePreview(ctx context.Context, imagePath string) (map[string]interface{}, error) {
	log.Printf("EnhancePreview: создание предпросмотра для %s", imagePath)

	// Загружаем оригинал
	origResp, err := s.cld.Upload.Upload(ctx, imagePath, uploader.UploadParams{
		Tags: []string{"original", "preview"},
	})
	if err != nil {
		log.Printf("EnhancePreview: ошибка загрузки оригинала: %v", err)
		return nil, err
	}

	// Трансформация аналогичная EnhanceImage
	transformationStr := "e_improve:outdoor,e_brightness:15,e_contrast:20,e_vibrance:10,e_sharpen:15,b_gen_remove,c_pad,w_1000,h_1000,b_white,q_auto:best"

	// Создаём улучшенную версию
	enhResp, err := s.cld.Upload.Upload(ctx, imagePath, uploader.UploadParams{
		Transformation: transformationStr,
		Tags:           []string{"enhanced", "preview", "ai_processed"},
	})
	if err != nil {
		log.Printf("EnhancePreview: ошибка создания улучшенного предпросмотра: %v", err)
		return nil, err
	}

	log.Printf("EnhancePreview: успешно создан предпросмотр, оригинал: %s, улучшенный: %s",
		origResp.SecureURL, enhResp.SecureURL)

	// Возвращаем ссылки на обе версии для сравнения
	return map[string]interface{}{
		"original":    origResp.SecureURL,
		"enhanced":    enhResp.SecureURL,
		"original_id": origResp.PublicID,
		"enhanced_id": enhResp.PublicID,
	}, nil
}

// RemoveBackground удаляет фон с изображения
func (s *CloudinaryService) RemoveBackground(ctx context.Context, imagePath string) (map[string]interface{}, error) {
	resp, err := s.cld.Upload.Upload(ctx, imagePath, uploader.UploadParams{
		BackgroundRemoval: "cloudinary_ai", // AI-удаление фона
		Tags:              []string{"nobg", "product"},
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"url":       resp.SecureURL,
		"public_id": resp.PublicID,
	}, nil
}
