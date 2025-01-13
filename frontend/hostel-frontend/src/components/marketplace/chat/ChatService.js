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

    async connect() {
        if (this.ws?.readyState === WebSocket.OPEN || this.isConnecting) {
            console.log('WebSocket уже подключен или подключается');
            return;
        }

        this.isConnecting = true;

        try {
            const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
            const host = process.env.NODE_ENV === 'development' ? 
                'localhost:3000' : 
                window.location.host;
            
            const wsUrl = `${protocol}//${host}/ws/chat`;
            console.log('Подключение к WebSocket:', wsUrl);

            this.ws = new WebSocket(wsUrl);

            this.ws.onopen = () => {
                console.log('WebSocket соединение установлено');
                this.isConnecting = false;
                this.reconnectAttempts = 0;

                this.ws.send(JSON.stringify({
                    type: 'auth',
                    user_id: this.userId
                }));
            };

            this.ws.onmessage = (event) => {
                try {
                    const message = JSON.parse(event.data);
                    if (message.error) {
                        console.error('Ошибка WebSocket:', message.error);
                        return;
                    }
                    this.messageHandlers.forEach(handler => handler(message));
                } catch (error) {
                    console.error('Ошибка обработки сообщения:', error);
                }
            };

            this.ws.onclose = (event) => {
                console.log('WebSocket соединение закрыто:', event.code, event.reason);
                this.isConnecting = false;
                this.handleReconnect();
            };

            this.ws.onerror = (error) => {
                console.error('WebSocket ошибка:', error);
                this.isConnecting = false;
            };

        } catch (error) {
            console.error('Ошибка при создании WebSocket:', error);
            this.isConnecting = false;
            this.handleReconnect();
        }
    }

    handleReconnect() {
        if (this.reconnectAttempts >= this.maxReconnectAttempts) {
            console.log('Достигнуто максимальное количество попыток переподключения');
            return;
        }

        const delay = Math.min(1000 * Math.pow(2, this.reconnectAttempts), 30000);
        console.log(`Попытка переподключения через ${delay}мс`);

        clearTimeout(this.reconnectTimeout);
        this.reconnectTimeout = setTimeout(() => {
            this.reconnectAttempts++;
            this.connect();
        }, delay);
    }

    disconnect() {
        clearTimeout(this.reconnectTimeout);
        this.isConnecting = false;
        this.reconnectAttempts = this.maxReconnectAttempts; // Предотвращаем автоматическое переподключение
        
        if (this.ws) {
            this.ws.close();
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

    async getMessageHistory(chatId) {
        try {
            const response = await axios.get(`/api/v1/marketplace/chat/${chatId}/messages`);
            return Array.isArray(response.data?.data) ? response.data.data : [];
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