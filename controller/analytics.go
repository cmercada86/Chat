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

func SearchChats(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	x, _ := ioutil.ReadAll(r.Body)
	state := model.UUID(vars["state"])
	room := vars["room"]

	_, isAuth := auth.CheckAuth(state)
	if !isAuth {
		errorHandler(w, r, 403, "Not Authorized!")
		return
	}
	searchString := string(x)
	chats, err := repository.SearchChat(room, searchString)
	if err != nil {

		errorHandler(w, r, 500, "")
		return
	}

	//Call lambda python script to search chats
	json.NewEncoder(w).Encode(chats)
}
