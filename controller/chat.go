package controller

import (
	"Chat/cache"
	"Chat/model"
	"Chat/repository"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

func AddMessage(w http.ResponseWriter, r *http.Request) {
	x, _ := ioutil.ReadAll(r.Body)

	state := model.UUID(r.FormValue("state"))

	if !cache.Contains(state) {
		log.Println("invalid state!")
		return
	}
	user := cache.Get(state)
	room := r.FormValue("room")

	message := string(x)

	repository.AddChatMessage(room, user.ID, message)
}

func GetMessages(w http.ResponseWriter, r *http.Request) {

	state := model.UUID(r.FormValue("state"))

	if !cache.Contains(state) {
		log.Println("invalid state!")
		return
	}

	room := r.FormValue("room")
	chats, err := repository.GetChatMessages(room)
	if err != nil {
		log.Println("Error retrieving chats: ", err)
		return
	}

	json.NewEncoder(w).Encode(chats)
}

func GetRoomNames(w http.ResponseWriter, r *http.Request) {
	state := model.UUID(r.FormValue("state"))

	if !cache.Contains(state) {
		log.Println("invalid state!")
		return
	}

	rooms, err := repository.GetRoomNames()
	if err != nil {
		log.Println("Error retrieving chats: ", err)
		return
	}

	json.NewEncoder(w).Encode(rooms)
}
