package domain

import (
	"context"

	"github.com/John-Dembaremba/pagination-technics/internal/model"
	"github.com/John-Dembaremba/pagination-technics/pkg"
)

type dataGenInterface interface {
	Generate(num int64) []model.UserGenData
}

type repoInterface interface {
	Create(ctx context.Context, users []model.UserGenData) error
}

type SeedHandler struct {
	Generator dataGenInterface
	Repo      repoInterface
}

func (s SeedHandler) Seed(ctx context.Context, itemsNum int) error {
	// tracer span instance
	tracerHander := pkg.TracerConfigHandler{}
	ctx, span := tracerHander.TracerSpan(ctx, "seeding-domain", "domain: Seed")
	defer span.End()

	usersData := s.Generator.Generate(int64(itemsNum))
	err := s.Repo.Create(ctx, usersData)
	return err
}
