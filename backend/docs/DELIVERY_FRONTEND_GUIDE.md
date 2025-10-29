# Frontend Delivery Integration Guide

## 📋 Table of Contents

- [Overview](#overview)
- [Architecture](#architecture)
- [Components](#components)
- [Redux State Management](#redux-state-management)
- [API Service](#api-service)
- [Usage Examples](#usage-examples)
- [Best Practices](#best-practices)
- [Troubleshooting](#troubleshooting)

---

## Overview

Frontend delivery integration построен на современном стеке React 19 + Next.js 15 с использованием Redux Toolkit для state management и BFF proxy для безопасной коммуникации с backend.

### Key Features

- ✅ **Unified Component**: `UnifiedDeliverySelector` для всех контекстов
- ✅ **Redux State**: Централизованное управление состоянием с кэшированием
- ✅ **BFF Proxy**: Безопасная коммуникация через `/api/v2`
- ✅ **TypeScript**: Полная типизация для type safety
- ✅ **Caching**: TTL-based кэш на 5 минут
- ✅ **Error Handling**: Graceful degradation и fallback UI

### Technology Stack

- **React:** 19.x
- **Next.js:** 15.x (App Router)
- **Redux Toolkit:** State management + caching
- **TypeScript:** Full type safety
- **Tailwind CSS + DaisyUI:** Styling
- **Heroicons:** Icons

---

## Architecture

### Component Hierarchy

```
┌────────────────────────────────────┐
│   Cart / Checkout / Admin Pages   │
└────────────────┬───────────────────┘
                 │
┌────────────────▼───────────────────┐
│   UnifiedDeliverySelector          │
│   (Main delivery selection UI)     │
└────────────────┬───────────────────┘
                 │
        ┌────────┴────────┐
        │                 │
┌───────▼──────┐  ┌──────▼────────┐
│ Redux Store  │  │ deliveryService│
│ (deliverySlice)│  │ (API wrapper)  │
└───────┬──────┘  └──────┬────────┘
        │                 │
        │         ┌───────▼─────────┐
        │         │   apiClient     │
        │         │ (BFF /api/v2/*) │
        │         └───────┬─────────┘
        │                 │
        └─────────────────▼─────────────────┐
                   Backend API               │
                   (/api/v1/delivery/*)      │
                   └───────────────────────────┘
```

### Data Flow

```
1. User interaction
   └─▶ Component dispatches Redux action

2. Redux Thunk
   ├─▶ Check cache (5 min TTL)
   └─▶ If miss: Call deliveryService

3. deliveryService
   └─▶ apiClient.post('/delivery/calculate-universal')

4. BFF Proxy (/api/v2/*)
   ├─▶ Add JWT from httpOnly cookie
   └─▶ Forward to backend (/api/v1/*)

5. Backend
   └─▶ gRPC call to delivery microservice

6. Response flows back
   ├─▶ Redux updates state
   └─▶ Component re-renders
```

---

## Components

### UnifiedDeliverySelector

Универсальный компонент для выбора доставки, работает в трех контекстах: cart, checkout, admin.

**Location:** `frontend/svetu/src/components/delivery/UnifiedDeliverySelector.tsx`

#### Props

```typescript
interface Props {
  // Параметры расчета доставки
  calculationRequest: CalculationRequest;

  // Callback при выборе варианта доставки
  onQuoteSelected?: (quote: DeliveryQuote) => void;

  // ID выбранного quote (для highlight)
  selectedQuoteId?: number;

  // Автоматически рассчитывать при mount (default: true)
  autoCalculate?: boolean;

  // Показывать сравнение провайдеров (default: true)
  showComparison?: boolean;

  // CSS класс для стилизации
  className?: string;
}
```

#### CalculationRequest Type

```typescript
interface CalculationRequest {
  // Откуда (город отправителя)
  from_location: {
    city: string;
    postal_code: string;
  };

  // Куда (город получателя)
  to_location: {
    city: string;
    postal_code: string;
  };

  // Товары для доставки
  items: Array<{
    weight: number;      // kg
    length?: number;     // cm
    width?: number;      // cm
    height?: number;     // cm
    quantity: number;
  }>;

  // Опционально: конкретный провайдер
  provider_id?: string;

  // Опционально: страхование
  insurance_value?: number;  // RSD

  // Опционально: наложенный платеж
  cod_amount?: number;  // RSD
}
```

#### Usage Example (Cart)

```typescript
'use client';

import { useState } from 'react';
import { UnifiedDeliverySelector } from '@/components/delivery';
import { DeliveryQuote } from '@/types/delivery';
import { useAppDispatch } from '@/store/hooks';
import { selectQuote } from '@/store/slices/deliverySlice';

export default function CartPage() {
  const dispatch = useAppDispatch();
  const [selectedQuote, setSelectedQuote] = useState<DeliveryQuote | null>(null);

  // Prepare calculation request from cart items
  const calculationRequest = {
    from_location: {
      city: 'Belgrade',
      postal_code: '11000',
    },
    to_location: {
      city: 'Novi Sad',
      postal_code: '21000',
    },
    items: cartItems.map(item => ({
      weight: item.weight,
      length: item.length,
      width: item.width,
      height: item.height,
      quantity: item.quantity,
    })),
  };

  const handleQuoteSelect = (quote: DeliveryQuote) => {
    setSelectedQuote(quote);

    // Save to Redux for persistence between pages
    dispatch(selectQuote({
      storefrontId: 'storefront-123',
      quote,
    }));
  };

  return (
    <div className="container mx-auto p-4">
      <h1>Shopping Cart</h1>

      {/* Cart items list */}
      <CartItemsList items={cartItems} />

      {/* Delivery selection */}
      <div className="mt-8">
        <h2 className="text-xl font-bold mb-4">Select Delivery</h2>
        <UnifiedDeliverySelector
          calculationRequest={calculationRequest}
          onQuoteSelected={handleQuoteSelect}
          selectedQuoteId={selectedQuote?.provider_id}
          autoCalculate={true}
          showComparison={true}
        />
      </div>

      {/* Checkout button */}
      <button
        className="btn btn-primary mt-4"
        disabled={!selectedQuote}
        onClick={() => router.push('/checkout')}
      >
        Proceed to Checkout
      </button>
    </div>
  );
}
```

#### Usage Example (Checkout)

```typescript
'use client';

import { UnifiedDeliverySelector } from '@/components/delivery';
import { useAppSelector } from '@/store/hooks';
import { selectSelectedQuote } from '@/store/slices/deliverySlice';

export default function CheckoutPage() {
  // Get previously selected quote from Redux
  const selectedQuote = useAppSelector(selectSelectedQuote('storefront-123'));

  return (
    <div className="container mx-auto p-4">
      <h1>Checkout</h1>

      {/* Show selected delivery (read-only) */}
      {selectedQuote && (
        <div className="alert alert-info">
          <div>
            <strong>Delivery:</strong> {selectedQuote.provider_name}
          </div>
          <div>
            <strong>Cost:</strong> {selectedQuote.total_cost} RSD
          </div>
        </div>
      )}

      {/* Order summary */}
      <OrderSummary quote={selectedQuote} />

      {/* Payment form */}
      <PaymentForm />
    </div>
  );
}
```

#### UI States

**1. Loading State:**

```tsx
<div className="card bg-base-100 shadow-lg">
  <div className="card-body p-6 text-center">
    <ArrowPathIcon className="w-12 h-12 mx-auto text-primary animate-spin" />
    <h3 className="text-lg font-semibold">Calculating delivery rates...</h3>
    <p className="text-base-content/60">
      Comparing offers from all providers...
    </p>
  </div>
</div>
```

**2. Success State (Multiple Quotes):**

```tsx
<div className="space-y-4">
  {quotes.map((quote) => (
    <div
      key={quote.provider_id}
      className={`card cursor-pointer hover:shadow-xl transition ${
        selectedQuoteId === quote.provider_id
          ? 'ring-2 ring-primary'
          : ''
      }`}
      onClick={() => onQuoteSelected(quote)}
    >
      <div className="card-body">
        <div className="flex justify-between items-center">
          <div className="flex items-center gap-4">
            <img
              src={quote.logo_url}
              alt={quote.provider_name}
              className="w-16 h-16 object-contain"
            />
            <div>
              <h3 className="font-bold">{quote.provider_name}</h3>
              <p className="text-sm text-base-content/60">
                Delivery in {quote.estimated_delivery_days} days
              </p>
            </div>
          </div>
          <div className="text-right">
            <div className="text-2xl font-bold">
              {quote.total_cost} {quote.currency}
            </div>
            {quote.is_recommended && (
              <span className="badge badge-success">Recommended</span>
            )}
          </div>
        </div>
      </div>
    </div>
  ))}
</div>
```

**3. Error State:**

```tsx
<div className="card bg-base-100 shadow-lg">
  <div className="card-body">
    <div className="alert alert-error">
      <ExclamationTriangleIcon className="w-5 h-5" />
      <div>
        <div className="font-semibold">Calculation Error</div>
        <div className="text-sm">{error}</div>
      </div>
    </div>
    <button
      className="btn btn-outline mt-4"
      onClick={() => calculateRates()}
    >
      Retry
    </button>
  </div>
</div>
```

### CartDeliveryCalculator

Упрощенный компонент для быстрого расчета доставки в корзине.

**Location:** `frontend/svetu/src/components/delivery/CartDeliveryCalculator.tsx`

#### Props

```typescript
interface Props {
  storefrontId: number;
  onCostCalculated?: (cost: number) => void;
}
```

#### Usage

```typescript
import { CartDeliveryCalculator } from '@/components/delivery';

<CartDeliveryCalculator
  storefrontId={123}
  onCostCalculated={(cost) => {
    console.log('Delivery cost:', cost);
    setDeliveryCost(cost);
  }}
/>
```

### DeliveryInfo

Компонент для отображения tracking информации заказа.

**Location:** `frontend/svetu/src/components/tracking/DeliveryInfo.tsx`

#### Props

```typescript
interface Props {
  trackingNumber: string;
  autoRefresh?: boolean;  // Auto-refresh every 60s
  className?: string;
}
```

#### Usage

```typescript
import { DeliveryInfo } from '@/components/tracking';

<DeliveryInfo
  trackingNumber="PE1234567890RS"
  autoRefresh={true}
/>
```

#### Features

- Real-time status updates
- Event timeline
- ETA display
- Courier information (if available)
- Auto-refresh option

---

## Redux State Management

### deliverySlice

Централизованное управление состоянием delivery системы с кэшированием.

**Location:** `frontend/svetu/src/store/slices/deliverySlice.ts`

### State Structure

```typescript
interface DeliveryState {
  // Providers
  providers: DeliveryProvider[];
  providersLoading: boolean;
  providersError: string | null;

  // Selected quotes by storefront
  selectedQuotes: Record<string, DeliveryQuote>;

  // Calculation cache with TTL
  calculations: Record<string, CachedCalculation>;
  calculationsLoading: Record<string, boolean>;
  calculationsError: Record<string, string>;

  // Tracking info
  tracking: Record<string, TrackingInfo>;
  trackingLoading: Record<string, boolean>;
  trackingError: Record<string, string>;
}

interface CachedCalculation {
  data: CalculationResponse;
  timestamp: number;
  params: CalculationRequest;
}
```

### Actions

#### 1. fetchProviders

Загружает список доступных провайдеров (обычно 1 раз при загрузке приложения).

```typescript
import { useAppDispatch, useAppSelector } from '@/store/hooks';
import { fetchProviders, selectProviders } from '@/store/slices/deliverySlice';

// Component
const dispatch = useAppDispatch();
const providers = useAppSelector(selectProviders);

// Fetch providers on mount
useEffect(() => {
  dispatch(fetchProviders());
}, [dispatch]);

// Use providers
console.log('Available providers:', providers);
```

#### 2. calculateRate

Рассчитывает стоимость доставки с автоматическим кэшированием.

```typescript
import { calculateRate, selectCalculation } from '@/store/slices/deliverySlice';

const dispatch = useAppDispatch();

// Define request
const request: CalculationRequest = {
  from_location: { city: 'Belgrade', postal_code: '11000' },
  to_location: { city: 'Novi Sad', postal_code: '21000' },
  items: [{ weight: 2.5, quantity: 1 }],
};

// Check cache first
const cached = useAppSelector(selectCalculation(request));
if (cached) {
  console.log('Using cached calculation:', cached);
} else {
  // Fetch new calculation
  dispatch(calculateRate({ request }));
}
```

**Cache behavior:**
- TTL: 5 minutes
- Key: Hash of request params
- Automatic invalidation after TTL

#### 3. selectQuote

Сохраняет выбранный quote для конкретного storefront.

```typescript
import { selectQuote, selectSelectedQuote } from '@/store/slices/deliverySlice';

const dispatch = useAppDispatch();

// Save selected quote
dispatch(selectQuote({
  storefrontId: 'storefront-123',
  quote: {
    provider_id: 'post_express',
    provider_name: 'Post Express',
    total_cost: 500,
    estimated_delivery_days: 2,
    currency: 'RSD',
  },
}));

// Retrieve selected quote (on another page)
const selectedQuote = useAppSelector(selectSelectedQuote('storefront-123'));
console.log('Selected quote:', selectedQuote);
```

#### 4. trackShipment

Отслеживает статус shipment.

```typescript
import { trackShipment, selectTracking } from '@/store/slices/deliverySlice';

const dispatch = useAppDispatch();

// Track shipment
dispatch(trackShipment('PE1234567890RS'));

// Get tracking info
const trackingInfo = useAppSelector(selectTracking('PE1234567890RS'));
console.log('Tracking:', trackingInfo);
```

### Selectors

#### selectProviders

```typescript
const providers = useAppSelector(selectProviders);
// Returns: DeliveryProvider[]
```

#### selectProvidersLoading

```typescript
const loading = useAppSelector(selectProvidersLoading);
// Returns: boolean
```

#### selectSelectedQuote

```typescript
const quote = useAppSelector(selectSelectedQuote('storefront-123'));
// Returns: DeliveryQuote | undefined
```

#### selectAllSelectedQuotes

```typescript
const allQuotes = useAppSelector(selectAllSelectedQuotes);
// Returns: Record<string, DeliveryQuote>
```

#### selectCalculation

```typescript
const calculation = useAppSelector(selectCalculation(request));
// Returns: CalculationResponse | null
```

#### selectCalculationLoading

```typescript
const loading = useAppSelector(selectCalculationLoading(request));
// Returns: boolean
```

#### selectTracking

```typescript
const tracking = useAppSelector(selectTracking('PE1234567890RS'));
// Returns: TrackingInfo | undefined
```

### Cache Management

#### Clear calculations cache

```typescript
import { clearCalculationsCache } from '@/store/slices/deliverySlice';

dispatch(clearCalculationsCache());
```

#### Clear specific quote

```typescript
import { clearQuote } from '@/store/slices/deliverySlice';

dispatch(clearQuote('storefront-123'));
```

#### Clear all quotes

```typescript
import { clearAllQuotes } from '@/store/slices/deliverySlice';

dispatch(clearAllQuotes());
```

---

## API Service

### deliveryService

API wrapper для всех delivery endpoints через BFF proxy.

**Location:** `frontend/svetu/src/services/delivery.ts`

**ВАЖНО:** Всегда используй `deliveryService`, НЕ прямые `apiClient` calls!

### Methods

#### 1. calculateRate

```typescript
import { deliveryService } from '@/services/delivery';

const response = await deliveryService.calculateRate({
  provider: 'post_express',
  from_city: 'Belgrade',
  to_city: 'Novi Sad',
  weight: 2.5,
  cash_on_delivery: true,
  cod_amount: 5000,
});

if (response.data) {
  console.log('Total cost:', response.data.total_cost);
  console.log('Estimated delivery:', response.data.estimated_delivery);
}

if (response.error) {
  console.error('Error:', response.error.message);
}
```

#### 2. getProviders

```typescript
const response = await deliveryService.getProviders();

if (response.data) {
  response.data.forEach(provider => {
    console.log(provider.name, provider.code, provider.enabled);
  });
}
```

#### 3. createShipment

```typescript
const response = await deliveryService.createShipment({
  order_id: 123,
  provider_code: 'post_express',
  from_address: {
    name: 'Store Name',
    phone: '+381601234567',
    street: 'Main Street 1',
    city: 'Belgrade',
    postalCode: '11000',
    country: 'RS',
  },
  to_address: {
    name: 'Customer Name',
    phone: '+381607654321',
    street: 'Customer Street 5',
    city: 'Novi Sad',
    postalCode: '21000',
    country: 'RS',
  },
  packages: [
    {
      weight: 2.5,
      length: 30,
      width: 20,
      height: 10,
      description: 'Order #123',
    },
  ],
});

if (response.data) {
  console.log('Tracking number:', response.data.tracking_number);
  console.log('Label URL:', response.data.label_url);
}
```

#### 4. trackShipment

```typescript
const response = await deliveryService.trackShipment('PE1234567890RS');

if (response.data) {
  console.log('Status:', response.data.status);
  console.log('Current location:', response.data.current_location);
  console.log('Events:', response.data.events);
}
```

#### 5. cancelShipment

```typescript
const response = await deliveryService.cancelShipment(123);

if (!response.error) {
  console.log('Shipment cancelled successfully');
}
```

### Error Handling

Все методы возвращают `ApiResponse<T>`:

```typescript
type ApiResponse<T> = {
  data?: T;
  error?: {
    code: string;
    message: string;
    details?: any;
  };
};
```

**Example:**

```typescript
const response = await deliveryService.calculateRate(request);

if (response.error) {
  // Handle error
  switch (response.error.code) {
    case 'DELIVERY_SERVICE_UNAVAILABLE':
      showWarning('Delivery service temporarily unavailable');
      break;
    case 'INVALID_DELIVERY_ADDRESS':
      showError('Please check delivery address');
      break;
    default:
      showError(response.error.message);
  }
  return;
}

// Use data
const { total_cost, estimated_delivery } = response.data!;
```

---

## Usage Examples

### Example 1: Cart Page with Delivery Selection

```typescript
'use client';

import { useState, useEffect } from 'react';
import { useAppDispatch, useAppSelector } from '@/store/hooks';
import { UnifiedDeliverySelector } from '@/components/delivery';
import {
  selectQuote,
  selectSelectedQuote,
  fetchProviders,
} from '@/store/slices/deliverySlice';

export default function CartPage() {
  const dispatch = useAppDispatch();
  const selectedQuote = useAppSelector(selectSelectedQuote('cart'));
  const [cartItems, setCartItems] = useState([]);

  // Load providers on mount
  useEffect(() => {
    dispatch(fetchProviders());
  }, [dispatch]);

  // Prepare calculation request
  const calculationRequest = {
    from_location: {
      city: 'Belgrade',
      postal_code: '11000',
    },
    to_location: {
      city: 'Novi Sad',
      postal_code: '21000',
    },
    items: cartItems.map(item => ({
      weight: item.weight || 1.0,
      quantity: item.quantity,
    })),
  };

  const handleQuoteSelect = (quote) => {
    dispatch(selectQuote({ storefrontId: 'cart', quote }));
  };

  const totalWithDelivery =
    cartTotal + (selectedQuote?.total_cost || 0);

  return (
    <div className="container mx-auto p-4">
      <h1 className="text-2xl font-bold mb-6">Shopping Cart</h1>

      {/* Cart items */}
      <div className="grid grid-cols-1 gap-4 mb-8">
        {cartItems.map(item => (
          <CartItem key={item.id} item={item} />
        ))}
      </div>

      {/* Delivery selection */}
      {cartItems.length > 0 && (
        <div className="mb-8">
          <h2 className="text-xl font-bold mb-4">Select Delivery</h2>
          <UnifiedDeliverySelector
            calculationRequest={calculationRequest}
            onQuoteSelected={handleQuoteSelect}
            selectedQuoteId={selectedQuote?.provider_id}
            autoCalculate={true}
            showComparison={true}
          />
        </div>
      )}

      {/* Summary */}
      <div className="card bg-base-200">
        <div className="card-body">
          <div className="flex justify-between">
            <span>Subtotal:</span>
            <span>{cartTotal} RSD</span>
          </div>
          {selectedQuote && (
            <div className="flex justify-between">
              <span>Delivery ({selectedQuote.provider_name}):</span>
              <span>{selectedQuote.total_cost} RSD</span>
            </div>
          )}
          <div className="divider" />
          <div className="flex justify-between text-xl font-bold">
            <span>Total:</span>
            <span>{totalWithDelivery} RSD</span>
          </div>
        </div>
      </div>

      {/* Checkout button */}
      <button
        className="btn btn-primary btn-block mt-4"
        disabled={!selectedQuote}
        onClick={() => router.push('/checkout')}
      >
        Proceed to Checkout
      </button>
    </div>
  );
}
```

### Example 2: Order Tracking Page

```typescript
'use client';

import { useEffect } from 'react';
import { useParams } from 'next/navigation';
import { useAppDispatch, useAppSelector } from '@/store/hooks';
import { trackShipment, selectTracking, selectTrackingLoading } from '@/store/slices/deliverySlice';
import { DeliveryInfo } from '@/components/tracking';

export default function OrderTrackingPage() {
  const params = useParams();
  const trackingNumber = params.trackingNumber as string;

  const dispatch = useAppDispatch();
  const trackingInfo = useAppSelector(selectTracking(trackingNumber));
  const loading = useAppSelector(selectTrackingLoading(trackingNumber));

  // Load tracking on mount
  useEffect(() => {
    if (trackingNumber) {
      dispatch(trackShipment(trackingNumber));
    }
  }, [dispatch, trackingNumber]);

  // Auto-refresh every 60 seconds
  useEffect(() => {
    if (!trackingNumber) return;

    const interval = setInterval(() => {
      dispatch(trackShipment(trackingNumber));
    }, 60000);

    return () => clearInterval(interval);
  }, [dispatch, trackingNumber]);

  if (loading && !trackingInfo) {
    return <div className="loading loading-spinner loading-lg" />;
  }

  if (!trackingInfo) {
    return (
      <div className="alert alert-error">
        <span>Tracking information not found</span>
      </div>
    );
  }

  return (
    <div className="container mx-auto p-4">
      <h1 className="text-2xl font-bold mb-6">
        Track Order: {trackingNumber}
      </h1>

      <DeliveryInfo
        trackingNumber={trackingNumber}
        autoRefresh={true}
      />

      {/* Additional order details */}
      <div className="mt-8 card bg-base-200">
        <div className="card-body">
          <h2 className="card-title">Shipment Details</h2>
          <div className="grid grid-cols-2 gap-4">
            <div>
              <strong>Status:</strong> {trackingInfo.status}
            </div>
            <div>
              <strong>Location:</strong> {trackingInfo.current_location}
            </div>
            <div>
              <strong>Estimated Delivery:</strong>{' '}
              {new Date(trackingInfo.estimated_delivery).toLocaleDateString()}
            </div>
            {trackingInfo.actual_delivery && (
              <div>
                <strong>Delivered:</strong>{' '}
                {new Date(trackingInfo.actual_delivery).toLocaleDateString()}
              </div>
            )}
          </div>
        </div>
      </div>

      {/* Event timeline */}
      <div className="mt-8">
        <h2 className="text-xl font-bold mb-4">Tracking History</h2>
        <div className="timeline timeline-vertical">
          {trackingInfo.events.map((event, index) => (
            <div key={index} className="timeline-item">
              <div className="timeline-start">
                {new Date(event.timestamp).toLocaleString()}
              </div>
              <div className="timeline-middle">
                <div className="w-4 h-4 rounded-full bg-primary" />
              </div>
              <div className="timeline-end timeline-box">
                <div className="font-bold">{event.location}</div>
                <div className="text-sm">{event.description}</div>
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}
```

### Example 3: Admin Shipment Management

```typescript
'use client';

import { useState, useEffect } from 'react';
import { deliveryService } from '@/services/delivery';
import { apiClient } from '@/services/api-client';

export default function AdminShipmentsPage() {
  const [shipments, setShipments] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadShipments();
  }, []);

  const loadShipments = async () => {
    setLoading(true);
    const response = await apiClient.get('/admin/delivery/shipments');
    if (response.data) {
      setShipments(response.data);
    }
    setLoading(false);
  };

  const handleCancelShipment = async (shipmentId: number) => {
    if (!confirm('Are you sure you want to cancel this shipment?')) {
      return;
    }

    const response = await deliveryService.cancelShipment(shipmentId);

    if (response.error) {
      alert('Failed to cancel: ' + response.error.message);
      return;
    }

    alert('Shipment cancelled successfully');
    loadShipments();
  };

  const handleRetryShipment = async (orderId: number) => {
    // Retry creating shipment for failed order
    const response = await apiClient.post(
      `/admin/orders/${orderId}/retry-shipment`
    );

    if (response.error) {
      alert('Failed to retry: ' + response.error.message);
      return;
    }

    alert('Shipment created successfully');
    loadShipments();
  };

  if (loading) {
    return <div className="loading loading-spinner loading-lg" />;
  }

  return (
    <div className="container mx-auto p-4">
      <h1 className="text-2xl font-bold mb-6">Shipment Management</h1>

      <div className="overflow-x-auto">
        <table className="table table-zebra">
          <thead>
            <tr>
              <th>ID</th>
              <th>Order ID</th>
              <th>Tracking Number</th>
              <th>Provider</th>
              <th>Status</th>
              <th>Created</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            {shipments.map(shipment => (
              <tr key={shipment.id}>
                <td>{shipment.id}</td>
                <td>{shipment.order_id}</td>
                <td>
                  <a
                    href={`/admin/tracking/${shipment.tracking_number}`}
                    className="link"
                  >
                    {shipment.tracking_number}
                  </a>
                </td>
                <td>{shipment.provider_code}</td>
                <td>
                  <span className={`badge ${getStatusBadgeClass(shipment.status)}`}>
                    {shipment.status}
                  </span>
                </td>
                <td>{new Date(shipment.created_at).toLocaleDateString()}</td>
                <td>
                  <div className="flex gap-2">
                    {shipment.label_url && (
                      <a
                        href={shipment.label_url}
                        target="_blank"
                        className="btn btn-sm btn-ghost"
                      >
                        Label
                      </a>
                    )}
                    {shipment.status === 'PENDING' && (
                      <button
                        className="btn btn-sm btn-error"
                        onClick={() => handleCancelShipment(shipment.id)}
                      >
                        Cancel
                      </button>
                    )}
                  </div>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
```

---

## Best Practices

### 1. Always use BFF proxy (apiClient)

```typescript
// ✅ CORRECT
import { deliveryService } from '@/services/delivery';
const response = await deliveryService.calculateRate(request);

// ❌ WRONG - Never do direct fetch!
const response = await fetch('http://localhost:3000/api/v1/delivery/...');
```

### 2. Use Redux for shared state

```typescript
// ✅ CORRECT - Share quote between pages
dispatch(selectQuote({ storefrontId: 'cart', quote }));

// Later on checkout page:
const quote = useAppSelector(selectSelectedQuote('cart'));

// ❌ WRONG - Local state lost on navigation
const [quote, setQuote] = useState(null);
```

### 3. Leverage caching

```typescript
// ✅ CORRECT - Check cache first
const cached = useAppSelector(selectCalculation(request));
if (cached) {
  return cached;
}

dispatch(calculateRate({ request }));

// ❌ WRONG - Always fetch (wasteful)
dispatch(calculateRate({ request }));
```

### 4. Handle errors gracefully

```typescript
// ✅ CORRECT - Show fallback UI
if (error) {
  return (
    <div className="alert alert-warning">
      <span>Delivery calculation unavailable. You can proceed without delivery.</span>
      <button onClick={retry}>Retry</button>
    </div>
  );
}

// ❌ WRONG - Block entire checkout
if (error) {
  throw new Error('Cannot proceed without delivery');
}
```

### 5. Performance optimization

```typescript
// ✅ CORRECT - Lazy load heavy components
import dynamic from 'next/dynamic';

const UnifiedDeliverySelector = dynamic(
  () => import('@/components/delivery/UnifiedDeliverySelector'),
  { ssr: false, loading: () => <Skeleton /> }
);

// ✅ CORRECT - Prefetch providers on homepage
useEffect(() => {
  dispatch(fetchProviders()); // Load once, cache forever
}, []);

// ✅ CORRECT - Debounce address input
const debouncedCalculate = useDebounce(() => {
  dispatch(calculateRate({ request }));
}, 500);
```

### 6. Type safety

```typescript
// ✅ CORRECT - Use TypeScript types
import type { DeliveryQuote, CalculationRequest } from '@/types/delivery';

const handleSelect = (quote: DeliveryQuote) => {
  // TypeScript knows quote structure
  console.log(quote.total_cost);
};

// ❌ WRONG - Untyped
const handleSelect = (quote: any) => {
  // No type safety
};
```

### 7. Loading states

```typescript
// ✅ CORRECT - Show loading indicator
const loading = useAppSelector(selectCalculationLoading(request));

if (loading) {
  return <LoadingSpinner />;
}

// ❌ WRONG - No feedback to user
// Just wait silently...
```

### 8. Clean up on unmount

```typescript
// ✅ CORRECT - Clear cache when leaving cart
useEffect(() => {
  return () => {
    dispatch(clearQuote('cart'));
  };
}, []);
```

---

## Troubleshooting

### Issue 1: "Calculation not working"

**Symptoms:**
- UnifiedDeliverySelector показывает loading бесконечно
- No quotes appear

**Solutions:**

1. Check request params:
```typescript
console.log('Request:', calculationRequest);
// Verify from_location, to_location, items are correct
```

2. Check Redux state:
```typescript
// Redux DevTools → delivery slice
// Look for calculationsError
```

3. Check network:
```typescript
// Browser DevTools → Network tab
// Find POST /api/v2/delivery/calculate-universal
// Check request/response
```

4. Check backend logs:
```bash
# Should see:
# [INFO] Calculating delivery rate: provider=post_express
```

### Issue 2: "Selected quote not persisting"

**Symptoms:**
- Quote selected in cart, but lost on checkout page

**Cause:** Not using Redux

**Solution:**
```typescript
// Cart page - Save to Redux
dispatch(selectQuote({ storefrontId: 'cart', quote }));

// Checkout page - Retrieve from Redux
const quote = useAppSelector(selectSelectedQuote('cart'));
```

### Issue 3: "Cache not working"

**Symptoms:**
- Same calculation making multiple API calls

**Cause:** Request params slightly different (deep comparison)

**Debug:**
```typescript
import { getCacheKey } from '@/store/slices/deliverySlice';

const key1 = getCacheKey(request1);
const key2 = getCacheKey(request2);

console.log('Same cache key?', key1 === key2);
```

**Solution:**
- Ensure consistent request structure
- Don't include timestamp or random fields

### Issue 4: "BFF proxy 401 Unauthorized"

**Symptoms:**
- Frontend shows "Authentication required"
- Backend returns 401

**Cause:** JWT token missing or expired

**Solutions:**

1. Check cookie:
```typescript
// Browser DevTools → Application → Cookies
// Look for 'access_token'
```

2. Re-login:
```typescript
// Force re-authentication
window.location.href = '/login';
```

3. Check BFF proxy logs:
```
[BFF Proxy] Has access_token: false  ← Problem!
```

### Issue 5: "Tracking not updating"

**Symptoms:**
- DeliveryInfo shows old status

**Solutions:**

1. Manual refresh:
```typescript
dispatch(trackShipment(trackingNumber));
```

2. Enable auto-refresh:
```typescript
<DeliveryInfo
  trackingNumber={trackingNumber}
  autoRefresh={true}  // Refresh every 60s
/>
```

3. Clear tracking cache:
```typescript
dispatch(clearTrackingCache());
```

---

## Performance Optimization

### Bundle Size

**Lazy load components:**

```typescript
import dynamic from 'next/dynamic';

const UnifiedDeliverySelector = dynamic(
  () => import('@/components/delivery/UnifiedDeliverySelector'),
  {
    ssr: false,
    loading: () => <DeliverySkeleton />,
  }
);
```

### API Calls

**Minimize redundant calls:**

```typescript
// ✅ Good - Fetch providers once on app load
useEffect(() => {
  dispatch(fetchProviders());
}, []); // Empty deps = run once

// ❌ Bad - Fetch on every render
dispatch(fetchProviders()); // NO deps = infinite loop!
```

### Cache Strategy

**5-minute TTL is configurable:**

```typescript
// In deliverySlice.ts
const CACHE_TTL = 5 * 60 * 1000; // 5 minutes

// Customize per use case:
// Cart: 5 min (price may change)
// Checkout: No cache (always fresh)
// Admin: 1 hour (rarely changes)
```

---

## Related Documentation

- [Delivery Service Integration Guide](./DELIVERY_SERVICE_INTEGRATION.md) - Backend integration architecture
- [Delivery Microservice API Reference](./DELIVERY_MICROSERVICE_API.md) - gRPC API documentation

---

**Last updated:** 2025-10-29
**Version:** 1.0.0
