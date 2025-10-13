interface LoginCredentials {
  email: string;
  password: string;
}

interface RegisterData {
  email: string;
  password: string;
  name: string;
  terms_accepted?: boolean;
}

interface User {
  id: number;
  email: string;
  name?: string;
  email_verified?: boolean;
  phone_verified?: boolean;
  two_factor_enabled?: boolean;
  is_admin?: boolean;
  roles?: string[];
  created_at?: string;
  updated_at?: string;
}

interface AuthResponse {
  success: boolean;
  user?: User;
  error?: string;
  message?: string;
}

class AuthNewService {
  async login(credentials: LoginCredentials): Promise<AuthResponse> {
    try {
      const response = await fetch('/api/auth/login', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(credentials),
      });

      const data = await response.json();

      if (!response.ok) {
        throw new Error(data.error || data.message || 'Login failed');
      }

      return data;
    } catch (error) {
      console.error('Login error:', error);
      throw error;
    }
  }

  async register(data: RegisterData): Promise<AuthResponse> {
    try {
      const response = await fetch('/api/auth/register', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          ...data,
          terms_accepted: true, // Всегда отправляем true для auth-service
        }),
      });

      const result = await response.json();

      if (!response.ok) {
        throw new Error(
          result.error || result.message || 'Registration failed'
        );
      }

      return result;
    } catch (error) {
      console.error('Registration error:', error);
      throw error;
    }
  }

  async logout(): Promise<void> {
    try {
      await fetch('/api/auth/logout', {
        method: 'POST',
      });
    } catch (error) {
      console.error('Logout error:', error);
      // Игнорируем ошибки logout
    }
  }

  async getSession(): Promise<User | null> {
    try {
      const response = await fetch('/api/auth/session', {
        method: 'GET',
      });

      if (!response.ok) {
        return null;
      }

      const data = await response.json();
      return data.user || null;
    } catch (error) {
      console.error('Session error:', error);
      return null;
    }
  }

  async refreshToken(): Promise<AuthResponse> {
    try {
      const response = await fetch('/api/auth/refresh', {
        method: 'POST',
      });

      const data = await response.json();

      if (!response.ok) {
        throw new Error(data.error || 'Token refresh failed');
      }

      return data;
    } catch (error) {
      console.error('Refresh error:', error);
      throw error;
    }
  }

  async updateProfile(data: Partial<User>): Promise<User> {
    try {
      const response = await fetch('/api/v2/users/profile', {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        credentials: 'include', // Отправляем httpOnly cookies
        body: JSON.stringify(data),
      });

      if (!response.ok) {
        const error = await response.json();
        throw new Error(
          error.error || error.message || 'Failed to update profile'
        );
      }

      const result = await response.json();
      return result.user || result.data || result;
    } catch (error) {
      console.error('Update profile error:', error);
      throw error;
    }
  }
}

export const authService = new AuthNewService();

// Для обратной совместимости со старым кодом
export class AuthService {
  static async login(credentials: LoginCredentials) {
    return authService.login(credentials);
  }

  static async register(data: RegisterData) {
    return authService.register(data);
  }

  static async logout() {
    return authService.logout();
  }

  static async getSession() {
    const user = await authService.getSession();
    return {
      authenticated: !!user,
      user: user || undefined,
    };
  }

  static async refreshToken() {
    return authService.refreshToken();
  }
}

export default authService;
