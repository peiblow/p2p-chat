package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/google/uuid"
)

func startClient(ipAddress string) {
	time.Sleep(1 * time.Second)
	
	if ipAddress == "" {
		fmt.Println("[CLIENT] No peer address provided")
		return
	}

	fmt.Println(ipAddress)
	conn, err := connectToPeer(ipAddress)
	if err != nil {
		fmt.Println("[CLIENT] dial error:", err)
		return
	}
	
	defer conn.Close()
	fmt.Println("[CLIENT] Connected to server!")
	
	handleMessages(conn)
}

func handleMessages(conn net.Conn) {
	fmt.Println("[CLIENT] Type a message and press ENTER. Ctrl+C to exit.")
	
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		text, _ := reader.ReadString('\n')

		err := sendMessage(conn, uuid.New().String(), conn.RemoteAddr().String(), text)
		if err != nil {
			fmt.Println("[CLIENT] send error:", err)
			return
		}
	}
}
