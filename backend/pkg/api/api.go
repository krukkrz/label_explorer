package api

import (
	"github.com/gorilla/mux"
	"github.com/krukkrz/label_explorer/pkg/api/handlers"
)

func NewRouter(handlers *handlers.Handlers) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/artists", handlers.Artists).Methods("GET")
	r.HandleFunc("/styles", handlers.Styles).Methods("GET")

	return r
}
