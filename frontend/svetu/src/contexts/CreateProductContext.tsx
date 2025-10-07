'use client';

import React, { createContext, useContext, useReducer, ReactNode } from 'react';
import type { components } from '@/types/generated/api';

type CreateProductRequest =
  components['schemas']['models.CreateProductRequest'];
type MarketplaceCategory =
  components['schemas']['models.MarketplaceCategory'];

interface ProductVariantCreate {
  sku?: string;
  barcode?: string;
  price?: number;
  compare_at_price?: number;
  cost_price?: number;
  stock_quantity: number;
  low_stock_threshold?: number;
  variant_attributes: Record<string, any>;
  weight?: number;
  dimensions?: Record<string, any>;
  is_default: boolean;
}

interface VariantSettings {
  track_inventory: boolean;
  continue_selling: boolean;
  require_shipping: boolean;
  taxable_product: boolean;
  weight_unit?: string;
  selected_attributes: string[];
}

interface ProductState {
  // Шаги
  currentStep: number;
  completedSteps: Set<number>;

  // Данные продукта
  category?: MarketplaceCategory;
  productData: CreateProductRequest;
  attributes: Record<number, any>;
  images: File[];

  // Данные о местоположении
  location?: {
    useStorefrontLocation: boolean;
    individualAddress?: string;
    latitude?: number;
    longitude?: number;
    city?: string;
    region?: string;
    country?: string;
    privacyLevel?: 'exact' | 'street' | 'district' | 'city';
    showOnMap?: boolean;
  };

  // Данные о вариантах
  hasVariants: boolean;
  variants: ProductVariantCreate[];
  variantSettings: VariantSettings;

  // Метаданные
  isDraft: boolean;
  isValid: boolean;
  errors: Record<string, string>;
  loading: boolean;
}

type ProductAction =
  | { type: 'SET_STEP'; payload: number }
  | { type: 'COMPLETE_STEP'; payload: number }
  | { type: 'SET_CATEGORY'; payload: MarketplaceCategory }
  | { type: 'SET_PRODUCT_DATA'; payload: Partial<CreateProductRequest> }
  | { type: 'SET_ATTRIBUTE'; payload: { id: number; value: any } }
  | { type: 'SET_IMAGES'; payload: File[] }
  | { type: 'SET_LOCATION'; payload: ProductState['location'] }
  | { type: 'SET_HAS_VARIANTS'; payload: boolean }
  | { type: 'SET_VARIANTS'; payload: ProductVariantCreate[] }
  | { type: 'ADD_VARIANT'; payload: ProductVariantCreate }
  | {
      type: 'UPDATE_VARIANT';
      payload: { index: number; variant: ProductVariantCreate };
    }
  | { type: 'REMOVE_VARIANT'; payload: number }
  | { type: 'SET_VARIANT_SETTINGS'; payload: Partial<VariantSettings> }
  | { type: 'SET_LOADING'; payload: boolean }
  | { type: 'SET_ERROR'; payload: { field: string; message: string } }
  | { type: 'CLEAR_ERROR'; payload: string }
  | { type: 'CLEAR_ALL_ERRORS' }
  | { type: 'RESET' }
  | { type: 'LOAD_DRAFT'; payload: ProductState };

const initialState: ProductState = {
  currentStep: 0,
  completedSteps: new Set(),
  productData: {
    name: '',
    description: '',
    price: 0,
    currency: 'RSD',
    category_id: 0,
    stock_quantity: 0,
    is_active: true,
    attributes: {},
  },
  attributes: {},
  images: [],
  location: {
    useStorefrontLocation: true,
    privacyLevel: 'exact',
    showOnMap: true,
  },
  hasVariants: false,
  variants: [],
  variantSettings: {
    track_inventory: true,
    continue_selling: false,
    require_shipping: true,
    taxable_product: true,
    selected_attributes: [],
  },
  isDraft: false,
  isValid: false,
  errors: {},
  loading: false,
};

