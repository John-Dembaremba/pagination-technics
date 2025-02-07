package repo

import (
	"database/sql"
	"fmt"

	"github.com/John-Dembaremba/pagination-technics/internal/model"
	"github.com/lib/pq"
)

type RepositoryHandler struct {
	Db *sql.DB
}

// Create inserts multiple UserGenData records into the 'users' table
// using a transaction and a prepared COPY statement for efficient bulk insert.
// Returns an error if any operation fails.
func (r RepositoryHandler) Create(users []model.UserGenData) error {
	// Open transaction
	tx, err := r.Db.Begin()
	if err != nil {
		return fmt.Errorf("failed to open transaction: %v", err)
	}
	// Ensure rollback on failure
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Prepare COPY statement
	stmt, err := tx.Prepare(pq.CopyIn("users", "name", "surname"))
	if err != nil {
		return fmt.Errorf("failed to prepare COPY statement: %v", err)
	}
	defer stmt.Close()

	// Bulk insert users
	for _, user := range users {
		if _, err = stmt.Exec(user.FirstName, user.Surname); err != nil {
			return fmt.Errorf("failed to insert data: %v", err)
		}
	}

	// Flush remaining data
	if _, err = stmt.Exec(); err != nil {
		return fmt.Errorf("failed to flush data: %v", err)
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

func (r RepositoryHandler) LimitOffsetRead(offset, limit int) {

}
