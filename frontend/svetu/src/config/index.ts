import { Config, ImageHost, EnvVariables } from './types';

class ConfigManager {
  private config: Config;

  constructor() {
    this.config = this.loadConfig();
  }

  private loadConfig(): Config {
    const env = process.env as unknown as EnvVariables;
    const imagePathPattern =
      env.NEXT_PUBLIC_IMAGE_PATH_PATTERN || '/listings/**';

    return {
      api: {
        url: env.NEXT_PUBLIC_API_URL || 'http://localhost:3000',
      },
      storage: {
        minioUrl: env.NEXT_PUBLIC_MINIO_URL || 'http://localhost:9000',
        imageHosts: this.parseImageHosts(env.NEXT_PUBLIC_IMAGE_HOSTS),
        imagePathPattern,
      },
      env: {
        isProduction: env.NODE_ENV === 'production',
        isDevelopment: env.NODE_ENV === 'development',
      },
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
      // Создаем конфигурации для разных путей
      const pathnames = ['/listings/**', '/chat-files/**'];

      return pathnames.map((path) => {
        const imageHost: ImageHost = {
          protocol: protocol as 'http' | 'https',
          hostname,
          pathname: path,
        };

        // Добавляем порт только если он указан и не является стандартным
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

    // Объединяем списки хостов
    return [...parsedHosts, ...googleHosts];
  }

  public getConfig(): Config {
    return this.config;
  }

  // Удобные методы для доступа к часто используемым значениям
  public getApiUrl(): string {
    // В разработке используем пустую строку, чтобы запросы шли через proxy
    if (this.config.env.isDevelopment && typeof window !== 'undefined') {
      return '';
    }
    return this.config.api.url;
  }

  public getMinioUrl(): string {
    return this.config.storage.minioUrl;
  }

  public getImageHosts(): ImageHost[] {
    return this.config.storage.imageHosts;
  }

  public isProduction(): boolean {
    return this.config.env.isProduction;
  }

  public isDevelopment(): boolean {
    return this.config.env.isDevelopment;
  }

  // Метод для получения базового URL для изображений
  public getImageBaseUrl(): string {
    if (this.config.env.isProduction) {
      return 'https://svetu.rs';
    }
    return this.config.storage.minioUrl;
  }

  // Метод для построения полного URL изображения
  public buildImageUrl(path: string): string {
    // Если URL уже полный, возвращаем как есть
    if (path.startsWith('http')) {
      return path;
    }

    // Убеждаемся, что путь начинается со слеша
    const normalizedPath = path.startsWith('/') ? path : `/${path}`;

    // Если путь начинается с /listings/ или /chat-files/, используем MinIO
    if (
      normalizedPath.startsWith('/listings/') ||
      normalizedPath.startsWith('/chat-files/')
    ) {
      return `${this.getImageBaseUrl()}${normalizedPath}`;
    }

    // Иначе используем API URL (для обратной совместимости)
    return `${this.config.api.url}${normalizedPath}`;
  }
}

// Создаем единственный экземпляр конфигурации
const configManager = new ConfigManager();

// Экспортируем как объект config для удобства использования
export const config = configManager.getConfig();

// Экспортируем также сам менеджер для доступа к методам
export default configManager;
