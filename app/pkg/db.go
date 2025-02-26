package pkg

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/lib/pq"
)

func DataDogDbMock() (*sql.DB, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		return nil, mock, err
	}
	return db, mock, nil
}

func NewPgDb(cntName, dbName, dbUser, dbPsw, dbPort string) (*sql.DB, error) {
	connStr := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable", dbUser, dbPsw, cntName, dbPort, dbName)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// Verify the connection
	if err := db.Ping(); err != nil {
		return nil, err
	}
	log.Println("Database connection established")

	return db, nil
}

func RunMigration(db *sql.DB, query string) error {
	_, err := db.Exec(query)

	return err
}
