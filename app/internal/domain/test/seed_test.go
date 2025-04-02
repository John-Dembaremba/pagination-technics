package test

import (
	"context"
	"database/sql"
	"log"
	"testing"

	"github.com/John-Dembaremba/pagination-technics/internal/domain"
	"github.com/John-Dembaremba/pagination-technics/internal/repo"
	"github.com/John-Dembaremba/pagination-technics/pkg"
)

func assertCreatedUser(t testing.TB, db *sql.DB, itemsNum int) {
	t.Helper()
	var count int
	err := db.QueryRow("SELECT COUNT(id) FROM users").Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	if count != itemsNum {
		t.Errorf("expected number of users: %v, got %v", itemsNum, count)
	}
}

func TestSeedHandler(t *testing.T) {
	testCases := []struct {
		name        string
		itemsNum    int
		expectedErr any
	}{
		{
			name:        "test-1000",
			itemsNum:    1000,
			expectedErr: nil,
		},
	}

	dbAttributes := pkg.DbAttributes{
		DbName:     "pagination-app",
		DbUserName: "user",
		DbPassword: "mypassword",
		MappedPort: "5432",
	}

	ctx := context.Background()
	testContainer := dbAttributes.PgTestContainerSetup(ctx)
	schemaFile := "../../../../pkg/schema.sql"
	db := dbAttributes.DbSetup(ctx, testContainer, schemaFile)
	defer pkg.TearDown(db, testContainer)

	repoInterface := repo.RepositoryHandler{Db: db}
	handler := domain.SeedHandler{
		Generator: domain.DataGenHandler{},
		Repo:      repoInterface,
	}

	for _, tc := range testCases {

		t.Run(tc.name, func(t *testing.T) {
			err := handler.Seed(ctx, tc.itemsNum)

			if err != tc.expectedErr {
				switch tc.expectedErr {
				case nil:
					if err != nil {
						t.Errorf("expected error: nil, got %v", err)
					}
				default:
					if err != nil {
						if err.Error() != tc.expectedErr {
							t.Errorf("expected error: %v, got %v", tc.expectedErr, err)
						}
					}

				}
			}
			assertCreatedUser(t, db, tc.itemsNum)
		})
	}
}
