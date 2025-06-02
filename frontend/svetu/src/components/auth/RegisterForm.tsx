import React from 'react';
import { useTranslations } from 'next-intl';
import FormField from '@/components/FormField';
import GoogleIcon from '@/components/GoogleIcon';
import { FormData, FormErrors } from '@/hooks/useAuthForm';

interface RegisterFormProps {
  formData: FormData;
  errors: FormErrors;
  isLoading: boolean;
  successMessage: string;
  onFieldChange: (field: keyof FormData, value: string) => void;
  onSubmit: (e: React.FormEvent) => void;
  onGoogleLogin: () => void;
  onSwitchToLogin: () => void;
  onClose: () => void;
  canSubmit: boolean;
}

const RegisterForm: React.FC<RegisterFormProps> = ({
  formData,
  errors,
  isLoading,
  successMessage,
  onFieldChange,
  onSubmit,
  onGoogleLogin,
  onSwitchToLogin,
  onClose,
  canSubmit,
}) => {
  const t = useTranslations('auth');

  // Show success message if registration completed
  if (successMessage) {
    return (
      <div className="text-center space-y-6">
        <div className="alert alert-success">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            className="stroke-current shrink-0 h-6 w-6"
            fill="none"
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth="2"
              d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
            />
          </svg>
          <span>{t('registerForm.successMessage')}</span>
        </div>

        <p className="text-base-content/70">
          {t('registerForm.successDescription')}
        </p>

        <div className="flex gap-3 justify-center">
          <button onClick={onSwitchToLogin} className="btn btn-primary">
            {t('registerForm.switchToLogin')}
          </button>
          <button onClick={onClose} className="btn btn-ghost">
            {t('registerForm.close')}
          </button>
        </div>
      </div>
    );
  }

  return (
    <form onSubmit={onSubmit} className="space-y-6">
      {/* Name Field */}
      <FormField
        label={t('registerForm.name')}
        required
        error={errors.name ? t(errors.name) : undefined}
      >
        <input
          type="text"
          className={`input input-bordered w-full ${
            errors.name ? 'input-error' : ''
          }`}
          placeholder={t('registerForm.namePlaceholder')}
          value={formData.name}
          onChange={(e) => onFieldChange('name', e.target.value)}
          disabled={isLoading}
          autoComplete="given-name"
          required
          minLength={2}
        />
      </FormField>

      {/* Email Field */}
      <FormField
        label={t('registerForm.email')}
        required
        error={errors.email ? t(errors.email) : undefined}
      >
        <input
          type="email"
          className={`input input-bordered w-full ${
            errors.email ? 'input-error' : ''
          }`}
          placeholder={t('registerForm.emailPlaceholder')}
          value={formData.email}
          onChange={(e) => onFieldChange('email', e.target.value)}
          disabled={isLoading}
          autoComplete="email"
          required
        />
      </FormField>

      {/* Phone Field (Optional) */}
      <FormField
        label={t('registerForm.phone')}
        error={errors.phone ? t(errors.phone) : undefined}
      >
        <input
          type="tel"
          className={`input input-bordered w-full ${
            errors.phone ? 'input-error' : ''
          }`}
          placeholder={t('registerForm.phonePlaceholder')}
          value={formData.phone}
          onChange={(e) => onFieldChange('phone', e.target.value)}
          disabled={isLoading}
          autoComplete="tel"
        />
      </FormField>

      {/* Password Field */}
      <FormField
        label={t('registerForm.password')}
        required
        error={errors.password ? t(errors.password) : undefined}
      >
        <input
          type="password"
          className={`input input-bordered w-full ${
            errors.password ? 'input-error' : ''
          }`}
          placeholder={t('registerForm.passwordPlaceholder')}
          value={formData.password}
          onChange={(e) => onFieldChange('password', e.target.value)}
          disabled={isLoading}
          autoComplete="new-password"
          required
          minLength={6}
        />
      </FormField>

      {/* Confirm Password Field */}
      <FormField
        label={t('registerForm.confirmPassword')}
        required
        error={errors.confirmPassword ? t(errors.confirmPassword) : undefined}
      >
        <input
          type="password"
          className={`input input-bordered w-full ${
            errors.confirmPassword ? 'input-error' : ''
          }`}
          placeholder={t('registerForm.confirmPasswordPlaceholder')}
          value={formData.confirmPassword}
          onChange={(e) => onFieldChange('confirmPassword', e.target.value)}
          disabled={isLoading}
          autoComplete="new-password"
          required
          minLength={6}
        />
      </FormField>

      {/* General Error */}
      {errors.general && (
        <div className="alert alert-error">
          <span>{t(errors.general)}</span>
        </div>
      )}

      {/* Register Button */}
      <button
        type="submit"
        className={`btn btn-primary w-full ${isLoading ? 'loading' : ''}`}
        disabled={isLoading || !canSubmit}
      >
        {isLoading ? t('registerForm.submitting') : t('registerForm.submit')}
      </button>

      {/* Divider */}
      <div className="divider">{t('registerForm.or')}</div>

      {/* Google Register Button */}
      <button
        type="button"
        onClick={onGoogleLogin}
        className="btn btn-outline w-full"
        disabled={isLoading}
      >
        <GoogleIcon />
        {t('registerForm.googleRegister')}
      </button>

      {/* Switch to Login */}
      <div className="text-center">
        <span className="text-base-content/70">
          {t('registerForm.loginText')}{' '}
        </span>
        <button
          type="button"
          onClick={onSwitchToLogin}
          className="link link-primary"
          disabled={isLoading}
        >
          {t('registerForm.switchToLogin')}
        </button>
      </div>
    </form>
  );
};

export default React.memo(RegisterForm);
