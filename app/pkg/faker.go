package pkg

import (
	"github.com/John-Dembaremba/pagination-technics/internal/model"
	"github.com/icrowley/fake"
)

func NewUserGenData() model.UserGenData {
	return model.UserGenData{
		FirstName: fake.FirstName(),
		Surname:   fake.LastName(),
	}
}
