import { createSlice, createAsyncThunk } from '@reduxjs/toolkit';
import { apiClient } from '@/services/api-client';

// Типы
export interface Category {
  id: number;
  slug: string;
  name: string;
  description?: string;
  parent_id?: number;
  iconName?: string; // Хранить имя иконки как строку
  color?: string;
  count?: number;
  listing_count?: number;
  translations?: {
    en: string;
    ru: string;
    sr: string;
  };
}

interface CategoriesState {
  categories: Category[];
  popularCategories: Category[];
  isLoadingCategories: boolean;
  isLoadingPopular: boolean;
  error: string | null;
  lastLocale: string | null;
}

// Начальное состояние
const initialState: CategoriesState = {
  categories: [],
  popularCategories: [],
  isLoadingCategories: false,
  isLoadingPopular: false,
  error: null,
  lastLocale: null,
};

// Маппинг имен иконок для популярных категорий (только строки)
const iconNameMap: { [key: string]: string } = {
  'real-estate': 'BsHouseDoor',
  automotive: 'FaCar',
  electronics: 'BsLaptop',
  fashion: 'FaTshirt',
  jobs: 'BsBriefcase',
  services: 'BsTools',
  'hobbies-entertainment': 'BsPalette',
  'home-garden': 'BsHandbag',
  industrial: 'BsTools',
  'food-beverages': 'BsPhone',
  'books-stationery': 'BsGem',
  'antiques-art': 'BsPalette',
};

const colorMap: { [key: string]: string } = {
  'real-estate': 'text-blue-600',
  automotive: 'text-red-600',
  electronics: 'text-purple-600',
  fashion: 'text-pink-600',
  jobs: 'text-green-600',
  services: 'text-orange-600',
  'hobbies-entertainment': 'text-indigo-600',
  'home-garden': 'text-yellow-600',
  industrial: 'text-gray-600',
  'food-beverages': 'text-teal-600',
  'books-stationery': 'text-cyan-600',
  'antiques-art': 'text-rose-600',
};

// Async thunks
export const fetchCategories = createAsyncThunk(
  'categories/fetchCategories',
  async () => {
    const response = await apiClient.get('/marketplace/categories');
    if (response.data.success && response.data.data) {
      return response.data.data;
    }
    return [];
  }
);

export const fetchPopularCategories = createAsyncThunk(
  'categories/fetchPopularCategories',
  async (
    {
      locale,
      forceRefresh = false,
    }: { locale: string; forceRefresh?: boolean },
    { getState }
  ) => {
    const state = getState() as { categories: CategoriesState };

    // Если не принудительное обновление и локаль не изменилась, не делаем запрос
    if (
      !forceRefresh &&
      state.categories.lastLocale === locale &&
      state.categories.popularCategories.length > 0
    ) {
      return state.categories.popularCategories;
    }

    const response = await apiClient.get(
      `/marketplace/popular-categories?lang=${locale}&limit=8`
    );

    if (response.data.success && response.data.data) {
      // Добавляем имена иконок и цвета к категориям
      const categoriesWithIcons = response.data.data.map((cat: any) => ({
        ...cat,
        iconName: iconNameMap[cat.slug] || 'BsHandbag',
        color: colorMap[cat.slug] || 'text-gray-600',
        count: cat.count || 0,
      }));

      return { categories: categoriesWithIcons, locale };
    }

    return { categories: [], locale };
  }
);

// Slice
const categoriesSlice = createSlice({
  name: 'categories',
  initialState,
  reducers: {
    clearCategories: (state) => {
      state.categories = [];
      state.popularCategories = [];
      state.error = null;
      state.lastLocale = null;
    },
  },
  extraReducers: (builder) => {
    // Fetch categories
    builder
      .addCase(fetchCategories.pending, (state) => {
        state.isLoadingCategories = true;
        state.error = null;
      })
      .addCase(fetchCategories.fulfilled, (state, action) => {
        state.isLoadingCategories = false;
        state.categories = action.payload;
      })
      .addCase(fetchCategories.rejected, (state, action) => {
        state.isLoadingCategories = false;
        state.error = action.error.message || 'Failed to load categories';
      })
      // Fetch popular categories
      .addCase(fetchPopularCategories.pending, (state) => {
        state.isLoadingPopular = true;
        state.error = null;
      })
      .addCase(fetchPopularCategories.fulfilled, (state, action) => {
        state.isLoadingPopular = false;
        // Проверяем формат данных в payload
        if ('categories' in action.payload && action.payload.categories) {
          state.popularCategories = action.payload.categories;
          state.lastLocale = action.payload.locale;
        } else {
          // Если вернулись уже кэшированные данные
          state.popularCategories = action.payload as Category[];
        }
      })
      .addCase(fetchPopularCategories.rejected, (state, action) => {
        state.isLoadingPopular = false;
        state.error =
          action.error.message || 'Failed to load popular categories';
      });
  },
});

export const { clearCategories } = categoriesSlice.actions;
export default categoriesSlice.reducer;
