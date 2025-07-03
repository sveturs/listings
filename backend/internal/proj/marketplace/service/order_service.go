package service

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"

	"backend/internal/domain/models"
	"backend/internal/proj/marketplace/repository"
	"backend/internal/storage"
	// Удаляем импорт payments service, используем наш локальный
)

// OrderService предоставляет бизнес-логику для работы с заказами
type OrderService struct {
	orderRepo       repository.OrderRepositoryInterface
	listingRepo     repository.MarketplaceRepositoryInterface
	userRepo        repository.UserRepositoryInterface
	paymentService  *PaymentService
	notificationSvc NotificationService // Интерфейс для уведомлений
	platformFeeRate float64             // Комиссия платформы в процентах
}

// NotificationService интерфейс для отправки уведомлений
type NotificationService interface {
	SendOrderCreated(ctx context.Context, order *models.MarketplaceOrder) error
	SendOrderPaid(ctx context.Context, order *models.MarketplaceOrder) error
	SendOrderShipped(ctx context.Context, order *models.MarketplaceOrder) error
	SendOrderDelivered(ctx context.Context, order *models.MarketplaceOrder) error
	SendProtectionExpiring(ctx context.Context, order *models.MarketplaceOrder) error
	SendPaymentReleased(ctx context.Context, order *models.MarketplaceOrder) error
}

// NewOrderService создает новый сервис заказов
func NewOrderService(
	orderRepo repository.OrderRepositoryInterface,
	listingRepo repository.MarketplaceRepositoryInterface,
	userRepo repository.UserRepositoryInterface,
	paymentService *PaymentService,
	notificationSvc NotificationService,
	platformFeeRate float64,
) *OrderService {
	return &OrderService{
		orderRepo:       orderRepo,
		listingRepo:     listingRepo,
		userRepo:        userRepo,
		paymentService:  paymentService,
		notificationSvc: notificationSvc,
		platformFeeRate: platformFeeRate,
	}
}

// NewSimpleOrderService создает упрощенный сервис заказов только с storage
func NewSimpleOrderService(storage storage.Storage) OrderServiceInterface {
	// Создаем адаптеры для репозиториев (как в service.go:27-30)
	if marketplaceOrderRepo := storage.MarketplaceOrder(); marketplaceOrderRepo != nil {
		orderRepoAdapter := repository.NewPostgresOrderAdapter(marketplaceOrderRepo)
		listingRepoAdapter := repository.NewPostgresMarketplaceAdapter(storage)
		userRepoAdapter := &SimpleUserRepository{storage: storage}

		return &OrderService{
			orderRepo:       orderRepoAdapter,
			listingRepo:     listingRepoAdapter,
			userRepo:        userRepoAdapter,
			paymentService:  NewPaymentService(), // Создаем платежный сервис
			notificationSvc: nil,                 // Инициализируем позже если нужно
			platformFeeRate: 0.05,                // 5% комиссия по умолчанию
		}
	}

	// Если репозиторий недоступен, возвращаем заглушку
	return &EmptyOrderService{}
}

// EmptyOrderService заглушка для случаев когда репозиторий недоступен
type EmptyOrderService struct{}

func (e *EmptyOrderService) CreateOrder(ctx context.Context, order *models.MarketplaceOrder) (*models.MarketplaceOrder, error) {
	return nil, fmt.Errorf("order service not available")
}

func (e *EmptyOrderService) GetOrder(ctx context.Context, orderID int) (*models.MarketplaceOrder, error) {
	return nil, fmt.Errorf("order service not available")
}

func (e *EmptyOrderService) GetOrdersByUser(ctx context.Context, userID int, isPurchaser bool) ([]models.MarketplaceOrder, error) {
	return nil, fmt.Errorf("order service not available")
}

func (e *EmptyOrderService) UpdateOrderStatus(ctx context.Context, orderID int, status string) error {
	return fmt.Errorf("order service not available")
}

func (e *EmptyOrderService) CreateOrderFromRequest(ctx context.Context, req CreateOrderRequest) (*models.MarketplaceOrder, *PaymentResult, error) {
	return nil, nil, fmt.Errorf("order service not available")
}

