import { WsClient } from "./WsClient";

export enum ConnectionState {
  Disconnected = 0,
  Connecting = 1,
  Connected = 2,
  Reconnecting = 3,
}

export class Chat {
  state: ConnectionState;
  client: WsClient;

  constructor() {
    this.state = ConnectionState.Disconnected;
    this.client = new WsClient();
  }

  async connect(url: string, token: string) {
    try {
      if (this.state === ConnectionState.Connected) {
        console.warn(`already connected to chat!`);
        return Promise.resolve();
      }

      await this.client.connect(url, token);
    } catch (error) {
      console.error("ERROR AT CONNECTION FROM CHAT", { error });
    }
  }
}
