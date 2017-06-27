package server

import (
	"net"
	"strconv"

	"alex-shch/logger"
)

type InStream interface {
	Msgs() <-chan []byte
}

func NewInStream(conn net.Conn, log logger.Logger) InStream {
	stream := &inStream{
		log:  log,
		msgs: make(chan []byte, 8),
		conn: conn,
		buf:  make([]byte, 64*1024),
	}

	go stream.waitForMsg()

	return stream
}

type inStream struct {
	log    logger.Logger
	msgs   chan []byte
	conn   net.Conn
	buf    []byte
	offset int
}

func (self *inStream) Msgs() <-chan []byte {
	return self.msgs
}

func (self *inStream) waitForMsg() {

	for { // TODO wait for exit event

		// receive header
		for self.offset < _HDR_SIZE {
			recvSize, err := self.conn.Read(self.buf[self.offset:])
			if err != nil {
				self.log.Error("Error reading: ", err.Error())
				self.conn.Close()
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
				self.log.Error("Error reading: ", err.Error())
				self.conn.Close()
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
