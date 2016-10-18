package controller

import (
	"Chat/auth"
	"Chat/model"
	"Chat/repository"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func AddMessage(w http.ResponseWriter, r *http.Request) {
	x, _ := ioutil.ReadAll(r.Body)

	vars := mux.Vars(r)
	state := model.UUID(vars["state"])
	room := vars["room"]

	user, isAuth := auth.CheckAuth(state)
	if !isAuth {
		log.Println("Add Message: Not authorized!")
		return
	}

	message := string(x)

	repository.AddChatMessage(room, user.ID, message)
}

func GetMessages(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	state := model.UUID(vars["state"])
	room := vars["room"]

	_, isAuth := auth.CheckAuth(state)
	if !isAuth {
		log.Println("Get Messages: Not authorized!")
		return
	}

	chats, err := repository.GetChatMessages(room)
	if err != nil {
		log.Println("Error retrieving chats: ", err)
		return
	}

	json.NewEncoder(w).Encode(chats)
}

func GetRoomNames(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	state := model.UUID(vars["state"])

	_, isAuth := auth.CheckAuth(state)
	if !isAuth {
		log.Println("Get rooms: Not authorized!")
		return
	}

	rooms, err := repository.GetRoomNames()
	if err != nil {
		log.Println("Error retrieving chats: ", err)
		return
	}

	json.NewEncoder(w).Encode(rooms)
}
