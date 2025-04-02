package repo

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/John-Dembaremba/pagination-technics/internal/model"
	"github.com/John-Dembaremba/pagination-technics/pkg"
	"github.com/lib/pq"
)

type RepositoryHandler struct {
	Db *sql.DB
}

// Create inserts multiple UserGenData records into the 'users' table
// using a transaction and a prepared COPY statement for efficient bulk insert.
// Returns an error if any operation fails.
func (r RepositoryHandler) Create(ctx context.Context, users []model.UserGenData) error {
	// tracer span instance
	tracerHander := pkg.TracerConfigHandler{}
	ctx, span := tracerHander.TracerSpan(ctx, "create-users-repo", "repo: Create")
	defer span.End()

	// Open transaction
	tx, err := r.Db.Begin()
	if err != nil {
		errTrans := fmt.Errorf("failed to open transaction: %v", err)
		span.RecordError(errTrans)
		return errTrans
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
		errPrepSt := fmt.Errorf("failed to prepare COPY statement: %v", err)
		span.RecordError(errPrepSt)
		return errPrepSt
	}
	defer stmt.Close()

	// Bulk insert users
	for _, user := range users {
		if _, err = stmt.Exec(user.Name, user.Surname); err != nil {
			errInsert := fmt.Errorf("failed to insert data: %v", err)
			span.RecordError(errInsert)
			return errInsert
		}
	}

	// Flush remaining data
	if _, err = stmt.Exec(); err != nil {
		errFlush := fmt.Errorf("failed to flush data: %v", err)
		span.RecordError(errFlush)
		return errFlush
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		errCommit := fmt.Errorf("failed to commit transaction: %v", err)
		span.RecordError(errCommit)
		return errCommit
	}

	return nil
}

func (r RepositoryHandler) LimitOffsetRead(ctx context.Context, offset, limit int) (model.UsersData, error) {
	// tracer span instance
	tracerHander := pkg.TracerConfigHandler{}
	ctx, span := tracerHander.TracerSpan(ctx, "limit-offset-repo", "repo: LimitOffsetRead")
	defer span.End()

	query := `SELECT id, name, surname FROM users ORDER BY id LIMIT $1 OFFSET $2;`

	var usersData model.UsersData
	rows, err := r.Db.Query(query, limit, offset)

	if err != nil {
		errQueryExec := fmt.Errorf("LimitOffsetReadquery exec failed with error: %v", err)
		span.RecordError(errQueryExec) // Record error in span
		return usersData, errQueryExec
	}

	defer rows.Close()
	for rows.Next() {
		var userData model.UserData
		if err := rows.Scan(&userData.ID, &userData.Name, &userData.Surname); err != nil {
			errQueryScan := fmt.Errorf("LimitOffsetReadquery scan failed with error: %v", err)
			span.RecordError(errQueryScan) // Record error in span
			return usersData, errQueryScan
		}
		usersData = append(usersData, userData)
	}

	return usersData, nil
}

func (r RepositoryHandler) TotalUsers(ctx context.Context) (int, error) {
	// tracer span instance
	tracerHander := pkg.TracerConfigHandler{}
	ctx, span := tracerHander.TracerSpan(ctx, "total-users-repo", "repo: TotalUsers")
	defer span.End()

	var count int
	if err := r.Db.QueryRow("SELECT COUNT(id) FROM users").Scan(&count); err != nil {
		errQueryExec := fmt.Errorf("TotalUsers query exec failed with error: %v", err)
		span.RecordError(errQueryExec)
		return count, errQueryExec
	}
	return count, nil
}

func (r RepositoryHandler) CursorBasedRead(ctx context.Context, cursor, limit int) (model.UsersData, error) {
	if cursor <= 1 {
		return initCursor(ctx, limit, r.Db)
	}
	return actualCursor(ctx, cursor, limit, r.Db)
}

// handles only cursor less or equal 1
func initCursor(ctx context.Context, limit int, db *sql.DB) (model.UsersData, error) {
	// tracer span instance
	tracerHander := pkg.TracerConfigHandler{}
	ctx, span := tracerHander.TracerSpan(ctx, "cursor-repo", "repo: initCursor")
	defer span.End()

	var usersData model.UsersData
	query := "SELECT id, name, surname FROM users ORDER BY id DESC LIMIT $1;"
	rows, err := db.Query(query, limit)
	if err != nil {
		errQueryExec := fmt.Errorf("CursorBasedRead-initCursor query exec failed with error: %v", err)
		span.RecordError(errQueryExec) // Record error in span
		return usersData, errQueryExec
	}

	defer rows.Close()
	for rows.Next() {
		var userData model.UserData
		if err := rows.Scan(&userData.ID, &userData.UserGenData.Name, &userData.UserGenData.Surname); err != nil {
			errQueryScan := fmt.Errorf("CursorBasedRead-initCursor query scan failed with error: %v", err)
			span.RecordError(errQueryScan) // Record error in span
			return usersData, errQueryScan
		}

		usersData = append(usersData, userData)

	}
	return usersData, nil
}

// handles only cursor greater than 1
func actualCursor(ctx context.Context, cursor, limit int, db *sql.DB) (model.UsersData, error) {
	// tracer span instance
	tracerHander := pkg.TracerConfigHandler{}
	ctx, span := tracerHander.TracerSpan(ctx, "cursor-repo", "repo: actualCursor")
	defer span.End()

	var usersData model.UsersData

	query := "SELECT id, name, surname FROM users WHERE id < $1 ORDER BY id DESC LIMIT $2;"
	rows, err := db.Query(query, cursor, limit)
	if err != nil {
		errQueryExec := fmt.Errorf("CursorBasedRead-actualCursor query exec failed with error: %v", err)
		span.RecordError(errQueryExec) // Record error in span
		return usersData, errQueryExec
	}

	defer rows.Close()
	for rows.Next() {
		var userData model.UserData
		if err := rows.Scan(&userData.ID, &userData.UserGenData.Name, &userData.UserGenData.Surname); err != nil {
			if err := rows.Scan(&userData.ID, &userData.UserGenData.Name, &userData.UserGenData.Surname); err != nil {
				errQueryScan := fmt.Errorf("CursorBasedRead-actualCursor query scan failed with error: %v", err)
				span.RecordError(errQueryScan) // Record error in span
				return usersData, errQueryScan
			}
		}

		usersData = append(usersData, userData)

	}
	return usersData, nil
}
