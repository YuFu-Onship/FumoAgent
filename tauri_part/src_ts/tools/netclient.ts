export type OnOpen = () => void;
export class NetClient {
  private WS_URL = "ws://localhost:9527/ws";
  private socket: WebSocket | null = null;
  private connectCount: number = 0;
  private maxConnectCount: number = 10;
  // private reConnectLimitTime: number = 10000;
  private isConnect: boolean = false;
  public methons: any;
  public onOpen = () => {};
  public onClose = () => {};
  public onError = () => {};

  constructor(methons: any) {
    this.methons = methons;
  }

  // 启动客户端
  public async StartClient() {
    this.connect();
    // this.startHeartbeat();
  }

  // 链接服务器
  private connect() {
    console.log("正在尝试连接...");
    this.socket = new WebSocket(this.WS_URL);
    this.isConnect = false;

    this.socket.onopen = () => {
      this.connectCount = 0;
      console.log("连接成功");
      this.isConnect = true;
      if (this.onOpen) {
        this.onOpen();
      }
    };

    // 接受服务器消息
    this.receive_server(this.socket, this.methons);

    this.socket.onerror = (err) => {
      console.log("连接错误", err);
      this.isConnect = false;
      this.onError();
    };

    // 服务端关闭情况
    this.socket.onclose = () => {
      console.log("连接关闭");
      this.isConnect = false;
      this.onClose();
      this.handleReconnect();
    };
  }

  // 接收并处理服务器消息
  private receive_server(socket: WebSocket, methons: any) {
    socket.onmessage = (event) => {
      if (event.data === "pong") return;
      console.log("收到消息:", event.data);
      try {
        methons(JSON.parse(event.data));
      } catch (err) {
        console.log("解析消息错误", err, "\n消息为:", event.data);
      }
    };
  }

  // 重连函数
  private handleReconnect() {
    if (this.connectCount >= this.maxConnectCount) {
      console.log("重连失败");
      return;
    }

    this.connectCount++;
    const delay = Math.min(1000 * this.connectCount, 10000);
    console.log(`第 ${this.connectCount} 次重连，${delay}ms 后`);
    setTimeout(() => {
      this.connect();
    }, delay);
  }

  // 心跳检测
  // private startHeartbeat() {
  //   setInterval(() => {
  //     // 避免在连接未建立时发送。
  //     if (this.socket?.readyState === WebSocket.OPEN) {
  //       this.socket.send("ping");
  //     }
  //   }, 3000);
  // }

  // 发送文本消息
  public SendMessage(message: string) {
    if (!this.socket) return;
    if (this.socket.readyState === WebSocket.OPEN) {
      this.socket.send(message);
      console.log("发送消息:", message);
    } else {
      console.log("无法发送消息，连接未打开");
    }
  }

  // 发送json消息
  public SendJson(data: any) {
    try {
      const jsonString = JSON.stringify(data);
      this.SendMessage(jsonString);
    } catch (err) {
      console.error("JSON 序列化失败，无法发送:", err);
    }
  }

  public IsConnect(): boolean {
    return this.isConnect;
  }
  public reConnect() {
    this.connectCount = 0;
    this.connect();
  }
}
