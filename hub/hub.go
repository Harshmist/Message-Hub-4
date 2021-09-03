package hub

import (
	"fmt"
	"message-hub/monitoring"
	"net"
	"strings"
)

var (
	SubChannel   = make(chan [2]interface{}, 2)
	NewChannel   = make(chan [2]interface{}, 2)
	PubChannel   = make(chan [2]interface{}, 2)
	JoinChan     = make(chan User)
	UDPChan      = make(chan [2]interface{}, 2)
	requestsChan = make(chan bool)
	allUsers     []User
	rooms        = make(map[string][]User, 10)
)

type User struct {
	Name    string
	Address chan string
	UDPAdd  *net.UDPAddr
}

// process user submitted new room
func NewSub(fields []string, user User) {
	var subSlice [2]interface{}
	roomName := fields[1]
	subSlice[0] = roomName
	subSlice[1] = user

	SubChannel <- subSlice
}

// creates new room
func NewRoom(fields []string, user User) {
	roomName := fields[1]
	var newRoom [2]interface{}
	newRoom[0] = roomName
	newRoom[1] = user
	NewChannel <- newRoom

}

// process sent by user, message and room to send to
func NewPub(fields []string, user User) {
	var msgSlice [2]interface{}
	var userName string
	room := fields[1]
	if user.Name != "" {
		userName = user.Name
	} else {
		userName = user.UDPAdd.String()
	}

	message := fmt.Sprintf("%v wrote in room %v: %v", userName, room, strings.Join(fields[2:], " "))
	msgSlice[0] = room
	msgSlice[1] = message
	PubChannel <- msgSlice
}

// create list of rooms available as string and send back to user
func List(user User) {
	var newString string
	var UDPSend [2]interface{}
	for k := range rooms {
		newString = fmt.Sprintf("%v %v |", newString, k)
	}
	sendString := "Available categories are: " + newString
	if user.Address != nil {
		user.Address <- sendString
	} else {
		UDPSend[0] = user.UDPAdd
		UDPSend[1] = sendString
		UDPChan <- UDPSend
	}
}

// logic for receiving on the various channels
func ChannelRouter() {
	for {
		select {
		case newUser := <-JoinChan:
			allUsers = append(allUsers, newUser)
			fmt.Println("User connected")
			monitoring.TotalUsers.Add(1)
		case newRoom := <-NewChannel:
			var creatingUser User
			creatingUser = newRoom[1].(User)
			newRoomName := newRoom[0].(string)
			rooms[newRoomName] = make([]User, 0, 10)
			rooms[newRoomName] = append(rooms[newRoomName], creatingUser)
			fmt.Println("new room created")
		case newSub := <-SubChannel:
			roomName := newSub[0].(string)
			subUser := newSub[1].(User)
			rooms[roomName] = append(rooms[roomName], subUser)
		case newMsg := <-PubChannel:
			roomName := newMsg[0].(string)
			message := newMsg[1].(string)

			for _, v := range rooms[roomName] {
				if v.Address != nil {
					v.Address <- message
				} else {
					var UDPMsg [2]interface{}
					UDPMsg[0] = v.UDPAdd
					UDPMsg[1] = message
					UDPChan <- UDPMsg
				}
			}
		case completed := <-requestsChan:
			if completed {
				monitoring.TotalRequests.Add(1)
			} else {
				monitoring.TotalRequests.Add(1)
				monitoring.FailedRequests.Add(1)
			}
		}
	}
}
