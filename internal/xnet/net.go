package xnet

import (
	"errors"
	"net"
)

func NewListener(addr string) (net.Listener, error) {
	if len(addr) == 0 {
		return nil, errors.New("empty address")
	}

	return net.Listen(ParseNetwork(addr), addr)
}

func ParseNetwork(addr string) string {
	switch addr[0] {
	case '/', '.', '@':
		return "unix"

	default:
		return "tcp"
	}
}
