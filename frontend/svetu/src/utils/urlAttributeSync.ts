import { ReadonlyURLSearchParams } from 'next/navigation';

export interface AttributeURLParams {
  [key: string]: string | string[] | undefined;
}

export interface ParsedAttributes {
  attributes: Record<string, string[]>;
  other: Record<string, string>;
}

/**
 * Parse attribute filters from URL search params
 * @param searchParams - URL search params
 * @returns Parsed attributes and other parameters
 */
export function parseAttributesFromURL(
  searchParams: ReadonlyURLSearchParams | URLSearchParams
): ParsedAttributes {
  const attributes: Record<string, string[]> = {};
  const other: Record<string, string> = {};

  searchParams.forEach((value, key) => {
    if (key.startsWith('attr_')) {
      // Attribute filter
      const attrId = key.replace('attr_', '');
      attributes[attrId] = value.split(',').filter((v) => v.length > 0);
    } else {
      // Other parameter
      other[key] = value;
    }
  });

  return { attributes, other };
}

/**
 * Build URL search params from filters
 * @param attributes - Attribute filters
 * @param otherParams - Other search parameters
 * @returns URL search params string
 */
export function buildAttributeURLParams(
  attributes: Record<string, string[]>,
  otherParams: Record<string, string | number | boolean | undefined>
): string {
  const params = new URLSearchParams();

  // Add other parameters
  Object.entries(otherParams).forEach(([key, value]) => {
    if (value !== undefined && value !== null && value !== '') {
      params.set(key, String(value));
    }
  });

  // Add attribute filters
  Object.entries(attributes).forEach(([attrId, values]) => {
    if (values.length > 0) {
      params.set(`attr_${attrId}`, values.join(','));
    }
  });

  return params.toString();
}

/**
 * Merge current URL params with new attribute filters
 * @param currentParams - Current URL search params
 * @param newAttributes - New attribute filters to apply
 * @param clearOthers - Clear other attributes not in newAttributes
 * @returns Merged URL params string
 */
export function mergeAttributeParams(
  currentParams: ReadonlyURLSearchParams | URLSearchParams,
  newAttributes: Record<string, string[]>,
  clearOthers: boolean = false
): string {
  const { attributes: currentAttributes, other } =
    parseAttributesFromURL(currentParams);

  let mergedAttributes: Record<string, string[]>;

  if (clearOthers) {
    // Replace all attributes with new ones
    mergedAttributes = newAttributes;
  } else {
    // Merge attributes
    mergedAttributes = { ...currentAttributes };
    Object.entries(newAttributes).forEach(([key, values]) => {
      if (values.length === 0) {
        delete mergedAttributes[key];
      } else {
        mergedAttributes[key] = values;
      }
    });
  }

  return buildAttributeURLParams(mergedAttributes, other);
}

/**
 * Get attribute value from URL
 * @param searchParams - URL search params
 * @param attributeId - Attribute ID
 * @returns Attribute values or empty array
 */
export function getAttributeFromURL(
  searchParams: ReadonlyURLSearchParams | URLSearchParams,
  attributeId: string
): string[] {
  const value = searchParams.get(`attr_${attributeId}`);
  return value ? value.split(',').filter((v) => v.length > 0) : [];
}

/**
 * Check if any attribute filters are active
 * @param searchParams - URL search params
 * @returns True if any attribute filters are present
 */
export function hasActiveAttributeFilters(
  searchParams: ReadonlyURLSearchParams | URLSearchParams
): boolean {
  for (const key of searchParams.keys()) {
    if (key.startsWith('attr_')) {
      return true;
    }
  }
  return false;
}

/**
 * Count active attribute filters
 * @param searchParams - URL search params
 * @returns Number of active attribute filters
 */
export function countActiveAttributeFilters(
  searchParams: ReadonlyURLSearchParams | URLSearchParams
): number {
  let count = 0;
  searchParams.forEach((value, key) => {
    if (key.startsWith('attr_')) {
      const values = value.split(',').filter((v) => v.length > 0);
      count += values.length;
    }
  });
  return count;
}

/**
 * Create a shareable URL with current filters
 * @param baseURL - Base URL
 * @param searchParams - Current search params
 * @returns Full shareable URL
 */
export function createShareableURL(
  baseURL: string,
  searchParams: ReadonlyURLSearchParams | URLSearchParams
): string {
  const params = searchParams.toString();
  return params ? `${baseURL}?${params}` : baseURL;
}

/**
 * Compress attribute filters for shorter URLs
 * @param attributes - Attribute filters
 * @returns Compressed string representation
 */
export function compressAttributes(
  attributes: Record<string, string[]>
): string {
  const compressed = Object.entries(attributes)
    .map(([id, values]) => `${id}:${values.join('|')}`)
    .join(';');
  return compressed ? btoa(compressed) : '';
}

/**
 * Decompress attribute filters from compressed string
 * @param compressed - Compressed attributes string
 * @returns Attribute filters
 */
export function decompressAttributes(
  compressed: string
): Record<string, string[]> {
  if (!compressed) return {};

  try {
    const decoded = atob(compressed);
    const attributes: Record<string, string[]> = {};

    decoded.split(';').forEach((part) => {
      const [id, valuesStr] = part.split(':');
      if (id && valuesStr) {
        attributes[id] = valuesStr.split('|');
      }
    });

    return attributes;
  } catch {
    return {};
  }
}

/**
 * Clean URL by removing empty parameters
 * @param searchParams - URL search params
 * @returns Cleaned URL params string
 */
export function cleanURLParams(
  searchParams: ReadonlyURLSearchParams | URLSearchParams
): string {
  const cleaned = new URLSearchParams();

  searchParams.forEach((value, key) => {
    if (value && value.trim() !== '') {
      cleaned.set(key, value);
    }
  });

  return cleaned.toString();
}
