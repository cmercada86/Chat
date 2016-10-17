package route

import (
	"Chat/controller"
	"net/http"

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
		Pattern:     "/connect",
		HandlerFunc: controller.Connect,
	},
	Route{
		Name:        "chat",
		Method:      "POST",
		Pattern:     "/chat",
		HandlerFunc: controller.AddMessage,
	},
	Route{
		Name:        "chat",
		Method:      "GET",
		Pattern:     "/chat",
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

	return router
}
