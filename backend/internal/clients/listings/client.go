package listings

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/rs/zerolog"
	pb "github.com/sveturs/listings/api/proto/listings/v1"
)

const (
	defaultTimeout    = 30 * time.Second
	maxRetries        = 3
	initialBackoff    = 100 * time.Millisecond
	maxBackoff        = 2 * time.Second
	backoffMultiplier = 2.0
	// Circuit breaker: открывается после 5 последовательных ошибок
	circuitBreakerThreshold = 5
	// Circuit breaker: переходит в half-open после 30 секунд
	circuitBreakerTimeout = 30 * time.Second
)

// Client представляет gRPC клиент для listings микросервиса
type Client struct {
	conn            *grpc.ClientConn
	client          pb.ListingsServiceClient
	logger          zerolog.Logger
	failureCount    int
	lastFailureTime time.Time
	isOpen          bool
}

// NewClient создает новый gRPC клиент для listings микросервиса
func NewClient(serverURL string, logger zerolog.Logger) (*Client, error) {
	logger.Info().
		Str("url", serverURL).
		Msg("Connecting to listings gRPC service")

	// grpc.NewClient заменяет deprecated grpc.DialContext
	// Соединение устанавливается лениво при первом вызове
	conn, err := grpc.NewClient(
		serverURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		logger.Error().
			Err(err).
			Str("url", serverURL).
			Msg("Failed to create listings gRPC client")
		return nil, fmt.Errorf("failed to create listings gRPC client: %w", err)
	}

	logger.Info().Msg("Successfully created listings gRPC client")

	return &Client{
		conn:   conn,
		client: pb.NewListingsServiceClient(conn),
		logger: logger.With().Str("component", "listings_grpc_client").Logger(),
	}, nil
}

// Close закрывает соединение с gRPC сервером
func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// GetListing получает listing по ID
func (c *Client) GetListing(ctx context.Context, req *pb.GetListingRequest) (*pb.GetListingResponse, error) {
	if c.isCircuitBreakerOpen() {
		c.logger.Warn().Msg("Circuit breaker is open, rejecting GetListing request")
		return nil, ErrServiceUnavailable
	}

	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var resp *pb.GetListingResponse
	var err error

	backoff := initialBackoff
	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			c.logger.Info().Int("attempt", attempt+1).Msg("Retrying GetListing")
			time.Sleep(backoff)
			backoff = time.Duration(float64(backoff) * backoffMultiplier)
			if backoff > maxBackoff {
				backoff = maxBackoff
			}
		}

		resp, err = c.client.GetListing(ctx, req)
		if err == nil {
			c.onSuccess()
			c.logger.Debug().Int64("listing_id", req.Id).Msg("GetListing successful")
			return resp, nil
		}

		if !c.shouldRetry(err) {
			c.onFailure()
			c.logger.Error().Err(err).Msg("Non-retryable error in GetListing")
			return nil, MapGRPCError(err)
		}

		c.logger.Warn().Err(err).Int("attempt", attempt+1).Msg("Retryable error in GetListing")
	}

	c.onFailure()
	c.logger.Error().Err(err).Msg("GetListing failed after all retries")
	return nil, fmt.Errorf("failed to get listing after %d attempts: %w", maxRetries, MapGRPCError(err))
}

// CreateListing создает новый listing
func (c *Client) CreateListing(ctx context.Context, req *pb.CreateListingRequest) (*pb.CreateListingResponse, error) {
	if c.isCircuitBreakerOpen() {
		c.logger.Warn().Msg("Circuit breaker is open, rejecting CreateListing request")
		return nil, ErrServiceUnavailable
	}

	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var resp *pb.CreateListingResponse
	var err error

	backoff := initialBackoff
	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			c.logger.Info().Int("attempt", attempt+1).Msg("Retrying CreateListing")
			time.Sleep(backoff)
			backoff = time.Duration(float64(backoff) * backoffMultiplier)
			if backoff > maxBackoff {
				backoff = maxBackoff
			}
		}

		resp, err = c.client.CreateListing(ctx, req)
		if err == nil {
			c.onSuccess()
			c.logger.Info().
				Int64("listing_id", resp.Listing.Id).
				Str("title", req.Title).
				Msg("Listing created successfully")
			return resp, nil
		}

		if !c.shouldRetry(err) {
			c.onFailure()
			c.logger.Error().Err(err).Msg("Non-retryable error in CreateListing")
			return nil, MapGRPCError(err)
		}

		c.logger.Warn().Err(err).Int("attempt", attempt+1).Msg("Retryable error in CreateListing")
	}

	c.onFailure()
	c.logger.Error().Err(err).Msg("CreateListing failed after all retries")
	return nil, fmt.Errorf("failed to create listing after %d attempts: %w", maxRetries, MapGRPCError(err))
}

