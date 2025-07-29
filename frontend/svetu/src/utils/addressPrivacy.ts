/**
 * Утилиты для обеспечения приватности адреса
 */

export interface AddressPrivacyOptions {
  showHouseNumber?: boolean;
  showStreet?: boolean;
  showCity?: boolean;
  showRegion?: boolean;
  showCountry?: boolean;
}

const DEFAULT_PRIVACY_OPTIONS: AddressPrivacyOptions = {
  showHouseNumber: false, // По умолчанию скрываем номер дома
  showStreet: true,
  showCity: true,
  showRegion: true,
  showCountry: true,
};

/**
 * Форматирует адрес с учетом настроек приватности
 */
export function formatAddressWithPrivacy(
  fullAddress: string,
  options: AddressPrivacyOptions = DEFAULT_PRIVACY_OPTIONS
): string {
  if (!fullAddress) return '';

  // Пытаемся разобрать адрес на компоненты
  // Типичный формат: "Улица Номер, Город Индекс, Регион, Страна"
  const parts = fullAddress.split(',').map(part => part.trim());
  
  if (parts.length === 0) return fullAddress;

  let result: string[] = [];

  // Первая часть обычно содержит улицу и номер дома
  if (parts[0] && options.showStreet) {
    if (!options.showHouseNumber) {
      // Удаляем номер дома (обычно в конце строки)
      const streetWithoutNumber = parts[0].replace(/\s+\d+[а-яa-z]?$/i, '').trim();
      if (streetWithoutNumber) {
        result.push(streetWithoutNumber);
      }
    } else {
      result.push(parts[0]);
    }
  }

  // Остальные части адреса
  for (let i = 1; i < parts.length; i++) {
    const part = parts[i];
    
    // Пропускаем почтовый индекс
    if (/^\d{5,6}$/.test(part)) {
      continue;
    }
    
    // Добавляем остальные части если нужно
    if (i === 1 && options.showCity) {
      // Убираем индекс из города если есть
      const cityWithoutIndex = part.replace(/\s*\d{5,6}\s*/, ' ').trim();
      if (cityWithoutIndex) {
        result.push(cityWithoutIndex);
      }
    } else if (i === 2 && options.showRegion) {
      result.push(part);
    } else if (i === 3 && options.showCountry) {
      result.push(part);
    }
  }

  return result.join(', ');
}

/**
 * Извлекает только улицу из полного адреса
 */
export function extractStreetFromAddress(fullAddress: string): string {
  if (!fullAddress) return '';
  
  const firstComma = fullAddress.indexOf(',');
  const streetPart = firstComma > 0 ? fullAddress.substring(0, firstComma) : fullAddress;
  
  // Удаляем номер дома
  return streetPart.replace(/\s+\d+[а-яa-z]?$/i, '').trim();
}

/**
 * Проверяет, содержит ли адрес номер дома
 */
export function hasHouseNumber(address: string): boolean {
  return /\s+\d+[а-яa-z]?(?:\s|,|$)/i.test(address);
}