// movie.go
package models

import (
	"database/sql"
	"errors"
	"strings"
	"time"
)

var (
	ErrNoMovie   = errors.New("models: no matching movie found")
	ErrDuplicate = errors.New("models: duplicate movie title")
)

type Movie struct {
	ID        int
	Title     string
	Genre     string
	Created   time.Time
	DepID     int
	Name      string
	Quantinty int
}

type MovieModel struct {
	DB *sql.DB
}

func (m *MovieModel) Add(title, genre string) error {
	stmt := `INSERT INTO movies (title, genre, created) VALUES (?, ?, UTC_TIMESTAMP())`

	_, err := m.DB.Exec(stmt, title, genre)
	if err != nil {
		if isDuplicateError(err) {
			return ErrDuplicate
		}
		return err
	}
	return nil
}

func (m *MovieModel) Update(title, genre string, id int) error {
	stmt := `UPDATE movies SET title=?, genre=?, created=UTC_TIMESTAMP() WHERE id=?`
	_, err := m.DB.Exec(stmt, title, genre, id)
	if err != nil {
		if isDuplicateError(err) {
			return ErrDuplicate
		}
		return err
	}
	return nil
}

func (m *MovieModel) Delete(id int) error {
	stmt := `DELETE FROM movies WHERE id=?`
	_, err := m.DB.Exec(stmt, id)
	if err != nil {
		if isDuplicateError(err) {
			return ErrDuplicate
		}
		return err
	}
	return nil
}

func (m *MovieModel) Get(id int) (*Movie, error) {
	stmt := `SELECT id, title, genre, created FROM movies WHERE id = ?`
	row := m.DB.QueryRow(stmt, id)

	movie := &Movie{}
	err := row.Scan(&movie.ID, &movie.Title, &movie.Genre, &movie.Created)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoMovie
		}
		return nil, err
	}

	return movie, nil
}

func (m *MovieModel) All() ([]*Movie, error) {
	stmt := `SELECT id, title, genre, created FROM movies ORDER BY created DESC`
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []*Movie

	for rows.Next() {
		movie := &Movie{}
		err := rows.Scan(&movie.ID, &movie.Title, &movie.Genre, &movie.Created)
		if err != nil {
			return nil, err
		}
		movies = append(movies, movie)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return movies, nil
}

func isDuplicateError(err error) bool {
	return err != nil && strings.Contains(err.Error(), "Error 1062:")
}