// UpdateListing обновляет существующий listing
func (c *Client) UpdateListing(ctx context.Context, req *pb.UpdateListingRequest) (*pb.UpdateListingResponse, error) {
	if c.isCircuitBreakerOpen() {
		c.logger.Warn().Msg("Circuit breaker is open, rejecting UpdateListing request")
		return nil, ErrServiceUnavailable
	}

	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var resp *pb.UpdateListingResponse
	var err error

	backoff := initialBackoff
	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			c.logger.Info().Int("attempt", attempt+1).Msg("Retrying UpdateListing")
			time.Sleep(backoff)
			backoff = time.Duration(float64(backoff) * backoffMultiplier)
			if backoff > maxBackoff {
				backoff = maxBackoff
			}
		}

		resp, err = c.client.UpdateListing(ctx, req)
		if err == nil {
			c.onSuccess()
			c.logger.Info().Int64("listing_id", req.Id).Msg("Listing updated successfully")
			return resp, nil
		}

		if !c.shouldRetry(err) {
			c.onFailure()
			c.logger.Error().Err(err).Msg("Non-retryable error in UpdateListing")
			return nil, MapGRPCError(err)
		}

		c.logger.Warn().Err(err).Int("attempt", attempt+1).Msg("Retryable error in UpdateListing")
	}

	c.onFailure()
	c.logger.Error().Err(err).Msg("UpdateListing failed after all retries")
	return nil, fmt.Errorf("failed to update listing after %d attempts: %w", maxRetries, MapGRPCError(err))
}

// DeleteListing удаляет listing (soft delete)
func (c *Client) DeleteListing(ctx context.Context, req *pb.DeleteListingRequest) (*pb.DeleteListingResponse, error) {
	if c.isCircuitBreakerOpen() {
		c.logger.Warn().Msg("Circuit breaker is open, rejecting DeleteListing request")
		return nil, ErrServiceUnavailable
	}

	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var resp *pb.DeleteListingResponse
	var err error

	backoff := initialBackoff
	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			c.logger.Info().Int("attempt", attempt+1).Msg("Retrying DeleteListing")
			time.Sleep(backoff)
			backoff = time.Duration(float64(backoff) * backoffMultiplier)
			if backoff > maxBackoff {
				backoff = maxBackoff
			}
		}

		resp, err = c.client.DeleteListing(ctx, req)
		if err == nil {
			c.onSuccess()
			c.logger.Info().Int64("listing_id", req.Id).Msg("Listing deleted successfully")
			return resp, nil
		}

		if !c.shouldRetry(err) {
			c.onFailure()
			c.logger.Error().Err(err).Msg("Non-retryable error in DeleteListing")
			return nil, MapGRPCError(err)
		}

		c.logger.Warn().Err(err).Int("attempt", attempt+1).Msg("Retryable error in DeleteListing")
	}

	c.onFailure()
	c.logger.Error().Err(err).Msg("DeleteListing failed after all retries")
	return nil, fmt.Errorf("failed to delete listing after %d attempts: %w", maxRetries, MapGRPCError(err))
}

// SearchListings выполняет поиск по listings
func (c *Client) SearchListings(ctx context.Context, req *pb.SearchListingsRequest) (*pb.SearchListingsResponse, error) {
	if c.isCircuitBreakerOpen() {
		c.logger.Warn().Msg("Circuit breaker is open, rejecting SearchListings request")
		return nil, ErrServiceUnavailable
	}

	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var resp *pb.SearchListingsResponse
	var err error

	backoff := initialBackoff
	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			c.logger.Info().Int("attempt", attempt+1).Msg("Retrying SearchListings")
			time.Sleep(backoff)
			backoff = time.Duration(float64(backoff) * backoffMultiplier)
			if backoff > maxBackoff {
				backoff = maxBackoff
			}
		}

		resp, err = c.client.SearchListings(ctx, req)
		if err == nil {
			c.onSuccess()
			c.logger.Debug().
				Int("count", len(resp.Listings)).
				Str("query", req.Query).
				Msg("SearchListings successful")
			return resp, nil
		}

		if !c.shouldRetry(err) {
			c.onFailure()
			c.logger.Error().Err(err).Msg("Non-retryable error in SearchListings")
			return nil, MapGRPCError(err)
		}

		c.logger.Warn().Err(err).Int("attempt", attempt+1).Msg("Retryable error in SearchListings")
	}

	c.onFailure()
	c.logger.Error().Err(err).Msg("SearchListings failed after all retries")
	return nil, fmt.Errorf("failed to search listings after %d attempts: %w", maxRetries, MapGRPCError(err))
}

