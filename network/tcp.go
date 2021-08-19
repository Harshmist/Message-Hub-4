package network

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"message-hub/functions"
	"message-hub/variables"
	"net"
	"strings"
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
func RequestHandler(conn net.Conn, AddressChan chan string) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) < 1 {
			continue
		}
		switch fields[0] {

		case "SUB":
			functions.NewSub(fields, AddressChan)

		case "NEW":
			functions.NewRoom(fields, AddressChan)
		case "PUB":
			functions.NewPub(fields, AddressChan)
		}
	}
}
func Broadcaster(conn net.Conn, AddressChan chan string) {

	for {
		select {
		case msg := <-AddressChan:
			message := fmt.Sprintf("%v\n", msg)
			io.WriteString(conn, message)
		}
	}
}