func (e *EmptyOrderService) GetBuyerOrders(ctx context.Context, buyerID int64, page, limit int) ([]*models.MarketplaceOrder, int, error) {
	return nil, 0, fmt.Errorf("order service not available")
}

func (e *EmptyOrderService) GetSellerOrders(ctx context.Context, sellerID int64, page, limit int) ([]*models.MarketplaceOrder, int, error) {
	return nil, 0, fmt.Errorf("order service not available")
}

func (e *EmptyOrderService) GetOrderDetails(ctx context.Context, orderID int64, userID int64) (*models.MarketplaceOrder, error) {
	return nil, fmt.Errorf("order service not available")
}

func (e *EmptyOrderService) MarkAsShipped(ctx context.Context, orderID int64, sellerID int64, shippingMethod string, trackingNumber string) error {
	return fmt.Errorf("order service not available")
}

func (e *EmptyOrderService) ConfirmDelivery(ctx context.Context, orderID int64, buyerID int64) error {
	return fmt.Errorf("order service not available")
}

func (e *EmptyOrderService) OpenDispute(ctx context.Context, orderID int64, userID int64, reason string) error {
	return fmt.Errorf("order service not available")
}

func (e *EmptyOrderService) ConfirmPayment(ctx context.Context, orderID int64) error {
	return fmt.Errorf("order service not available")
}

// Добавляем методы с правильными сигнатурами для интерфейса OrderServiceInterface

// CreateOrder создает заказ из модели (для интерфейса)
func (s *OrderService) CreateOrder(ctx context.Context, order *models.MarketplaceOrder) (*models.MarketplaceOrder, error) {
	// Преобразуем модель в запрос и вызываем основной метод
	req := CreateOrderRequest{
		BuyerID:   int64(order.BuyerID),
		ListingID: int64(order.ListingID),
	}

	createdOrder, _, err := s.CreateOrderFromRequest(ctx, req)
	return createdOrder, err
}

// GetOrder получает заказ по ID (для интерфейса)
func (s *OrderService) GetOrder(ctx context.Context, orderID int) (*models.MarketplaceOrder, error) {
	if s.orderRepo == nil {
		return nil, fmt.Errorf("order repository not available")
	}
	return s.orderRepo.GetByID(ctx, int64(orderID))
}

// GetOrdersByUser получает заказы пользователя (для интерфейса)
func (s *OrderService) GetOrdersByUser(ctx context.Context, userID int, isPurchaser bool) ([]models.MarketplaceOrder, error) {
	if s.orderRepo == nil {
		return nil, fmt.Errorf("order repository not available")
	}

	var orders []*models.MarketplaceOrder
	var err error

	if isPurchaser {
		orders, _, err = s.orderRepo.GetBuyerOrders(ctx, int64(userID), 100, 0)
	} else {
		orders, _, err = s.orderRepo.GetSellerOrders(ctx, int64(userID), 100, 0)
	}

	if err != nil {
		return nil, err
	}

	// Преобразуем []*models.MarketplaceOrder в []models.MarketplaceOrder
	result := make([]models.MarketplaceOrder, len(orders))
	for i, order := range orders {
		if order != nil {
			result[i] = *order
		}
	}

	return result, nil
}

// UpdateOrderStatus обновляет статус заказа (для интерфейса)
func (s *OrderService) UpdateOrderStatus(ctx context.Context, orderID int, status string) error {
	if s.orderRepo == nil {
		return fmt.Errorf("order repository not available")
	}

	// Преобразуем строку в enum статуса
	orderStatus := models.MarketplaceOrderStatus(status)
	return s.orderRepo.UpdateStatus(ctx, int64(orderID), orderStatus, "Updated via API", nil)
}

