package client

import (
	"fmt"
	"net"

	logPkg "alex-shch/tcp-msg-server/cmd/log"
)

var log = logPkg.NewLogger(logPkg.DEBUG)

type Client struct {
	conn net.Conn
}

func NewClient(addr string) (*Client, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Error(err)
	}

	return &Client{
		conn: conn,
	}, nil
}

func (self *Client) Disconnect() {
	self.conn.Close()
}

func (self *Client) SendMessage(msg []byte) error {
	header := fmt.Sprintf("%08x", len(msg))

	sentHdr, err := self.conn.Write([]byte(header))
	if err != nil {
		log.Error(err)
		return err
	}
	sentBody, err := self.conn.Write([]byte(msg))
	if err != nil {
		log.Error(err)
		return err
	}

	log.Debugf("sent %d bytes", sentHdr+sentBody)

	return nil
}
