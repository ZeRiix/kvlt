package routes

import (
	"github.com/gorilla/mux"
)

type RouteHandler func(*mux.Router)

func SetupRouter() *mux.Router {
	router := mux.NewRouter()

	routeHandlers := []RouteHandler{
		getValueRoute,
		setValueRoute,
	}

	for _, registerRoute := range routeHandlers {
		registerRoute(router)
	}

	return router
}