// ListListings получает список listings с фильтрацией
func (c *Client) ListListings(ctx context.Context, req *pb.ListListingsRequest) (*pb.ListListingsResponse, error) {
	if c.isCircuitBreakerOpen() {
		c.logger.Warn().Msg("Circuit breaker is open, rejecting ListListings request")
		return nil, ErrServiceUnavailable
	}

	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var resp *pb.ListListingsResponse
	var err error

	backoff := initialBackoff
	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			c.logger.Info().Int("attempt", attempt+1).Msg("Retrying ListListings")
			time.Sleep(backoff)
			backoff = time.Duration(float64(backoff) * backoffMultiplier)
			if backoff > maxBackoff {
				backoff = maxBackoff
			}
		}

		resp, err = c.client.ListListings(ctx, req)
		if err == nil {
			c.onSuccess()
			c.logger.Debug().Int("count", len(resp.Listings)).Msg("ListListings successful")
			return resp, nil
		}

		if !c.shouldRetry(err) {
			c.onFailure()
			c.logger.Error().Err(err).Msg("Non-retryable error in ListListings")
			return nil, MapGRPCError(err)
		}

		c.logger.Warn().Err(err).Int("attempt", attempt+1).Msg("Retryable error in ListListings")
	}

	c.onFailure()
	c.logger.Error().Err(err).Msg("ListListings failed after all retries")
	return nil, fmt.Errorf("failed to list listings after %d attempts: %w", maxRetries, MapGRPCError(err))
}

// GetAllCategories получает все категории из microservice
func (c *Client) GetAllCategories(ctx context.Context) (*pb.CategoriesResponse, error) {
	if c.isCircuitBreakerOpen() {
		c.logger.Warn().Msg("Circuit breaker is open, rejecting GetAllCategories request")
		return nil, ErrServiceUnavailable
	}

	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var resp *pb.CategoriesResponse
	var err error

	backoff := initialBackoff
	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			c.logger.Info().Int("attempt", attempt+1).Msg("Retrying GetAllCategories")
			time.Sleep(backoff)
			backoff = time.Duration(float64(backoff) * backoffMultiplier)
			if backoff > maxBackoff {
				backoff = maxBackoff
			}
		}

		// GetAllCategories doesn't take any parameters (uses google.protobuf.Empty)
		resp, err = c.client.GetAllCategories(ctx, &emptypb.Empty{})
		if err == nil {
			c.onSuccess()
			c.logger.Debug().Int("count", len(resp.Categories)).Msg("GetAllCategories successful")
			return resp, nil
		}

		if !c.shouldRetry(err) {
			c.onFailure()
			c.logger.Error().Err(err).Msg("Non-retryable error in GetAllCategories")
			return nil, MapGRPCError(err)
		}

		c.logger.Warn().Err(err).Int("attempt", attempt+1).Msg("Retryable error in GetAllCategories")
	}

	c.onFailure()
	c.logger.Error().Err(err).Msg("GetAllCategories failed after all retries")
	return nil, fmt.Errorf("failed to get all categories after %d attempts: %w", maxRetries, MapGRPCError(err))
}

// AddToFavorites добавляет listing в избранное пользователя
func (c *Client) AddToFavorites(ctx context.Context, userID, listingID int) error {
	if c.isCircuitBreakerOpen() {
		c.logger.Warn().Msg("Circuit breaker is open, rejecting AddToFavorites request")
		return ErrServiceUnavailable
	}

	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	req := &pb.AddToFavoritesRequest{
		UserId:    int64(userID),
		ListingId: int64(listingID),
	}

	var err error
	backoff := initialBackoff

	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			c.logger.Info().Int("attempt", attempt+1).Msg("Retrying AddToFavorites")
			time.Sleep(backoff)
			backoff = time.Duration(float64(backoff) * backoffMultiplier)
			if backoff > maxBackoff {
				backoff = maxBackoff
			}
		}

		_, err = c.client.AddToFavorites(ctx, req)
		if err == nil {
			c.onSuccess()
			c.logger.Debug().
				Int("user_id", userID).
				Int("listing_id", listingID).
				Msg("AddToFavorites successful")
			return nil
		}

		if !c.shouldRetry(err) {
			c.onFailure()
			c.logger.Error().Err(err).Msg("Non-retryable error in AddToFavorites")
			return MapGRPCError(err)
		}

		c.logger.Warn().Err(err).Int("attempt", attempt+1).Msg("Retryable error in AddToFavorites")
	}

	c.onFailure()
	c.logger.Error().Err(err).Msg("AddToFavorites failed after all retries")
	return fmt.Errorf("failed to add to favorites after %d attempts: %w", maxRetries, MapGRPCError(err))
}

