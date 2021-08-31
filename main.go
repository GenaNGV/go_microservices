package main

import (
	"auth/controler"
	"auth/enviroment"
	"auth/utils"
	"github.com/gorilla/mux"
	"net/http"
)

func createRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/login", controler.Login).Methods("POST")
	r.HandleFunc("/status", controler.Check).Methods("POST")
	return r
}

Eev *Envi
func main() {
	utils.Initialize("students.log")
	enviroment.Initialize()

	router := createRouter()
	err := http.ListenAndServe(":3002", router)
	if err != nil {
		panic(err.Error())
	}
}
