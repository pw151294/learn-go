package main

import (
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"iflytek.com/weipan4/learn-go/logger/zap"
	"sync"
	"time"
)

const (
	readTimeout = 60 * time.Second
	retryLimit  = 3
)

type WebSocketClient struct {
	wskConn   *websocket.Conn
	isOpen    bool
	onMessage func([]byte) error
	onData    func([]byte) error
	onError   func(error)
	sync.RWMutex
}

type WebSocketClientOptions func(*WebSocketClient)

func (wsc *WebSocketClient) ReadLoop() {
	defer func() {
		if r := recover(); r != nil {
			zap.GetLogger().Error("receive message failed", "error", r)
		}
	}()

	for {
		if !wsc.isOpen {
			break
		}
		messageType, rawMessage, err := wsc.wskConn.ReadMessage()
		if err != nil {
			wsc.onError(err)
			break
		} else if messageType != websocket.TextMessage && messageType != websocket.BinaryMessage && messageType != websocket.CloseMessage {
			wsc.onError(fmt.Errorf("invalid message type: %v, only support text or binary msg", messageType))
		} else if messageType == websocket.BinaryMessage {
			if wsc.onData != nil {
				if err = wsc.wskConn.SetReadDeadline(time.Now().Add(readTimeout)); err != nil {
					wsc.onError(err)
					break
				}
				if err = wsc.onData(rawMessage); err != nil {
					wsc.onError(err)
					break
				}
			}
		} else if messageType == websocket.TextMessage {
			if err = wsc.onMessage(rawMessage); err != nil {
				wsc.onError(err)
				break
			}
		} else {
			zap.GetLogger().Info("receive close message, ready to close......")
			break
		}
	}

	// 关闭连接
	zap.GetLogger().Info("close websocket connection......")
	if err := wsc.CloseWskConn(); err != nil {
		zap.GetLogger().Error("close websocket connection failed", "message", err)
	}
}

/*func (wsc *WebSocketClient) MessageLoop() {
	defer func() {
		if r := recover(); r != nil {
			zap.GetLogger().Error("receive message exception", "error", r)
		}
		if wsc.wskConn != nil {
			wsc.wskConn.Close()
		}
	}()
	wsc.wskConn.SetReadDeadline(time.Now().Add(heartbeatInterval))
	wsc.wskConn.SetPingHandler(func(appData string) error {
		zap.GetLogger().Info("receive ping message", "data", appData)
		return wsc.wskConn.WriteMessage(websocket.PongMessage, []byte("pong"))
	})
	for {
		if !wsc.isOpen {
			break
		}
		messageType, rawMessage, err := wsc.wskConn.ReadMessage()
		if err != nil {
			wsc.onError(err)
			break
		} else if messageType != websocket.TextMessage {
			wsc.onError(fmt.Errorf("invalid message type: %v only support text or binary msg", messageType))
		} else {
			if wsc.onMessage != nil {
				if err = wsc.onMessage(rawMessage); err != nil {
					wsc.onError(err)
				}
			}
		}
	}
}
*/

func (wsc *WebSocketClient) OpenWskConn(wskConn *websocket.Conn) {
	wsc.Lock()
	defer wsc.Unlock()
	wsc.wskConn = wskConn
	wsc.isOpen = true
	wsc.wskConn.SetPingHandler(func(appData string) error {
		zap.GetLogger().Info("receive ping request from client, begin handling......", "request", appData)
		time.Sleep(1 * time.Second)
		if err := wsc.wskConn.WriteMessage(websocket.PongMessage, []byte("pong")); err != nil {
			return err
		}
		zap.GetLogger().Info("handling ping request completed")
		return nil
	})
}

func (wsc *WebSocketClient) CloseWskConn() error {
	if !wsc.isOpen {
		return nil
	}
	if wsc.wskConn == nil {
		return errors.New("websocket connection is nil")
	}

	wsc.Lock()
	defer wsc.Unlock()
	wsc.isOpen = false
	return wsc.wskConn.Close()
}

func (wsc *WebSocketClient) SendMessage(message []byte) error {
	return wsc.wskConn.WriteMessage(websocket.TextMessage, message)
}

func (wsc *WebSocketClient) SendData(data []byte) error {
	return wsc.wskConn.WriteMessage(websocket.BinaryMessage, data)
}

func NewWebSocketClient(opts ...WebSocketClientOptions) *WebSocketClient {
	wsc := &WebSocketClient{
		isOpen: false,
	}
	if len(opts) > 0 {
		for _, opt := range opts {
			opt(wsc)
		}
	}

	return wsc
}

func WithMessageHandler(handler func([]byte) error) WebSocketClientOptions {
	return func(c *WebSocketClient) {
		c.onMessage = handler
	}
}

func WithDataHandler(handler func([]byte) error) WebSocketClientOptions {
	return func(c *WebSocketClient) {
		c.onData = handler
	}
}

func WithErrorHandler(handler func(error)) WebSocketClientOptions {
	return func(c *WebSocketClient) {
		c.onError = handler
	}
}

func WithChannelOpen() WebSocketClientOptions {
	return func(c *WebSocketClient) {
		c.isOpen = true
	}
}

func WithChannelClose() WebSocketClientOptions {
	return func(c *WebSocketClient) {
		c.isOpen = false
	}
}
