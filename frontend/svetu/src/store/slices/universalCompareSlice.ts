import { createSlice, PayloadAction } from '@reduxjs/toolkit';

// Универсальный интерфейс для сравнения
export interface UniversalCompareItem {
  id: number;
  category: string; // 'cars', 'real_estate', 'electronics', etc.
  title: string;
  price: number;
  currency?: string;
  image?: string;
  location?: string;
  // Универсальные атрибуты - динамические для разных категорий
  attributes: Record<string, any>;
  // Поля для сравнения (определяются категорией)
  compareFields?: string[];
}

// Конфигурация для разных категорий
export const COMPARE_CONFIG = {
  cars: {
    maxItems: 3,
    storageKey: 'compare_cars',
    compareFields: [
      'year',
      'make',
      'model',
      'mileage',
      'fuelType',
      'transmission',
      'engineSize',
      'power',
      'bodyType',
      'driveType',
      'condition',
      'previousOwners',
      'warranty',
      'features',
    ],
    requiredFields: ['year', 'make', 'model'],
  },
  real_estate: {
    maxItems: 4,
    storageKey: 'compare_real_estate',
    compareFields: [
      'propertyType',
      'area',
      'rooms',
      'bedrooms',
      'bathrooms',
      'floor',
      'totalFloors',
      'buildingYear',
      'renovation',
      'heating',
      'parking',
      'balcony',
      'elevator',
      'security',
      'features',
    ],
    requiredFields: ['propertyType', 'area', 'rooms'],
  },
  electronics: {
    maxItems: 5,
    storageKey: 'compare_electronics',
    compareFields: [
      'brand',
      'model',
      'category',
      'condition',
      'warranty',
      'specs',
      'color',
      'storage',
      'memory',
      'screenSize',
      'processor',
      'features',
    ],
    requiredFields: ['brand', 'model'],
  },
  marketplace: {
    maxItems: 3,
    storageKey: 'compare_marketplace',
    compareFields: ['category', 'condition', 'brand', 'features'],
    requiredFields: ['title'],
  },
};

interface CompareState {
  // Разделяем элементы по категориям
  itemsByCategory: Record<string, UniversalCompareItem[]>;
  // Текущая активная категория для сравнения
  activeCategory: string | null;
  // Конфигурация для каждой категории
  config: typeof COMPARE_CONFIG;
  // Состояние панели сравнения
  isPanelOpen: boolean;
}

// Загрузка из localStorage для конкретной категории
const loadFromStorage = (category: string): UniversalCompareItem[] => {
  if (typeof window === 'undefined') return [];

  const config = COMPARE_CONFIG[category as keyof typeof COMPARE_CONFIG];
  if (!config) return [];

  const saved = localStorage.getItem(config.storageKey);
  if (!saved) return [];

  try {
    const items = JSON.parse(saved);
    return Array.isArray(items) ? items.slice(0, config.maxItems) : [];
  } catch {
    return [];
  }
};

// Сохранение в localStorage для конкретной категории
const saveToStorage = (category: string, items: UniversalCompareItem[]) => {
  if (typeof window === 'undefined') return;

  const config = COMPARE_CONFIG[category as keyof typeof COMPARE_CONFIG];
  if (!config) return;

  localStorage.setItem(config.storageKey, JSON.stringify(items));
};

// Загрузка всех категорий из localStorage
const loadAllFromStorage = (): Record<string, UniversalCompareItem[]> => {
  const result: Record<string, UniversalCompareItem[]> = {};

  Object.keys(COMPARE_CONFIG).forEach((category) => {
    const items = loadFromStorage(category);
    if (items.length > 0) {
      result[category] = items;
    }
  });

  return result;
};

const initialState: CompareState = {
  itemsByCategory: {},
  activeCategory: null,
  config: COMPARE_CONFIG,
  isPanelOpen: false,
};

