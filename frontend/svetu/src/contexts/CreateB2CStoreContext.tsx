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
import { storefrontApi } from '@/services/b2cStoreApi';
import type { B2CStoreCreateDTO } from '@/types/b2c';
import dynamic from 'next/dynamic';

const UpgradePrompt = dynamic(() => import('@/components/b2c/UpgradePrompt'), {
  ssr: false,
});

interface B2CStoreFormData {
  // Basic Info
  name: string;
  slug: string;
  description: string;
  businessType: string;

  // Images
  logoFile?: File;
  bannerFile?: File;
  logoUrl?: string;
  bannerUrl?: string;

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
  formData: B2CStoreFormData;
  updateFormData: (data: Partial<B2CStoreFormData>) => void;
  resetFormData: () => void;
  isSubmitting: boolean;
  submitStorefront: () => Promise<{
    success: boolean;
    storefrontId?: number;
    error?: string;
  }>;
}

const initialFormData: B2CStoreFormData = {
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
  const t = useTranslations('create_storefront');
  const tMisc = useTranslations('misc');
  const [formData, setFormData] = useState<B2CStoreFormData>(initialFormData);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [showUpgradePrompt, setShowUpgradePrompt] = useState(false);

  const updateFormData = useCallback((data: Partial<B2CStoreFormData>) => {
    setFormData((prev) => ({ ...prev, ...data }));
  }, []);

  const resetFormData = useCallback(() => {
    setFormData(initialFormData);
  }, []);

  const submitStorefront = useCallback(async () => {
    setIsSubmitting(true);
    try {
      // Transform form data to API format
      const requestData: B2CStoreCreateDTO = {
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
        // Upload images if provided
        try {
          const uploadPromises = [];

          if (formData.logoFile) {
            uploadPromises.push(
              storefrontApi
                .uploadLogo(response.id, formData.logoFile)
                .catch((err) => console.error('Failed to upload logo:', err))
            );
          }

          if (formData.bannerFile) {
            uploadPromises.push(
              storefrontApi
                .uploadBanner(response.id, formData.bannerFile)
                .catch((err) => console.error('Failed to upload banner:', err))
            );
          }

          if (uploadPromises.length > 0) {
            await Promise.all(uploadPromises);
          }
        } catch (error) {
          console.error('Error uploading images:', error);
          // Don't fail the whole process if image upload fails
        }

        toast.success(t('success'));
        resetFormData();
        return { success: true, storefrontId: response.id };
      } else {
        throw new Error('Failed to create storefront');
      }
    } catch (error) {
      let errorMessage = error instanceof Error ? error.message : t('error');

      // Проверяем, является ли это ошибкой лимита витрин
      const isLimitError = errorMessage === 'storefronts.error.limit_reached';

      // Попытаемся перевести ключ ошибки, если он похож на ключ перевода
      if (errorMessage && errorMessage.includes('.')) {
        try {
          // Проверяем, является ли это ключом перевода из misc
          // Backend отправляет "storefronts.error.limit_reached"
          // но в misc.json это находится по пути "storefronts.error.limit_reached"
          const translatedError = tMisc(errorMessage as any);
          if (translatedError && translatedError !== errorMessage) {
            errorMessage = translatedError;
          }
        } catch {
          // Если не удалось перевести, используем оригинальное сообщение
        }
      }

      // Если это ошибка лимита, показываем upgrade prompt вместо toast
      if (isLimitError) {
        setShowUpgradePrompt(true);
      } else {
        toast.error(errorMessage);
      }

      return { success: false, error: errorMessage };
    } finally {
      setIsSubmitting(false);
    }
  }, [formData, t, tMisc, resetFormData]);

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
      {showUpgradePrompt && (
        <UpgradePrompt
          currentPlan="starter"
          onClose={() => setShowUpgradePrompt(false)}
        />
      )}
    </CreateStorefrontContext.Provider>
  );
};
