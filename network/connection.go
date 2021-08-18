package network

import (
	"io"
	"log"
	"net"
)

var (
	JoinChan = make(chan chan string)
)

type User struct {
	Name    string
	Address chan string
}

func TcpListener() {
	listener, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		io.WriteString(conn, "Welcome to the message hub!\n Write [CMD] for a list of commands\n")
		var AddressChan = make(chan string)

		go RequestHandler(conn, AddressChan)
		go Broadcaster(conn, AddressChan)
		JoinChan <- AddressChan

	}
}
