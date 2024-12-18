// backend/internal/storage/postgres/car_image.go
package postgres

import (
    "backend/internal/domain/models"
	
    "context"
)

func (db *Storage) AddCarImage(ctx context.Context, image *models.CarImage) (int, error) {
    var imageID int
    err := db.pool.QueryRow(ctx, `
        INSERT INTO car_images (car_id, file_path, file_name, file_size, content_type, is_main)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id
    `,
        image.CarID, image.FilePath, image.FileName,
        image.FileSize, image.ContentType, image.IsMain,
    ).Scan(&imageID)

    return imageID, err
}

func (db *Storage) GetCarImages(ctx context.Context, carID string) ([]models.CarImage, error) {
    rows, err := db.pool.Query(ctx, `
        SELECT id, car_id, file_path, file_name, file_size, content_type, is_main, created_at
        FROM car_images
        WHERE car_id = $1
        ORDER BY is_main DESC, created_at DESC
    `, carID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var images []models.CarImage
    for rows.Next() {
        var img models.CarImage
        err := rows.Scan(
            &img.ID, &img.CarID, &img.FilePath,
            &img.FileName, &img.FileSize, &img.ContentType,
            &img.IsMain, &img.CreatedAt,
        )
        if err != nil {
            continue
        }
        images = append(images, img)
    }

    return images, rows.Err()
}

func (db *Storage) DeleteCarImage(ctx context.Context, imageID string) (string, error) {
    var filePath string
    err := db.pool.QueryRow(ctx,
        "DELETE FROM car_images WHERE id = $1 RETURNING file_path", imageID,
    ).Scan(&filePath)
    return filePath, err
}