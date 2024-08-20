package sqlite

import (
	"database/sql"
	"fmt"
	"time"

	. "github.com/slh335/shoppinglistserver"
	"github.com/slh335/shoppinglistserver/crypto"
)

type AuthService struct {
	DB *sql.DB
}

func (m *AuthService) Register(username, password string) (user User, err error) {
	passwordHash, err := crypto.HashPassword(password)
	if err != nil {
		return user, err
	}

	stmt := "INSERT INTO users (username, password_hash) VALUES (?, ?)"
	res, err := m.DB.Exec(stmt, username, passwordHash)
	if err != nil {
		return user, err
	}

	lastInsertId, _ := res.LastInsertId()
	user = User{
		Id:           int(lastInsertId),
		Username:     username,
		PasswordHash: passwordHash,
	}
	return user, nil
}

func (m *AuthService) Login(username, password string) (user User, err error) {
	stmt := "SELECT * FROM users WHERE username=?"
	row := m.DB.QueryRow(stmt, username)

	err = row.Scan(&user.Id, &user.Username, &user.PasswordHash)
	if err != nil {
		return user, err
	}

	match := crypto.VerifyPassword(password, user.PasswordHash)
	if !match {
		return user, fmt.Errorf("error: invalid password")
	}
	return user, nil
}

func (s *AuthService) NewSession(user User, validDays int) (session Session, err error) {
	token := crypto.GenerateToken(64)
	createdAt := time.Now()
	var expiresAt time.Time

	if validDays > 0 {
		stmt := "INSERT INTO sessions (token, user_id, created_at, expires_at) VALUES (?, ?, ?, ?)"
		expiresAt = time.Now().Add(24 * time.Hour * time.Duration(validDays))
		_, err = s.DB.Exec(stmt, token, user.Id, createdAt, expiresAt)
	} else {
		stmt := "INSERT INTO sessions (token, user_id, created_at) VALUES (?, ?, ?)"
		_, err = s.DB.Exec(stmt, token, user.Id, createdAt)
	}
	if err != nil {
		return Session{}, err
	}

	return Session{
		Token: token,
		User: User{
			Id:       user.Id,
			Username: user.Username,
		},
		CreatedAt: createdAt,
		ExpiresAt: expiresAt,
	}, nil
}

func (s *AuthService) VerifySession(token string) (user User, err error) {
	if token == "" {
		return user, fmt.Errorf("error: no token provided")
	}
	stmt := `
		SELECT users.*
		FROM sessions
		INNER JOIN users ON sessions.user_id=users.id
		WHERE sessions.token=?`
	row := s.DB.QueryRow(stmt, token)
	err = row.Scan(&user.Id, &user.Username, &user.PasswordHash)
	if err != nil {
		return user, err
	}
	return user, nil
}
