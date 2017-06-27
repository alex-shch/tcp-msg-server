package tcpserver

import (
	"fmt"
	"io"
	"net"

	"alex-shch/logger"
)

type outStream struct {
	log    logger.Logger
	msgs   chan []byte
	conn   net.Conn
	offset int
	exit   bool
}

func (self *outStream) waitForMsg() {
	header := [8]byte{}

	for msg := range self.msgs {
		msgSize := len(msg)
		strHeader := fmt.Sprintf("%08x", msgSize)
		copy(header[:], strHeader)

		self.offset = 0
		for self.offset < _HDR_SIZE {
			sent, err := self.conn.Write(header[self.offset:_HDR_SIZE])
			if err != nil {
				self.handleReadError(err)
				return
			}

			self.offset += sent
		}

		self.offset = 0
		for self.offset < msgSize {
			sent, err := self.conn.Write(msg[self.offset:msgSize])
			if err != nil {
				self.handleReadError(err)
				return
			}

			self.offset += sent
		}
	}
}

func (self *outStream) handleReadError(err error) {
	if err == io.EOF {
		self.log.Debugf("Remote host %s close connection", self.conn.RemoteAddr())
	} else if self.exit {
		self.log.Debugf("Cancel write from %s", self.conn.RemoteAddr())
	} else {
		self.log.Error("Error writing: ", err.Error())
		self.conn.Close()
	}
}
