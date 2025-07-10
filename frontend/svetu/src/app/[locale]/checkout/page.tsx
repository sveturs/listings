'use client';

import React, { useState, useEffect } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';
import { useSelector, useDispatch } from 'react-redux';
import { useTranslations } from 'next-intl';
import {
  fetchCart,
  selectCart,
  selectCartLoading,
} from '@/store/slices/cartSlice';
import { ordersService } from '@/services/orders';
import SafeImage from '@/components/SafeImage';
import { useAuth } from '@/contexts/AuthContext';
import { tokenManager } from '@/utils/tokenManager';
import type { AppDispatch } from '@/store';
import type { components } from '@/types/generated/api';

type CreateOrderRequest =
  components['schemas']['backend_internal_domain_models.CreateOrderRequest'];
type ShippingAddress =
  components['schemas']['backend_internal_domain_models.ShippingAddress'];

export default function CheckoutPage() {
  const t = useTranslations('checkout');
  const router = useRouter();
  const searchParams = useSearchParams();
  const dispatch = useDispatch<AppDispatch>();
  const { user } = useAuth();

  const cart = useSelector(selectCart);
  const cartLoading = useSelector(selectCartLoading);

  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const [shippingAddress, setShippingAddress] = useState<ShippingAddress>({
    full_name: '',
    street: '',
    house_number: '',
    apartment: '',
    city: '',
    postal_code: '',
    country: '',
    phone: '',
    email: '',
    notes: '',
  });

  const [billingAddress, setBillingAddress] = useState<ShippingAddress>({
    full_name: '',
    street: '',
    house_number: '',
    apartment: '',
    city: '',
    postal_code: '',
    country: '',
    phone: '',
    email: '',
    notes: '',
  });

  const [useSameAddress, setUseSameAddress] = useState(true);
  const [paymentMethod, setPaymentMethod] = useState('allsecure');
  const [customerNotes, setCustomerNotes] = useState('');

  const storefrontId = searchParams.get('storefront');

  useEffect(() => {
    if (storefrontId) {
      dispatch(fetchCart(Number(storefrontId)));
    }
  }, [dispatch, storefrontId]);

  // Проверка авторизации
  useEffect(() => {
    if (!user) {
      // Сохраняем URL для возврата после входа
      const returnUrl = `/checkout?storefront=${storefrontId}`;
      router.push(`/?returnUrl=${encodeURIComponent(returnUrl)}`);
    }
  }, [user, router, storefrontId]);

  const calculateTotals = () => {
    if (!cart?.items) return { subtotal: 0, shipping: 0, tax: 0, total: 0 };

    const subtotal = cart.items.reduce(
      (sum: number, item: any) => sum + Number(item.total_price || 0),
      0
    );
    const shipping = 5.0; // Fixed shipping for now
    const tax = subtotal * 0.1; // 10% tax
    const total = subtotal + shipping + tax;

    return { subtotal, shipping, tax, total };
  };

  const { subtotal, shipping, tax, total } = calculateTotals();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!cart?.items?.length || !storefrontId) return;

    console.log('[CheckoutPage] Submitting order, user:', user);
    console.log(
      '[CheckoutPage] Token exists:',
      !!tokenManager.getAccessToken()
    );

    setLoading(true);
    setError(null);

    try {
      const orderData: CreateOrderRequest = {
        storefront_id: Number(storefrontId),
        items: cart.items.map((item: any) => ({
          product_id: item.product_id!,
          quantity: item.quantity!,
          price_per_unit: item.price_per_unit!,
          variant_id: item.variant_id,
        })),
        shipping_address: shippingAddress,
        billing_address: useSameAddress ? shippingAddress : billingAddress,
        payment_method: paymentMethod,
        customer_notes: customerNotes || undefined,
        shipping_method: 'standard',
      };

      const order = await ordersService.createOrder(orderData);

      // Redirect to payment
      router.push(`/payment/process/${order.id}`);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to create order');
    } finally {
      setLoading(false);
    }
  };

  if (cartLoading) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="flex justify-center">
          <span className="loading loading-spinner loading-lg"></span>
        </div>
      </div>
    );
  }

  // Не показываем страницу если пользователь не авторизован
  if (!user) {
    return null;
  }

  if (!cart?.items?.length) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="text-center">
          <h1 className="mb-4 text-2xl font-bold">{t('title')}</h1>
          <p className="mb-8 text-base-content/70">{t('emptyCart')}</p>
          <button onClick={() => router.back()} className="btn btn-primary">
            {t('continueShopping')}
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="container mx-auto px-4 py-8">
      <h1 className="mb-8 text-3xl font-bold">{t('title')}</h1>

      <div className="grid gap-8 lg:grid-cols-2">
        {/* Order Form */}
        <div>
          <form onSubmit={handleSubmit} className="space-y-6">
            {/* Shipping Address */}
            <div className="card bg-base-100 shadow">
              <div className="card-body">
                <h2 className="card-title">{t('shippingAddress')}</h2>

                <div className="grid gap-4 md:grid-cols-2">
                  <div className="md:col-span-2">
                    <label className="label">
                      <span className="label-text">{t('street')}</span>
                    </label>
                    <input
                      type="text"
                      className="input input-bordered w-full"
                      value={shippingAddress.street}
                      onChange={(e) =>
                        setShippingAddress({
                          ...shippingAddress,
                          street: e.target.value,
                        })
                      }
                      required
                    />
                  </div>

                  <div>
                    <label className="label">
                      <span className="label-text">{t('city')}</span>
                    </label>
                    <input
                      type="text"
                      className="input input-bordered w-full"
                      value={shippingAddress.city}
                      onChange={(e) =>
                        setShippingAddress({
                          ...shippingAddress,
                          city: e.target.value,
                        })
                      }
                      required
                    />
                  </div>

                  <div>
                    <label className="label">
                      <span className="label-text">{t('postalCode')}</span>
                    </label>
                    <input
                      type="text"
                      className="input input-bordered w-full"
                      value={shippingAddress.postal_code}
                      onChange={(e) =>
                        setShippingAddress({
                          ...shippingAddress,
                          postal_code: e.target.value,
                        })
                      }
                      required
                    />
                  </div>

                  <div>
                    <label className="label">
                      <span className="label-text">{t('country')}</span>
                    </label>
                    <input
                      type="text"
                      className="input input-bordered w-full"
                      value={shippingAddress.country}
                      onChange={(e) =>
                        setShippingAddress({
                          ...shippingAddress,
                          country: e.target.value,
                        })
                      }
                      required
                    />
                  </div>

                  <div>
                    <label className="label">
                      <span className="label-text">{t('phone')}</span>
                    </label>
                    <input
                      type="tel"
                      className="input input-bordered w-full"
                      value={shippingAddress.phone}
                      onChange={(e) =>
                        setShippingAddress({
                          ...shippingAddress,
                          phone: e.target.value,
                        })
                      }
                    />
                  </div>
                </div>
              </div>
            </div>

            {/* Billing Address */}
            <div className="card bg-base-100 shadow">
              <div className="card-body">
                <div className="flex items-center justify-between">
                  <h2 className="card-title">{t('billingAddress')}</h2>
                  <label className="label cursor-pointer">
                    <span className="label-text mr-2">
                      {t('sameAsShipping')}
                    </span>
                    <input
                      type="checkbox"
                      className="checkbox"
                      checked={useSameAddress}
                      onChange={(e) => setUseSameAddress(e.target.checked)}
                    />
                  </label>
                </div>

                {!useSameAddress && (
                  <div className="grid gap-4 md:grid-cols-2">
                    <div className="md:col-span-2">
                      <label className="label">
                        <span className="label-text">{t('street')}</span>
                      </label>
                      <input
                        type="text"
                        className="input input-bordered w-full"
                        value={billingAddress.street}
                        onChange={(e) =>
                          setBillingAddress({
                            ...billingAddress,
                            street: e.target.value,
                          })
                        }
                        required
                      />
                    </div>

                    <div>
                      <label className="label">
                        <span className="label-text">{t('city')}</span>
                      </label>
                      <input
                        type="text"
                        className="input input-bordered w-full"
                        value={billingAddress.city}
                        onChange={(e) =>
                          setBillingAddress({
                            ...billingAddress,
                            city: e.target.value,
                          })
                        }
                        required
                      />
                    </div>

                    <div>
                      <label className="label">
                        <span className="label-text">{t('postalCode')}</span>
                      </label>
                      <input
                        type="text"
                        className="input input-bordered w-full"
                        value={billingAddress.postal_code}
                        onChange={(e) =>
                          setBillingAddress({
                            ...billingAddress,
                            postal_code: e.target.value,
                          })
                        }
                        required
                      />
                    </div>
                  </div>
                )}
              </div>
            </div>

            {/* Payment Method */}
            <div className="card bg-base-100 shadow">
              <div className="card-body">
                <h2 className="card-title">{t('paymentMethod')}</h2>

                <div className="space-y-2">
                  <label className="label cursor-pointer">
                    <span className="label-text">AllSecure Payment</span>
                    <input
                      type="radio"
                      name="payment"
                      className="radio"
                      value="allsecure"
                      checked={paymentMethod === 'allsecure'}
                      onChange={(e) => setPaymentMethod(e.target.value)}
                    />
                  </label>
                </div>
              </div>
            </div>

            {/* Customer Notes */}
            <div className="card bg-base-100 shadow">
              <div className="card-body">
                <h2 className="card-title">{t('orderNotes')}</h2>
                <textarea
                  className="textarea textarea-bordered"
                  placeholder={t('orderNotesPlaceholder')}
                  value={customerNotes}
                  onChange={(e) => setCustomerNotes(e.target.value)}
                  rows={3}
                />
              </div>
            </div>

            {error && (
              <div className="alert alert-error">
                <span>{error}</span>
              </div>
            )}

            <button
              type="submit"
              className="btn btn-primary w-full"
              disabled={loading}
            >
              {loading && (
                <span className="loading loading-spinner loading-sm mr-2" />
              )}
              {t('placeOrder')}
            </button>
          </form>
        </div>

        {/* Order Summary */}
        <div>
          <div className="card bg-base-100 shadow">
            <div className="card-body">
              <h2 className="card-title">{t('orderSummary')}</h2>

              <div className="space-y-4">
                {cart.items.map((item: any) => (
                  <div key={item.id} className="flex items-center gap-4">
                    {item.product?.images?.[0] && (
                      <SafeImage
                        src={
                          item.product.images[0].thumbnail_url ||
                          item.product.images[0].image_url ||
                          ''
                        }
                        alt={item.product.name || 'Product'}
                        width={48}
                        height={48}
                        className="h-12 w-12 rounded object-cover"
                      />
                    )}

                    <div className="flex-1">
                      <h4 className="font-medium">{item.product?.name}</h4>
                      <p className="text-sm text-base-content/70">
                        {t('quantity')}: {item.quantity}
                      </p>
                    </div>

                    <div className="text-right">
                      <p className="font-medium">
                        {Number(item.total_price || 0).toFixed(2)}{' '}
                        {item.product?.currency}
                      </p>
                    </div>
                  </div>
                ))}
              </div>

              <div className="divider"></div>

              <div className="space-y-2">
                <div className="flex justify-between">
                  <span>{t('subtotal')}</span>
                  <span>
                    {subtotal.toFixed(2)} {cart.items[0]?.product?.currency}
                  </span>
                </div>

                <div className="flex justify-between">
                  <span>{t('shipping')}</span>
                  <span>
                    {shipping.toFixed(2)} {cart.items[0]?.product?.currency}
                  </span>
                </div>

                <div className="flex justify-between">
                  <span>{t('tax')}</span>
                  <span>
                    {tax.toFixed(2)} {cart.items[0]?.product?.currency}
                  </span>
                </div>

                <div className="divider"></div>

                <div className="flex justify-between text-lg font-bold">
                  <span>{t('total')}</span>
                  <span>
                    {total.toFixed(2)} {cart.items[0]?.product?.currency}
                  </span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
