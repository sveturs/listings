export interface DeliveryAttributes {
  weight_kg: number;
  dimensions: {
    length_cm: number;
    width_cm: number;
    height_cm: number;
  };
  volume_m3?: number;
  is_fragile: boolean;
  requires_special_handling: boolean;
  stackable: boolean;
  max_stack_weight_kg?: number;
  packaging_type: 'box' | 'envelope' | 'pallet' | 'custom';
  hazmat_class?: string | null;
}

export interface DeliveryProvider {
  id: number;
  code: string;
  name: string;
  logo_url?: string;
  is_active: boolean;
  supports_cod: boolean;
  supports_insurance: boolean;
  supports_tracking: boolean;
  capabilities: Record<string, any>;
}

export interface DeliveryQuote {
  provider_id: number;
  provider_name: string;
  provider_code: string;
  total_price: number;
  delivery_cost: number;
  insurance_cost?: number;
  cod_fee?: number;
  estimated_days: number;
  cost_breakdown: {
    base_price: number;
    weight_surcharge?: number;
    fragile_surcharge?: number;
    oversized_surcharge?: number;
    special_handling_surcharge?: number;
    [key: string]: number | undefined;
  };
  services: {
    name: string;
    price: number;
    included: boolean;
  }[];
  pickup_date?: string;
  estimated_delivery_date?: string;
}

export interface CalculationRequest {
  from_location: {
    city: string;
    postal_code?: string;
    address?: string;
  };
  to_location: {
    city: string;
    postal_code?: string;
    address?: string;
  };
  items: {
    product_id: number;
    product_type: 'listing' | 'storefront_product';
    quantity: number;
    attributes?: DeliveryAttributes;
  }[];
  provider_id?: number;
  insurance_value?: number;
  cod_amount?: number;
}

export interface CalculationResponse {
  success: boolean;
  data?: {
    providers: DeliveryQuote[];
    cheapest?: DeliveryQuote;
    fastest?: DeliveryQuote;
    recommended?: DeliveryQuote;
  };
  message?: string;
}

export interface CategoryDefaults {
  category_id: number;
  default_weight_kg?: number;
  default_length_cm?: number;
  default_width_cm?: number;
  default_height_cm?: number;
  default_packaging_type?: string;
  is_typically_fragile: boolean;
}

export interface ValidationErrors {
  weight_kg?: string;
  dimensions?: {
    length_cm?: string;
    width_cm?: string;
    height_cm?: string;
  };
  packaging_type?: string;
  [key: string]: any;
}
