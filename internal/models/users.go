package models

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}

type UserModel struct {
	DB *sql.DB
}

// Insert a new user to the users table.
func (m *UserModel) Insert(name, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `
	INSERT INTO users (name, email, hashed_password, created)
	VALUES(?, ?, ?, datetime('now', 'utc'));`
	args := []any{name, email, string(hashedPassword)}

	if _, err := m.DB.Exec(stmt, args...); err != nil {
		var sqliteError sqlite3.Error

		if errors.As(err, &sqliteError) {
			if sqliteError.Code == sqlite3.ErrConstraint &&
				strings.Contains(sqliteError.Error(), "UNIQUE constraint failed: users.email") {
				return ErrDuplicateEmail
			}
		}
		return err
	}

	return nil
}

// Verify whether a user exists with the provided email address and passsword
// Retuns the relevant user ID if they do.
func (m *UserModel) Authenticate(email, password string) (int, error) {
	var id int
	var hashedPassword []byte

	stmt := "SELECT id, hashed_password FROM users WHERE email = ?"
	if err := m.DB.QueryRow(stmt, email).Scan(&id, &hashedPassword); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrNoRecord
		}

		return 0, err
	}

	if err := bcrypt.CompareHashAndPassword(hashedPassword, []byte(password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrNoRecord
		}

		return 0, err
	}

	return id, nil
}

// Checks if a user exists with a specific ID.
func (m *UserModel) Exists(id int) (bool, error) {
	return false, nil
}
