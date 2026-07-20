import type { WebSocketMessage } from "./dispatcher";

type Listener = (message: WebSocketMessage) => void;

const API_BASE_URL =
  process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:8080";
const WS_URL = API_BASE_URL.replace(/^http/, "ws");

class WebSocketManager {
  private socket: WebSocket | null = null;
  private token: string | null = null;
  private listeners = new Set<Listener>();
  private reconnectAttempt = 0;
  private reconnectTimer: number | undefined;
  private manualClose = false;

  connect(accessToken: string) {
    console.log("Connect called");
    this.token = accessToken;
    this.manualClose = false;
    if (
      this.socket?.readyState === WebSocket.OPEN ||
      this.socket?.readyState === WebSocket.CONNECTING
    )
      return;

    this.socket = new WebSocket(
      `${WS_URL}/ws?token=${encodeURIComponent(accessToken)}`,
    );
    this.socket.onopen = () => {
      console.log("OPEN");
      this.reconnectAttempt = 0;
    };

    this.socket.onmessage = (event) => {
      console.log("RAW WS:", event.data);

      const message = JSON.parse(event.data) as WebSocketMessage;

      console.log("PARSED WS:", message);

      this.emit(message);
    };

    this.socket.onclose = (event) => {
      console.log(event, event.reason, event.code);
      this.socket = null;
      this.scheduleReconnect();
    };
    this.socket.onerror = () => this.socket?.close();
  }

  disconnect() {
    this.manualClose = true;
    window.clearTimeout(this.reconnectTimer);
    this.socket?.close(1000, "logout");
    this.socket = null;
  }

  send(type: string, payload: unknown) {
    if (this.socket?.readyState !== WebSocket.OPEN)
      throw new Error("Websocket is not connected");
    this.socket.send(JSON.stringify({ type, payload }));
  }

  subscribe(listener: Listener) {
    this.listeners.add(listener);
    return () => {
      this.listeners.delete(listener);
    };
  }

  private emit(message: WebSocketMessage) {
    for (const listener of this.listeners) listener(message);
  }

  private scheduleReconnect() {
    if (this.manualClose || !this.token) return;
    const delay = Math.min(30_000, 1000 * 2 ** this.reconnectAttempt++);
    this.reconnectTimer = window.setTimeout(
      () => this.token && this.connect(this.token),
      delay,
    );
  }
}

export const websocketManager = new WebSocketManager();
