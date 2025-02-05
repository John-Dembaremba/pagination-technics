package domain

import (
	"log"

	"github.com/John-Dembaremba/pagination-technics/internal/model"
	"github.com/John-Dembaremba/pagination-technics/pkg"
)

type DataGenHandler struct{}

func (DataGenHandler) Generate(num int64) []model.UserGenData {
	var data []model.UserGenData
	log.Println("Starting generating data ...")
	for range num {
		fakerHandler := pkg.NewUserGenData()
		data = append(data, fakerHandler)
	}

	log.Println("Completed generating data.")
	return data
}
