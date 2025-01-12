// frontend/hostel-frontend/src/components/marketplace/chat/ChatService.js
import axios from '../../../api/axios'; 
class ChatService {
    constructor(userId) {
        this.userId = userId;
        this.ws = null;
        this.messageHandlers = new Set();
        this.reconnectAttempts = 0;
        this.maxReconnectAttempts = 5;
        this.backendUrl = process.env.REACT_APP_BACKEND_URL || '';
        this.isConnecting = false;
    }

    connect() {
        if (this.ws?.readyState === WebSocket.OPEN || this.isConnecting) return;

        this.isConnecting = true;

        let wsUrl;
        if (process.env.NODE_ENV === 'development') {
            wsUrl = 'ws://localhost:3000/ws/chat';
        } else {
            const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
            const host = window.location.host;
            wsUrl = `${protocol}//${host}/ws/chat`;
        }

        console.log('Connecting to WebSocket:', wsUrl);

        try {
            this.ws = new WebSocket(wsUrl);

            this.ws.onopen = () => {
                console.log('WebSocket соединение установлено');
                this.isConnecting = false;
                this.reconnectAttempts = 0;

                // Небольшая задержка перед отправкой авторизации
                setTimeout(() => {
                    if (this.ws?.readyState === WebSocket.OPEN) {
                        this.ws.send(JSON.stringify({
                            type: 'auth',
                            user_id: this.userId
                        }));
                    }
                }, 100);
            };

            this.ws.onmessage = (event) => {
                try {
                    const message = JSON.parse(event.data);
                    this.messageHandlers.forEach(handler => handler(message));
                } catch (error) {
                    console.error('Ошибка при обработке сообщения:', error);
                }
            };

            this.ws.onclose = () => {
                console.log('WebSocket соединение закрыто');
                this.isConnecting = false;
                if (this.reconnectAttempts < this.maxReconnectAttempts) {
                    const delay = Math.min(1000 * Math.pow(2, this.reconnectAttempts), 30000);
                    setTimeout(() => this.connect(), delay);
                    this.reconnectAttempts++;
                }
            };

            this.ws.onerror = (error) => {
                console.error('WebSocket ошибка:', error);
                this.isConnecting = false;
            };

        } catch (error) {
            console.error('Ошибка при создании WebSocket:', error);
            this.isConnecting = false;
        }
    }

    disconnect() {
        this.isConnecting = false;
        if (this.ws) {
            this.ws.close();
            this.ws = null;
        }
    }


    setCurrentChat(chatId) {
        this.currentChatId = chatId;
        if (this.ws?.readyState === WebSocket.OPEN) {
            this.ws.send(JSON.stringify({
                type: 'get_messages',
                chat_id: chatId
            }));
        }
    }
 
    async sendMessageHTTP(message) {
        try {
            const response = await axios.post('/api/v1/marketplace/chat/messages', message, {
                withCredentials: true
            });
            return response.data;
        } catch (error) {
            console.error('Error sending message via HTTP:', error);
            throw error;
        }
    }
    async sendMessage(message) {
        if (!this.ws || this.ws.readyState !== WebSocket.OPEN) {
            try {
                this.connect();
                // Ждем подключения
                await new Promise((resolve, reject) => {
                    const timeout = setTimeout(() => reject(new Error('Connection timeout')), 5000);
                    const checkConnection = setInterval(() => {
                        if (this.ws?.readyState === WebSocket.OPEN) {
                            clearTimeout(timeout);
                            clearInterval(checkConnection);
                            resolve();
                        }
                    }, 100);
                });
            } catch (error) {
                console.error('WebSocket connection failed:', error);
                return this.sendMessageHTTP(message);
            }
        }
    
        try {
            this.ws.send(JSON.stringify({
                type: 'message',
                ...message
            }));
        } catch (error) {
            console.error('WebSocket send failed:', error);
            return this.sendMessageHTTP(message);
        }
    }

    onMessage(handler) {
        this.messageHandlers.add(handler);
        return () => this.messageHandlers.delete(handler);
    }
}

export default ChatService;