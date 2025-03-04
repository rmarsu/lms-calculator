package server

import (
	"context"
	"net/http"
	"os"
	"time"
)

type Server struct {
	Port       string
	httpServer *http.Server
}

type ServerConfig struct {
	Port               string
	ReadTimeout        time.Duration
	WriteTimeout       time.Duration
	MaxHeaderMegabytes int
}

func NewDefaultServerConfig() *ServerConfig {
	return &ServerConfig{
		Port:               os.Getenv("ORCHESTRATOR_PORT"),
		ReadTimeout:        10 * time.Second,
		WriteTimeout:       10 * time.Second,
		MaxHeaderMegabytes: 1,
	}
}

func NewServer(cfg *ServerConfig, handler http.Handler) *Server {
	return &Server{
		Port: cfg.Port,
		httpServer: &http.Server{
			Addr:           ":" + cfg.Port,
			Handler:        handler,
			ReadTimeout:    cfg.ReadTimeout,
			WriteTimeout:   cfg.WriteTimeout,
			MaxHeaderBytes: cfg.MaxHeaderMegabytes << 20,
			IdleTimeout:    time.Second * 5,
		},
	}
}

func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
