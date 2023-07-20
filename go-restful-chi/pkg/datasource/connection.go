package datasource

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

func NewPgClient(connString string) *sql.DB {
	conn, err := sql.Open("postgres", connString)
	if err != nil {
		fmt.Println(err)
	}

	if err = conn.Ping(); err != nil {
		log.Fatal(err)
	}

	return conn
}
