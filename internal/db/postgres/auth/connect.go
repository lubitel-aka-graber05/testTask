package auth

import (
	"database/sql"
	"log/slog"

	_ "github.com/lib/pq"
)

type Database struct {
	db *sql.DB
}

func ConnectUserDB(log *slog.Logger, connData string) (*Database, error) {
	const (
		driverName = "postgres"

		namedFnError = "ConnectUserDB"
	)

	db, err := sql.Open(driverName, connData)
	if err != nil {
		log.Error(namedFnError, "sql.Open", err)

		return nil, err
	}
	if err = db.Ping(); err != nil {
		log.Error(namedFnError, "db.Ping", err)

		return nil, err
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS auth (username varchar(50), pass varchar(20))")
	if err != nil {
		log.Error(namedFnError, "db.Exec", err)

		return nil, err
	}

	return &Database{db: db}, nil
}
