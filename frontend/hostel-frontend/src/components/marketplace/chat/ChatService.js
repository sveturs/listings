import axios from '../../../api/axios';
class ChatService {
    constructor(userId) {
        this.userId = userId;
        this.ws = null;
        this.messageHandlers = new Set();
        this.reconnectAttempts = 0;
        this.maxReconnectAttempts = 5;
        this.reconnectTimeout = null;
        this.isConnecting = false;
    }

    connect() {
        if (this.ws?.readyState === WebSocket.OPEN || this.isConnecting) {
            return;
        }
    
        this.isConnecting = true;
        clearTimeout(this.reconnectTimer);
    
        try {
            const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
            // Используем window.location.host для локальной разработки
            const host = process.env.NODE_ENV === 'development' ? 
                'localhost:3000' : // для разработки
                window.location.host; // для продакшена
    
            const wsUrl = `${protocol}//${host}/ws/chat`;
            console.log('Попытка подключения к WebSocket:', wsUrl);
    
            this.ws = new WebSocket(wsUrl);
    
            this.ws.onopen = () => {
                console.log('WebSocket соединение установлено');
                this.isConnecting = false;
                this.reconnectAttempts = 0;
    
                // После успешного подключения отправляем авторизационные данные
                this.ws.send(JSON.stringify({
                    type: 'auth',
                    user_id: this.userId
                }));
            };
    
            this.ws.onmessage = (event) => {
                try {
                    const message = JSON.parse(event.data);
                    this.messageHandlers.forEach(handler => handler(message));
                } catch (error) {
                    console.error('Ошибка обработки сообщения:', error);
                }
            };
    
            this.ws.onerror = (error) => {
                console.error('WebSocket ошибка:', error);
                this.isConnecting = false;
            };
    
            this.ws.onclose = (event) => {
                console.log('WebSocket соединение закрыто:', event.code, event.reason);
                this.isConnecting = false;
                if (event.code !== 1000 && event.code !== 1005) {
                    this.handleReconnect();
                }
            };
    
        } catch (error) {
            console.error('Ошибка при создании WebSocket:', error);
            this.isConnecting = false;
            this.handleReconnect();
        }
    }

    handleReconnect() {
        if (this.reconnectAttempts >= this.maxReconnectAttempts || this.isConnecting) {
            return;
        }

        clearTimeout(this.reconnectTimer);
        this.reconnectTimer = setTimeout(() => {
            this.reconnectAttempts++;
            this.connect();
        }, Math.min(1000 * Math.pow(2, this.reconnectAttempts), 10000));
    }

    disconnect() {
        clearTimeout(this.reconnectTimer);
        this.reconnectAttempts = this.maxReconnectAttempts; 
        this.isConnecting = false;

        if (this.ws) {
            try {
                this.ws.close();
            } catch (e) {
                console.error('Ошибка при закрытии WebSocket:', e);
            }
            this.ws = null;
        }
    }

    async sendMessage(message) {
        try {
            // Пробуем отправить через WebSocket
            if (this.ws?.readyState === WebSocket.OPEN) {
                this.ws.send(JSON.stringify({
                    type: 'message',
                    ...message
                }));
                return;
            }

            // Если WebSocket недоступен, отправляем через HTTP
            console.log('WebSocket недоступен, отправка через HTTP');
            const response = await axios.post('/api/v1/marketplace/chat/messages', message);

            // После успешной отправки через HTTP пробуем переподключить WebSocket
            if (this.reconnectAttempts < this.maxReconnectAttempts) {
                this.connect();
            }

            return response.data;

        } catch (error) {
            console.error('Ошибка отправки сообщения:', error);
            throw error;
        }
    }

    onMessage(handler) {
        this.messageHandlers.add(handler);
        return () => this.messageHandlers.delete(handler);
    }

    async getMessageHistory(chatId, listingId) {
        if (!chatId || !listingId) {
            console.error('Отсутствует chatId или listingId:', { chatId, listingId });
            return [];
        }
    
        try {
            // Добавляем retry логику
            let attempts = 3;
            while (attempts > 0) {
                try {
                    const response = await axios.get(`/api/v1/marketplace/chat/${listingId}/messages`);
                    
                    if (response.data?.data) {
                        return response.data.data.sort((a, b) =>
                            new Date(a.created_at) - new Date(b.created_at)
                        );
                    }
                    break;
                } catch (error) {
                    attempts--;
                    if (attempts === 0) throw error;
                    await new Promise(resolve => setTimeout(resolve, 1000));
                }
            }
            return [];
        } catch (error) {
            console.error('Ошибка получения истории сообщений:', error);
            return [];
        }
    }
    async markMessagesAsRead(messageIds) {
        try {
            await axios.put('/api/v1/marketplace/chat/messages/read', {
                message_ids: messageIds
            });
        } catch (error) {
            console.error('Ошибка отметки сообщений как прочитанных:', error);
        }
    }
}

export default ChatService;