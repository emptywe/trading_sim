package binancews

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type Income struct {
	Name  string
	Value float64
}

func (ws *WSClient) connectToWS() (err error) {
	var resp *http.Response
	ws.wsConn, resp, err = websocket.DefaultDialer.Dial(wsUrl, nil)
	if err != nil {
		b, _ := io.ReadAll(resp.Body)
		return errors.New(fmt.Sprintf("handshake failed error: %v %v", err, string(b)))
	}
	zap.S().Infof("WS connected to %s", wsUrl)
	return nil
}

func (ws *WSClient) Subscribe(params []string) (err error) {
	reqc := connRequest{Method: subscribe, Params: params, Id: 1}

	reqm, err := json.Marshal(reqc)
	if err != nil {
		return err
	}
	err = ws.wsConn.WriteMessage(websocket.TextMessage, reqm)
	if err != nil {
		return err
	}
	return nil
}

func (ws *WSClient) Unsubscribe(params []string) (err error) {
	reqd := connRequest{Method: unsubscribe, Params: params, Id: 312}
	stop, err := json.Marshal(reqd)
	if err != nil {
		zap.S().Errorf("can't encode interrupt request, error: %v", err)
	}
	err = ws.wsConn.WriteMessage(websocket.TextMessage, stop)
	if err != nil {
		zap.S().Errorf("write close: %v", err)
		return err
	}
	return nil
}

func (ws *WSClient) disconnect() error {
	if err := ws.wsConn.Close(); err != nil {
		return err
	}
	select {
	case <-time.After(time.Second):
	}
	return nil
}

func (ws *WSClient) reconnect(params []string) error {
	if err := ws.Unsubscribe(params); err != nil {
		zap.S().Error(err)
	}
	if err := ws.disconnect(); err != nil {
		zap.S().Error(err)
	}
	zap.S().Infof("reconnect %v", params)
	if err := ws.connectToWS(); err != nil {
		zap.S().Error(err)
		return err
	}
	if err := ws.Subscribe(params); err != nil {
		zap.S().Error(err)
		return err
	}
	return nil
}

func (ws *WSClient) WSHandlerBinance(params []string) {
	start := time.Now()
	res := new(DataPrice)
	Interrupt := make(chan os.Signal, 1)
	signal.Notify(Interrupt, os.Interrupt)
	if err := ws.connectToWS(); err != nil {
		zap.S().Errorf("can't establish connection: %v", err)
		return
	}
	if err := ws.Subscribe(params); err != nil {
		zap.S().Errorf("can't subscribe to %v", params)
		return
	}
	defer ws.wsConn.Close()
	ticker := time.NewTicker(time.Millisecond * tradeUpd)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:

			_, message, err := ws.wsConn.ReadMessage()
			if err != nil {
				zap.S().Errorf("error read: %v %v %s", err, params, string(message))
				err = ws.reconnect(params)
				if err != nil {
					return
				}
			}

			err = json.Unmarshal(message, res)
			if err != nil {
				zap.S().Errorf("can't decode data: %v", err)
			}
			ws.Data <- *res
			if err != nil {
				zap.S().Errorf("can't update currency: %s", err.Error())
			}
			if time.Since(start) > time.Hour*23+time.Minute*55 {
				err = ws.reconnect(params)
				if err != nil {
					return
				}
				start = time.Now()
			}
		case <-Interrupt:
			if err := ws.Unsubscribe(params); err != nil {
				zap.S().Error(err)
			}
			if err := ws.disconnect(); err != nil {
				zap.S().Error(err)
				return
			}
			return
		}
	}
}
