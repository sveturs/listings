import { useEffect, useRef, useState, useCallback } from 'react';

export type ConnectionState = 'connecting' | 'connected' | 'disconnected' | 'error';

interface UseWebSocketOptions {
  onOpen?: () => void;
  onClose?: () => void;
  onError?: (error: Event) => void;
  onMessage?: (event: MessageEvent) => void;
  shouldReconnect?: (closeEvent: CloseEvent) => boolean;
  reconnectInterval?: number;
  maxReconnectAttempts?: number;
}

interface UseWebSocketReturn {
  lastMessage: MessageEvent | null;
  connectionState: ConnectionState;
  sendMessage: (message: string) => void;
  connect: () => void;
  disconnect: () => void;
}

export function useWebSocket(
  url: string,
  options: UseWebSocketOptions = {}
): UseWebSocketReturn {
  const {
    onOpen,
    onClose,
    onError,
    onMessage,
    shouldReconnect = () => true,
    reconnectInterval = 3000,
    maxReconnectAttempts = 5
  } = options;

  const [lastMessage, setLastMessage] = useState<MessageEvent | null>(null);
  const [connectionState, setConnectionState] = useState<ConnectionState>('connecting');
  
  const websocket = useRef<WebSocket | null>(null);
  const reconnectTimeoutId = useRef<NodeJS.Timeout | null>(null);
  const reconnectAttempts = useRef(0);
  const isManualClose = useRef(false);

  const connect = useCallback(() => {
    if (!url) {
      console.warn('WebSocket URL is not provided');
      return;
    }

    try {
      setConnectionState('connecting');
      websocket.current = new WebSocket(url);

      websocket.current.onopen = (event) => {
        console.log('WebSocket connected:', url);
        setConnectionState('connected');
        reconnectAttempts.current = 0;
        onOpen?.();
      };

      websocket.current.onmessage = (event) => {
        setLastMessage(event);
        onMessage?.(event);
      };

      websocket.current.onclose = (event) => {
        console.log('WebSocket closed:', event.code, event.reason);
        setConnectionState('disconnected');
        onClose?.();

        // Попытка переподключения если это не ручное закрытие
        if (!isManualClose.current && shouldReconnect(event) && reconnectAttempts.current < maxReconnectAttempts) {
          reconnectAttempts.current++;
          console.log(`Attempting to reconnect (${reconnectAttempts.current}/${maxReconnectAttempts})...`);
          
          reconnectTimeoutId.current = setTimeout(() => {
            connect();
          }, reconnectInterval);
        } else if (reconnectAttempts.current >= maxReconnectAttempts) {
          console.error('Max reconnection attempts reached');
          setConnectionState('error');
        }
      };

      websocket.current.onerror = (error) => {
        console.error('WebSocket error:', error);
        setConnectionState('error');
        onError?.(error);
      };

    } catch (error) {
      console.error('Failed to create WebSocket connection:', error);
      setConnectionState('error');
    }
  }, [url, onOpen, onClose, onError, onMessage, shouldReconnect, reconnectInterval, maxReconnectAttempts]);

  const disconnect = useCallback(() => {
    isManualClose.current = true;
    
    if (reconnectTimeoutId.current) {
      clearTimeout(reconnectTimeoutId.current);
      reconnectTimeoutId.current = null;
    }
    
    if (websocket.current) {
      websocket.current.close();
      websocket.current = null;
    }
    
    setConnectionState('disconnected');
  }, []);

  const sendMessage = useCallback((message: string) => {
    if (websocket.current && websocket.current.readyState === WebSocket.OPEN) {
      websocket.current.send(message);
    } else {
      console.warn('WebSocket is not connected. Cannot send message:', message);
    }
  }, []);

  // Подключение при монтировании
  useEffect(() => {
    isManualClose.current = false;
    connect();

    return () => {
      isManualClose.current = true;
      
      if (reconnectTimeoutId.current) {
        clearTimeout(reconnectTimeoutId.current);
      }
      
      if (websocket.current) {
        websocket.current.close();
      }
    };
  }, [connect]);

  // Переподключение при изменении URL
  useEffect(() => {
    if (websocket.current && websocket.current.url !== url) {
      disconnect();
      setTimeout(() => {
        isManualClose.current = false;
        connect();
      }, 100);
    }
  }, [url, connect, disconnect]);

  return {
    lastMessage,
    connectionState,
    sendMessage,
    connect,
    disconnect
  };
}