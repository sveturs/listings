import { useMemo } from 'react';
import { FormData, FormErrors, FormMode } from './useAuthForm';
import {
  VALIDATION_PATTERNS,
  VALIDATION_CONSTRAINTS,
  VALIDATION_MESSAGES,
  type ValidationErrorKey,
} from '@/constants/validation';

// Validation functions
const validateEmail = (email: string): ValidationErrorKey | undefined => {
  if (!email.trim()) {
    return VALIDATION_MESSAGES.EMAIL_REQUIRED;
  }
  if (!VALIDATION_PATTERNS.EMAIL.test(email)) {
    return VALIDATION_MESSAGES.EMAIL_INVALID;
  }
  return undefined;
};

const validatePassword = (password: string): ValidationErrorKey | undefined => {
  if (!password) {
    return VALIDATION_MESSAGES.PASSWORD_REQUIRED;
  }
  if (password.length < VALIDATION_CONSTRAINTS.PASSWORD_MIN_LENGTH) {
    return VALIDATION_MESSAGES.PASSWORD_TOO_SHORT;
  }
  return undefined;
};

const validateConfirmPassword = (
  password: string,
  confirmPassword: string
): ValidationErrorKey | undefined => {
  if (!confirmPassword) {
    return VALIDATION_MESSAGES.PASSWORD_REQUIRED;
  }
  if (password !== confirmPassword) {
    return VALIDATION_MESSAGES.PASSWORDS_DO_NOT_MATCH;
  }
  return undefined;
};

const validateName = (name: string): ValidationErrorKey | undefined => {
  if (!name.trim()) {
    return VALIDATION_MESSAGES.NAME_REQUIRED;
  }
  if (name.trim().length < VALIDATION_CONSTRAINTS.NAME_MIN_LENGTH) {
    return VALIDATION_MESSAGES.NAME_TOO_SHORT;
  }
  if (name.trim().length > VALIDATION_CONSTRAINTS.NAME_MAX_LENGTH) {
    return VALIDATION_MESSAGES.NAME_TOO_LONG;
  }
  return undefined;
};

const validatePhone = (phone: string): ValidationErrorKey | undefined => {
  if (phone && !VALIDATION_PATTERNS.PHONE.test(phone)) {
    return VALIDATION_MESSAGES.PHONE_INVALID;
  }
  return undefined;
};

// Main validation function
const validateFormData = (formData: FormData, mode: FormMode): FormErrors => {
  const errors: FormErrors = {};

  // Email validation (required for both modes)
  const emailError = validateEmail(formData.email);
  if (emailError) {
    errors.email = emailError;
  }

  // Password validation (required for both modes)
  const passwordError = validatePassword(formData.password);
  if (passwordError) {
    errors.password = passwordError;
  }

  // Mode-specific validation
  if (mode === 'register') {
    // Name validation (required for register)
    const nameError = validateName(formData.name);
    if (nameError) {
      errors.name = nameError;
    }

    // Confirm password validation (required for register)
    const confirmPasswordError = validateConfirmPassword(
      formData.password,
      formData.confirmPassword
    );
    if (confirmPasswordError) {
      errors.confirmPassword = confirmPasswordError;
    }

    // Phone validation (optional for register)
    const phoneError = validatePhone(formData.phone);
    if (phoneError) {
      errors.phone = phoneError;
    }
  }

  return errors;
};

// Field-specific validation for real-time feedback
const validateField = (
  field: keyof FormData,
  value: string,
  formData: FormData,
  mode: FormMode
): ValidationErrorKey | undefined => {
  switch (field) {
    case 'email':
      return validateEmail(value);
    case 'password':
      return validatePassword(value);
    case 'confirmPassword':
      return mode === 'register'
        ? validateConfirmPassword(formData.password, value)
        : undefined;
    case 'name':
      return mode === 'register' ? validateName(value) : undefined;
    case 'phone':
      return mode === 'register' ? validatePhone(value) : undefined;
    default:
      return undefined;
  }
};

export const useFormValidation = (formData: FormData, mode: FormMode) => {
  // Validate all fields
  const allErrors = useMemo(() => {
    return validateFormData(formData, mode);
  }, [formData, mode]);

  // Check if form is valid
  const isFormValid = useMemo(() => {
    return Object.keys(allErrors).length === 0;
  }, [allErrors]);

  // Get validation for specific field
  const getFieldValidation = useMemo(() => {
    return (field: keyof FormData) => {
      const error = validateField(field, formData[field], formData, mode);
      const isValid = !error && formData[field].length > 0;
      return { error, isValid };
    };
  }, [formData, mode]);

  // Check if specific field is valid
  const isFieldValid = useMemo(() => {
    return (field: keyof FormData) => {
      const { isValid } = getFieldValidation(field);
      return isValid;
    };
  }, [getFieldValidation]);

  // Get required fields for current mode
  const getRequiredFields = useMemo(() => {
    const baseFields = ['email', 'password'] as const;
    if (mode === 'register') {
      return [...baseFields, 'name', 'confirmPassword'] as const;
    }
    return baseFields;
  }, [mode]);

  // Check if all required fields are filled
  const areRequiredFieldsFilled = useMemo(() => {
    const requiredFields = getRequiredFields;
    return requiredFields.every((field) => formData[field].trim().length > 0);
  }, [formData, getRequiredFields]);

  // Check if form can be submitted
  const canSubmit = useMemo(() => {
    return areRequiredFieldsFilled && isFormValid;
  }, [areRequiredFieldsFilled, isFormValid]);

  return {
    allErrors,
    isFormValid,
    getFieldValidation,
    isFieldValid,
    getRequiredFields,
    areRequiredFieldsFilled,
    canSubmit,
    validateField: (
      field: keyof FormData,
      value: string
    ): ValidationErrorKey | undefined =>
      validateField(field, value, formData, mode),
  };
};