// RemoveFromFavorites удаляет listing из избранного пользователя
func (c *Client) RemoveFromFavorites(ctx context.Context, userID, listingID int) error {
	if c.isCircuitBreakerOpen() {
		c.logger.Warn().Msg("Circuit breaker is open, rejecting RemoveFromFavorites request")
		return ErrServiceUnavailable
	}

	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	req := &pb.RemoveFromFavoritesRequest{
		UserId:    int64(userID),
		ListingId: int64(listingID),
	}

	var err error
	backoff := initialBackoff

	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			c.logger.Info().Int("attempt", attempt+1).Msg("Retrying RemoveFromFavorites")
			time.Sleep(backoff)
			backoff = time.Duration(float64(backoff) * backoffMultiplier)
			if backoff > maxBackoff {
				backoff = maxBackoff
			}
		}

		_, err = c.client.RemoveFromFavorites(ctx, req)
		if err == nil {
			c.onSuccess()
			c.logger.Debug().
				Int("user_id", userID).
				Int("listing_id", listingID).
				Msg("RemoveFromFavorites successful")
			return nil
		}

		if !c.shouldRetry(err) {
			c.onFailure()
			c.logger.Error().Err(err).Msg("Non-retryable error in RemoveFromFavorites")
			return MapGRPCError(err)
		}

		c.logger.Warn().Err(err).Int("attempt", attempt+1).Msg("Retryable error in RemoveFromFavorites")
	}

	c.onFailure()
	c.logger.Error().Err(err).Msg("RemoveFromFavorites failed after all retries")
	return fmt.Errorf("failed to remove from favorites after %d attempts: %w", maxRetries, MapGRPCError(err))
}

// GetUserFavorites получает список ID избранных listings пользователя
func (c *Client) GetUserFavorites(ctx context.Context, userID int) ([]int, error) {
	if c.isCircuitBreakerOpen() {
		c.logger.Warn().Msg("Circuit breaker is open, rejecting GetUserFavorites request")
		return nil, ErrServiceUnavailable
	}

	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	req := &pb.GetUserFavoritesRequest{
		UserId: int64(userID),
		Limit:  0, // 0 = без лимита
		Offset: 0,
	}

	var resp *pb.GetUserFavoritesResponse
	var err error
	backoff := initialBackoff

	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			c.logger.Info().Int("attempt", attempt+1).Msg("Retrying GetUserFavorites")
			time.Sleep(backoff)
			backoff = time.Duration(float64(backoff) * backoffMultiplier)
			if backoff > maxBackoff {
				backoff = maxBackoff
			}
		}

		resp, err = c.client.GetUserFavorites(ctx, req)
		if err == nil {
			c.onSuccess()

			// Конвертируем int64 в int
			listingIDs := make([]int, len(resp.ListingIds))
			for i, id := range resp.ListingIds {
				listingIDs[i] = int(id)
			}

			c.logger.Debug().
				Int("user_id", userID).
				Int("count", len(listingIDs)).
				Msg("GetUserFavorites successful")
			return listingIDs, nil
		}

		if !c.shouldRetry(err) {
			c.onFailure()
			c.logger.Error().Err(err).Msg("Non-retryable error in GetUserFavorites")
			return nil, MapGRPCError(err)
		}

		c.logger.Warn().Err(err).Int("attempt", attempt+1).Msg("Retryable error in GetUserFavorites")
	}

	c.onFailure()
	c.logger.Error().Err(err).Msg("GetUserFavorites failed after all retries")
	return nil, fmt.Errorf("failed to get user favorites after %d attempts: %w", maxRetries, MapGRPCError(err))
}

// IsFavorite проверяет, находится ли listing в избранном пользователя
func (c *Client) IsFavorite(ctx context.Context, userID, listingID int) (bool, error) {
	if c.isCircuitBreakerOpen() {
		c.logger.Warn().Msg("Circuit breaker is open, rejecting IsFavorite request")
		return false, ErrServiceUnavailable
	}

	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	req := &pb.IsFavoriteRequest{
		UserId:    int64(userID),
		ListingId: int64(listingID),
	}

	var resp *pb.IsFavoriteResponse
	var err error
	backoff := initialBackoff

	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			c.logger.Info().Int("attempt", attempt+1).Msg("Retrying IsFavorite")
			time.Sleep(backoff)
			backoff = time.Duration(float64(backoff) * backoffMultiplier)
			if backoff > maxBackoff {
				backoff = maxBackoff
			}
		}

		resp, err = c.client.IsFavorite(ctx, req)
		if err == nil {
			c.onSuccess()
			c.logger.Debug().
				Int("user_id", userID).
				Int("listing_id", listingID).
				Bool("is_favorite", resp.IsFavorite).
				Msg("IsFavorite successful")
			return resp.IsFavorite, nil
		}

		if !c.shouldRetry(err) {
			c.onFailure()
			c.logger.Error().Err(err).Msg("Non-retryable error in IsFavorite")
			return false, MapGRPCError(err)
		}

		c.logger.Warn().Err(err).Int("attempt", attempt+1).Msg("Retryable error in IsFavorite")
	}

	c.onFailure()
	c.logger.Error().Err(err).Msg("IsFavorite failed after all retries")
	return false, fmt.Errorf("failed to check favorite after %d attempts: %w", maxRetries, MapGRPCError(err))
}

