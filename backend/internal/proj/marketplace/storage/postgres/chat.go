// backend/internal/proj/marketplace/storage/postgres/chat.go
package postgres

import (
	"backend/internal/domain/models"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/jackc/pgx/v5"
)

// Реализуем методы хранилища
func (s *Storage) GetChat(ctx context.Context, chatID int, userID int) (*models.MarketplaceChat, error) {
	chat := &models.MarketplaceChat{}
	chat.Listing = &models.MarketplaceListing{}

	// Сначала пытаемся получить чат с листингом
	err := s.pool.QueryRow(ctx, `
        SELECT
            c.id, COALESCE(c.listing_id, 0), c.buyer_id, c.seller_id,
            c.last_message_at, c.created_at, c.updated_at, c.is_archived,
            CASE 
                WHEN c.listing_id IS NULL THEN 'Личное сообщение'
                WHEN l.id IS NULL THEN 'Удаленное объявление'
                ELSE l.title
            END as listing_title,
            (
                SELECT COUNT(*)
                FROM marketplace_messages m
                WHERE m.chat_id = c.id AND m.receiver_id = $2 AND NOT m.is_read
            ) as unread_count
        FROM marketplace_chats c
        LEFT JOIN marketplace_listings l ON c.listing_id = l.id
        WHERE c.id = $1 AND (c.buyer_id = $2 OR c.seller_id = $2)
    `, chatID, userID).Scan(
		&chat.ID, &chat.ListingID, &chat.BuyerID, &chat.SellerID,
		&chat.LastMessageAt, &chat.CreatedAt, &chat.UpdatedAt, &chat.IsArchived,
		&chat.Listing.Title,
		&chat.UnreadCount,
	)

	if err != nil {
		return nil, fmt.Errorf("error getting chat: %w", err)
	}

	// Получаем имена и аватарки пользователей
	var sellerName, sellerPicture, buyerName, buyerPicture sql.NullString

	err = s.pool.QueryRow(ctx, `
	    SELECT
	        seller.name, seller.picture_url,
	        buyer.name, buyer.picture_url
	    FROM marketplace_chats c
	    LEFT JOIN users seller ON c.seller_id = seller.id
	    LEFT JOIN users buyer ON c.buyer_id = buyer.id
	    WHERE c.id = $1
	`, chatID).Scan(
		&sellerName, &sellerPicture,
		&buyerName, &buyerPicture,
	)

	if err != nil {
		return nil, fmt.Errorf("error getting user details: %w", err)
	}

	// Добавляем информацию о продавце и покупателе
	chat.Seller = &models.User{
		ID:         chat.SellerID,
		Name:       sellerName.String,
		PictureURL: sellerPicture.String,
	}

	chat.Buyer = &models.User{
		ID:         chat.BuyerID,
		Name:       buyerName.String,
		PictureURL: buyerPicture.String,
	}

	// Определяем, кто другой пользователь (не текущий)
	if chat.BuyerID == userID {
		chat.OtherUser = chat.Seller
	} else {
		chat.OtherUser = chat.Buyer
	}

	return chat, nil
}