function productReducer(
  state: ProductState,
  action: ProductAction
): ProductState {
  switch (action.type) {
    case 'SET_STEP':
      return { ...state, currentStep: action.payload };

    case 'COMPLETE_STEP':
      return {
        ...state,
        completedSteps: new Set([...state.completedSteps, action.payload]),
      };

    case 'SET_CATEGORY':
      // Clear completed steps after category (steps 2, 3, 4)
      const newCompletedSteps = new Set(state.completedSteps);
      newCompletedSteps.delete(2); // Attributes step
      newCompletedSteps.delete(3); // Photos step
      newCompletedSteps.delete(4); // Preview step

      return {
        ...state,
        category: action.payload,
        attributes: {}, // Clear attributes when category changes
        completedSteps: newCompletedSteps,
        productData: {
          ...state.productData,
          category_id: action.payload.id || 0,
          attributes: {}, // Clear attributes in productData too
        },
      };

    case 'SET_PRODUCT_DATA':
      return {
        ...state,
        productData: { ...state.productData, ...action.payload },
      };

    case 'SET_ATTRIBUTE':
      return {
        ...state,
        attributes: {
          ...state.attributes,
          [action.payload.id]: action.payload.value,
        },
        productData: {
          ...state.productData,
          attributes: {
            ...state.productData.attributes,
            [action.payload.id]: action.payload.value,
          },
        },
      };

    case 'SET_IMAGES':
      return { ...state, images: action.payload };

    case 'SET_LOCATION':
      return { ...state, location: action.payload };

    case 'SET_HAS_VARIANTS':
      return {
        ...state,
        hasVariants: action.payload,
        variants: action.payload ? state.variants : [],
      };

    case 'SET_VARIANTS':
      return { ...state, variants: action.payload };

    case 'ADD_VARIANT':
      return { ...state, variants: [...state.variants, action.payload] };

    case 'UPDATE_VARIANT':
      const updatedVariants = [...state.variants];
      updatedVariants[action.payload.index] = action.payload.variant;
      return { ...state, variants: updatedVariants };

    case 'REMOVE_VARIANT':
      return {
        ...state,
        variants: state.variants.filter((_, index) => index !== action.payload),
      };

    case 'SET_VARIANT_SETTINGS':
      return {
        ...state,
        variantSettings: { ...state.variantSettings, ...action.payload },
      };

    case 'SET_LOADING':
      return { ...state, loading: action.payload };

    case 'SET_ERROR':
      return {
        ...state,
        errors: {
          ...state.errors,
          [action.payload.field]: action.payload.message,
        },
      };

    case 'CLEAR_ERROR':
      const { [action.payload]: _, ...restErrors } = state.errors;
      return { ...state, errors: restErrors };

    case 'CLEAR_ALL_ERRORS':
      return { ...state, errors: {} };

    case 'RESET':
      return initialState;

    case 'LOAD_DRAFT':
      return { ...action.payload };

    default:
      return state;
  }
}

interface CreateProductContextType {
  state: ProductState;
  dispatch: React.Dispatch<ProductAction>;
  // Хелперы
  goToStep: (step: number) => void;
  nextStep: () => void;
  prevStep: () => void;
  completeStep: (step: number) => void;
  setCategory: (category: MarketplaceCategory) => void;
  setProductData: (data: Partial<CreateProductRequest>) => void;
  setAttribute: (id: number, value: any) => void;
  setImages: (images: File[]) => void;
  setLocation: (location: ProductState['location']) => void;
  setHasVariants: (hasVariants: boolean) => void;
  setVariants: (variants: ProductVariantCreate[]) => void;
  addVariant: (variant: ProductVariantCreate) => void;
  updateVariant: (index: number, variant: ProductVariantCreate) => void;
  removeVariant: (index: number) => void;
  setVariantSettings: (settings: Partial<VariantSettings>) => void;
  setError: (field: string, message: string) => void;
  clearError: (field: string) => void;
  clearAllErrors: () => void;
  validateStep: (step: number) => boolean;
  canGoToStep: (step: number) => boolean;
}

const CreateProductContext = createContext<
  CreateProductContextType | undefined
>(undefined);

