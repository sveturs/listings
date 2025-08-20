# Storefront Product Chat Integration TODO

## Current Status
✅ **Completed:**
- Fixed translation namespace issue (Chat → chat)
- Added support for storefront_product_id in MessageInput component
- Backend already supports storefront_product_id in chat tables
- Added placeholder text for storefront products in chat

## Pending Implementation

### 1. Backend API Endpoint for Storefront Products
**Priority: HIGH**

Need to create endpoint: `GET /api/v1/storefronts/products/{id}`

```go
// Handler method needed in storefront package
func (h *Handler) GetStorefrontProduct(c *fiber.Ctx) error {
    productID, err := strconv.Atoi(c.Params("id"))
    if err != nil {
        return utils.ErrorResponse(c, fiber.StatusBadRequest, "invalid.productId")
    }
    
    product, err := h.service.GetStorefrontProduct(c.Context(), productID)
    if err != nil {
        return utils.ErrorResponse(c, fiber.StatusNotFound, "product.notFound")
    }
    
    return utils.SuccessResponse(c, fiber.StatusOK, product, "product.retrieved")
}
```

### 2. Database Query for Storefront Products
**Priority: HIGH**

The `storefront_products` table exists and needs to be queried properly:

```sql
SELECT 
    sp.id,
    sp.storefront_id,
    sp.name as title,
    sp.description,
    sp.price,
    sp.currency,
    sp.status,
    sp.created_at,
    sp.updated_at,
    s.user_id,
    u.name as seller_name
FROM storefront_products sp
JOIN storefronts s ON sp.storefront_id = s.id
JOIN users u ON s.user_id = u.id
WHERE sp.id = $1 AND sp.status = 'active'
```

### 3. Frontend Updates
**Priority: MEDIUM**

Update ChatWindow.tsx to properly fetch storefront product info:

```typescript
// Replace placeholder with actual API call
if (isNewChat && initialStorefrontProductId && !listingInfo && !isContactChat) {
    const apiUrl = configManager.getApiUrl();
    fetch(`${apiUrl}/api/v1/storefronts/products/${initialStorefrontProductId}`, {
        credentials: 'include',
        headers: {
            'Content-Type': 'application/json',
        },
    })
    .then(res => res.json())
    .then(result => {
        const data = result.data || result;
        if (data && data.id) {
            setListingInfo({
                id: data.id,
                title: data.name || data.title,
                images: data.images || [],
                user_id: data.user_id || initialSellerId || 0,
            });
        }
    })
    .catch(err => console.error('Error loading storefront product:', err));
}
```

### 4. Image Support for Storefront Products
**Priority: LOW**

The `storefront_product_images` table exists and should be included in the product query:

```sql
-- Add to product query
LEFT JOIN (
    SELECT 
        storefront_product_id,
        JSON_AGG(
            JSON_BUILD_OBJECT(
                'id', id,
                'public_url', image_url,
                'thumbnail_url', thumbnail_url
            ) ORDER BY display_order
        ) as images
    FROM storefront_product_images
    GROUP BY storefront_product_id
) spi ON sp.id = spi.storefront_product_id
```

### 5. Link to Storefront Product Page
**Priority: LOW**

Update the link in ChatWindow.tsx to point to the correct storefront product page:

```typescript
// Line 506 in ChatWindow.tsx
<Link
    href={`/${locale}/storefronts/products/${initialStorefrontProductId || listingInfo?.id}`}
    className="btn btn-ghost btn-xs sm:btn-sm"
    title={t('viewProduct')}
>
```

## Testing Checklist
- [ ] Create backend endpoint for getting storefront product by ID
- [ ] Test that storefront product info loads in chat
- [ ] Verify images display correctly if available
- [ ] Test link to storefront product page works
- [ ] Ensure chat creation with storefront_product_id works
- [ ] Verify chat messages are properly associated with storefront products