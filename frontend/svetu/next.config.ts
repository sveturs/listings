import type { NextConfig } from 'next';
import createNextIntlPlugin from 'next-intl/plugin';
import configManager from './src/config';

const withNextIntl = createNextIntlPlugin();

const nextConfig: NextConfig = {
  images: {
    remotePatterns: configManager.getImageHosts(),
  },
  async rewrites() {
    // Rewrites нужны только для локальной разработки
    if (process.env.NODE_ENV === 'development') {
      return [
        {
          source: '/chat-files/:path*',
          destination: 'http://localhost:9000/chat-files/:path*',
        },
        // Проксируем API запросы на backend, кроме /api/auth/*
        {
          source: '/api/:path((?!auth).*)',
          destination: 'http://localhost:3000/api/:path*',
        },
        // Проксируем auth запросы (не API)
        {
          source: '/auth/:path*',
          destination: 'http://localhost:3000/auth/:path*',
        },
        // Проксируем WebSocket для чата
        {
          source: '/ws/:path*',
          destination: 'http://localhost:3000/ws/:path*',
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
    ];
  },
};

export default withNextIntl(nextConfig);
