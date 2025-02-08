package pagination

import (
	"database/sql"
	"math"

	"github.com/John-Dembaremba/pagination-technics/internal/model"
	"github.com/John-Dembaremba/pagination-technics/internal/repo"
)

type repoInterface interface {
	LimitOffsetRead(offset, limit int) (model.UsersData, error)
	TotalUsers() (int, error)
}

type LimitOffSetHandler struct {
	Repo repoInterface
}

// NewLimitOffSetHandler initializes a LimitOffSetHandler with the given database connection.
func NewLimitOffSetHandler(db *sql.DB) LimitOffSetHandler {
	repoHandler := repo.RepositoryHandler{Db: db}
	return LimitOffSetHandler{
		Repo: repoHandler,
	}
}

// RetrieveUsers fetches paginated user data based on the given page and limit values.
func (h LimitOffSetHandler) RetrieveUsers(page, limit int) (model.UsersPaginationMetaData, error) {
	if page < 1 {
		page = 1
	}

	var data model.UsersPaginationMetaData
	var pg model.Pagination

	offset := (page - 1) * limit
	usersData, err := h.Repo.LimitOffsetRead(offset, limit)
	if err != nil {
		return data, err
	}

	totalUsers, err := h.Repo.TotalUsers()
	if err != nil {
		return data, err
	}

	totalPages := int(math.Ceil(float64(totalUsers) / float64(limit)))
	nextPage := getNextPage(page+1, totalPages)
	prevPage := getPrevPage(page-1, 1)

	pg.CurrentPage = page
	pg.TotalPages = totalPages
	pg.NextPage = nextPage
	pg.PrevPage = prevPage

	data.Pagination = pg
	data.Users = usersData
	return data, nil
}

// getNextPage returns the next page number, ensuring it does not exceed the total pages.
func getNextPage(currentPage, totalPages int) int {
	if currentPage < totalPages {
		return currentPage
	}
	return totalPages
}

// getPrevPage returns the previous page number, ensuring it does not go below the first page.
func getPrevPage(currentPage, firstPage int) int {
	if currentPage > firstPage {
		return currentPage
	}
	return firstPage
}
