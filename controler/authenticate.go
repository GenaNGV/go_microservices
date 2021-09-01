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

	user, err := service.Authenticate(email, password)

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func Check(w http.ResponseWriter, r *http.Request) {

	log.Trace("login user")

	authorization := r.Header.Get("Authorization")

	w.Header().Set("Content-Type", "application/json")

	user, err := service.Check(authorization)

	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}
