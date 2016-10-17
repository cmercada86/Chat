package controller

import (
	"Chat/auth"
	"Chat/cache"
	"Chat/model"
	"Chat/repository"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

func Connect(w http.ResponseWriter, r *http.Request) {

	x, _ := ioutil.ReadAll(r.Body)

	state := model.UUID(r.FormValue("state"))

	if !cache.Contains(state) {
		log.Println("invalid state!")
		return
	}
	//get auth code	string
	authCode := string(x)

	user, err := auth.ExchangeAuthCodeForUser(authCode)
	if err != nil {

	}

	cache.Set(state, user)

	repository.AddOrUpdateUserInfo(user)

	json.NewEncoder(w).Encode(user)

}
