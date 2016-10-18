package controller

import (
	"Chat/auth"
	"Chat/model"
	"Chat/repository"
	"log"
	"time"

	"github.com/gorilla/mux"

	"golang.org/x/net/websocket"
)

type wsMessage struct {
	Type    string      `json:"type"`
	Message interface{} `json:"message"`
}

var room string

func InitRealTime(ws *websocket.Conn) {
	vars := mux.Vars(ws.Request())
	state := model.UUID(vars["state"])
	room = vars["room"]
	log.Println(state, room)
	user, isAuth := auth.CheckAuth(state)
	if !isAuth {
		log.Println("Web socket:Not authorized!")
		return
	}

	listener := model.NewListener(user, room)

	repository.AddListener(listener)

	users, err := repository.GetCurrentListeners()
	if err != nil {
		log.Println("Error getting current users: ", err)
	} else {
		for _, curUser := range users {
			//dont need to send yourself!
			if curUser.ID != user.ID {
				websocket.Message.Send(ws, model.ObjectToJsonString(wsMessage{
					Type:    "user",
					Message: &curUser,
				}))
			}
		}
	}
	//Register listener to db listener

	for i := 0; i < 1440; {
		select {
		case chat := <-listener.ChatChannel:
			//send chat

			if err := websocket.Message.Send(ws, model.ObjectToJsonString(wsMessage{
				Type:    "chat",
				Message: &chat,
			})); err != nil {
				//remove from db listener
				i = 1440
			}

		case dm := <-listener.DMchannel:
			//send dm
			if err := websocket.Message.Send(ws, model.ObjectToJsonString(wsMessage{
				Type:    "dm",
				Message: &dm,
			})); err != nil {
				//remove from db listener
				i = 1440
			}
		//case user := <-listener.UserChannel:
		case room := <-listener.RoomChannel:
			if err := websocket.Message.Send(ws, model.ObjectToJsonString(wsMessage{
				Type:    "room",
				Message: &room,
			})); err != nil {
				//remove from db listener
				i = 1440
			}
		case <-time.After(time.Second * 60):

			i++
			// case <-closeChan:

			// 	i = 1440
			// 	//remove from db listener
		}
	}
}
