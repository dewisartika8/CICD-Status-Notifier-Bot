import { io, Socket } from 'socket.io-client';
import { WebSocketMessage } from '@/types';

class WebSocketService {
  private socket: Socket | null = null;
  private url: string;
  private reconnectAttempts = 0;
  private maxReconnectAttempts = 5;
  private reconnectInterval = 5000;
  private listeners: { [event: string]: ((data: any) => void)[] } = {};

  constructor(url?: string) {
    this.url = url || import.meta.env.VITE_WS_URL || 'ws://localhost:8080';
  }

  connect(): Promise<void> {
    return new Promise((resolve, reject) => {
      try {
        this.socket = io(this.url, {
          transports: ['websocket'],
          autoConnect: false,
          reconnection: true,
          reconnectionAttempts: this.maxReconnectAttempts,
          reconnectionDelay: this.reconnectInterval,
        });

        this.socket.on('connect', () => {
          console.log('WebSocket connected');
          this.reconnectAttempts = 0;
          this.emit('connection', { connected: true });
          resolve();
        });

        this.socket.on('disconnect', (reason) => {
          console.log('WebSocket disconnected:', reason);
          this.emit('connection', { connected: false, reason });
        });

        this.socket.on('connect_error', (error) => {
          console.error('WebSocket connection error:', error);
          this.emit('connection', { connected: false, error: error.message });
          
          if (this.reconnectAttempts >= this.maxReconnectAttempts) {
            reject(new Error('Max reconnection attempts reached'));
          }
        });

        this.socket.on('reconnect', (attemptNumber) => {
          console.log('WebSocket reconnected after', attemptNumber, 'attempts');
          this.emit('connection', { connected: true, reconnected: true });
        });

        this.socket.on('reconnect_error', (error) => {
          console.error('WebSocket reconnection error:', error);
          this.reconnectAttempts++;
          this.emit('connection', { connected: false, error: error.message });
        });

        // Listen for custom events
        this.socket.on('project_update', (data: WebSocketMessage) => {
          this.emit('project_update', data);
        });

        this.socket.on('build_event', (data: WebSocketMessage) => {
          this.emit('build_event', data);
        });

        this.socket.on('notification', (data: WebSocketMessage) => {
          this.emit('notification', data);
        });

        this.socket.connect();
      } catch (error) {
        reject(error);
      }
    });
  }

  disconnect(): void {
    if (this.socket) {
      this.socket.disconnect();
      this.socket = null;
    }
    this.listeners = {};
  }

  on(event: string, callback: (data: any) => void): void {
    if (!this.listeners[event]) {
      this.listeners[event] = [];
    }
    this.listeners[event].push(callback);
  }

  off(event: string, callback?: (data: any) => void): void {
    if (!this.listeners[event]) return;

    if (callback) {
      this.listeners[event] = this.listeners[event].filter(cb => cb !== callback);
    } else {
      delete this.listeners[event];
    }
  }

  private emit(event: string, data: any): void {
    if (this.listeners[event]) {
      this.listeners[event].forEach(callback => callback(data));
    }
  }

  send(event: string, data: any): void {
    if (this.socket && this.socket.connected) {
      this.socket.emit(event, data);
    } else {
      console.warn('WebSocket is not connected. Cannot send message:', event, data);
    }
  }

  isConnected(): boolean {
    return this.socket ? this.socket.connected : false;
  }

  // Subscribe to project updates
  subscribeToProject(projectId: string): void {
    this.send('subscribe_project', { projectId });
  }

  // Unsubscribe from project updates
  unsubscribeFromProject(projectId: string): void {
    this.send('unsubscribe_project', { projectId });
  }

  // Subscribe to all projects
  subscribeToAllProjects(): void {
    this.send('subscribe_all_projects', {});
  }

  // Get connection status
  getConnectionStatus(): { connected: boolean; reconnectAttempts: number } {
    return {
      connected: this.isConnected(),
      reconnectAttempts: this.reconnectAttempts,
    };
  }
}

// Create singleton instance
const webSocketService = new WebSocketService();

export default webSocketService;

// Export types for convenience
export type { WebSocketMessage };