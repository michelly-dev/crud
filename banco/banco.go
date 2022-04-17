package banco

import (
	"database/sql"

	_ "github.com/lib/pq"
)

func Conectar() (*sql.DB, error) {
	connStr := "host=localhost user=postgres password=13798246 dbname=devbook sslmode=disable"

	db, err := sql.Open("postgres", connStr)

	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, err
}
