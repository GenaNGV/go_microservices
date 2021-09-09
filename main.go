package main

import (
	"auth/controler"
	"auth/enviroment"
	"auth/middleware"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func createRouter() *mux.Router {

	r := mux.NewRouter()
	r.HandleFunc("/login", controler.Login).Methods("POST")
	r.HandleFunc("/parse", controler.Parse).Methods("POST")

	r.Use(middleware.Validate)

	return r
}

func main() {

	enviroment.Env = enviroment.NewEnvironment()

	router := createRouter()
	log.Fatal(http.ListenAndServe(":3002", router))
}
