package test

import (
	"errors"
	"testing"

	"github.com/John-Dembaremba/pagination-technics/internal/domain"
	"github.com/John-Dembaremba/pagination-technics/internal/model"
)

type dataGenInterfaceMock struct{}

func (dataGenInterfaceMock) Generate(num int64) []model.UserGenData {
	var d []model.UserGenData
	for range num {
		d = append(d, model.UserGenData{})
	}
	return d
}

type repoInterfaceMock struct{}

func (repoInterfaceMock) Create(users []model.UserGenData) error {
	if len(users) == 10 {
		return errors.New("some database error")
	}
	return nil
}

func TestSeedHandler(t *testing.T) {
	testCases := []struct {
		name        string
		itemsNum    int
		expectedErr any
	}{
		{
			name:        "happy path",
			itemsNum:    1000,
			expectedErr: nil,
		},
		{
			name:        "unhappy path",
			itemsNum:    10, // trigger failure
			expectedErr: "some database error",
		},
	}

	handler := domain.SeedHandler{
		Generator: dataGenInterfaceMock{},
		Repo:      repoInterfaceMock{},
	}

	for _, tc := range testCases {

		t.Run(tc.name, func(t *testing.T) {
			err := handler.Seed(tc.itemsNum)

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
		})
	}
}
