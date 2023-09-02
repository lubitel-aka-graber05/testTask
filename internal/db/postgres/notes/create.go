package notes

import (
	"fmt"
)

func (d *Database) AddNote(userName, noteBody string) error {
	const namedFnError = "AddNoteFn: %w"

	_, err := d.db.Exec("CREATE TABLE IF NOT EXISTS note (username varchar(50), body text)")
	if err != nil {
		// log.Error(namedFnError, "db.Exec", err)

		return fmt.Errorf(namedFnError, err)
	}

	_, err = d.db.Exec(`INSERT INTO note (username, body) VALUES($1, $2)`, userName, noteBody)
	if err != nil {
		return fmt.Errorf(namedFnError, err)
	}

	return nil
}
