package main

import (
	"net/http"
	"origin/midwr"

	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/auth/", midwr.GetToken).Methods(http.MethodPost)
	router.HandleFunc("/auth/", midwr.Render_login).Methods(http.MethodGet)
	return router
}

func Server(router *mux.Router) *http.Server {
	return &http.Server{
		Addr:    ":8282",
		Handler: router,
	}
}

func main() {
	Server(Router()).ListenAndServe()
}
