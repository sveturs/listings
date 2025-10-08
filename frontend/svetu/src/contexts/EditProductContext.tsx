'use client';

import React, {
  createContext,
  useContext,
  useReducer,
  ReactNode,
  useCallback,
} from 'react';
import type { components } from '@/types/generated/api';

type StorefrontProduct = components['schemas']['models.StorefrontProduct'];
type UpdateProductRequest =
  components['schemas']['models.UpdateProductRequest'];
type MarketplaceCategory = components['schemas']['models.MarketplaceCategory'];

interface ProductEditState {
  // Шаги
  currentStep: number;
  completedSteps: Set<number>;

  // Исходные данные продукта
  originalProduct?: StorefrontProduct;

  // Данные для редактирования
  category?: MarketplaceCategory;
  productData: UpdateProductRequest;
  attributes: Record<number, any>;

  // Изображения
  existingImages: Array<{
    id: number;
    url: string;
    is_main: boolean;
  }>;
  newImages: File[];
  imagesToDelete: number[];

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

  // Метаданные
  isLoading: boolean;
  isSaving: boolean;
  isValid: boolean;
  errors: Record<string, string>;
  hasUnsavedChanges: boolean;
}

type ProductEditAction =
  | { type: 'SET_STEP'; payload: number }
  | { type: 'COMPLETE_STEP'; payload: number }
  | { type: 'LOAD_PRODUCT'; payload: StorefrontProduct }
  | { type: 'SET_CATEGORY'; payload: MarketplaceCategory }
  | { type: 'SET_PRODUCT_DATA'; payload: Partial<UpdateProductRequest> }
  | { type: 'SET_ATTRIBUTE'; payload: { id: number; value: any } }
  | { type: 'SET_EXISTING_IMAGES'; payload: ProductEditState['existingImages'] }
  | { type: 'SET_NEW_IMAGES'; payload: File[] }
  | { type: 'ADD_NEW_IMAGE'; payload: File }
  | { type: 'REMOVE_NEW_IMAGE'; payload: number }
  | { type: 'MARK_IMAGE_FOR_DELETION'; payload: number }
  | { type: 'RESTORE_IMAGE'; payload: number }
  | {
      type: 'SET_MAIN_IMAGE';
      payload: { type: 'existing' | 'new'; index: number };
    }
  | { type: 'SET_LOCATION'; payload: ProductEditState['location'] }
  | { type: 'SET_LOADING'; payload: boolean }
  | { type: 'SET_SAVING'; payload: boolean }
  | { type: 'SET_ERROR'; payload: { field: string; message: string } }
  | { type: 'CLEAR_ERROR'; payload: string }
  | { type: 'CLEAR_ALL_ERRORS' }
  | { type: 'SET_UNSAVED_CHANGES'; payload: boolean }
  | { type: 'RESET' };

const initialState: ProductEditState = {
  currentStep: 0,
  completedSteps: new Set(),
  productData: {
    name: '',
    description: '',
    price: 0,
    category_id: 0,
    stock_quantity: 0,
    is_active: true,
    attributes: {},
    has_individual_location: false,
    location_privacy: 'exact',
    show_on_map: true,
  },
  attributes: {},
  existingImages: [],
  newImages: [],
  imagesToDelete: [],
  location: {
    useStorefrontLocation: true,
    privacyLevel: 'exact',
    showOnMap: true,
  },
  isLoading: false,
  isSaving: false,
  isValid: false,
  errors: {},
  hasUnsavedChanges: false,
};

