window.ENV = {
    REACT_APP_BACKEND_URL: process.env.NODE_ENV === 'development'
        ? 'http://localhost:3000'
        : '',
    REACT_APP_API_PREFIX: '/api',
    REACT_APP_AUTH_PREFIX: '/auth',
    REACT_APP_WEBSOCKET_URL: process.env.NODE_ENV === 'development'
        ? 'ws://localhost:3000/ws'
        : 'wss://svetu.rs/ws',
    REACT_APP_WS_URL: process.env.NODE_ENV === 'development'
        ? 'ws://localhost:3000/ws/chat'
        : 'wss://svetu.rs/ws/chat',
    REACT_APP_MINIO_URL: process.env.NODE_ENV === 'development'
        ? 'http://localhost:9000'
        : ''
};