import { createSlice, createAsyncThunk, PayloadAction } from '@reduxjs/toolkit';
import { apiClient } from '@/services/api-client';

// TODO: Update these types once payment endpoints are added to the backend
interface Transaction {
  id: string;
  payment_id: string;
  type: 'payment' | 'withdrawal' | 'refund';
  amount: number;
  currency: string;
  status: 'pending' | 'completed' | 'failed' | 'cancelled';
  created_at: string;
  updated_at: string;
  description?: string;
}

interface PaymentMethod {
  id: string;
  name: string;
  type: 'card' | 'cash_on_delivery' | 'bank_transfer';
  enabled: boolean;
  fee_percentage?: number;
  fee_fixed?: number;
}

interface PaymentState {
  // Current checkout
  checkoutData: {
    listingId: string | null;
    amount: number;
    currency: string;
    paymentMethod: string | null;
    commission: number;
    total: number;
  } | null;

  // Transactions
  transactions: Transaction[];
  transactionsLoading: boolean;
  transactionsError: string | null;

  // Wallet
  wallet: {
    balance: number;
    pendingBalance: number;
    currency: string;
  } | null;
  walletLoading: boolean;

  // Payment methods
  paymentMethods: PaymentMethod[];
  paymentMethodsLoading: boolean;

  // Current payment process
  paymentProcessing: boolean;
  paymentError: string | null;
  lastPaymentId: string | null;
}

const initialState: PaymentState = {
  checkoutData: null,
  transactions: [],
  transactionsLoading: false,
  transactionsError: null,
  wallet: null,
  walletLoading: false,
  paymentMethods: [],
  paymentMethodsLoading: false,
  paymentProcessing: false,
  paymentError: null,
  lastPaymentId: null,
};

// Async thunks
export const createPayment = createAsyncThunk(
  'payment/create',
  async (data: {
    listingId: string;
    amount: number;
    paymentMethod: string;
    buyerInfo: {
      name: string;
      email: string;
      phone: string;
      address?: string;
    };
  }) => {
    const response = await apiClient.post('/api/v1/payments/create', {
      listing_id: data.listingId,
      amount: data.amount,
      payment_method: data.paymentMethod,
      buyer_info: data.buyerInfo,
    });
    return response.data;
  }
);

export const fetchTransactions = createAsyncThunk(
  'payment/fetchTransactions',
  async (params?: { status?: string; limit?: number; offset?: number }) => {
    const queryParams = new URLSearchParams();
    if (params?.status) queryParams.append('status', params.status);
    if (params?.limit) queryParams.append('limit', params.limit.toString());
    if (params?.offset) queryParams.append('offset', params.offset.toString());

    const url = `/api/v1/payments/transactions${queryParams.toString() ? `?${queryParams.toString()}` : ''}`;
    const response = await apiClient.get(url);
    return response.data;
  }
);

export const fetchWallet = createAsyncThunk('payment/fetchWallet', async () => {
  const response = await apiClient.get('/api/v1/payments/wallet');
  return response.data;
});

export const fetchPaymentMethods = createAsyncThunk(
  'payment/fetchPaymentMethods',
  async () => {
    const response = await apiClient.get('/api/v1/payments/methods');
    return response.data;
  }
);

export const requestWithdrawal = createAsyncThunk(
  'payment/requestWithdrawal',
  async (data: { amount: number; method: string; details: any }) => {
    const response = await apiClient.post('/api/v1/payments/withdraw', data);
    return response.data;
  }
);

export const confirmPayment = createAsyncThunk(
  'payment/confirm',
  async (paymentId: string) => {
    const response = await apiClient.post(
      `/api/v1/payments/${paymentId}/confirm`
    );
    return response.data;
  }
);

export const refundPayment = createAsyncThunk(
  'payment/refund',
  async (data: { paymentId: string; reason: string }) => {
    const response = await apiClient.post(
      `/api/v1/payments/${data.paymentId}/refund`,
      {
        reason: data.reason,
      }
    );
    return response.data;
  }
);

const paymentSlice = createSlice({
  name: 'payment',
  initialState,
  reducers: {
    setCheckoutData: (
      state,
      action: PayloadAction<{
        listingId: string;
        amount: number;
        currency: string;
      }>
    ) => {
      const { amount } = action.payload;
      // Calculate commission based on category (simplified for now)
      const commissionRate = 0.05; // 5% default
      const commission = amount * commissionRate;
      const total = amount + commission;

      state.checkoutData = {
        ...action.payload,
        paymentMethod: null,
        commission,
        total,
      };
    },

    setPaymentMethod: (state, action: PayloadAction<string>) => {
      if (state.checkoutData) {
        state.checkoutData.paymentMethod = action.payload;
        // Add 2% for cash on delivery
        if (action.payload === 'cash_on_delivery') {
          const extraCharge = state.checkoutData.amount * 0.02;
          state.checkoutData.total =
            state.checkoutData.amount +
            state.checkoutData.commission +
            extraCharge;
        }
      }
    },

    clearCheckout: (state) => {
      state.checkoutData = null;
      state.paymentError = null;
    },

    clearPaymentError: (state) => {
      state.paymentError = null;
    },
  },

  extraReducers: (builder) => {
    // Create payment
    builder
      .addCase(createPayment.pending, (state) => {
        state.paymentProcessing = true;
        state.paymentError = null;
      })
      .addCase(createPayment.fulfilled, (state, action) => {
        state.paymentProcessing = false;
        state.lastPaymentId = action.payload.data.payment_id;
      })
      .addCase(createPayment.rejected, (state, action) => {
        state.paymentProcessing = false;
        state.paymentError = action.error.message || 'Payment failed';
      });

    // Fetch transactions
    builder
      .addCase(fetchTransactions.pending, (state) => {
        state.transactionsLoading = true;
        state.transactionsError = null;
      })
      .addCase(fetchTransactions.fulfilled, (state, action) => {
        state.transactionsLoading = false;
        state.transactions = action.payload.data.transactions;
      })
      .addCase(fetchTransactions.rejected, (state, action) => {
        state.transactionsLoading = false;
        state.transactionsError =
          action.error.message || 'Failed to fetch transactions';
      });

    // Fetch wallet
    builder
      .addCase(fetchWallet.pending, (state) => {
        state.walletLoading = true;
      })
      .addCase(fetchWallet.fulfilled, (state, action) => {
        state.walletLoading = false;
        state.wallet = action.payload.data;
      })
      .addCase(fetchWallet.rejected, (state) => {
        state.walletLoading = false;
      });

    // Fetch payment methods
    builder
      .addCase(fetchPaymentMethods.pending, (state) => {
        state.paymentMethodsLoading = true;
      })
      .addCase(fetchPaymentMethods.fulfilled, (state, action) => {
        state.paymentMethodsLoading = false;
        state.paymentMethods = action.payload.data.methods;
      })
      .addCase(fetchPaymentMethods.rejected, (state) => {
        state.paymentMethodsLoading = false;
      });
  },
});

export const {
  setCheckoutData,
  setPaymentMethod,
  clearCheckout,
  clearPaymentError,
} = paymentSlice.actions;

export default paymentSlice.reducer;
