export interface BuildStatusUpdate {
  build_id: string
  status: 'pending' | 'running' | 'completed' | 'failed'
  progress?: number
  message?: string
  started_at?: string
  completed_at?: string
}

export class BuildWebSocketService {
  private ws: WebSocket | null = null
  private listeners: Map<string, (update: BuildStatusUpdate) => void> = new Map()
  private reconnectAttempts = 0
  private maxReconnectAttempts = 5
  private reconnectDelay = 1000

  connect() {
    if (this.ws?.readyState === WebSocket.OPEN) {
      return
    }

    const wsUrl = import.meta.env.VITE_WS_URL || 'ws://localhost:8000/ws/builds'
    this.ws = new WebSocket(wsUrl)

    this.ws.onopen = () => {
      console.log('Build WebSocket connected')
      this.reconnectAttempts = 0
    }

    this.ws.onmessage = (event) => {
      try {
        const update: BuildStatusUpdate = JSON.parse(event.data)
        this.notifyListeners(update)
      } catch (error) {
        console.error('Failed to parse WebSocket message:', error)
      }
    }

    this.ws.onclose = () => {
      console.log('Build WebSocket disconnected')
      this.attemptReconnect()
    }

    this.ws.onerror = (error) => {
      console.error('Build WebSocket error:', error)
    }
  }

  disconnect() {
    if (this.ws) {
      this.ws.close()
      this.ws = null
    }
    this.listeners.clear()
  }

  subscribe(buildId: string, callback: (update: BuildStatusUpdate) => void) {
    this.listeners.set(buildId, callback)
    this.connect()
  }

  unsubscribe(buildId: string) {
    this.listeners.delete(buildId)
    if (this.listeners.size === 0) {
      this.disconnect()
    }
  }

  private notifyListeners(update: BuildStatusUpdate) {
    const listener = this.listeners.get(update.build_id)
    if (listener) {
      listener(update)
    }
  }

  private attemptReconnect() {
    if (this.reconnectAttempts >= this.maxReconnectAttempts) {
      console.error('Max reconnection attempts reached')
      return
    }

    if (this.listeners.size === 0) {
      return
    }

    this.reconnectAttempts++
    const delay = this.reconnectDelay * Math.pow(2, this.reconnectAttempts - 1)
    
    setTimeout(() => {
      console.log(`Attempting to reconnect (${this.reconnectAttempts}/${this.maxReconnectAttempts})`)
      this.connect()
    }, delay)
  }
}

export const buildWebSocket = new BuildWebSocketService()
