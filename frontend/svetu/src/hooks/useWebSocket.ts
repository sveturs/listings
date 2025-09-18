import { useEffect, useRef, useState, useCallback } from 'react';

export type ConnectionState =
  | 'connecting'
  | 'connected'
  | 'disconnected'
  | 'error';

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
    maxReconnectAttempts = 5,
  } = options;

  const [lastMessage, setLastMessage] = useState<MessageEvent | null>(null);
  const [connectionState, setConnectionState] =
    useState<ConnectionState>('disconnected');

  const websocket = useRef<WebSocket | null>(null);
  const reconnectTimeoutId = useRef<NodeJS.Timeout | null>(null);
  const reconnectAttempts = useRef(0);
  const isManualClose = useRef(false);

  // Store URL in ref to avoid recreating connect function
  const urlRef = useRef(url);
  const optionsRef = useRef(options);

  const connect = useCallback(() => {
    const currentUrl = urlRef.current;
    if (!currentUrl) {
      console.warn('WebSocket URL is not provided');
      return;
    }

    // Prevent multiple simultaneous connections
    if (websocket.current && websocket.current.readyState === WebSocket.CONNECTING) {
      console.log('WebSocket is already connecting...');
      return;
    }

    // Close existing connection if any
    if (websocket.current &&
        (websocket.current.readyState === WebSocket.OPEN ||
         websocket.current.readyState === WebSocket.CONNECTING)) {
      websocket.current.close();
    }

    try {
      setConnectionState('connecting');
      websocket.current = new WebSocket(currentUrl);

      websocket.current.onopen = (_event) => {
        console.log('WebSocket connected:', currentUrl);
        setConnectionState('connected');
        reconnectAttempts.current = 0;
        optionsRef.current.onOpen?.();
      };

      websocket.current.onmessage = (event) => {
        setLastMessage(event);
        optionsRef.current.onMessage?.(event);
      };

      websocket.current.onclose = (event) => {
        console.log('WebSocket closed:', event.code, event.reason);
        setConnectionState('disconnected');
        optionsRef.current.onClose?.();

        // Попытка переподключения если это не ручное закрытие
        if (
          !isManualClose.current &&
          optionsRef.current.shouldReconnect?.(event) !== false &&
          reconnectAttempts.current < (optionsRef.current.maxReconnectAttempts || 5)
        ) {
          reconnectAttempts.current++;
          console.log(
            `Attempting to reconnect (${reconnectAttempts.current}/${optionsRef.current.maxReconnectAttempts || 5})...`
          );

          reconnectTimeoutId.current = setTimeout(() => {
            connect();
          }, optionsRef.current.reconnectInterval || 3000);
        } else if (reconnectAttempts.current >= (optionsRef.current.maxReconnectAttempts || 5)) {
          console.error('Max reconnection attempts reached');
          setConnectionState('error');
        }
      };

      websocket.current.onerror = (error) => {
        console.error('WebSocket error:', error);
        setConnectionState('error');
        optionsRef.current.onError?.(error);
      };
    } catch (error) {
      console.error('Failed to create WebSocket connection:', error);
      setConnectionState('error');
    }
  }, []);

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

  // Update refs when props change
  useEffect(() => {
    urlRef.current = url;
  }, [url]);

  useEffect(() => {
    optionsRef.current = options;
  }, [options]);

  // Initial connection
  useEffect(() => {
    if (!url) return;

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
  }, [connect]); // Only run on mount and when connect changes (which is stable)

  // Handle URL changes
  useEffect(() => {
    if (!url) return;

    // Only reconnect if URL actually changed and we had a previous connection
    if (websocket.current && urlRef.current !== url) {
      disconnect();
      // Small delay to ensure cleanup completes
      const timer = setTimeout(() => {
        isManualClose.current = false;
        connect();
      }, 100);

      return () => clearTimeout(timer);
    }
  }, [url, connect, disconnect]); // Dependencies are stable callbacks

  return {
    lastMessage,
    connectionState,
    sendMessage,
    connect,
    disconnect,
  };
}
