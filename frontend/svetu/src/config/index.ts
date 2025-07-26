import { env } from 'next-runtime-env';
import {
  Config,
  ImageHost,
  publicEnvSchema,
  serverEnvSchema,
  ConfigValidationError,
} from './types';

class ConfigManager {
  private config: Config | null = null;
  private validationErrors: ConfigValidationError[] = [];

  constructor() {
    // Инициализация будет ленивой
  }

  /**
   * Получает значение переменной окружения с поддержкой runtime
   */
  private getEnvValue(key: string, defaultValue: string = ''): string {
    if (typeof window === 'undefined') {
      // Server-side: используем process.env
      return process.env[key] || defaultValue;
    }
    // Client-side: используем runtime env
    return env(key) || defaultValue;
  }

  /**
   * Валидирует конфигурацию и собирает ошибки
   */
  private validateConfig(publicEnv: any, serverEnv: any): boolean {
    this.validationErrors = [];

    // Валидация публичных переменных
    const publicResult = publicEnvSchema.safeParse(publicEnv);
    if (!publicResult.success) {
      publicResult.error.issues.forEach((issue) => {
        this.validationErrors.push({
          field: issue.path.join('.'),
          message: issue.message,
        });
      });
    }

    // Валидация серверных переменных (только на сервере)
    if (typeof window === 'undefined') {
      const serverResult = serverEnvSchema.safeParse(serverEnv);
      if (!serverResult.success) {
        serverResult.error.issues.forEach((issue) => {
          this.validationErrors.push({
            field: issue.path.join('.'),
            message: issue.message,
          });
        });
      }
    }

    return this.validationErrors.length === 0;
  }

  private loadConfig(): Config {
    const isServer = typeof window === 'undefined';

    // Собираем публичные переменные
    const publicEnv = {
      NEXT_PUBLIC_API_URL: this.getEnvValue(
        'NEXT_PUBLIC_API_URL',
        'http://localhost:3000'
      ),
      NEXT_PUBLIC_MINIO_URL: this.getEnvValue(
        'NEXT_PUBLIC_MINIO_URL',
        'http://localhost:9000'
      ),
      NEXT_PUBLIC_IMAGE_HOSTS: this.getEnvValue('NEXT_PUBLIC_IMAGE_HOSTS'),
      NEXT_PUBLIC_IMAGE_PATH_PATTERN: this.getEnvValue(
        'NEXT_PUBLIC_IMAGE_PATH_PATTERN'
      ),
      NEXT_PUBLIC_WEBSOCKET_URL: this.getEnvValue('NEXT_PUBLIC_WEBSOCKET_URL'),
      NEXT_PUBLIC_ENABLE_PAYMENTS: this.getEnvValue(
        'NEXT_PUBLIC_ENABLE_PAYMENTS'
      ),
    };

    // Собираем серверные переменные
    const serverEnv = {
      INTERNAL_API_URL: process.env.INTERNAL_API_URL,
      NODE_ENV: process.env.NODE_ENV,
    };

    // Валидация в production
    if (process.env.NODE_ENV === 'production') {
      const isValid = this.validateConfig(publicEnv, serverEnv);
      if (!isValid) {
        console.error('Configuration validation failed:');
        this.validationErrors.forEach((error) => {
          console.error(`  ${error.field}: ${error.message}`);
        });
        // В production можем выбросить ошибку или использовать дефолты
        if (isServer) {
          throw new Error('Invalid server configuration');
        }
      }
    }

    return {
      api: {
        url: publicEnv.NEXT_PUBLIC_API_URL,
        internalUrl: isServer ? serverEnv.INTERNAL_API_URL : undefined,
        websocketUrl: publicEnv.NEXT_PUBLIC_WEBSOCKET_URL,
      },
      storage: {
        minioUrl: publicEnv.NEXT_PUBLIC_MINIO_URL,
        imageHosts: this.parseImageHosts(publicEnv.NEXT_PUBLIC_IMAGE_HOSTS),
        imagePathPattern:
          publicEnv.NEXT_PUBLIC_IMAGE_PATH_PATTERN || '/listings/**',
      },
      env: {
        isProduction: serverEnv.NODE_ENV === 'production',
        isDevelopment: serverEnv.NODE_ENV === 'development',
        isServer,
      },
      features: {
        enableChat: !!publicEnv.NEXT_PUBLIC_WEBSOCKET_URL,
        enablePayments: publicEnv.NEXT_PUBLIC_ENABLE_PAYMENTS === 'true',
      },
      claudeApiKey: this.getEnvValue('NEXT_PUBLIC_CLAUDE_API_KEY'),
    };
  }

