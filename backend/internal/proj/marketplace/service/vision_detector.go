// backend/internal/proj/marketplace/service/vision_detector.go
package service

import (
	"context"
	"log"
	"os"

	vision "cloud.google.com/go/vision/apiv1"
)

func DetectFaceInImage(ctx context.Context, imagePath string) (bool, error) {
	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		log.Printf("DetectFaceInImage: ошибка создания клиента Vision API: %v", err)
		return false, err
	}
	defer func() {
		if err := client.Close(); err != nil {
			log.Printf("DetectFaceInImage: ошибка закрытия клиента Vision API: %v", err)
		}
	}()

	file, err := os.Open(imagePath)
	if err != nil {
		log.Printf("DetectFaceInImage: ошибка открытия файла: %v", err)
		return false, err
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("Failed to close file: %v", err)
		}
	}()

	image, err := vision.NewImageFromReader(file)
	if err != nil {
		log.Printf("DetectFaceInImage: ошибка создания изображения: %v", err)
		return false, err
	}

	faces, err := client.DetectFaces(ctx, image, nil, 5)
	if err != nil {
		log.Printf("DetectFaceInImage: ошибка обнаружения лиц: %v", err)
		return false, err
	}

	if len(faces) > 0 {
		log.Printf("DetectFaceInImage: обнаружено %d лиц", len(faces))
		return true, nil
	}

	log.Printf("DetectFaceInImage: лиц не обнаружено")
	return false, nil
}
