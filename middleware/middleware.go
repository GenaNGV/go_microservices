package middleware

import (
	"auth/service"
	"encoding/json"
	"net/http"
)

func Validate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/login" {
			next.ServeHTTP(w, req)
		} else {
			token := req.Header.Get("Authorization")

			user, err := service.TokenDetail(token)

			if err != nil || user == nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusForbidden)
				_ = json.NewEncoder(w).Encode(err.Error())
				return
			}

			next.ServeHTTP(w, req)
		}
	})
}
