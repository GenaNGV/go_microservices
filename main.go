package main

import (
	"auth/dao"
	"auth/router"
	"github.com/gorilla/mux"
	"net/http"
)

func createRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/login", router.Authenticate).Methods("POST")
	return r
}

func main() {
	dao.DatabaseConnect()

	router := createRouter()
	err := http.ListenAndServe(":3002", router)
	if err != nil {
		panic(err.Error())
	}
}