// CreateOrderRequest запрос на создание заказа
// CreateOrderFromRequest создает новый заказ и инициирует платеж
func (s *OrderService) CreateOrderFromRequest(ctx context.Context, req CreateOrderRequest) (*models.MarketplaceOrder, *PaymentResult, error) {
	// 1. Получаем листинг
	listing, err := s.listingRepo.GetListingByID(ctx, int(req.ListingID))
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to get listing")
	}

	// Проверки
	if listing.Status != "active" {
		return nil, nil, errors.New("listing is not active")
	}
	if listing.UserID == int(req.BuyerID) {
		return nil, nil, errors.New("cannot buy own listing")
	}

	// 2. Создаем заказ
	order := &models.MarketplaceOrder{
		BuyerID:              req.BuyerID,
		SellerID:             int64(listing.UserID),
		ListingID:            req.ListingID,
		ItemPrice:            listing.Price,
		PlatformFeeRate:      s.platformFeeRate,
		Status:               models.MarketplaceOrderStatusPending,
		ProtectionPeriodDays: 7, // По умолчанию 7 дней
	}

	// Рассчитываем комиссии
	order.CalculateFees()

	// Сохраняем заказ
	err = s.orderRepo.Create(ctx, order)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to create order")
	}

	// 3. Создаем платеж с preauthorization
	paymentReq := CreatePaymentRequest{
		UserID:        int(req.BuyerID),
		Amount:        decimal.NewFromFloat(order.ItemPrice),
		Currency:      "RSD",
		PaymentMethod: req.PaymentMethod,
		Description:   fmt.Sprintf("Order #%d: %s", order.ID, listing.Title),
		ReturnURL:     req.ReturnURL,
		MetaData: map[string]interface{}{
			"order_id":     order.ID,
			"listing_id":   listing.ID,
			"seller_id":    listing.UserID,
			"capture_mode": "manual", // Важно! Используем manual capture
		},
		// Новые поля для delayed capture
		CaptureMode:   "manual",
		AutoCaptureIn: 0, // Не используем авто-захват
	}

	paymentResult, err := s.paymentService.CreatePayment(ctx, paymentReq)
	if err != nil {
		// Отменяем заказ если платеж не создался
		s.orderRepo.UpdateStatus(ctx, order.ID, models.MarketplaceOrderStatusCancelled, "Payment creation failed", nil)
		return nil, nil, errors.Wrap(err, "failed to create payment")
	}

	// 4. Связываем заказ с платежом
	order.PaymentTransactionID = &paymentResult.TransactionID
	// TODO: Обновить payment_transaction_id в БД

	// 5. Отправляем уведомление продавцу
	go s.notificationSvc.SendOrderCreated(ctx, order)

	return order, paymentResult, nil
}

// ConfirmPayment подтверждает оплату заказа
func (s *OrderService) ConfirmPayment(ctx context.Context, orderID int64) error {
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return errors.Wrap(err, "failed to get order")
	}

	if order.Status != models.MarketplaceOrderStatusPending {
		return errors.New("order is not in pending status")
	}

	// Обновляем статус на оплачен
	err = s.orderRepo.UpdateStatus(ctx, orderID, models.MarketplaceOrderStatusPaid, "Payment confirmed", nil)
	if err != nil {
		return errors.Wrap(err, "failed to update order status")
	}

	// Отправляем уведомления
	go s.notificationSvc.SendOrderPaid(ctx, order)

	return nil
}

// MarkAsShipped отмечает заказ как отправленный
func (s *OrderService) MarkAsShipped(ctx context.Context, orderID int64, sellerID int64, shippingMethod string, trackingNumber string) error {
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return errors.Wrap(err, "failed to get order")
	}

	// Проверки
	if order.SellerID != sellerID {
		return errors.New("unauthorized: not the seller")
	}
	if order.Status != models.MarketplaceOrderStatusPaid {
		return errors.New("order must be paid before shipping")
	}

	// Обновляем информацию о доставке
	err = s.orderRepo.UpdateShippingInfo(ctx, orderID, shippingMethod, trackingNumber)
	if err != nil {
		return errors.Wrap(err, "failed to update shipping info")
	}

	// Обновляем статус
	err = s.orderRepo.UpdateStatus(ctx, orderID, models.MarketplaceOrderStatusShipped, "Shipped by seller", &sellerID)
	if err != nil {
		return errors.Wrap(err, "failed to update order status")
	}

	// Добавляем системное сообщение
	message := &models.OrderMessage{
		OrderID:     orderID,
		SenderID:    sellerID,
		MessageType: models.OrderMessageTypeShippingUpdate,
		Content:     fmt.Sprintf("Заказ отправлен. Способ доставки: %s. Трек-номер: %s", shippingMethod, trackingNumber),
		Metadata: map[string]interface{}{
			"shipping_method": shippingMethod,
			"tracking_number": trackingNumber,
		},
	}
	s.orderRepo.AddMessage(ctx, message)

	// Уведомляем покупателя
	go s.notificationSvc.SendOrderShipped(ctx, order)

	return nil
}

