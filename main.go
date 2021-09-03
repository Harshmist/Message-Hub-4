package main

import (
	_ "expvar"
	"fmt"
	"message-hub/hub"
	"message-hub/monitoring"
	"message-hub/network"
	"net/http"
	_ "net/http/pprof"
	"time"
)

var startTime time.Time

func init() {
	startTime = time.Now()
}

func main() {
	go monitoring.TimeMonitoring(startTime)
	go network.TcpListener()
	go hub.ChannelRouter()
	go network.StartUDP()
	// http handler for the pprof profiling
	go func() {
		http.ListenAndServe(":8004", nil)
	}()
	fmt.Println(time.Now())

	fmt.Scanln()

}
