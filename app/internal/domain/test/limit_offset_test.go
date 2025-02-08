package test

import (
	"context"
	"log"
	"reflect"
	"testing"

	"github.com/John-Dembaremba/pagination-technics/internal/domain"
	"github.com/John-Dembaremba/pagination-technics/internal/domain/pagination"
	"github.com/John-Dembaremba/pagination-technics/internal/model"
	"github.com/John-Dembaremba/pagination-technics/internal/repo"
	"github.com/John-Dembaremba/pagination-technics/pkg"
)

func TestLimitOffSetHandler(t *testing.T) {
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

	if err := seedHandler.Seed(100); err != nil {
		log.Fatalf("Failed to load test data with error: %v", err)
	}

	repoHandler := repo.RepositoryHandler{Db: db}
	handler := pagination.LimitOffSetHandler{Repo: repoHandler}

	t.Run("pagination data", func(t *testing.T) {
		const (
			currentPage = "current page"
			nextPage    = "next page"
			prevPage    = "prev page"
			totalPages  = "total pages"
		)
		testCases := []struct {
			name     string
			page     int
			limit    int
			expected model.UsersPaginationMetaData
		}{
			{
				name:  currentPage,
				page:  1,
				limit: 10,
				expected: model.UsersPaginationMetaData{
					Users: pkg.FakeUsersData(5),
					Pagination: model.Pagination{
						CurrentPage: 1,
						NextPage:    2,
						PrevPage:    1,
						TotalPages:  5,
					},
				},
			},
			{
				name:  totalPages,
				page:  1,
				limit: 10,
				expected: model.UsersPaginationMetaData{
					Users: pkg.FakeUsersData(5),
					Pagination: model.Pagination{
						CurrentPage: 1,
						NextPage:    2,
						PrevPage:    1,
						TotalPages:  10,
					},
				},
			},
			{
				name:  nextPage,
				page:  3,
				limit: 10,
				expected: model.UsersPaginationMetaData{
					Users: pkg.FakeUsersData(5),
					Pagination: model.Pagination{
						CurrentPage: 1,
						NextPage:    4,
						PrevPage:    1,
						TotalPages:  10,
					},
				},
			},
			{
				name:  prevPage,
				page:  8,
				limit: 10,
				expected: model.UsersPaginationMetaData{
					Users: pkg.FakeUsersData(5),
					Pagination: model.Pagination{
						CurrentPage: 1,
						NextPage:    4,
						PrevPage:    7,
						TotalPages:  10,
					},
				},
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				got, err := handler.RetrieveUsers(tc.page, tc.limit)
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}

				if !reflect.DeepEqual(tc.expected.Pagination, got.Pagination) {

					switch tc.name {
					case currentPage:
						if tc.expected.Pagination.CurrentPage != got.Pagination.CurrentPage {
							t.Errorf("expected current page: %v, got %v", tc.expected.Pagination.CurrentPage, got.Pagination.CurrentPage)
						}
					case totalPages:
						if tc.expected.Pagination.TotalPages != got.Pagination.TotalPages {
							t.Errorf("expected total page: %v, got %v", tc.expected.Pagination.TotalPages, got.Pagination.TotalPages)
						}
					case nextPage:
						if tc.expected.Pagination.NextPage != got.Pagination.NextPage {
							t.Errorf("expected next page: %v, got %v", tc.expected.Pagination.NextPage, got.Pagination.NextPage)
						}
					case prevPage:
						if tc.expected.Pagination.PrevPage != got.Pagination.PrevPage {
							t.Errorf("expected prev page: %v, got %v", tc.expected.Pagination.PrevPage, got.Pagination.PrevPage)
						}

					}

				}

			})
		}

	})

	t.Run("user data", func(t *testing.T) {

		testCases := []struct {
			name  string
			page  int
			limit int
		}{
			{
				name:  "page 1",
				page:  1,
				limit: 10,
			},
			{
				name:  "page 0 error",
				page:  0,
				limit: 10,
			},
		}

		assertUserData := func(t testing.TB, page, limit int, got model.UsersData) {
			t.Helper()

			offset := (page - 1) * limit
			query := "SELECT id, name, surname FROM users ORDER BY id LIMIT $1 OFFSET $2;"

			rows, err := db.Query(query, limit, offset)
			if err != nil {
				t.Errorf("Query execution failed with error: %v", err)
			}
			defer rows.Close()

			var usersData model.UsersData
			for rows.Next() {
				var userData model.UserData
				if err := rows.Scan(&userData.ID, &userData.UserGenData.Name, &userData.UserGenData.Surname); err != nil {
					t.Errorf("Query Scanner failed with error: %v", err)
				}
				usersData = append(usersData, userData)
			}

			if !reflect.DeepEqual(usersData, got) {
				t.Errorf("expected data: %v, got %v", usersData, got)
			}
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				data, err := handler.RetrieveUsers(tc.page, tc.limit)

				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
				if tc.page != 0 {
					assertUserData(t, tc.page, tc.limit, data.Users)

				}
			})
		}
	})
}
