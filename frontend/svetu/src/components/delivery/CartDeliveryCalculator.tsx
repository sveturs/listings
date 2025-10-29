'use client';

import { useState, useEffect } from 'react';
import {
  ShoppingCartIcon,
  MapPinIcon,
  TruckIcon,
  CurrencyDollarIcon,
  ExclamationTriangleIcon,
  InformationCircleIcon,
  CheckIcon,
  ClockIcon,
} from '@heroicons/react/24/outline';
import {
  DeliveryQuote,
  CalculationRequest,
  DeliveryAttributes,
} from '@/types/delivery';
import UnifiedDeliverySelector from './UnifiedDeliverySelector';
import { deliveryService } from '@/services/delivery';

interface CartItem {
  id: number;
  product_id: number;
  product_type: 'listing' | 'storefront_product';
  name: string;
  price: number;
  quantity: number;
  seller_id: number;
  seller_name: string;
  attributes?: DeliveryAttributes;
  category_id?: number;
}

interface Props {
  cartItems: CartItem[];
  deliveryAddress?: {
    city: string;
    postal_code?: string;
    address?: string;
  };
  onDeliverySelected?: (
    quote: DeliveryQuote,
    splitByProvider: { [providerId: number]: CartItem[] }
  ) => void;
  className?: string;
}

