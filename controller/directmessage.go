package controller

import (
	"Chat/cache"
	"Chat/model"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

func SendDirectMessage(w http.ResponseWriter, r *http.Request) {
	x, _ := ioutil.ReadAll(r.Body)

	state := model.UUID(r.FormValue("state"))

	if !cache.Contains(state) {
		log.Println("invalid state!")
		return
	}

	user := cache.Get(state)

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
	state := model.UUID(r.FormValue("state"))

	if !cache.Contains(state) {
		log.Println("invalid state!")
		return
	}

	user := cache.Get(state)

	dms, err := repository.GetDirectMessages(user.ID)
	if err != nil {
		log.Println("Error retrieving dms: ", err)
		return
	}

	json.NewEncoder(w).Encode(dms)

}