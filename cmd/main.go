package main

import (
	"fmt"
	"sync"

	"alex-shch/logger"
	"alex-shch/tcp-msg-server"
	"alex-shch/tcp-msg-server/cmd/client"
)

// implementation server.Handlers
type Handlers struct {
}

func (Handlers) OnConnect(inMsgs <-chan []byte, outMsgs chan<- []byte) {
	go func() {
		msg := <-inMsgs
		fmt.Println("serv <-- msg: ", string(msg))

		outMsg := string(msg) + " * 2 = " + string(msg) + string(msg)
		fmt.Println("serv --> msg: ", outMsg)
		outMsgs <- []byte(outMsg)
	}()
}

func main() {
	wg := &sync.WaitGroup{}

	log := logger.NewLogger(logger.INFO)

	server, err := server.NewServer("localhost:4567", Handlers{}, log, wg)
	if err != nil {
		panic(err)
	}

	go server.Run()

	// client part
	client, err := client.NewClient("localhost:4567", log)
	if err != nil {
		panic(err)
	}
	outMsg := "123abc"
	fmt.Println("client -> msg: ", outMsg)
	//client.SendMessage([]byte(outMsg))
	client.Out.Msgs() <- []byte(outMsg)
	msg := <-client.In.Msgs()
	fmt.Println("cient <- msg: ", string(msg))
	client.Disconnect()

	//fmt.Scanln()

	server.Stop()
	wg.Wait()
}
