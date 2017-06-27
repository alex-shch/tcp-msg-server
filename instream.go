package tcpserver

import (
	"io"
	"net"
	"strconv"

	"github.com/alex-shch/logger"
)

type inStream struct {
	log    logger.Logger
	msgs   chan []byte
	conn   net.Conn
	buf    []byte
	offset int
	exit   bool
}

func (self *inStream) waitForMsg() {
	for {
		// receive header
		for self.offset < _HDR_SIZE {
			recvSize, err := self.conn.Read(self.buf[self.offset:])
			if err != nil {
				self.handleReadError(err)
				return
			}
			self.offset += recvSize
		}

		// parse header
		msgSize64, err := strconv.ParseInt(string(self.buf[0:_HDR_SIZE]), 16, 32)
		if err != nil {
			self.log.Error("Error parsing: ", err)
			return
		}
		msgSize := int(msgSize64)

		// receive body
		for self.offset < _HDR_SIZE+msgSize {
			recvSize, err := self.conn.Read(self.buf[self.offset:])
			if err != nil {
				self.handleReadError(err)
				return
			}
			self.offset += recvSize
		}

		// store body
		buf := make([]byte, msgSize)
		copy(buf, self.buf[_HDR_SIZE:_HDR_SIZE+msgSize])
		self.msgs <- buf

		// prepare to receive next message
		pkgSize := _HDR_SIZE + msgSize
		if self.offset > pkgSize {
			copy(self.buf, self.buf[pkgSize:self.offset])
		}
		self.offset -= pkgSize
	}
}

func (self *inStream) handleReadError(err error) {
	if err == io.EOF {
		self.log.Debugf("Remote host %s close connection", self.conn.RemoteAddr())
	} else if self.exit {
		self.log.Debugf("Cancel read from %s", self.conn.RemoteAddr())
	} else {
		self.log.Error("Error reading: ", err.Error())
		self.conn.Close()
	}
}
