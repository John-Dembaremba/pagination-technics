package test

import (
	"reflect"
	"testing"

	"github.com/John-Dembaremba/pagination-technics/internal/repo"
	"github.com/John-Dembaremba/pagination-technics/pkg"
)

func TestLimitOffsetRead(t *testing.T) {
	db, mock, err := pkg.DataDogDbMock()
	if err != nil {
		t.Errorf("db mock failed with error: %v", err)
	}
	defer db.Close()
	repoH := repo.RepositoryHandler{Db: db}
	offset, limit := 0, 10
	usersData := pkg.FakeUsersData(3)
	rows := mock.NewRows([]string{"id", "name", "surname"})

	for _, user := range usersData {
		rows.AddRow(
			user.ID,
			user.UserGenData.Name,
			user.UserGenData.Surname,
		)
	}
	query := "SELECT id, name, surname FROM users ORDER BY id LIMIT $1 OFFSET $2;"

	t.Run("success query", func(t *testing.T) {
		mock.ExpectQuery(query).WithArgs(limit, offset).WillReturnRows(rows)
		got, err := repoH.LimitOffsetRead(offset, limit)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}

		if err != nil {
			t.Errorf("expected no error but got %v", err)
		}

		if !reflect.DeepEqual(got, usersData) {
			t.Errorf("expected data: %v, got %v", usersData, got)
		}

	})
}