func (s *Storage) GetChats(ctx context.Context, userID int) ([]models.MarketplaceChat, error) {
	// Единый запрос с LEFT JOIN для получения всех данных включая изображения
	query := `
	WITH unread_counts AS (
		SELECT
			c.id as chat_id,
			COUNT(*) as unread_count
		FROM marketplace_chats c
		JOIN marketplace_messages m ON m.chat_id = c.id
		WHERE m.receiver_id = $1 AND NOT m.is_read
		GROUP BY c.id
	),
	chat_images AS (
		SELECT 
			c.id as chat_id,
			json_agg(
				json_build_object(
					'id', mi.id,
					'listing_id', mi.listing_id,
					'file_path', mi.file_path,
					'file_name', mi.file_name,
					'file_size', mi.file_size,
					'content_type', mi.content_type,
					'is_main', mi.is_main,
					'storage_type', mi.storage_type,
					'storage_bucket', mi.storage_bucket,
					'public_url', mi.public_url,
					'created_at', mi.created_at
				) ORDER BY mi.is_main DESC, mi.id ASC
			) as images
		FROM marketplace_chats c
		LEFT JOIN marketplace_images mi ON mi.listing_id = c.listing_id
		WHERE c.buyer_id = $1 OR c.seller_id = $1
		GROUP BY c.id
	)
	SELECT
		c.id, COALESCE(c.listing_id, 0), c.buyer_id, c.seller_id,
		c.last_message_at, c.created_at, c.updated_at, c.is_archived,
		CASE 
			WHEN c.listing_id IS NULL THEN '__DIRECT_MESSAGE__'
			WHEN l.id IS NULL THEN '__DELETED_LISTING__'
			ELSE l.title
		END as listing_title,
		COALESCE(l.price, 0) as listing_price,
		COALESCE(uc.unread_count, 0) as unread_count,
		-- Информация о пользователе
		CASE
			WHEN c.buyer_id = $1 THEN seller.name
			ELSE buyer.name
		END as other_user_name,
		CASE
			WHEN c.buyer_id = $1 THEN seller.picture_url
			ELSE buyer.picture_url
		END as other_user_picture,
		-- Изображения листинга
		COALESCE(ci.images, '[]'::json) as listing_images
	FROM marketplace_chats c
	LEFT JOIN marketplace_listings l ON c.listing_id = l.id
	LEFT JOIN unread_counts uc ON c.id = uc.chat_id
	LEFT JOIN users buyer ON c.buyer_id = buyer.id
	LEFT JOIN users seller ON c.seller_id = seller.id
	LEFT JOIN chat_images ci ON c.id = ci.chat_id
	WHERE c.buyer_id = $1 OR c.seller_id = $1
	ORDER BY c.last_message_at DESC`

	rows, err := s.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("error querying chats: %w", err)
	}
	defer rows.Close()

	var chats []models.MarketplaceChat
	for rows.Next() {
		var (
			chat             models.MarketplaceChat
			otherUserName    sql.NullString
			otherUserPicture sql.NullString
			imagesJSON       json.RawMessage
		)
		chat.Listing = &models.MarketplaceListing{}

		err := rows.Scan(
			&chat.ID, &chat.ListingID, &chat.BuyerID, &chat.SellerID,
			&chat.LastMessageAt, &chat.CreatedAt, &chat.UpdatedAt, &chat.IsArchived,
			&chat.Listing.Title, &chat.Listing.Price,
			&chat.UnreadCount,
			&otherUserName,
			&otherUserPicture,
			&imagesJSON,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning chat: %w", err)
		}

		// Парсим изображения из JSON
		var images []models.MarketplaceImage
		if err := json.Unmarshal(imagesJSON, &images); err != nil {
			log.Printf("Error unmarshalling images for chat %d: %v", chat.ID, err)
			images = []models.MarketplaceImage{}
		}
		chat.Listing.Images = images

		// Создаем структуру other_user
		chat.OtherUser = &models.User{
			Name:       otherUserName.String,
			PictureURL: otherUserPicture.String,
		}

		// Добавляем информацию о продавце и покупателе для полноты данных
		if chat.BuyerID == userID {
			chat.Buyer = &models.User{ID: userID}
			chat.Seller = chat.OtherUser
			chat.Seller.ID = chat.SellerID
		} else {
			chat.Seller = &models.User{ID: userID}
			chat.Buyer = chat.OtherUser
			chat.Buyer.ID = chat.BuyerID
		}

		chats = append(chats, chat)
	}

	return chats, nil
}

