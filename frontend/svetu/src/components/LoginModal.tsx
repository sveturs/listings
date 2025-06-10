import React, { useCallback, useEffect, useMemo } from 'react';
import { useTranslations } from 'next-intl';
import { useDebounce } from '@/hooks/useDebounce';
import { useAuthForm, FormMode, FormData } from '@/hooks/useAuthForm';
import { useFormValidation } from '@/hooks/useFormValidation';
import { useAuth } from '@/contexts/AuthContext';
import { AuthService } from '@/services/auth';
import {
  VALIDATION_MESSAGES,
  ValidationErrorKey,
} from '@/constants/validation';
import LoginForm from './auth/LoginForm';
import RegisterForm from './auth/RegisterForm';

// Map server error messages to client error keys
const mapServerErrorToClientError = (
  serverError: string
): ValidationErrorKey => {
  const errorMap: Record<string, ValidationErrorKey> = {
    'users.register.error.name_required':
      VALIDATION_MESSAGES.AUTH_NAME_REQUIRED,
    'users.register.error.email_required':
      VALIDATION_MESSAGES.AUTH_EMAIL_REQUIRED,
    'users.register.error.password_required':
      VALIDATION_MESSAGES.AUTH_PASSWORD_REQUIRED,
    'users.register.error.password_too_short':
      VALIDATION_MESSAGES.AUTH_PASSWORD_TOO_SHORT,
    'users.register.error.email_exists': VALIDATION_MESSAGES.EMAIL_EXISTS,
    'users.login.error.invalid_credentials':
      VALIDATION_MESSAGES.INVALID_CREDENTIALS,
  };

  return (
    errorMap[serverError] ||
    (serverError.includes('login')
      ? VALIDATION_MESSAGES.LOGIN_FAILED
      : VALIDATION_MESSAGES.REGISTER_FAILED)
  );
};

interface LoginModalProps {
  isOpen: boolean;
  onClose: () => void;
  initialMode?: FormMode;
}

