package netserver

import (
	"encoding/json"
	"fmt"
	"log"
	"myapp/src/config"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type CmdHandler func(args map[string]string)

type JsonMeta struct {
	Name string            `json:"name"`
	Args map[string]string `json:"args"`
}

type WsServer struct {
	port     string
	upgrader websocket.Upgrader
	conn     *websocket.Conn
	handlers map[string]CmdHandler

	onConnect func()
	onHeart   func()

	mu            sync.Mutex // 消息锁
	isClientAlive bool       // 判断客户端是否活着
}

func New_WsServer(port string) *WsServer {
	self := WsServer{
		port: port,
		// 配置upgrader, 允许跨域请求
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}

	self.handlers = make(map[string]CmdHandler)
	self.isClientAlive = false
	return &self
}

// 开始服务器
func (self *WsServer) Run() {
	// 启动服务器
	http.HandleFunc("/ws", self.HandleConnections)
	http.HandleFunc("/live2d/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// 处理预检请求
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// 2. 计算本地文件的绝对路径
		baseDir := filepath.Join(config.RootPath, "/live2d")
		relPath := strings.TrimPrefix(r.URL.Path, "/live2d/")
		filePath := filepath.Join(baseDir, relPath)

		// 3. 检查文件是否存在，防止恶意探测
		_, err := os.Stat(filePath)
		if os.IsNotExist(err) {
			http.NotFound(w, r)
			return
		}
		http.ServeFile(w, r, filePath)
	})

	err := http.ListenAndServe(":"+self.port, nil)
	if err != nil {
		log.Println("ListenAndServe error:", err)
	}
}

// 注册服务端的执行插件,服务器接收到相关功能后触发该功能
func (self *WsServer) RegisterHandler(name string, handler CmdHandler) {
	self.handlers[name] = handler
}

// 启动链接
func (self *WsServer) HandleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := self.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	self.conn = conn

	// 关闭conn,确保资源释放
	defer func() {
		self.mu.Lock()
		self.conn = nil
		self.isClientAlive = false
		self.mu.Unlock()
		conn.Close()
	}()

	// 执行一些外部的初始化函数
	if self.onConnect != nil {
		go self.onConnect()
	}

	// 客户端存活确认 ------------------------------------------------
	pongWait := 10 * time.Second
	pongDead := 15 * time.Second
	conn.SetReadDeadline(time.Now().Add(pongDead)) // 设置死亡时间, 超过此时间后认为客户端死了
	conn.SetPongHandler(func(string) error {
		self.mu.Lock()
		self.isClientAlive = true
		self.mu.Unlock()
		conn.SetReadDeadline(time.Now().Add(pongDead))
		if self.onHeart != nil {
			go self.onHeart()
		}
		return nil
	})
	go func(c *websocket.Conn) {
		ticker := time.NewTicker(pongWait) // 启动自动重复定时器
		defer ticker.Stop()

		for {
			<-ticker.C

			self.mu.Lock()
			if self.conn == nil {
				self.mu.Unlock()
				return
			}

			// 发送标准底层 Ping 消息
			err := c.WriteMessage(websocket.PingMessage, nil)
			self.mu.Unlock()

			if err != nil {
				return
			}
		}
	}(conn)

	// 消息接收 --------------------------------------------------------
	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				self.mu.Lock()
				self.isClientAlive = false
				self.mu.Unlock()
			}
			break
		}

		conn.SetReadDeadline(time.Now().Add(pongDead))
		self.mu.Lock()
		self.isClientAlive = true
		self.mu.Unlock()

		var meta JsonMeta
		err = json.Unmarshal(p, &meta)
		if err != nil {
			continue
		}

		if handler, exists := self.handlers[meta.Name]; exists {
			fmt.Println(meta.Name)
			go handler(meta.Args)
		}
	}
}

// 服务器启动时的回调函数
func (self *WsServer) OnConnect(fn func()) {
	self.onConnect = fn
}

// 给予外界示例服务器
func (self *WsServer) API_GetConn() *websocket.Conn {
	return self.conn
}

// 向客户端发送对应类型的消息
func (self *WsServer) SendCommand(name string, arguement map[string]any) {
	if !self.isClientAlive {
		return
	}
	cont := CommandToJson(Command{Name: name, Args: arguement})
	self.SendMessage(cont)
}

// 向客户端发送文本信息
func (self *WsServer) SendMessage(msg string) {
	self.mu.Lock()
	defer self.mu.Unlock()
	if self.conn == nil {
		log.Println("WebSocket connection not established, skipping message")
		return
	}
	err := self.conn.WriteMessage(websocket.TextMessage, []byte(msg))
	if err != nil {
		log.Printf("Send error:%s\n", err)
	}
}

// 启动文件url,需要跨域头
func corsFileServer(fs http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 添加跨域头
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		// 调用原本的 FileServer
		fs.ServeHTTP(w, r)
	})
}

// api: 获得client是否存活
func (self *WsServer) API_GetAlive() bool {
	return self.isClientAlive
}

// api: 每次心跳时调用的函数
func (self *WsServer) OnHeart(fn func()) {
	self.onHeart = fn
}