// ConfirmDelivery подтверждает получение заказа
func (s *OrderService) ConfirmDelivery(ctx context.Context, orderID int64, buyerID int64) error {
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return errors.Wrap(err, "failed to get order")
	}

	// Проверки
	if order.BuyerID != buyerID {
		return errors.New("unauthorized: not the buyer")
	}
	if order.Status != models.MarketplaceOrderStatusShipped {
		return errors.New("order must be shipped before delivery confirmation")
	}

	// Обновляем статус
	err = s.orderRepo.UpdateStatus(ctx, orderID, models.MarketplaceOrderStatusDelivered, "Confirmed by buyer", &buyerID)
	if err != nil {
		return errors.Wrap(err, "failed to update order status")
	}

	// Устанавливаем дату окончания защитного периода
	order.DeliveredAt = &[]time.Time{time.Now()}[0]
	order.SetProtectionExpiry()
	// TODO: Обновить protection_expires_at в БД

	// Уведомляем продавца
	go s.notificationSvc.SendOrderDelivered(ctx, order)

	return nil
}

// CompleteOrder завершает заказ и инициирует выплату продавцу
func (s *OrderService) CompleteOrder(ctx context.Context, orderID int64) error {
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return errors.Wrap(err, "failed to get order")
	}

	// Проверяем можно ли завершить заказ
	if !order.CanBeCaptured() {
		return errors.New("order cannot be completed yet")
	}

	// Захватываем платеж
	if order.PaymentTransactionID != nil {
		err = s.paymentService.CapturePayment(ctx, *order.PaymentTransactionID)
		if err != nil {
			return errors.Wrap(err, "failed to capture payment")
		}
	}

	// Обновляем статус заказа
	err = s.orderRepo.UpdateStatus(ctx, orderID, models.MarketplaceOrderStatusCompleted, "Auto-completed after protection period", nil)
	if err != nil {
		return errors.Wrap(err, "failed to update order status")
	}

	// Инициируем выплату продавцу (будет в следующей фазе)
	// TODO: s.paymentService.CreatePayout(...)

	// Уведомляем о выплате
	go s.notificationSvc.SendPaymentReleased(ctx, order)

	return nil
}

// OpenDispute открывает спор по заказу
func (s *OrderService) OpenDispute(ctx context.Context, orderID int64, userID int64, reason string) error {
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return errors.Wrap(err, "failed to get order")
	}

	// Проверки
	if order.BuyerID != userID && order.SellerID != userID {
		return errors.New("unauthorized: not a party of this order")
	}
	if order.Status != models.MarketplaceOrderStatusShipped && order.Status != models.MarketplaceOrderStatusDelivered {
		return errors.New("dispute can only be opened for shipped or delivered orders")
	}

	// Обновляем статус
	err = s.orderRepo.UpdateStatus(ctx, orderID, models.MarketplaceOrderStatusDisputed, reason, &userID)
	if err != nil {
		return errors.Wrap(err, "failed to update order status")
	}

	// Добавляем сообщение о споре
	messageType := models.OrderMessageTypeDisputeOpened
	initiator := "покупателем"
	if userID == order.SellerID {
		initiator = "продавцом"
	}

	message := &models.OrderMessage{
		OrderID:     orderID,
		SenderID:    userID,
		MessageType: messageType,
		Content:     fmt.Sprintf("Открыт спор %s. Причина: %s", initiator, reason),
		Metadata: map[string]interface{}{
			"dispute_reason": reason,
		},
	}
	s.orderRepo.AddMessage(ctx, message)

	// TODO: Уведомить службу поддержки

	return nil
}

