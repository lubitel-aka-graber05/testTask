package auth

import (
	"database/sql"
	"errors"
	"fmt"
)

func (d *Database) AddUser(userName, pass string) error {
	const addFnError = "addUserFn: %w"

	var chekedUserName string

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
