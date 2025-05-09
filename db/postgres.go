package db

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type PostgresDB struct {
	DB *sql.DB
}

func NewPostgresDB(url string) (*PostgresDB, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &PostgresDB{DB: db}, nil
}

func (p *PostgresDB) Close() {
	p.DB.Close()
}
