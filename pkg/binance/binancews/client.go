package binancews

import "github.com/gorilla/websocket"

type WSClient struct {
	wsConn *websocket.Conn
	Data   chan DataPrice
}

func NewBinanceWSClient() *WSClient {
	return &WSClient{Data: make(chan DataPrice, 10)}
}
