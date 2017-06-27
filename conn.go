package server

import (
	"net"

	"alex-shch/logger"
)

const (
	_HDR_SIZE = 8
)

type _ConnHandler struct {
	// TODO идентификатор соединения, клиента, версия протокола
	in  inStream
	out outStream
}

func newConnHandler(conn net.Conn, log logger.Logger) *_ConnHandler {
	// TODO придумать способ корректно закрывать соединения
	hdlr := &_ConnHandler{
		in: inStream{
			log:  log,
			msgs: make(chan []byte, 8),
			conn: conn,
			buf:  make([]byte, 64*1024), // TODO проверить, надо ли выделять заранее
		},
		out: outStream{
			log:  log,
			msgs: make(chan []byte, 8),
			conn: conn,
		},
	}
	return hdlr
}

func (self *_ConnHandler) run() {
	go self.in.waitForMsg()
	go self.out.waitForMsg()
}
