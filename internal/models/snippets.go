package models

import (
	"database/sql"
	"errors"
	"time"
)

type SnippetModelInterface interface {
	Insert(title string, content string, expires int) (int, error)
	Get(id int) (*Snippet, error)
	Latest() ([]*Snippet, error)
}

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type SnippetModel struct {
	DB *sql.DB
}

func (m *SnippetModel) Insert(title, content string, expires int) (int, error) {
	stmt := `
  INSERT INTO snippets (title, content, created, expires)
  VALUES (?, ?, CURRENT_TIMESTAMP, datetime(CURRENT_TIMESTAMP, '+' || ? || ' days'))
  `

	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *SnippetModel) Get(id int) (*Snippet, error) {
	if id < 1 {
		return nil, ErrNoRecord
	}

	stmt := `
  SELECT id, title, content, created, expires FROM snippets
  WHERE expires > CURRENT_TIMESTAMP AND id = ?
  `

	s := &Snippet{}
	args := []any{&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires}

	if err := m.DB.QueryRow(stmt, id).Scan(args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		}

		return nil, err
	}

	return s, nil
}

// Return the 10 most recently created snippets
func (m *SnippetModel) Latest() ([]*Snippet, error) {
	stmt := `
  SELECT id, title, content, created, expires FROM snippets
  WHERE expires > CURRENT_TIMESTAMP ORDER BY id DESC LIMIT 10
  `

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	snippets := []*Snippet{}

	for rows.Next() {
		s := &Snippet{}
		args := []any{&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires}

		if err := rows.Scan(args...); err != nil {
			return nil, err
		}

		snippets = append(snippets, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
