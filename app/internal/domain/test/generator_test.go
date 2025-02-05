package test

import (
	"testing"

	"github.com/John-Dembaremba/pagination-technics/internal/domain"
)

func TestDataGenHandler(t *testing.T) {
	// Arrange
	testCases := []struct {
		name   string
		GenNum int
	}{
		{
			name:   "happy path",
			GenNum: 100000,
		},
	}
	handler := domain.DataGenHandler{}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			data := handler.Generate(int64(tc.GenNum))

			// Assert
			if tc.GenNum != len(data) {
				t.Errorf("expected data Length of %v, got %v", tc.GenNum, len(data))
			}
		})

	}

}
