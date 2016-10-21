package controller

import (
	"fmt"
	"net/http"
)

func errorHandler(w http.ResponseWriter, r *http.Request, status int, message string) {
	w.WriteHeader(status)
	if status == http.StatusNotFound {
		fmt.Fprint(w, message)
	}
}
