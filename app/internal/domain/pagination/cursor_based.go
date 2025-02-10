package pagination

import (
	"database/sql"

	"github.com/John-Dembaremba/pagination-technics/internal/model"
	"github.com/John-Dembaremba/pagination-technics/internal/repo"
)

type cursoBasedRepoInterface interface {
	CursorBasedRead(cursor, limit int) (model.UsersData, error)
	TotalUsers() (int, error)
}

type CursorBasedHandler struct {
	Repo cursoBasedRepoInterface
}

func NewCursorBasedHandler(db *sql.DB) CursorBasedHandler {
	repoHandler := repo.RepositoryHandler{Db: db}
	return CursorBasedHandler{
		Repo: repoHandler,
	}
}

func (h CursorBasedHandler) Retrieve(cursor, limit int) (model.UsersCursorBasedMetaData, error) {
	var pgMetaData model.UsersCursorBasedMetaData

	usersData, err := h.Repo.CursorBasedRead(cursor, limit)
	if err != nil {
		return pgMetaData, err
	}

	nextCursor := usersData[len(usersData)-1].ID
	pgMetaData.Users = usersData
	pgMetaData.NextCursor = nextCursor
	return pgMetaData, nil
}
