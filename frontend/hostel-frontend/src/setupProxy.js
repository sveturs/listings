const { createProxyMiddleware } = require('http-proxy-middleware');

/**
 * Упрощенный setupProxy для локальной разработки - направляет все
 * /api, /ws и /auth запросы на backend
 */
module.exports = function(app) {
  // Используем простое прямое проксирование на бэкенд
  console.log('Setting up proxy for development environment');
  
  // Определяем бэкенд URL
  // Порядок: сначала пробуем localhost, потом контейнер
  const BACKEND_URL = 'http://localhost:3000';
  
  // Общие настройки прокси
  const proxyOptions = {
    target: BACKEND_URL,
    changeOrigin: true,
    logLevel: 'debug',
    onProxyReq: (proxyReq, req) => {
      console.log(`[PROXY] ${req.method} ${req.url} -> ${BACKEND_URL}${proxyReq.path}`);
    },
    onError: (err, req, res) => {
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