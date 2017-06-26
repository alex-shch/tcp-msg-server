package main

import (
	"fmt"
	"sync"

	"alex-shch/tcp-msg-server"
	"alex-shch/tcp-msg-server/cmd/client"
	"alex-shch/tcp-msg-server/cmd/log"
)

// implementation server.Handlers
type Handlers struct {
}

func (Handlers) OnConnect(inMsgs <-chan []byte, outMsgs chan<- []byte) {
}

func main() {
	wg := &sync.WaitGroup{}

	log := log.NewLogger(log.DEBUG)

	server, err := server.NewServer("localhost:4567", Handlers{}, log, wg)
	if err != nil {
		panic(err)
	}

	go server.Run()

	client, err := client.NewClient("localhost:4567")
	if err != nil {
		panic(err)
	}
	client.SendMessage([]byte("123abc"))

	//sendMsg()

	fmt.Scanln()

	client.Disconnect()

	server.Stop()
	wg.Wait()
}

func sendMsg() {
	//conn := net.
}
