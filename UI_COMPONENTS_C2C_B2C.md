# 🎨 UI компоненты для C2C и B2C

## 📱 Компоненты страницы товара

### 1. Универсальная страница товара/объявления

```typescript
// app/[locale]/marketplace/[id]/page.tsx

interface ListingPageProps {
  listing: MarketplaceListing;
}

export default function ListingPage({ listing }: ListingPageProps) {
  const isStorefrontProduct = listing.storefront_id != null;
  
  return (
    <div className="container">
      {/* Общие элементы */}
      <ImageGallery images={listing.images} />
      <ProductInfo listing={listing} />
      
      {/* Разные действия для C2C и B2C */}
      {isStorefrontProduct ? (
        <B2CActions listing={listing} />
      ) : (
        <C2CActions listing={listing} />
      )}
      
      {/* Информация о продавце */}
      {isStorefrontProduct ? (
        <StorefrontInfo storefront={listing.storefront} />
      ) : (
        <PrivateSellerInfo seller={listing.user} />
      )}
    </div>
  );
}
```

### 2. C2C Actions Component

```typescript
// components/marketplace/C2CActions.tsx

export const C2CActions: React.FC<{ listing: MarketplaceListing }> = ({ listing }) => {
  const { user } = useAuth();
  const router = useRouter();
  const t = useTranslations();
  
  if (!user) {
    return (
      <div className="card bg-base-200 p-6">
        <h3 className="text-lg font-semibold mb-4">
          {t('c2c.contactSellerTitle')}
        </h3>
        <p className="text-base-content/70 mb-4">
          {t('c2c.loginToContact')}
        </p>
        <button 
          onClick={() => router.push('/auth/login')}
          className="btn btn-primary"
        >
          {t('auth.login')}
        </button>
      </div>
    );
  }
  
  return (
    <div className="card bg-base-200 p-6">
      <h3 className="text-lg font-semibold mb-4">
        {t('c2c.howToBuy')}
      </h3>
      
      {/* Способы оплаты */}
      <div className="mb-6">
        <h4 className="font-medium mb-2">{t('c2c.paymentMethods')}</h4>
        <div className="space-y-2">
          <label className="flex items-center gap-3 p-3 bg-base-100 rounded-lg">
            <input type="radio" name="payment" className="radio radio-primary" />
            <div>
              <div className="font-medium">💵 {t('c2c.cash')}</div>
              <div className="text-sm text-base-content/60">
                {t('c2c.cashDescription')}
              </div>
            </div>
          </label>
          
          <label className="flex items-center gap-3 p-3 bg-base-100 rounded-lg">
            <input type="radio" name="payment" className="radio radio-primary" />
            <div>
              <div className="font-medium">💳 {t('c2c.bankTransfer')}</div>
              <div className="text-sm text-base-content/60">
                {t('c2c.bankTransferDescription')}
              </div>
            </div>
          </label>
          
          {listing.shipping_available && (
            <label className="flex items-center gap-3 p-3 bg-base-100 rounded-lg">
              <input type="radio" name="payment" className="radio radio-primary" />
              <div>
                <div className="font-medium">📦 {t('c2c.cod')}</div>
                <div className="text-sm text-base-content/60">
                  {t('c2c.codDescription')}
                </div>
              </div>
            </label>
          )}
        </div>
      </div>
      
      {/* CTA кнопка */}
      <button 
        onClick={() => router.push(`/chat?listing_id=${listing.id}&seller_id=${listing.user_id}`)}
        className="btn btn-primary btn-lg w-full"
      >
        <svg className="w-5 h-5 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} 
            d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
        </svg>
        {t('c2c.contactSeller')}
      </button>
      
      {/* Trust indicators */}
      <div className="mt-4 p-4 bg-warning/10 rounded-lg">
        <h5 className="font-medium text-warning-content mb-2">
          {t('c2c.safetyTips')}
        </h5>
        <ul className="text-sm space-y-1 text-warning-content/80">
          <li>• {t('c2c.meetInPublic')}</li>
          <li>• {t('c2c.inspectBeforePay')}</li>
          <li>• {t('c2c.useSecurePayment')}</li>
          <li>• {t('c2c.keepChatHistory')}</li>
        </ul>
      </div>
    </div>
  );
};
```