function productEditReducer(
  state: ProductEditState,
  action: ProductEditAction
): ProductEditState {
  switch (action.type) {
    case 'SET_STEP':
      return { ...state, currentStep: action.payload };

    case 'COMPLETE_STEP':
      return {
        ...state,
        completedSteps: new Set([...state.completedSteps, action.payload]),
      };

    case 'LOAD_PRODUCT':
      const product = action.payload;
      return {
        ...state,
        originalProduct: product,
        category: product.category,
        productData: {
          name: product.name || '',
          description: product.description || '',
          price: product.price || 0,
          category_id: product.category_id || 0,
          stock_quantity: product.stock_quantity || 0,
          is_active: product.is_active ?? true,
          attributes: product.attributes || {},
          has_individual_location: product.has_individual_location || false,
          individual_address: product.individual_address || '',
          individual_latitude: product.individual_latitude || undefined,
          individual_longitude: product.individual_longitude || undefined,
          location_privacy: (product.location_privacy || 'exact') as
            | 'exact'
            | 'street'
            | 'district'
            | 'city',
          show_on_map: product.show_on_map ?? true,
          sku: product.sku || '',
          barcode: product.barcode || '',
        },
        attributes: product.attributes || {},
        existingImages:
          product.images?.map((img) => ({
            id: img.id || 0,
            url: img.public_url || img.image_url || '',
            is_main: img.is_default || false,
          })) || [],
        location: {
          useStorefrontLocation: !product.has_individual_location,
          individualAddress: product.individual_address || '',
          latitude: product.individual_latitude || undefined,
          longitude: product.individual_longitude || undefined,
          privacyLevel: (product.location_privacy as any) || 'exact',
          showOnMap: product.show_on_map ?? true,
        },
        isLoading: false,
        hasUnsavedChanges: false,
      };

    case 'SET_CATEGORY':
      const newCompletedSteps = new Set(state.completedSteps);
      newCompletedSteps.delete(3); // Attributes step

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
        hasUnsavedChanges: true,
      };

    case 'SET_PRODUCT_DATA':
      return {
        ...state,
        productData: { ...state.productData, ...action.payload },
        hasUnsavedChanges: true,
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
        hasUnsavedChanges: true,
      };

    case 'SET_EXISTING_IMAGES':
      return {
        ...state,
        existingImages: action.payload,
        hasUnsavedChanges: true,
      };

    case 'SET_NEW_IMAGES':
      return {
        ...state,
        newImages: action.payload,
        hasUnsavedChanges: true,
      };

    case 'ADD_NEW_IMAGE':
      return {
        ...state,
        newImages: [...state.newImages, action.payload],
        hasUnsavedChanges: true,
      };

    case 'REMOVE_NEW_IMAGE':
      return {
        ...state,
        newImages: state.newImages.filter(
          (_, index) => index !== action.payload
        ),
        hasUnsavedChanges: true,
      };

    case 'MARK_IMAGE_FOR_DELETION':
      return {
        ...state,
        imagesToDelete: [...state.imagesToDelete, action.payload],
        hasUnsavedChanges: true,
      };

    case 'RESTORE_IMAGE':
      return {
        ...state,
        imagesToDelete: state.imagesToDelete.filter(
          (id) => id !== action.payload
        ),
        hasUnsavedChanges: true,
      };

    case 'SET_MAIN_IMAGE':
      if (action.payload.type === 'existing') {
        return {
          ...state,
          existingImages: state.existingImages.map((img, index) => ({
            ...img,
            is_main: index === action.payload.index,
          })),
          hasUnsavedChanges: true,
        };
      }
      // Для новых изображений главное будет установлено при загрузке
      return state;

    case 'SET_LOCATION':
      return {
        ...state,
        location: action.payload,
        productData: {
          ...state.productData,
          has_individual_location: !action.payload?.useStorefrontLocation,
          individual_address: action.payload?.individualAddress || '',
          individual_latitude: action.payload?.latitude,
          individual_longitude: action.payload?.longitude,
          location_privacy: (action.payload?.privacyLevel || 'exact') as
            | 'exact'
            | 'street'
            | 'district'
            | 'city',
          show_on_map: action.payload?.showOnMap ?? true,
        },
        hasUnsavedChanges: true,
      };

    case 'SET_LOADING':
      return { ...state, isLoading: action.payload };

    case 'SET_SAVING':
      return { ...state, isSaving: action.payload };

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

    case 'SET_UNSAVED_CHANGES':
      return { ...state, hasUnsavedChanges: action.payload };

    case 'RESET':
      return initialState;

    default:
      return state;
  }
}

interface EditProductContextType {
  state: ProductEditState;
  dispatch: React.Dispatch<ProductEditAction>;
  // Хелперы
  goToStep: (step: number) => void;
  nextStep: () => void;
  prevStep: () => void;
  completeStep: (step: number) => void;
  loadProduct: (product: StorefrontProduct) => void;
  setCategory: (category: MarketplaceCategory) => void;
  setProductData: (data: Partial<UpdateProductRequest>) => void;
  setAttribute: (id: number, value: any) => void;
  setExistingImages: (images: ProductEditState['existingImages']) => void;
  setNewImages: (images: File[]) => void;
  addNewImage: (image: File) => void;
  removeNewImage: (index: number) => void;
  markImageForDeletion: (id: number) => void;
  restoreImage: (id: number) => void;
  setMainImage: (type: 'existing' | 'new', index: number) => void;
  setLocation: (location: ProductEditState['location']) => void;
  setError: (field: string, message: string) => void;
  clearError: (field: string) => void;
  clearAllErrors: () => void;
  validateStep: (step: number) => boolean;
  canGoToStep: (step: number) => boolean;
}

const EditProductContext = createContext<EditProductContextType | undefined>(
  undefined
);

