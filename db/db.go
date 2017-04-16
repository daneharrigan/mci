package db

import (
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"

	"github.com/daneharrigan/mci/config"
)

var (
	cfg = config.New()
	db  *sql.DB
)

func init() {
	var err error
	db, err = sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Panicf("fn=Open error=%q", err)
	}
}

type Record interface {
	Values() []interface{}
	Columns() []string
	TableName() string
	IdentifierName() string
	IdentifierValue() string
	Touch()
}

func NewUUID() string {
	return uuid.New().String()
}

func IsNotFound(err error) bool {
	if err == nil {
		return false
	}

	return err.Error() == "sql: no rows in result set"
}

func findBy(m Record, name, value string) error {
	columns := strings.Join(m.Columns(), ", ")
	query := fmt.Sprintf("SELECT %s FROM %q WHERE %q = $1 LIMIT 1",
		columns, m.TableName(), name)

	if cfg.Debug {
		log.Printf("DEBUG: %s", query)
	}

	row := db.QueryRow(query, value)
	return row.Scan(m.Values()...)
}

func create(m Record) error {
	var placeholders []string
	for i := 0; i < len(m.Columns()); i++ {
		placeholders = append(placeholders, fmt.Sprintf("$%d", i+1))
	}

	columns := strings.Join(m.Columns(), ", ")
	values := strings.Join(placeholders, ", ")
	query := fmt.Sprintf("INSERT INTO %q (%s) VALUES(%s) RETURNING %s",
		m.TableName(), columns, values, columns)

	if cfg.Debug {
		log.Printf("DEBUG: %s", query)
	}

	m.Touch()
	row := db.QueryRow(query, m.Values()...)
	return row.Scan(m.Values()...)
}

func update(m Record) error {
	var placeholders []string
	for i, c := range m.Columns() {
		placeholders = append(placeholders, fmt.Sprintf("%q = $%d", c, i+1))
	}

	assignments := strings.Join(placeholders, ", ")
	columns := strings.Join(m.Columns(), ", ")
	query := fmt.Sprintf("UPDATE %q SET %s WHERE %q = $%d RETURNING %s",
		m.TableName(), assignments, m.IdentifierName(), len(m.Columns()), columns)

	if cfg.Debug {
		log.Printf("DEBUG: %s", query)
	}

	m.Touch()
	row := db.QueryRow(query, m.Values()...)
	return row.Scan(m.Values()...)
}

func destroy(m Record) error {
	query := fmt.Sprintf("DELETE FROM %q WHERE %q = $1",
		m.TableName(), m.IdentifierName())

	if cfg.Debug {
		log.Printf("DEBUG: %s", query)
	}

	_, err := db.Exec(query, m.IdentifierValue())
	return err
}

func touch(m interface{}) {
	v := reflect.Indirect(reflect.ValueOf(m))
	now := reflect.ValueOf(time.Now())

	f := v.FieldByName("ID")
	if f.String() == "" {
		f.SetString(NewUUID())
	}

	f = v.FieldByName("CreatedAt")
	if f.Interface().(time.Time).IsZero() {
		f.Set(now)
	}

	f = v.FieldByName("UpdatedAt")
	f.Set(now)
}