func (s *Storage) GetMessages(ctx context.Context, listingID int, userID int, offset int, limit int) ([]models.MarketplaceMessage, error) {
	if offset < 0 {
		offset = 0
	}
	if limit <= 0 {
		limit = 20
	}

	// Получаем chatID из контекста, безопасное приведение типов
	var chatID int
	if chatIDValue := ctx.Value("chat_id"); chatIDValue != nil {
		switch v := chatIDValue.(type) {
		case int:
			chatID = v
		case string:
			var err error
			chatID, err = strconv.Atoi(v)
			if err != nil {
				chatID = 0
			}
		default:
			chatID = 0
		}
	}

	log.Printf("GetMessages storage: chatID from context = %d, listingID = %d, userID = %d", chatID, listingID, userID)

	var query string
	var args []interface{}

	if chatID > 0 {
		// Если известен ID чата, получаем сообщения напрямую по chatID
		// Это позволяет получать сообщения даже если листинг больше не существует
		query = `
		WITH ordered_messages AS (
			SELECT
				m.id, m.chat_id, m.listing_id, m.sender_id, m.receiver_id,
				m.content, m.is_read, m.created_at,
				sender.name as sender_name,
				sender.picture_url as sender_picture,
				receiver.name as receiver_name,
				receiver.picture_url as receiver_picture,
				m.has_attachments, m.attachments_count
			FROM marketplace_messages m
			JOIN marketplace_chats c ON m.chat_id = c.id
			JOIN users sender ON m.sender_id = sender.id
			JOIN users receiver ON m.receiver_id = receiver.id
			WHERE m.chat_id = $1
			AND (c.buyer_id = $2 OR c.seller_id = $2)
			ORDER BY m.created_at DESC
			LIMIT $3 OFFSET $4
		),
		message_attachments AS (
			SELECT 
				om.id as message_id,
				json_agg(
					json_build_object(
						'id', ca.id,
						'message_id', ca.message_id,
						'file_type', ca.file_type,
						'file_path', ca.file_path,
						'file_name', ca.file_name,
						'file_size', ca.file_size,
						'content_type', ca.content_type,
						'storage_type', ca.storage_type,
						'storage_bucket', ca.storage_bucket,
						'public_url', ca.public_url,
						'thumbnail_url', ca.thumbnail_url,
						'metadata', ca.metadata,
						'created_at', ca.created_at
					) ORDER BY ca.created_at ASC
				) as attachments
			FROM ordered_messages om
			LEFT JOIN chat_attachments ca ON ca.message_id = om.id
			WHERE om.has_attachments = true
			GROUP BY om.id
		)
		SELECT 
			om.*,
			COALESCE(ma.attachments, '[]'::json) as attachments_json
		FROM ordered_messages om
		LEFT JOIN message_attachments ma ON om.id = ma.message_id
		ORDER BY om.created_at ASC`
		args = []interface{}{chatID, userID, limit, offset}
	} else {
		// Если ID чата не известен, ищем чат по listingID и userID
		query = `
		WITH chat AS (
			SELECT c.id
			FROM marketplace_chats c
			WHERE c.listing_id = $1
			AND (c.buyer_id = $2 OR c.seller_id = $2)
			LIMIT 1
		),
		ordered_messages AS (
			SELECT
				m.id, m.chat_id, m.listing_id, m.sender_id, m.receiver_id,
				m.content, m.is_read, m.created_at,
				sender.name as sender_name,
				sender.picture_url as sender_picture,
				receiver.name as receiver_name,
				receiver.picture_url as receiver_picture,
				m.has_attachments, m.attachments_count
			FROM marketplace_messages m
			JOIN chat c ON m.chat_id = c.id
			JOIN users sender ON m.sender_id = sender.id
			JOIN users receiver ON m.receiver_id = receiver.id
			ORDER BY m.created_at DESC
			LIMIT $3 OFFSET $4
		),
		message_attachments AS (
			SELECT 
				om.id as message_id,
				json_agg(
					json_build_object(
						'id', ca.id,
						'message_id', ca.message_id,
						'file_type', ca.file_type,
						'file_path', ca.file_path,
						'file_name', ca.file_name,
						'file_size', ca.file_size,
						'content_type', ca.content_type,
						'storage_type', ca.storage_type,
						'storage_bucket', ca.storage_bucket,
						'public_url', ca.public_url,
						'thumbnail_url', ca.thumbnail_url,
						'metadata', ca.metadata,
						'created_at', ca.created_at
					) ORDER BY ca.created_at ASC
				) as attachments
			FROM ordered_messages om
			LEFT JOIN chat_attachments ca ON ca.message_id = om.id
			WHERE om.has_attachments = true
			GROUP BY om.id
		)
		SELECT 
			om.*,
			COALESCE(ma.attachments, '[]'::json) as attachments_json
		FROM ordered_messages om
		LEFT JOIN message_attachments ma ON om.id = ma.message_id
		ORDER BY om.created_at ASC`
		args = []interface{}{listingID, userID, limit, offset}
	}

	rows, err := s.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error querying messages: %w", err)
	}
	defer rows.Close()

	var messages []models.MarketplaceMessage
	for rows.Next() {
		var msg models.MarketplaceMessage
		var listingID sql.NullInt64
		var attachmentsJSON json.RawMessage
		msg.Sender = &models.User{}
		msg.Receiver = &models.User{}

		err := rows.Scan(
			&msg.ID, &msg.ChatID, &listingID, &msg.SenderID, &msg.ReceiverID,
			&msg.Content, &msg.IsRead, &msg.CreatedAt,
			&msg.Sender.Name, &msg.Sender.PictureURL,
			&msg.Receiver.Name, &msg.Receiver.PictureURL,
			&msg.HasAttachments, &msg.AttachmentsCount,
			&attachmentsJSON,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning message: %w", err)
		}

		// Устанавливаем ListingID только если он не NULL
		if listingID.Valid {
			msg.ListingID = int(listingID.Int64)
		} else {
			msg.ListingID = 0
		}

		// Парсим вложения из JSON
		if msg.HasAttachments && msg.AttachmentsCount > 0 {
			var attachments []models.ChatAttachment
			if err := json.Unmarshal(attachmentsJSON, &attachments); err != nil {
				log.Printf("Error unmarshalling attachments for message %d: %v", msg.ID, err)
				attachments = []models.ChatAttachment{}
			}
			msg.Attachments = attachments
		}

		messages = append(messages, msg)
	}

	return messages, nil
}