### 3. B2C Actions Component

```typescript
// components/marketplace/B2CActions.tsx

export const B2CActions: React.FC<{ listing: StorefrontProduct }> = ({ listing }) => {
  const [quantity, setQuantity] = useState(1);
  const [selectedVariant, setSelectedVariant] = useState(null);
  const { addToCart, isLoading } = useCart();
  const t = useTranslations();
  
  return (
    <div className="card bg-base-200 p-6">
      {/* Варианты товара */}
      {listing.variants && listing.variants.length > 0 && (
        <div className="mb-6">
          <h4 className="font-medium mb-3">{t('product.selectVariant')}</h4>
          <div className="grid grid-cols-2 gap-2">
            {listing.variants.map((variant) => (
              <button
                key={variant.id}
                onClick={() => setSelectedVariant(variant)}
                className={`btn ${
                  selectedVariant?.id === variant.id 
                    ? 'btn-primary' 
                    : 'btn-outline'
                }`}
              >
                {variant.name}
                {variant.price_diff && (
                  <span className="text-xs ml-1">
                    {variant.price_diff > 0 ? '+' : ''}{variant.price_diff} $
                  </span>
                )}
              </button>
            ))}
          </div>
        </div>
      )}
      
      {/* Количество */}
      <div className="mb-6">
        <h4 className="font-medium mb-3">{t('product.quantity')}</h4>
        <div className="flex items-center gap-4">
          <div className="join">
            <button 
              className="btn join-item"
              onClick={() => setQuantity(Math.max(1, quantity - 1))}
            >
              -
            </button>
            <input 
              type="number" 
              className="input join-item w-20 text-center" 
              value={quantity}
              onChange={(e) => setQuantity(parseInt(e.target.value) || 1)}
              min="1"
              max={listing.stock_quantity}
            />
            <button 
              className="btn join-item"
              onClick={() => setQuantity(Math.min(listing.stock_quantity, quantity + 1))}
            >
              +
            </button>
          </div>
          <span className="text-sm text-base-content/60">
            {t('product.available', { count: listing.stock_quantity })}
          </span>
        </div>
      </div>
      
      {/* Цена */}
      <div className="mb-6">
        <div className="flex items-baseline gap-2">
          <span className="text-3xl font-bold text-primary">
            {calculatePrice(listing, selectedVariant, quantity)} {listing.currency}
          </span>
          {quantity > 1 && (
            <span className="text-sm text-base-content/60">
              ({listing.price} $ × {quantity})
            </span>
          )}
        </div>
      </div>
      
      {/* Кнопки действий */}
      <div className="space-y-3">
        <button 
          onClick={() => addToCart(listing, selectedVariant, quantity)}
          className="btn btn-primary btn-lg w-full"
          disabled={isLoading || listing.stock_status === 'out_of_stock'}
        >
          {isLoading && <span className="loading loading-spinner" />}
          <svg className="w-5 h-5 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} 
              d="M3 3h2l.4 2M7 13h10l4-8H5.4M7 13L5.4 5M7 13l-2.293 2.293c-.63.63-.184 1.707.707 1.707H17m0 0a2 2 0 100 4 2 2 0 000-4zm-8 2a2 2 0 11-4 0 2 2 0 014 0z" />
          </svg>
          {t('cart.addToCart')}
        </button>
        
        <button 
          onClick={() => buyNow(listing, selectedVariant, quantity)}
          className="btn btn-outline btn-lg w-full"
        >
          {t('product.buyNow')}
        </button>
      </div>
      
      {/* Преимущества покупки в магазине */}
      <div className="mt-6 space-y-2">
        <div className="flex items-center gap-2 text-sm">
          <svg className="w-4 h-4 text-success" fill="currentColor" viewBox="0 0 20 20">
            <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clipRule="evenodd" />
          </svg>
          <span>{t('b2c.securePayment')}</span>
        </div>
        <div className="flex items-center gap-2 text-sm">
          <svg className="w-4 h-4 text-success" fill="currentColor" viewBox="0 0 20 20">
            <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clipRule="evenodd" />
          </svg>
          <span>{t('b2c.buyerProtection')}</span>
        </div>
        <div className="flex items-center gap-2 text-sm">
          <svg className="w-4 h-4 text-success" fill="currentColor" viewBox="0 0 20 20">
            <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clipRule="evenodd" />
          </svg>
          <span>{t('b2c.trackingAvailable')}</span>
        </div>
      </div>
    </div>
  );
};
```

