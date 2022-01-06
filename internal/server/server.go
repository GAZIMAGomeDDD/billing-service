package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/sirupsen/logrus"
)

type Server struct {
	httpServer *http.Server
	address    string
}

func NewServer(handler http.Handler, addr string) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:              addr,
			Handler:           handler,
			ReadTimeout:       30 * time.Second,
			WriteTimeout:      30 * time.Second,
			ReadHeaderTimeout: 30 * time.Second,
		},
	}
}

func (s *Server) Run() error {
	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt)

	go func() {
		for range c {
			s.Stop()
			break
		}
	}()

	logrus.Info("starting server..")

	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop() {
	logrus.Info("shutting down..")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	s.httpServer.Shutdown(ctx)
	cancel()
}
