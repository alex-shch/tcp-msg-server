package tcpserver

import (
	"net"
	"sync"

	"alex-shch/logger"
)

type TcpServer interface {
	Run()
	Stop()
}

type ConnectionCallback interface {
	OnConnect(connId uint64, inMsgs <-chan []byte, outMsgs chan<- []byte, done <-chan struct{})
}

type _TcpServer struct {
	idCounter   uint64
	log         logger.Logger
	wg          *sync.WaitGroup
	listener    net.Listener
	callback    ConnectionCallback
	connections map[uint64]*_ConnHandler
}

func NewServer(addr string, callback ConnectionCallback, log logger.Logger, wg *sync.WaitGroup) (TcpServer, error) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Error("Error listening:", err.Error())
		return nil, err
	}

	log.Info("Listening on ", addr)

	return &_TcpServer{
		log:         log,
		wg:          wg,
		listener:    l,
		callback:    callback,
		connections: make(map[uint64]*_ConnHandler),
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

		self.idCounter++

		connHandler := newConnHandler(self.idCounter, conn, self.log)
		self.connections[self.idCounter] = connHandler

		self.wg.Add(1)
		go func() {
			// server waits for close all clients connections
			defer self.wg.Done()
			<-connHandler.done
		}()

		go self.callback.OnConnect(
			self.idCounter,
			connHandler.in.msgs,
			connHandler.out.msgs,
			connHandler.done,
		)
	}
}

func (self *_TcpServer) Stop() {
	self.listener.Close()

	// TODO сделать потокобезопасным
	for _, h := range self.connections {
		h.Disconnect()
	}
}
