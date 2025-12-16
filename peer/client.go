package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/google/uuid"
)

func handleMessages() {
	fmt.Println("[CLIENT] Type a message and press ENTER. Ctrl+C to exit.")
	
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		text, _ := reader.ReadString('\n')

		server.handleServerMessages(uuid.New().String(), server.Address, text)
	}
}
