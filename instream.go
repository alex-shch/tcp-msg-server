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

	go stream.readHeader()

	return stream
}

type inStream struct {
	log     logger.Logger
	msgs    chan []byte
	conn    net.Conn
	buf     []byte
	offset  int
	msgSize int
}

func (self *inStream) Msgs() <-chan []byte {
	return self.msgs
}

func (self *inStream) readHeader() {
	self.log.Debug("readHeader")

	reqLen, err := self.conn.Read(self.buf[self.offset:])
	if err != nil {
		self.log.Error("Error reading: ", err.Error())
		self.conn.Close()
		return
	}
	self.log.Debugf("received %d bytes", reqLen)

	self.offset += reqLen

	if self.offset < _HDR_SIZE {
		go self.readHeader()
	} else {
		self.processHeader()
	}
}

func (self *inStream) processHeader() {
	self.log.Debug("processHeader")

	size, err := strconv.ParseInt(string(self.buf[0:_HDR_SIZE]), 16, 32)
	if err != nil {
		self.log.Error("Error parsing: ", err)
		return
	}
	self.msgSize = int(size)

	if self.offset < _HDR_SIZE+self.msgSize {
		go self.readBody()
	} else {
		self.processBody()
	}
}

func (self *inStream) processBody() {
	self.log.Debug("processBody")

	buf := make([]byte, self.msgSize)
	copy(buf, self.buf[_HDR_SIZE:_HDR_SIZE+self.msgSize])
	self.msgs <- buf

	pkgSize := _HDR_SIZE + self.msgSize
	if self.offset == pkgSize {
		self.offset = 0
		go self.readHeader()
		return
	}

	self.log.Debugf("pkgSize: %d, offset: %d", pkgSize, self.offset)
	copy(self.buf, self.buf[pkgSize:self.offset])
	self.offset -= pkgSize
	if self.offset < _HDR_SIZE {
		go self.readHeader()
	} else {
		self.processHeader()
	}
}

func (self *inStream) readBody() {
	self.log.Debug("readBody")

	reqLen, err := self.conn.Read(self.buf[self.offset:])
	if err != nil {
		self.log.Error("Error reading: ", err.Error())
		self.conn.Close()
		return
	}
	self.log.Debugf("received %d bytes", reqLen)

	self.offset += reqLen

	if self.offset < _HDR_SIZE+self.msgSize {
		go self.readBody()
	} else {
		self.processBody()
	}
}
