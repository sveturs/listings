import { createSlice, createAsyncThunk, PayloadAction } from '@reduxjs/toolkit';
import type { StorefrontProduct } from '@/types/storefront';
import type { components } from '@/types/generated/api';
import { productApi } from '@/services/productApi';
import { toast } from 'react-hot-toast';

type BulkOperationError =
  components['schemas']['models.BulkOperationError'];

interface ProductState {
  products: StorefrontProduct[];
  selectedIds: number[];
  loading: boolean;
  error: string | null;

  // Фильтры
  filters: {
    search: string;
    categoryId: number | null;
    minPrice: number | null;
    maxPrice: number | null;
    stockStatus: 'all' | 'in_stock' | 'low_stock' | 'out_of_stock';
    isActive: boolean | null;
  };

  // Пагинация
  pagination: {
    page: number;
    limit: number;
    total: number;
    hasMore: boolean;
  };

  // Массовые операции
  bulkOperation: {
    isProcessing: boolean;
    progress: number;
    total: number;
    errors: BulkOperationError[];
    successCount: number;
    currentOperation: 'idle' | 'delete' | 'update' | 'status' | 'export';
  };

  // UI состояния
  ui: {
    isSelectMode: boolean;
    viewMode: 'grid' | 'list' | 'table';
    sortBy: 'name' | 'price' | 'created_at' | 'stock_quantity';
    sortOrder: 'asc' | 'desc';
  };
}

const initialState: ProductState = {
  products: [],
  selectedIds: [],
  loading: false,
  error: null,

  filters: {
    search: '',
    categoryId: null,
    minPrice: null,
    maxPrice: null,
    stockStatus: 'all',
    isActive: null,
  },

  pagination: {
    page: 1,
    limit: 20,
    total: 0,
    hasMore: true,
  },

  bulkOperation: {
    isProcessing: false,
    progress: 0,
    total: 0,
    errors: [],
    successCount: 0,
    currentOperation: 'idle',
  },

  ui: {
    isSelectMode: false,
    viewMode: 'grid',
    sortBy: 'created_at',
    sortOrder: 'desc',
  },
};

// Async thunks для массовых операций
export const bulkDeleteProducts = createAsyncThunk(
  'products/bulkDelete',
  async ({
    storefrontSlug,
    productIds,
  }: {
    storefrontSlug: string;
    productIds: number[];
  }) => {
    const response = await productApi.bulkDelete(storefrontSlug, productIds);

    // Показываем уведомления
    if (response?.deleted && response.deleted.length > 0) {
      toast.success(`Удалено товаров: ${response.deleted.length}`);
    }
    if (response?.failed && response.failed.length > 0) {
      toast.error(`Не удалось удалить: ${response.failed.length}`);
    }

    return response;
  }
);

export const bulkUpdateStatus = createAsyncThunk(
  'products/bulkUpdateStatus',
  async ({
    storefrontSlug,
    productIds,
    isActive,
  }: {
    storefrontSlug: string;
    productIds: number[];
    isActive: boolean;
  }) => {
    const response = await productApi.bulkUpdateStatus(
      storefrontSlug,
      productIds,
      isActive
    );

    // Показываем уведомления
    if (response?.updated && response.updated.length > 0) {
      toast.success(`Обновлено товаров: ${response.updated.length}`);
    }
    if (response?.failed && response.failed.length > 0) {
      toast.error(`Не удалось обновить: ${response.failed.length}`);
    }

    return response;
  }
);

export const exportProducts = createAsyncThunk(
  'products/export',
  async ({
    storefrontSlug,
    productIds,
    format,
  }: {
    storefrontSlug: string;
    productIds?: number[];
    format: 'csv' | 'xml';
  }) => {
    try {
      if (format === 'csv') {
        await productApi.exportToCSV(storefrontSlug, productIds);
      } else {
        await productApi.exportToXML(storefrontSlug, productIds);
      }
      toast.success('Экспорт начат, файл будет загружен автоматически');
    } catch (error) {
      toast.error('Ошибка при экспорте товаров');
      throw error;
    }
  }
);

