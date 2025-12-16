package main

import (
	"fmt"
	"time"
)

var server *Server

func main() {
	server = NewServer()
	go startServer(server)
	
	time.Sleep(2 * time.Second)
	bootstrapRegister(server.Address)
	
	time.Sleep(2 * time.Second)
	server.fetchPeerList()

	go handleMessages()

	fmt.Println("[MAIN] Server running, waiting for connections...")

	select {}
}