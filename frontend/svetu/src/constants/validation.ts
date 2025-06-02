/**
 * Validation constants and patterns used across the application
 */

// Regex patterns for validation
export const VALIDATION_PATTERNS = {
  EMAIL: /^[^\s@]+@[^\s@]+\.[^\s@]+$/,
  PHONE: /^\+?[\d\s\-\(\)]{10,}$/,
} as const;

// Validation constraints
export const VALIDATION_CONSTRAINTS = {
  PASSWORD_MIN_LENGTH: 6,
  NAME_MIN_LENGTH: 2,
  NAME_MAX_LENGTH: 100,
  CITY_MAX_LENGTH: 100,
  COUNTRY_MAX_LENGTH: 100,
} as const;

// Translation keys for validation messages
export const VALIDATION_MESSAGES = {
  // Email validation
  EMAIL_REQUIRED: 'validation.emailRequired',
  EMAIL_INVALID: 'validation.emailInvalid',

  // Password validation
  PASSWORD_REQUIRED: 'validation.passwordRequired',
  PASSWORD_TOO_SHORT: 'validation.passwordTooShort',
  PASSWORDS_DO_NOT_MATCH: 'validation.passwordsDoNotMatch',

  // Name validation
  NAME_REQUIRED: 'validation.nameRequired',
  NAME_TOO_SHORT: 'validation.nameMinLength',
  NAME_TOO_LONG: 'validation.nameMaxLength',

  // Phone validation
  PHONE_INVALID: 'validation.phoneInvalid',

  // Profile validation (without profile prefix since useTranslations('profile') adds it)
  PROFILE_NAME_REQUIRED: 'validation.nameRequired',
  PROFILE_NAME_TOO_SHORT: 'validation.nameMinLength',
  PROFILE_NAME_TOO_LONG: 'validation.nameMaxLength',
  PROFILE_PHONE_INVALID: 'validation.phoneInvalid',
  CITY_TOO_LONG: 'validation.cityMaxLength',
  COUNTRY_TOO_LONG: 'validation.countryMaxLength',

  // Auth errors
  LOGIN_FAILED: 'errors.loginFailed',
  REGISTER_FAILED: 'errors.registerFailed',
  TOO_MANY_ATTEMPTS: 'errors.tooManyAttempts',
  ACCOUNT_TEMPORARILY_LOCKED: 'errors.accountTemporarilyLocked',
  INVALID_CREDENTIALS: 'errors.invalidCredentials',
  EMAIL_EXISTS: 'errors.emailExists',
  REGISTER_TOO_MANY_ATTEMPTS: 'register.error.tooManyAttempts',

  // Server error mappings
  AUTH_EMAIL_REQUIRED: 'errors.emailRequired',
  AUTH_PASSWORD_REQUIRED: 'errors.passwordRequired',
  AUTH_NAME_REQUIRED: 'errors.nameRequired',
  AUTH_PASSWORD_TOO_SHORT: 'errors.passwordTooShort',
} as const;

// Type for validation error keys
export type ValidationErrorKey =
  (typeof VALIDATION_MESSAGES)[keyof typeof VALIDATION_MESSAGES];
