package controller

import (
	"Chat/auth"
	"Chat/model"
	"Chat/repository"
	"encoding/json"
	"io/ioutil"
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
		errorHandler(w, r, 403, "Not Authorized!")
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
		errorHandler(w, r, 403, "Not Authorized!")
		return
	}

	chats, err := repository.GetChatMessages(room)
	if err != nil {
		errorHandler(w, r, 500, "")
		return
	}

	json.NewEncoder(w).Encode(chats)
}

func GetRoomNames(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	state := model.UUID(vars["state"])

	_, isAuth := auth.CheckAuth(state)
	if !isAuth {
		errorHandler(w, r, 403, "Not Authorized!")
		return
	}

	rooms, err := repository.GetRoomNames()
	if err != nil {
		errorHandler(w, r, 500, "")
		return
	}

	json.NewEncoder(w).Encode(rooms)
}
