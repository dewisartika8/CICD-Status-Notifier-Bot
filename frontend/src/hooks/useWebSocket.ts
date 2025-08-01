import { useEffect, useCallback } from 'react';
import { useAppDispatch, useAppSelector } from '@/store';
import { 
  setWebsocketConnected, 
  setLastWebsocketMessage 
} from '@/store/slices/uiSlice';
import { 
  updateProjectBuildEvent,
  addProject,
} from '@/store/slices/projectSlice';
import { WebSocketMessage } from '@/types';

export interface UseWebSocketReturn {
  isConnected: boolean;
  connect: () => void;
  disconnect: () => void;
  sendMessage: (message: any) => void;
  connectionStatus: string;
}

// WebSocket Service Class
class WebSocketService {
  private ws: WebSocket | null = null;
  private reconnectAttempts = 0;
  private maxReconnectAttempts = 5;
  private reconnectInterval = 3000;
  private url: string;
  private isManualClose = false;

  constructor(url: string) {
    this.url = url;
  }

  connect(
    onOpen?: () => void,
    onMessage?: (data: any) => void,
    onError?: (error: Event) => void,
    onClose?: () => void
  ) {
    try {
      this.ws = new WebSocket(this.url);
      
      this.ws.onopen = () => {
        console.log('WebSocket connected');
        this.reconnectAttempts = 0;
        onOpen?.();
      };

      this.ws.onmessage = (event) => {
        try {
          const data = JSON.parse(event.data);
          onMessage?.(data);
        } catch (error) {
          console.error('Failed to parse WebSocket message:', error);
        }
      };

      this.ws.onerror = (error) => {
        console.error('WebSocket error:', error);
        onError?.(error);
      };

      this.ws.onclose = () => {
        console.log('WebSocket closed');
        onClose?.();
        
        if (!this.isManualClose && this.reconnectAttempts < this.maxReconnectAttempts) {
          this.scheduleReconnect(onOpen, onMessage, onError, onClose);
        }
      };
    } catch (error) {
      console.error('Failed to create WebSocket connection:', error);
    }
  }

  private scheduleReconnect(
    onOpen?: () => void,
    onMessage?: (data: any) => void,
    onError?: (error: Event) => void,
    onClose?: () => void
  ) {
    this.reconnectAttempts++;
    console.log(`Attempting to reconnect... (${this.reconnectAttempts}/${this.maxReconnectAttempts})`);
    
    setTimeout(() => {
      this.connect(onOpen, onMessage, onError, onClose);
    }, this.reconnectInterval);
  }

  send(data: any) {
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify(data));
    } else {
      console.warn('WebSocket is not connected');
    }
  }

  close() {
    this.isManualClose = true;
    if (this.ws) {
      this.ws.close();
      this.ws = null;
    }
  }

  getConnectionState() {
    return this.ws?.readyState || WebSocket.CLOSED;
  }

  getConnectionStatus() {
    switch (this.getConnectionState()) {
      case WebSocket.CONNECTING:
        return 'connecting';
      case WebSocket.OPEN:
        return 'connected';
      case WebSocket.CLOSING:
        return 'closing';
      case WebSocket.CLOSED:
      default:
        return 'disconnected';
    }
  }
}

// WebSocket service instance
const WS_URL = import.meta.env.VITE_WS_URL || 'ws://localhost:8080/ws';
const webSocketService = new WebSocketService(WS_URL);

const useWebSocket = (): UseWebSocketReturn => {
  const dispatch = useAppDispatch();
  const { websocketConnected } = useAppSelector(state => state.ui);

  const handleWebSocketMessage = useCallback((data: WebSocketMessage) => {
    console.log('WebSocket message received:', data);
    
    dispatch(setLastWebsocketMessage(JSON.stringify(data)));
    
    switch (data.type) {
      case 'project_update':
        if (data.payload.project) {
          dispatch(addProject(data.payload.project));
        }
        break;
        
      case 'build_event':
        if (data.payload.projectId && data.payload.buildEvent) {
          dispatch(updateProjectBuildEvent({
            projectId: data.payload.projectId,
            buildEvent: data.payload.buildEvent
          }));
        }
        break;
        
      case 'notification':
        // Handle notification events
        console.log('Notification received:', data.payload);
        break;
        
      default:
        console.log('Unknown WebSocket message type:', data.type);
    }
  }, [dispatch]);

  const handleWebSocketOpen = useCallback(() => {
    dispatch(setWebsocketConnected(true));
  }, [dispatch]);

  const handleWebSocketClose = useCallback(() => {
    dispatch(setWebsocketConnected(false));
  }, [dispatch]);

  const handleWebSocketError = useCallback((error: Event) => {
    console.error('WebSocket error:', error);
    dispatch(setWebsocketConnected(false));
  }, [dispatch]);

  const connect = useCallback(() => {
    webSocketService.connect(
      handleWebSocketOpen,
      handleWebSocketMessage,
      handleWebSocketError,
      handleWebSocketClose
    );
  }, [handleWebSocketOpen, handleWebSocketMessage, handleWebSocketError, handleWebSocketClose]);

  const disconnect = useCallback(() => {
    webSocketService.close();
  }, []);

  const sendMessage = useCallback((message: any) => {
    webSocketService.send(message);
  }, []);

  useEffect(() => {
    connect();
    
    return () => {
      disconnect();
    };
  }, [connect, disconnect]);

  return {
    isConnected: websocketConnected,
    connect,
    disconnect,
    sendMessage,
    connectionStatus: webSocketService.getConnectionStatus(),
  };
};

export { useWebSocket };
export default useWebSocket;