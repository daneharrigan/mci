package db

import (
	"database/sql"
	"fmt"
	"strings"
)

type Users struct {
	err error
}

func (u *Users) All() <-chan *User {
	user := new(User)
	columns := strings.Join(user.Columns(), ", ")
	query := fmt.Sprintf("SELECT %s FROM %q ORDER BY created_at",
		columns, user.TableName())
	rows, err := db.Query(query)

	if err != nil {
		u.err = err
		return nil
	}

	ch := make(chan *User)
	go func(u *Users, rows *sql.Rows, ch chan<- *User) {
		defer rows.Close()
		defer close(ch)
		for rows.Next() {
			user := new(User)
			if err := rows.Scan(user.Values()...); err != nil {
				u.err = err
				return
			}

			ch <- user
		}

		u.err = rows.Err()
	}(u, rows, ch)

	return ch
}

func (u *Users) Err() error {
	return u.err
}
