package pagination

import (
	"context"
	"database/sql"

	"github.com/John-Dembaremba/pagination-technics/internal/model"
	"github.com/John-Dembaremba/pagination-technics/internal/repo"
	"github.com/John-Dembaremba/pagination-technics/pkg"
)

type cursoBasedRepoInterface interface {
	CursorBasedRead(ctx context.Context, cursor, limit int) (model.UsersData, error)
	TotalUsers(ctx context.Context) (int, error)
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
func (h CursorBasedHandler) Retrieve(ctx context.Context, cursor, limit int) (model.UsersCursorBasedMetaData, error) {
	// tracer span instance
	tracerHander := pkg.TracerConfigHandler{}
	ctx, span := tracerHander.TracerSpan(ctx, "cursor-domain", "pagination: retrieve")
	defer span.End()

	var pgMetaData model.UsersCursorBasedMetaData

	usersData, err := h.Repo.CursorBasedRead(ctx, cursor, limit)
	if err != nil {
		span.RecordError(err) // Record error in span
		return pgMetaData, err
	}

	nextCursor := usersData[len(usersData)-1].ID
	pgMetaData.Users = usersData
	pgMetaData.NextCursor = nextCursor
	return pgMetaData, nil
}
