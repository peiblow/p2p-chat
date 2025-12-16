package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
)

type Bootstrap struct {
	peers map[string]bool
	mu sync.Mutex
}

func NewBootstrap() *Bootstrap {
	return &Bootstrap{
		peers: make(map[string]bool),
	}
}

func (b *Bootstrap) Start() {
	ln, err := net.Listen("tcp", ":9001")
	if err != nil {
		panic(err)
	}

	fmt.Println("[BOOTSTRAP] Listening on :9001")

	for {
		conn, _ := ln.Accept()
		go b.handleConnection(conn)
	}
}

func (b *Bootstrap) handleConnection(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			return
		}

		msg := strings.TrimSpace(message)

		if strings.HasPrefix(msg, "REGISTER_PEER") {
			peer := strings.TrimSpace(strings.TrimPrefix(msg, "REGISTER_PEER"))
			b.mu.Lock()
			b.peers[peer] = true
			b.mu.Unlock()

			conn.Write([]byte("REGISTERED\n"))

			// send to all peers about the new peer
			if len(b.peers) > 1 {
				fmt.Println("[BOOTSTRAP] Notifying peers about new peer:", peer)
				b.mu.Lock()
				for p := range b.peers {
					if p != peer {
						peerConn, err := net.Dial("tcp", p)
						if err != nil {
							fmt.Println("[BOOTSTRAP] Error notifying peer:", err)
							continue
						}
						peerConn.Write([]byte("NEW_PEER " + peer + "\n"))
						peerConn.Close()
					}
				}
				b.mu.Unlock()
			}

			fmt.Printf("[BOOTSTRAP] Registered peer: %s\n", peer)
		} else if msg == "GET_PEERS" {
			b.mu.Lock()
			peersList := make([]string, 0)
			for p := range b.peers {
					peersList = append(peersList, p)
			}
			b.mu.Unlock()

			response := strings.Join(peersList, ",") + "\n"
			conn.Write([]byte(response))
			continue
		} else if msg == "PING" {
			conn.Write([]byte("PONG\n"))
		} else if msg == "QUIT" {
			conn.Write([]byte("GOODBYE\n"))
			conn.Close()
			return
		} else {
			conn.Write([]byte("UNKNOWN_COMMAND\n"))
		}
	}
}

func main() {
	bootstrap := NewBootstrap()
	go bootstrap.Start()

	select {}
}