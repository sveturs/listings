// Package service provides business logic layer for the listings microservice.
package service

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/rs/zerolog"

	"github.com/vondi-global/listings/internal/domain"
	"github.com/vondi-global/listings/internal/repository/postgres"
)

// CartService defines business logic operations for shopping cart management
type CartService interface {
	// Cart operations
	AddToCart(ctx context.Context, req *AddToCartRequest) (*domain.Cart, error)
	UpdateCartItem(ctx context.Context, req *UpdateCartItemRequest) (*domain.Cart, error)
	UpdateCartItemByItemID(ctx context.Context, cartItemID int64, quantity int32, userID *int64, sessionID *string) (*domain.Cart, error)
	RemoveFromCart(ctx context.Context, cartID, itemID int64) error
	RemoveFromCartByItemID(ctx context.Context, cartItemID int64, userID *int64, sessionID *string) error
	GetCart(ctx context.Context, userID *int64, sessionID *string, storefrontID int64) (*domain.Cart, error)
	ClearCart(ctx context.Context, cartID int64) error
	GetUserCarts(ctx context.Context, userID int64) ([]*domain.Cart, error)

	// Helper methods
	MergeSessionCartToUser(ctx context.Context, sessionID string, userID int64) error
	RecalculateCart(ctx context.Context, cartID int64) (*domain.Cart, error)
	ValidateCartItems(ctx context.Context, cartID int64) ([]PriceChangeItem, error)
}

// AddToCartRequest contains parameters for adding items to cart
type AddToCartRequest struct {
	UserID       *int64  // NULL for anonymous cart
	SessionID    *string // NULL for authenticated cart
	StorefrontID int64   // Required
	ListingID    int64   // Product/Listing ID
	VariantID    *int64  // Optional variant ID
	Quantity     int32   // Quantity to add
}

// UpdateCartItemRequest contains parameters for updating cart item quantity
type UpdateCartItemRequest struct {
	CartID   int64
	ItemID   int64
	Quantity int32
}

// cartService implements CartService
type cartService struct {
	cartRepo       postgres.CartRepository
	productsRepo   *postgres.Repository
	storefrontRepo *postgres.Repository
	db             *postgres.Repository
	logger         zerolog.Logger
}

// NewCartService creates a new cart service
func NewCartService(
	cartRepo postgres.CartRepository,
	productsRepo *postgres.Repository,
	storefrontRepo *postgres.Repository,
	db *postgres.Repository,
	logger zerolog.Logger,
) CartService {
	return &cartService{
		cartRepo:       cartRepo,
		productsRepo:   productsRepo,
		storefrontRepo: storefrontRepo,
		db:             db,
		logger:         logger.With().Str("component", "cart_service").Logger(),
	}
}

