package main

import (
	"net/http"
	"origin-main/controllers"

	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/auth", controllers.GetToken).Methods(http.MethodPost)
	router.HandleFunc("/auth", controllers.Render_login).Methods(http.MethodGet)
	router.HandleFunc("/signup", controllers.Render_signup).Methods(http.MethodGet)
	router.HandleFunc("/signup", controllers.NewUser).Methods(http.MethodPost)
	router.Handle("/test", controllers.CheckSession(http.HandlerFunc(controllers.TestRequestToApi)))
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
