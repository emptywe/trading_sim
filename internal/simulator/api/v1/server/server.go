package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"go.uber.org/zap"
)

const maxPorts = 65535

var Host string
var Port int

type Server struct {
	HttpServer *http.Server
}

func (s *Server) run(handler http.Handler) error {

	s.HttpServer = &http.Server{
		Addr:           Host + ":" + fmt.Sprintf("%d", Port),
		Handler:        handler,
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}
	return s.HttpServer.ListenAndServe()
}

func (s *Server) shutdown(ctx context.Context) error {
	return s.HttpServer.Shutdown(ctx)
}

func (s *Server) StartNewServer(handler http.Handler) {
	if err := s.run(handler); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			if strings.Contains(err.Error(), "address already in use") && Port < maxPorts {
				Port++
				go s.StartNewServer(handler)
				zap.S().Errorf("Can't run http simulator: %s", err)
				return
			}
			zap.S().Fatalf("Can't run Trading Simulator: %s", err)
		}
	}
}

func (s *Server) WaitServer() {
	ts := time.Now()
	for {
		if s.HttpServer != nil {
			zap.S().Infof("Trading Simulator Start on %s", s.HttpServer.Addr)
			return
		} else if time.Since(ts) > time.Second {
			zap.S().Fatal("unable to start sever")
			return
		} else {
			time.Sleep(time.Second)
		}
	}
}

func (s *Server) StopServer() {
	zap.S().Info("Trading Simulator Shutdown")
	if err := s.shutdown(context.Background()); err != nil {
		zap.S().Error("can't shout down server: " + err.Error())
	}
}
