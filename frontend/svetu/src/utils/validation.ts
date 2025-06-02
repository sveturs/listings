/**
 * Validation utilities for forms
 */

import {
  VALIDATION_PATTERNS,
  VALIDATION_CONSTRAINTS,
  VALIDATION_MESSAGES,
} from '@/constants/validation';

export interface ValidationError {
  field: string;
  message: string;
}

/**
 * Validate user profile data
 */
export function validateProfileData(data: {
  name: string;
  phone?: string;
  city?: string;
  country?: string;
}): ValidationError[] {
  const errors: ValidationError[] = [];

  // Name validation
  if (!data.name || !data.name.trim()) {
    errors.push({
      field: 'name',
      message: VALIDATION_MESSAGES.PROFILE_NAME_REQUIRED,
    });
  } else if (data.name.trim().length < VALIDATION_CONSTRAINTS.NAME_MIN_LENGTH) {
    errors.push({
      field: 'name',
      message: VALIDATION_MESSAGES.PROFILE_NAME_TOO_SHORT,
    });
  } else if (data.name.trim().length > VALIDATION_CONSTRAINTS.NAME_MAX_LENGTH) {
    errors.push({
      field: 'name',
      message: VALIDATION_MESSAGES.PROFILE_NAME_TOO_LONG,
    });
  }

  // Phone validation (optional but must be valid if provided)
  if (data.phone && data.phone.trim()) {
    if (!isValidPhone(data.phone.trim())) {
      errors.push({
        field: 'phone',
        message: VALIDATION_MESSAGES.PROFILE_PHONE_INVALID,
      });
    }
  }

  // City validation (optional but reasonable length if provided)
  if (
    data.city &&
    data.city.trim().length > VALIDATION_CONSTRAINTS.CITY_MAX_LENGTH
  ) {
    errors.push({
      field: 'city',
      message: VALIDATION_MESSAGES.CITY_TOO_LONG,
    });
  }

  // Country validation (optional but reasonable length if provided)
  if (
    data.country &&
    data.country.trim().length > VALIDATION_CONSTRAINTS.COUNTRY_MAX_LENGTH
  ) {
    errors.push({
      field: 'country',
      message: VALIDATION_MESSAGES.COUNTRY_TOO_LONG,
    });
  }

  return errors;
}

/**
 * Validate phone number using centralized pattern with enhanced checks
 */
function isValidPhone(phone: string): boolean {
  // First check with the standard pattern
  if (!VALIDATION_PATTERNS.PHONE.test(phone)) {
    return false;
  }

  // Remove all non-digit characters except +
  const cleanPhone = phone.replace(/[^\d+]/g, '');

  // More strict international phone validation
  // Must start with + or digit, have 7-15 digits total
  const phoneRegex = /^(\+?[1-9]\d{6,14})$/;

  // Additional checks
  if (!phoneRegex.test(cleanPhone)) {
    return false;
  }

  // Check for obvious invalid patterns
  const invalidPatterns = [
    /^(\+?0+)$/, // All zeros
    /^(\+?1{7,})$/, // Too many ones
    /^(\+?9{7,})$/, // Too many nines
  ];

  return !invalidPatterns.some((pattern) => pattern.test(cleanPhone));
}

/**
 * Check if form data has any changes compared to original
 */
export function hasFormChanges(
  original: Record<string, string>,
  current: Record<string, string>,
  fields: string[]
): boolean {
  return fields.some((field) => {
    const originalValue = (original[field] || '').trim();
    const currentValue = (current[field] || '').trim();
    return originalValue !== currentValue;
  });
}
