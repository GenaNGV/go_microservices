package router

import (
	dao "auth/dao"
	"auth/model"
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

func Authenticate(w http.ResponseWriter, r *http.Request) {

	log.Print("Authenticate...")
	err := r.ParseForm()

	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	u := new(model.User)

	w.Header().Set("Content-Type", "application/json")
	if res := dao.DB.Where(&model.User{Email: email}).First(&u); res.RowsAffected <= 0 {
		w.Write([]byte("{'error': true, 'general': 'Invalid Credentials.'}"))
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		w.Write([]byte("{'error': true, 'general': 'Invalid Credentials.'}"))
		return
	}

	userBytes, err := json.Marshal(u)

	w.Write(userBytes)
}
