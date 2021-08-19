package network

import (
	"bufio"
	"message-hub/functions"
	"net"
	"strings"
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
			functions.NewSub(fields, AddressChan)

		case "NEW":
			functions.NewRoom(fields, AddressChan)
		case "PUB":
			functions.NewPub(fields, AddressChan)
		}
	}
}
