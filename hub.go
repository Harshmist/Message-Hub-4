package main

import (
	"message-hub/hub"
	"message-hub/network"
)

var (
	allUsers []hub.User
	rooms    = make(map[string][]hub.User, 10)
)

func main() {
	go network.TcpListener()
	go func() {
		for {
			select {
			case newUser := <-hub.JoinChan:
				allUsers = append(allUsers, newUser)
			case newRoom := <-hub.NewChannel:
				var creatingUser hub.User
				creatingUser.Address = newRoom[1].(chan string)
				newRoomName := newRoom[0].(string)
				rooms[newRoomName] = make([]hub.User, 0, 10)
				rooms[newRoomName] = append(rooms[newRoomName], creatingUser)
			case newSub := <-hub.SubChannel:
				roomName := newSub[0].(string)
				subUser := newSub[1].(hub.User)
				rooms[roomName] = append(rooms[roomName], subUser)
			case newMsg := <-hub.PubChannel:
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
