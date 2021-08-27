package main

import (
	"fmt"
	"message-hub/hub"
	"message-hub/network"
)

func main() {
	go network.TcpListener()
	go hub.ChannelRouter()
	go network.StartUDP()

	fmt.Scanln()

}
