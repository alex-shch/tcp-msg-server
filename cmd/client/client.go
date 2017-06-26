package client

import (
	"net"

	"alex-shch/logger"
	"alex-shch/tcp-msg-server"
)

type Client struct {
	conn net.Conn
	In   server.InStream
	Out  server.OutStream
}

func NewClient(addr string, log logger.Logger) (*Client, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Error(err)
	}

	client := &Client{
		conn: conn,
		In:   server.NewInStream(conn, log),
		Out:  server.NewOutStream(conn, log),
	}

	return client, nil
}

func (self *Client) Disconnect() {
	self.conn.Close()
}
