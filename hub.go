package main

import (
	"message-hub/network"
	"message-hub/variables"
)

var (
	allUsers []variables.User
	rooms    = make(map[string][]variables.User, 10)
)

func main() {
	go network.TcpListener()
	go func() {
		for {
			select {
			case newUser := <-variables.JoinChan:
				allUsers = append(allUsers, newUser)
			case newRoom := <-variables.NewChannel:
				var creatingUser variables.User
				creatingUser.Address = newRoom[1].(chan string)
				newRoomName := newRoom[0].(string)
				rooms[newRoomName] = make([]variables.User, 0, 10)
				rooms[newRoomName] = append(rooms[newRoomName], creatingUser)
			case newSub := <-variables.SubChannel:
				roomName := newSub[0].(string)
				subUser := newSub[1].(variables.User)
				rooms[roomName] = append(rooms[roomName], subUser)
			case newMsg := <-variables.PubChannel:
				roomName := newMsg[0].(string)
				message := newMsg[1].(string)

				for _, v := range rooms[roomName] {
					v.Address <- message
				}
			}
		}
	}()

	select {}

}
