package auth

import (
	"database/sql"
	"errors"
	"fmt"
)

func (d *Database) AddUser(userName, pass string) error {
	const addFnError = "addUserFn: %w"

	var chekedUserName string

	_, err := d.db.Exec("CREATE TABLE IF NOT EXISTS auth (username varchar(50), pass varchar(20))")
	if err != nil {
		// log.Error(addFnError, "db.Exec", err)

		return err
	}

	if err := d.db.QueryRow("SELECT username FROM auth WHERE username=$1", userName).Scan(&chekedUserName); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("user with current name (%s) exists", chekedUserName)
		}
	}

	if _, err := d.db.Exec(`INSERT INTO auth (username, pass) VALUES ($1, $2)`, userName, pass); err != nil {
		return fmt.Errorf(addFnError, err)
	}

	return nil
}
