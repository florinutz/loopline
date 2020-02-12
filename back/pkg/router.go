package pkg

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (c *Controller) generateRouter(middlewares ...mux.MiddlewareFunc) *mux.Router {
	router := mux.NewRouter()

	sub := router.PathPrefix("/").Headers("Content-Type", "application/json").Subrouter()

	for _, mid := range middlewares {
		router.Use(mid)
	}

	sub.HandleFunc("/", c.List).Methods(http.MethodGet)
	sub.HandleFunc("/", c.Create).Methods(http.MethodPost)
	sub.HandleFunc("/", c.Delete).Methods(http.MethodDelete)

	return router
}
