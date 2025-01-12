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
    onMessage(handler) {
        this.messageHandlers.add(handler);
        // При добавлении нового обработчика можно обновить сообщения
        if (this.ws?.readyState === WebSocket.OPEN) {
            this.ws.send(JSON.stringify({
                type: 'get_messages',
                chat_id: this.currentChatId // Добавьте это свойство при выборе чата
            }));
        }
        return () => this.messageHandlers.delete(handler);
    }
    
    // Добавьте метод для установки текущего чата
    setCurrentChat(chatId) {
        this.currentChatId = chatId;
        if (this.ws?.readyState === WebSocket.OPEN) {
            this.ws.send(JSON.stringify({
                type: 'get_messages',
                chat_id: chatId
            }));
        }
    }
    async sendMessage(message) {
        if (!this.ws || this.ws.readyState !== WebSocket.OPEN) {
            this.connect();
            await new Promise(resolve => setTimeout(resolve, 100));
        }

        try {
            if (this.ws?.readyState === WebSocket.OPEN) {
                this.ws.send(JSON.stringify(message));
            } else {
                throw new Error('WebSocket не подключен');
            }
        } catch (error) {
            console.error('Ошибка при отправке через WebSocket:', error);
            await this.sendMessageHTTP(message);
        }
    }

    async sendMessageHTTP(message) {
        const response = await fetch('/api/v1/marketplace/chat/messages', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            credentials: 'include',
            body: JSON.stringify(message)
        });

        if (!response.ok) {
            throw new Error('Ошибка отправки сообщения');
        }

        return response.json();
    }

    onMessage(handler) {
        this.messageHandlers.add(handler);
        return () => this.messageHandlers.delete(handler);
    }
}

export default ChatService;