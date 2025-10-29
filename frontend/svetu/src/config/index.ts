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
      NEXT_PUBLIC_MINIO_BUCKET: this.getEnvValue(
        'NEXT_PUBLIC_MINIO_BUCKET',
        'listings'
      ),
      NEXT_PUBLIC_IMAGE_HOSTS: this.getEnvValue('NEXT_PUBLIC_IMAGE_HOSTS'),
      NEXT_PUBLIC_IMAGE_PATH_PATTERN: this.getEnvValue(
        'NEXT_PUBLIC_IMAGE_PATH_PATTERN'
      ),
      NEXT_PUBLIC_WEBSOCKET_URL: this.getEnvValue('NEXT_PUBLIC_WEBSOCKET_URL'),
      NEXT_PUBLIC_ENABLE_PAYMENTS: this.getEnvValue(
        'NEXT_PUBLIC_ENABLE_PAYMENTS'
      ),
      NEXT_PUBLIC_MAPBOX_ACCESS_TOKEN: this.getEnvValue(
        'NEXT_PUBLIC_MAPBOX_ACCESS_TOKEN'
      ),
      NEXT_PUBLIC_ENABLE_MAPBOX: this.getEnvValue('NEXT_PUBLIC_ENABLE_MAPBOX'),
      NEXT_PUBLIC_CLAUDE_API_KEY: this.getEnvValue(
        'NEXT_PUBLIC_CLAUDE_API_KEY'
      ),
      NEXT_PUBLIC_FRONTEND_URL: this.getEnvValue('NEXT_PUBLIC_FRONTEND_URL'),
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
        minioBucket: publicEnv.NEXT_PUBLIC_MINIO_BUCKET || 'listings',
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
        enableMapbox: publicEnv.NEXT_PUBLIC_ENABLE_MAPBOX === 'true',
      },
      mapbox: {
        accessToken: publicEnv.NEXT_PUBLIC_MAPBOX_ACCESS_TOKEN,
        enabled: publicEnv.NEXT_PUBLIC_ENABLE_MAPBOX === 'true',
      },
      claudeApiKey: publicEnv.NEXT_PUBLIC_CLAUDE_API_KEY,
      frontendUrl: publicEnv.NEXT_PUBLIC_FRONTEND_URL,
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

  /**
   * @deprecated Use apiClient instead of direct fetch with getApiUrl()
   *
   * ВАЖНО: Этот метод устарел и будет удален!
   *
   * ❌ НЕ ИСПОЛЬЗУЙ:
   * ```typescript
   * const apiUrl = configManager.getApiUrl();
   * const response = await fetch(`${apiUrl}/api/v1/...`);
   * ```
   *
   * ✅ ИСПОЛЬЗУЙ apiClient:
   * ```typescript
   * import { apiClient } from '@/services/api-client';
   * const response = await apiClient.get('/...');
   * ```
   *
   * Или специализированные сервисы (deliveryService, etc.)
   *
   * Исключение: Server-side код (getApiUrl({ internal: true }))
   */
  public getApiUrl(options?: { internal?: boolean }): string {
    const config = this.getConfig();

    // Для серверных запросов используем внутренний URL если доступен
    if (config.env.isServer && options?.internal && config.api.internalUrl) {
      return config.api.internalUrl;
    }

    // Всегда используем URL из конфигурации - это проще и надежнее
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
    // Всегда используем MinIO URL для изображений
    // В production это будет https://s3.svetu.rs
    // В development это будет http://localhost:9000
    return config.storage.minioUrl;
  }

  public buildImageUrl(path: string): string {
    // Backend теперь возвращает полные URL в поле image_url/public_url
    // Этот метод остается только для обратной совместимости

    // Если путь уже является полным URL, возвращаем его как есть
    if (path.startsWith('http')) {
      return path;
    }

    // Обратная совместимость для старых данных (backend иногда возвращает относительные пути)
    const normalizedPath = path.startsWith('/') ? path : `/${path}`;
    const config = this.getConfig();
    const minioUrl = config.storage.minioUrl;
    const bucketName = config.storage.minioBucket;

    // Простое построение URL для обратной совместимости
    // Поддерживаем пути: /listings/..., /products/..., /chat-files/..., storefronts/...
    if (
      normalizedPath.match(
        /^\/(listings\/.*|products\/\d+\/.*|chat-files\/.*)$/
      )
    ) {
      // Проверяем, не начинается ли путь уже с имени bucket
      // Если путь /listings/... и bucket тоже listings, не дублируем
      if (normalizedPath.startsWith(`/${bucketName}/`)) {
        return `${minioUrl}${normalizedPath}`;
      }
      return `${minioUrl}/${bucketName}${normalizedPath}`;
    }

    // Для остальных путей - предполагаем что это API пути
    return `${config.api.url}${normalizedPath}`;
  }

  /**
   * Получает Mapbox access token
   */
  public getMapboxToken(): string | undefined {
    return this.getConfig().mapbox.accessToken;
  }

  /**
   * Проверяет, включен ли Mapbox
   */
  public isMapboxEnabled(): boolean {
    return this.getConfig().mapbox.enabled;
  }

  /**
   * Получает Claude API key
   */
  public getClaudeApiKey(): string | undefined {
    return this.getConfig().claudeApiKey;
  }

  /**
   * Получает Frontend URL
   */
  public getFrontendUrl(): string | undefined {
    return this.getConfig().frontendUrl;
  }

  /**
   * Получает WebSocket URL
   */
  public getWebSocketUrl(): string | undefined {
    return this.getConfig().api.websocketUrl;
  }
}

// Создаем единственный экземпляр
const configManager = new ConfigManager();

// Экспортируем объект config для обратной совместимости
export const config = configManager.getConfig();

// Экспортируем функцию buildImageUrl для удобства
export const buildImageUrl = (path: string) =>
  configManager.buildImageUrl(path);

// Экспортируем менеджер как default
export default configManager;
