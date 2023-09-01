package addnote

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"testTask/internal/db/postgres/auth"
	"testTask/internal/db/postgres/notes"
	authuser "testTask/internal/handlers/middleware/authUser"
)

type Request struct {
	UserName string `json:"username"`
	NoteBody string `json:"body"`
}

type Response struct {
	// Status     string `json:"status"`
	StatusCode int    `json:"statuscode"`
	Error      string `json:"error,omitempty"`
}

func AddNoteHandler(log *slog.Logger, auth *auth.Database, noteDB *notes.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authuser.BasicAuthHandler(w, r, log, auth)

		req := Request{}
		res := Response{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Error("AddNoteHandler", "json.Decode", err)
			res.Error = err.Error()
			res.StatusCode = http.StatusInternalServerError
			if err = json.NewEncoder(w).Encode(res); err != nil {
				log.Error(err.Error())

				return
			}

			return
		}

		if err := noteDB.AddNote(req.UserName, req.NoteBody); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			res.Error = err.Error()
			res.StatusCode = http.StatusInternalServerError
			if err = json.NewEncoder(w).Encode(res); err != nil {
				log.Error("AddNoteHandler", "json.Encode", err)

				return
			}
		}

		w.WriteHeader(http.StatusOK)
		res.StatusCode = http.StatusOK
		log.Info("Note add succesful", "Remote address", r.RemoteAddr)
		if err := json.NewEncoder(w).Encode(res); err != nil {
			log.Error("AddNoteHandler", "json.Encode", err)

			return
		}
	}
}
