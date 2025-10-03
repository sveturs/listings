'use client';

import React, { createContext, useContext, useReducer, ReactNode } from 'react';

type View = 'upload' | 'process' | 'enhance' | 'variants' | 'publish';

// Типы для вариантов (из CreateProductContext)
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

interface AIProductState {
  // Workflow
  currentView: View;
  isProcessing: boolean;
  error: string | null;

  // Изображения
  images: string[]; // blob URLs
  imageFiles: File[];

  // AI данные
  aiData: {
    title: string;
    titleVariants: string[];
    selectedTitleIndex: number;
    description: string;
    category: string;
    categoryId: number;
    categoryProbabilities: Array<{ name: string; probability: number }>;
    price: number;
    priceRange: { min: number; max: number };
    currency: string;
    attributes: Record<string, any>;
    suggestedVariants: Array<{
      attributes: Record<string, string>;
      sku?: string;
      price?: number;
      stockQuantity: number;
    }>;
    stockQuantity: number;
    condition: 'new' | 'used' | 'refurbished';
    keywords: string[];
    translations: Record<string, { title: string; description: string }>;
    location: {
      address: string;
      latitude: number;
      longitude: number;
      city: string;
      region: string;
      source: 'exif' | 'storefront' | 'manual';
    } | null;

    // НОВОЕ: Варианты товара
    hasVariants: boolean;
    variants: ProductVariantCreate[];
    variantSettings: VariantSettings;
  };

  // Настройки
  useStorefrontLocation: boolean;
  locationPrivacyLevel: 'exact' | 'street' | 'district' | 'city';
  showOnMap: boolean;
  previewLanguage: string;
}

type AIProductAction =
  | { type: 'SET_VIEW'; payload: View }
  | { type: 'SET_PROCESSING'; payload: boolean }
  | { type: 'SET_ERROR'; payload: string | null }
  | { type: 'SET_IMAGES'; payload: { urls: string[]; files: File[] } }
  | { type: 'SET_AI_DATA'; payload: Partial<AIProductState['aiData']> }
  | { type: 'SELECT_TITLE_VARIANT'; payload: number }
  | { type: 'SET_USE_STOREFRONT_LOCATION'; payload: boolean }
  | {
      type: 'SET_LOCATION_PRIVACY_LEVEL';
      payload: 'exact' | 'street' | 'district' | 'city';
    }
  | { type: 'SET_SHOW_ON_MAP'; payload: boolean }
  | { type: 'SET_HAS_VARIANTS'; payload: boolean }
  | { type: 'SET_VARIANTS'; payload: ProductVariantCreate[] }
  | { type: 'ADD_VARIANT'; payload: ProductVariantCreate }
  | {
      type: 'UPDATE_VARIANT';
      payload: { index: number; variant: ProductVariantCreate };
    }
  | { type: 'REMOVE_VARIANT'; payload: number }
  | { type: 'SET_VARIANT_SETTINGS'; payload: Partial<VariantSettings> }
  | { type: 'RESET' };

const initialState: AIProductState = {
  currentView: 'upload',
  isProcessing: false,
  error: null,
  images: [],
  imageFiles: [],
  aiData: {
    title: '',
    titleVariants: [],
    selectedTitleIndex: 0,
    description: '',
    category: '',
    categoryId: 0,
    categoryProbabilities: [],
    price: 0,
    priceRange: { min: 0, max: 0 },
    currency: 'RSD',
    attributes: {},
    suggestedVariants: [],
    stockQuantity: 0,
    condition: 'new',
    keywords: [],
    translations: {},
    location: null,

    // Варианты
    hasVariants: false,
    variants: [],
    variantSettings: {
      track_inventory: true,
      continue_selling: false,
      require_shipping: true,
      taxable_product: true,
      selected_attributes: [],
    },
  },
  useStorefrontLocation: true,
  locationPrivacyLevel: 'exact',
  showOnMap: true,
  previewLanguage: 'ru',
};

function aiProductReducer(
  state: AIProductState,
  action: AIProductAction
): AIProductState {
  switch (action.type) {
    case 'SET_VIEW':
      return { ...state, currentView: action.payload };

    case 'SET_PROCESSING':
      return { ...state, isProcessing: action.payload };

    case 'SET_ERROR':
      return { ...state, error: action.payload };

    case 'SET_IMAGES':
      return {
        ...state,
        images: action.payload.urls,
        imageFiles: action.payload.files,
      };

    case 'SET_AI_DATA':
      return {
        ...state,
        aiData: { ...state.aiData, ...action.payload },
      };

    case 'SELECT_TITLE_VARIANT':
      return {
        ...state,
        aiData: {
          ...state.aiData,
          selectedTitleIndex: action.payload,
        },
      };

    case 'SET_USE_STOREFRONT_LOCATION':
      return { ...state, useStorefrontLocation: action.payload };

    case 'SET_LOCATION_PRIVACY_LEVEL':
      return { ...state, locationPrivacyLevel: action.payload };

    case 'SET_SHOW_ON_MAP':
      return { ...state, showOnMap: action.payload };

    case 'SET_HAS_VARIANTS':
      return {
        ...state,
        aiData: {
          ...state.aiData,
          hasVariants: action.payload,
          variants: action.payload ? state.aiData.variants : [],
        },
      };

    case 'SET_VARIANTS':
      return {
        ...state,
        aiData: { ...state.aiData, variants: action.payload },
      };

    case 'ADD_VARIANT':
      return {
        ...state,
        aiData: {
          ...state.aiData,
          variants: [...state.aiData.variants, action.payload],
        },
      };

    case 'UPDATE_VARIANT': {
      const updatedVariants = [...state.aiData.variants];
      updatedVariants[action.payload.index] = action.payload.variant;
      return {
        ...state,
        aiData: { ...state.aiData, variants: updatedVariants },
      };
    }

    case 'REMOVE_VARIANT':
      return {
        ...state,
        aiData: {
          ...state.aiData,
          variants: state.aiData.variants.filter(
            (_, i) => i !== action.payload
          ),
        },
      };

    case 'SET_VARIANT_SETTINGS':
      return {
        ...state,
        aiData: {
          ...state.aiData,
          variantSettings: {
            ...state.aiData.variantSettings,
            ...action.payload,
          },
        },
      };

    case 'RESET':
      return initialState;

    default:
      return state;
  }
}

