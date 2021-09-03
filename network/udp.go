package network

import (
	"fmt"
	"message-hub/hub"
	"net"
	"strings"
)

var users = make(map[string]string)

func StartUDP() {
	var user hub.User
	service := "localhost:8002"

	s, err := net.ResolveUDPAddr("udp4", service)
	if err != nil {
		fmt.Println(err)
		return
	}

	connection, err := net.ListenUDP("udp4", s)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer connection.Close()
	buffer := make([]byte, 1024)
	go UdpBroadcast(connection)
	for {

		n, addr, err := connection.ReadFromUDP(buffer)
		input := strings.TrimSpace(string(buffer[:n]))
		fields := strings.Split(string(input), " ")
		user.Name = users[addr.String()]
		user.UDPAdd = addr

		if input == "STOP" {
			connection.WriteToUDP([]byte("Closing connection...\n"), addr)
			connection.Close()
		}
		switch fields[0] {
		case "LIST":
			hub.List(user)

		case "SUB":
			connection.WriteToUDP([]byte(fmt.Sprintf("Subscribed to category %v\n", fields[1])), addr)
			hub.NewSub(fields, user)
		case "NEW":
			hub.NewRoom(fields, user)
		case "PUB":
			hub.NewPub(fields, user)
		case "NICK":
			addrString := addr.String()
			users[addrString] = fields[1]
		}

		if err != nil {
			fmt.Print(err)
		}
	}
}

func UdpBroadcast(connection *net.UDPConn) {
	for {
		select {
		case msg := <-hub.UDPChan:
			addr := msg[0].(*net.UDPAddr)
			message := msg[1].(string)
			connection.WriteToUDP([]byte(fmt.Sprintf("%v\n", message)), addr)
		}
	}
}
