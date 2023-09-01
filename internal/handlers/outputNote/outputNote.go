package outputnote

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"testTask/internal/db/postgres/auth"
	"testTask/internal/db/postgres/notes"
	authuser "testTask/internal/handlers/middleware/authUser"
	"testTask/internal/handlers/middleware/yandexspeller"
)

type RequestByName struct {
	UserName string `json:"username"`
}

type ResponseByName struct {
	UserName string   `json:"username,omitempty"`
	Error    string   `json:"error,omitempty"`
	NoteBody []string `json:"body,omitempty"`
	Status   int      `json:"status"`
}

type ResponseAll struct {
	Notes  map[string][]string `json:"notes,omitempty"`
	Error  string              `json:"error,omitempty"`
	Status int                 `json:"status"`
}

func OutputByNameHandler(log *slog.Logger, auth *auth.Database, notesDB *notes.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authuser.BasicAuthHandler(w, r, log, auth)

		req := RequestByName{}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Error(err.Error())

			return
		}

		outData, err := notesDB.OutputByName(req.UserName)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Error(err.Error())

			res := ResponseByName{Status: http.StatusInternalServerError, Error: err.Error()}
			if err = json.NewEncoder(w).Encode(res); err != nil {
				log.Error(err.Error())

				return
			}

			return
		}

		outData = yandexspeller.YandexSpeller(log, outData, &req.UserName)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		res := ResponseByName{UserName: req.UserName, NoteBody: outData[req.UserName], Status: http.StatusOK}

		if err := json.NewEncoder(w).Encode(res); err != nil {
			log.Error(err.Error())

			return
		}
	}
}
