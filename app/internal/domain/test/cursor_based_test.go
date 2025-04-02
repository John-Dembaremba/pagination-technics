package test

import (
	"context"
	"log"
	"testing"

	"github.com/John-Dembaremba/pagination-technics/internal/domain"
	"github.com/John-Dembaremba/pagination-technics/internal/domain/pagination"
	"github.com/John-Dembaremba/pagination-technics/internal/model"
	"github.com/John-Dembaremba/pagination-technics/internal/repo"
	"github.com/John-Dembaremba/pagination-technics/pkg"
)

func TestCursorBasedRead(t *testing.T) {
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
	seedHandler := domain.SeedHandler{
		Generator: domain.DataGenHandler{},
		Repo:      repoInterface,
	}

	if err := seedHandler.Seed(ctx, 100); err != nil {
		log.Fatalf("Failed to load test data with error: %v", err)
	}

	repoHandler := repo.RepositoryHandler{Db: db}
	handler := pagination.CursorBasedHandler{Repo: repoHandler}

	t.Run("pagination data", func(t *testing.T) {

		testCases := []struct {
			name      string
			cursor    int
			limit     int
			isSuccess bool
			expected  model.UsersCursorBasedMetaData
		}{
			{
				name:      "success - next cursor (10)",
				cursor:    10,
				limit:     10,
				isSuccess: true,
				expected: model.UsersCursorBasedMetaData{
					Users:      model.UsersData{},
					NextCursor: 1,
				},
			},
			{
				name:      "success - next cursor (10)",
				cursor:    20,
				limit:     10,
				isSuccess: true,
				expected: model.UsersCursorBasedMetaData{
					Users:      model.UsersData{},
					NextCursor: 10,
				},
			},
			{
				name:      "success - next cursor (90)",
				cursor:    100,
				limit:     10,
				isSuccess: true,
				expected: model.UsersCursorBasedMetaData{
					Users:      model.UsersData{},
					NextCursor: 90,
				},
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {

				got, err := handler.Retrieve(ctx, tc.cursor, tc.limit)
				if err != nil {
					t.Errorf("expected error nil, got %v", err)
				}

				if tc.isSuccess {
					if tc.expected.NextCursor != got.NextCursor {
						t.Errorf("expected next cursor: %v, got %v", tc.expected.NextCursor, got.NextCursor)
					}
				} else {
					if tc.expected.NextCursor == got.NextCursor {
						t.Errorf("expected next cursor: %v, got %v to be unequal", tc.expected.NextCursor, got.NextCursor)
					}
				}
			})
		}
	})
}