// GetMessageAttachments загружает вложения для сообщения
func (s *Storage) GetMessageAttachments(ctx context.Context, messageID int) ([]*models.ChatAttachment, error) {
	query := `
		SELECT 
			id, message_id, file_type, file_path, file_name,
			file_size, content_type, storage_type, storage_bucket,
			public_url, thumbnail_url, metadata, created_at
		FROM chat_attachments
		WHERE message_id = $1
		ORDER BY created_at ASC
	`

	rows, err := s.pool.Query(ctx, query, messageID)
	if err != nil {
		return nil, fmt.Errorf("error getting message attachments: %w", err)
	}
	defer rows.Close()

	var attachments []*models.ChatAttachment

	for rows.Next() {
		attachment := &models.ChatAttachment{}
		err := rows.Scan(
			&attachment.ID,
			&attachment.MessageID,
			&attachment.FileType,
			&attachment.FilePath,
			&attachment.FileName,
			&attachment.FileSize,
			&attachment.ContentType,
			&attachment.StorageType,
			&attachment.StorageBucket,
			&attachment.PublicURL,
			&attachment.ThumbnailURL,
			&attachment.Metadata,
			&attachment.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning attachment: %w", err)
		}

		attachments = append(attachments, attachment)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating attachments: %w", err)
	}

	return attachments, nil
}

