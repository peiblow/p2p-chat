package main

import (
	"fmt"
	"time"
)

func main() {
	server := NewServer()
	go startServer()
	
	time.Sleep(2 * time.Second)
	myIp, _ := getMyIpAddress()
	bootstrapRegister(myIp)
	
	time.Sleep(2 * time.Second)
	server.fetchPeerList()

	fmt.Println("[MAIN] Peer list:", server.PeerList)
	if len(server.PeerList) > 0{
		go startClient(server.PeerList[0].Address)
	} else {
		fmt.Println("[MAIN] No other peers available, waiting for connections...")
	}

	select {}
}