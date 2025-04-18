package httpd

import (
	"github.com/buglloc/yubictld/internal/touchctl"
	"github.com/buglloc/yubictld/internal/ykman"
)

type Option func(server *Server)

func WithAddr(addr string) Option {
	return func(s *Server) {
		s.addr = addr
	}
}

func WithToucher(t touchctl.Toucher) Option {
	return func(s *Server) {
		s.touch = t
	}
}

func WithYkMan(yk *ykman.YkMan) Option {
	return func(s *Server) {
		s.yk = yk
	}
}
