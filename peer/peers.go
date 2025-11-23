package main

import (
	"fmt"
	"net"
	"sort"
	"strings"
	"time"
)

type PeerInfo struct {
	Address string
	Latency int64
}

func NewPeerInfo() *PeerInfo {
	return &PeerInfo{
		Address: "",
		Latency: 0,
	}
}

func (s *Server) fetchPeerList() ([]PeerInfo, error) {
	conn, err := net.Dial("tcp", "localhost:9001")
	if err != nil {
		fmt.Println("[SERVER] Erro ao conectar no bootstrap:", err)
		return nil, err
	}
	defer conn.Close()
	
	_, err = conn.Write([]byte("GET_PEERS\n"))
	if err != nil {
		fmt.Println("[SERVER] Erro ao solicitar lista de peers:", err)
		return nil, err
	}

	buf := make([]byte, 4096)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("[SERVER] Erro ao ler resposta:", err)
		return nil, err
	}

	response := strings.TrimSpace(string(buf[:n]))
	fmt.Println("[PEER] Bootstrap response:", response)

	if response == "" {
		return []PeerInfo{}, nil
	}

	peers := strings.Split(response, ",")

	for _, peer := range peers {
		if isMyIpInPeerList(peer) {
			continue
		}

		peerInfo := NewPeerInfo()
		latency, err := measureLatencyToPeer(peer)
		if err != nil {
			fmt.Println("[PEER] Error measuring latency to", peer, ":", err)
			continue
		}
		peerInfo.Address = peer
		peerInfo.Latency = latency
		s.PeerList = append(s.PeerList, *peerInfo)
	}

	fmt.Println("[PEER] Peers available:", s.PeerList)
	sort.Slice(s.PeerList, func(i, j int) bool {
		return s.PeerList[i].Latency < s.PeerList[j].Latency
	})
	return s.PeerList, nil
}

func isMyIpInPeerList(peerAdr string) bool {
	myIp, err := getMyIpAddress()
	if err != nil {
		fmt.Println("[CHECK] Error getting my IP address:", err)
		return false
	}
	return peerAdr == myIp + ":8080"
}

func measureLatencyToPeer(peer string) (int64, error) {
	start := time.Now()
	conn, err := net.DialTimeout("tcp", peer, 2*time.Second)
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	latency := time.Since(start).Milliseconds()
	return latency, nil
}


