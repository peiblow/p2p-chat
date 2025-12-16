package main

import (
	"fmt"
	"net"
	"strings"
	"time"
)

func connectToPeer(address string) (net.Conn, error) {
	address = strings.TrimSpace(address)
	fmt.Printf("[DEBUG] Connecting to peer at address: %q\n", address)
	return net.DialTimeout("tcp", address, 3*time.Second)	
}

func sendMessage(conn net.Conn, msgUuid string, originPeer string, message string) error {
	_, err := conn.Write([]byte("NEW_MESSAGE " + msgUuid + " " + originPeer + " " + message))
	return err
}