export default function CartDeliveryCalculator({
  cartItems,
  deliveryAddress,
  onDeliverySelected,
  className = '',
}: Props) {
  const [selectedQuote, setSelectedQuote] = useState<DeliveryQuote | null>(
    null
  );
  const [loading, setLoading] = useState(false);
  const [itemsWithAttributes, setItemsWithAttributes] = useState<CartItem[]>(
    []
  );
  const [_splitByProvider, setSplitByProvider] = useState<{
    [providerId: number]: CartItem[];
  }>({});
  const [showAddressForm, setShowAddressForm] = useState(!deliveryAddress);
  const [addressForm, setAddressForm] = useState({
    city: deliveryAddress?.city || '',
    postal_code: deliveryAddress?.postal_code || '',
    address: deliveryAddress?.address || '',
  });

  useEffect(() => {
    if (cartItems.length > 0) {
      loadItemAttributes();
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [cartItems]);

  const loadItemAttributes = async () => {
    setLoading(true);

    try {
      const itemsWithAttrs = await Promise.all(
        cartItems.map(async (item) => {
          if (item.attributes) {
            return item; // Already has attributes
          }

          try {
            const response = await deliveryService.getProductAttributes(
              item.product_id.toString(),
              item.product_type
            );

            return {
              ...item,
              attributes: response.data || await getDefaultAttributes(item.category_id),
            };
          } catch (error) {
            console.error(
              `Failed to load attributes for product ${item.product_id}:`,
              error
            );
            return {
              ...item,
              attributes: await getDefaultAttributes(item.category_id),
            };
          }
        })
      );

      setItemsWithAttributes(itemsWithAttrs);
    } catch (error) {
      console.error('Failed to load item attributes:', error);
      // Use items as is, with default attributes where needed
      setItemsWithAttributes(
        cartItems.map((item) => ({
          ...item,
          attributes: item.attributes || {
            weight_kg: 1,
            dimensions: { length_cm: 30, width_cm: 20, height_cm: 10 },
            is_fragile: false,
            requires_special_handling: false,
            stackable: true,
            packaging_type: 'box' as const,
          },
        }))
      );
    } finally {
      setLoading(false);
    }
  };

  const getDefaultAttributes = async (
    categoryId?: number
  ): Promise<DeliveryAttributes> => {
    if (!categoryId) {
      return {
        weight_kg: 1,
        dimensions: { length_cm: 30, width_cm: 20, height_cm: 10 },
        is_fragile: false,
        requires_special_handling: false,
        stackable: true,
        packaging_type: 'box',
      };
    }

    try {
      const response = await deliveryService.getCategoryDefaults(categoryId.toString());

      if (response.data) {
        const defaults = response.data;
        return {
          weight_kg: defaults.default_weight_kg || 1,
          dimensions: {
            length_cm: defaults.default_length_cm || 30,
            width_cm: defaults.default_width_cm || 20,
            height_cm: defaults.default_height_cm || 10,
          },
          is_fragile: defaults.is_typically_fragile || false,
          requires_special_handling: false,
          stackable: !defaults.is_typically_fragile,
          packaging_type: (defaults.default_packaging_type as any) || 'box',
        };
      }
    } catch (error) {
      console.error('Failed to load category defaults:', error);
    }

    // Fallback defaults
    return {
      weight_kg: 1,
      dimensions: { length_cm: 30, width_cm: 20, height_cm: 10 },
      is_fragile: false,
      requires_special_handling: false,
      stackable: true,
      packaging_type: 'box',
    };
  };

  const buildCalculationRequest = (): CalculationRequest | null => {
    if (!addressForm.city || itemsWithAttributes.length === 0) {
      return null;
    }

    return {
      from_location: {
        city: 'Нови Сад', // Default sender location
        postal_code: '21000',
      },
      to_location: {
        city: addressForm.city,
        postal_code: addressForm.postal_code,
        address: addressForm.address,
      },
      items: itemsWithAttributes.map((item) => ({
        product_id: item.product_id,
        product_type: item.product_type,
        quantity: item.quantity,
        attributes: item.attributes,
      })),
    };
  };

  const handleQuoteSelected = (quote: DeliveryQuote) => {
    setSelectedQuote(quote);

    // Split items by provider (for multi-vendor scenarios)
    const split = { [quote.provider_id]: itemsWithAttributes };
    setSplitByProvider(split);

    onDeliverySelected?.(quote, split);
  };

  const handleAddressSubmit = () => {
    if (addressForm.city) {
      setShowAddressForm(false);
    }
  };

  const getTotalWeight = () => {
    return itemsWithAttributes.reduce((total, item) => {
      const weight = item.attributes?.weight_kg || 1;
      return total + weight * item.quantity;
    }, 0);
  };

  const getTotalValue = () => {
    return cartItems.reduce((total, item) => {
      return total + item.price * item.quantity;
    }, 0);
  };

  const getSellerCount = () => {
    return new Set(cartItems.map((item) => item.seller_id)).size;
  };

  const hasFragileItems = () => {
    return itemsWithAttributes.some((item) => item.attributes?.is_fragile);
  };

  if (cartItems.length === 0) {
    return (
      <div className={`card bg-base-100 shadow-lg ${className}`}>
        <div className="card-body p-6 text-center">
          <ShoppingCartIcon className="w-16 h-16 mx-auto text-base-content/30 mb-4" />
          <h3 className="text-lg font-semibold mb-2">Корзина пуста</h3>
          <p className="text-base-content/60">
            Добавьте товары в корзину для расчета доставки
          </p>
        </div>
      </div>
    );
  }

  return (
    <div className={`space-y-6 ${className}`}>
      {/* Cart Summary */}
      <div className="card bg-base-100 shadow-lg">
        <div className="card-body p-6">
          <div className="flex items-center gap-3 mb-4">
            <div className="p-2 bg-primary/10 rounded-lg">
              <ShoppingCartIcon className="w-6 h-6 text-primary" />
            </div>
            <div>
              <h3 className="text-lg font-semibold">Ваш заказ</h3>
              <p className="text-sm text-base-content/60">
                {cartItems.length} товаров от {getSellerCount()} продавцов
              </p>
            </div>
          </div>

          <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-4">
            <div className="text-center">
              <div className="text-2xl font-bold text-primary">
                {cartItems.length}
              </div>
              <div className="text-xs text-base-content/60">товаров</div>
            </div>
            <div className="text-center">
              <div className="text-2xl font-bold">
                {getTotalWeight().toFixed(1)} кг
              </div>
              <div className="text-xs text-base-content/60">общий вес</div>
            </div>
            <div className="text-center">
              <div className="text-2xl font-bold">
                {getTotalValue().toFixed(0)} RSD
              </div>
              <div className="text-xs text-base-content/60">стоимость</div>
            </div>
            <div className="text-center">
              <div className="text-2xl font-bold">{getSellerCount()}</div>
              <div className="text-xs text-base-content/60">продавцов</div>
            </div>
          </div>

          {/* Special notices */}
          {hasFragileItems() && (
            <div className="alert alert-warning alert-sm">
              <ExclamationTriangleIcon className="w-4 h-4" />
              <span className="text-sm">В заказе есть хрупкие товары</span>
            </div>
          )}

          {getSellerCount() > 1 && (
            <div className="alert alert-info alert-sm">
              <InformationCircleIcon className="w-4 h-4" />
              <span className="text-sm">
                Товары от разных продавцов могут доставляться отдельно
              </span>
            </div>
          )}
        </div>
      </div>

      {/* Delivery Address */}
      <div className="card bg-base-100 shadow-lg">
        <div className="card-body p-6">
          <div className="flex items-center justify-between mb-4">
            <div className="flex items-center gap-3">
              <div className="p-2 bg-primary/10 rounded-lg">
                <MapPinIcon className="w-6 h-6 text-primary" />
              </div>
              <div>
                <h3 className="text-lg font-semibold">Адрес доставки</h3>
                <p className="text-sm text-base-content/60">
                  Укажите куда доставить заказ
                </p>
              </div>
            </div>
            {!showAddressForm && addressForm.city && (
              <button
                className="btn btn-sm btn-ghost"
                onClick={() => setShowAddressForm(true)}
              >
                Изменить
              </button>
            )}
          </div>

          {showAddressForm ? (
            <div className="space-y-4">
              <div className="grid md:grid-cols-2 gap-4">
                <div className="form-control">
                  <label className="label">
                    <span className="label-text font-medium">Город *</span>
                  </label>
                  <input
                    type="text"
                    className="input input-bordered focus:input-primary"
                    placeholder="Белград, Ниш, Крагуевац..."
                    value={addressForm.city}
                    onChange={(e) =>
                      setAddressForm((prev) => ({
                        ...prev,
                        city: e.target.value,
                      }))
                    }
                  />
                </div>
                <div className="form-control">
                  <label className="label">
                    <span className="label-text">Почтовый индекс</span>
                  </label>
                  <input
                    type="text"
                    className="input input-bordered focus:input-primary"
                    placeholder="11000"
                    value={addressForm.postal_code}
                    onChange={(e) =>
                      setAddressForm((prev) => ({
                        ...prev,
                        postal_code: e.target.value,
                      }))
                    }
                  />
                </div>
              </div>
              <div className="form-control">
                <label className="label">
                  <span className="label-text">Адрес (улица, дом)</span>
                </label>
                <input
                  type="text"
                  className="input input-bordered focus:input-primary"
                  placeholder="Опционально, для более точного расчета"
                  value={addressForm.address}
                  onChange={(e) =>
                    setAddressForm((prev) => ({
                      ...prev,
                      address: e.target.value,
                    }))
                  }
                />
              </div>
              <button
                className="btn btn-primary"
                onClick={handleAddressSubmit}
                disabled={!addressForm.city}
              >
                <CheckIcon className="w-4 h-4" />
                Подтвердить адрес
              </button>
            </div>
          ) : (
            <div className="flex items-center gap-2 text-base-content/70">
              <MapPinIcon className="w-4 h-4" />
              <span>
                {addressForm.city}
                {addressForm.postal_code && `, ${addressForm.postal_code}`}
                {addressForm.address && `, ${addressForm.address}`}
              </span>
            </div>
          )}
        </div>
      </div>

      {/* Delivery Options */}
      {!showAddressForm && addressForm.city && !loading && (
        <UnifiedDeliverySelector
          calculationRequest={buildCalculationRequest()!}
          onQuoteSelected={handleQuoteSelected}
          selectedQuoteId={selectedQuote?.provider_id}
          autoCalculate={true}
          showComparison={true}
        />
      )}

      {/* Selected Delivery Summary */}
      {selectedQuote && (
        <div className="card bg-gradient-to-r from-success/5 to-success/10 border border-success/20">
          <div className="card-body p-6">
            <div className="flex items-center gap-3 mb-4">
              <div className="p-2 bg-success/20 rounded-lg">
                <TruckIcon className="w-6 h-6 text-success" />
              </div>
              <div>
                <h3 className="text-lg font-semibold text-success">
                  Выбранная доставка
                </h3>
                <p className="text-sm text-base-content/60">
                  Готово к оформлению заказа
                </p>
              </div>
            </div>

            <div className="grid md:grid-cols-3 gap-4">
              <div className="text-center">
                <div className="text-xl font-bold">
                  {selectedQuote.provider_name}
                </div>
                <div className="text-xs text-base-content/60">провайдер</div>
              </div>
              <div className="text-center">
                <div className="text-xl font-bold flex items-center justify-center gap-1">
                  <ClockIcon className="w-4 h-4" />
                  {selectedQuote.estimated_days} дней
                </div>
                <div className="text-xs text-base-content/60">
                  время доставки
                </div>
              </div>
              <div className="text-center">
                <div className="text-xl font-bold text-success flex items-center justify-center gap-1">
                  <CurrencyDollarIcon className="w-4 h-4" />
                  {selectedQuote.total_price.toFixed(0)} RSD
                </div>
                <div className="text-xs text-base-content/60">
                  стоимость доставки
                </div>
              </div>
            </div>
          </div>
        </div>
      )}

      {loading && (
        <div className="card bg-base-100 shadow-lg">
          <div className="card-body p-6 text-center">
            <span className="loading loading-spinner loading-lg"></span>
            <p className="mt-4">Загрузка параметров товаров...</p>
          </div>
        </div>
      )}
    </div>
  );
}
