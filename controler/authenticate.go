package controler

import (
	service "auth/service"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func Login(w http.ResponseWriter, r *http.Request) {

	log.Trace("authenticate user")
	err := r.ParseForm()

	if err != nil {
		log.WithError(err).Error("unable to parse form")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	email := r.FormValue("email")
	password := r.FormValue("password")

	auth, err := service.Login(email, password)

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(auth)
}

func Status(w http.ResponseWriter, r *http.Request) {

	log.Trace("login user")

	token := r.Header.Get("Authorization")

	w.Header().Set("Content-Type", "application/json")

	user, err := service.Status(token)

	if err != nil || user == nil {
		w.WriteHeader(http.StatusForbidden)
		_ = json.NewEncoder(w).Encode(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(user)
}
