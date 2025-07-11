'use client';

import React, { createContext, useContext, useReducer, ReactNode } from 'react';
import type { components } from '@/types/generated/api';

type CreateProductRequest =
  components['schemas']['backend_internal_domain_models.CreateProductRequest'];
type MarketplaceCategory =
  components['schemas']['backend_internal_domain_models.MarketplaceCategory'];

interface ProductState {
  // Шаги
  currentStep: number;
  completedSteps: Set<number>;

  // Данные продукта
  category?: MarketplaceCategory;
  productData: CreateProductRequest;
  attributes: Record<number, any>;
  images: File[];

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
      return {
        ...state,
        category: action.payload,
        productData: {
          ...state.productData,
          category_id: action.payload.id || 0,
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
    const next = Math.min(state.currentStep + 1, 4); // Максимум 5 шагов (0-4)
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
      case 2: // Атрибуты
        return true; // Атрибуты опциональны
      case 3: // Фотографии
        return state.images.length > 0;
      case 4: // Превью
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
