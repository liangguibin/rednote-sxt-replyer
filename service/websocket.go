package service

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/liangguibin/rednote-sxt-replyer/constant"
	"github.com/liangguibin/rednote-sxt-replyer/model"
	"github.com/liangguibin/rednote-sxt-replyer/store"
	"github.com/liangguibin/rednote-sxt-replyer/util"
	"net/http"
	"time"
)

// CreateWebSocketHeaders 创建 websocket 请求头
func CreateWebSocketHeaders() http.Header {
	headers := http.Header{}
	headers.Set("accept-encoding", "gzip, deflate, br, zstd")
	headers.Set("accept-language", "zh-CN,zh;q=0.9")
	headers.Set("cache-control", "no-cache")
	headers.Set("host", "zelda.xiaohongshu.com")
	headers.Set("origin", "https://sxt.xiaohongshu.com")
	headers.Set("pragma", "no-cache")
	headers.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Safari/537.36")
	return headers
}

// InitWebSocket 初始化 WebSocket - 包含重试机制
func InitWebSocket() {
	// 拨号器
	dialer := &websocket.Dialer{
		HandshakeTimeout: 45 * time.Second,
		ReadBufferSize:   1024,
		WriteBufferSize:  1024,
	}
	// 连接
	for {
		conn, _, err := dialer.Dial(constant.WebsocketUrl, CreateWebSocketHeaders())
		if err != nil {
			fmt.Println(util.GetTime(), " 连接失败，10秒后重试: ", err)
			time.Sleep(10 * time.Second)
			continue
		}

		store.WsConn = conn
		fmt.Println(util.GetTime(), " 已连接到 WebSocket 服务器")
		return
	}
}

// CloseWebSocket 关闭 WebSocket
func CloseWebSocket() {
	if store.WsConn != nil {
		err := store.WsConn.Close()
		if err != nil {
			fmt.Println(util.GetTime(), " WebSocket 连接关闭失败: ", err)
		} else {
			fmt.Println(util.GetTime(), " WebSocket 连接已关闭")
		}
	}
}

// SendMessage 发送消息
func SendMessage(message *model.WebSocketMessage) {
	// 检查连接状态
	if store.WsConn == nil {
		fmt.Println(util.GetTime(), " 连接不存在，无法发送消息")
		return
	}
	messageJson, _ := json.Marshal(message)
	// 加锁 避免并发写入
	store.WsWriteMutex.Lock()
	defer store.WsWriteMutex.Unlock()

	err := store.WsConn.WriteMessage(websocket.TextMessage, messageJson)
	if err != nil {
		fmt.Println(util.GetTime(), " 消息发送失败: ", err)
	}
}

// Authentication 发送认证消息
func Authentication() {
	message := model.WebSocketMessage{
		AppId: constant.AppId,
		Token: constant.Token,
		Type:  constant.One,
	}
	SendMessage(&message)
}

// StartHeartbeat 启动心跳机制 - 客户端
func StartHeartbeat() *time.Ticker {
	ticker := time.NewTicker(30 * time.Second)

	go func() {
		for {
			select {
			case <-ticker.C:
				// 检查连接状态
				if store.WsConn == nil {
					fmt.Println(util.GetTime(), " 连接不存在，跳过本次心跳")
					continue
				}
				// 心跳消息
				message := model.WebSocketMessage{
					Type: constant.Four,
				}
				SendMessage(&message)
			}
		}
	}()

	return ticker
}

// OnMessage 启用协程监听消息
func OnMessage() {
	go func() {
		seq := 1
		for {
			// 检查连接状态
			if store.WsConn == nil {
				fmt.Println(util.GetTime(), " 连接不存在，尝试重连...")
				InitWebSocket()
				seq = 1
				Authentication()
				continue
			}
			// 读取消息 消息为空则 ReadMessage() 阻塞 直至收到消息
			_, msg, err := store.WsConn.ReadMessage()
			if err != nil {
				fmt.Println(util.GetTime(), " 读取消息时出错，尝试重连... ", err)
				// 清理连接
				if store.WsConn != nil {
					CloseWebSocket()
					store.WsConn = nil
				}
				InitWebSocket()
				seq = 1
				Authentication()
				continue
			}
			seq += 1
			// 打印接收到的消息
			store.Logger.Info("收到消息: ", string(msg))
			// 解析消息
			var formatMsg model.WebSocketMessage
			_ = json.Unmarshal(msg, &formatMsg)
			// 处理消息
			// 主题订阅
			if formatMsg.Type == 129 {
				topic := store.BUserId
				topicEncoded, _ := util.AESECBEncrypt(topic, formatMsg.SecureKey)

				message := model.WebSocketMessage{
					Encrypt: true,
					Seq:     constant.Two,
					Topic:   topicEncoded,
					Type:    constant.Ten,
				}
				SendMessage(&message)
			}
			// 客户端信息
			if formatMsg.Type == 138 {
				additionalInfo := model.WebSocketAdditionalInfo{
					SellerId: store.CsProviderId,
					UserId:   store.BUserId,
				}
				userAgent := model.WebSocketUserAgent{
					AppName:    "walle-ad",
					AppVersion: "0.37.2",
				}
				info := model.WebSocketInfo{
					AdditionalInfo: additionalInfo,
					UserAgent:      userAgent,
				}
				message := model.WebSocketMessage{
					Info: info,
					Seq:  constant.Three,
					Type: constant.Twelve,
				}
				SendMessage(&message)
			}
			// 心跳机制 - 服务端
			if formatMsg.Type == constant.Four {
				message := model.WebSocketMessage{
					Type: 132,
				}
				SendMessage(&message)
			}
			// 接收消息
			if formatMsg.Type == constant.Two {
				message := model.WebSocketMessage{
					Ack:  formatMsg.Seq,
					Seq:  seq,
					Type: 130,
				}
				SendMessage(&message)
				// 处理消息
				HandleMessage(msg)
			}
		}
	}()
}
