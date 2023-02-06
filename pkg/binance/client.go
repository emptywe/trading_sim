package binance

import "github.com/gorilla/websocket"

type Client struct {
	wsConn *websocket.Conn
	Data   chan DataPrice
}

func NewBinanceClient() *Client {
	return &Client{Data: make(chan DataPrice, 10)}
}
