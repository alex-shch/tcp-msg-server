package tcpserver

import (
	"net"
	"sync"

	"github.com/alex-shch/logger"
)

const (
	_HDR_SIZE = 8
)

type ConnHandler interface {
	InMsgs() <-chan []byte
	Send(msg []byte) error
	Disconnect()
}

type _ConnHandler struct {
	id   uint64
	in   inStream
	out  outStream
	log  logger.Logger
	done chan struct{}
}

func NewConnHandler(conn net.Conn, log logger.Logger) ConnHandler {
	return newConnHandler(0, conn, log)
}

func newConnHandler(id uint64, conn net.Conn, log logger.Logger) *_ConnHandler {
	connHandler := &_ConnHandler{
		id: id,
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
		log:  log,
		done: make(chan struct{}),
	}

	connHandler.run()

	return connHandler
}

func (self *_ConnHandler) InMsgs() <-chan []byte {
	return self.in.msgs
}

func (self *_ConnHandler) Send(msg []byte) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()

	self.out.msgs <- msg
	return nil
}

func (self *_ConnHandler) Disconnect() {
	self.out.exit = true
	self.in.exit = true
	self.in.conn.Close()
	<-self.done
}

func (self *_ConnHandler) run() {
	wg := &sync.WaitGroup{}
	wg.Add(2) // wait for completed 2 routines

	go func() {
		defer wg.Done()

		// when read routine completed, close message channels
		defer close(self.out.msgs)
		defer close(self.in.msgs)

		self.in.waitForMsg()
	}()

	go func() {
		defer wg.Done()
		self.out.waitForMsg()
	}()

	go func() {
		wg.Wait() // wain in+out routines

		// done event when in+out routines completed
		close(self.done)
	}()
}
