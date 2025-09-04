// Утилиты для локализации цветов в зависимости от языка интерфейса

export interface LocalizedColor {
  hex: string;
  name: string;
  nameEn: string;
  nameRu: string;
  nameSr: string;
}

// Расширенная палитра цветов с переводами
export const LOCALIZED_COLORS: LocalizedColor[] = [
  {
    hex: '#000000',
    name: 'Black',
    nameEn: 'Black',
    nameRu: 'Черный',
    nameSr: 'Црна',
  },
  {
    hex: '#FFFFFF',
    name: 'White',
    nameEn: 'White',
    nameRu: 'Белый',
    nameSr: 'Бела',
  },
  {
    hex: '#FF0000',
    name: 'Red',
    nameEn: 'Red',
    nameRu: 'Красный',
    nameSr: 'Црвена',
  },
  {
    hex: '#00FF00',
    name: 'Green',
    nameEn: 'Green',
    nameRu: 'Зеленый',
    nameSr: 'Зелена',
  },
  {
    hex: '#0000FF',
    name: 'Blue',
    nameEn: 'Blue',
    nameRu: 'Синий',
    nameSr: 'Плава',
  },
  {
    hex: '#FFFF00',
    name: 'Yellow',
    nameEn: 'Yellow',
    nameRu: 'Желтый',
    nameSr: 'Жута',
  },
  {
    hex: '#FF00FF',
    name: 'Pink',
    nameEn: 'Pink',
    nameRu: 'Розовый',
    nameSr: 'Ружичаста',
  },
  {
    hex: '#00FFFF',
    name: 'Cyan',
    nameEn: 'Cyan',
    nameRu: 'Голубой',
    nameSr: 'Цијан',
  },
  {
    hex: '#800080',
    name: 'Purple',
    nameEn: 'Purple',
    nameRu: 'Фиолетовый',
    nameSr: 'Љубичаста',
  },
  {
    hex: '#FFA500',
    name: 'Orange',
    nameEn: 'Orange',
    nameRu: 'Оранжевый',
    nameSr: 'Наранџаста',
  },
  {
    hex: '#808080',
    name: 'Gray',
    nameEn: 'Gray',
    nameRu: 'Серый',
    nameSr: 'Сива',
  },
  {
    hex: '#8B4513',
    name: 'Brown',
    nameEn: 'Brown',
    nameRu: 'Коричневый',
    nameSr: 'Браон',
  },
  {
    hex: '#FFB6C1',
    name: 'Light Pink',
    nameEn: 'Light Pink',
    nameRu: 'Светло-розовый',
    nameSr: 'Светло ружичаста',
  },
  {
    hex: '#90EE90',
    name: 'Light Green',
    nameEn: 'Light Green',
    nameRu: 'Светло-зеленый',
    nameSr: 'Светло зелена',
  },
  {
    hex: '#87CEEB',
    name: 'Sky Blue',
    nameEn: 'Sky Blue',
    nameRu: 'Небесно-голубой',
    nameSr: 'Небеско плава',
  },
  {
    hex: '#DDA0DD',
    name: 'Plum',
    nameEn: 'Plum',
    nameRu: 'Светло-фиолетовый',
    nameSr: 'Светло љубичаста',
  },
  {
    hex: '#A52A2A',
    name: 'Brown Red',
    nameEn: 'Brown Red',
    nameRu: 'Коричнево-красный',
    nameSr: 'Браон црвена',
  },
  {
    hex: '#228B22',
    name: 'Forest Green',
    nameEn: 'Forest Green',
    nameRu: 'Лесной зеленый',
    nameSr: 'Шумско зелена',
  },
  {
    hex: '#4169E1',
    name: 'Royal Blue',
    nameEn: 'Royal Blue',
    nameRu: 'Королевский синий',
    nameSr: 'Краљевско плава',
  },
  {
    hex: '#FFD700',
    name: 'Gold',
    nameEn: 'Gold',
    nameRu: 'Золотой',
    nameSr: 'Златна',
  },
  {
    hex: '#C0C0C0',
    name: 'Silver',
    nameEn: 'Silver',
    nameRu: 'Серебряный',
    nameSr: 'Сребрна',
  },
  {
    hex: '#FF69B4',
    name: 'Hot Pink',
    nameEn: 'Hot Pink',
    nameRu: 'Ярко-розовый',
    nameSr: 'Јарко ружичаста',
  },
  {
    hex: '#32CD32',
    name: 'Lime Green',
    nameEn: 'Lime Green',
    nameRu: 'Лимонно-зеленый',
    nameSr: 'Лимун зелена',
  },
  {
    hex: '#FF4500',
    name: 'Orange Red',
    nameEn: 'Orange Red',
    nameRu: 'Оранжево-красный',
    nameSr: 'Наранџасто црвена',
  },
];