  private parseImageHosts(hostsString?: string): ImageHost[] {
    const defaultHosts =
      'http:localhost:9000,https:svetu.rs:443,http:localhost:3000';
    const hosts = hostsString || defaultHosts;

    // Добавляем Google domains для аватарок
    const googleHosts: ImageHost[] = [
      {
        protocol: 'https',
        hostname: 'lh3.googleusercontent.com',
        pathname: '/**',
      },
      {
        protocol: 'https',
        hostname: '*.googleusercontent.com',
        pathname: '/**',
      },
    ];

    const parsedHosts = hosts.split(',').flatMap((host) => {
      const [protocol, hostname, port] = host.split(':');
      const pathnames = ['/listings/**', '/chat-files/**'];

      return pathnames.map((path) => {
        const imageHost: ImageHost = {
          protocol: protocol as 'http' | 'https',
          hostname,
          pathname: path,
        };

        if (
          port &&
          !(protocol === 'http' && port === '80') &&
          !(protocol === 'https' && port === '443')
        ) {
          imageHost.port = port;
        }

        return imageHost;
      });
    });

    return [...parsedHosts, ...googleHosts];
  }

  public getConfig(): Config {
    if (!this.config) {
      this.config = this.loadConfig();
    }
    return this.config;
  }

  /**
   * Получает ошибки валидации (если есть)
   */
  public getValidationErrors(): ConfigValidationError[] {
    return this.validationErrors;
  }

  /**
   * Сбрасывает кеш конфигурации (полезно для тестов)
   */
  public resetConfig(): void {
    this.config = null;
    this.validationErrors = [];
  }

  // Метод для получения правильного API URL в зависимости от контекста
  public getApiUrl(options?: { internal?: boolean }): string {
    const config = this.getConfig();

    // Для серверных запросов используем внутренний URL если доступен
    if (config.env.isServer && options?.internal && config.api.internalUrl) {
      return config.api.internalUrl;
    }

    // В development на клиенте используем относительные пути
    // чтобы запросы шли через прокси Next.js и сохранялись cookies
    if (config.env.isDevelopment && !config.env.isServer) {
      // Возвращаем пустую строку чтобы использовались относительные пути
      // Например: /api/v1/auth/refresh вместо http://100.88.44.15:3000/api/v1/auth/refresh
      return '';
    }

    return config.api.url;
  }

  public getMinioUrl(): string {
    return this.getConfig().storage.minioUrl;
  }

  public getImageHosts(): ImageHost[] {
    return this.getConfig().storage.imageHosts;
  }

  public isProduction(): boolean {
    return this.getConfig().env.isProduction;
  }

  public isDevelopment(): boolean {
    return this.getConfig().env.isDevelopment;
  }

  public isFeatureEnabled(feature: keyof Config['features']): boolean {
    return this.getConfig().features[feature];
  }

  public getImageBaseUrl(): string {
    const config = this.getConfig();
    if (config.env.isProduction) {
      return 'https://svetu.rs';
    }
    return config.storage.minioUrl;
  }

  public buildImageUrl(path: string): string {
    if (path.startsWith('http')) {
      return path;
    }

    const normalizedPath = path.startsWith('/') ? path : `/${path}`;

    if (
      normalizedPath.startsWith('/listings/') ||
      normalizedPath.startsWith('/chat-files/')
    ) {
      return `${this.getImageBaseUrl()}${normalizedPath}`;
    }

    return `${this.getConfig().api.url}${normalizedPath}`;
  }
}

// Создаем единственный экземпляр
const configManager = new ConfigManager();

// Экспортируем объект config для обратной совместимости
export const config = configManager.getConfig();

// Экспортируем менеджер как default
export default configManager;
