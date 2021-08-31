package main

import (
	"auth/controler"
	"auth/enviroment"
	"github.com/gorilla/mux"
	"net/http"
)

func createRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/login", controler.Login).Methods("POST")
	r.HandleFunc("/status", controler.Check).Methods("POST")
	return r
}

func main() {

	enviroment.Env = enviroment.NewEnvironment()

	router := createRouter()
	err := http.ListenAndServe(":3002", router)
	if err != nil {
		panic(err.Error())
	}
}
