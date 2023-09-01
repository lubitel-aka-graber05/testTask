package notes

import (
	"fmt"
)

func (d *Database) OutputByName(userName string) (map[string][]string, error) {
	const namedFnError = "OutputByName: %w"

	output := make(map[string][]string)

	rows, err := d.db.Query("SELECT username, body FROM note WHERE username=$1", userName)
	if err != nil {
		return nil, fmt.Errorf(namedFnError, err)
	}
	defer rows.Close()

	for rows.Next() {
		var userName, noteBody string

		if err = rows.Scan(&userName, &noteBody); err != nil {
			return nil, fmt.Errorf(namedFnError, err)
		}

		if output[userName] == nil {
			output[userName] = make([]string, 0)
		}
		output[userName] = append(output[userName], noteBody)
	}

	

	return output, nil
}

func (d *Database) OutputAll() (map[string][]string, error) {
	const namedFnError = "OutputAll: %w"

	output := make(map[string][]string)
	
	rows, err := d.db.Query("SELECT username, body FROM note")
	if err != nil {
		return nil, fmt.Errorf(namedFnError, err)
	}
	defer rows.Close()

	for rows.Next() {
		var userName, noteBody string

		if err = rows.Scan(&userName, &noteBody); err != nil {
			return nil, fmt.Errorf(namedFnError, err)
		}

		if output[userName] == nil {
			output[userName] = make([]string, 0)
		}
		output[userName] = append(output[userName], noteBody)
	}

	return output, nil
}
