import axios from '../../../api/axios';

class ChatService {
    constructor(userId) {
        this.userId = userId;
        this.ws = null;
        this.messageHandlers = new Set();
        this.reconnectAttempts = 0;
        this.maxReconnectAttempts = 10; // Увеличено максимальное количество попыток
        this.reconnectTimer = null;
        this.isConnecting = false;
        this.isActive = true;
        this.chatListHandlers = new Set();
        this.lastPingTime = 0;
        this.pingInterval = null;
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
                
                // Устанавливаем интервал для отправки ping-сообщений
                this.startPingInterval();
            };

            this.ws.onmessage = (event) => {
                if (!this.isActive || !this.ws) return;

                try {
                    const message = JSON.parse(event.data);
                    
                    // Обработка pong-сообщений
                    if (message.type === 'pong') {
                        console.log('Получен pong от сервера');
                        this.lastPingTime = Date.now(); // Обновляем время последнего ping/pong
                        return;
                    }
                    
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

        // Остановим ping-интервал при переподключении
        this.stopPingInterval();

        // Увеличиваем задержку только до определенного предела и линейный рост вместо экспоненциального
        const delay = Math.min(2000 + (this.reconnectAttempts * 1000), 15000);
        console.log(`Попытка переподключения через ${delay}ms (попытка ${this.reconnectAttempts + 1})`);

        clearTimeout(this.reconnectTimer);
        this.reconnectTimer = setTimeout(() => {
            this.reconnectAttempts++;
            
            // Если WebSocket все еще существует, закрываем его перед новым подключением
            if (this.ws) {
                try {
                    this.ws.close();
                } catch (e) {
                    console.error('Ошибка при закрытии старого WebSocket:', e);
                }
                this.ws = null;
            }
            
            this.connect();
        }, delay);
    }

    startPingInterval() {
        // Очищаем предыдущий интервал, если он существует
        this.stopPingInterval();
        
        // Установка интервала для отправки ping каждые 30 секунд
        this.pingInterval = setInterval(() => {
            if (this.ws && this.ws.readyState === WebSocket.OPEN) {
                try {
                    console.log('Отправка ping');
                    this.ws.send(JSON.stringify({ type: 'ping', timestamp: Date.now() }));
                    this.lastPingTime = Date.now();
                    
                    // Проверяем, получен ли ответ на предыдущий ping в течение 15 секунд
                    setTimeout(() => {
                        // Если последний ping был давно и не получен ответ
                        const pingTimeout = 15000; // 15 секунд
                        if (Date.now() - this.lastPingTime > pingTimeout && this.ws?.readyState === WebSocket.OPEN) {
                            console.log('Не получен ответ на ping, переподключение...');
                            // Закрываем соединение и пытаемся переподключиться
                            try {
                                this.ws.close(1000);
                            } catch (err) {
                                console.error('Ошибка при закрытии неотвечающего WebSocket:', err);
                            }
                            this.ws = null;
                            this.handleReconnect();
                        }
                    }, 15000);
                    
                } catch (e) {
                    console.error('Ошибка при отправке ping:', e);
                    this.handleReconnect();
                }
            } else {
                // Если соединение закрыто, пытаемся переподключиться
                this.handleReconnect();
            }
        }, 30000); // 30 секунд
    }
    
    stopPingInterval() {
        if (this.pingInterval) {
            clearInterval(this.pingInterval);
            this.pingInterval = null;
        }
    }

    disconnect() {
        console.log('Отключение WebSocket...');
        this.isActive = false;
        clearTimeout(this.reconnectTimer);
        this.stopPingInterval();

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