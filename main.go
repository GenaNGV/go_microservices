package main

import (
	"auth/controler"
	"auth/enviroment"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func createRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/login", controler.Login).Methods("POST")
	r.HandleFunc("/status", controler.Status).Methods("GET")
	return r
}

func main() {

	enviroment.Env = enviroment.NewEnvironment()

	router := createRouter()
	log.Fatal(http.ListenAndServe(":3002", router))
}
