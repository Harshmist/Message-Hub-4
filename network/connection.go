package network

import (
	"io"
	"log"
	"message-hub/variables"
	"net"
)

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
		io.WriteString(conn, "Welcome to the message hub!\n")
		var AddressChan = make(chan string)
		var user variables.User
		user.Address = AddressChan

		go RequestHandler(conn, AddressChan)
		go Broadcaster(conn, AddressChan)
		variables.JoinChan <- user

	}
}
