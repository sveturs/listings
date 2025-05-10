import { createProxyMiddleware } from 'http-proxy-middleware';
import { Express, Request } from 'express';
import { IncomingMessage, ServerResponse } from 'http';
import { ClientRequest } from 'http';

// Определяем собственный интерфейс для параметров прокси, так как
// типы в библиотеке могут не полностью соответствовать реальным возможностям
interface CustomProxyOptions {
  target: string;
  changeOrigin?: boolean;
  logLevel?: 'debug' | 'info' | 'warn' | 'error' | 'silent';
  ws?: boolean;
  pathRewrite?: { [path: string]: string };
  onProxyReq?: (proxyReq: ClientRequest, req: IncomingMessage) => void;
  onError?: (err: Error, req: IncomingMessage, res: ServerResponse) => void;
  [key: string]: any;
}

/**
 * Упрощенный setupProxy для локальной разработки - направляет все
 * /api, /ws и /auth запросы на backend
 */
module.exports = function(app: Express): void {
  // Используем простое прямое проксирование на бэкенд
  console.log('Setting up proxy for development environment');
  
  // Определяем бэкенд URL
  // Порядок: сначала пробуем localhost, потом контейнер
  const BACKEND_URL = 'http://localhost:3000';
  
  // Общие настройки прокси
  const proxyOptions: CustomProxyOptions = {
    target: BACKEND_URL,
    changeOrigin: true,
    logLevel: 'debug',
    onProxyReq: (proxyReq: ClientRequest, req: IncomingMessage): void => {
      console.log(`[PROXY] ${(req as Request).method} ${req.url} -> ${BACKEND_URL}${(proxyReq as any).path}`);
    },
    onError: (err: Error, req: IncomingMessage, res: ServerResponse): void => {
      console.error('[PROXY ERROR]', err);
      res.writeHead(500, { 'Content-Type': 'text/plain' });
      res.end(`Proxy Error: ${err.message}`);
    }
  };
  
  // API эндпоинты
  app.use('/api', createProxyMiddleware({
    ...proxyOptions,
    pathRewrite: { '^/api': '/api' }
  }));
  
  // WebSocket соединения
  app.use('/ws', createProxyMiddleware({
    ...proxyOptions,
    ws: true,
    pathRewrite: { '^/ws': '/ws' }
  }));
  
  // Auth эндпоинты - КРИТИЧНО для Google OAuth!
  app.use('/auth', createProxyMiddleware({
    ...proxyOptions,
    pathRewrite: { '^/auth': '/auth' }
  }));
  
  console.log('Proxy setup complete!');
};