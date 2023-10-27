export class WsClient {
  // 3 seconds
  websocketTimeout = 3000;

  ws?: WebSocket;

  constructor() {}

  async connect(url: string, token: string) {
    return new Promise<void>((resolve, reject) => {
      const wsTimeout = setTimeout(() => {
        //this.close() restart state
        reject(new Error("server ws connection has timed out"));
      }, this.websocketTimeout);

      const params = new URLSearchParams();

      params.set("access_token", token);

      this.ws = new WebSocket(`${url}?${params.toString()}`);
      this.ws.binaryType = "arraybuffer";

      this.ws.onopen = () => {
        clearTimeout(wsTimeout);
      };

      this.ws.onerror = (event) => {
        console.log("WEBSOCKET ERROR!", { event });
      };

      this.ws.onmessage = (message) => {
        console.log("WEBSOCKET MESSSAGE!", { message });
      };

      console.log("WS CONNECTED SUCCESSFULLY");

      resolve();
    });
  }
}
