package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

type Server struct {
	Address	 string
	PeerList []PeerInfo
	Messages  map[string]bool
}

func NewServer() *Server {
	return &Server{
		PeerList: []PeerInfo{},
	}
}

func getMyIpAddress() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "unknown", fmt.Errorf("no valid IPv4 address found")
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}
	return "unknown", fmt.Errorf("no valid IPv4 address found")
}

func getBootstrapAddress() string {
	host := os.Getenv("BOOTSTRAP_HOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("BOOTSTRAP_PORT")
	if port == "" {
		port = "9001"
	}
	return host + ":" + port
}

func startServer(s *Server) {
	ln, err := net.Listen("tcp", ":8081")
	if err != nil {
		panic(err)
	}

	myIp, _ := getMyIpAddress()
	s.Address = myIp + ":8081"
	fmt.Println("[SERVER] Listening on 8081...")

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("[SERVER] accept error:", err)
			continue
		}

		fmt.Println("[SERVER] New connection:", conn.RemoteAddr())
		
		go handleConnection(conn, s)
	}
}

func handleConnection(conn net.Conn, server *Server) {
	defer conn.Close()
	buffer := make([]byte, 1024)

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("[SERVER] read error:", err)
			return
		}

		fmt.Println("[SERVER] Received:", string(buffer[:n]))

		msg := string(buffer[:n])
		if len(msg) >= 12 && msg[:12] == "NEW_MESSAGE " {
			parts := strings.SplitN(msg, " ", 4)
			if len(parts) < 4 {
				fmt.Println("[SERVER] Invalid NEW_MESSAGE format")
				continue
			}

			msgUuid := parts[1]
			originPeer := parts[2]
			message := parts[3]

			server.handleServerMessages(msgUuid, originPeer, message)
		} 
		
		if len(msg) >= 9 && msg[:9] == "NEW_PEER " {
			peerAddr := msg[9:]
			exists := false

			fmt.Println("[SERVER] Registering new peer:", peerAddr)

			for _, peer := range server.PeerList {
				if peer.Address == peerAddr {
					exists = true
					break
				}
			}

			if exists {
				fmt.Println("[SERVER] Peer already in list:", peerAddr)
				continue
			} else {
				server.PeerList = append(server.PeerList, PeerInfo{Address: peerAddr})
				fmt.Println("[SERVER] Updated peer list:", server.PeerList)
			}
		} else {
			fmt.Println("[SERVER] Unknown message format")
		}
	}
}

func bootstrapRegister(ipAddress string) {
	bootstrapAddr := getBootstrapAddress()
	conn, err := net.Dial("tcp", bootstrapAddr)
	if err != nil {
			fmt.Println("[SERVER] Erro ao conectar no bootstrap:", err)
			return
	}

	defer conn.Close()
  fmt.Println("[SERVER] Conectado ao bootstrap")

	_, err = conn.Write([]byte("REGISTER_PEER " + ipAddress + "\n"))
	if err != nil {
			fmt.Println("[SERVER] Erro ao enviar registro:", err)
			return
	}
	fmt.Println("[SERVER] Registro enviado, aguardando resposta...")

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
			fmt.Println("[SERVER] Erro ao ler resposta:", err)
			return
	}

	fmt.Println("[SERVER] Bootstrap respondeu:", string(buf[:n]))
}

func (s *Server) handleServerMessages(msgId string, msgPeerOrigin string, message string) {
	peers := s.PeerList
	if len(peers) == 0 {
		fmt.Println("[SERVER] No peers available to forward the message.")
		return
	}

	// Simple round-robin forwarding
	nextPeer := peers[0]
	peers = append(peers[1:], nextPeer)
	
	if (strings.TrimSpace(nextPeer.Address) == strings.TrimSpace(msgPeerOrigin)) {
		if len(peers) == 1 {
			fmt.Println("[SERVER] Only origin peer available, not forwarding.")
			return
		}
		
		fmt.Println("[SERVER] Skipping origin peer:", nextPeer.Address)
	}

	s.PeerList = peers
	conn, err := connectToPeer(nextPeer.Address)
	if err != nil {
		fmt.Println("[SERVER] Error connecting to peer:", err)
		return
	}
	defer conn.Close()

	err = sendMessage(conn, msgId, msgPeerOrigin, message)
	if err != nil {
		fmt.Println("[SERVER] Error sending message to peer:", err)
		return
	}

	fmt.Println("[SERVER] Message forwarded to peer:", nextPeer.Address)
}
