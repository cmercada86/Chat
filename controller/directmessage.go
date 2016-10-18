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

func SendDirectMessage(w http.ResponseWriter, r *http.Request) {
	x, _ := ioutil.ReadAll(r.Body)

	vars := mux.Vars(r)
	state := model.UUID(vars["state"])

	user, isAuth := auth.CheckAuth(state)
	if !isAuth {
		log.Println("Not authorized!")
		return
	}

	var dm model.DirectMessage

	if err := json.Unmarshal(x, &dm); err != nil {
		log.Println("Unable to unmarshal direct message: ", err)
	}

	if user.ID != dm.SenderID {
		//
	}

	repository.InsertDirectMessage(dm)
}

func ReceiveDirectMessages(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	state := model.UUID(vars["state"])

	user, isAuth := auth.CheckAuth(state)
	if !isAuth {
		log.Println("Not authorized!")
		return
	}

	dms, err := repository.GetDirectMessages(user.ID)
	if err != nil {
		log.Println("Error retrieving dms: ", err)
		return
	}

	json.NewEncoder(w).Encode(dms)

}
