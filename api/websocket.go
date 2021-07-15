package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
    "github.com/gorilla/websocket"
	// "net/http"
	"strconv"
	"strings"
	"../model"
)

const (
	MSG_NAME_CREATE = "__create"
	MSG_NAME_INVITE = "__invite"
	MSG_NAME_RING   = "__ring"
	MSG_NAME_CANCEL = "__cancel"
	MSG_NAME_REJECT = "__reject"
	MSG_NAME_JOIN   = "__join"
	MSG_NAME_ICECAN = "__ice_candidate"
	MSG_NAME_OFFER  = "__offer"
	MSG_NAME_ANSWER = "__answer"
	MSG_NAME_LEAVE  = "__leave"
	MSG_NAME_AUDIO  = "__audio"
	MSG_NAME_DISCONN= "__disconnect"
	MSG_NAME_PEERS  = "__peers"
	MSG_NAME_NEWPEER= "__new_peer"
)

var (
	websocketUpgrader = websocket.Upgrader{}
)

type EventMessage struct {
	Name string                 `json:"eventName"`
	Data map[string]interface{} `json:"data"`
}

func handleMessage(user *model.User, conn *websocket.Conn, userID string) {
	defer func() {
		fmt.Println("user logout：", userID)
		coreMgr.DataMgr.UserDelete(userID)
	}()

	for {
    	msg := EventMessage{}
		err := conn.ReadJSON(&msg)
		if err != nil {
			errStr := err.Error()
			if strings.HasPrefix(errStr, "websocket: close") {
				fmt.Println(errStr)
				break
			}
			showError(err)
			continue
		}
		fmt.Println(msg.Name)
		switch msg.Name {
		case MSG_NAME_CREATE: {
			createRoom(msg.Data)
		}
		case MSG_NAME_INVITE: {
			invite(msg.Data)
		}
		case MSG_NAME_RING: {
			ring(msg.Data)
		}
		case MSG_NAME_CANCEL: {
			cancel(msg.Data)
		}
		case MSG_NAME_REJECT: {
			reject(msg.Data)
		}
		case MSG_NAME_JOIN: {
			join(msg.Data)
		}
		case MSG_NAME_ICECAN: {
			iceCandidate(msg.Data)
		}
		case MSG_NAME_OFFER: {
			offer(msg.Data)
		}
		case MSG_NAME_ANSWER: {
			answer(msg.Data)
		}
		case MSG_NAME_LEAVE: {
			leave(msg.Data)
		}
		case MSG_NAME_AUDIO: {
			transAudio(msg.Data)
		}
		case MSG_NAME_DISCONN: {
			disconnet(msg.Data)
			return
		}
		}
	}
}

func sendMessage(conn *websocket.Conn, msgName string, data map[string]interface{}) error {
	msg := EventMessage{}
	msg.Name = msgName
	msg.Data = data
	return conn.WriteJSON(&msg)
}

func sendMessageToUser(userID, msgName string, data map[string]interface{}) error {
	user := coreMgr.DataMgr.UserGet(userID)
	if user == nil {
		return model.Err_User_Not_Found
	}

	msg := EventMessage{}
	msg.Name = msgName
	msg.Data = data
	return user.WriteJSON(&msg)
}

func createRoom(data map[string]interface{}) {
	roomID := data["room"].(string)
	userID := data["userID"].(string)

	room := coreMgr.DataMgr.RoomGet(roomID)
	if room != nil {
		return
	}
	roomSize := (int)(data["roomSize"].(float64))
	room = model.NewRoom(roomID, userID, roomSize)
	coreMgr.DataMgr.RoomSet(roomID, room)

	user := coreMgr.DataMgr.UserGet(userID)
	if user != nil {
		room.AddUser(user)

		sendData := map[string]interface{}{}
		sendData["connections"] = ""
		sendData["you"] = userID
		sendData["roomSize"] = roomSize

		sendMessageToUser(userID, MSG_NAME_PEERS, sendData)
	}
}

func invite(data map[string]interface{}) {
	userList  := data["userList"].(string)
	inviteID  := data["inviteID"].(string)
	roomID    := data["room"].(string)
	audioFlag := data["audioOnly"]
	userIDs   := strings.Split(userList, ",")

	audioOnly := false
	if audioFlag != nil {
		audioOnly = true
	}

	fmt.Printf("invite: Room: %s, User %s invite %s, audioOnly: %v\n", roomID, inviteID, userList, audioOnly)
	for _, userID := range userIDs {
		sendMessageToUser(userID, MSG_NAME_INVITE, data)
	}
}

