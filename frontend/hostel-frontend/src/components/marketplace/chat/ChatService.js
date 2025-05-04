import axios from '../../../api/axios';

class ChatService {
    constructor(userId) {
        this.userId = userId;
        this.ws = null;
        this.messageHandlers = new Set();
        this.reconnectAttempts = 0;
        this.maxReconnectAttempts = 5;
        this.reconnectTimer = null;
        this.isConnecting = false;
        this.isActive = true;
        this.chatListHandlers = new Set();
    }

    connect() {
        // Проверяем, активен ли сервис и нет ли уже активного подключения
        if (!this.isActive || this.ws?.readyState === WebSocket.OPEN || this.ws?.readyState === WebSocket.CONNECTING || this.isConnecting) {
            return;
        }

        this.isConnecting = true;
        clearTimeout(this.reconnectTimer);

        try {
            // Используем window.ENV для WebSocket URL если доступно
            let wsUrl;
            if (window.ENV && window.ENV.REACT_APP_WS_URL) {
                wsUrl = window.ENV.REACT_APP_WS_URL;
            } else {
                const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
                const host = process.env.NODE_ENV === 'development' ? 'localhost:3000' : window.location.host;
                wsUrl = `${protocol}//${host}/ws/chat`;
            }

            console.log('Попытка подключения к WebSocket:', wsUrl);
            this.ws = new WebSocket(wsUrl);

            this.ws.onopen = () => {
                if (!this.isActive) {
                    this.ws?.close();
                    return;
                }

                console.log('WebSocket соединение установлено');
                this.isConnecting = false;
                this.reconnectAttempts = 0;

                this.ws.send(JSON.stringify({
                    type: 'auth',
                    user_id: this.userId
                }));
            };

            this.ws.onmessage = (event) => {
                if (!this.isActive || !this.ws) return;

                try {
                    const message = JSON.parse(event.data);
                    this.messageHandlers.forEach(handler => handler(message));
                } catch (error) {
                    console.error('Ошибка обработки сообщения:', error);
                }
            };

            this.ws.onerror = (error) => {
                if (!this.isActive) return;
                console.error('WebSocket ошибка:', error);
                this.isConnecting = false;

                // Переподключаемся только если соединение было разорвано
                if (this.ws?.readyState === WebSocket.CLOSED) {
                    this.handleReconnect();
                }
            };

            this.ws.onclose = (event) => {
                if (!this.isActive) return;
                console.log('WebSocket соединение закрыто:', event.code, event.reason);
                this.isConnecting = false;
                this.ws = null;

                // Переподключаемся только при определенных кодах закрытия
                if (event.code === 1006 || event.code === 1001 || event.code === 1011) {
                    this.handleReconnect();
                }
            };

        } catch (error) {
            console.error('Ошибка при создании WebSocket:', error);
            this.isConnecting = false;
            if (this.isActive) {
                this.handleReconnect();
            }
        }
    }

    handleReconnect() {
        if (!this.isActive || this.reconnectAttempts >= this.maxReconnectAttempts || this.isConnecting) {
            return;
        }

        const delay = Math.min(1000 * Math.pow(2, this.reconnectAttempts), 10000);
        console.log(`Попытка переподключения через ${delay}ms (попытка ${this.reconnectAttempts + 1})`);

        clearTimeout(this.reconnectTimer);
        this.reconnectTimer = setTimeout(() => {
            this.reconnectAttempts++;
            this.connect();
        }, delay);
    }

    disconnect() {
        console.log('Отключение WebSocket...');
        this.isActive = false;
        clearTimeout(this.reconnectTimer);

        if (this.ws) {
            try {
                // Отправляем "чистое" закрытие
                this.ws.close(1000, 'Нормальное закрытие');
            } catch (e) {
                console.error('Ошибка при закрытии WebSocket:', e);
            }
            this.ws = null;
        }

        this.messageHandlers.clear();
        this.reconnectAttempts = 0;
        this.isConnecting = false;
    }


    async sendMessage(message) {
        try {
            if (this.ws?.readyState === WebSocket.OPEN) {
                this.ws.send(JSON.stringify({
                    type: 'message',
                    ...message
                }));
                return;
            }

            console.log('WebSocket недоступен, отправка через HTTP');
            const response = await axios.post('/api/v1/marketplace/chat/messages', message);

            if (this.isActive && this.reconnectAttempts < this.maxReconnectAttempts) {
                this.connect();
            }

            return response.data;

        } catch (error) {
            console.error('Ошибка отправки сообщения:', error);
            throw error;
        }
    }
    onChatListUpdate(handler) {
        this.chatListHandlers.add(handler);
        return () => this.chatListHandlers.delete(handler);
    }
    async updateChatsList() {
        try {
            const response = await axios.get('/api/v1/marketplace/chat');
            const chats = response.data?.data || [];
            this.chatListHandlers.forEach(handler => handler(chats));
        } catch (error) {
            console.error('Error updating chats list:', error);
        }
    }
    onMessage(handler) {
        this.messageHandlers.add(handler);
        return () => {
            this.messageHandlers.delete(handler);
        };
    }

    async getMessageHistory(chatId, listingId) {
        if (!chatId || !listingId) {
            console.error('Missing required params:', { chatId, listingId });
            return [];
        }

        try {
            const response = await axios.get('/api/v1/marketplace/chat/messages', {
                params: {
                    chat_id: chatId,
                    listing_id: listingId
                }
            });

            if (response.data?.data) {
                return response.data.data.sort((a, b) =>
                    new Date(a.created_at) - new Date(b.created_at)
                );
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
    handleNewMessage = async (message) => {
        this.messageHandlers.forEach(handler => handler(message));
        await this.updateChatsList();
    }
}

export default ChatService;