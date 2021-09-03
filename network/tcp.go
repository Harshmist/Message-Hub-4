package network

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"message-hub/hub"
	"message-hub/monitoring"
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
		scanner := bufio.NewScanner(conn)
		io.WriteString(conn, "Welcome to the message hub!\n")
		io.WriteString(conn, "Server message: What is your name?\n")
		scanner.Scan()
		var AddressChan = make(chan string)
		var user hub.User
		user.Name = scanner.Text()
		user.Address = AddressChan

		go RequestHandler(conn, user)
		go Broadcaster(conn, user)
		hub.JoinChan <- user

	}
}
func RequestHandler(conn net.Conn, user hub.User) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) < 1 {
			monitoring.FailedRequest()
			io.WriteString(conn, "Please use a command\n")
			continue
		}
		switch fields[0] {
		case "LIST":
			hub.List(user)
			monitoring.TotalRequests.Add(1)
		case "SUB":
			hub.NewSub(fields, user)

		case "NEW":
			hub.NewRoom(fields, user)
		case "PUB":
			hub.NewPub(fields, user)

		case "STOP":
			user.Address <- "Closing connection"
			io.WriteString(conn, "STOP")
			conn.Close()
			monitoring.TotalRequests.Add(1)

		default:
			monitoring.FailedRequest()
			io.WriteString(conn, "Command not recognised\n")
		}
	}
}
func Broadcaster(conn net.Conn, user hub.User) {

	for {
		select {
		case msg := <-user.Address:
			message := fmt.Sprintf("%v\n", msg)
			io.WriteString(conn, message)
		}
	}
}