// shouldRetry определяет, нужно ли повторять запрос при данной ошибке
func (c *Client) shouldRetry(err error) bool {
	st, ok := status.FromError(err)
	if !ok {
		// Неизвестная ошибка - пробуем повторить
		return true
	}

	code := st.Code()
	switch code {
	// Non-retryable errors - ошибки валидации и бизнес-логики
	case codes.InvalidArgument, codes.NotFound, codes.AlreadyExists,
		codes.PermissionDenied, codes.Unauthenticated, codes.FailedPrecondition,
		codes.OutOfRange, codes.Unimplemented:
		return false
	// Retryable errors - временные сетевые проблемы
	case codes.Unavailable, codes.DeadlineExceeded, codes.ResourceExhausted,
		codes.Aborted, codes.Canceled:
		return true
	// Success and other codes - не повторяем
	case codes.OK, codes.Unknown, codes.Internal, codes.DataLoss:
		return false
	default:
		return true
	}
}

// isCircuitBreakerOpen проверяет, открыт ли circuit breaker
func (c *Client) isCircuitBreakerOpen() bool {
	if !c.isOpen {
		return false
	}

	// Переход в half-open state после таймаута
	if time.Since(c.lastFailureTime) > circuitBreakerTimeout {
		c.logger.Info().Msg("Circuit breaker transitioning to half-open state")
		c.isOpen = false
		c.failureCount = 0
		return false
	}

	return true
}

// onSuccess обрабатывает успешный ответ
func (c *Client) onSuccess() {
	if c.failureCount > 0 || c.isOpen {
		c.logger.Info().Msg("Circuit breaker closed after successful request")
	}
	c.failureCount = 0
	c.isOpen = false
}

// onFailure обрабатывает ошибку
func (c *Client) onFailure() {
	c.failureCount++
	c.lastFailureTime = time.Now()

	if c.failureCount >= circuitBreakerThreshold {
		c.isOpen = true
		c.logger.Warn().
			Int("failure_count", c.failureCount).
			Msg("Circuit breaker opened due to consecutive failures")
	}
}

// ========== Products API Methods ==========

// GetProduct retrieves a single product by ID
func (c *Client) GetProduct(ctx context.Context, productID int64, storefrontID *int64) (*pb.Product, error) {
	if c.isCircuitBreakerOpen() {
		c.logger.Warn().Msg("Circuit breaker is open, rejecting GetProduct request")
		return nil, ErrServiceUnavailable
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	req := &pb.GetProductRequest{
		ProductId: productID,
	}
	if storefrontID != nil {
		req.StorefrontId = storefrontID
	}

	var resp *pb.ProductResponse
	var err error
	backoff := initialBackoff

	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			c.logger.Info().Int("attempt", attempt+1).Msg("Retrying GetProduct")
			time.Sleep(backoff)
			backoff = time.Duration(float64(backoff) * backoffMultiplier)
			if backoff > maxBackoff {
				backoff = maxBackoff
			}
		}

		resp, err = c.client.GetProduct(ctx, req)
		if err == nil {
			c.onSuccess()
			c.logger.Debug().Int64("product_id", productID).Msg("GetProduct successful")
			return resp.Product, nil
		}

		if !c.shouldRetry(err) {
			c.onFailure()
			c.logger.Error().Err(err).Msg("Non-retryable error in GetProduct")
			return nil, MapGRPCError(err)
		}

		c.logger.Warn().Err(err).Int("attempt", attempt+1).Msg("Retryable error in GetProduct")
	}

	c.onFailure()
	c.logger.Error().Err(err).Msg("GetProduct failed after all retries")
	return nil, fmt.Errorf("failed to get product after %d attempts: %w", maxRetries, MapGRPCError(err))
}