// AddToCart adds an item to the shopping cart
func (s *cartService) AddToCart(ctx context.Context, req *AddToCartRequest) (*domain.Cart, error) {
	s.logger.Debug().
		Interface("user_id", req.UserID).
		Interface("session_id", req.SessionID).
		Int64("storefront_id", req.StorefrontID).
		Int64("listing_id", req.ListingID).
		Int32("quantity", req.Quantity).
		Msg("adding to cart")

	// Validate input
	if err := s.validateAddToCartRequest(req); err != nil {
		return nil, err
	}

	// Fetch listing/product to validate it exists and get price
	listing, err := s.productsRepo.GetProductByID(ctx, req.ListingID, &req.StorefrontID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &ErrListingNotFound{ListingID: req.ListingID}
		}
		s.logger.Error().Err(err).Int64("listing_id", req.ListingID).Msg("failed to get listing")
		return nil, fmt.Errorf("failed to get listing: %w", err)
	}

	// Check if listing is active
	if !listing.IsActive {
		return nil, &ErrListingInactive{ListingID: req.ListingID}
	}

	// Check stock availability
	if listing.StockQuantity < req.Quantity {
		return nil, &ErrInsufficientStock{
			ListingID:      req.ListingID,
			ListingName:    listing.Name,
			RequestedQty:   req.Quantity,
			AvailableStock: listing.StockQuantity,
		}
	}

	// Validate storefront matches
	if listing.StorefrontID != req.StorefrontID {
		return nil, &ErrStorefrontMismatch{
			CartStorefrontID: req.StorefrontID,
			ItemStorefrontID: listing.StorefrontID,
		}
	}

	// Get or create cart
	var cart *domain.Cart
	if req.UserID != nil {
		cart, err = s.cartRepo.GetByUserAndStorefront(ctx, *req.UserID, req.StorefrontID)
	} else if req.SessionID != nil {
		cart, err = s.cartRepo.GetBySessionAndStorefront(ctx, *req.SessionID, req.StorefrontID)
	}

	if err != nil && err.Error() != "cart not found" {
		s.logger.Error().Err(err).Msg("failed to get cart")
		return nil, fmt.Errorf("failed to get cart: %w", err)
	}

	// Create new cart if not found
	if cart == nil {
		cart = &domain.Cart{
			UserID:       req.UserID,
			SessionID:    req.SessionID,
			StorefrontID: req.StorefrontID,
		}

		if err := s.cartRepo.Create(ctx, cart); err != nil {
			s.logger.Error().Err(err).Msg("failed to create cart")
			return nil, fmt.Errorf("failed to create cart: %w", err)
		}
	}

	// Check if item already exists in cart (same listing + variant)
	existingItem := s.findCartItem(cart.Items, req.ListingID, req.VariantID)

	if existingItem != nil {
		// Update quantity
		existingItem.Quantity += req.Quantity
		existingItem.PriceSnapshot = listing.Price // Update price snapshot

		if err := s.cartRepo.UpdateItem(ctx, existingItem); err != nil {
			s.logger.Error().Err(err).Int64("item_id", existingItem.ID).Msg("failed to update cart item")
			return nil, fmt.Errorf("failed to update cart item: %w", err)
		}
	} else {
		// Add new item
		newItem := &domain.CartItem{
			CartID:        cart.ID,
			ListingID:     req.ListingID,
			VariantID:     req.VariantID,
			Quantity:      req.Quantity,
			PriceSnapshot: listing.Price,
		}

		if err := s.cartRepo.AddItem(ctx, newItem); err != nil {
			s.logger.Error().Err(err).Msg("failed to add cart item")
			return nil, fmt.Errorf("failed to add cart item: %w", err)
		}

		// Reload cart items
		cart.Items, err = s.cartRepo.GetItemsByCartID(ctx, cart.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to reload cart items: %w", err)
		}
	}

	s.logger.Info().Int64("cart_id", cart.ID).Int64("listing_id", req.ListingID).Msg("item added to cart")
	return cart, nil
}

