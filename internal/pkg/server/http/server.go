package http

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Server struct {
	*gin.Engine
	srv *http.Server
}

func NewServer() *Server {
	engine := gin.Default()
	return &Server{
		Engine: engine,
		srv: &http.Server{
			Handler: engine,
		},
	}
}

func (s *Server) Start(addr string) error {
	s.srv.Addr = addr
	return s.srv.ListenAndServe()
}

func (s *Server) AddRoute(method, path string, handler gin.HandlerFunc) {
	s.Handle(method, path, handler)
}

func (s *Server) Use(middleware ...gin.HandlerFunc) {
	s.Use(middleware...)
}

func (s *Server) Group(path string) *gin.RouterGroup {
	return s.Group(path)
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
