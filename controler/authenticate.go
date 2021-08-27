package controler

import (
	service "auth/service"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func Authenticate(w http.ResponseWriter, r *http.Request) {

	log.Trace("Authenticate user")
	err := r.ParseForm()

	if err != nil {
		log.Error("Unable to parse form ", fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	email := r.FormValue("email")
	password := r.FormValue("password")

	user, err := service.Authenticate(email, password)

	if err != nil {
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	userBytes, err := json.Marshal(user)

	_, _ = w.Write(userBytes)
}
