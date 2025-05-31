import configManager from '@/config';
import type {
  SessionResponse,
  UserUpdate,
  UpdateProfileRequest,
} from '@/types/auth';

const API_BASE = configManager.getApiUrl();

export class AuthService {
  private static abortControllers = new Map<string, AbortController>();

  static cleanup(): void {
    this.abortControllers.forEach((controller) => controller.abort());
    this.abortControllers.clear();
  }

  private static getAbortController(key: string): AbortController {
    // Cancel any existing request with the same key
    const existing = this.abortControllers.get(key);
    if (existing) {
      existing.abort();
    }

    // Create new controller
    const controller = new AbortController();
    this.abortControllers.set(key, controller);
    return controller;
  }

  static async getSession(): Promise<SessionResponse> {
    const controller = this.getAbortController('session');

    try {
      const response = await fetch(`${API_BASE}/auth/session`, {
        method: 'GET',
        credentials: 'include',
        signal: controller.signal,
      });

      if (!response.ok) {
        throw new Error('Failed to fetch session');
      }

      const data = await response.json();
      this.abortControllers.delete('session');
      return data;
    } catch (error) {
      if (error instanceof Error && error.name === 'AbortError') {
        console.log('Session request was cancelled');
        return { authenticated: false };
      }
      console.error('Session fetch error:', error);
      throw error;
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

  static async loginWithGoogle(
    returnTo?: string,
    redirect = true
  ): Promise<string> {
    const params = returnTo ? `?returnTo=${encodeURIComponent(returnTo)}` : '';
    const url = `${API_BASE}/auth/google${params}`;

    if (redirect && typeof window !== 'undefined') {
      window.location.href = url;
    }

    return url;
  }

  static async updateProfile(
    data: UpdateProfileRequest
  ): Promise<UserUpdate | null> {
    const controller = this.getAbortController('updateProfile');

    try {
      const response = await fetch(`${API_BASE}/api/v1/users/me`, {
        method: 'PUT',
        credentials: 'include',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(data),
        signal: controller.signal,
      });

      if (!response.ok) {
        throw new Error('Failed to update profile');
      }

      const result = await response.json();
      this.abortControllers.delete('updateProfile');
      return result;
    } catch (error) {
      if (error instanceof Error && error.name === 'AbortError') {
        console.log('Update profile request was cancelled');
        return null;
      }
      console.error('Profile update error:', error);
      throw error;
    }
  }
}