// GetProductsBySKUs retrieves products by list of SKUs (для корзины)
func (c *Client) GetProductsBySKUs(ctx context.Context, skus []string, storefrontID *int64) ([]*pb.Product, error) {
	if c.isCircuitBreakerOpen() {
		c.logger.Warn().Msg("Circuit breaker is open, rejecting GetProductsBySKUs request")
		return nil, ErrServiceUnavailable
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second) // Больший timeout для batch операции
	defer cancel()

	req := &pb.GetProductsBySKUsRequest{
		Skus: skus,
	}
	if storefrontID != nil {
		req.StorefrontId = storefrontID
	}

	var resp *pb.ProductsResponse
	var err error
	backoff := initialBackoff

	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			c.logger.Info().Int("attempt", attempt+1).Msg("Retrying GetProductsBySKUs")
			time.Sleep(backoff)
			backoff = time.Duration(float64(backoff) * backoffMultiplier)
			if backoff > maxBackoff {
				backoff = maxBackoff
			}
		}

		resp, err = c.client.GetProductsBySKUs(ctx, req)
		if err == nil {
			c.onSuccess()
			c.logger.Debug().
				Int("count", len(resp.Products)).
				Int("requested_skus", len(skus)).
				Msg("GetProductsBySKUs successful")
			return resp.Products, nil
		}

		if !c.shouldRetry(err) {
			c.onFailure()
			c.logger.Error().Err(err).Msg("Non-retryable error in GetProductsBySKUs")
			return nil, MapGRPCError(err)
		}

		c.logger.Warn().Err(err).Int("attempt", attempt+1).Msg("Retryable error in GetProductsBySKUs")
	}

	c.onFailure()
	c.logger.Error().Err(err).Msg("GetProductsBySKUs failed after all retries")
	return nil, fmt.Errorf("failed to get products by SKUs after %d attempts: %w", maxRetries, MapGRPCError(err))
}

// GetProductsByIDs retrieves products by list of IDs (для корзины)
func (c *Client) GetProductsByIDs(ctx context.Context, productIDs []int64, storefrontID *int64) ([]*pb.Product, error) {
	if c.isCircuitBreakerOpen() {
		c.logger.Warn().Msg("Circuit breaker is open, rejecting GetProductsByIDs request")
		return nil, ErrServiceUnavailable
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second) // Больший timeout для batch операции
	defer cancel()

	req := &pb.GetProductsByIDsRequest{
		ProductIds: productIDs,
	}
	if storefrontID != nil {
		req.StorefrontId = storefrontID
	}

	var resp *pb.ProductsResponse
	var err error
	backoff := initialBackoff

	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			c.logger.Info().Int("attempt", attempt+1).Msg("Retrying GetProductsByIDs")
			time.Sleep(backoff)
			backoff = time.Duration(float64(backoff) * backoffMultiplier)
			if backoff > maxBackoff {
				backoff = maxBackoff
			}
		}

		resp, err = c.client.GetProductsByIDs(ctx, req)
		if err == nil {
			c.onSuccess()
			c.logger.Debug().
				Int("count", len(resp.Products)).
				Int("requested_ids", len(productIDs)).
				Msg("GetProductsByIDs successful")
			return resp.Products, nil
		}

		if !c.shouldRetry(err) {
			c.onFailure()
			c.logger.Error().Err(err).Msg("Non-retryable error in GetProductsByIDs")
			return nil, MapGRPCError(err)
		}

		c.logger.Warn().Err(err).Int("attempt", attempt+1).Msg("Retryable error in GetProductsByIDs")
	}

	c.onFailure()
	c.logger.Error().Err(err).Msg("GetProductsByIDs failed after all retries")
	return nil, fmt.Errorf("failed to get products by IDs after %d attempts: %w", maxRetries, MapGRPCError(err))
}

// ListProducts retrieves paginated list of products
func (c *Client) ListProducts(ctx context.Context, storefrontID int64, page, pageSize int, isActiveOnly bool) ([]*pb.Product, int32, error) {
	if c.isCircuitBreakerOpen() {
		c.logger.Warn().Msg("Circuit breaker is open, rejecting ListProducts request")
		return nil, 0, ErrServiceUnavailable
	}

	// Validate pagination parameters
	if page < 1 || page > 100000 {
		return nil, 0, fmt.Errorf("invalid page parameter: %d", page)
	}
	if pageSize < 1 || pageSize > 1000 {
		return nil, 0, fmt.Errorf("invalid pageSize parameter: %d", pageSize)
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	isActiveOnlyPtr := &isActiveOnly
	req := &pb.ListProductsRequest{
		StorefrontId: storefrontID,
		Page:         int32(page),     // #nosec G115 - validated above
		PageSize:     int32(pageSize), // #nosec G115 - validated above
		IsActiveOnly: isActiveOnlyPtr,
	}

	var resp *pb.ProductsResponse
	var err error
	backoff := initialBackoff

	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			c.logger.Info().Int("attempt", attempt+1).Msg("Retrying ListProducts")
			time.Sleep(backoff)
			backoff = time.Duration(float64(backoff) * backoffMultiplier)
			if backoff > maxBackoff {
				backoff = maxBackoff
			}
		}

		resp, err = c.client.ListProducts(ctx, req)
		if err == nil {
			c.onSuccess()
			c.logger.Debug().
				Int64("storefront_id", storefrontID).
				Int("count", len(resp.Products)).
				Int32("total", resp.TotalCount).
				Msg("ListProducts successful")
			return resp.Products, resp.TotalCount, nil
		}

		if !c.shouldRetry(err) {
			c.onFailure()
			c.logger.Error().Err(err).Msg("Non-retryable error in ListProducts")
			return nil, 0, MapGRPCError(err)
		}

		c.logger.Warn().Err(err).Int("attempt", attempt+1).Msg("Retryable error in ListProducts")
	}

	c.onFailure()
	c.logger.Error().Err(err).Msg("ListProducts failed after all retries")
	return nil, 0, fmt.Errorf("failed to list products after %d attempts: %w", maxRetries, MapGRPCError(err))
}