// GetMessagesCount возвращает общее количество сообщений в чате
func (s *Storage) GetMessagesCount(ctx context.Context, chatID int) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM marketplace_messages WHERE chat_id = $1`
	err := s.pool.QueryRow(ctx, query, chatID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("error counting messages: %w", err)
	}
	return count, nil
}

// In chat.go
func (s *Storage) CreateMessage(ctx context.Context, msg *models.MarketplaceMessage) error {
	if msg == nil {
		return fmt.Errorf("message cannot be nil")
	}

	// Initialize user structs if nil
	if msg.Sender == nil {
		msg.Sender = &models.User{}
	}
	if msg.Receiver == nil {
		msg.Receiver = &models.User{}
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Validate input data
	if msg.SenderID == 0 || msg.ReceiverID == 0 {
		return fmt.Errorf("invalid message data: sender_id=%d, receiver_id=%d",
			msg.SenderID, msg.ReceiverID)
	}

	// Для прямых сообщений (без привязки к объявлению) разрешаем ListingID == 0
	// Но если есть ChatID, проверяем, что чат существует

	// Проверяем из контекста, существует ли листинг
	var listingExists bool = true
	if listingExistsValue := ctx.Value("listing_exists"); listingExistsValue != nil {
		if exists, ok := listingExistsValue.(bool); ok {
			listingExists = exists
		}
	}

	var sellerID, buyerID int

	// Для прямых сообщений без ChatID и ListingID сначала пытаемся найти существующий чат
	if msg.ChatID == 0 && msg.ListingID == 0 {
		// Определяем buyerID и sellerID для прямого чата
		if msg.SenderID < msg.ReceiverID {
			buyerID = msg.SenderID
			sellerID = msg.ReceiverID
		} else {
			buyerID = msg.ReceiverID
			sellerID = msg.SenderID
		}

		// Ищем существующий прямой чат
		err = tx.QueryRow(ctx, `
			SELECT id
			FROM marketplace_chats
			WHERE listing_id IS NULL
			AND ((buyer_id = $1 AND seller_id = $2) OR (buyer_id = $2 AND seller_id = $1))
			LIMIT 1
		`, buyerID, sellerID).Scan(&msg.ChatID)

		if err != nil && err != pgx.ErrNoRows {
			return fmt.Errorf("error checking existing direct chat: %w", err)
		}
		// Если нашли чат, msg.ChatID теперь заполнен
	}

	// Определяем sellerID и buyerID
	if msg.ChatID > 0 {
		// Если это существующий чат, получаем информацию из таблицы чатов
		// вместо обращения к таблице листингов
		err = tx.QueryRow(ctx, `
            SELECT buyer_id, seller_id
            FROM marketplace_chats
            WHERE id = $1
        `, msg.ChatID).Scan(&buyerID, &sellerID)

		if err != nil {
			if err == pgx.ErrNoRows {
				return fmt.Errorf("chat not found: %d", msg.ChatID)
			}
			return fmt.Errorf("error getting chat info: %w", err)
		}
	} else {
		// Только для нового чата проверяем листинг
		if msg.ListingID > 0 && listingExists {
			// Get seller ID and check listing existence
			err = tx.QueryRow(ctx, `
                SELECT user_id
                FROM marketplace_listings
                WHERE id = $1
            `, msg.ListingID).Scan(&sellerID)

			if err != nil {
				if err == pgx.ErrNoRows {
					return fmt.Errorf("listing not found: %d", msg.ListingID)
				}
				return fmt.Errorf("error getting listing seller: %w", err)
			}

			// Determine buyer
			if msg.SenderID == sellerID {
				buyerID = msg.ReceiverID
			} else {
				buyerID = msg.SenderID
			}
		} else {
			// Это прямое сообщение без привязки к объявлению или
			// существующий чат с удаленным листингом
			if msg.ReceiverID == msg.SenderID {
				return fmt.Errorf("sender and receiver cannot be the same")
			}

			// Для прямых сообщений используем SenderID и ReceiverID напрямую
			// Условно назначаем seller и buyer для совместимости со структурой БД
			if msg.SenderID < msg.ReceiverID {
				sellerID = msg.SenderID
				buyerID = msg.ReceiverID
			} else {
				sellerID = msg.ReceiverID
				buyerID = msg.SenderID
			}
		}
	}

	// Add NULL handling for user information
	type userInfo struct {
		name       sql.NullString
		pictureURL sql.NullString
	}
	var senderInfo, receiverInfo userInfo

	// Create or get existing chat with proper NULL handling
	if msg.ChatID > 0 {
		// Если ChatID уже известен, обновляем last_message_at и устанавливаем is_archived = false
		_, err = tx.Exec(ctx, `
            UPDATE marketplace_chats
            SET last_message_at = CURRENT_TIMESTAMP, is_archived = false
            WHERE id = $1
        `, msg.ChatID)

		if err != nil {
			return fmt.Errorf("error updating chat: %w", err)
		}

		// Получаем информацию о пользователях
		err = tx.QueryRow(ctx, `
            WITH user_info AS (
                SELECT id, name, picture_url
                FROM users
                WHERE id IN ($1, $2)
            )
            SELECT
                (SELECT name FROM user_info WHERE id = $1),
                (SELECT picture_url FROM user_info WHERE id = $1),
                (SELECT name FROM user_info WHERE id = $2),
                (SELECT picture_url FROM user_info WHERE id = $2)
        `, msg.SenderID, msg.ReceiverID).Scan(
			&senderInfo.name,
			&senderInfo.pictureURL,
			&receiverInfo.name,
			&receiverInfo.pictureURL,
		)

		if err != nil {
			return fmt.Errorf("error getting user info: %w", err)
		}
	} else {
		// Создаем новый чат или получаем существующий
		err = tx.QueryRow(ctx, `
            WITH user_info AS (
                SELECT id, name, picture_url
                FROM users
                WHERE id IN ($1, $2)
            ),
            chat_insert AS (
                INSERT INTO marketplace_chats (
                    listing_id,
                    buyer_id,
                    seller_id,
                    last_message_at
                ) VALUES (NULLIF($3, 0), $4, $5, CURRENT_TIMESTAMP)
                ON CONFLICT (listing_id, buyer_id, seller_id)
                DO UPDATE SET
                    last_message_at = CURRENT_TIMESTAMP,
                    is_archived = false
                RETURNING id
            )
            SELECT
                ci.id,
                (SELECT name FROM user_info WHERE id = $1),
                (SELECT picture_url FROM user_info WHERE id = $1),
                (SELECT name FROM user_info WHERE id = $2),
                (SELECT picture_url FROM user_info WHERE id = $2)
            FROM chat_insert ci
        `,
			msg.SenderID, msg.ReceiverID,
			msg.ListingID, buyerID, sellerID,
		).Scan(
			&msg.ChatID,
			&senderInfo.name,
			&senderInfo.pictureURL,
			&receiverInfo.name,
			&receiverInfo.pictureURL,
		)

		if err != nil {
			return fmt.Errorf("error creating/getting chat: %w", err)
		}
	}

	// Set user information with NULL handling
	msg.Sender.Name = senderInfo.name.String
	msg.Sender.PictureURL = senderInfo.pictureURL.String
	msg.Receiver.Name = receiverInfo.name.String
	msg.Receiver.PictureURL = receiverInfo.pictureURL.String

	// Create message
	err = tx.QueryRow(ctx, `
        INSERT INTO marketplace_messages (
            chat_id,
            listing_id,
            sender_id,
            receiver_id,
            content,
            is_read,
            original_language
        ) VALUES ($1, NULLIF($2, 0), $3, $4, $5, false, $6)
        RETURNING id, created_at
    `,
		msg.ChatID,
		msg.ListingID,
		msg.SenderID,
		msg.ReceiverID,
		msg.Content,
		msg.OriginalLanguage,
	).Scan(&msg.ID, &msg.CreatedAt)

	if err != nil {
		return fmt.Errorf("error creating message: %w", err)
	}

	return tx.Commit(ctx)
}
func (s *Storage) MarkMessagesAsRead(ctx context.Context, messageIDs []int, userID int) error {
	_, err := s.pool.Exec(ctx, `
        UPDATE marketplace_messages
        SET is_read = true
        WHERE id = ANY($1) AND receiver_id = $2
    `, messageIDs, userID)
	if err != nil {
		return fmt.Errorf("error marking messages as read: %w", err)
	}
	return nil
}

func (s *Storage) ArchiveChat(ctx context.Context, chatID int, userID int) error {
	result, err := s.pool.Exec(ctx, `
        UPDATE marketplace_chats
        SET is_archived = true
        WHERE id = $1 AND (buyer_id = $2 OR seller_id = $2)
    `, chatID, userID)
	if err != nil {
		return fmt.Errorf("error archiving chat: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("chat not found or permission denied")
	}

	return nil
}
