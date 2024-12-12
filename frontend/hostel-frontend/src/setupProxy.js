const { createProxyMiddleware } = require('http-proxy-middleware');

module.exports = function(app) {
  // Proxy для API запросов
  app.use(
    '/api',
    createProxyMiddleware({
      target: 'http://localhost:3000',
      changeOrigin: true,
    })
  );

  // Proxy для WebSocket соединений
  app.use(
    '/ws',
    createProxyMiddleware({
      target: 'http://localhost:3001',
      ws: true,
      changeOrigin: true,
    })
  );
};