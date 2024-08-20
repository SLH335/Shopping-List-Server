package sqlite

import (
	"database/sql"

	. "github.com/slh335/shoppinglistserver"
)

type UserService struct {
	DB *sql.DB
}

func (m *UserService) GetUser(username string) (user User, err error) {
	stmt := "SELECT * FROM users WHERE username=?"
	row := m.DB.QueryRow(stmt, username)

	err = row.Scan(&user.Id, &user.Username, &user.PasswordHash)
	if err != nil {
		return user, err
	}
	return user, nil
}
