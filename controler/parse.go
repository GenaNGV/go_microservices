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
	searchTerm := r.FormValue("chars")
	rule := r.FormValue("rule")
	log.WithFields(log.Fields{"file": file, "rule": rule}).Trace("Parsing file")

	w.Header().Set("Content-Type", "application/json")

	terms := strings.Split(searchTerm, ",")

	user, _ := service.TokenDetail(r.Header.Get("Authorization"))

	jobInfo, err := service.Parse(file, terms, user, rule)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(jobInfo)

}
