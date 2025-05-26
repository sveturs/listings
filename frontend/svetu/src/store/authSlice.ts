import { createSlice, createAsyncThunk, PayloadAction } from '@reduxjs/toolkit';
import { AuthState, User, LoginRequest, RegisterRequest } from '@/types/auth';
import authService from '@/services/auth.service';

const initialState: AuthState = {
  user: null,
  token: null,
  isLoading: false,
  isAuthenticated: false,
  error: null,
};

// Async thunks
export const loginUser = createAsyncThunk(
  'auth/login',
  async (credentials: LoginRequest) => {
    const response = await authService.login(credentials);
    return response;
  }
);

export const registerUser = createAsyncThunk(
  'auth/register',
  async (data: RegisterRequest) => {
    const response = await authService.register(data);
    return response;
  }
);

export const logoutUser = createAsyncThunk(
  'auth/logout',
  async () => {
    await authService.logout();
  }
);

export const checkAuth = createAsyncThunk(
  'auth/checkAuth',
  async () => {
    const response = await authService.checkSession();
    return response;
  }
);

export const googleLogin = createAsyncThunk(
  'auth/googleLogin',
  async (params?: string) => {
    // This will handle the redirect to Google OAuth
    const baseUrl = process.env.NEXT_PUBLIC_API_URL || '';
    const authUrl = '/auth/google';
    window.location.href = `${baseUrl}${authUrl}${params || ''}`;
  }
);

const authSlice = createSlice({
  name: 'auth',
  initialState,
  reducers: {
    updateUser: (state, action: PayloadAction<User>) => {
      state.user = action.payload;
    },
    clearError: (state) => {
      state.error = null;
    },
    setAuthFromSession: (state, action: PayloadAction<{ user: User; token?: string }>) => {
      state.user = action.payload.user;
      state.token = action.payload.token || null;
      state.isAuthenticated = true;
      state.isLoading = false;
    },
  },
  extraReducers: (builder) => {
    // Login
    builder
      .addCase(loginUser.pending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(loginUser.fulfilled, (state, action) => {
        state.isLoading = false;
        state.isAuthenticated = true;
        state.user = action.payload.user;
        state.token = action.payload.token;
        state.error = null;
      })
      .addCase(loginUser.rejected, (state, action) => {
        state.isLoading = false;
        state.isAuthenticated = false;
        state.user = null;
        state.token = null;
        state.error = action.error.message || 'Login failed';
      });
    
    // Register
    builder
      .addCase(registerUser.pending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(registerUser.fulfilled, (state, action) => {
        state.isLoading = false;
        state.isAuthenticated = true;
        state.user = action.payload.user;
        state.token = action.payload.token;
        state.error = null;
      })
      .addCase(registerUser.rejected, (state, action) => {
        state.isLoading = false;
        state.error = action.error.message || 'Registration failed';
      });
    
    // Logout
    builder
      .addCase(logoutUser.fulfilled, (state) => {
        state.user = null;
        state.token = null;
        state.isAuthenticated = false;
        state.error = null;
      });
    
    // Check Auth
    builder
      .addCase(checkAuth.pending, (state) => {
        state.isLoading = true;
      })
      .addCase(checkAuth.fulfilled, (state, action) => {
        if (action.payload.authenticated && action.payload.user) {
          state.user = action.payload.user;
          state.isAuthenticated = true;
        }
        state.isLoading = false;
      })
      .addCase(checkAuth.rejected, (state) => {
        state.isLoading = false;
      });
  },
});

export const {
  updateUser,
  clearError,
  setAuthFromSession,
} = authSlice.actions;

export default authSlice.reducer;