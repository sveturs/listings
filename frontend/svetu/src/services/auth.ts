import configManager from '@/config';
import type {
  SessionResponse,
  UserProfile,
  UpdateProfileRequest,
} from '@/types/auth';

const API_BASE = configManager.getApiUrl();

export class AuthService {
  static async getSession(): Promise<SessionResponse> {
    try {
      const response = await fetch(`${API_BASE}/auth/session`, {
        method: 'GET',
        credentials: 'include',
      });

      if (!response.ok) {
        throw new Error('Failed to fetch session');
      }

      return await response.json();
    } catch (error) {
      console.error('Session fetch error:', error);
      return { authenticated: false };
    }
  }

  static async logout(): Promise<void> {
    try {
      await fetch(`${API_BASE}/auth/logout`, {
        method: 'GET',
        credentials: 'include',
      });
    } catch (error) {
      console.error('Logout error:', error);
    }
  }

  static async loginWithGoogle(returnTo?: string): Promise<void> {
    const params = returnTo ? `?returnTo=${encodeURIComponent(returnTo)}` : '';
    window.location.href = `${API_BASE}/auth/google${params}`;
  }

  static async getProfile(): Promise<UserProfile | null> {
    try {
      const response = await fetch(`${API_BASE}/api/v1/users/me`, {
        method: 'GET',
        credentials: 'include',
      });

      if (!response.ok) {
        return null;
      }

      return await response.json();
    } catch (error) {
      console.error('Profile fetch error:', error);
      return null;
    }
  }

  static async updateProfile(
    data: UpdateProfileRequest
  ): Promise<UserProfile | null> {
    try {
      const response = await fetch(`${API_BASE}/api/v1/users/me`, {
        method: 'PUT',
        credentials: 'include',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(data),
      });

      if (!response.ok) {
        throw new Error('Failed to update profile');
      }

      return await response.json();
    } catch (error) {
      console.error('Profile update error:', error);
      return null;
    }
  }
}