// Получение цветов с локализованными названиями для конкретного языка
export function getLocalizedColors(locale: string): Array<{
  hex: string;
  name: string;
  originalName: string;
}> {
  return LOCALIZED_COLORS.map((color) => {
    let localizedName = color.nameEn; // По умолчанию английский

    switch (locale) {
      case 'ru':
        localizedName = color.nameRu;
        break;
      case 'sr':
        localizedName = color.nameSr;
        break;
      case 'en':
      default:
        localizedName = color.nameEn;
        break;
    }

    return {
      hex: color.hex,
      name: localizedName,
      originalName: color.nameEn, // Для обратной совместимости
    };
  });
}

// Поиск цвета по названию на любом языке
export function findColorByName(name: string): LocalizedColor | undefined {
  const lowerName = name.toLowerCase().trim();

  return LOCALIZED_COLORS.find(
    (color) =>
      color.nameEn.toLowerCase() === lowerName ||
      color.nameRu.toLowerCase() === lowerName ||
      color.nameSr.toLowerCase() === lowerName ||
      color.name.toLowerCase() === lowerName
  );
}

// Получение hex по названию цвета
export function getColorHex(colorName: string): string {
  if (colorName.startsWith('#')) {
    return colorName;
  }

  const color = findColorByName(colorName);
  return color?.hex || '#000000';
}

// Получение локализованного названия цвета
export function getLocalizedColorName(
  colorValue: string,
  locale: string
): string {
  // Если это hex код, пытаемся найти соответствующий цвет
  if (colorValue.startsWith('#')) {
    const color = LOCALIZED_COLORS.find(
      (c) => c.hex.toLowerCase() === colorValue.toLowerCase()
    );
    if (color) {
      switch (locale) {
        case 'ru':
          return color.nameRu;
        case 'sr':
          return color.nameSr;
        default:
          return color.nameEn;
      }
    }
    return colorValue; // Возвращаем hex если не нашли
  }

  // Если это название, пытаемся найти и вернуть локализованное
  const color = findColorByName(colorValue);
  if (color) {
    switch (locale) {
      case 'ru':
        return color.nameRu;
      case 'sr':
        return color.nameSr;
      default:
        return color.nameEn;
    }
  }

  return colorValue; // Возвращаем как есть если не нашли
}

// Проверка является ли строка валидным hex цветом
export function isValidHexColor(color: string): boolean {
  return /^#[0-9A-Fa-f]{6}$/.test(color);
}

// Получение контрастного цвета (черный или белый) для фона
export function getContrastColor(hexColor: string): string {
  // Убираем #
  const hex = hexColor.replace('#', '');

  // Переводим в RGB
  const r = parseInt(hex.substr(0, 2), 16);
  const g = parseInt(hex.substr(2, 2), 16);
  const b = parseInt(hex.substr(4, 2), 16);

  // Считаем яркость
  const brightness = (r * 299 + g * 587 + b * 114) / 1000;

  // Возвращаем контрастный цвет
  return brightness > 128 ? '#000000' : '#ffffff';
}

// Группировка цветов по категориям для UI
export function getColorCategories(locale: string) {
  const colors = getLocalizedColors(locale);

  return {
    basic: colors.filter((c) =>
      ['Black', 'White', 'Red', 'Green', 'Blue', 'Yellow'].includes(
        c.originalName
      )
    ),
    vibrant: colors.filter((c) =>
      ['Pink', 'Orange', 'Purple', 'Cyan', 'Hot Pink', 'Lime Green'].includes(
        c.originalName
      )
    ),
    neutral: colors.filter((c) =>
      ['Gray', 'Brown', 'Silver', 'Gold'].includes(c.originalName)
    ),
    pastel: colors.filter((c) =>
      ['Light Pink', 'Light Green', 'Sky Blue', 'Plum'].includes(c.originalName)
    ),
    dark: colors.filter((c) =>
      ['Brown Red', 'Forest Green', 'Royal Blue', 'Orange Red'].includes(
        c.originalName
      )
    ),
  };
}
