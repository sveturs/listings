import { z } from 'zod';

// Схема валидации для публичных переменных
export const publicEnvSchema = z.object({
  NEXT_PUBLIC_API_URL: z.string().url().min(1),
  NEXT_PUBLIC_MINIO_URL: z.string().url().min(1),
  NEXT_PUBLIC_IMAGE_HOSTS: z.string().optional(),
  NEXT_PUBLIC_IMAGE_PATH_PATTERN: z.string().optional(),
  NEXT_PUBLIC_WEBSOCKET_URL: z.string().url().optional().or(z.literal('')),
  NEXT_PUBLIC_ENABLE_PAYMENTS: z.string().optional(),
});

// Схема для серверных переменных
export const serverEnvSchema = z.object({
  INTERNAL_API_URL: z.string().url().optional(),
  NODE_ENV: z.enum(['development', 'production', 'test']).optional(),
});

// Типы из схем
export type PublicEnvVariables = z.infer<typeof publicEnvSchema>;
export type ServerEnvVariables = z.infer<typeof serverEnvSchema>;

// Существующий интерфейс ImageHost остается
export interface ImageHost {
  protocol: 'http' | 'https';
  hostname: string;
  port?: string;
  pathname: string;
}

// Расширенный интерфейс конфигурации
export interface Config {
  // API Configuration
  api: {
    url: string;
    internalUrl?: string; // Для SSR запросов
    websocketUrl?: string;
  };

  // MinIO/Storage Configuration
  storage: {
    minioUrl: string;
    imageHosts: ImageHost[];
    imagePathPattern: string;
  };

  // Environment
  env: {
    isProduction: boolean;
    isDevelopment: boolean;
    isServer: boolean;
  };

  // Feature flags
  features: {
    enableChat: boolean;
    enablePayments: boolean;
  };
}

// Дополнительный тип для ошибок валидации
export interface ConfigValidationError {
  field: string;
  message: string;
}

// Оставляем для обратной совместимости
export interface EnvVariables {
  NEXT_PUBLIC_API_URL?: string;
  NEXT_PUBLIC_MINIO_URL?: string;
  NEXT_PUBLIC_IMAGE_HOSTS?: string;
  NEXT_PUBLIC_IMAGE_PATH_PATTERN?: string;
  NODE_ENV?: string;
}
