package main

import (
	"fmt"
	"net"
	"time"
)

func startClient(ipAddress string) {
	time.Sleep(1 * time.Second)
	
	if ipAddress == "" {
		fmt.Println("[CLIENT] No peer address provided")
		return
	}

	// ipAddress já vem com a porta incluída (ex: 172.20.0.2:8080)
	conn, err := net.Dial("tcp", ipAddress)
	if err != nil {
		fmt.Println("[CLIENT] dial error:", err)
		return
	}
	defer conn.Close()

	fmt.Println("[CLIENT] Connected to server!")

	conn.Write([]byte("Hello from client!"))
}
