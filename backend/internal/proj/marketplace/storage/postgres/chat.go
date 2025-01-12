package postgres

import (
    "context"
    "backend/internal/domain/models"
    "fmt"
)

// Реализуем методы хранилища
func (s *Storage) GetChat(ctx context.Context, chatID int, userID int) (*models.MarketplaceChat, error) {
    chat := &models.MarketplaceChat{}
    err := s.pool.QueryRow(ctx, `
        SELECT 
            c.id, c.listing_id, c.buyer_id, c.seller_id,
            c.last_message_at, c.created_at, c.updated_at, c.is_archived,
            l.title,
            (
                SELECT COUNT(*)
                FROM marketplace_messages m
                WHERE m.chat_id = c.id AND m.receiver_id = $2 AND NOT m.is_read
            ) as unread_count
        FROM marketplace_chats c
        JOIN marketplace_listings l ON c.listing_id = l.id
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
    return chat, nil
}

func (s *Storage) GetChats(ctx context.Context, userID int) ([]models.MarketplaceChat, error) {
    rows, err := s.pool.Query(ctx, `
        WITH unread_counts AS (
            SELECT 
                c.id as chat_id,
                COUNT(*) as unread_count
            FROM marketplace_chats c
            JOIN marketplace_messages m ON m.chat_id = c.id
            WHERE m.receiver_id = $1 AND NOT m.is_read
            GROUP BY c.id
        )
        SELECT 
            c.id, c.listing_id, c.buyer_id, c.seller_id,
            c.last_message_at, c.created_at, c.updated_at, c.is_archived,
            l.title,
            COALESCE(uc.unread_count, 0) as unread_count
        FROM marketplace_chats c
        JOIN marketplace_listings l ON c.listing_id = l.id
        LEFT JOIN unread_counts uc ON c.id = uc.chat_id
        WHERE c.buyer_id = $1 OR c.seller_id = $1
        ORDER BY c.last_message_at DESC
    `, userID)
    if err != nil {
        return nil, fmt.Errorf("error querying chats: %w", err)
    }
    defer rows.Close()

    var chats []models.MarketplaceChat
    for rows.Next() {
        var chat models.MarketplaceChat
        chat.Listing = &models.MarketplaceListing{}

        err := rows.Scan(
            &chat.ID, &chat.ListingID, &chat.BuyerID, &chat.SellerID,
            &chat.LastMessageAt, &chat.CreatedAt, &chat.UpdatedAt, &chat.IsArchived,
            &chat.Listing.Title,
            &chat.UnreadCount,
        )
        if err != nil {
            return nil, fmt.Errorf("error scanning chat: %w", err)
        }

        chats = append(chats, chat)
    }

    return chats, nil
}

func (s *Storage) CreateMessage(ctx context.Context, msg *models.MarketplaceMessage) error {
    tx, err := s.pool.Begin(ctx)
    if err != nil {
        return fmt.Errorf("error starting transaction: %w", err)
    }
    defer tx.Rollback(ctx)

    // Создаем или получаем существующий чат
    var chatID int
    err = tx.QueryRow(ctx, `
        INSERT INTO marketplace_chats (
            listing_id, buyer_id, seller_id, last_message_at
        ) VALUES (
            $1, 
            CASE WHEN $2 = (SELECT user_id FROM marketplace_listings WHERE id = $1) 
            THEN $3 ELSE $2 END,
            (SELECT user_id FROM marketplace_listings WHERE id = $1),
            CURRENT_TIMESTAMP
        )
        ON CONFLICT (listing_id, buyer_id, seller_id) 
        DO UPDATE SET 
            last_message_at = CURRENT_TIMESTAMP
        RETURNING id
    `, msg.ListingID, msg.SenderID, msg.ReceiverID).Scan(&chatID)
    if err != nil {
        return fmt.Errorf("error creating/getting chat: %w", err)
    }

    // Создаем сообщение с chat_id
    err = tx.QueryRow(ctx, `
        INSERT INTO marketplace_messages (
            chat_id, listing_id, sender_id, receiver_id, content
        ) VALUES ($1, $2, $3, $4, $5)
        RETURNING id, created_at
    `,
        chatID, msg.ListingID, msg.SenderID, msg.ReceiverID, msg.Content,
    ).Scan(&msg.ID, &msg.CreatedAt)
    if err != nil {
        return fmt.Errorf("error creating message: %w", err)
    }

    msg.ChatID = chatID
    return tx.Commit(ctx)
}

func (s *Storage) GetMessages(ctx context.Context, listingID int, userID int, offset int, limit int) ([]models.MarketplaceMessage, error) {
    messages := []models.MarketplaceMessage{} // Инициализируем пустой массив

    rows, err := s.pool.Query(ctx, `
        SELECT 
            m.id, m.chat_id, m.listing_id, m.sender_id, m.receiver_id,
            m.content, m.is_read, m.created_at,
            u.name as sender_name, u.picture_url as sender_picture
        FROM marketplace_messages m
        JOIN users u ON m.sender_id = u.id
        JOIN marketplace_chats c ON m.chat_id = c.id
        WHERE m.listing_id = $1 
        AND (c.buyer_id = $2 OR c.seller_id = $2)
        ORDER BY m.created_at DESC
        LIMIT $3 OFFSET $4
    `, listingID, userID, limit, offset)
    if err != nil {
        return messages, fmt.Errorf("error querying messages: %w", err)
    }
    defer rows.Close()

    for rows.Next() {
        var msg models.MarketplaceMessage
        msg.Sender = &models.User{}
        
        err := rows.Scan(
            &msg.ID, &msg.ChatID, &msg.ListingID, &msg.SenderID, &msg.ReceiverID,
            &msg.Content, &msg.IsRead, &msg.CreatedAt,
            &msg.Sender.Name, &msg.Sender.PictureURL,
        )
        if err != nil {
            return messages, fmt.Errorf("error scanning message: %w", err)
        }
        
        messages = append(messages, msg)
    }

    return messages, nil
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