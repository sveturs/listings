// Server-side API client for server actions
import configManager from '@/config';

const INTERNAL_API_URL =
  process.env.INTERNAL_API_URL ||
  process.env.NEXT_PUBLIC_API_URL ||
  configManager.getApiUrl({ internal: true });

interface ApiClientOptions {
  headers?: Record<string, string>;
}

export const apiClientServer = {
  async get(path: string, options?: ApiClientOptions) {
    const response = await fetch(`${INTERNAL_API_URL}${path}`, {
      headers: {
        'Content-Type': 'application/json',
        ...options?.headers,
      },
    });
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    return {
      data: await response.json(),
    };
  },

  async post(path: string, data: any, options?: ApiClientOptions) {
    const response = await fetch(`${INTERNAL_API_URL}${path}`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        ...options?.headers,
      },
      body: JSON.stringify(data),
    });
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    return {
      data: await response.json(),
    };
  },

  async put(path: string, data: any, options?: ApiClientOptions) {
    const response = await fetch(`${INTERNAL_API_URL}${path}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
        ...options?.headers,
      },
      body: JSON.stringify(data),
    });
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    return {
      data: await response.json(),
    };
  },

  async delete(path: string, options?: ApiClientOptions) {
    const response = await fetch(`${INTERNAL_API_URL}${path}`, {
      method: 'DELETE',
      headers: {
        'Content-Type': 'application/json',
        ...options?.headers,
      },
    });
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    return {
      data: await response.json(),
    };
  },
};
