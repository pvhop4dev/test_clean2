package grpc

import (
	"net"

	"clean-arch-go/internal/domain/service"

	"google.golang.org/grpc"
)

type Server struct {
	server  *grpc.Server
	authSvc service.AuthService
	bookSvc service.BookService
}

func NewServer(authSvc service.AuthService, bookSvc service.BookService) *Server {
	s := &Server{
		server:  grpc.NewServer(),
		authSvc: authSvc,
		bookSvc: bookSvc,
	}

	// Đăng ký các service gRPC ở đây
	// pb.RegisterAuthServiceServer(s.server, s)
	// pb.RegisterBookServiceServer(s.server, s)

	return s
}

func (s *Server) Start(addr string) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	return s.server.Serve(lis)
}

func (s *Server) Stop() {
	s.server.GracefulStop()
}
