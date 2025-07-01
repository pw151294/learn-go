package ws_handler

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"iflytek.com/weipan4/learn-go/logger/zap"
	"iflytek.com/weipan4/learn-go/net/host"
	"net/http"
	"sync"
	"time"
)

const (
	heartbeatInterval = 10 * time.Second
	retryLimit        = 3
)

type WebsocketHandler struct {
	wsConn      *websocket.Conn
	dialer      *websocket.Dialer
	pingOutChan chan struct{}
	wsUrl       string
	isOpen      bool
	onMessage   func([]byte)
	onError     func(error)
	sync.RWMutex
}

type WebSocketHandlerOptions func(handler *WebsocketHandler)

func NewWebSocketHandler(opts ...WebSocketHandlerOptions) *WebsocketHandler {
	wsHandler := &WebsocketHandler{
		pingOutChan: make(chan struct{}),
	}
	if len(opts) > 0 {
		for _, opt := range opts {
			opt(wsHandler)
		}
	}

	return wsHandler
}

func (wsh *WebsocketHandler) Initialize() error {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}
	wsh.dialer = &websocket.Dialer{
		TLSClientConfig: tlsConfig,
		Proxy:           http.ProxyFromEnvironment,
	}
	zap.GetLogger().Info("opening websocket connection", "target url", wsh.wsUrl)
	header := http.Header{}
	conn, resp, err := wsh.dialer.Dial(wsh.wsUrl, header)
	if err != nil {
		if resp != nil {
			zap.GetLogger().Warn("failed to dial websocket", "status", resp.Status, "message", err)
		} else {
			zap.GetLogger().Warn("failed to dial websocket", "message", err)
		}
		return err
	}

	wsh.wsConn = conn
	wsh.isOpen = true
	zap.GetLogger().Info("successfully opened websocket connection", "target url", wsh.wsUrl)
	return nil
}

func (wsh *WebsocketHandler) Close() error {
	if !wsh.isOpen {
		return nil
	}

	wsh.pingOutChan <- struct{}{}
	wsh.isOpen = false
	if wsh.wsConn == nil {
		return errors.New("no websocket connection")
	}
	zap.GetLogger().Info("closing websocket connection", "target url", wsh.wsUrl)
	if err := wsh.wsConn.Close(); err != nil {
		zap.GetLogger().Warn("failed to close websocket connection", "message", err)
		return err
	}

	zap.GetLogger().Info("successfully closed websocket connection", "target url", wsh.wsUrl)
	return nil
}

func (wsh *WebsocketHandler) SendMessage(messageType int, data []byte) error {
	if !wsh.isOpen {
		return errors.New("websocket connection is closes")
	}
	if len(data) == 0 {
		return errors.New("websocket data length is 0")
	}
	wsh.Lock()
	defer wsh.Unlock()
	return wsh.wsConn.WriteMessage(messageType, data)
}

func (wsh *WebsocketHandler) StartPing() {
	defer func() {
		if r := recover(); r != nil {
			zap.GetLogger().Error("recovered from panic in ping", "message", r)
		}
	}()

	zap.GetLogger().Info("websocket channel send ping message")
	wsh.wsConn.SetPongHandler(func(resp string) error { // 设置心跳响应处理器
		zap.GetLogger().Info("begin handling websocket heartbeat response", "resp", resp)
		time.Sleep(1 * time.Second) // 模拟处理心跳响应的流程
		zap.GetLogger().Info("handling websocket heartbeat response completed")
		return nil
	})

	// 连接成功即开始发送心跳信息
	message := struct {
		Endpoint  string `json:"endpoint"`
		Timestamp int64  `json:"timestamp"`
		Payload   []byte `json:"payload"`
	}{}
	message.Endpoint = host.HostInfo.GetIP()
	message.Timestamp = time.Now().UnixNano()
	message.Payload = []byte("hello")
	msgBytes, _ := json.Marshal(message)
	if err := wsh.SendMessage(websocket.PingMessage, msgBytes); err != nil {
		zap.GetLogger().Warn("send websocket ping message failed", "message", err)
	}
	// 开启定时任务
	ticker := time.NewTicker(heartbeatInterval)
	for {
		if !wsh.isOpen {
			return
		}
		select {
		case <-ticker.C:
			message.Timestamp = time.Now().UnixNano()
			msgBytes, _ = json.Marshal(message)
			if err := wsh.SendMessage(websocket.PingMessage, msgBytes); err != nil {
				zap.GetLogger().Warn("send websocket ping message failed", "message", err)
			}
		case <-wsh.pingOutChan:
			ticker.Stop()
		}
	}
}

func (wsh *WebsocketHandler) ReceiveMessageLoop() {
	// 开始读取数据
	retryCount := 0
	for {
		if !wsh.isOpen {
			zap.GetLogger().Info("ending message receive listening because websocket connection is closed")
			break
		}
		messageType, rawMessage, err := wsh.wsConn.ReadMessage()
		if err != nil {
			retryCount++
			if retryCount >= retryLimit {
				zap.GetLogger().Error(fmt.Sprintf("reach the retry limit: %d", retryLimit), "message", err)
				wsh.onError(err)
				break
			}
			zap.GetLogger().Warn("read message failed, retrying for next read", "retry times", retryCount)
		} else if messageType != websocket.TextMessage {
			zap.GetLogger().Error("invalid message type, only accept UTF-8 or binary encoded text",
				"message type", messageType)
		} else {
			retryCount = 0
			zap.GetLogger().Info("received message successfully!", "message", string(rawMessage))
			wsh.onMessage(rawMessage)
		}
	}

	// 关闭连接
	zap.GetLogger().Info("close websocket connection......")
	if err := wsh.Close(); err != nil {
		zap.GetLogger().Error("failed to close websocket connection", "message", err)
	}
}

func WithWsUrl(url string) WebSocketHandlerOptions {
	return func(handler *WebsocketHandler) {
		handler.wsUrl = url
	}
}

func WithOnMessage(f func([]byte)) WebSocketHandlerOptions {
	return func(handler *WebsocketHandler) {
		handler.onMessage = f
	}
}

func WithOnError(f func(error)) WebSocketHandlerOptions {
	return func(handler *WebsocketHandler) {
		handler.onError = f
	}
}

func WithChannelOpen() WebSocketHandlerOptions {
	return func(handler *WebsocketHandler) {
		handler.isOpen = true
	}
}

func WithChannelClose() WebSocketHandlerOptions {
	return func(handler *WebsocketHandler) {
		handler.isOpen = false
	}
}
