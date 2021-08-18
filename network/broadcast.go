package network

import (
	"fmt"
	"io"
	"net"
)

func Broadcaster(conn net.Conn, AddressChan chan string) {

	for {
		select {
		case msg := <-AddressChan:
			message := fmt.Sprintf("%v\n", msg)
			io.WriteString(conn, message)
		}
	}
}
