const { createProxyMiddleware } = require('http-proxy-middleware');

module.exports = function(app) {
    app.use(
        '/api',
        createProxyMiddleware({
            target: process.env.NODE_ENV === 'development'
            ? 'http://backend:3000'  // Используем имя сервиса из docker-compose
            : 'https://svetu.rs',
            changeOrigin: true,
            onProxyReq: function(proxyReq, req, res) {
                console.log('Proxying to:', proxyReq.path);
            },
            onError: function(err, req, res) {
                console.error('Proxy Error:', err);
                if (err.code === 'ECONNREFUSED') {
                    // Пробуем альтернативный URL если основной недоступен
                    const altProxy = createProxyMiddleware({
                        target: 'http://localhost:3000',
                        changeOrigin: true
                    });
                    return altProxy(req, res);
                }
                res.status(500).send('Proxy Error');
            }
        })
    );

    // Добавляем прокси для WebSocket соединений
    app.use(
        '/ws',
        createProxyMiddleware({
            target: process.env.NODE_ENV === 'development'
                ? 'http://backend:3000'
                : 'https://svetu.rs',
            ws: true,
            changeOrigin: true
        })
    );
};