// GetVariant retrieves a single variant by ID
func (c *Client) GetVariant(ctx context.Context, variantID int64, productID *int64) (*pb.ProductVariant, error) {
	if c.isCircuitBreakerOpen() {
		c.logger.Warn().Msg("Circuit breaker is open, rejecting GetVariant request")
		return nil, ErrServiceUnavailable
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	req := &pb.GetVariantRequest{
		VariantId: variantID,
	}
	if productID != nil {
		req.ProductId = productID
	}

	var resp *pb.VariantResponse
	var err error
	backoff := initialBackoff

	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			c.logger.Info().Int("attempt", attempt+1).Msg("Retrying GetVariant")
			time.Sleep(backoff)
			backoff = time.Duration(float64(backoff) * backoffMultiplier)
			if backoff > maxBackoff {
				backoff = maxBackoff
			}
		}

		resp, err = c.client.GetVariant(ctx, req)
		if err == nil {
			c.onSuccess()
			c.logger.Debug().Int64("variant_id", variantID).Msg("GetVariant successful")
			return resp.Variant, nil
		}

		if !c.shouldRetry(err) {
			c.onFailure()
			c.logger.Error().Err(err).Msg("Non-retryable error in GetVariant")
			return nil, MapGRPCError(err)
		}

		c.logger.Warn().Err(err).Int("attempt", attempt+1).Msg("Retryable error in GetVariant")
	}

	c.onFailure()
	c.logger.Error().Err(err).Msg("GetVariant failed after all retries")
	return nil, fmt.Errorf("failed to get variant after %d attempts: %w", maxRetries, MapGRPCError(err))
}

// GetVariantsByProductID retrieves all variants for a product
func (c *Client) GetVariantsByProductID(ctx context.Context, productID int64, isActiveOnly bool) ([]*pb.ProductVariant, error) {
	if c.isCircuitBreakerOpen() {
		c.logger.Warn().Msg("Circuit breaker is open, rejecting GetVariantsByProductID request")
		return nil, ErrServiceUnavailable
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	isActiveOnlyPtr := &isActiveOnly
	req := &pb.GetVariantsByProductIDRequest{
		ProductId:    productID,
		IsActiveOnly: isActiveOnlyPtr,
	}

	var resp *pb.ProductVariantsResponse
	var err error
	backoff := initialBackoff

	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			c.logger.Info().Int("attempt", attempt+1).Msg("Retrying GetVariantsByProductID")
			time.Sleep(backoff)
			backoff = time.Duration(float64(backoff) * backoffMultiplier)
			if backoff > maxBackoff {
				backoff = maxBackoff
			}
		}

		resp, err = c.client.GetVariantsByProductID(ctx, req)
		if err == nil {
			c.onSuccess()
			c.logger.Debug().
				Int64("product_id", productID).
				Int("count", len(resp.Variants)).
				Msg("GetVariantsByProductID successful")
			return resp.Variants, nil
		}

		if !c.shouldRetry(err) {
			c.onFailure()
			c.logger.Error().Err(err).Msg("Non-retryable error in GetVariantsByProductID")
			return nil, MapGRPCError(err)
		}

		c.logger.Warn().Err(err).Int("attempt", attempt+1).Msg("Retryable error in GetVariantsByProductID")
	}

	c.onFailure()
	c.logger.Error().Err(err).Msg("GetVariantsByProductID failed after all retries")
	return nil, fmt.Errorf("failed to get variants by product ID after %d attempts: %w", maxRetries, MapGRPCError(err))
}

// ========================================
// Stock Management Methods (for Orders Service)
// ========================================

