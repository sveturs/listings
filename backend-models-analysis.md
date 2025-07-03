# Backend Models and Naming Conventions Analysis

## Overview
This document provides a comprehensive analysis of the Go models and their JSON tags in the SveTu Backend system, focusing on field naming patterns and conventions.

## Naming Convention Summary

### Primary Pattern: **snake_case**
The backend predominantly uses **snake_case** for JSON field names. This is the standard convention across all models.

### Key ID Field Patterns

#### User-related IDs
- `user_id` - Standard pattern used in most models
- Examples:
  - `MarketplaceListing.UserID` → `json:"user_id"`
  - `Storefront.UserID` → `json:"user_id"`
  - `StorefrontStaff.UserID` → `json:"user_id"`
  - `UserBalance.UserID` → `json:"user_id"`
  - `BalanceTransaction.UserID` → `json:"user_id"`
  - `StorefrontInventoryMovement.UserID` → `json:"user_id"`

#### Other ID Fields
- `id` - Primary key of entities
- `listing_id` - Reference to marketplace listings
- `storefront_id` - Reference to storefronts
- `category_id` - Reference to categories
- `chat_id` - Reference to chats
- `sender_id` / `receiver_id` - Message participants
- `buyer_id` / `seller_id` - Transaction participants
- `product_id` - Reference to products
- `attribute_id` - Reference to attributes

## Core Models Structure

### 1. **User Model** (`models.go`)
```go
type User struct {
    ID         int       `json:"id"`
    Name       string    `json:"name"`
    Email      string    `json:"email"`
    GoogleID   string    `json:"google_id"`
    PictureURL string    `json:"picture_url"`
    Phone      *string   `json:"phone,omitempty"`
    Password   *string   `json:"-"`
    Provider   string    `json:"provider"`
    CreatedAt  time.Time `json:"created_at"`
}
```

### 2. **MarketplaceListing Model** (`models.go`)
```go
type MarketplaceListing struct {
    ID              int       `json:"id"`
    UserID          int       `json:"user_id"`
    CategoryID      int       `json:"category_id"`
    Title           string    `json:"title"`
    Description     string    `json:"description"`
    Price           float64   `json:"price"`
    Status          string    `json:"status"`
    ViewsCount      int       `json:"views_count"`
    CreatedAt       time.Time `json:"created_at"`
    UpdatedAt       time.Time `json:"updated_at"`
    StorefrontID    *int      `json:"storefront_id,omitempty"`
    // ... additional fields
}
```

### 3. **Storefront Models** (`storefront.go`)
```go
type Storefront struct {
    ID               int              `json:"id"`
    UserID           int              `json:"user_id"`
    Slug             string           `json:"slug"`
    Name             string           `json:"name"`
    IsActive         bool             `json:"is_active"`
    IsVerified       bool             `json:"is_verified"`
    ProductsCount    int              `json:"products_count"`
    SubscriptionPlan SubscriptionPlan `json:"subscription_plan"`
    CreatedAt        time.Time        `json:"created_at"`
    UpdatedAt        time.Time        `json:"updated_at"`
    // ... additional fields
}
```

### 4. **StorefrontProduct Model** (`storefront_product.go`)
Uses both `json` and `db` tags:
```go
type StorefrontProduct struct {
    ID            int     `json:"id" db:"id"`
    StorefrontID  int     `json:"storefront_id" db:"storefront_id"`
    Name          string  `json:"name" db:"name"`
    Price         float64 `json:"price" db:"price"`
    CategoryID    int     `json:"category_id" db:"category_id"`
    StockQuantity int     `json:"stock_quantity" db:"stock_quantity"`
    IsActive      bool    `json:"is_active" db:"is_active"`
    ViewCount     int     `json:"view_count" db:"view_count"`
    // ... additional fields
}
```

