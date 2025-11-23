package main

func main() {
	server := NewServer()
	go startServer()
	
	myIp, _ := getMyIpAddress()
	bootstrapRegister(myIp)
	
	server.fetchPeerList()

	if len(server.PeerList) > 0{
		go startClient(server.PeerList[0].Address)
	} else {
		go startClient("")
	}

	select {}
}