interface CreateAIProductContextType {
  state: AIProductState;
  dispatch: React.Dispatch<AIProductAction>;
  // Helper methods
  setView: (view: View) => void;
  setProcessing: (isProcessing: boolean) => void;
  setError: (error: string | null) => void;
  setImages: (urls: string[], files: File[]) => void;
  setAIData: (data: Partial<AIProductState['aiData']>) => void;
  selectTitleVariant: (index: number) => void;
  reset: () => void;
  // Location methods
  setUseStorefrontLocation: (use: boolean) => void;
  setLocationPrivacyLevel: (
    level: 'exact' | 'street' | 'district' | 'city'
  ) => void;
  setShowOnMap: (show: boolean) => void;
  // Variant methods
  setHasVariants: (hasVariants: boolean) => void;
  setVariants: (variants: ProductVariantCreate[]) => void;
  addVariant: (variant: ProductVariantCreate) => void;
  updateVariant: (index: number, variant: ProductVariantCreate) => void;
  removeVariant: (index: number) => void;
  setVariantSettings: (settings: Partial<VariantSettings>) => void;
}

const CreateAIProductContext = createContext<
  CreateAIProductContextType | undefined
>(undefined);

export function CreateAIProductProvider({ children }: { children: ReactNode }) {
  const [state, dispatch] = useReducer(aiProductReducer, initialState);

  const setView = (view: View) => dispatch({ type: 'SET_VIEW', payload: view });
  const setProcessing = (isProcessing: boolean) =>
    dispatch({ type: 'SET_PROCESSING', payload: isProcessing });
  const setError = (error: string | null) =>
    dispatch({ type: 'SET_ERROR', payload: error });
  const setImages = (urls: string[], files: File[]) =>
    dispatch({ type: 'SET_IMAGES', payload: { urls, files } });
  const setAIData = (data: Partial<AIProductState['aiData']>) =>
    dispatch({ type: 'SET_AI_DATA', payload: data });
  const selectTitleVariant = (index: number) =>
    dispatch({ type: 'SELECT_TITLE_VARIANT', payload: index });
  const reset = () => dispatch({ type: 'RESET' });

  // Location methods
  const setUseStorefrontLocation = (use: boolean) =>
    dispatch({ type: 'SET_USE_STOREFRONT_LOCATION', payload: use });
  const setLocationPrivacyLevel = (
    level: 'exact' | 'street' | 'district' | 'city'
  ) => dispatch({ type: 'SET_LOCATION_PRIVACY_LEVEL', payload: level });
  const setShowOnMap = (show: boolean) =>
    dispatch({ type: 'SET_SHOW_ON_MAP', payload: show });

  // Variant methods
  const setHasVariants = (hasVariants: boolean) =>
    dispatch({ type: 'SET_HAS_VARIANTS', payload: hasVariants });
  const setVariants = (variants: ProductVariantCreate[]) =>
    dispatch({ type: 'SET_VARIANTS', payload: variants });
  const addVariant = (variant: ProductVariantCreate) =>
    dispatch({ type: 'ADD_VARIANT', payload: variant });
  const updateVariant = (index: number, variant: ProductVariantCreate) =>
    dispatch({ type: 'UPDATE_VARIANT', payload: { index, variant } });
  const removeVariant = (index: number) =>
    dispatch({ type: 'REMOVE_VARIANT', payload: index });
  const setVariantSettings = (settings: Partial<VariantSettings>) =>
    dispatch({ type: 'SET_VARIANT_SETTINGS', payload: settings });

  const value: CreateAIProductContextType = {
    state,
    dispatch,
    setView,
    setProcessing,
    setError,
    setImages,
    setAIData,
    selectTitleVariant,
    reset,
    setUseStorefrontLocation,
    setLocationPrivacyLevel,
    setShowOnMap,
    setHasVariants,
    setVariants,
    addVariant,
    updateVariant,
    removeVariant,
    setVariantSettings,
  };

  return (
    <CreateAIProductContext.Provider value={value}>
      {children}
    </CreateAIProductContext.Provider>
  );
}

export function useCreateAIProduct(): CreateAIProductContextType {
  const context = useContext(CreateAIProductContext);
  if (context === undefined) {
    throw new Error(
      'useCreateAIProduct must be used within CreateAIProductProvider'
    );
  }
  return context;
}