export function EditProductProvider({ children }: { children: ReactNode }) {
  const [state, dispatch] = useReducer(productEditReducer, initialState);

  const goToStep = useCallback((step: number) => {
    dispatch({ type: 'SET_STEP', payload: step });
  }, []);

  const nextStep = useCallback(() => {
    dispatch({ type: 'SET_STEP', payload: Math.min(state.currentStep + 1, 5) });
  }, [state.currentStep]);

  const prevStep = useCallback(() => {
    dispatch({ type: 'SET_STEP', payload: Math.max(state.currentStep - 1, 0) });
  }, [state.currentStep]);

  const completeStep = useCallback((step: number) => {
    dispatch({ type: 'COMPLETE_STEP', payload: step });
  }, []);

  const loadProduct = useCallback((product: StorefrontProduct) => {
    dispatch({ type: 'LOAD_PRODUCT', payload: product });
  }, []);

  const setCategory = useCallback((category: MarketplaceCategory) => {
    dispatch({ type: 'SET_CATEGORY', payload: category });
  }, []);

  const setProductData = useCallback((data: Partial<UpdateProductRequest>) => {
    dispatch({ type: 'SET_PRODUCT_DATA', payload: data });
  }, []);

  const setAttribute = useCallback((id: number, value: any) => {
    dispatch({ type: 'SET_ATTRIBUTE', payload: { id, value } });
  }, []);

  const setExistingImages = useCallback(
    (images: ProductEditState['existingImages']) => {
      dispatch({ type: 'SET_EXISTING_IMAGES', payload: images });
    },
    []
  );

  const setNewImages = useCallback((images: File[]) => {
    dispatch({ type: 'SET_NEW_IMAGES', payload: images });
  }, []);

  const addNewImage = useCallback((image: File) => {
    dispatch({ type: 'ADD_NEW_IMAGE', payload: image });
  }, []);

  const removeNewImage = useCallback((index: number) => {
    dispatch({ type: 'REMOVE_NEW_IMAGE', payload: index });
  }, []);

  const markImageForDeletion = useCallback((id: number) => {
    dispatch({ type: 'MARK_IMAGE_FOR_DELETION', payload: id });
  }, []);

  const restoreImage = useCallback((id: number) => {
    dispatch({ type: 'RESTORE_IMAGE', payload: id });
  }, []);

  const setMainImage = useCallback(
    (type: 'existing' | 'new', index: number) => {
      dispatch({ type: 'SET_MAIN_IMAGE', payload: { type, index } });
    },
    []
  );

  const setLocation = useCallback((location: ProductEditState['location']) => {
    dispatch({ type: 'SET_LOCATION', payload: location });
  }, []);

  const setError = useCallback((field: string, message: string) => {
    dispatch({ type: 'SET_ERROR', payload: { field, message } });
  }, []);

  const clearError = useCallback((field: string) => {
    dispatch({ type: 'CLEAR_ERROR', payload: field });
  }, []);

  const clearAllErrors = useCallback(() => {
    dispatch({ type: 'CLEAR_ALL_ERRORS' });
  }, []);

  const validateStep = useCallback(
    (step: number): boolean => {
      switch (step) {
        case 0: // Категория
          return !!state.category;
        case 1: // Основная информация
          return !!(
            state.productData.name &&
            state.productData.description &&
            state.productData.price &&
            state.productData.price > 0
          );
        case 2: // Местоположение
          if (state.location?.useStorefrontLocation) {
            return true;
          }
          return !!(
            state.location?.latitude &&
            state.location?.longitude &&
            state.location?.individualAddress
          );
        case 3: // Атрибуты
          return true; // Атрибуты опциональны
        case 4: // Фотографии
          return (
            state.existingImages.filter(
              (img) => !state.imagesToDelete.includes(img.id)
            ).length > 0 || state.newImages.length > 0
          );
        case 5: // Превью
          return true;
        default:
          return false;
      }
    },
    [
      state.category,
      state.productData,
      state.location,
      state.existingImages,
      state.imagesToDelete,
      state.newImages,
    ]
  );

  const canGoToStep = useCallback(
    (step: number): boolean => {
      if (step <= state.currentStep) return true;

      for (let i = 0; i < step; i++) {
        if (!validateStep(i)) return false;
      }
      return true;
    },
    [state.currentStep, validateStep]
  );

  const value: EditProductContextType = {
    state,
    dispatch,
    goToStep,
    nextStep,
    prevStep,
    completeStep,
    loadProduct,
    setCategory,
    setProductData,
    setAttribute,
    setExistingImages,
    setNewImages,
    addNewImage,
    removeNewImage,
    markImageForDeletion,
    restoreImage,
    setMainImage,
    setLocation,
    setError,
    clearError,
    clearAllErrors,
    validateStep,
    canGoToStep,
  };

  return (
    <EditProductContext.Provider value={value}>
      {children}
    </EditProductContext.Provider>
  );
}

export function useEditProduct(): EditProductContextType {
  const context = useContext(EditProductContext);
  if (context === undefined) {
    throw new Error(
      'useEditProduct must be used within an EditProductProvider'
    );
  }
  return context;
}
