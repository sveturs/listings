import { apiClient } from '@/lib/api-client';
import { LoginRequest, LoginResponse, RegisterRequest, RegisterResponse, User } from '@/types/auth';
import { ApiResponse } from '@/types/api';

class AuthService {
  async login(credentials: LoginRequest): Promise<LoginResponse> {
    const response = await apiClient.post<ApiResponse<LoginResponse>>('/auth/login', credentials);
    const { user, token } = response.data.data;
    
    // Store token in apiClient
    apiClient.setAuthToken(token);
    
    return { user, token };
  }

  async register(data: RegisterRequest): Promise<RegisterResponse> {
    const response = await apiClient.post<ApiResponse<RegisterResponse>>('/auth/register', data);
    const { user, token } = response.data.data;
    
    // Store token in apiClient
    apiClient.setAuthToken(token);
    
    return { user, token };
  }

  async logout(): Promise<void> {
    try {
      await apiClient.get('/auth/logout');
    } catch (error) {
      // Even if logout fails on server, clear local data
      console.error('Logout error:', error);
    } finally {
      // Clear local storage
      if (typeof window !== 'undefined') {
        localStorage.removeItem('authToken');
        localStorage.removeItem('user');
        localStorage.removeItem('user_session');
        localStorage.removeItem('user_session_token');
      }
      // Clear token from apiClient
      apiClient.setAuthToken(null);
    }
  }

  async checkSession(): Promise<{ authenticated: boolean; user?: User }> {
    try {
      // First check URL for session token (from OAuth redirect)
      if (typeof window !== 'undefined') {
        const urlParams = new URLSearchParams(window.location.search);
        const sessionTokenFromUrl = urlParams.get('session_token');
        
        if (sessionTokenFromUrl) {
          localStorage.setItem('user_session_token', sessionTokenFromUrl);
          // Clean up URL
          window.history.replaceState({}, document.title, window.location.pathname);
        }
      }

      const response = await apiClient.get<ApiResponse<{ authenticated: boolean; user?: User }>>('/auth/session');
      
      if (response.data?.data?.authenticated && response.data?.data?.user) {
        // Save user session
        if (typeof window !== 'undefined') {
          localStorage.setItem('user_session', JSON.stringify(response.data.data.user));
        }
      }
      
      return response.data?.data || { authenticated: false };
    } catch (error) {
      console.error('Error checking auth status:', error);
      return { authenticated: false };
    }
  }

  async getCurrentUser(): Promise<User> {
    const response = await apiClient.get<ApiResponse<User>>('/auth/me');
    return response.data.data;
  }

  async updateProfile(data: Partial<User>): Promise<User> {
    const response = await apiClient.put<ApiResponse<User>>('/auth/profile', data);
    return response.data.data;
  }

  async changePassword(oldPassword: string, newPassword: string): Promise<void> {
    await apiClient.post('/auth/change-password', { oldPassword, newPassword });
  }

  async requestPasswordReset(email: string): Promise<void> {
    await apiClient.post('/auth/forgot-password', { email });
  }

  async resetPassword(token: string, newPassword: string): Promise<void> {
    await apiClient.post('/auth/reset-password', { token, newPassword });
  }
}

const authService = new AuthService();

export { authService };
export default authService;