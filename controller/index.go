package controller

import (
	"Chat/auth"
	"Chat/cache"
	"Chat/config"
	"Chat/model"
	"Chat/template"
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request) {

	state, _ := auth.NewUUID()

	cache.Set(state, model.User{})
	cache.Expire(state)

	template.Execute(w, "index.html", auth.ClientInfo{
		ClientID: config.GetConfig().GogClientID,
		State:    state,
	})

}
