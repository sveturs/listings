'use client';

import React, {
  createContext,
  useContext,
  useState,
  useCallback,
  ReactNode,
} from 'react';
import { useTranslations } from 'next-intl';
import { toast } from '@/utils/toast';
import { storefrontApi } from '@/services/storefrontApi';
import type { StorefrontCreateDTO } from '@/types/storefront';

interface StorefrontFormData {
  // Basic Info
  name: string;
  slug: string;
  description: string;
  businessType: string;

  // Business Details
  registrationNumber?: string;
  taxNumber?: string;
  vatNumber?: string;
  website?: string;
  email?: string;
  phone?: string;

  // Location
  address: string;
  city: string;
  postalCode: string;
  country: string;
  latitude?: number;
  longitude?: number;

  // Business Hours
  businessHours: Array<{
    dayOfWeek: number;
    openTime: string;
    closeTime: string;
    isClosed: boolean;
  }>;

  // Payment & Delivery
  paymentMethods: string[];
  deliveryOptions: Array<{
    providerName: string;
    deliveryTimeMinutes: number;
    deliveryCostRSD: number;
    freeDeliveryThresholdRSD?: number;
    maxDistanceKm?: number;
  }>;

  // Staff
  staff: Array<{
    email: string;
    role: string;
    canManageProducts: boolean;
    canManageOrders: boolean;
    canManageSettings: boolean;
  }>;
}

interface CreateStorefrontContextType {
  formData: StorefrontFormData;
  updateFormData: (data: Partial<StorefrontFormData>) => void;
  resetFormData: () => void;
  isSubmitting: boolean;
  submitStorefront: () => Promise<{
    success: boolean;
    storefrontId?: number;
    error?: string;
  }>;
}

const initialFormData: StorefrontFormData = {
  name: '',
  slug: '',
  description: '',
  businessType: 'retail',
  address: '',
  city: '',
  postalCode: '',
  country: 'RS',
  businessHours: [
    { dayOfWeek: 1, openTime: '09:00', closeTime: '18:00', isClosed: false },
    { dayOfWeek: 2, openTime: '09:00', closeTime: '18:00', isClosed: false },
    { dayOfWeek: 3, openTime: '09:00', closeTime: '18:00', isClosed: false },
    { dayOfWeek: 4, openTime: '09:00', closeTime: '18:00', isClosed: false },
    { dayOfWeek: 5, openTime: '09:00', closeTime: '18:00', isClosed: false },
    { dayOfWeek: 6, openTime: '09:00', closeTime: '15:00', isClosed: false },
    { dayOfWeek: 0, openTime: '09:00', closeTime: '13:00', isClosed: true },
  ],
  paymentMethods: [],
  deliveryOptions: [],
  staff: [],
};

const CreateStorefrontContext =
  createContext<CreateStorefrontContextType | null>(null);

export const useCreateStorefrontContext = () => {
  const context = useContext(CreateStorefrontContext);
  if (!context) {
    throw new Error(
      'useCreateStorefrontContext must be used within CreateStorefrontProvider'
    );
  }
  return context;
};

interface CreateStorefrontProviderProps {
  children: ReactNode;
}

export const CreateStorefrontProvider: React.FC<
  CreateStorefrontProviderProps
> = ({ children }) => {
  const t = useTranslations();
  const [formData, setFormData] = useState<StorefrontFormData>(initialFormData);
  const [isSubmitting, setIsSubmitting] = useState(false);

  const updateFormData = useCallback((data: Partial<StorefrontFormData>) => {
    setFormData((prev) => ({ ...prev, ...data }));
  }, []);

  const resetFormData = useCallback(() => {
    setFormData(initialFormData);
  }, []);

  const submitStorefront = useCallback(async () => {
    setIsSubmitting(true);
    try {
      // Transform form data to API format
      const requestData: StorefrontCreateDTO = {
        name: formData.name,
        slug: formData.slug,
        description: formData.description,
        website: formData.website,
        email: formData.email,
        phone: formData.phone,
        location: {
          full_address: formData.address,
          city: formData.city,
          postal_code: formData.postalCode,
          country: formData.country,
          user_lat: formData.latitude,
          user_lng: formData.longitude,
        },
        settings: {
          businessType: formData.businessType,
          registrationNumber: formData.registrationNumber,
          taxNumber: formData.taxNumber,
          vatNumber: formData.vatNumber,
          businessHours: formData.businessHours,
          paymentMethods: formData.paymentMethods,
          deliveryOptions: formData.deliveryOptions,
          staff: formData.staff,
        },
      };

      const response = await storefrontApi.createStorefront(requestData);

      if (response?.id) {
        toast.success(t('create_storefront.success'));
        resetFormData();
        return { success: true, storefrontId: response.id };
      } else {
        throw new Error('Failed to create storefront');
      }
    } catch (error) {
      const errorMessage =
        error instanceof Error ? error.message : t('create_storefront.error');
      toast.error(errorMessage);
      return { success: false, error: errorMessage };
    } finally {
      setIsSubmitting(false);
    }
  }, [formData, t, resetFormData]);

  return (
    <CreateStorefrontContext.Provider
      value={{
        formData,
        updateFormData,
        resetFormData,
        isSubmitting,
        submitStorefront,
      }}
    >
      {children}
    </CreateStorefrontContext.Provider>
  );
};
