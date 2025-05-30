/**
 * Validation utilities for forms
 */

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
      message: 'validation.nameRequired',
    });
  } else if (data.name.trim().length < 2) {
    errors.push({
      field: 'name',
      message: 'validation.nameMinLength',
    });
  } else if (data.name.trim().length > 100) {
    errors.push({
      field: 'name',
      message: 'validation.nameMaxLength',
    });
  }

  // Phone validation (optional but must be valid if provided)
  if (data.phone && data.phone.trim()) {
    if (!isValidPhone(data.phone.trim())) {
      errors.push({
        field: 'phone',
        message: 'validation.phoneInvalid',
      });
    }
  }

  // City validation (optional but reasonable length if provided)
  if (data.city && data.city.trim().length > 100) {
    errors.push({
      field: 'city',
      message: 'validation.cityMaxLength',
    });
  }

  // Country validation (optional but reasonable length if provided)
  if (data.country && data.country.trim().length > 100) {
    errors.push({
      field: 'country',
      message: 'validation.countryMaxLength',
    });
  }

  return errors;
}

/**
 * Validate phone number with improved regex
 */
function isValidPhone(phone: string): boolean {
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
