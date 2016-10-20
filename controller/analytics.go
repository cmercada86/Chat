package controller

import (
	"Chat/auth"
	"Chat/model"
	"io/ioutil"
	"log"
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

	log.Println(room, searchString)
	//Call lambda python script to search chats

}
