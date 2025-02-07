package pkg

import (
	"database/sql"

	"github.com/DATA-DOG/go-sqlmock"
)

func DataDogDbMock() (*sql.DB, *sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		return nil, nil, err
	}
	return db, &mock, nil
}
