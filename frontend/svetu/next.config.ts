import type { NextConfig } from 'next';
import createNextIntlPlugin from 'next-intl/plugin';
import { version } from './package.json';
// import configManager from './src/config';

const withNextIntl = createNextIntlPlugin();

const nextConfig: NextConfig = {
  // output: 'standalone' - отключено для dev сервера (нужно только для Docker)
  // Moved from experimental to top-level in Next.js 15
  serverExternalPackages: [],
  env: {
    NEXT_PUBLIC_APP_VERSION: version,
  },
  eslint: {
    // During production builds, do not run ESLint
    ignoreDuringBuilds: true,
  },
  // Remove console.log in production builds
  compiler: {
    removeConsole:
      process.env.NODE_ENV === 'production'
        ? {
            exclude: ['error', 'warn'], // Keep error and warn for debugging
          }
        : false,
  },
  images: {
    remotePatterns: [
      // Динамически создаем список хостов из переменных окружения
      ...(
        process.env.NEXT_PUBLIC_IMAGE_HOSTS ||
        'http:localhost:9000,https:svetu.rs:443,http:localhost:3000,http:100.88.44.15:9000,http:100.88.44.15:3000'
      )
        .split(',')
        .flatMap((host) => {
          const [protocol, hostname, port] = host.split(':');
          // Разрешаем все пути для s3.svetu.rs
          const pathnames =
            hostname === 's3.svetu.rs'
              ? ['/**'] // Все пути разрешены для нашего S3
              : [
                  '/listings/**',
                  '/development-listings/**',
                  '/chat-files/**',
                  '/uploads/**',
                ];

          return pathnames.map((path) => {
            const config: any = {
              protocol: protocol as 'http' | 'https',
              hostname,
              pathname: path,
            };

            if (
              port &&
              !(protocol === 'http' && port === '80') &&
              !(protocol === 'https' && port === '443')
            ) {
              config.port = port;
            }

            return config;
          });
        }),
      // Наш S3 домен - разрешаем все пути
      {
        protocol: 'https',
        hostname: 's3.svetu.rs',
        pathname: '/**',
      },
      // Google domains для аватарок
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
      // Unsplash images для тестовых данных
      {
        protocol: 'https',
        hostname: 'images.unsplash.com',
        pathname: '/**',
      },
      // Picsum для тестовых изображений
      {
        protocol: 'https',
        hostname: 'picsum.photos',
        pathname: '/**',
      },
      // UI Avatars для генерируемых аватаров
      {
        protocol: 'https',
        hostname: 'ui-avatars.com',
        pathname: '/**',
      },
      // DigitalVision - внешний storefront
      {
        protocol: 'https',
        hostname: 'digitalvision.rs',
        pathname: '/**',
      },
    ],
  },
  async rewrites() {
    // Rewrites нужны только для локальной разработки
    if (process.env.NODE_ENV === 'development') {
      const apiUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:3000';
      const minioUrl =
        process.env.NEXT_PUBLIC_MINIO_URL || 'http://localhost:9000';

      return [
        {
          source: '/listings/:path*',
          destination: `${minioUrl}/listings/:path*`,
        },
        {
          source: '/products/:path*',
          destination: `${minioUrl}/products/:path*`,
        },
        {
          source: '/chat-files/:path*',
          destination: `${minioUrl}/chat-files/:path*`,
        },
        {
          source: '/uploads/:path*',
          destination: `${apiUrl}/uploads/:path*`,
        },
        // Проксируем все API запросы на backend, КРОМЕ /api/v2 (BFF proxy)
        {
          source: '/api/:path((?!v2).*)*',
          destination: `${apiUrl}/api/:path*`,
        },
        // Проксируем auth запросы (не API) с учетом локали
        // Исключаем callback страницы которые обрабатываются на frontend
        {
          source: '/:locale/auth/:path((?!callback|oauth/google/callback).*)',
          destination: `${apiUrl}/auth/:path*`,
        },
        // Проксируем auth запросы (не API) без локали
        // Исключаем callback страницы
        {
          source: '/auth/:path((?!callback).*)',
          destination: `${apiUrl}/auth/:path*`,
        },
        // Проксируем WebSocket для чата
        {
          source: '/ws/:path*',
          destination: `${apiUrl}/ws/:path*`,
        },
        // Проксируем запросы к Nominatim для избежания CORS
        {
          source: '/geocode/search',
          destination: 'https://nominatim.openstreetmap.org/search',
        },
        {
          source: '/geocode/reverse',
          destination: 'https://nominatim.openstreetmap.org/reverse',
        },
      ];
    }
    return [];
  },
  async headers() {
    return [
      {
        source: '/chat-files/:path*',
        headers: [
          {
            key: 'Cache-Control',
            value: 'public, max-age=31536000, immutable',
          },
        ],
      },
      // Кэширование модулей переводов
      {
        source: '/_next/static/chunks/src_messages_:path*.js',
        headers: [
          {
            key: 'Cache-Control',
            value: 'public, max-age=31536000, immutable',
          },
        ],
      },
      // Кэширование JSON файлов переводов
      {
        source: '/messages/:locale/:module.json',
        headers: [
          {
            key: 'Cache-Control',
            value: 'public, max-age=604800, stale-while-revalidate=86400',
          },
        ],
      },
    ];
  },
};

export default withNextIntl(nextConfig);
