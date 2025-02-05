package test

import (
	"testing"

	"github.com/John-Dembaremba/pagination-technics/internal/domain"
	"github.com/John-Dembaremba/pagination-technics/internal/model"
)

type dataGenInterfaceMock struct{}

func (dataGenInterfaceMock) Generate(num int64) []model.UserGenData {
	var d []model.UserGenData
	return d
}

type repoInterfaceMock struct{}

func (repoInterfaceMock) Create(users []model.UserGenData) error {
	return nil
}

func TestSeedHandler(t *testing.T) {
	testCases := []struct {
		name        string
		itemsNum    int
		expectedErr any
	}{}

	handler := domain.SeedHandler{
		Generator: dataGenInterfaceMock{},
		Repo:      repoInterfaceMock{},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

		})
	}
}
