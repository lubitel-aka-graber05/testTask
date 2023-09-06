package adduser

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"testTask/internal/db/postgres/auth"
)

type Request struct {
	UserName string `json:"username"`
	Pass     string `json:"pass"`
}

type Response struct {
	Error  string `json:"error,omitempty"`
	Status int    `json:"statuscode"`
}

func AddUserHandler(log *slog.Logger, auth *auth.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := Request{}
		res := Response{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Error("AddUserHandler", "json.Decode", err)
			res.Status = http.StatusInternalServerError
			res.Error = err.Error()
			if err = json.NewEncoder(w).Encode(res); err != nil {
				log.Error("AddUserHandler", "json.Encode", err)

				return
			}

			return
		}

		if strings.EqualFold(req.UserName, "") || strings.EqualFold(req.Pass, "") {
			w.WriteHeader(http.StatusNetworkAuthenticationRequired)
			res.Status = http.StatusNetworkAuthenticationRequired
			res.Error = "no auth data"
			log.Error("AddUserHandler", "Authentication", res.Error)
			if err := json.NewEncoder(w).Encode(res); err != nil {
				log.Error("AddUserHandler", "json.Encode", err)

				return
			}
		}

		if err := auth.AddUser(req.UserName, req.Pass); err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		w.WriteHeader(http.StatusOK)
		res.Status = http.StatusOK
		if err := json.NewEncoder(w).Encode(res); err != nil {
			log.Error("AddUserHandler", "json.Encode", err)

			return
		}

		log.Info("Create user successful", "Username", req.UserName, "Remote Address", r.RemoteAddr, "Remote host", r.Host)
	}
}
