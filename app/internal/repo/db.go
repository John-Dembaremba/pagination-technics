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
		if _, err = stmt.Exec(user.Name, user.Surname); err != nil {
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

func (r RepositoryHandler) LimitOffsetRead(offset, limit int) (model.UsersData, error) {
	query := `SELECT id, name, surname FROM users ORDER BY id LIMIT $1 OFFSET $2;`

	var usersData model.UsersData
	rows, err := r.Db.Query(query, limit, offset)
	if err != nil {
		return usersData, fmt.Errorf("LimitOffsetRead query exec failed with error: %v", err)
	}

	defer rows.Close()
	for rows.Next() {
		var userData model.UserData
		if err := rows.Scan(&userData.ID, &userData.Name, &userData.Surname); err != nil {
			return usersData, fmt.Errorf("LimitOffsetRead query scan failed with error: %v", err)
		}
		usersData = append(usersData, userData)
	}

	return usersData, nil
}

func (r RepositoryHandler) TotalUsers() (int, error) {
	var count int
	if err := r.Db.QueryRow("SELECT COUNT(id) FROM users").Scan(&count); err != nil {
		return count, fmt.Errorf("TotalUsers query exec failed with error: %v", err)
	}
	return count, nil
}