// CheckStockAvailability проверяет доступность товаров перед созданием заказа
func (c *Client) CheckStockAvailability(ctx context.Context, items []*pb.StockItem) (*pb.CheckStockAvailabilityResponse, error) {
	if c.isCircuitBreakerOpen() {
		c.logger.Warn().Msg("Circuit breaker is open, rejecting CheckStockAvailability request")
		return nil, ErrServiceUnavailable
	}

	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	req := &pb.CheckStockAvailabilityRequest{
		Items: items,
	}

	var resp *pb.CheckStockAvailabilityResponse
	var err error
	backoff := initialBackoff

	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			c.logger.Info().Int("attempt", attempt+1).Msg("Retrying CheckStockAvailability")
			time.Sleep(backoff)
			backoff = time.Duration(float64(backoff) * backoffMultiplier)
			if backoff > maxBackoff {
				backoff = maxBackoff
			}
		}

		resp, err = c.client.CheckStockAvailability(ctx, req)
		if err == nil {
			c.onSuccess()
			c.logger.Debug().
				Int("items_count", len(items)).
				Bool("all_available", resp.AllAvailable).
				Msg("CheckStockAvailability successful")
			return resp, nil
		}

		if !c.shouldRetry(err) {
			c.onFailure()
			c.logger.Error().Err(err).Msg("Non-retryable error in CheckStockAvailability")
			return nil, MapGRPCError(err)
		}

		c.logger.Warn().Err(err).Int("attempt", attempt+1).Msg("Retryable error in CheckStockAvailability")
	}

	c.onFailure()
	c.logger.Error().Err(err).Msg("CheckStockAvailability failed after all retries")
	return nil, fmt.Errorf("failed to check stock availability after %d attempts: %w", maxRetries, MapGRPCError(err))
}

// DecrementStock атомарно уменьшает количество товаров при создании заказа
func (c *Client) DecrementStock(ctx context.Context, items []*pb.StockItem, orderID string) (*pb.DecrementStockResponse, error) {
	if c.isCircuitBreakerOpen() {
		c.logger.Warn().Msg("Circuit breaker is open, rejecting DecrementStock request")
		return nil, ErrServiceUnavailable
	}

	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	req := &pb.DecrementStockRequest{
		Items:   items,
		OrderId: &orderID,
	}

	var resp *pb.DecrementStockResponse
	var err error
	backoff := initialBackoff

	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			c.logger.Info().Int("attempt", attempt+1).Msg("Retrying DecrementStock")
			time.Sleep(backoff)
			backoff = time.Duration(float64(backoff) * backoffMultiplier)
			if backoff > maxBackoff {
				backoff = maxBackoff
			}
		}

		resp, err = c.client.DecrementStock(ctx, req)
		if err == nil {
			c.onSuccess()
			c.logger.Info().
				Str("order_id", orderID).
				Int("items_count", len(items)).
				Bool("success", resp.Success).
				Msg("DecrementStock successful")
			return resp, nil
		}

		if !c.shouldRetry(err) {
			c.onFailure()
			c.logger.Error().Err(err).Msg("Non-retryable error in DecrementStock")
			return nil, MapGRPCError(err)
		}

		c.logger.Warn().Err(err).Int("attempt", attempt+1).Msg("Retryable error in DecrementStock")
	}

	c.onFailure()
	c.logger.Error().Err(err).Msg("DecrementStock failed after all retries")
	return nil, fmt.Errorf("failed to decrement stock after %d attempts: %w", maxRetries, MapGRPCError(err))
}

// RollbackStock восстанавливает количество товаров при откате заказа (compensating transaction)
func (c *Client) RollbackStock(ctx context.Context, items []*pb.StockItem, orderID string) error {
	if c.isCircuitBreakerOpen() {
		c.logger.Warn().Msg("Circuit breaker is open, rejecting RollbackStock request")
		return ErrServiceUnavailable
	}

	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	req := &pb.RollbackStockRequest{
		Items:   items,
		OrderId: &orderID,
	}

	var resp *pb.RollbackStockResponse
	var err error
	backoff := initialBackoff

	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			c.logger.Info().Int("attempt", attempt+1).Msg("Retrying RollbackStock")
			time.Sleep(backoff)
			backoff = time.Duration(float64(backoff) * backoffMultiplier)
			if backoff > maxBackoff {
				backoff = maxBackoff
			}
		}

		resp, err = c.client.RollbackStock(ctx, req)
		if err == nil {
			c.onSuccess()
			c.logger.Info().
				Str("order_id", orderID).
				Int("items_count", len(items)).
				Bool("success", resp.Success).
				Msg("RollbackStock successful")
			return nil
		}

		if !c.shouldRetry(err) {
			c.onFailure()
			c.logger.Error().Err(err).Msg("Non-retryable error in RollbackStock")
			return MapGRPCError(err)
		}

		c.logger.Warn().Err(err).Int("attempt", attempt+1).Msg("Retryable error in RollbackStock")
	}

	c.onFailure()
	c.logger.Error().Err(err).Msg("RollbackStock failed after all retries")
	return fmt.Errorf("failed to rollback stock after %d attempts: %w", maxRetries, MapGRPCError(err))
}
