package main

import (
	"fmt"
	"message-hub/network"
)

var (
	allUsers []network.User
	rooms    = make(map[string][]network.User, 10)
)

func main() {
	go network.TcpListener()
	go func() {

		for {
			select {
			case newUser := <-network.JoinChan:
				var user network.User
				user.Address = newUser
				allUsers = append(allUsers, user)
			case newRoom := <-network.NewChannel:
				var creatingUser network.User
				creatingUser.Address = newRoom[1].(chan string)
				newRoomName := newRoom[0].(string)
				rooms[newRoomName] = make([]network.User, 0, 10)
				rooms[newRoomName] = append(rooms[newRoomName], creatingUser)
			case newSub := <-network.SubChannel:
				roomName := newSub[0].(string)
				subUser := newSub[1].(network.User)
				rooms[roomName] = append(rooms[roomName], subUser)
			case newMsg := <-network.PubChannel:
				roomName := newMsg[0].(string)
				message := newMsg[1].(string)

				for _, v := range rooms[roomName] {
					v.Address <- message
				}
			}
		}

	}()

	fmt.Scanln()

}
