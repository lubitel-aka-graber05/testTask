package notes

import "fmt"

func (d *Database) AddNote(userName, noteBody string) error {
	const namedFnError = "AddNoteFn: %w"

	_, err := d.db.Exec(`INSERT INTO note (username, body) VALUES($1, $2)`, userName, noteBody)
	if err != nil {
		return fmt.Errorf(namedFnError, err)
	}

	return nil
}
