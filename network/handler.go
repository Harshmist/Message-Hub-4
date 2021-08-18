package network

import (
	"bufio"
	"net"
	"strings"
)

var (
	SubChannel = make(chan [2]interface{}, 2)
	NewChannel = make(chan [2]interface{}, 2)
	PubChannel = make(chan [2]interface{}, 2)
)

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
			var subUser User
			var subSlice [2]interface{}
			roomName := fields[1]
			subUser.Address = AddressChan
			subSlice[0] = roomName
			subSlice[1] = subUser
			SubChannel <- subSlice

		case "NEW":
			roomName := fields[1]
			var newRoom [2]interface{}
			newRoom[0] = roomName
			newRoom[1] = AddressChan
			NewChannel <- newRoom

		case "PUB":
			var msgSlice [2]interface{}
			room := fields[1]
			message := strings.Join(fields[2:], " ")
			msgSlice[0] = room
			msgSlice[1] = message
			PubChannel <- msgSlice
		}
	}
}
