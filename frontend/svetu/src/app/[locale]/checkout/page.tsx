'use client';

import React, { useState, useEffect } from 'react';
import { useSelector } from 'react-redux';
import { useRouter } from 'next/navigation';
import { useLocale, useTranslations } from 'next-intl';
import { useForm } from 'react-hook-form';
import { z } from 'zod';
import { zodResolver } from '@hookform/resolvers/zod';
import { motion, AnimatePresence } from 'framer-motion';
import Link from 'next/link';
import SafeImage from '@/components/SafeImage';
import {
  selectCartItems as selectLocalCartItems,
  selectCartTotal as selectLocalCartTotal,
  selectCartItemsCount as selectLocalCartItemsCount,
} from '@/store/slices/localCartSlice';
import {
  selectCart,
  selectCartTotal as selectApiCartTotal,
  selectAllCartsItemsCount as selectApiCartItemsCount,
} from '@/store/slices/cartSlice';
import { PageTransition } from '@/components/ui/PageTransition';
import { useAuthContext } from '@/contexts/AuthContext';
import type { components } from '@/types/generated/api';
import { useAppDispatch } from '@/store/hooks';
import { clearCart as clearLocalCart } from '@/store/slices/localCartSlice';
import { clearCart as clearApiCart } from '@/store/slices/cartSlice';
import { apiClient } from '@/services/api-client';

type CreateOrderRequest =
  components['schemas']['backend_internal_domain_models.CreateOrderRequest'];

// Validation schemas for each step
const customerInfoSchema = z.object({
  firstName: z.string().min(2, 'First name is required'),
  lastName: z.string().min(2, 'Last name is required'),
  email: z.string().email('Invalid email address'),
  phone: z.string().min(10, 'Phone number is required'),
});

const shippingAddressSchema = z.object({
  street: z.string().min(5, 'Street address is required'),
  city: z.string().min(2, 'City is required'),
  postalCode: z.string().min(5, 'Postal code is required'),
  country: z.string().min(2, 'Country is required'),
  shippingMethod: z.enum(['standard', 'express', 'pickup']),
});

const paymentSchema = z.object({
  paymentMethod: z.enum(['cod', 'bank_transfer', 'card']),
  acceptTerms: z.boolean().refine((val) => val === true, {
    message: 'You must accept the terms and conditions',
  }),
});

type CustomerInfo = z.infer<typeof customerInfoSchema>;
type ShippingAddress = z.infer<typeof shippingAddressSchema>;
type PaymentInfo = z.infer<typeof paymentSchema>;

const steps = ['customer', 'shipping', 'payment', 'review'] as const;
type Step = (typeof steps)[number];