const productSlice = createSlice({
  name: 'products',
  initialState,
  reducers: {
    // Управление выбором
    toggleProductSelection: (state, action: PayloadAction<number>) => {
      const id = action.payload;
      const index = state.selectedIds.indexOf(id);

      if (index !== -1) {
        state.selectedIds.splice(index, 1);
      } else {
        state.selectedIds.push(id);
      }
    },

    selectAll: (state) => {
      state.selectedIds = state.products.filter((p) => p.id).map((p) => p.id!);
    },

    clearSelection: (state) => {
      state.selectedIds = [];
    },

    selectByFilter: (
      state,
      action: PayloadAction<(product: StorefrontProduct) => boolean>
    ) => {
      const filterFn = action.payload;
      const newSelectedIds: number[] = [];

      state.products.forEach((product) => {
        if (product.id && filterFn(product)) {
          newSelectedIds.push(product.id);
        }
      });

      state.selectedIds = newSelectedIds;
    },

    // UI управление
    toggleSelectMode: (state) => {
      state.ui.isSelectMode = !state.ui.isSelectMode;
      if (!state.ui.isSelectMode) {
        state.selectedIds = [];
      }
    },

    setViewMode: (state, action: PayloadAction<'grid' | 'list' | 'table'>) => {
      state.ui.viewMode = action.payload;
    },

    setSortBy: (
      state,
      action: PayloadAction<{
        sortBy: 'name' | 'price' | 'created_at' | 'stock_quantity';
        sortOrder: 'asc' | 'desc';
      }>
    ) => {
      state.ui.sortBy = action.payload.sortBy;
      state.ui.sortOrder = action.payload.sortOrder;
    },

    // Фильтры
    setFilters: (
      state,
      action: PayloadAction<Partial<ProductState['filters']>>
    ) => {
      state.filters = { ...state.filters, ...action.payload };
      state.pagination.page = 1; // Сброс на первую страницу при изменении фильтров
    },

    resetFilters: (state) => {
      state.filters = initialState.filters;
      state.pagination.page = 1;
    },

    // Продукты
    setProducts: (state, action: PayloadAction<StorefrontProduct[]>) => {
      state.products = action.payload;
    },

    appendProducts: (state, action: PayloadAction<StorefrontProduct[]>) => {
      state.products = [...state.products, ...action.payload];
    },

    updateProduct: (state, action: PayloadAction<StorefrontProduct>) => {
      if (action.payload.id) {
        const index = state.products.findIndex(
          (p) => p.id === action.payload.id
        );
        if (index !== -1) {
          state.products[index] = action.payload;
        }
      }
    },

    removeProducts: (state, action: PayloadAction<number[]>) => {
      const idsToRemove = new Set(action.payload);
      state.products = state.products.filter(
        (p) => p.id && !idsToRemove.has(p.id)
      );

      // Также удаляем из выбранных
      state.selectedIds = state.selectedIds.filter(
        (id) => !idsToRemove.has(id)
      );
    },

    // Прогресс операций
    setBulkOperationProgress: (
      state,
      action: PayloadAction<{
        progress: number;
        total: number;
        successCount?: number;
      }>
    ) => {
      state.bulkOperation.progress = action.payload.progress;
      state.bulkOperation.total = action.payload.total;
      if (action.payload.successCount !== undefined) {
        state.bulkOperation.successCount = action.payload.successCount;
      }
    },

    addBulkOperationError: (
      state,
      action: PayloadAction<BulkOperationError>
    ) => {
      state.bulkOperation.errors.push(action.payload);
    },

    resetBulkOperation: (state) => {
      state.bulkOperation = initialState.bulkOperation;
    },

    // Пагинация
    setPagination: (
      state,
      action: PayloadAction<Partial<ProductState['pagination']>>
    ) => {
      state.pagination = { ...state.pagination, ...action.payload };
    },
  },

  extraReducers: (builder) => {
    // Bulk delete
    builder
      .addCase(bulkDeleteProducts.pending, (state) => {
        state.bulkOperation.isProcessing = true;
        state.bulkOperation.currentOperation = 'delete';
        state.bulkOperation.errors = [];
        state.bulkOperation.successCount = 0;
      })
      .addCase(bulkDeleteProducts.fulfilled, (state, action) => {
        state.bulkOperation.isProcessing = false;
        state.bulkOperation.currentOperation = 'idle';

        // Удаляем успешно удаленные продукты
        if (action.payload?.deleted && action.payload.deleted.length > 0) {
          const idsToRemove = new Set(action.payload.deleted);
          state.products = state.products.filter(
            (p) => p.id && !idsToRemove.has(p.id)
          );

          // Также удаляем из выбранных
          state.selectedIds = state.selectedIds.filter(
            (id) => !idsToRemove.has(id)
          );
        }

        // Сохраняем ошибки
        state.bulkOperation.errors = action.payload?.failed || [];
        state.bulkOperation.successCount = action.payload?.deleted?.length || 0;
      })
      .addCase(bulkDeleteProducts.rejected, (state, action) => {
        state.bulkOperation.isProcessing = false;
        state.bulkOperation.currentOperation = 'idle';
        state.error = action.error.message || 'Failed to delete products';
      });

    // Bulk update status
    builder
      .addCase(bulkUpdateStatus.pending, (state) => {
        state.bulkOperation.isProcessing = true;
        state.bulkOperation.currentOperation = 'status';
        state.bulkOperation.errors = [];
        state.bulkOperation.successCount = 0;
      })
      .addCase(bulkUpdateStatus.fulfilled, (state, action) => {
        state.bulkOperation.isProcessing = false;
        state.bulkOperation.currentOperation = 'idle';

        // Обновляем статус продуктов
        if (action.payload?.updated && action.payload.updated.length > 0) {
          const updatedIds = new Set(action.payload.updated);
          state.products = state.products.map((product) => {
            if (product.id && updatedIds.has(product.id)) {
              return { ...product, is_active: action.meta.arg.isActive };
            }
            return product;
          });
        }

        // Сохраняем ошибки
        state.bulkOperation.errors = action.payload?.failed || [];
        state.bulkOperation.successCount = action.payload?.updated?.length || 0;
      })
      .addCase(bulkUpdateStatus.rejected, (state, action) => {
        state.bulkOperation.isProcessing = false;
        state.bulkOperation.currentOperation = 'idle';
        state.error = action.error.message || 'Failed to update product status';
      });
  },
});

export const {
  toggleProductSelection,
  selectAll,
  clearSelection,
  selectByFilter,
  toggleSelectMode,
  setViewMode,
  setSortBy,
  setFilters,
  resetFilters,
  setProducts,
  appendProducts,
  updateProduct,
  removeProducts,
  setBulkOperationProgress,
  addBulkOperationError,
  resetBulkOperation,
  setPagination,
} = productSlice.actions;

export default productSlice.reducer;
