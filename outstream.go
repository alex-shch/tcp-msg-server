package server

import (
	"net"
)

type outStream struct {
	log  Logger
	msgs chan []byte
	conn net.Conn
}
