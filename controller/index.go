package controller

import (
	"Chat/auth"
	"Chat/config"
	"Chat/template"
	"net/http"
)

func index(w http.ResponseWriter, r *http.Request) {

	template.Execute(w, "index.html", auth.ClientInfo{
		ClientID: config.GetConfig().GogClientID,
		State:    "TEST2",
	})

}
