// Конфигурация вариативных атрибутов для категорий
export interface VariantAttributeConfig {
  categorySlug: string;
  variantAttributes: string[]; // Имена атрибутов из product_variant_attributes
}

// Маппинг категорий на доступные вариативные атрибуты
export const categoryVariantAttributes: VariantAttributeConfig[] = [
  // Одежда
  {
    categorySlug: 'womens-clothing',
    variantAttributes: ['color', 'size', 'material', 'pattern', 'style'],
  },
  {
    categorySlug: 'mens-clothing',
    variantAttributes: ['color', 'size', 'material', 'pattern', 'style'],
  },
  {
    categorySlug: 'kids-clothing',
    variantAttributes: ['color', 'size', 'material', 'pattern'],
  },
  {
    categorySlug: 'sports-clothing',
    variantAttributes: ['color', 'size', 'material'],
  },

  // Обувь
  {
    categorySlug: 'shoes',
    variantAttributes: ['color', 'size', 'material', 'style'],
  },

  // Аксессуары
  {
    categorySlug: 'bags',
    variantAttributes: ['color', 'size', 'material', 'style', 'pattern'],
  },
  {
    categorySlug: 'accessories',
    variantAttributes: ['color', 'size', 'material', 'style', 'pattern'],
  },

  // Электроника
  {
    categorySlug: 'smartphones',
    variantAttributes: ['color', 'memory', 'storage'],
  },
  {
    categorySlug: 'computers',
    variantAttributes: ['color', 'memory', 'storage', 'connectivity'],
  },
  {
    categorySlug: 'gaming-consoles',
    variantAttributes: ['color', 'storage', 'bundle'],
  },
  {
    categorySlug: 'electronics-accessories',
    variantAttributes: ['color', 'connectivity', 'bundle'],
  },

  // Бытовая техника
  {
    categorySlug: 'home-appliances',
    variantAttributes: ['color', 'capacity', 'power'],
  },

  // Мебель
  {
    categorySlug: 'furniture',
    variantAttributes: ['color', 'material', 'style'],
  },

  // Кухонная утварь
  {
    categorySlug: 'kitchenware',
    variantAttributes: ['color', 'capacity', 'material'],
  },
];

// Функция для получения вариативных атрибутов категории
export function getVariantAttributesForCategory(
  categorySlug: string
): string[] {
  const config = categoryVariantAttributes.find(
    (c) => c.categorySlug === categorySlug
  );
  return config?.variantAttributes || [];
}

// Функция проверки, является ли атрибут вариативным для категории
export function isVariantAttribute(
  categorySlug: string,
  attributeName: string
): boolean {
  const variantAttrs = getVariantAttributesForCategory(categorySlug);
  return variantAttrs.includes(attributeName.toLowerCase());
}
