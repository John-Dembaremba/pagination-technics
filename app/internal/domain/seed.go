package domain

import (
	"github.com/John-Dembaremba/pagination-technics/internal/model"
)

type dataGenInterface interface {
	Generate(num int64) []model.UserGenData
}

type repoInterface interface {
	Create(users []model.UserGenData) error
}

type SeedHandler struct {
	Generator dataGenInterface
	Repo      repoInterface
}

func (s SeedHandler) Seed(itemsNum int) error {
	usersData := s.Generator.Generate(int64(itemsNum))
	err := s.Repo.Create(usersData)
	return err
}
