package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"

	"github.com/emptywe/trading_sim/internal/storage"
	"github.com/emptywe/trading_sim/model"
)

type ForeignConnService struct {
	repo storage.ForeignConn
}

type Income struct {
	Name  string
	Value float64
}

func NewForeignConnService(repo storage.ForeignConn) *ForeignConnService {
	return &ForeignConnService{repo: repo}
}

// func (s *ForeignConnService) WShandlerBinance(cur string) {
func (s *ForeignConnService) WShandlerBinance(cur string) {
	start := time.Now()

	res := new(model.ConnResponse)
	Interrupt := make(chan os.Signal, 1)

	signal.Notify(Interrupt, os.Interrupt)

	url := fmt.Sprintf("wss://stream.binance.com:9443/ws")
	//logrus.Printf("connecting to %s", url)

	conn, err := s.connectAndSubscribe(cur, url)
	if err != nil {
		logrus.Println("Can't establish connection: ", err)
		return
	}

	defer conn.Close()
	ticker := time.NewTicker(time.Millisecond * 500)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:

			_, message, err := conn.ReadMessage()
			if err != nil {
				logrus.Println("read", err)
				return
			}

			err = json.Unmarshal(message, res)
			if err != nil {
				logrus.Printf("can't decode data: %s", err.Error())
			}

			fdata, err := strconv.ParseFloat(res.Price, 64)
			if err != nil {
				//logrus.Printf("wrong data format, errpr:%s", err.Error())
				//fmt.Println(string(message))
			}
			ndata := strings.ToLower(res.Symbol)
			data := Income{Name: ndata, Value: fdata}
			//fmt.Println(data)
			err = s.repo.UpdateCurrency(data.Name, data.Value)
			if err != nil {
				logrus.Printf("can't update currency: %s", err.Error())
			}
			if time.Now().Unix()-start.Unix() > int64(time.Hour*23) {
				err = s.disconnectAndUnsubscribe(cur, conn)
				if err != nil {
					logrus.Println(err)
					return
				}
				conn, err = s.connectAndSubscribe(cur, url)
				if err != nil {
					logrus.Println(err)
					return
				}
				start = time.Now()
			}

		case <-Interrupt:
			err = s.disconnectAndUnsubscribe(cur, conn)
			if err != nil {
				logrus.Println("close error:", err)
				return
			}
		}

	}

}

func (s *ForeignConnService) connectAndSubscribe(cur, url string) (*websocket.Conn, error) {
	conn, resp, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("handshake failed with status %d, error: %s", resp.StatusCode, err.Error()))
	}

	reqc := model.ConnRequest{Method: "SUBSCRIBE", Params: []string{cur + "@trade"}, Id: 1}

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

func (s *ForeignConnService) disconnectAndUnsubscribe(cur string, conn *websocket.Conn) error {
	reqd := model.ConnRequest{Method: "UNSUBSCRIBE", Params: []string{cur + "@trade"}, Id: 312}
	logrus.Println("interrupt")

	stop, err := json.Marshal(reqd)
	if err != nil {
		logrus.Printf("can't encode interrupt request, error: %s", err.Error())
	}
	err = conn.WriteMessage(websocket.TextMessage, stop)
	if err != nil {
		log.Println("write close:", err)
		return err
	}

	select {
	case <-time.After(time.Second):
	}
	logrus.Println("WSDone")
	return nil
}
