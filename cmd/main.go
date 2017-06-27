package main

import (
	"fmt"
	"net"
	"sync"

	"alex-shch/logger"
	"alex-shch/tcp-msg-server"
)

// implementation server.Handlers
type ConnectionCallback struct {
}

func (ConnectionCallback) OnConnect(id uint64, inMsgs <-chan []byte, outMsgs chan<- []byte, done <-chan struct{}) {
	fmt.Printf("serv (%d) connect\n", id)
	go func() {
		msg := <-inMsgs
		fmt.Printf("serv (%d) <-- msg: %s\n", id, string(msg))

		outMsg := string(msg) + " * 2 = " + string(msg) + string(msg)
		fmt.Printf("serv (%d) --> msg: %s\n", id, outMsg)
		outMsgs <- []byte(outMsg)
	}()

	go func() {
		<-done
		fmt.Printf("serv (%d) disconnect\n", id)
	}()
}

func main() {
	wg := &sync.WaitGroup{}

	log := logger.NewLogger(logger.INFO)

	serverAddr := "localhost:4567"
	server, err := tcpserver.NewServer(serverAddr, ConnectionCallback{}, log, wg)
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
		connHandler.OutMsgs() <- []byte(outMsg)
		msg := <-connHandler.InMsgs()
		fmt.Println("cient <- msg: ", string(msg))
		connHandler.Disconnect()
		fmt.Println()
	}

	//fmt.Scanln()

	server.Stop()
	wg.Wait()
}
