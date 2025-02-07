package pkg

import (
	"github.com/John-Dembaremba/pagination-technics/internal/model"
	"github.com/icrowley/fake"
)

func NewUserGenData() model.UserGenData {
	return model.UserGenData{
		Name:    fake.FirstName(),
		Surname: fake.LastName(),
	}
}

func FakeUsersData(n int) model.UsersData {
	var usersData model.UsersData
	for i := 1; i <= n; i++ {
		usersData = append(usersData, model.UserData{
			ID:          i,
			UserGenData: NewUserGenData(),
		})
	}
	return usersData
}