export default function CheckoutPage() {
  const t = useTranslations('checkout');
  const locale = useLocale();
  const router = useRouter();
  const dispatch = useAppDispatch();
  const { user, isAuthenticated } = useAuthContext();

  // Выбираем данные из правильного slice в зависимости от авторизации
  const localItems = useSelector(selectLocalCartItems);
  const localTotal = useSelector(selectLocalCartTotal);
  const localItemsCount = useSelector(selectLocalCartItemsCount);

  const backendCart = useSelector(selectCart);
  const apiTotal = useSelector(selectApiCartTotal);
  const apiItemsCount = useSelector(selectApiCartItemsCount);

  // Используем данные в зависимости от авторизации
  const items = isAuthenticated
    ? (backendCart?.items || []).map((item) => ({
        productId: item.product_id || 0,
        variantId: item.variant_id,
        quantity: item.quantity || 0,
        name: item.product?.name || 'Product',
        variantName: item.variant?.name,
        price:
          typeof item.price_per_unit === 'string'
            ? parseFloat(item.price_per_unit)
            : item.price_per_unit || 0,
        currency: 'RSD',
        image: item.product?.images?.[0]?.image_url,
        storefrontId: backendCart?.storefront_id || 0,
        storefrontName: backendCart?.storefront?.name || 'Store',
        storefrontSlug:
          backendCart?.storefront?.slug ||
          String(backendCart?.storefront_id || 0),
      }))
    : localItems;

  const total = isAuthenticated ? apiTotal : localTotal;
  const itemsCount = isAuthenticated ? apiItemsCount : localItemsCount;

  const [currentStep, setCurrentStep] = useState<Step>('customer');
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [formData, setFormData] = useState<{
    customer?: CustomerInfo;
    shipping?: ShippingAddress;
    payment?: PaymentInfo;
  }>({});

  // Redirect if cart is empty
  useEffect(() => {
    // Небольшая задержка чтобы избежать мгновенного перенаправления при загрузке
    const timer = setTimeout(() => {
      if (itemsCount === 0 && !isSubmitting) {
        console.log(
          'Checkout: redirecting to cart, itemsCount:',
          itemsCount,
          'isAuthenticated:',
          isAuthenticated
        );
        router.push(`/${locale}/cart`);
      }
    }, 100);

    return () => clearTimeout(timer);
  }, [itemsCount, locale, router, isSubmitting, isAuthenticated]);

  // Pre-fill customer info if authenticated
  useEffect(() => {
    if (isAuthenticated && user) {
      const [firstName, ...lastNameParts] = (user.name || '').split(' ');
      setFormData((prev) => ({
        ...prev,
        customer: {
          firstName: firstName || '',
          lastName: lastNameParts.join(' ') || '',
          email: user.email || '',
          phone: user.phone || '',
        },
      }));
    }
  }, [isAuthenticated, user]);

  const customerForm = useForm<CustomerInfo>({
    resolver: zodResolver(customerInfoSchema),
    defaultValues: formData.customer || {},
  });

  const shippingForm = useForm<ShippingAddress>({
    resolver: zodResolver(shippingAddressSchema),
    defaultValues: formData.shipping || {
      shippingMethod: 'standard',
      country: 'Serbia',
    },
  });

  const paymentForm = useForm<PaymentInfo>({
    resolver: zodResolver(paymentSchema),
    defaultValues: formData.payment || {
      paymentMethod: 'cod',
    },
  });

  const getStepIndex = (step: Step) => steps.indexOf(step);
  const currentStepIndex = getStepIndex(currentStep);

  const handleNextStep = async () => {
    let isValid = false;
    let data: any = {};

    switch (currentStep) {
      case 'customer':
        isValid = await customerForm.trigger();
        if (isValid) {
          data = customerForm.getValues();
          setFormData((prev) => ({ ...prev, customer: data }));
          setCurrentStep('shipping');
        }
        break;
      case 'shipping':
        isValid = await shippingForm.trigger();
        if (isValid) {
          data = shippingForm.getValues();
          setFormData((prev) => ({ ...prev, shipping: data }));
          setCurrentStep('payment');
        }
        break;
      case 'payment':
        isValid = await paymentForm.trigger();
        if (isValid) {
          data = paymentForm.getValues();
          setFormData((prev) => ({ ...prev, payment: data }));
          setCurrentStep('review');
        }
        break;
    }
  };

  const handlePreviousStep = () => {
    const currentIndex = getStepIndex(currentStep);
    if (currentIndex > 0) {
      setCurrentStep(steps[currentIndex - 1]);
    }
  };

  const calculateShipping = () => {
    if (formData.shipping?.shippingMethod === 'pickup') return 0;
    if (formData.shipping?.shippingMethod === 'express') return 500;
    return total >= 5000 ? 0 : 300;
  };

  const shipping = calculateShipping();
  const finalTotal = total + shipping;

  const handleSubmitOrder = async () => {
    setIsSubmitting(true);

    try {
      // Group items by storefront for API
      const itemsByStorefront = items.reduce(
        (acc, item) => {
          if (!acc[item.storefrontId]) {
            acc[item.storefrontId] = [];
          }
          acc[item.storefrontId].push({
            product_id: item.productId,
            variant_id: item.variantId,
            quantity: item.quantity,
            price: item.price,
          });
          return acc;
        },
        {} as Record<number, any[]>
      );

      // Create orders for each storefront
      for (const [storefrontId, storefrontItems] of Object.entries(
        itemsByStorefront
      )) {
        const orderData: CreateOrderRequest = {
          storefront_id: parseInt(storefrontId),
          items: storefrontItems,
          shipping_address: {
            street: formData.shipping!.street,
            city: formData.shipping!.city,
            postal_code: formData.shipping!.postalCode,
            country: formData.shipping!.country,
          },
          billing_address: {
            street: formData.shipping!.street,
            city: formData.shipping!.city,
            postal_code: formData.shipping!.postalCode,
            country: formData.shipping!.country,
          },
          payment_method: formData.payment!.paymentMethod,
          shipping_method: formData.shipping!.shippingMethod,
          customer_notes: '',
        };

        // Создаем заказ через API
        const response = await apiClient.post('/api/v1/orders', orderData);

        if (response.error) {
          throw new Error(response.error.message || 'Failed to create order');
        }

        console.log('Order created successfully:', response.data);

        // Очищаем корзину в зависимости от авторизации
        if (isAuthenticated) {
          // Для авторизованных пользователей очищаем API корзину
          dispatch(clearApiCart(parseInt(storefrontId)));
        } else {
          // Для неавторизованных пользователей очищаем локальную корзину
          dispatch(clearLocalCart());
        }

        // Перенаправляем на страницу успешного заказа
        if (response.data && response.data.data && response.data.data.id) {
          router.push(
            `/${locale}/checkout/success?orderId=${response.data.data.id}`
          );
        } else {
          router.push(`/${locale}/checkout/success`);
        }
      }
    } catch (error) {
      console.error('Failed to create order:', error);
      // Show error toast instead of setting state
      const errorMessage =
        error instanceof Error ? error.message : 'Failed to create order';
      alert(errorMessage); // Replace with toast notification in production
      setIsSubmitting(false);
    }
  };

  // Group items by storefront for display
  const itemsByStorefront = items.reduce(
    (acc, item) => {
      const storefrontId = item.storefrontId;
      if (!acc[storefrontId]) {
        acc[storefrontId] = {
          items: [],
          name: item.storefrontName || `Store ${storefrontId}`,
          subtotal: 0,
        };
      }
      acc[storefrontId].items.push(item);
      acc[storefrontId].subtotal += item.price * item.quantity;
      return acc;
    },
    {} as Record<
      number,
      { items: (typeof items)[0][]; name: string; subtotal: number }
    >
  );

  return (
    <PageTransition>
      <div className="min-h-screen bg-base-100 pt-24">
        <div className="container mx-auto px-4 py-8">
          {/* Progress Steps */}
          <div className="mb-8">
            <ul className="steps steps-horizontal w-full">
              {steps.map((step, index) => (
                <li
                  key={step}
                  className={`step ${index <= currentStepIndex ? 'step-primary' : ''}`}
                  onClick={() => {
                    if (index < currentStepIndex) {
                      setCurrentStep(step);
                    }
                  }}
                >
                  <span className="hidden sm:inline">{t(`steps.${step}`)}</span>
                  <span className="sm:hidden">{index + 1}</span>
                </li>
              ))}
            </ul>
          </div>

          <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
            {/* Main Content */}
            <div className="lg:col-span-2">
              <AnimatePresence mode="wait">
                <motion.div
                  key={currentStep}
                  initial={{ opacity: 0, x: 20 }}
                  animate={{ opacity: 1, x: 0 }}
                  exit={{ opacity: 0, x: -20 }}
                  transition={{ duration: 0.3 }}
                >
                  {/* Customer Information */}
                  {currentStep === 'customer' && (
                    <div className="card bg-base-200">
                      <div className="card-body">
                        <h2 className="card-title">{t('customer.title')}</h2>
                        <form className="space-y-4">
                          <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                            <div className="form-control">
                              <label className="label">
                                <span className="label-text">
                                  {t('customer.firstName')}
                                </span>
                              </label>
                              <input
                                type="text"
                                className="input input-bordered"
                                {...customerForm.register('firstName')}
                              />
                              {customerForm.formState.errors.firstName && (
                                <label className="label">
                                  <span className="label-text-alt text-error">
                                    {
                                      customerForm.formState.errors.firstName
                                        .message
                                    }
                                  </span>
                                </label>
                              )}
                            </div>
                            <div className="form-control">
                              <label className="label">
                                <span className="label-text">
                                  {t('customer.lastName')}
                                </span>
                              </label>
                              <input
                                type="text"
                                className="input input-bordered"
                                {...customerForm.register('lastName')}
                              />
                              {customerForm.formState.errors.lastName && (
                                <label className="label">
                                  <span className="label-text-alt text-error">
                                    {
                                      customerForm.formState.errors.lastName
                                        .message
                                    }
                                  </span>
                                </label>
                              )}
                            </div>
                          </div>
                          <div className="form-control">
                            <label className="label">
                              <span className="label-text">
                                {t('customer.email')}
                              </span>
                            </label>
                            <input
                              type="email"
                              className="input input-bordered"
                              {...customerForm.register('email')}
                            />
                            {customerForm.formState.errors.email && (
                              <label className="label">
                                <span className="label-text-alt text-error">
                                  {customerForm.formState.errors.email.message}
                                </span>
                              </label>
                            )}
                          </div>
                          <div className="form-control">
                            <label className="label">
                              <span className="label-text">
                                {t('customer.phone')}
                              </span>
                            </label>
                            <input
                              type="tel"
                              className="input input-bordered"
                              {...customerForm.register('phone')}
                            />
                            {customerForm.formState.errors.phone && (
                              <label className="label">
                                <span className="label-text-alt text-error">
                                  {customerForm.formState.errors.phone.message}
                                </span>
                              </label>
                            )}
                          </div>
                        </form>
                      </div>
                    </div>
                  )}

                  {/* Shipping Information */}
                  {currentStep === 'shipping' && (
                    <div className="card bg-base-200">
                      <div className="card-body">
                        <h2 className="card-title">{t('shipping.title')}</h2>
                        <form className="space-y-4">
                          <div className="form-control">
                            <label className="label">
                              <span className="label-text">
                                {t('shipping.street')}
                              </span>
                            </label>
                            <input
                              type="text"
                              className="input input-bordered"
                              {...shippingForm.register('street')}
                            />
                            {shippingForm.formState.errors.street && (
                              <label className="label">
                                <span className="label-text-alt text-error">
                                  {shippingForm.formState.errors.street.message}
                                </span>
                              </label>
                            )}
                          </div>
                          <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                            <div className="form-control">
                              <label className="label">
                                <span className="label-text">
                                  {t('shipping.city')}
                                </span>
                              </label>
                              <input
                                type="text"
                                className="input input-bordered"
                                {...shippingForm.register('city')}
                              />
                              {shippingForm.formState.errors.city && (
                                <label className="label">
                                  <span className="label-text-alt text-error">
                                    {shippingForm.formState.errors.city.message}
                                  </span>
                                </label>
                              )}
                            </div>
                            <div className="form-control">
                              <label className="label">
                                <span className="label-text">
                                  {t('shipping.postalCode')}
                                </span>
                              </label>
                              <input
                                type="text"
                                className="input input-bordered"
                                {...shippingForm.register('postalCode')}
                              />
                              {shippingForm.formState.errors.postalCode && (
                                <label className="label">
                                  <span className="label-text-alt text-error">
                                    {
                                      shippingForm.formState.errors.postalCode
                                        .message
                                    }
                                  </span>
                                </label>
                              )}
                            </div>
                          </div>
                          <div className="form-control">
                            <label className="label">
                              <span className="label-text">
                                {t('shipping.country')}
                              </span>
                            </label>
                            <input
                              type="text"
                              className="input input-bordered"
                              {...shippingForm.register('country')}
                            />
                            {shippingForm.formState.errors.country && (
                              <label className="label">
                                <span className="label-text-alt text-error">
                                  {
                                    shippingForm.formState.errors.country
                                      .message
                                  }
                                </span>
                              </label>
                            )}
                          </div>

                          {/* Shipping Method */}
                          <div className="form-control">
                            <label className="label">
                              <span className="label-text">
                                {t('shipping.method')}
                              </span>
                            </label>
                            <div className="space-y-2">
                              <label className="label cursor-pointer justify-start gap-4">
                                <input
                                  type="radio"
                                  className="radio radio-primary"
                                  value="standard"
                                  {...shippingForm.register('shippingMethod')}
                                />
                                <div className="flex-1">
                                  <span className="label-text">
                                    {t('shipping.standard')}
                                  </span>
                                  <p className="text-sm text-base-content/60">
                                    {t('shipping.standardDesc')}
                                  </p>
                                </div>
                                <span className="font-semibold">
                                  {total >= 5000 ? t('free') : '300 RSD'}
                                </span>
                              </label>
                              <label className="label cursor-pointer justify-start gap-4">
                                <input
                                  type="radio"
                                  className="radio radio-primary"
                                  value="express"
                                  {...shippingForm.register('shippingMethod')}
                                />
                                <div className="flex-1">
                                  <span className="label-text">
                                    {t('shipping.express')}
                                  </span>
                                  <p className="text-sm text-base-content/60">
                                    {t('shipping.expressDesc')}
                                  </p>
                                </div>
                                <span className="font-semibold">500 RSD</span>
                              </label>
                              <label className="label cursor-pointer justify-start gap-4">
                                <input
                                  type="radio"
                                  className="radio radio-primary"
                                  value="pickup"
                                  {...shippingForm.register('shippingMethod')}
                                />
                                <div className="flex-1">
                                  <span className="label-text">
                                    {t('shipping.pickup')}
                                  </span>
                                  <p className="text-sm text-base-content/60">
                                    {t('shipping.pickupDesc')}
                                  </p>
                                </div>
                                <span className="font-semibold">
                                  {t('free')}
                                </span>
                              </label>
                            </div>
                          </div>
                        </form>
                      </div>
                    </div>
                  )}

                  {/* Payment Information */}
                  {currentStep === 'payment' && (
                    <div className="card bg-base-200">
                      <div className="card-body">
                        <h2 className="card-title">{t('payment.title')}</h2>
                        <form className="space-y-4">
                          <div className="form-control">
                            <label className="label">
                              <span className="label-text">
                                {t('payment.method')}
                              </span>
                            </label>
                            <div className="space-y-2">
                              <label className="label cursor-pointer justify-start gap-4">
                                <input
                                  type="radio"
                                  className="radio radio-primary"
                                  value="cod"
                                  {...paymentForm.register('paymentMethod')}
                                />
                                <div className="flex-1">
                                  <span className="label-text">
                                    {t('payment.cod')}
                                  </span>
                                  <p className="text-sm text-base-content/60">
                                    {t('payment.codDesc')}
                                  </p>
                                </div>
                              </label>
                              <label className="label cursor-pointer justify-start gap-4">
                                <input
                                  type="radio"
                                  className="radio radio-primary"
                                  value="bank_transfer"
                                  {...paymentForm.register('paymentMethod')}
                                />
                                <div className="flex-1">
                                  <span className="label-text">
                                    {t('payment.bankTransfer')}
                                  </span>
                                  <p className="text-sm text-base-content/60">
                                    {t('payment.bankTransferDesc')}
                                  </p>
                                </div>
                              </label>
                              <label className="label cursor-pointer justify-start gap-4 opacity-50">
                                <input
                                  type="radio"
                                  className="radio radio-primary"
                                  value="card"
                                  disabled
                                  {...paymentForm.register('paymentMethod')}
                                />
                                <div className="flex-1">
                                  <span className="label-text">
                                    {t('payment.card')}
                                  </span>
                                  <p className="text-sm text-base-content/60">
                                    {t('payment.cardDisabled')}
                                  </p>
                                </div>
                              </label>
                            </div>
                          </div>

                          <div className="form-control">
                            <label className="label cursor-pointer justify-start gap-4">
                              <input
                                type="checkbox"
                                className="checkbox checkbox-primary"
                                {...paymentForm.register('acceptTerms')}
                              />
                              <span className="label-text">
                                {t('payment.acceptTerms')}{' '}
                                <Link
                                  href="/terms"
                                  className="link link-primary"
                                >
                                  {t('payment.termsLink')}
                                </Link>
                              </span>
                            </label>
                            {paymentForm.formState.errors.acceptTerms && (
                              <label className="label">
                                <span className="label-text-alt text-error">
                                  {
                                    paymentForm.formState.errors.acceptTerms
                                      .message
                                  }
                                </span>
                              </label>
                            )}
                          </div>
                        </form>
                      </div>
                    </div>
                  )}

                  {/* Order Review */}
                  {currentStep === 'review' && (
                    <div className="space-y-6">
                      {/* Customer Info Review */}
                      <div className="card bg-base-200">
                        <div className="card-body">
                          <div className="flex justify-between items-center">
                            <h3 className="font-semibold">
                              {t('review.customer')}
                            </h3>
                            <button
                              onClick={() => setCurrentStep('customer')}
                              className="btn btn-ghost btn-sm"
                            >
                              {t('edit')}
                            </button>
                          </div>
                          <div className="text-sm space-y-1">
                            <p>
                              {formData.customer?.firstName}{' '}
                              {formData.customer?.lastName}
                            </p>
                            <p>{formData.customer?.email}</p>
                            <p>{formData.customer?.phone}</p>
                          </div>
                        </div>
                      </div>

                      {/* Shipping Info Review */}
                      <div className="card bg-base-200">
                        <div className="card-body">
                          <div className="flex justify-between items-center">
                            <h3 className="font-semibold">
                              {t('review.shipping')}
                            </h3>
                            <button
                              onClick={() => setCurrentStep('shipping')}
                              className="btn btn-ghost btn-sm"
                            >
                              {t('edit')}
                            </button>
                          </div>
                          <div className="text-sm space-y-1">
                            <p>{formData.shipping?.street}</p>
                            <p>
                              {formData.shipping?.city},{' '}
                              {formData.shipping?.postalCode}
                            </p>
                            <p>{formData.shipping?.country}</p>
                            <p className="font-medium mt-2">
                              {t(
                                `shipping.${formData.shipping?.shippingMethod}`
                              )}
                            </p>
                          </div>
                        </div>
                      </div>

                      {/* Payment Info Review */}
                      <div className="card bg-base-200">
                        <div className="card-body">
                          <div className="flex justify-between items-center">
                            <h3 className="font-semibold">
                              {t('review.payment')}
                            </h3>
                            <button
                              onClick={() => setCurrentStep('payment')}
                              className="btn btn-ghost btn-sm"
                            >
                              {t('edit')}
                            </button>
                          </div>
                          <div className="text-sm">
                            <p className="font-medium">
                              {t(`payment.${formData.payment?.paymentMethod}`)}
                            </p>
                          </div>
                        </div>
                      </div>

                      {/* Order Items Review */}
                      {Object.entries(itemsByStorefront).map(
                        ([storefrontId, group]) => (
                          <div key={storefrontId} className="card bg-base-200">
                            <div className="card-body">
                              <h3 className="font-semibold mb-4">
                                {group.name}
                              </h3>
                              <div className="space-y-3">
                                {group.items.map((item) => (
                                  <div
                                    key={`${item.productId}-${item.variantId || ''}`}
                                    className="flex gap-3"
                                  >
                                    <div className="relative w-16 h-16 rounded-lg overflow-hidden bg-base-100 flex-shrink-0">
                                      <SafeImage
                                        src={item.image}
                                        alt={item.name}
                                        fill
                                        className="object-cover"
                                        sizes="64px"
                                      />
                                    </div>
                                    <div className="flex-1">
                                      <p className="font-medium">{item.name}</p>
                                      {item.variantName && (
                                        <p className="text-sm text-base-content/60">
                                          {item.variantName}
                                        </p>
                                      )}
                                      <p className="text-sm">
                                        {item.quantity} ×{' '}
                                        {item.price.toFixed(2)} {item.currency}
                                      </p>
                                    </div>
                                    <div className="text-right">
                                      <p className="font-semibold">
                                        {(item.price * item.quantity).toFixed(
                                          2
                                        )}{' '}
                                        {item.currency}
                                      </p>
                                    </div>
                                  </div>
                                ))}
                              </div>
                            </div>
                          </div>
                        )
                      )}
                    </div>
                  )}

                  {/* Navigation Buttons */}
                  <div className="flex justify-between mt-6">
                    {currentStep !== 'customer' && (
                      <button
                        onClick={handlePreviousStep}
                        className="btn btn-outline"
                      >
                        {t('previous')}
                      </button>
                    )}
                    {currentStep !== 'review' ? (
                      <button
                        onClick={handleNextStep}
                        className="btn btn-primary ml-auto"
                      >
                        {t('next')}
                      </button>
                    ) : (
                      <button
                        onClick={handleSubmitOrder}
                        className="btn btn-primary ml-auto"
                        disabled={isSubmitting}
                      >
                        {isSubmitting && (
                          <span className="loading loading-spinner" />
                        )}
                        {t('placeOrder')}
                      </button>
                    )}
                  </div>
                </motion.div>
              </AnimatePresence>
            </div>

            {/* Order Summary Sidebar */}
            <div className="lg:col-span-1">
              <div className="card bg-base-200 sticky top-24">
                <div className="card-body">
                  <h2 className="card-title mb-4">{t('summary.title')}</h2>

                  <div className="space-y-2">
                    <div className="flex justify-between">
                      <span>{t('summary.items', { count: itemsCount })}</span>
                      <span>{total.toFixed(2)} RSD</span>
                    </div>
                    <div className="flex justify-between">
                      <span>{t('summary.shipping')}</span>
                      <span>
                        {shipping === 0 ? (
                          <span className="text-success">{t('free')}</span>
                        ) : (
                          `${shipping} RSD`
                        )}
                      </span>
                    </div>
                    <div className="divider"></div>
                    <div className="flex justify-between text-lg font-bold">
                      <span>{t('summary.total')}</span>
                      <span className="text-primary">
                        {finalTotal.toFixed(2)} RSD
                      </span>
                    </div>
                  </div>

                  {/* Security badges */}
                  <div className="mt-6 pt-6 border-t border-base-300">
                    <div className="space-y-2 text-sm text-base-content/60">
                      <div className="flex items-center gap-2">
                        <svg
                          className="w-4 h-4 text-success"
                          fill="currentColor"
                          viewBox="0 0 20 20"
                        >
                          <path
                            fillRule="evenodd"
                            d="M2.166 4.999A11.954 11.954 0 0010 1.944 11.954 11.954 0 0017.834 5c.11.65.166 1.32.166 2.001 0 5.225-3.34 9.67-8 11.317C5.34 16.67 2 12.225 2 7c0-.682.057-1.35.166-2.001zm11.541 3.708a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z"
                            clipRule="evenodd"
                          />
                        </svg>
                        <span>{t('summary.secureCheckout')}</span>
                      </div>
                      <div className="flex items-center gap-2">
                        <svg
                          className="w-4 h-4 text-success"
                          fill="currentColor"
                          viewBox="0 0 20 20"
                        >
                          <path d="M8 16.5a1.5 1.5 0 11-3 0 1.5 1.5 0 013 0zM15 16.5a1.5 1.5 0 11-3 0 1.5 1.5 0 013 0z" />
                          <path d="M3 4a1 1 0 00-1 1v10a1 1 0 001 1h1.05a2.5 2.5 0 014.9 0H10a1 1 0 001-1V5a1 1 0 00-1-1H3zM14 7a1 1 0 00-1 1v6.05A2.5 2.5 0 0115.95 16H17a1 1 0 001-1v-5a1 1 0 00-.293-.707l-2-2A1 1 0 0015 7h-1z" />
                        </svg>
                        <span>{t('summary.freeReturns')}</span>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </PageTransition>
  );
}
