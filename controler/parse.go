package controler

import (
	"auth/service"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

func Parse(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()

	if err != nil {
		log.WithError(err).Error("unable to parse form")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	file := r.FormValue("file")
	chars := r.FormValue("chars")
	log.WithFields(log.Fields{"file": file}).Trace("Parsing file")

	w.Header().Set("Content-Type", "application/json")

	arr := strings.Split(chars, ",")

	user, _ := service.TokenDetail(r.Header.Get("Authorization"))

	jobInfo, err := service.Parse(file, arr, user)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(jobInfo)

}
