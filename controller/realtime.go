package controller

import (
	"Chat/auth"
	"Chat/model"
	"log"
	"net/http"
	"time"

	"golang.org/x/net/websocket"
)

func InitRealTime(ws *websocket.Conn) {
	state := model.UUID(r.FormValue("state"))
	room := r.FormValue("room")

	user, isAuth := auth.CheckAuth(state)
	if !isAuth {
		log.Println("Not authorized!")
		return
	}

	listener := model.NewListener()
	closeChan := w.(http.CloseNotifier).CloseNotify()
	//Register listener to db listener

	for i := 0; i < 1440; {
		select {
		case chat <- listener.ChatChannel:
			//send chat
			if err := websocket.Message.Send(ws, model.ChatToJsonString(chat)); err != nil {
				//remove from db listener
				i = 1440
			}
		case dm <- listener.DMchannel:
			//send dm
			if err := websocket.Message.Send(ws, model.DirectMessageToJsonSring(dm)); err != nil {
				//remove from db listener
				i = 1440
			}
		case <-time.After(time.Second * 60):

			f.Flush()
			i++
		case <-closeChan:
			f.Flush()
			i = 1440
			//remove from db listener
		}
	}
}