export function CreateProductProvider({ children }: { children: ReactNode }) {
  const [state, dispatch] = useReducer(productReducer, initialState);

  const goToStep = (step: number) => {
    dispatch({ type: 'SET_STEP', payload: step });
  };

  const nextStep = () => {
    console.log('CreateProductContext nextStep called');
    console.log('Current step:', state.currentStep);
    const next = Math.min(state.currentStep + 1, 6); // Максимум 7 шагов (0-6)
    console.log('Next step:', next);
    dispatch({ type: 'SET_STEP', payload: next });
  };

  const prevStep = () => {
    const prev = Math.max(state.currentStep - 1, 0);
    dispatch({ type: 'SET_STEP', payload: prev });
  };

  const completeStep = (step: number) => {
    dispatch({ type: 'COMPLETE_STEP', payload: step });
  };

  const setCategory = (category: MarketplaceCategory) => {
    dispatch({ type: 'SET_CATEGORY', payload: category });
  };

  const setProductData = (data: Partial<CreateProductRequest>) => {
    dispatch({ type: 'SET_PRODUCT_DATA', payload: data });
  };

  const setAttribute = (id: number, value: any) => {
    dispatch({ type: 'SET_ATTRIBUTE', payload: { id, value } });
  };

  const setImages = (images: File[]) => {
    dispatch({ type: 'SET_IMAGES', payload: images });
  };

  const setLocation = (location: ProductState['location']) => {
    dispatch({ type: 'SET_LOCATION', payload: location });
  };

  const setHasVariants = (hasVariants: boolean) => {
    dispatch({ type: 'SET_HAS_VARIANTS', payload: hasVariants });
  };

  const setVariants = (variants: ProductVariantCreate[]) => {
    dispatch({ type: 'SET_VARIANTS', payload: variants });
  };

  const addVariant = (variant: ProductVariantCreate) => {
    dispatch({ type: 'ADD_VARIANT', payload: variant });
  };

  const updateVariant = (index: number, variant: ProductVariantCreate) => {
    dispatch({ type: 'UPDATE_VARIANT', payload: { index, variant } });
  };

  const removeVariant = (index: number) => {
    dispatch({ type: 'REMOVE_VARIANT', payload: index });
  };

  const setVariantSettings = (settings: Partial<VariantSettings>) => {
    dispatch({ type: 'SET_VARIANT_SETTINGS', payload: settings });
  };

  const setError = (field: string, message: string) => {
    dispatch({ type: 'SET_ERROR', payload: { field, message } });
  };

  const clearError = (field: string) => {
    dispatch({ type: 'CLEAR_ERROR', payload: field });
  };

  const clearAllErrors = () => {
    dispatch({ type: 'CLEAR_ALL_ERRORS' });
  };

  const validateStep = (step: number): boolean => {
    switch (step) {
      case 0: // Категория
        return !!state.category;
      case 1: // Основная информация
        return !!(
          state.productData.name &&
          state.productData.description &&
          state.productData.price > 0
        );
      case 2: // Местоположение
        if (state.location?.useStorefrontLocation) {
          return true; // Если используется адрес витрины, дополнительная валидация не нужна
        }
        // Если индивидуальный адрес, проверяем наличие координат
        return !!(
          state.location?.latitude &&
          state.location?.longitude &&
          state.location?.individualAddress
        );
      case 3: // Атрибуты
        return true; // Атрибуты опциональны
      case 4: // Варианты
        return true; // Варианты опциональны
      case 5: // Фотографии
        return state.images.length > 0;
      case 6: // Превью
        return true;
      default:
        return false;
    }
  };

  const canGoToStep = (step: number): boolean => {
    if (step <= state.currentStep) return true;

    // Проверяем все предыдущие шаги
    for (let i = 0; i < step; i++) {
      if (!validateStep(i)) return false;
    }
    return true;
  };

  const value: CreateProductContextType = {
    state,
    dispatch,
    goToStep,
    nextStep,
    prevStep,
    completeStep,
    setCategory,
    setProductData,
    setAttribute,
    setImages,
    setLocation,
    setHasVariants,
    setVariants,
    addVariant,
    updateVariant,
    removeVariant,
    setVariantSettings,
    setError,
    clearError,
    clearAllErrors,
    validateStep,
    canGoToStep,
  };

  return (
    <CreateProductContext.Provider value={value}>
      {children}
    </CreateProductContext.Provider>
  );
}

export function useCreateProduct(): CreateProductContextType {
  const context = useContext(CreateProductContext);
  if (context === undefined) {
    throw new Error(
      'useCreateProduct must be used within a CreateProductProvider'
    );
  }
  return context;
}
