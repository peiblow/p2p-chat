package main

import (
	"fmt"
	"net"
	"os"
)

type Server struct {
	PeerList []PeerInfo
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

func startServer() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	fmt.Println("[SERVER] Listening on 8080...")

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("[SERVER] accept error:", err)
			continue
		}

		fmt.Println("[SERVER] New connection:", conn.RemoteAddr())
		
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	buffer := make([]byte, 1024)

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("[SERVER] read error:", err)
			return
		}

		fmt.Println("[SERVER] Received:", string(buffer[:n]))
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

	_, err = conn.Write([]byte("REGISTER_PEER " + ipAddress + ":8080\n"))
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
