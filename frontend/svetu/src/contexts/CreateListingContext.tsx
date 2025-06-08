'use client';

import { createContext, useContext, useReducer, ReactNode } from 'react';

export interface CreateListingState {
  // Базовая информация
  category?: {
    id: number;
    name: string;
  };
  title: string;
  description: string;
  price: number;
  currency: 'RSD' | 'EUR' | 'HRK' | 'MKD';
  condition: 'new' | 'used' | 'refurbished';

  // Региональная система доверия
  trust: {
    phoneVerified: boolean;
    preferredMeetingType: 'personal' | 'pickup' | 'delivery';
    meetingLocations: string[];
    availableHours: string;
    localReputation: number;
  };

  // Платежи и доставка
  payment: {
    methods: string[];
    codEnabled: boolean;
    codPrice: number;
    personalMeeting: boolean;
    deliveryOptions: string[];
    negotiablePrice: boolean;
    bundleDeals: boolean;
  };

  // Локализация
  localization: {
    script: 'cyrillic' | 'latin' | 'mixed';
    language: 'sr' | 'hr' | 'sl' | 'mk';
    traditionalUnits: boolean;
    regionalPhrases: string[];
  };

  // Pijaca 2.0 features
  pijaca: {
    vendorStallStyle: string;
    neighborhoodGroup?: string;
    regularCustomers: boolean;
    traditionalStyle: boolean;
  };

  // Изображения
  images: string[];
  mainImageIndex: number;

  // Местоположение
  location?: {
    latitude: number;
    longitude: number;
    address: string;
    city: string;
    region: string;
    country: string;
  };

  // Атрибуты
  attributes: Record<string, any>;

  // Метаданные
  isPublished: boolean;
  isDraft: boolean;
  lastSaved?: Date;
}

type CreateListingAction =
  | { type: 'SET_CATEGORY'; payload: { id: number; name: string } }
  | {
      type: 'SET_BASIC_INFO';
      payload: Partial<
        Pick<
          CreateListingState,
          'title' | 'description' | 'price' | 'currency' | 'condition'
        >
      >;
    }
  | { type: 'SET_TRUST'; payload: Partial<CreateListingState['trust']> }
  | { type: 'SET_PAYMENT'; payload: Partial<CreateListingState['payment']> }
  | {
      type: 'SET_LOCALIZATION';
      payload: Partial<CreateListingState['localization']>;
    }
  | { type: 'SET_PIJACA'; payload: Partial<CreateListingState['pijaca']> }
  | { type: 'SET_IMAGES'; payload: string[] }
  | { type: 'SET_MAIN_IMAGE'; payload: number }
  | { type: 'SET_LOCATION'; payload: CreateListingState['location'] }
  | { type: 'SET_ATTRIBUTES'; payload: Record<string, any> }
  | { type: 'SAVE_DRAFT' }
  | { type: 'PUBLISH' }
  | { type: 'RESET' };

const initialState: CreateListingState = {
  title: '',
  description: '',
  price: 0,
  currency: 'RSD',
  condition: 'used',

  trust: {
    phoneVerified: false,
    preferredMeetingType: 'personal',
    meetingLocations: [],
    availableHours: '',
    localReputation: 0,
  },

  payment: {
    methods: ['cod'],
    codEnabled: true,
    codPrice: 250,
    personalMeeting: true,
    deliveryOptions: [],
    negotiablePrice: false,
    bundleDeals: false,
  },

  localization: {
    script: 'cyrillic',
    language: 'sr',
    traditionalUnits: false,
    regionalPhrases: [],
  },

  pijaca: {
    vendorStallStyle: 'traditional',
    regularCustomers: false,
    traditionalStyle: true,
  },

  images: [],
  mainImageIndex: 0,
  attributes: {},
  isPublished: false,
  isDraft: false,
};

function createListingReducer(
  state: CreateListingState,
  action: CreateListingAction
): CreateListingState {
  switch (action.type) {
    case 'SET_CATEGORY':
      return { ...state, category: action.payload };

    case 'SET_BASIC_INFO':
      return { ...state, ...action.payload };

    case 'SET_TRUST':
      return { ...state, trust: { ...state.trust, ...action.payload } };

    case 'SET_PAYMENT':
      return { ...state, payment: { ...state.payment, ...action.payload } };

    case 'SET_LOCALIZATION':
      return {
        ...state,
        localization: { ...state.localization, ...action.payload },
      };

    case 'SET_PIJACA':
      return { ...state, pijaca: { ...state.pijaca, ...action.payload } };

    case 'SET_IMAGES':
      return { ...state, images: action.payload };

    case 'SET_MAIN_IMAGE':
      return { ...state, mainImageIndex: action.payload };

    case 'SET_LOCATION':
      return { ...state, location: action.payload };

    case 'SET_ATTRIBUTES':
      return {
        ...state,
        attributes: { ...state.attributes, ...action.payload },
      };

    case 'SAVE_DRAFT':
      return { ...state, isDraft: true, lastSaved: new Date() };

    case 'PUBLISH':
      return {
        ...state,
        isPublished: true,
        isDraft: false,
        lastSaved: new Date(),
      };

    case 'RESET':
      return initialState;

    default:
      return state;
  }
}

interface CreateListingContextType {
  state: CreateListingState;
  dispatch: React.Dispatch<CreateListingAction>;

  // Вспомогательные функции
  setCategory: (category: { id: number; name: string }) => void;
  setBasicInfo: (
    info: Partial<
      Pick<
        CreateListingState,
        'title' | 'description' | 'price' | 'currency' | 'condition'
      >
    >
  ) => void;
  setTrust: (trust: Partial<CreateListingState['trust']>) => void;
  setPayment: (payment: Partial<CreateListingState['payment']>) => void;
  setLocalization: (
    localization: Partial<CreateListingState['localization']>
  ) => void;
  setPijaca: (pijaca: Partial<CreateListingState['pijaca']>) => void;
  saveDraft: () => void;
  publish: () => void;
  reset: () => void;
}

const CreateListingContext = createContext<
  CreateListingContextType | undefined
>(undefined);

export function CreateListingProvider({ children }: { children: ReactNode }) {
  const [state, dispatch] = useReducer(createListingReducer, initialState);

  const contextValue: CreateListingContextType = {
    state,
    dispatch,

    setCategory: (category) =>
      dispatch({ type: 'SET_CATEGORY', payload: category }),
    setBasicInfo: (info) => dispatch({ type: 'SET_BASIC_INFO', payload: info }),
    setTrust: (trust) => dispatch({ type: 'SET_TRUST', payload: trust }),
    setPayment: (payment) =>
      dispatch({ type: 'SET_PAYMENT', payload: payment }),
    setLocalization: (localization) =>
      dispatch({ type: 'SET_LOCALIZATION', payload: localization }),
    setPijaca: (pijaca) => dispatch({ type: 'SET_PIJACA', payload: pijaca }),
    saveDraft: () => dispatch({ type: 'SAVE_DRAFT' }),
    publish: () => dispatch({ type: 'PUBLISH' }),
    reset: () => dispatch({ type: 'RESET' }),
  };

  return (
    <CreateListingContext.Provider value={contextValue}>
      {children}
    </CreateListingContext.Provider>
  );
}

export function useCreateListing() {
  const context = useContext(CreateListingContext);
  if (context === undefined) {
    throw new Error(
      'useCreateListing must be used within a CreateListingProvider'
    );
  }
  return context;
}
