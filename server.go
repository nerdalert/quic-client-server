package main

import (
	"context"
	"fmt"

	"github.com/lucas-clemente/quic-go"
)

type Server struct {
	Addr    string
	Handler Handler
}

func NewServer(addr string) *Server {
	return &Server{
		Addr: addr,
	}
}

func (s *Server) SetHandler(handler Handler) {
	s.Handler = handler
}

func (s *Server) StartServer(ctx context.Context) error {
	listener, err := quic.ListenAddr(s.Addr, getTLSConfig(), &quic.Config{
		EnableDatagrams: true,
	})
	if err != nil {
		return err
	}
	for {
		conn, err := listener.Accept(ctx)
		if err != nil {
			return err
		}
		go func() {
			err := handleMsg(conn, s.Handler)
			if err != nil {
				fmt.Printf("handler err: %v", err)
			}
		}()
	}
}
