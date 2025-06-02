import React from 'react';
import { useTranslations } from 'next-intl';
import FormField from '@/components/FormField';
import GoogleIcon from '@/components/GoogleIcon';
import { FormData, FormErrors } from '@/hooks/useAuthForm';

interface LoginFormProps {
  formData: FormData;
  errors: FormErrors;
  isLoading: boolean;
  onFieldChange: (field: keyof FormData, value: string) => void;
  onSubmit: (e: React.FormEvent) => void;
  onGoogleLogin: () => void;
  onSwitchToRegister: () => void;
  canSubmit: boolean;
}

const LoginForm: React.FC<LoginFormProps> = ({
  formData,
  errors,
  isLoading,
  onFieldChange,
  onSubmit,
  onGoogleLogin,
  onSwitchToRegister,
  canSubmit,
}) => {
  const t = useTranslations('auth');

  return (
    <form onSubmit={onSubmit} className="space-y-6">
      {/* Email Field */}
      <FormField
        label={t('loginForm.email')}
        required
        error={errors.email ? t(errors.email) : undefined}
      >
        <input
          type="email"
          className={`input input-bordered w-full ${
            errors.email ? 'input-error' : ''
          }`}
          placeholder={t('loginForm.emailPlaceholder')}
          value={formData.email}
          onChange={(e) => onFieldChange('email', e.target.value)}
          disabled={isLoading}
          autoComplete="email"
          required
        />
      </FormField>

      {/* Password Field */}
      <FormField
        label={t('loginForm.password')}
        required
        error={errors.password ? t(errors.password) : undefined}
      >
        <input
          type="password"
          className={`input input-bordered w-full ${
            errors.password ? 'input-error' : ''
          }`}
          placeholder={t('loginForm.passwordPlaceholder')}
          value={formData.password}
          onChange={(e) => onFieldChange('password', e.target.value)}
          disabled={isLoading}
          autoComplete="current-password"
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

      {/* Login Button */}
      <button
        type="submit"
        className={`btn btn-primary w-full ${isLoading ? 'loading' : ''}`}
        disabled={isLoading || !canSubmit}
      >
        {isLoading ? t('loginForm.loggingIn') : t('loginForm.submit')}
      </button>

      {/* Divider */}
      <div className="divider">{t('loginForm.or')}</div>

      {/* Google Login Button */}
      <button
        type="button"
        onClick={onGoogleLogin}
        className="btn btn-outline w-full"
        disabled={isLoading}
      >
        <GoogleIcon />
        {t('loginForm.googleLogin')}
      </button>

      {/* Switch to Register */}
      <div className="text-center">
        <span className="text-base-content/70">
          {t('loginForm.registerText')}{' '}
        </span>
        <button
          type="button"
          onClick={onSwitchToRegister}
          className="link link-primary"
          disabled={isLoading}
        >
          {t('loginForm.register')}
        </button>
      </div>
    </form>
  );
};

export default React.memo(LoginForm);