### 5. **Chat Models** (`marketplace_chat.go`)
```go
type MarketplaceMessage struct {
    ID         int       `json:"id"`
    ChatID     int       `json:"chat_id"`
    ListingID  int       `json:"listing_id"`
    SenderID   int       `json:"sender_id"`
    ReceiverID int       `json:"receiver_id"`
    Content    string    `json:"content"`
    IsRead     bool      `json:"is_read"`
    CreatedAt  time.Time `json:"created_at"`
}

type MarketplaceChat struct {
    ID            int       `json:"id"`
    ListingID     int       `json:"listing_id"`
    BuyerID       int       `json:"buyer_id"`
    SellerID      int       `json:"seller_id"`
    LastMessageAt time.Time `json:"last_message_at"`
    IsArchived    bool      `json:"is_archived"`
}
```

### 6. **Balance Models** (`balance.go`)
```go
type UserBalance struct {
    UserID        int       `json:"user_id"`
    Balance       float64   `json:"balance"`
    FrozenBalance float64   `json:"frozen_balance"`
    Currency      string    `json:"currency"`
    UpdatedAt     time.Time `json:"updated_at"`
}
```

## Special JSON Tag Patterns

### 1. **Omitempty Pattern**
Used for optional fields:
- `json:"phone,omitempty"`
- `json:"storefront_id,omitempty"`
- `json:"latitude,omitempty"`

### 2. **Hidden Fields Pattern**
Using `-` to exclude from JSON:
- `Password *string json:"-"`

### 3. **Dual Tagging Pattern**
Some models use both `json` and `db` tags (for database mapping):
- `json:"storefront_id" db:"storefront_id"`

## Common Field Types and Conventions

### Timestamps
- `created_at` - Creation timestamp
- `updated_at` - Last update timestamp
- `deleted_at` - Soft delete timestamp
- `last_message_at` - Domain-specific timestamps

### Boolean Fields
- `is_active` - Active status
- `is_verified` - Verification status
- `is_read` - Read status
- `is_archived` - Archive status
- `has_attachments` - Presence indicators

### Count Fields
- `views_count` - View counts
- `products_count` - Product counts
- `reviews_count` - Review counts
- `unread_count` - Unread message counts

### Relationship Fields
Always use snake_case with `_id` suffix:
- `user_id`
- `listing_id`
- `category_id`
- `storefront_id`

## API Endpoint Patterns

From handler analysis, the API follows RESTful conventions with snake_case in:
- Query parameters
- Request body fields
- Response body fields

## Swagger Documentation Patterns

Swagger annotations use the same snake_case convention:
```go
// @Param user_id path int true "User ID"
// @Success 200 {object} utils.SuccessResponseSwag{data=models.UserBalance}
```

## Key Observations

1. **Consistency**: The backend strictly follows snake_case for all JSON fields
2. **No camelCase**: Unlike many Go APIs, this backend does NOT use camelCase
3. **Database Alignment**: Field names align with PostgreSQL column naming conventions
4. **Type Safety**: Uses proper Go types with JSON/DB tags for marshaling
5. **Null Handling**: Uses pointers for nullable fields with omitempty tags

## Recommendations for Frontend Integration

When working with this backend from the frontend:
1. Always expect and send snake_case field names
2. User ID fields will always be `user_id`, never `userId` or `userID`
3. All timestamp fields end with `_at` suffix
4. All ID reference fields end with `_id` suffix
5. Boolean flags typically start with `is_` or `has_` prefix
6. Count fields typically end with `_count` suffix

## Common Model Patterns

### Request DTOs
```go
type CreateMessageRequest struct {
    ListingID  int    `json:"listing_id"`
    ChatID     int    `json:"chat_id"`
    ReceiverID int    `json:"receiver_id" validate:"required"`
    Content    string `json:"content" validate:"required"`
}
```

### Response Wrappers
```go
type SuccessResponse struct {
    Success bool        `json:"success"`
    Data    interface{} `json:"data,omitempty"`
    Message string      `json:"message,omitempty"`
}
```

### Pagination
```go
type PaginatedResponse struct {
    Data       interface{} `json:"data"`
    Page       int         `json:"page"`
    PageSize   int         `json:"page_size"`
    Total      int         `json:"total"`
    TotalPages int         `json:"total_pages"`
}
```