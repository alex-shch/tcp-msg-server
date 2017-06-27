package tcpserver

import (
	"net"
	"sync"

	"github.com/alex-shch/logger"
)

type TcpServer interface {
	Run()
	Stop()
	Done() <-chan struct{}
}

type ConnectionCallback interface {
	//OnConnect(connId uint64, inMsgs <-chan []byte, outMsgs chan<- []byte, done <-chan struct{})
	OnConnect(connId uint64, inMsgs <-chan []byte, send func([]byte) error, done <-chan struct{})
}

type _TcpServer struct {
	idCounter uint64
	log       logger.Logger
	wg        *sync.WaitGroup
	listener  net.Listener
	callback  ConnectionCallback

	exit chan struct{} // for server start shutdown
	done chan struct{} // wnen server has been shutting down
}

func NewServer(addr string, callback ConnectionCallback, log logger.Logger) (TcpServer, error) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Error("Error listening:", err.Error())
		return nil, err
	}

	log.Debug("Listening on ", addr)

	return &_TcpServer{
		log:      log,
		wg:       &sync.WaitGroup{},
		listener: l,
		callback: callback,
		exit:     make(chan struct{}),
		done:     make(chan struct{}),
	}, nil
}

func (self *_TcpServer) Run() {
	self.wg.Add(1)
	defer self.wg.Done()

	for {
		conn, err := self.listener.Accept()
		if err != nil {
			if _, closed := <-self.exit; closed {
				self.log.Error("Error accepting: ", err.Error())
			} else {
				self.log.Debug("Cancel accepting")
			}
			break
		}

		self.idCounter++

		self.log.Debugf("accept connection (%d) from %s", self.idCounter, conn.RemoteAddr())

		connHandler := newConnHandler(self.idCounter, conn, self.log)

		// for close connection when server is down
		self.wg.Add(1)
		go func() {
			defer self.wg.Done()
			select {
			case <-self.exit: // server is shutting down
				connHandler.Disconnect()
			case <-connHandler.done: // connection closed
			}
		}()

		// live connection counter
		self.wg.Add(1)
		go func() {
			defer self.wg.Done()
			<-connHandler.done // connection closed
			self.log.Debugf("close connection (%d) from %s", self.idCounter, conn.RemoteAddr())
		}()

		// new connection callback
		self.wg.Add(1)
		go func() {
			defer self.wg.Done()
			self.callback.OnConnect(
				self.idCounter,
				connHandler.in.msgs,
				//connHandler.out.msgs,
				connHandler.Send,
				connHandler.done,
			)
		}()
	}
}

func (self *_TcpServer) Stop() {
	close(self.exit)
	self.listener.Close()
	go func() {
		self.wg.Wait()
		close(self.done)
	}()
}

func (self *_TcpServer) Done() <-chan struct{} {
	return self.done
}