const LoginModal: React.FC<LoginModalProps> = ({
  isOpen,
  onClose,
  initialMode = 'login',
}) => {
  const t = useTranslations('auth');
  const { updateUser } = useAuth();

  // Form state management
  const {
    mode,
    formData,
    errors,
    isLoading,
    successMessage,
    setField,
    setMode,
    setErrors,
    setError,
    setLoading,
    setSuccess,
    clearSuccess,
    resetForm,
  } = useAuthForm(initialMode);

  // Form validation
  const { canSubmit, validateField } = useFormValidation(formData, mode);

  // Debounced values for real-time validation
  const debouncedEmail = useDebounce(formData.email, 500);
  const debouncedPassword = useDebounce(formData.password, 500);
  const debouncedName = useDebounce(formData.name, 500);
  const debouncedConfirmPassword = useDebounce(formData.confirmPassword, 500);
  const debouncedPhone = useDebounce(formData.phone, 500);

  // Memoized validation handler to optimize dependencies
  const handleFieldValidation = useCallback(
    (field: keyof FormData, value: string, currentError?: string) => {
      if (value) {
        const validationError = validateField(field, value);
        // Only update if error has changed
        if (validationError !== currentError && validationError) {
          setError(field, validationError);
        }
      }
    },
    [validateField, setError]
  );

  // Memoized fields configuration
  const fieldsToValidate = useMemo(
    () => [
      {
        field: 'email' as keyof FormData,
        value: debouncedEmail,
        originalValue: formData.email,
        error: errors.email,
      },
      {
        field: 'password' as keyof FormData,
        value: debouncedPassword,
        originalValue: formData.password,
        error: errors.password,
      },
      ...(mode === 'register'
        ? [
            {
              field: 'name' as keyof FormData,
              value: debouncedName,
              originalValue: formData.name,
              error: errors.name,
            },
            {
              field: 'confirmPassword' as keyof FormData,
              value: debouncedConfirmPassword,
              originalValue: formData.confirmPassword,
              error: errors.confirmPassword,
            },
            {
              field: 'phone' as keyof FormData,
              value: debouncedPhone,
              originalValue: formData.phone,
              error: errors.phone,
            },
          ]
        : []),
    ],
    [
      debouncedEmail,
      debouncedPassword,
      debouncedName,
      debouncedConfirmPassword,
      debouncedPhone,
      formData.email,
      formData.password,
      formData.name,
      formData.confirmPassword,
      formData.phone,
      errors.email,
      errors.password,
      errors.name,
      errors.confirmPassword,
      errors.phone,
      mode,
    ]
  );

  // Optimized real-time validation effect
  useEffect(() => {
    fieldsToValidate.forEach(
      ({ field, value, originalValue, error: currentError }) => {
        // Skip if debounced value doesn't match current form value (user is still typing)
        if (value && value !== originalValue) return;

        handleFieldValidation(field, value, currentError);
      }
    );
  }, [fieldsToValidate, handleFieldValidation]);

  // Handle form submission
  const handleSubmit = useCallback(
    async (e: React.FormEvent) => {
      e.preventDefault();
      if (!canSubmit || isLoading) return;

      setLoading(true);
      setErrors({});

      try {
        if (mode === 'login') {
          const loginResponse = await AuthService.login({
            email: formData.email,
            password: formData.password,
          });

          // Update user state with login response
          if (
            loginResponse.user &&
            loginResponse.user.id &&
            loginResponse.user.name &&
            loginResponse.user.email
          ) {
            updateUser({
              id: loginResponse.user.id,
              name: loginResponse.user.name,
              email: loginResponse.user.email,
              provider: 'email', // Email login always uses 'email' provider
              picture_url: loginResponse.user.picture_url,
            });
          }

          // Small delay to ensure cookies are set
          setTimeout(() => {
            // Login successful, close modal
            onClose();
          }, 100);
        } else {
          const registerResponse = await AuthService.register({
            name: formData.name,
            email: formData.email,
            password: formData.password,
            phone: formData.phone || undefined,
          });

          // Update user state with register response
          if (
            registerResponse.user &&
            registerResponse.user.id &&
            registerResponse.user.name &&
            registerResponse.user.email
          ) {
            updateUser({
              id: registerResponse.user.id,
              name: registerResponse.user.name,
              email: registerResponse.user.email,
              provider: 'email', // Email registration always uses 'email' provider
              picture_url: registerResponse.user.picture_url,
            });
          }

          setSuccess(t('registerForm.successMessage'));

          // Close modal after successful registration
          setTimeout(() => {
            onClose();
          }, 1500);
        }
      } catch (error) {
        if (error instanceof Error) {
          // Map server error messages to appropriate client error keys
          const errorMessage = error.message;

          if (errorMessage === 'users.errors.tooManyAttempts') {
            setError('general', VALIDATION_MESSAGES.TOO_MANY_ATTEMPTS);
          } else if (errorMessage === 'users.errors.accountTemporarilyLocked') {
            setError('general', VALIDATION_MESSAGES.ACCOUNT_TEMPORARILY_LOCKED);
          } else if (errorMessage === 'users.register.error.tooManyAttempts') {
            setError('general', VALIDATION_MESSAGES.REGISTER_TOO_MANY_ATTEMPTS);
          } else if (errorMessage === 'users.login.error.invalid_credentials') {
            setError('general', VALIDATION_MESSAGES.INVALID_CREDENTIALS);
          } else if (errorMessage === 'users.register.error.email_exists') {
            setError('general', VALIDATION_MESSAGES.EMAIL_EXISTS);
          } else if (
            errorMessage.startsWith('users.register.error.') ||
            errorMessage.startsWith('users.login.error.')
          ) {
            // Handle other specific registration/login errors by mapping them
            const mappedError = mapServerErrorToClientError(errorMessage);
            setError('general', mappedError);
          } else {
            // Default fallback errors
            setError(
              'general',
              mode === 'login'
                ? VALIDATION_MESSAGES.LOGIN_FAILED
                : VALIDATION_MESSAGES.REGISTER_FAILED
            );
          }
        }
      } finally {
        setLoading(false);
      }
    },
    [
      canSubmit,
      isLoading,
      mode,
      formData,
      setLoading,
      setErrors,
      setError,
      setSuccess,
      onClose,
      t,
      updateUser,
    ]
  );

  // Handle Google OAuth
  const handleGoogleLogin = useCallback(async () => {
    if (isLoading) return;

    try {
      await AuthService.loginWithGoogle();
    } catch {
      setError('general', VALIDATION_MESSAGES.LOGIN_FAILED);
    }
  }, [isLoading, setError]);

  // Handle field changes
  const handleFieldChange = useCallback(
    (field: keyof FormData, value: string) => {
      setField(field, value);
    },
    [setField]
  );

  // Handle mode switching
  const handleSwitchMode = useCallback(
    (newMode: FormMode) => {
      setMode(newMode);
      clearSuccess();
    },
    [setMode, clearSuccess]
  );

  // Reset form when modal closes
  useEffect(() => {
    if (!isOpen) {
      resetForm();
      clearSuccess();
    }
  }, [isOpen, resetForm, clearSuccess]);

  // Handle modal close with escape key
  useEffect(() => {
    const handleEscape = (e: KeyboardEvent) => {
      if (e.key === 'Escape' && isOpen) {
        onClose();
      }
    };

    document.addEventListener('keydown', handleEscape);
    return () => document.removeEventListener('keydown', handleEscape);
  }, [isOpen, onClose]);

  // Don't render if not open
  if (!isOpen) return null;

  return (
    <div
      className="modal modal-open"
      onClick={onClose}
      role="dialog"
      aria-modal="true"
      aria-labelledby="modal-title"
    >
      <div
        className="modal-box w-11/12 max-w-md"
        onClick={(e) => e.stopPropagation()}
      >
        {/* Close button */}
        <button
          onClick={onClose}
          className="btn btn-sm btn-circle btn-ghost absolute right-2 top-2"
          aria-label={t('common.close')}
          disabled={isLoading}
        >
          ✕
        </button>

        {/* Modal title */}
        <h3 id="modal-title" className="font-bold text-lg mb-6">
          {mode === 'login' ? t('loginForm.title') : t('registerForm.title')}
        </h3>

        {/* Form content */}
        {mode === 'login' ? (
          <LoginForm
            formData={formData}
            errors={errors}
            isLoading={isLoading}
            onFieldChange={handleFieldChange}
            onSubmit={handleSubmit}
            onGoogleLogin={handleGoogleLogin}
            onSwitchToRegister={() => handleSwitchMode('register')}
            canSubmit={canSubmit}
          />
        ) : (
          <RegisterForm
            formData={formData}
            errors={errors}
            isLoading={isLoading}
            successMessage={successMessage}
            onFieldChange={handleFieldChange}
            onSubmit={handleSubmit}
            onGoogleLogin={handleGoogleLogin}
            onSwitchToLogin={() => handleSwitchMode('login')}
            onClose={onClose}
            canSubmit={canSubmit}
          />
        )}
      </div>
    </div>
  );
};

export default LoginModal;
