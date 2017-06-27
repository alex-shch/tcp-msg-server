package tcpserver

import (
	"fmt"
	"net"

	"alex-shch/logger"
)

type outStream struct {
	log     logger.Logger
	msgs    chan []byte
	conn    net.Conn
	header  [8]byte
	offset  int
	msgSize int
	lock    chan struct{}
}

func (self *outStream) waitForMsg() {
	// TODO exit event
	for msg := range self.msgs {
		self.msgSize = len(msg)
		header := fmt.Sprintf("%08x", self.msgSize)
		copy(self.header[:], header)

		self.offset = 0
		for self.offset < _HDR_SIZE {
			sent, err := self.conn.Write(self.header[self.offset:_HDR_SIZE])
			if err != nil {
				self.log.Error(err)
				self.conn.Close()
				return
			}

			self.offset += sent
		}

		self.offset = 0
		for self.offset < self.msgSize {
			sent, err := self.conn.Write(msg[self.offset:self.msgSize])
			if err != nil {
				self.log.Error(err)
				self.conn.Close()
				return
			}

			self.offset += sent
		}
	}
}
