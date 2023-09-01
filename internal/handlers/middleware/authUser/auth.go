package authuser

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/json"
	"log/slog"
	"net/http"

	"testTask/internal/db/postgres/auth"
)

type Request struct {
	UserName string `json:"username"`
	Pass     string `json:"pass"`
}

type Response struct {
	Error  string `json:"error,omitempty"`
	Status int    `json:"status"`
}

// type CheckAuthInfomer interface {
// 	AuthInfo(*http.Request) ([32]byte, [32]byte, error)
// }

func BasicAuthHandler(w http.ResponseWriter, r *http.Request, log *slog.Logger, authDB *auth.Database) {
	user, pass, ok := r.BasicAuth()
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		w.Header().Set("WWW-Authenticate", "Basic realm=\"restricted\"")
		log.Error("no authorization header")

		res := Response{Status: http.StatusUnauthorized}
		if err := json.NewEncoder(w).Encode(res); err != nil {
			log.Error(err.Error())

			return
		}

		return
	}

	userHash := sha256.Sum256([]byte(user))
	passHash := sha256.Sum256([]byte(pass))

	expectedUserHash, expectedPassHash, err := authDB.AuthInfo(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Error(err.Error())

		res := Response{Status: http.StatusInternalServerError, Error: err.Error()}

		if err = json.NewEncoder(w).Encode(res); err != nil {
			log.Error(err.Error())

			return
		}

		return
	}

	userOK := subtle.ConstantTimeCompare(userHash[:], expectedUserHash[:]) == 1
	passOK := subtle.ConstantTimeCompare(passHash[:], expectedPassHash[:]) == 1

	if !userOK && !passOK {
		w.WriteHeader(http.StatusUnauthorized)
		w.Header().Set("WWW-Authenticate", "Basic realm=\"restricted\"")
		res := Response{Status: http.StatusUnauthorized}
		log.Error("invalid authentication data")

		if err = json.NewEncoder(w).Encode(res); err != nil {
			log.Error(err.Error())

			return
		}

		return
	}
}
