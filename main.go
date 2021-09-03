package main

import (
	_ "expvar"
	"fmt"
	"message-hub/hub"
	"message-hub/network"
	"net/http"
	_ "net/http/pprof"
)

func main() {
	go network.TcpListener()
	go hub.ChannelRouter()
	go network.StartUDP()
	go func() {
		http.ListenAndServe(":8004", nil)
	}()

	fmt.Scanln()

}
