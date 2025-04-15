package pkg

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"
	"time"

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

	// Connection pooling
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(30 * time.Minute)

	log.Println("Database connection established")

	return db, nil
}

func RunMigration(db *sql.DB, query string) error {
	_, err := db.Exec(query)

	return err
}

func getMd5(userName, password string) string {
	// "pagy" "md5ac77bfe847b783150cc181043bd7d2d7"
	combined := password + userName

	hashV := md5.Sum([]byte(combined))
	hexString := "md5" + hex.EncodeToString(hashV[:])
	return hexString
}
