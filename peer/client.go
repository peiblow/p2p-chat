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
		myIp, err := getMyIpAddress()
		if err != nil {
			fmt.Println("[CLIENT] Error getting my IP address:", err)
			return
		}

		ipAddress = myIp
	}

	conn, err := net.Dial("tcp", ipAddress + ":8080")
	if err != nil {
		fmt.Println("[CLIENT] dial error:", err)
		return
	}
	defer conn.Close()

	fmt.Println("[CLIENT] Connected to server!")

	conn.Write([]byte("Hello from client!"))
}
