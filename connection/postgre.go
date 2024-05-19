package connection

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

type Postgre struct {
	db *sql.DB
}

func (p *Postgre) DB() *sql.DB {
	return p.db
}

func NewPostgre() *Postgre {
	host := os.Getenv("PG_HOST")
	port := os.Getenv("PG_PORT")
	user := os.Getenv("PG_USER")
	password := os.Getenv("PG_PASSWORD")
	dbName := os.Getenv("PG_DB")

	psqlconn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbName)
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		panic(err)
	}
	return &Postgre{
		db: db,
	}
}
