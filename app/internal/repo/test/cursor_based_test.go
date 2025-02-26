package test

import (
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/John-Dembaremba/pagination-technics/internal/model"
	"github.com/John-Dembaremba/pagination-technics/internal/repo"
	"github.com/John-Dembaremba/pagination-technics/pkg"
)

func TestCursorBasedRead(t *testing.T) {
	db, mock, err := pkg.DataDogDbMock()
	if err != nil {
		t.Errorf("db mock failed with error: %v", err)
	}
	defer db.Close()
	repoH := repo.RepositoryHandler{Db: db}
	usersData := pkg.FakeUsersData(3)
	rows := mock.NewRows([]string{"id", "name", "surname"})

	for _, user := range usersData {
		rows.AddRow(
			user.ID,
			user.UserGenData.Name,
			user.UserGenData.Surname,
		)
	}

	assertHelper := func(t testing.TB, mock sqlmock.Sqlmock, got model.UsersData) {
		t.Helper()
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}

		if !reflect.DeepEqual(got, usersData) {
			t.Errorf("expected data: %v, got %v", usersData, got)
		}
	}

	t.Run("success query", func(t *testing.T) {

		t.Run("init cursor", func(t *testing.T) {
			cursor, limit := 1, 10
			query := "SELECT id, name, surname FROM users ORDER BY id DESC LIMIT $1;"
			mock.ExpectQuery(query).WithArgs(limit).WillReturnRows(rows)
			got, err := repoH.CursorBasedRead(cursor, limit)
			if err != nil {
				t.Errorf("expected no error but got %v", err)
			}
			assertHelper(t, mock, got)
		})

		t.Run("actual cursor", func(t *testing.T) {
			cursor, limit := 5, 3
			query := "SELECT id, name, surname FROM users WHERE id < $1 ORDER BY id DESC LIMIT $2;"
			mock.ExpectQuery(query).WithArgs(cursor, limit).WillReturnRows(rows)
			_, err := repoH.CursorBasedRead(cursor, limit)
			if err != nil {
				t.Errorf("expected no error but got %v", err)
			}
			// assertHelper(t, mock, got)
		})

	})
}
