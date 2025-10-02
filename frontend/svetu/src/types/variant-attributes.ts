// Общие типы для variant attributes
export interface VariantAttributeBase {
  id?: number;
  name: string;
  display_name: string;
  is_required: boolean;
  sort_order: number;
  affects_stock: boolean;
  created_at?: string;
  updated_at?: string;
}

// Упрощённый тип для формы (старая структура)
export interface VariantAttributeSimple extends VariantAttributeBase {
  type: string;
}

// Полный тип для unified attributes (новая структура)
export interface VariantAttributeFull extends VariantAttributeBase {
  code: string;
  attribute_type: string;
  purpose: string;
  options?: any;
  validation_rules?: any;
  ui_settings?: any;
  is_searchable: boolean;
  is_filterable: boolean;
  is_variant_compatible: boolean;
  affects_price: boolean;
  is_active: boolean;
}

// Маппинг для категорий
export interface VariantMapping {
  id: number;
  variant_attribute_id: number;
  category_id: number;
  sort_order: number;
  is_required: boolean;
  attribute?: VariantAttributeFull;
  category?: {
    id: number;
    name: string;
  };
}