### 4. Seller Info Components

```typescript
// components/marketplace/PrivateSellerInfo.tsx

export const PrivateSellerInfo: React.FC<{ seller: User }> = ({ seller }) => {
  const { data: metrics } = useC2CSellerMetrics(seller.id);
  const t = useTranslations();
  
  return (
    <div className="card bg-base-100 shadow-xl">
      <div className="card-body">
        <h3 className="card-title">
          {t('seller.privateSellerInfo')}
        </h3>
        
        {/* Seller avatar and name */}
        <div className="flex items-center gap-4 mb-4">
          <div className="avatar">
            <div className="w-16 rounded-full">
              <img src={seller.avatar || '/default-avatar.png'} />
            </div>
          </div>
          <div>
            <h4 className="font-semibold">{seller.name}</h4>
            <p className="text-sm text-base-content/60">
              {t('seller.memberSince', { date: formatDate(seller.created_at) })}
            </p>
          </div>
        </div>
        
        {/* Trust indicators */}
        <div className="stats stats-vertical lg:stats-horizontal shadow">
          <div className="stat">
            <div className="stat-title">{t('seller.completedDeals')}</div>
            <div className="stat-value text-primary">
              {metrics?.completed_transactions || 0}
            </div>
          </div>
          
          <div className="stat">
            <div className="stat-title">{t('seller.responseTime')}</div>
            <div className="stat-value text-secondary">
              {metrics?.avg_response_time || '< 1ч'}
            </div>
          </div>
          
          <div className="stat">
            <div className="stat-title">{t('seller.trustScore')}</div>
            <div className="stat-value">
              <div className="flex items-center">
                {renderStars(metrics?.trust_score || 0)}
              </div>
            </div>
          </div>
        </div>
        
        {/* Verified badges */}
        <div className="flex flex-wrap gap-2 mt-4">
          {metrics?.verified_phone && (
            <span className="badge badge-success gap-2">
              <svg className="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
                <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clipRule="evenodd" />
              </svg>
              {t('seller.phoneVerified')}
            </span>
          )}
          
          {metrics?.verified_email && (
            <span className="badge badge-info gap-2">
              <svg className="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
                <path d="M2.003 5.884L10 9.882l7.997-3.998A2 2 0 0016 4H4a2 2 0 00-1.997 1.884z" />
                <path d="M18 8.118l-8 4-8-4V14a2 2 0 002 2h12a2 2 0 002-2V8.118z" />
              </svg>
              {t('seller.emailVerified')}
            </span>
          )}
        </div>
        
        {/* Preferred payment methods */}
        {metrics?.preferred_payment_methods && (
          <div className="mt-4">
            <h5 className="font-medium mb-2">
              {t('seller.preferredPayments')}
            </h5>
            <div className="flex flex-wrap gap-2">
              {metrics.preferred_payment_methods.map((method) => (
                <span key={method} className="badge badge-outline">
                  {getPaymentMethodIcon(method)} {t(`payment.${method}`)}
                </span>
              ))}
            </div>
          </div>
        )}
        
        {/* Recent reviews */}
        <RecentReviews sellerId={seller.id} limit={3} />
      </div>
    </div>
  );
};
```

