export interface TranslatedAttribute {
  id: number;
  name: string;
  display_name: string;
  attribute_type: string;
  options?: any;
  translations?: Record<string, string>;
  option_translations?: Record<string, Record<string, string>>;
}

export function getTranslatedAttribute(
  attribute: TranslatedAttribute | null | undefined,
  locale: string
) {
  if (!attribute) {
    return {
      displayName: '',
      getOptionLabel: (option: string) => option,
    };
  }

  // Получаем переведенное название атрибута
  // Для английского языка используем display_name если нет перевода
  const displayName =
    locale === 'en' && !attribute.translations?.[locale]
      ? attribute.display_name || attribute.name
      : attribute.translations?.[locale] ||
        attribute.display_name ||
        attribute.name;

  // Функция для получения переведенной опции
  const getOptionLabel = (option: string): string => {
    return attribute.option_translations?.[locale]?.[option] || option;
  };

  return {
    displayName,
    getOptionLabel,
  };
}
