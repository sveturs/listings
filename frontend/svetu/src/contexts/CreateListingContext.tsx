'use client';

import {
  createContext,
  useContext,
  useReducer,
  ReactNode,
  // useEffect, // ВРЕМЕННО ОТКЛЮЧЕНО
} from 'react';
import { useAuth } from '@/contexts/AuthContext';
// import { useListingDraft } from '@/hooks/useListingDraft'; // ВРЕМЕННО ОТКЛЮЧЕНО

export interface CreateListingState {
  // Базовая информация
  category?: {
    id: number;
    name: string;
    slug: string;
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
    addressMultilingual?: {
      sr: string;
      en: string;
      ru: string;
    };
  };

  // Атрибуты
  attributes: Record<string, any>;

  // Переводы
  translations?: Record<string, Record<string, string>>;

  // Язык оригинала
  originalLanguage?: string;

  // Метаданные
  isPublished: boolean;
  isDraft: boolean;
  lastSaved?: Date;
}

type CreateListingAction =
  | {
      type: 'SET_CATEGORY';
      payload: { id: number; name: string; slug: string };
    }
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
  | { type: 'RESET' }
  | { type: 'LOAD_FROM_DRAFT'; payload: CreateListingState };

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
  console.log(
    'CreateListingReducer - Action:',
    action.type,
    'payload' in action ? action.payload : 'no payload'
  );
  switch (action.type) {
    case 'SET_CATEGORY':
      const newStateWithCategory = { ...state, category: action.payload };
      console.log(
        'CreateListingReducer - SET_CATEGORY result:',
        newStateWithCategory.category
      );
      return newStateWithCategory;

    case 'SET_BASIC_INFO':
      const newStateWithBasicInfo = { ...state, ...action.payload };
      console.log('CreateListingReducer - SET_BASIC_INFO result:', {
        title: newStateWithBasicInfo.title,
        description: newStateWithBasicInfo.description,
        price: newStateWithBasicInfo.price,
        category: newStateWithBasicInfo.category,
      });
      return newStateWithBasicInfo;

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

    case 'LOAD_FROM_DRAFT':
      return action.payload;

    default:
      return state;
  }
}

interface CreateListingContextType {
  state: CreateListingState;
  dispatch: React.Dispatch<CreateListingAction>;

  // Вспомогательные функции
  setCategory: (category: { id: number; name: string; slug: string }) => void;
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

  // Черновики
  draftId?: string;
  isLoadingDraft: boolean;
  isSavingDraft: boolean;
  hasUnsavedChanges: boolean;
  lastSavedAt?: Date;
}

const CreateListingContext = createContext<
  CreateListingContextType | undefined
>(undefined);

interface CreateListingProviderProps {
  children: ReactNode;
  draftId?: string;
}

export function CreateListingProvider({
  children,
  draftId: _draftId,
}: CreateListingProviderProps) {
  const {} = useAuth();
  const [state, dispatch] = useReducer(createListingReducer, initialState);

  // ВРЕМЕННО ОТКЛЮЧЕНО: Интеграция с системой черновиков
  // const {
  //   draft,
  //   isLoading: isLoadingDraft,
  //   isSaving: isSavingDraft,
  //   hasChanges: hasUnsavedChanges,
  //   lastSaved: lastSavedAt,
  //   saveDraft: saveDraftToStorage,
  //   updateFormData,
  //   updateAttributes,
  //   updateStep,
  // } = useListingDraft({
  //   draftId,
  //   autoSaveInterval: 30000, // Автосохранение каждые 30 секунд
  // });

  // Заглушки для отключенной системы черновиков
  const _draft = null;
  const isLoadingDraft = false;
  const isSavingDraft = false;
  const hasUnsavedChanges = false;
  const lastSavedAt = null;
  const saveDraftToStorage = () => {};
  const _updateFormData = () => {};
  const _updateAttributes = () => {};
  const _updateStep = () => {};

  // ВРЕМЕННО ОТКЛЮЧЕНО: Загрузка и синхронизация черновиков
  // useEffect(() => {
  //   if (draft && !isLoadingDraft) {
  //     // Преобразуем данные черновика в состояние контекста
  //     const draftState: CreateListingState = {
  //       ...initialState,
  //       ...draft.formData,
  //       attributes: draft.attributes,
  //       isDraft: true,
  //       lastSaved: new Date(draft.metadata.updatedAt),
  //     };

  //     dispatch({ type: 'LOAD_FROM_DRAFT', payload: draftState });

  //     // Также загружаем текущий шаг если сохранен
  //     if (draft.metadata.currentStep && updateStep) {
  //       updateStep(draft.metadata.currentStep);
  //     }
  //   }
  // }, [draft, isLoadingDraft, updateStep]);

  // // Синхронизируем изменения состояния с черновиком
  // useEffect(() => {
  //   if (!draft || isLoadingDraft) return;

  //   // Создаем debounced функцию для обновления
  //   const timeoutId = setTimeout(() => {
  //     // Обновляем данные формы в черновике
  //     updateFormData({
  //       title: state.title,
  //       description: state.description,
  //       price: state.price,
  //       currency: state.currency,
  //       condition: state.condition,
  //       category: state.category,
  //       trust: state.trust,
  //       payment: state.payment,
  //       localization: state.localization,
  //       pijaca: state.pijaca,
  //       location: state.location,
  //     });

  //     // Обновляем атрибуты
  //     updateAttributes(state.attributes);
  //   }, 500); // Задержка 500мс для избежания частых обновлений

  //   return () => clearTimeout(timeoutId);
  // }, [
  //   state.title,
  //   state.description,
  //   state.price,
  //   state.currency,
  //   state.condition,
  //   state.category?.id,
  //   state.trust.phoneVerified,
  //   state.trust.preferredMeetingType,
  //   state.payment.methods.join(','),
  //   state.payment.codEnabled,
  //   state.localization.script,
  //   state.pijaca.vendorStallStyle,
  //   state.location?.city,
  //   JSON.stringify(state.attributes),
  //   draft,
  //   isLoadingDraft,
  //   updateFormData,
  //   updateAttributes,
  // ]);

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
    saveDraft: () => {
      dispatch({ type: 'SAVE_DRAFT' });
      saveDraftToStorage();
    },
    publish: () => dispatch({ type: 'PUBLISH' }),
    reset: () => dispatch({ type: 'RESET' }),

    // Состояние черновика
    draftId: undefined, // _draft?.metadata.id, // ВРЕМЕННО ОТКЛЮЧЕНО
    isLoadingDraft,
    isSavingDraft,
    hasUnsavedChanges,
    lastSavedAt: lastSavedAt || undefined,
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