### 5. Search Results Mixed View

```typescript
// components/marketplace/SearchResults.tsx

export const SearchResults: React.FC<{ results: MarketplaceListing[] }> = ({ results }) => {
  const t = useTranslations();
  
  return (
    <div className="space-y-4">
      {results.map((listing) => (
        <div key={listing.id} className="card card-side bg-base-100 shadow-xl">
          <figure className="w-48">
            <img src={listing.images[0]?.url || '/placeholder.jpg'} />
          </figure>
          
          <div className="card-body">
            <div className="flex justify-between items-start">
              <div className="flex-1">
                <h2 className="card-title">
                  {listing.title}
                  {listing.storefront_id && (
                    <span className="badge badge-primary badge-sm">
                      {t('listing.fromStore')}
                    </span>
                  )}
                </h2>
                <p className="text-2xl font-bold text-primary mt-2">
                  {listing.price} $
                </p>
                <p className="text-sm text-base-content/60 mt-1">
                  {listing.location}
                </p>
              </div>
              
              <div className="card-actions">
                {listing.storefront_id ? (
                  <div className="flex flex-col gap-2">
                    <AddToCartButton 
                      product={listing} 
                      size="sm"
                      className="btn-primary"
                    />
                    <Link 
                      href={`/marketplace/${listing.id}`}
                      className="btn btn-outline btn-sm"
                    >
                      {t('common.details')}
                    </Link>
                  </div>
                ) : (
                  <Link 
                    href={`/marketplace/${listing.id}`}
                    className="btn btn-primary"
                  >
                    {t('c2c.contactSeller')}
                  </Link>
                )}
              </div>
            </div>
            
            {/* Seller info */}
            <div className="mt-4 flex items-center gap-4 text-sm">
              {listing.storefront_id ? (
                <>
                  <span className="flex items-center gap-1">
                    <svg className="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
                      <path fillRule="evenodd" d="M6 2a1 1 0 00-1 1v1H4a2 2 0 00-2 2v10a2 2 0 002 2h12a2 2 0 002-2V6a2 2 0 00-2-2h-1V3a1 1 0 10-2 0v1H7V3a1 1 0 00-1-1zm0 5a1 1 0 000 2h8a1 1 0 100-2H6z" clipRule="evenodd" />
                    </svg>
                    {listing.storefront.name}
                  </span>
                  {listing.storefront.is_verified && (
                    <span className="badge badge-success badge-xs">
                      {t('storefront.verified')}
                    </span>
                  )}
                </>
              ) : (
                <>
                  <span className="flex items-center gap-1">
                    <svg className="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
                      <path fillRule="evenodd" d="M10 9a3 3 0 100-6 3 3 0 000 6zm-7 9a7 7 0 1114 0H3z" clipRule="evenodd" />
                    </svg>
                    {listing.user.name}
                  </span>
                  {listing.user.trust_score >= 4.5 && (
                    <span className="badge badge-warning badge-xs">
                      ⭐ {t('seller.trusted')}
                    </span>
                  )}
                </>
              )}
            </div>
          </div>
        </div>
      ))}
    </div>
  );
};
```

## 🎯 Ключевые отличия в UI

### C2C интерфейс:
1. **Нет корзины** - только кнопка "Связаться"
2. **Чат-центричный** процесс покупки
3. **Trust badges** для продавцов
4. **Предупреждения** о безопасности
5. **История сделок** продавца

### B2C интерфейс:
1. **Полноценная корзина** и checkout
2. **Варианты товаров** и количество
3. **Мгновенная покупка** (Buy Now)
4. **Отслеживание заказов**
5. **Возвраты и гарантии**

### Общие элементы:
1. **Галерея изображений**
2. **Описание и характеристики**
3. **Отзывы и рейтинги**
4. **Похожие товары**
5. **Местоположение** (с приватностью)