// UpdateCartItem updates the quantity of a cart item
func (s *cartService) UpdateCartItem(ctx context.Context, req *UpdateCartItemRequest) (*domain.Cart, error) {
	s.logger.Debug().
		Int64("cart_id", req.CartID).
		Int64("item_id", req.ItemID).
		Int32("quantity", req.Quantity).
		Msg("updating cart item")

	// Validate quantity
	if req.Quantity <= 0 {
		return nil, fmt.Errorf("%w: quantity must be greater than 0", ErrInvalidInput)
	}

	// Get cart
	cart, err := s.cartRepo.GetByID(ctx, req.CartID)
	if err != nil {
		if err.Error() == "cart not found" {
			return nil, ErrCartNotFound
		}
		s.logger.Error().Err(err).Int64("cart_id", req.CartID).Msg("failed to get cart")
		return nil, fmt.Errorf("failed to get cart: %w", err)
	}

	// Find item in cart
	var targetItem *domain.CartItem
	for _, item := range cart.Items {
		if item.ID == req.ItemID {
			targetItem = item
			break
		}
	}

	if targetItem == nil {
		return nil, ErrCartItemNotFound
	}

	// Fetch listing to validate stock
	listing, err := s.productsRepo.GetProductByID(ctx, targetItem.ListingID, &cart.StorefrontID)
	if err != nil {
		s.logger.Error().Err(err).Int64("listing_id", targetItem.ListingID).Msg("failed to get listing")
		return nil, fmt.Errorf("failed to get listing: %w", err)
	}

	// Check stock availability
	if listing.StockQuantity < req.Quantity {
		return nil, &ErrInsufficientStock{
			ListingID:      targetItem.ListingID,
			ListingName:    listing.Name,
			RequestedQty:   req.Quantity,
			AvailableStock: listing.StockQuantity,
		}
	}

	// Update quantity and price snapshot
	targetItem.Quantity = req.Quantity
	targetItem.PriceSnapshot = listing.Price

	if err := s.cartRepo.UpdateItem(ctx, targetItem); err != nil {
		s.logger.Error().Err(err).Int64("item_id", req.ItemID).Msg("failed to update cart item")
		return nil, fmt.Errorf("failed to update cart item: %w", err)
	}

	// Reload cart
	cart, err = s.cartRepo.GetByID(ctx, req.CartID)
	if err != nil {
		return nil, fmt.Errorf("failed to reload cart: %w", err)
	}

	s.logger.Info().Int64("cart_id", req.CartID).Int64("item_id", req.ItemID).Msg("cart item updated")
	return cart, nil
}

// RemoveFromCart removes an item from the cart
func (s *cartService) RemoveFromCart(ctx context.Context, cartID, itemID int64) error {
	s.logger.Debug().Int64("cart_id", cartID).Int64("item_id", itemID).Msg("removing from cart")

	// Verify cart exists
	cart, err := s.cartRepo.GetByID(ctx, cartID)
	if err != nil {
		if err.Error() == "cart not found" {
			return ErrCartNotFound
		}
		return fmt.Errorf("failed to get cart: %w", err)
	}

	// Verify item exists in cart
	itemExists := false
	for _, item := range cart.Items {
		if item.ID == itemID {
			itemExists = true
			break
		}
	}

	if !itemExists {
		return ErrCartItemNotFound
	}

	// Remove item
	if err := s.cartRepo.RemoveItem(ctx, cartID, itemID); err != nil {
		s.logger.Error().Err(err).Int64("cart_id", cartID).Int64("item_id", itemID).Msg("failed to remove cart item")
		return fmt.Errorf("failed to remove cart item: %w", err)
	}

	s.logger.Info().Int64("cart_id", cartID).Int64("item_id", itemID).Msg("item removed from cart")
	return nil
}

// GetCart retrieves a cart by user ID or session ID
func (s *cartService) GetCart(ctx context.Context, userID *int64, sessionID *string, storefrontID int64) (*domain.Cart, error) {
	s.logger.Debug().
		Interface("user_id", userID).
		Interface("session_id", sessionID).
		Int64("storefront_id", storefrontID).
		Msg("getting cart")

	// Validate input (must have either userID or sessionID)
	if (userID == nil && sessionID == nil) || (userID != nil && sessionID != nil) {
		return nil, fmt.Errorf("%w: must provide either user_id or session_id (not both)", ErrInvalidInput)
	}

	var cart *domain.Cart
	var err error

	if userID != nil {
		cart, err = s.cartRepo.GetByUserAndStorefront(ctx, *userID, storefrontID)
	} else {
		cart, err = s.cartRepo.GetBySessionAndStorefront(ctx, *sessionID, storefrontID)
	}

	if err != nil {
		if err.Error() == "cart not found" {
			return nil, ErrCartNotFound
		}
		s.logger.Error().Err(err).Msg("failed to get cart")
		return nil, fmt.Errorf("failed to get cart: %w", err)
	}

	return cart, nil
}

