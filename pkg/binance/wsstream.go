package binance

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"time"
)

type Income struct {
	Name  string
	Value float64
}

func connectAndSubscribe(cur, url string) (*websocket.Conn, error) {
	conn, resp, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("handshake failed with status %d, error: %v", resp.StatusCode, err))
	}

	reqc := connRequest{method: subscribe, params: []string{cur + trade}, id: 1}

	reqm, err := json.Marshal(reqc)
	if err != nil {
		conn.Close()
		return nil, err
	}

	err = conn.WriteMessage(websocket.TextMessage, reqm)
	if err != nil {
		conn.Close()
		return nil, err
	}

	return conn, nil
}

func disconnectAndUnsubscribe(cur string, conn *websocket.Conn) error {
	reqd := connRequest{method: unsubscribe, params: []string{cur + trade}, id: 312}
	zap.S().Info("interrupt")
	stop, err := json.Marshal(reqd)
	if err != nil {
		zap.S().Infof("can't encode interrupt request, error: %v", err)
	}
	err = conn.WriteMessage(websocket.TextMessage, stop)
	if err != nil {
		zap.S().Errorf("write close: %v", err)
		return err
	}

	select {
	case <-time.After(time.Second):
	}
	zap.S().Info("WSDone")
	return nil
}

func (c *Client) WSHandlerBinance(cur string) {
	start := time.Now()
	res := new(DataPrice)
	Interrupt := make(chan os.Signal, 1)
	signal.Notify(Interrupt, os.Interrupt)
	conn, err := connectAndSubscribe(cur, wsUrl)
	if err != nil {
		zap.S().Infof("Can't establish connection: %v", err)
		return
	}

	defer conn.Close()
	ticker := time.NewTicker(time.Millisecond * tradeUpd)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:

			_, message, err := conn.ReadMessage()
			if err != nil {
				zap.S().Infof("read: %v", err)
				return
			}

			err = json.Unmarshal(message, res)
			if err != nil {
				logrus.Printf("can't decode data: %s", err.Error())
			}
			c.Data <- *res
			if err != nil {
				logrus.Printf("can't update currency: %s", err.Error())
			}
			if time.Now().Unix()-start.Unix() > int64(time.Hour*23) {
				err = disconnectAndUnsubscribe(cur, conn)
				if err != nil {
					logrus.Println(err)
					return
				}
				conn, err = connectAndSubscribe(cur, wsUrl)
				if err != nil {
					logrus.Println(err)
					return
				}
				start = time.Now()
			}
		case <-Interrupt:
			err = disconnectAndUnsubscribe(cur, conn)
			if err != nil {
				logrus.Println("close error:", err)
				return
			}
		}
	}
}
