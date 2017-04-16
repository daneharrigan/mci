package db

import (
	"fmt"
	"log"
	"strings"
	"time"
)

func (m *User) FindBy(name, value string) error {
	return findBy(m, name, value)
}

func (m *User) Create() error {
	return create(m)
}

func (m *User) Update() error {
	return update(m)
}

func (m *User) Destroy() error {
	return destroy(m)
}

func (m *User) Comics(startedAt, endedAt time.Time) ([]*Comic, error) {
	comic := new(Comic)
	var columns []string
	for _, c := range comic.Columns() {
		columns = append(columns, fmt.Sprintf("%s.%s", comic.TableName(), c))
	}

	query := fmt.Sprintf(`
		SELECT %s
		FROM comics
		JOIN series ON series.id = comics.series_id
		JOIN user_series ON user_series.series_id = series.id
		JOIN users ON users.id = user_series.user_id
		WHERE
			users.id = $1 AND
			comics.released_at BETWEEN $2 AND $3`,
		strings.Join(columns, ", "))

	if cfg.Debug {
		log.Printf("DEBUG: %s", query)
	}

	rows, err := db.Query(query, m.ID, startedAt, endedAt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var comics []*Comic
	for rows.Next() {
		comic := new(Comic)
		if err := rows.Scan(comic.Values()...); err != nil {
			return nil, err
		}

		comics = append(comics, comic)
	}

	return comics, nil
}

// for Record interface
func (m *User) Values() []interface{} {
	return []interface{}{&m.ID, &m.Email, &m.Provider, &m.Picture, &m.CreatedAt, &m.UpdatedAt}
}

func (m *User) Columns() []string {
	return []string{"id", "email", "provider", "picture", "created_at", "updated_at"}
}

func (m *User) TableName() string {
	return "users"
}

func (m *User) IdentifierName() string {
	return "id"
}

func (m *User) IdentifierValue() string {
	return m.ID
}

func (m *User) Touch() {
	touch(m)
}
