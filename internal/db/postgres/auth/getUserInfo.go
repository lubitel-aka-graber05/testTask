package auth

import (
	"crypto/sha256"
	"fmt"
)

func (d *Database) AuthInfo(userName string) (userHash, passHash [32]byte, err error) {
	const checkFnError = "CheckAuthInfo: %w"

	rows, err := d.db.Query(`SELECT username, pass FROM auth WHERE username=$1`, userName)
	if err != nil {
		return userHash, passHash, fmt.Errorf(checkFnError, err)
	}

	var checkedName, checkedPass string

	for rows.Next() {
		if err = rows.Scan(&checkedName, &checkedPass); err != nil {
			return userHash, passHash, fmt.Errorf(checkFnError, err)
		}
	}

	userHash = sha256.Sum256([]byte(checkedName))
	passHash = sha256.Sum256([]byte(checkedPass))

	return userHash, passHash, nil
}
