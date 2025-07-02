package service

import (
	"context"
	"fmt"
	"time"

	"backend/internal/domain/models"
	"backend/internal/storage/postgres"
	"backend/pkg/logger"
)

// InventoryManager управляет резервированием и списанием товаров
type InventoryManager struct {
	inventoryRepo postgres.InventoryRepositoryInterface
	productRepo   ProductRepositoryInterface
	logger        logger.Logger
}

// NewInventoryManager создает новый менеджер инвентаря
func NewInventoryManager(
	inventoryRepo postgres.InventoryRepositoryInterface,
	productRepo ProductRepositoryInterface,
	logger logger.Logger,
) *InventoryManager {
	return &InventoryManager{
		inventoryRepo: inventoryRepo,
		productRepo:   productRepo,
		logger:        logger,
	}
}

// ReserveStock резервирует товар для заказа
func (im *InventoryManager) ReserveStock(ctx context.Context, productID int64, variantID *int64, quantity int, orderID int64) (*models.InventoryReservation, error) {
	im.logger.Info("Reserving stock", "product_id", productID, "variant_id", variantID, "quantity", quantity, "order_id", orderID)

	// Проверяем доступность товара
	stock, err := im.inventoryRepo.GetStock(ctx, productID, variantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get stock: %w", err)
	}

	availableStock := stock.AvailableQuantity
	if availableStock < quantity {
		return nil, fmt.Errorf("insufficient stock: requested %d, available %d", quantity, availableStock)
	}

	// Создаем резервирование
	reservation := &models.InventoryReservation{
		ProductID: productID,
		VariantID: variantID,
		OrderID:   orderID,
		Quantity:  quantity,
		Status:    models.ReservationStatusActive,
		ExpiresAt: time.Now().Add(2 * time.Hour), // Резерв на 2 часа
	}

	// Используем ReserveStock вместо CreateReservation
	err = im.inventoryRepo.ReserveStock(ctx, reservation)
	if err != nil {
		return nil, fmt.Errorf("failed to create reservation: %w", err)
	}

	im.logger.Info("Stock reserved successfully", "product_id", productID, "quantity", quantity)
	return reservation, nil
}

// CommitReservation подтверждает резервирование и списывает товар
func (im *InventoryManager) CommitReservation(ctx context.Context, reservationID int64) error {
	im.logger.Info("Committing reservation", "reservation_id", reservationID)

	// Используем ConfirmReservation для подтверждения
	err := im.inventoryRepo.ConfirmReservation(ctx, reservationID)
	if err != nil {
		return fmt.Errorf("failed to confirm reservation: %w", err)
	}

	im.logger.Info("Reservation committed successfully", "reservation_id", reservationID)
	return nil
}

// ReleaseReservation освобождает зарезервированный товар
func (im *InventoryManager) ReleaseReservation(ctx context.Context, reservationID int64) error {
	im.logger.Info("Releasing reservation", "reservation_id", reservationID)

	// Используем ReleaseReservation для освобождения
	err := im.inventoryRepo.ReleaseReservation(ctx, reservationID)
	if err != nil {
		return fmt.Errorf("failed to release reservation: %w", err)
	}

	im.logger.Info("Reservation released successfully", "reservation_id", reservationID)
	return nil
}

// CleanupExpiredReservations освобождает истекшие резервирования
func (im *InventoryManager) CleanupExpiredReservations(ctx context.Context) error {
	im.logger.Info("Cleaning up expired reservations")

	// Получаем истекшие резервирования
	expiredReservations, err := im.inventoryRepo.GetExpiredReservations(ctx)
	if err != nil {
		return fmt.Errorf("failed to get expired reservations: %w", err)
	}

	// Освобождаем каждое резервирование
	for _, reservation := range expiredReservations {
		if err := im.ReleaseReservation(ctx, reservation.ID); err != nil {
			im.logger.Error("Failed to release expired reservation", "reservation_id", reservation.ID, "error", err)
			// Продолжаем с остальными
		}
	}

	im.logger.Info("Cleanup completed", "released_count", len(expiredReservations))
	return nil
}