// ProcessAutoCaptureOrders обрабатывает заказы для автоматического capture
func (s *OrderService) ProcessAutoCaptureOrders(ctx context.Context) error {
	orders, err := s.orderRepo.GetOrdersForAutoCapture(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get orders for auto capture")
	}

	for _, order := range orders {
		// Пытаемся завершить каждый заказ
		err := s.CompleteOrder(ctx, order.ID)
		if err != nil {
			// Логируем ошибку но продолжаем с другими заказами
			// TODO: Добавить логирование
			continue
		}
	}

	return nil
}

// GetBuyerOrders получает заказы покупателя
func (s *OrderService) GetBuyerOrders(ctx context.Context, buyerID int64, page, limit int) ([]*models.MarketplaceOrder, int, error) {
	offset := (page - 1) * limit
	orders, total, err := s.orderRepo.GetBuyerOrders(ctx, buyerID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	// Загружаем связанные данные для каждого заказа
	for _, order := range orders {
		// Загружаем информацию о листинге
		if order.ListingID > 0 {
			listing, err := s.listingRepo.GetListingByID(ctx, int(order.ListingID))
			if err == nil && listing != nil {
				order.Listing = listing
			}
		}

		// Загружаем информацию о продавце
		if order.SellerID > 0 {
			seller, err := s.userRepo.GetByID(ctx, int(order.SellerID))
			if err == nil && seller != nil {
				order.Seller = seller
			}
		}

		// Загружаем информацию о покупателе
		if order.BuyerID > 0 {
			buyer, err := s.userRepo.GetByID(ctx, int(order.BuyerID))
			if err == nil && buyer != nil {
				order.Buyer = buyer
			}
		}
	}

	return orders, total, nil
}

// GetSellerOrders получает заказы продавца
func (s *OrderService) GetSellerOrders(ctx context.Context, sellerID int64, page, limit int) ([]*models.MarketplaceOrder, int, error) {
	offset := (page - 1) * limit
	orders, total, err := s.orderRepo.GetSellerOrders(ctx, sellerID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	// Загружаем связанные данные для каждого заказа
	for _, order := range orders {
		// Загружаем информацию о листинге
		if order.ListingID > 0 {
			listing, err := s.listingRepo.GetListingByID(ctx, int(order.ListingID))
			if err == nil && listing != nil {
				order.Listing = listing
			}
		}

		// Загружаем информацию о продавце
		if order.SellerID > 0 {
			seller, err := s.userRepo.GetByID(ctx, int(order.SellerID))
			if err == nil && seller != nil {
				order.Seller = seller
			}
		}

		// Загружаем информацию о покупателе
		if order.BuyerID > 0 {
			buyer, err := s.userRepo.GetByID(ctx, int(order.BuyerID))
			if err == nil && buyer != nil {
				order.Buyer = buyer
			}
		}
	}

	return orders, total, nil
}

// GetOrderDetails получает детали заказа с проверкой доступа
func (s *OrderService) GetOrderDetails(ctx context.Context, orderID int64, userID int64) (*models.MarketplaceOrder, error) {
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	// Проверка доступа
	if order.BuyerID != userID && order.SellerID != userID {
		return nil, errors.New("unauthorized: not a party of this order")
	}

	// Загружаем информацию о листинге
	if order.ListingID > 0 {
		listing, err := s.listingRepo.GetListingByID(ctx, int(order.ListingID))
		if err == nil && listing != nil {
			order.Listing = listing
		}
	}

	// Загружаем информацию о продавце
	if order.SellerID > 0 {
		seller, err := s.userRepo.GetByID(ctx, int(order.SellerID))
		if err == nil && seller != nil {
			order.Seller = seller
		}
	}

	// Загружаем информацию о покупателе
	if order.BuyerID > 0 {
		buyer, err := s.userRepo.GetByID(ctx, int(order.BuyerID))
		if err == nil && buyer != nil {
			order.Buyer = buyer
		}
	}

	// Загружаем сообщения
	messages, err := s.orderRepo.GetOrderMessages(ctx, orderID)
	if err == nil {
		order.Messages = messages
	}

	return order, nil
}

// SimpleUserRepository простая реализация UserRepositoryInterface
type SimpleUserRepository struct {
	storage storage.Storage
}

func (r *SimpleUserRepository) GetByID(ctx context.Context, userID int) (*models.User, error) {
	return r.storage.GetUserByID(ctx, userID)
}
