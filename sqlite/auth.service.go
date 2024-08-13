package sqlite

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"
	"time"

	. "github.com/slh335/shoppinglistserver"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	DB *sql.DB
}

func (m *AuthService) Register(username, password string) (id int, err error) {
	passwordHash, err := hashPassword(password)
	if err != nil {
		return -1, err
	}

	stmt := "INSERT INTO users (username, passwordHash) VALUES (?, ?)"
	res, err := m.DB.Exec(stmt, username, passwordHash)
	if err != nil {
		return -1, err
	}

	lastInsertId, _ := res.LastInsertId()
	return int(lastInsertId), nil
}

func (m *AuthService) Login(username, password string) (user User, err error) {
	stmt := "SELECT * FROM users WHERE username=?"
	row := m.DB.QueryRow(stmt, username)

	err = row.Scan(&user.Id, &user.Username, &user.PasswordHash)
	if err != nil {
		fmt.Println(err)
		return user, err
	}

	match := verifyPassword(password, user.PasswordHash)
	if !match {
		return user, fmt.Errorf("error: invalid password")
	}
	return user, nil
}

func (s *AuthService) NewSession(user User, validDays int) (token string, err error) {
	token = generateSessionToken(64)
	if validDays > 0 {
		stmt := "INSERT INTO sessions (token, userId, createdAt, expiresAt) VALUES (?, ?, ?, ?)"
		expiresAt := time.Now().Add(24 * time.Hour * time.Duration(validDays))
		_, err = s.DB.Exec(stmt, token, user.Id, time.Now(), expiresAt)
	} else {
		stmt := "INSERT INTO sessions (token, userId, createdAt) VALUES (?, ?, ?)"
		_, err = s.DB.Exec(stmt, token, user.Id, time.Now())
	}
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *AuthService) VerifySession(token string) (user User, err error) {
	if token == "" {
		return user, fmt.Errorf("error: no token provided")
	}
	stmt := `
		SELECT users.*
		FROM sessions
		INNER JOIN users ON sessions.userId=users.id
		WHERE sessions.token=?`
	row := s.DB.QueryRow(stmt, token)
	err = row.Scan(&user.Id, &user.Username, &user.PasswordHash)
	if err != nil {
		return user, err
	}
	return user, nil
}

func generateSessionToken(length int) (token string) {
	buf := make([]byte, length)
	_, err := rand.Read(buf)
	if err != nil {
		return ""
	}
	return base64.RawURLEncoding.EncodeToString(buf)
}

func hashPassword(password string) (hash string, err error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func verifyPassword(password, hash string) (match bool) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
