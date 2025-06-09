package grpc

import (
	"google.golang.org/grpc"
)

type Server struct {
	*grpc.Server
}

func NewServer() *Server {
	return &Server{
		Server: grpc.NewServer(),
	}
}
