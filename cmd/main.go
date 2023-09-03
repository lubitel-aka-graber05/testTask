package main

import (
	"net/http"

	"testTask/internal/config"
	"testTask/internal/db/postgres/auth"
	"testTask/internal/db/postgres/notes"
	addnote "testTask/internal/handlers/addNote"
	"testTask/internal/handlers/adduser"
	outputnote "testTask/internal/handlers/outputNote"
	"testTask/internal/logger"
)

func main() {
	log := logger.SetupLogger()

	const path = "configs/configs.yaml"

	cfg, err := config.CreateConfig(log, path)
	if err != nil {
		return
	}

	log.Info(
		"Start server configuration",
		"server address", cfg.ServerAddress,
		"server timeout", cfg.ServerTimeout,
	)

	authDB, err := auth.ConnectUserDB(log, cfg.AddressDB)
	if err != nil {
		log.Error(err.Error())

		return
	}

	notesDB, err := notes.ConnectNoteDB(log, cfg.AddressDB)
	if err != nil {
		log.Error(err.Error())

		return
	}

	router := http.NewServeMux()
	router.HandleFunc("/adduser", adduser.AddUserHandler(log, authDB))
	router.HandleFunc("/addnote", addnote.AddNoteHandler(log, authDB, notesDB))
	router.HandleFunc("/outbyname", outputnote.OutputByNameHandler(log, authDB, notesDB))

	server := http.Server{
		Addr:         cfg.ServerAddress,
		IdleTimeout:  cfg.ServerTimeout,
		WriteTimeout: cfg.ServerTimeout,
		ReadTimeout:  cfg.ServerTimeout,
		Handler:      router,
	}

	if err = server.ListenAndServe(); err != nil {
		log.Error(err.Error())

		return
	}
}
