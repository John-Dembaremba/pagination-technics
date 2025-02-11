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

// NewCursorBasedHandler initializes a CursorBasedHandler with a database connection.
func NewCursorBasedHandler(db *sql.DB) CursorBasedHandler {
	repoHandler := repo.RepositoryHandler{Db: db}
	return CursorBasedHandler{
		Repo: repoHandler,
	}
}

// Retrieve fetches a paginated list of users using cursor-based pagination.
// It returns a UsersCursorBasedMetaData struct containing the retrieved users
// and the next cursor for subsequent queries.
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
