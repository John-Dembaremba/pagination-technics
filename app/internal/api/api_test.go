package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"testing"

	"github.com/John-Dembaremba/pagination-technics/internal/domain"
	"github.com/John-Dembaremba/pagination-technics/internal/domain/pagination"
	"github.com/John-Dembaremba/pagination-technics/internal/model"
	"github.com/John-Dembaremba/pagination-technics/internal/repo"
	"github.com/John-Dembaremba/pagination-technics/pkg"
)

func TestGetUsers(t *testing.T) {
	dbAttributes := pkg.DbAttributes{
		DbName:     "pagination-app",
		DbUserName: "user",
		DbPassword: "mypassword",
		MappedPort: "5432",
	}

	ctx := context.Background()
	testContainer := dbAttributes.PgTestContainerSetup(ctx)
	schemaFile := "../../../pkg/schema.sql"
	db := dbAttributes.DbSetup(ctx, testContainer, schemaFile)
	defer pkg.TearDown(db, testContainer)

	repo := repo.RepositoryHandler{Db: db}
	seedH := domain.SeedHandler{
		Generator: domain.DataGenHandler{},
		Repo:      repo,
	}
	if err := seedH.Seed(1000); err != nil {
		log.Fatalf("Seeding failed with error: %v", err)
	}

	t.Run("cursor based", func(t *testing.T) {
		testCases := []struct {
			name         string
			cursor       string
			limit        string
			code         int
			expectedResp model.ResponseMeta
			isSuccess    bool
		}{
			{
				name:   "invalid cursor",
				cursor: "invalid cursor",
				limit:  "20",
				code:   400,
				expectedResp: model.ResponseMeta{
					Error:   "invalid cursor param",
					Success: "",
					Data:    model.UsersCursorBasedMetaData{},
				},
				isSuccess: false,
			},
			{
				name:   "invalid limit",
				cursor: "50",
				limit:  "invalid-limit",
				code:   400,
				expectedResp: model.ResponseMeta{
					Error:   "invalid limit param",
					Success: "",
					Data:    model.UsersCursorBasedMetaData{},
				},
				isSuccess: false,
			},
			{
				name:   "success",
				cursor: "50",
				limit:  "10",
				code:   200,
				expectedResp: model.ResponseMeta{
					Error:   "",
					Success: "retrieved successfully",
					Data:    model.UsersCursorBasedMetaData{},
				},
				isSuccess: true,
			},
		}

		handler := pagination.NewCursorBasedHandler(db)
		httpController := CursorBasedHttpController{Handler: handler}

		assertUserData := func(t testing.TB, got interface{}, cursor, limit int) {
			t.Helper()

			query := "SELECT id, name, surname FROM users WHERE id < $1 ORDER BY id DESC LIMIT $2;"

			var usersData model.UsersData
			rows, err := db.Query(query, cursor, limit)
			if err != nil {
				t.Errorf("helper func query scan failed with error: %v", err)
			}

			defer rows.Close()
			for rows.Next() {
				var userData model.UserData
				if err := rows.Scan(&userData.ID, &userData.UserGenData.Name, &userData.UserGenData.Surname); err != nil {
					t.Errorf("helper func query exec failed with error: %v", err)
				}
				usersData = append(usersData, userData)
			}

			var payload model.UsersCursorBasedMetaData
			dataBytes, err := json.Marshal(got)
			if err != nil {
				t.Fatalf("failed to marshal Data: %v", err)
			}

			if err := json.Unmarshal(dataBytes, &payload); err != nil {
				t.Fatalf("failed to unmarshal Data: %v", err)
			}

			if !reflect.DeepEqual(payload.Users, usersData) {
				t.Errorf("helper func - expected data: %v, got %v", usersData, payload.Users)
			}

		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				url := fmt.Sprintf("/users/cursor-based?cursor=%v&limit=%v", tc.cursor, tc.limit)

				req, err := http.NewRequest(http.MethodGet, url, nil)
				if err != nil {
					t.Fatalf("failed to create request with error: %v", err)
				}

				req.Header.Set("Content-Type", "application/json")
				resp := httptest.NewRecorder()
				var payload model.ResponseMeta

				httpController.GetUsers(resp, req)
				decoder := json.NewDecoder(resp.Body)

				if tc.code != resp.Code {
					t.Errorf("expected code: %v, got %v", tc.code, resp.Code)
				}

				if err := decoder.Decode(&payload); err != nil {
					t.Errorf("Error decoding JSON: %v", err)
				}

				if !tc.isSuccess {
					if tc.expectedResp.Error != payload.Error {
						t.Errorf("expected error message: %v, got %v", tc.expectedResp.Error, payload.Error)
					}

					if payload.Success != "" {
						t.Errorf("expected no success message, got %v", payload.Success)
					}
				} else {
					if tc.expectedResp.Success != payload.Success {
						t.Errorf("expected success message: %v, got %v", tc.expectedResp.Success, payload.Success)
					}

					if payload.Error != "" {
						t.Errorf("expected no error message, got %v", payload.Error)
					}

					cursorInt, _ := strconv.Atoi(tc.cursor)
					limitInt, _ := strconv.Atoi(tc.limit)
					assertUserData(t, payload.Data, cursorInt, limitInt)
				}

			})
		}

	})

}
