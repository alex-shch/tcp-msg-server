package main

import (
	"fmt"
	"net"

	"alex-shch/logger"
	"alex-shch/tcp-msg-server"
)

// implementation server.Handlers
type ConnectionCallback struct {
}

func (ConnectionCallback) OnConnect(connId uint64, inMsgs <-chan []byte, send func([]byte) error, done <-chan struct{}) {

	fmt.Printf("serv (%d) connect\n", connId)
	go func() {
		msg := <-inMsgs
		fmt.Printf("serv (%d) <-- msg: %s\n", connId, string(msg))

		outMsg := string(msg) + " * 2 = " + string(msg) + string(msg)
		fmt.Printf("serv (%d) --> msg: %s\n", connId, outMsg)
		if err := send([]byte(outMsg)); err != nil {
			fmt.Println("Send error: ", err)
		}
	}()

	go func() {
		<-done
		fmt.Printf("serv (%d) disconnect\n", connId)
	}()
}

func main() {
	log := logger.NewLogger(logger.DEBUG)

	serverAddr := "localhost:4567"
	server, err := tcpserver.NewServer(serverAddr, ConnectionCallback{}, log)
	if err != nil {
		panic(err)
	}

	go server.Run()

	// client part
	for i := 0; i < 1; i++ {
		conn, err := net.Dial("tcp", serverAddr)
		if err != nil {
			panic(err)
		}

		connHandler := tcpserver.NewConnHandler(conn, log)
		outMsg := "123abc"
		fmt.Println("client -> msg: ", outMsg)
		//connHandler.OutMsgs() <- []byte(outMsg)
		connHandler.Send([]byte(outMsg))

		msg := <-connHandler.InMsgs()
		fmt.Println("cient <- msg: ", string(msg))
		connHandler.Disconnect()
		fmt.Println()
	}

	//fmt.Scanln()

	server.Stop()
	<-server.Done()
}
