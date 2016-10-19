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
		Name:        "getchats",
		Method:      "GET",
		Pattern:     "/chat/{room}/{state}",
		HandlerFunc: controller.GetMessages,
	},
	Route{
		Name:        "rooms",
		Method:      "GET",
		Pattern:     "/rooms/{state}",
		HandlerFunc: controller.GetRoomNames,
	},
	Route{
		Name:        "dm",
		Method:      "POST",
		Pattern:     "/dm/{receiver_id}/{state}",
		HandlerFunc: controller.SendDirectMessage,
	},
}

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	router.Handle("/chat/ws/{state}", websocket.Handler(controller.InitRealTime))

	for _, route := range routes {

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}

	//... ListenAndServe, etc)

	return router
}
