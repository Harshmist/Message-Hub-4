package functions

import (
	"message-hub/variables"
	"strings"
)

func NewSub(fields []string, AddressChan chan string) {
	var subUser variables.User
	var subSlice [2]interface{}
	roomName := fields[1]
	subUser.Address = AddressChan
	subSlice[0] = roomName
	subSlice[1] = subUser

	variables.SubChannel <- subSlice
}

func NewRoom(fields []string, AddressChan chan string) {
	roomName := fields[1]
	var newRoom [2]interface{}
	newRoom[0] = roomName
	newRoom[1] = AddressChan
	variables.NewChannel <- newRoom
}

func NewPub(fields []string, AddressChan chan string) {
	var msgSlice [2]interface{}
	room := fields[1]
	message := strings.Join(fields[2:], " ")
	msgSlice[0] = room
	msgSlice[1] = message
	variables.PubChannel <- msgSlice
}
