import React, { PropsWithChildren } from 'react';
import { render, RenderOptions } from '@testing-library/react';
import { Provider } from 'react-redux';
import { configureStore } from '@reduxjs/toolkit';
import { NextIntlClientProvider } from 'next-intl';
// import type { RootState } from '@/store';
// import chatReducer from '@/store/slices/chatSlice';
import reviewsReducer from '@/store/slices/reviewsSlice';
import storefrontsReducer from '@/store/slices/storefrontSlice';
import importReducer from '@/store/slices/importSlice';
import productReducer from '@/store/slices/productSlice';
import paymentReducer from '@/store/slices/paymentSlice';
import cartReducer from '@/store/slices/cartSlice';

// Mock messages for next-intl
const messages = {
  marketplace: {
    create: {
      selected: 'Selected',
      selectOptions: 'Select options',
      noOptions: 'No options available',
      from: 'From',
      to: 'To',
      min: 'Minimum',
      max: 'Maximum',
      minGreaterThanMax: 'Minimum value cannot be greater than maximum',
      allowedRange: 'Allowed range',
      minAllowed: 'Minimum allowed',
      maxAllowed: 'Maximum allowed',
    },
  },
  admin: {
    translations: {
      status: 'Translation Status',
      translate: 'Translate',
      notTranslated: 'Not translated',
      verified: 'Verified',
      machineTranslated: 'Machine translated',
      manualTranslated: 'Manually translated',
      allTranslated: 'All languages translated',
    },
  },
  common: {
    loading: 'Loading...',
    edit: 'Edit',
    delete: 'Delete',
    active: 'Active',
    inactive: 'Inactive',
    activate: 'Activate',
    deactivate: 'Deactivate',
    filters: 'Filters',
    showInactive: 'Show inactive',
    noData: 'No data available',
    expand: 'Expand',
    collapse: 'Collapse',
  },
  sections: {
    attributes: 'Attributes',
  },
};

interface ExtendedRenderOptions extends Omit<RenderOptions, 'queries'> {
  preloadedState?: any;
  store?: ReturnType<typeof configureStore>;
  locale?: string;
}

export function renderWithProviders(
  ui: React.ReactElement,
  {
    preloadedState = {},
    store = configureStore({
      reducer: {
        reviews: reviewsReducer,
        storefronts: storefrontsReducer,
        import: importReducer,
        products: productReducer,
        payment: paymentReducer,
        cart: cartReducer,
      } as any,
      preloadedState,
      middleware: (getDefaultMiddleware) =>
        getDefaultMiddleware({
          serializableCheck: false,
        }),
    }),
    locale = 'en',
    ...renderOptions
  }: ExtendedRenderOptions = {}
) {
  function Wrapper({ children }: PropsWithChildren<object>) {
    return (
      <Provider store={store}>
        <NextIntlClientProvider messages={messages} locale={locale}>
          {children}
        </NextIntlClientProvider>
      </Provider>
    );
  }

  return { store, ...render(ui, { wrapper: Wrapper, ...renderOptions }) };
}

// Re-export everything
export * from '@testing-library/react';
export { renderWithProviders as render };
