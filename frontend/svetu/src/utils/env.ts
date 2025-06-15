import { env as runtimeEnv } from 'next-runtime-env';

/**
 * Безопасный доступ к runtime переменным окружения
 * На сервере использует process.env, на клиенте - runtime env
 */
export function getEnv(key: string, defaultValue?: string): string | undefined {
  if (typeof window === 'undefined') {
    // Server-side: используем process.env
    return process.env[key] || defaultValue;
  }

  // Client-side: используем runtime env
  return runtimeEnv(key) || defaultValue;
}

/**
 * Типизированный доступ к публичным переменным
 */
export const publicEnv = {
  get API_URL() {
    return getEnv('NEXT_PUBLIC_API_URL', 'http://localhost:3000');
  },
  get MINIO_URL() {
    return getEnv('NEXT_PUBLIC_MINIO_URL', 'http://localhost:9000');
  },
  get WEBSOCKET_URL() {
    return getEnv('NEXT_PUBLIC_WEBSOCKET_URL');
  },
  get IMAGE_HOSTS() {
    return getEnv('NEXT_PUBLIC_IMAGE_HOSTS');
  },
  get ENABLE_PAYMENTS() {
    return getEnv('NEXT_PUBLIC_ENABLE_PAYMENTS') === 'true';
  },
};