func ring(data map[string]interface{}) {
	roomID   := data["room"].(string)
	inviteID := data["toID"].(string)
	fmt.Printf("ring: invite user: %s to room: %s\n", inviteID, roomID)

	sendMessageToUser(inviteID, MSG_NAME_RING, data)
}

func cancel(data map[string]interface{}) {
	roomID   := data["room"].(string)
	userList := data["userList"].(string)
	userIDs  := strings.Split(userList, ",")

	fmt.Printf("cancel: Room: %s, invite %s\n", roomID, userList)
	for _, userID := range userIDs {
		sendMessageToUser(userID, MSG_NAME_CANCEL, data)
	}

	coreMgr.DataMgr.RoomDelete(roomID)
}

func reject(data map[string]interface{}) {
	roomID := data["room"].(string)
	toID   := data["toID"].(string)

	sendMessageToUser(toID, MSG_NAME_REJECT, data)
	// TODO 多人时不处理
	coreMgr.DataMgr.RoomDelete(roomID)
}

func join(data map[string]interface{}) {
	roomID := data["room"].(string)
	userID := data["userID"].(string)

	room := coreMgr.DataMgr.RoomGet(roomID)
	if room == nil {
		return
	}
	if room.GetMaxSize() == room.GetCurSize() {
		return
	}

	me := coreMgr.DataMgr.UserGet(userID)
	if me == nil {
		return
	}
	room.AddUser(me)

	userIDs := room.GetUserIDs()
	strUserIDs := ""
	for _, userID := range userIDs {
		strUserIDs += userID + ","
	}
	if len(strUserIDs) > 1 {
		strUserIDs = strUserIDs[:len(strUserIDs) - 1]
	}

	sendData := map[string]interface{}{}
	sendData["connections"] = strUserIDs
	sendData["you"] = userID
	sendData["roomSize"] = room.GetMaxSize()

	sendMessageToUser(userID, MSG_NAME_PEERS, sendData)

	peerData := map[string]interface{}{}
	peerData["userID"] = userID
	for _, userID := range userIDs {
		sendMessageToUser(userID, MSG_NAME_NEWPEER, peerData)
	}
}

func iceCandidate(data map[string]interface{}) {
	userID := data["userID"].(string)
	sendMessageToUser(userID, MSG_NAME_ICECAN, data)
}

func offer(data map[string]interface{}) {
	userID := data["userID"].(string)
	sendMessageToUser(userID, MSG_NAME_OFFER, data)
}

func answer(data map[string]interface{}) {
	userID := data["userID"].(string)
	sendMessageToUser(userID, MSG_NAME_ANSWER, data)
}

func leave(data map[string]interface{}) {
	roomID := data["room"].(string)
	fromID := data["fromID"].(string)

	room := coreMgr.DataMgr.RoomGet(roomID)
	if room == nil {
		return
	}
	room.RemoveUser(fromID)
	users := room.GetUsers()
	for _, user := range users {
		userID := user.GetUserID()
		if userID == fromID {
			continue
		}
		sendMessageToUser(user.GetUserID(), MSG_NAME_LEAVE, data)
	}
	if room.GetCurSize() < 2 {
		coreMgr.DataMgr.RoomDelete(roomID)
	}
}

func transAudio(data map[string]interface{}) {
	userID := data["userID"].(string)
	sendMessageToUser(userID, MSG_NAME_AUDIO, data)
}

func disconnet(data map[string]interface{}) {
	userID := data["userID"].(string)
	sendMessageToUser(userID, MSG_NAME_DISCONN, data)
}

func showError(err error) {
    if err != nil {
        fmt.Println(err)
    }
}

//==== API ====
func WS(c *gin.Context) {
	userID := c.Param("userid")
	device, _ := strconv.Atoi(c.Param("device"))
	devStr  := "Phone"
	if device != 0 {
		devStr = "PC"
	}

	conn, err := websocketUpgrader.Upgrade(c.Writer, c.Request, nil)
	showError(err)
	user := coreMgr.DataMgr.UserGet(userID)
	if user == nil {
		user = model.NewUser(userID, "", device)
		coreMgr.DataMgr.UserSet(userID, user)
	} else {
		user.SetDevice(device)
	}
	user.SetConn(conn)
	user.SetStatus(model.DEVICE_STATUS_IDLE)

	fmt.Printf("%s user login：%s\n", devStr, userID)

	data := map[string]interface{}{}
	data["userID"] = userID
	data["avatar"] = ""
	sendMessage(conn, "__login_success", data)

	go handleMessage(user, conn, userID)
}
