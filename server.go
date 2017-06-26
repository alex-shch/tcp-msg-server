package server

import (
	"net"
	"sync"
)

type TcpServer interface {
	Run()
	Stop()
}

type Handlers interface {
	OnConnect(inMsgs <-chan []byte, outMsgs chan<- []byte)
	// OnDisconnect()
}

type _TcpServer struct {
	log      Logger
	wg       *sync.WaitGroup
	listener net.Listener
	handlers Handlers
}

func NewServer(addr string, handlers Handlers, log Logger, wg *sync.WaitGroup) (TcpServer, error) {
	if log == nil {
		log = nullLogger{}
	}

	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Error("Error listening:", err.Error())
		return nil, err
	}

	log.Debug("Listening on ", addr)

	return &_TcpServer{
		log:      log,
		wg:       wg,
		listener: l,
		handlers: handlers,
	}, nil
}

func (self *_TcpServer) Run() {
	self.wg.Add(1)
	defer self.wg.Done()

	for {
		conn, err := self.listener.Accept()
		if err != nil {
			self.log.Error("Error accepting: ", err.Error())
			break
		}

		connHandler := newConnHandler(conn, self.log)
		connHandler.run()
		self.handlers.OnConnect(connHandler.in.msgs, connHandler.out.msgs)
	}
}

func (self *_TcpServer) Stop() {
	self.listener.Close()
}