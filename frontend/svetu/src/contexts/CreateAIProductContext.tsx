'use client';

import React, { createContext, useContext, useReducer, ReactNode } from 'react';

type View = 'upload' | 'process' | 'enhance' | 'publish';

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
  };

  // Настройки
  useStorefrontLocation: boolean;
  locationPrivacyLevel: 'exact' | 'street' | 'city';
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
  },
  useStorefrontLocation: true,
  locationPrivacyLevel: 'exact',
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
