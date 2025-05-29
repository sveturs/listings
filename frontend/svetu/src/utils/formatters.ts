// Utility functions for formatting data

export const formatPrice = (price?: number, currency: string = 'RSD'): string => {
  if (price === undefined || price === null) return '';
  
  return new Intl.NumberFormat('sr-RS', {
    style: 'currency',
    currency: currency,
    maximumFractionDigits: 0
  }).format(price);
};

export const formatDate = (dateString?: string, locale: string = 'sr-RS'): string => {
  if (!dateString) return '';
  
  try {
    const date = new Date(dateString);
    return date.toLocaleDateString(locale);
  } catch {
    return dateString;
  }
};

export const formatDateTime = (dateString?: string, locale: string = 'sr-RS'): string => {
  if (!dateString) return '';
  
  try {
    const date = new Date(dateString);
    return date.toLocaleString(locale);
  } catch {
    return dateString;
  }
};

export const formatNumber = (num?: number, locale: string = 'sr-RS'): string => {
  if (num === undefined || num === null) return '0';
  
  return new Intl.NumberFormat(locale).format(num);
};

export const formatPercentage = (value?: number): string => {
  if (value === undefined || value === null) return '0%';
  
  return `${Math.round(value)}%`;
};

export const formatDistance = (meters?: number): string => {
  if (!meters) return '';
  
  if (meters < 1000) {
    return `${Math.round(meters)}m`;
  }
  
  return `${(meters / 1000).toFixed(1)}km`;
};

export const truncateText = (text: string, maxLength: number): string => {
  if (text.length <= maxLength) return text;
  
  return text.substring(0, maxLength) + '...';
};