// ClearCart removes all items from a cart
func (s *cartService) ClearCart(ctx context.Context, cartID int64) error {
	s.logger.Debug().Int64("cart_id", cartID).Msg("clearing cart")

	// Verify cart exists
	_, err := s.cartRepo.GetByID(ctx, cartID)
	if err != nil {
		if err.Error() == "cart not found" {
			return ErrCartNotFound
		}
		return fmt.Errorf("failed to get cart: %w", err)
	}

	// Clear all items
	if err := s.cartRepo.ClearItems(ctx, cartID); err != nil {
		s.logger.Error().Err(err).Int64("cart_id", cartID).Msg("failed to clear cart")
		return fmt.Errorf("failed to clear cart: %w", err)
	}

	s.logger.Info().Int64("cart_id", cartID).Msg("cart cleared")
	return nil
}

// GetUserCarts retrieves all carts for a user (across all storefronts)
func (s *cartService) GetUserCarts(ctx context.Context, userID int64) ([]*domain.Cart, error) {
	s.logger.Debug().Int64("user_id", userID).Msg("getting user carts")

	carts, err := s.cartRepo.GetUserCarts(ctx, userID)
	if err != nil {
		s.logger.Error().Err(err).Int64("user_id", userID).Msg("failed to get user carts")
		return nil, fmt.Errorf("failed to get user carts: %w", err)
	}

	return carts, nil
}

// MergeSessionCartToUser merges a session cart into a user cart upon login
func (s *cartService) MergeSessionCartToUser(ctx context.Context, sessionID string, userID int64) error {
	s.logger.Debug().Str("session_id", sessionID).Int64("user_id", userID).Msg("merging session cart to user")

	// Get all session carts (may have carts from multiple storefronts)
	sessionCarts, err := s.cartRepo.GetUserCarts(ctx, 0) // This needs a different method
	if err != nil {
		// No session carts found is not an error
		s.logger.Debug().Msg("no session carts found to merge")
		return nil
	}

	for _, sessionCart := range sessionCarts {
		if sessionCart.SessionID == nil || *sessionCart.SessionID != sessionID {
			continue
		}

		// Get or create user cart for this storefront
		userCart, err := s.cartRepo.GetByUserAndStorefront(ctx, userID, sessionCart.StorefrontID)
		if err != nil && err.Error() != "cart not found" {
			return fmt.Errorf("failed to get user cart: %w", err)
		}

		if userCart == nil {
			// Create new user cart
			userCart = &domain.Cart{
				UserID:       &userID,
				StorefrontID: sessionCart.StorefrontID,
			}
			if err := s.cartRepo.Create(ctx, userCart); err != nil {
				return fmt.Errorf("failed to create user cart: %w", err)
			}
		}

		// Transfer items from session cart to user cart
		for _, item := range sessionCart.Items {
			// Check if item already exists in user cart
			existingItem := s.findCartItem(userCart.Items, item.ListingID, item.VariantID)

			if existingItem != nil {
				// Merge quantities
				existingItem.Quantity += item.Quantity
				if err := s.cartRepo.UpdateItem(ctx, existingItem); err != nil {
					return fmt.Errorf("failed to update cart item: %w", err)
				}
			} else {
				// Move item to user cart
				item.CartID = userCart.ID
				if err := s.cartRepo.AddItem(ctx, item); err != nil {
					return fmt.Errorf("failed to add cart item: %w", err)
				}
			}
		}

		// Delete session cart
		if err := s.cartRepo.Delete(ctx, sessionCart.ID); err != nil {
			s.logger.Warn().Err(err).Int64("session_cart_id", sessionCart.ID).Msg("failed to delete session cart")
		}
	}

	s.logger.Info().Str("session_id", sessionID).Int64("user_id", userID).Msg("session cart merged to user")
	return nil
}