// ReleaseOrderReservations освобождает все резервирования для заказа
func (im *InventoryManager) ReleaseOrderReservations(ctx context.Context, orderID int64) error {
	im.logger.Info("Releasing all reservations for order", "order_id", orderID)

	// Получаем все резервирования для заказа
	reservations, err := im.inventoryRepo.GetReservationsByOrder(ctx, orderID)
	if err != nil {
		return fmt.Errorf("failed to get order reservations: %w", err)
	}

	// Освобождаем каждое резервирование
	for _, reservation := range reservations {
		if reservation.Status == models.ReservationStatusActive {
			if err := im.ReleaseReservation(ctx, reservation.ID); err != nil {
				im.logger.Error("Failed to release reservation", "reservation_id", reservation.ID, "error", err)
				// Продолжаем с остальными
			}
		}
	}

	im.logger.Info("All reservations released for order", "order_id", orderID, "count", len(reservations))
	return nil
}

// CommitOrderReservations подтверждает все резервирования для заказа
func (im *InventoryManager) CommitOrderReservations(ctx context.Context, orderID int64) error {
	im.logger.Info("Committing all reservations for order", "order_id", orderID)

	// Получаем все резервирования для заказа
	reservations, err := im.inventoryRepo.GetReservationsByOrder(ctx, orderID)
	if err != nil {
		return fmt.Errorf("failed to get order reservations: %w", err)
	}

	// Подтверждаем каждое резервирование
	for _, reservation := range reservations {
		if reservation.Status == models.ReservationStatusActive {
			if err := im.CommitReservation(ctx, reservation.ID); err != nil {
				im.logger.Error("Failed to commit reservation", "reservation_id", reservation.ID, "error", err)
				// Откатываем предыдущие
				im.ReleaseOrderReservations(ctx, orderID)
				return fmt.Errorf("failed to commit reservation %d: %w", reservation.ID, err)
			}
		}
	}

	im.logger.Info("All reservations committed for order", "order_id", orderID, "count", len(reservations))
	return nil
}

// GetAvailableStock получает доступное количество товара
func (im *InventoryManager) GetAvailableStock(ctx context.Context, productID int64, variantID *int64) (int, error) {
	stock, err := im.inventoryRepo.GetStock(ctx, productID, variantID)
	if err != nil {
		return 0, fmt.Errorf("failed to get stock: %w", err)
	}

	return stock.AvailableQuantity, nil
}

// UpdateStock обновляет запас товара
func (im *InventoryManager) UpdateStock(ctx context.Context, productID int64, variantID *int64, quantity int) error {
	return im.inventoryRepo.UpdateStock(ctx, productID, variantID, quantity)
}

// RecordMovement записывает движение товара
func (im *InventoryManager) RecordMovement(ctx context.Context, productID int64, variantID *int64, movementType string, quantity int, referenceType string, referenceID int64, notes string) error {
	movement := &models.InventoryMovement{
		ProductID:     productID,
		VariantID:     variantID,
		Type:          movementType,
		Quantity:      quantity,
		ReferenceType: &referenceType,
		ReferenceID:   &referenceID,
		Notes:         &notes,
	}

	return im.inventoryRepo.RecordMovement(ctx, movement)
}

// GetLowStockItems получает товары с низким запасом
func (im *InventoryManager) GetLowStockItems(ctx context.Context, storefrontID int) ([]models.LowStockItem, error) {
	return im.inventoryRepo.GetLowStockItems(ctx, storefrontID)
}

// GetStockMovements получает историю движений товара
func (im *InventoryManager) GetStockMovements(ctx context.Context, productID int64, variantID *int64, limit int, offset int) ([]models.InventoryMovementDTO, error) {
	// Этот метод требует отдельной реализации в репозитории
	// Пока возвращаем пустой список
	return []models.InventoryMovementDTO{}, nil
}
