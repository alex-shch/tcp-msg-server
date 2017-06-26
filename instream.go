package server

import (
	"net"
	"strconv"
)

type inStream struct {
	log     Logger
	msgs    chan []byte
	conn    net.Conn
	buf     []byte
	offset  int
	msgSize int
}

func (self *inStream) readHeader() {
	self.log.Debug("readHeader")

	reqLen, err := self.conn.Read(self.buf[self.offset:])
	if err != nil {
		self.log.Error("Error reading: ", err.Error())
		self.conn.Close()
		return
	}
	self.log.Debugf("received %d bytes\n", reqLen)

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

	if len(self.buf) < _HDR_SIZE+self.msgSize {
		go self.readBody()
	} else {
		self.processBody()
	}
}

func (self *inStream) processBody() {
	self.log.Debug("processBody")

	buf := make([]byte, self.msgSize)
	copy(buf, self.buf[_HDR_SIZE:_HDR_SIZE+self.msgSize])
	self.log.Debug("<- msg: ", string(buf))
	self.msgs <- buf

	pkgSize := _HDR_SIZE + self.msgSize
	if self.offset == pkgSize {
		self.offset = 0
		go self.readHeader()
		return
	}

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
	self.log.Debugf("received %d bytes\n", reqLen)

	self.offset += reqLen

	if self.offset < _HDR_SIZE+self.msgSize {
		go self.readBody()
	} else {
		self.processBody()
	}
}