// RecalculateCart recalculates cart totals and updates price snapshots
func (s *cartService) RecalculateCart(ctx context.Context, cartID int64) (*domain.Cart, error) {
	s.logger.Debug().Int64("cart_id", cartID).Msg("recalculating cart")

	// Get cart with items
	cart, err := s.cartRepo.GetByID(ctx, cartID)
	if err != nil {
		if err.Error() == "cart not found" {
			return nil, ErrCartNotFound
		}
		return nil, fmt.Errorf("failed to get cart: %w", err)
	}

	// Update price snapshots for all items
	for _, item := range cart.Items {
		listing, err := s.productsRepo.GetProductByID(ctx, item.ListingID, &cart.StorefrontID)
		if err != nil {
			s.logger.Warn().Err(err).Int64("listing_id", item.ListingID).Msg("failed to get listing for recalculation")
			continue
		}

		// Update price snapshot if changed
		if listing.Price != item.PriceSnapshot {
			item.PriceSnapshot = listing.Price
			if err := s.cartRepo.UpdateItem(ctx, item); err != nil {
				s.logger.Warn().Err(err).Int64("item_id", item.ID).Msg("failed to update item price")
			}
		}
	}

	// Reload cart to get updated totals
	cart, err = s.cartRepo.GetByID(ctx, cartID)
	if err != nil {
		return nil, fmt.Errorf("failed to reload cart: %w", err)
	}

	s.logger.Info().Int64("cart_id", cartID).Msg("cart recalculated")
	return cart, nil
}

// ValidateCartItems validates cart items against current listing prices and stock
func (s *cartService) ValidateCartItems(ctx context.Context, cartID int64) ([]PriceChangeItem, error) {
	s.logger.Debug().Int64("cart_id", cartID).Msg("validating cart items")

	// Get cart with items
	cart, err := s.cartRepo.GetByID(ctx, cartID)
	if err != nil {
		if err.Error() == "cart not found" {
			return nil, ErrCartNotFound
		}
		return nil, fmt.Errorf("failed to get cart: %w", err)
	}

	priceChanges := []PriceChangeItem{}

	for _, item := range cart.Items {
		listing, err := s.productsRepo.GetProductByID(ctx, item.ListingID, &cart.StorefrontID)
		if err != nil {
			s.logger.Warn().Err(err).Int64("listing_id", item.ListingID).Msg("failed to get listing for validation")
			continue
		}

		// Check for price changes
		if listing.Price != item.PriceSnapshot {
			priceChanges = append(priceChanges, PriceChangeItem{
				ListingID:     item.ListingID,
				ListingName:   listing.Name,
				OldPrice:      item.PriceSnapshot,
				NewPrice:      listing.Price,
				PriceIncrease: listing.Price > item.PriceSnapshot,
			})
		}

		// Check stock availability
		if listing.StockQuantity < item.Quantity {
			s.logger.Warn().
				Int64("listing_id", item.ListingID).
				Int32("requested", item.Quantity).
				Int32("available", listing.StockQuantity).
				Msg("insufficient stock for cart item")
		}
	}

	return priceChanges, nil
}

// Helper methods

// validateAddToCartRequest validates the add to cart request
func (s *cartService) validateAddToCartRequest(req *AddToCartRequest) error {
	if (req.UserID == nil && req.SessionID == nil) || (req.UserID != nil && req.SessionID != nil) {
		return fmt.Errorf("%w: must provide either user_id or session_id (not both)", ErrInvalidInput)
	}

	if req.StorefrontID <= 0 {
		return fmt.Errorf("%w: storefront_id must be greater than 0", ErrInvalidInput)
	}

	if req.ListingID <= 0 {
		return fmt.Errorf("%w: listing_id must be greater than 0", ErrInvalidInput)
	}

	if req.Quantity <= 0 {
		return fmt.Errorf("%w: quantity must be greater than 0", ErrInvalidInput)
	}

	return nil
}

