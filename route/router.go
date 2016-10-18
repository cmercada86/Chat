package route

import (
	"Chat/controller"
	"net/http"

	"golang.org/x/net/websocket"

	"github.com/gorilla/mux"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

var routes = []Route{
	Route{
		Name:        "index",
		Method:      "GET",
		Pattern:     "/",
		HandlerFunc: controller.Index,
	},
	Route{
		Name:        "connect",
		Method:      "POST",
		Pattern:     "/connect/{state}",
		HandlerFunc: controller.Connect,
	},
	Route{
		Name:        "chat",
		Method:      "POST",
		Pattern:     "/chat/{room}/{state}",
		HandlerFunc: controller.AddMessage,
	},
	Route{
		Name:        "chat",
		Method:      "GET",
		Pattern:     "/chat/{room}/{state}",
		HandlerFunc: controller.GetMessages,
	},
}

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	for _, route := range routes {

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}

	router.Handle("/chat/ws/{room}/{state}", websocket.Handler(controller.InitRealTime))
	//... ListenAndServe, etc)

	return router
}
