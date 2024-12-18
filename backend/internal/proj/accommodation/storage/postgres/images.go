// backend/internal/proj/accommodation/storage/postgres/images.go
package postgres

import (
    "backend/internal/domain/models"
    "context"
)

func (s *Storage) AddRoomImage(ctx context.Context, image *models.RoomImage) (int, error) {
    var imageID int
    err := s.pool.QueryRow(ctx, `
        INSERT INTO room_images (room_id, file_path, file_name, file_size, content_type, is_main)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id
    `,
        image.RoomID, image.FilePath, image.FileName,
        image.FileSize, image.ContentType, image.IsMain,
    ).Scan(&imageID)

    return imageID, err
}

func (s *Storage) GetRoomImages(ctx context.Context, roomID string) ([]models.RoomImage, error) {
    rows, err := s.pool.Query(ctx, `
        SELECT id, room_id, file_path, file_name, file_size, content_type, is_main, created_at
        FROM room_images
        WHERE room_id = $1
        ORDER BY is_main DESC, created_at DESC
    `, roomID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var images []models.RoomImage
    for rows.Next() {
        var img models.RoomImage
        err := rows.Scan(
            &img.ID, &img.RoomID, &img.FilePath,
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

func (s *Storage) DeleteRoomImage(ctx context.Context, imageID string) (string, error) {
    var filePath string
    err := s.pool.QueryRow(ctx,
        "DELETE FROM room_images WHERE id = $1 RETURNING file_path",
        imageID,
    ).Scan(&filePath)
    return filePath, err
}

func (s *Storage) GetBedImages(ctx context.Context, bedID string) ([]models.RoomImage, error) {
    rows, err := s.pool.Query(ctx, `
        SELECT bi.id, b.room_id, bi.file_path, bi.file_name, bi.file_size, bi.content_type, bi.is_main, bi.created_at
        FROM bed_images bi
        JOIN beds b ON bi.bed_id = b.id
        WHERE bi.bed_id = $1
        ORDER BY bi.is_main DESC, bi.created_at DESC
    `, bedID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var images []models.RoomImage
    for rows.Next() {
        var img models.RoomImage
        err := rows.Scan(
            &img.ID, &img.RoomID, &img.FilePath,
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

func (s *Storage) AddBedImage(ctx context.Context, image *models.RoomImage) (int, error) {
    var imageID int
    err := s.pool.QueryRow(ctx, `
        INSERT INTO bed_images (bed_id, file_path, file_name, file_size, content_type, is_main)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id
    `,
        image.BedID, image.FilePath, image.FileName,
        image.FileSize, image.ContentType, image.IsMain,
    ).Scan(&imageID)

    return imageID, err
}