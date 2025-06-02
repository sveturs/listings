import { useReducer, useCallback } from 'react';
import type { ValidationErrorKey } from '@/constants/validation';

export type FormMode = 'login' | 'register';

export interface FormData {
  email: string;
  password: string;
  confirmPassword: string;
  name: string;
  phone: string;
}

export interface FormErrors {
  email?: ValidationErrorKey;
  password?: ValidationErrorKey;
  confirmPassword?: ValidationErrorKey;
  name?: ValidationErrorKey;
  phone?: ValidationErrorKey;
  general?: ValidationErrorKey;
}

interface AuthFormState {
  mode: FormMode;
  formData: FormData;
  errors: FormErrors;
  ui: {
    isLoading: boolean;
    successMessage: string;
  };
}

type AuthFormAction =
  | { type: 'SET_MODE'; payload: FormMode }
  | { type: 'SET_FIELD'; payload: { field: keyof FormData; value: string } }
  | { type: 'SET_ERRORS'; payload: FormErrors }
  | {
      type: 'SET_ERROR';
      payload: { field: keyof FormErrors; message: ValidationErrorKey };
    }
  | { type: 'CLEAR_ERROR'; payload: keyof FormErrors }
  | { type: 'SET_LOADING'; payload: boolean }
  | { type: 'SET_SUCCESS'; payload: string }
  | { type: 'RESET_FORM' }
  | { type: 'CLEAR_SUCCESS' };

const initialFormData: FormData = {
  email: '',
  password: '',
  confirmPassword: '',
  name: '',
  phone: '',
};

const initialAuthFormState: AuthFormState = {
  mode: 'login',
  formData: initialFormData,
  errors: {},
  ui: {
    isLoading: false,
    successMessage: '',
  },
};

const authFormReducer = (
  state: AuthFormState,
  action: AuthFormAction
): AuthFormState => {
  switch (action.type) {
    case 'SET_MODE':
      return {
        ...state,
        mode: action.payload,
        errors: {},
        formData: {
          ...initialFormData,
          email: state.formData.email, // Keep email when switching modes
        },
        ui: { ...state.ui, successMessage: '' },
      };

    case 'SET_FIELD':
      return {
        ...state,
        formData: {
          ...state.formData,
          [action.payload.field]: action.payload.value,
        },
        errors: { ...state.errors, [action.payload.field]: undefined },
      };

    case 'SET_ERRORS':
      return { ...state, errors: action.payload };

    case 'SET_ERROR':
      return {
        ...state,
        errors: {
          ...state.errors,
          [action.payload.field]: action.payload.message,
        },
      };

    case 'CLEAR_ERROR':
      const newErrors = { ...state.errors };
      delete newErrors[action.payload];
      return { ...state, errors: newErrors };

    case 'SET_LOADING':
      return { ...state, ui: { ...state.ui, isLoading: action.payload } };

    case 'SET_SUCCESS':
      return {
        ...state,
        ui: { ...state.ui, successMessage: action.payload },
        errors: {},
      };

    case 'CLEAR_SUCCESS':
      return { ...state, ui: { ...state.ui, successMessage: '' } };

    case 'RESET_FORM':
      return { ...initialAuthFormState, mode: state.mode };

    default:
      return state;
  }
};

export const useAuthForm = (initialMode: FormMode = 'login') => {
  const [state, dispatch] = useReducer(authFormReducer, {
    ...initialAuthFormState,
    mode: initialMode,
  });

  const setField = useCallback((field: keyof FormData, value: string) => {
    dispatch({ type: 'SET_FIELD', payload: { field, value } });
  }, []);

  const setMode = useCallback((mode: FormMode) => {
    dispatch({ type: 'SET_MODE', payload: mode });
  }, []);

  const setErrors = useCallback((errors: FormErrors) => {
    dispatch({ type: 'SET_ERRORS', payload: errors });
  }, []);

  const setError = useCallback(
    (field: keyof FormErrors, message: ValidationErrorKey) => {
      dispatch({ type: 'SET_ERROR', payload: { field, message } });
    },
    []
  );

  const clearError = useCallback((field: keyof FormErrors) => {
    dispatch({ type: 'CLEAR_ERROR', payload: field });
  }, []);

  const setLoading = useCallback((loading: boolean) => {
    dispatch({ type: 'SET_LOADING', payload: loading });
  }, []);

  const setSuccess = useCallback((message: string) => {
    dispatch({ type: 'SET_SUCCESS', payload: message });
  }, []);

  const clearSuccess = useCallback(() => {
    dispatch({ type: 'CLEAR_SUCCESS' });
  }, []);

  const resetForm = useCallback(() => {
    dispatch({ type: 'RESET_FORM' });
  }, []);

  const getFieldValue = useCallback(
    (field: keyof FormData) => state.formData[field],
    [state.formData]
  );

  const getFieldError = useCallback(
    (field: keyof FormErrors) => state.errors[field],
    [state.errors]
  );

  const hasErrors = useCallback(() => {
    return Object.keys(state.errors).length > 0;
  }, [state.errors]);

  const isFieldValid = useCallback(
    (field: keyof FormData) => {
      const value = state.formData[field];
      const error = state.errors[field];
      return value.length > 0 && !error;
    },
    [state.formData, state.errors]
  );

  return {
    // State
    mode: state.mode,
    formData: state.formData,
    errors: state.errors,
    isLoading: state.ui.isLoading,
    successMessage: state.ui.successMessage,

    // Actions
    setField,
    setMode,
    setErrors,
    setError,
    clearError,
    setLoading,
    setSuccess,
    clearSuccess,
    resetForm,

    // Computed
    getFieldValue,
    getFieldError,
    hasErrors,
    isFieldValid,
  };
};