// findCartItem finds an item in cart by listing ID and variant ID
func (s *cartService) findCartItem(items []*domain.CartItem, listingID int64, variantID *int64) *domain.CartItem {
	for _, item := range items {
		if item.ListingID == listingID {
			// Check if variant matches
			if (variantID == nil && item.VariantID == nil) ||
				(variantID != nil && item.VariantID != nil && *variantID == *item.VariantID) {
				return item
			}
		}
	}
	return nil
}

// UpdateCartItemByItemID updates cart item quantity by item ID with ownership verification
func (s *cartService) UpdateCartItemByItemID(ctx context.Context, cartItemID int64, quantity int32, userID *int64, sessionID *string) (*domain.Cart, error) {
	s.logger.Debug().
		Int64("cart_item_id", cartItemID).
		Int32("quantity", quantity).
		Interface("user_id", userID).
		Interface("session_id", sessionID).
		Msg("updating cart item by item ID")

	// Validate quantity
	if quantity <= 0 {
		return nil, fmt.Errorf("%w: quantity must be greater than 0", ErrInvalidInput)
	}

	// Get cart item to find cart_id
	cartItem, err := s.cartRepo.GetCartItemByID(ctx, cartItemID)
	if err != nil {
		s.logger.Error().Err(err).Int64("cart_item_id", cartItemID).Msg("failed to get cart item")
		return nil, fmt.Errorf("failed to get cart item: %w", err)
	}

	// Get cart to verify ownership
	cart, err := s.cartRepo.GetByID(ctx, cartItem.CartID)
	if err != nil {
		s.logger.Error().Err(err).Int64("cart_id", cartItem.CartID).Msg("failed to get cart")
		return nil, fmt.Errorf("failed to get cart: %w", err)
	}

	// Verify ownership (user_id or session_id matches cart owner)
	if userID != nil {
		if cart.UserID == nil || *cart.UserID != *userID {
			return nil, ErrUnauthorized
		}
	} else if sessionID != nil {
		if cart.SessionID == nil || *cart.SessionID != *sessionID {
			return nil, ErrUnauthorized
		}
	} else {
		return nil, fmt.Errorf("%w: must provide either user_id or session_id", ErrInvalidInput)
	}

	// Call existing UpdateCartItem method
	return s.UpdateCartItem(ctx, &UpdateCartItemRequest{
		CartID:   cartItem.CartID,
		ItemID:   cartItemID,
		Quantity: quantity,
	})
}

// RemoveFromCartByItemID removes cart item by item ID with ownership verification
func (s *cartService) RemoveFromCartByItemID(ctx context.Context, cartItemID int64, userID *int64, sessionID *string) error {
	s.logger.Debug().
		Int64("cart_item_id", cartItemID).
		Interface("user_id", userID).
		Interface("session_id", sessionID).
		Msg("removing cart item by item ID")

	// Get cart item to find cart_id
	cartItem, err := s.cartRepo.GetCartItemByID(ctx, cartItemID)
	if err != nil {
		s.logger.Error().Err(err).Int64("cart_item_id", cartItemID).Msg("failed to get cart item")
		return fmt.Errorf("failed to get cart item: %w", err)
	}

	// Get cart to verify ownership
	cart, err := s.cartRepo.GetByID(ctx, cartItem.CartID)
	if err != nil {
		s.logger.Error().Err(err).Int64("cart_id", cartItem.CartID).Msg("failed to get cart")
		return fmt.Errorf("failed to get cart: %w", err)
	}

	// Verify ownership (user_id or session_id matches cart owner)
	if userID != nil {
		if cart.UserID == nil || *cart.UserID != *userID {
			return ErrUnauthorized
		}
	} else if sessionID != nil {
		if cart.SessionID == nil || *cart.SessionID != *sessionID {
			return ErrUnauthorized
		}
	} else {
		return fmt.Errorf("%w: must provide either user_id or session_id", ErrInvalidInput)
	}

	// Call existing RemoveFromCart method
	return s.RemoveFromCart(ctx, cartItem.CartID, cartItemID)
}
