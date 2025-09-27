import { createSlice, PayloadAction } from '@reduxjs/toolkit';

// Универсальный интерфейс для обратной совместимости
interface CarComparison {
  id: number;
  title: string;
  price: number;
  year?: number;
  make?: string;
  model?: string;
  mileage?: number;
  fuelType?: string;
  transmission?: string;
  engineSize?: string;
  power?: string;
  bodyType?: string;
  color?: string;
  location?: string;
  imageUrl?: string;
  image?: string; // Добавлено для совместимости с универсальным интерфейсом
  vin?: string;
  driveType?: string;
  doors?: number;
  seats?: number;
  condition?: string;
  previousOwners?: number;
  warranty?: string;
  firstRegistration?: string;
  technicalInspection?: string;
  features?: string[];
  // Добавляем универсальные поля
  category?: string;
  attributes?: Record<string, any>;
}

interface CompareState {
  items: CarComparison[];
  maxItems: number;
  isOpen: boolean;
}

const STORAGE_KEY = 'car_comparison_items';

// Load from localStorage
const loadFromStorage = (): CarComparison[] => {
  if (typeof window === 'undefined') return [];

  const saved = localStorage.getItem(STORAGE_KEY);
  if (!saved) return [];

  try {
    const items = JSON.parse(saved);
    return Array.isArray(items) ? items.slice(0, 3) : [];
  } catch {
    return [];
  }
};

// Save to localStorage
const saveToStorage = (items: CarComparison[]) => {
  if (typeof window === 'undefined') return;
  localStorage.setItem(STORAGE_KEY, JSON.stringify(items));
};

const initialState: CompareState = {
  items: [],
  maxItems: 3,
  isOpen: false,
};

const compareSlice = createSlice({
  name: 'compare',
  initialState,
  reducers: {
    initializeCompare: (state) => {
      state.items = loadFromStorage();
    },

    addToCompare: (state, action: PayloadAction<CarComparison>) => {
      const existingItem = state.items.find(
        (item) => item.id === action.payload.id
      );

      if (!existingItem && state.items.length < state.maxItems) {
        state.items.push(action.payload);
        saveToStorage(state.items);
      }
    },

    removeFromCompare: (state, action: PayloadAction<number>) => {
      state.items = state.items.filter((item) => item.id !== action.payload);
      saveToStorage(state.items);

      // Close panel if no items left
      if (state.items.length === 0) {
        state.isOpen = false;
      }
    },

    clearCompare: (state) => {
      state.items = [];
      state.isOpen = false;
      saveToStorage([]);
    },

    toggleComparePanel: (state) => {
      if (state.items.length > 0) {
        state.isOpen = !state.isOpen;
      } else {
        state.isOpen = false;
      }
    },

    openComparePanel: (state) => {
      if (state.items.length > 0) {
        state.isOpen = true;
      }
    },

    closeComparePanel: (state) => {
      state.isOpen = false;
    },

    replaceCompareItem: (
      state,
      action: PayloadAction<{ oldId: number; newItem: CarComparison }>
    ) => {
      const index = state.items.findIndex(
        (item) => item.id === action.payload.oldId
      );
      if (index !== -1) {
        state.items[index] = action.payload.newItem;
        saveToStorage(state.items);
      }
    },
  },
});

export const {
  initializeCompare,
  addToCompare,
  removeFromCompare,
  clearCompare,
  toggleComparePanel,
  openComparePanel,
  closeComparePanel,
  replaceCompareItem,
} = compareSlice.actions;

export default compareSlice.reducer;