const universalCompareSlice = createSlice({
  name: 'universalCompare',
  initialState,
  reducers: {
    // Инициализация при загрузке приложения
    initializeCompare: (state) => {
      state.itemsByCategory = loadAllFromStorage();
      // Устанавливаем активную категорию, если есть элементы
      const categoriesWithItems = Object.keys(state.itemsByCategory);
      if (categoriesWithItems.length > 0) {
        state.activeCategory = categoriesWithItems[0];
      }
    },

    // Добавление элемента в сравнение
    addItem: (state, action: PayloadAction<UniversalCompareItem>) => {
      const { category, id } = action.payload;
      const config = state.config[category as keyof typeof COMPARE_CONFIG];

      if (!config) return;

      // Инициализируем массив для категории, если его нет
      if (!state.itemsByCategory[category]) {
        state.itemsByCategory[category] = [];
      }

      const categoryItems = state.itemsByCategory[category];

      // Проверяем, не добавлен ли уже этот элемент
      const exists = categoryItems.some((item) => item.id === id);

      if (!exists && categoryItems.length < config.maxItems) {
        categoryItems.push(action.payload);
        saveToStorage(category, categoryItems);

        // Устанавливаем эту категорию как активную
        state.activeCategory = category;
      }
    },

    // Удаление элемента из сравнения
    removeItem: (state, action: PayloadAction<number>) => {
      const id = action.payload;

      // Ищем во всех категориях
      Object.keys(state.itemsByCategory).forEach((category) => {
        state.itemsByCategory[category] = state.itemsByCategory[
          category
        ].filter((item) => item.id !== id);

        // Сохраняем изменения
        saveToStorage(category, state.itemsByCategory[category]);

        // Удаляем категорию, если в ней не осталось элементов
        if (state.itemsByCategory[category].length === 0) {
          delete state.itemsByCategory[category];
        }
      });

      // Обновляем активную категорию
      const remainingCategories = Object.keys(state.itemsByCategory);
      if (remainingCategories.length > 0) {
        if (
          state.activeCategory &&
          !state.itemsByCategory[state.activeCategory]
        ) {
          state.activeCategory = remainingCategories[0];
        }
      } else {
        state.activeCategory = null;
        state.isPanelOpen = false;
      }
    },

    // Удаление элемента из конкретной категории
    removeItemFromCategory: (
      state,
      action: PayloadAction<{ category: string; id: number }>
    ) => {
      const { category, id } = action.payload;

      if (state.itemsByCategory[category]) {
        state.itemsByCategory[category] = state.itemsByCategory[
          category
        ].filter((item) => item.id !== id);

        saveToStorage(category, state.itemsByCategory[category]);

        if (state.itemsByCategory[category].length === 0) {
          delete state.itemsByCategory[category];

          // Обновляем активную категорию
          const remainingCategories = Object.keys(state.itemsByCategory);
          if (remainingCategories.length > 0) {
            state.activeCategory = remainingCategories[0];
          } else {
            state.activeCategory = null;
            state.isPanelOpen = false;
          }
        }
      }
    },

    // Очистка всех элементов категории
    clearCategory: (state, action: PayloadAction<string>) => {
      const category = action.payload;

      if (state.itemsByCategory[category]) {
        delete state.itemsByCategory[category];
        saveToStorage(category, []);

        // Обновляем активную категорию
        const remainingCategories = Object.keys(state.itemsByCategory);
        if (remainingCategories.length > 0) {
          state.activeCategory = remainingCategories[0];
        } else {
          state.activeCategory = null;
          state.isPanelOpen = false;
        }
      }
    },

    // Очистка всех элементов
    clearAll: (state) => {
      Object.keys(state.itemsByCategory).forEach((category) => {
        saveToStorage(category, []);
      });

      state.itemsByCategory = {};
      state.activeCategory = null;
      state.isPanelOpen = false;
    },

    // Установка активной категории
    setActiveCategory: (state, action: PayloadAction<string>) => {
      if (state.itemsByCategory[action.payload]) {
        state.activeCategory = action.payload;
      }
    },

    // Управление панелью сравнения
    togglePanel: (state) => {
      if (Object.keys(state.itemsByCategory).length > 0) {
        state.isPanelOpen = !state.isPanelOpen;
      } else {
        state.isPanelOpen = false;
      }
    },

    openPanel: (state) => {
      if (Object.keys(state.itemsByCategory).length > 0) {
        state.isPanelOpen = true;
      }
    },

    closePanel: (state) => {
      state.isPanelOpen = false;
    },

    // Замена элемента в категории
    replaceItem: (
      state,
      action: PayloadAction<{
        category: string;
        oldId: number;
        newItem: UniversalCompareItem;
      }>
    ) => {
      const { category, oldId, newItem } = action.payload;

      if (state.itemsByCategory[category]) {
        const index = state.itemsByCategory[category].findIndex(
          (item) => item.id === oldId
        );

        if (index !== -1) {
          state.itemsByCategory[category][index] = newItem;
          saveToStorage(category, state.itemsByCategory[category]);
        }
      }
    },
  },
});

// Селекторы
export const selectCompareItems = (
  state: { universalCompare: CompareState },
  category?: string
) => {
  if (category) {
    return state.universalCompare.itemsByCategory[category] || [];
  }
  return state.universalCompare.activeCategory
    ? state.universalCompare.itemsByCategory[
        state.universalCompare.activeCategory
      ] || []
    : [];
};

export const selectAllCompareItems = (state: {
  universalCompare: CompareState;
}) => {
  return state.universalCompare.itemsByCategory;
};

export const selectCompareCount = (
  state: { universalCompare: CompareState },
  category?: string
) => {
  if (category) {
    return state.universalCompare.itemsByCategory[category]?.length || 0;
  }

  return Object.values(state.universalCompare.itemsByCategory).reduce(
    (sum, items) => sum + items.length,
    0
  );
};

export const selectIsInCompare = (
  state: { universalCompare: CompareState },
  id: number,
  category?: string
) => {
  if (category) {
    return (
      state.universalCompare.itemsByCategory[category]?.some(
        (item) => item.id === id
      ) || false
    );
  }

  return Object.values(state.universalCompare.itemsByCategory).some((items) =>
    items.some((item) => item.id === id)
  );
};

export const selectActiveCategory = (state: {
  universalCompare: CompareState;
}) => {
  return state.universalCompare.activeCategory;
};

export const selectCompareConfig = (
  state: { universalCompare: CompareState },
  category: string
) => {
  return state.universalCompare.config[category as keyof typeof COMPARE_CONFIG];
};

export const {
  initializeCompare,
  addItem,
  removeItem,
  removeItemFromCategory,
  clearCategory,
  clearAll,
  setActiveCategory,
  togglePanel,
  openPanel,
  closePanel,
  replaceItem,
} = universalCompareSlice.actions;

export default universalCompareSlice.reducer;
