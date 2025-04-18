package pkg

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/docker/go-connections/nat"
	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

type DbAttributes struct {
	DbName     string
	DbUserName string
	DbPassword string
	MappedPort string
}

// PgTestContainerSetup initializes a PostgreSQL test container with the given attributes.
func (d *DbAttributes) PgTestContainerSetup(ctx context.Context) *postgres.PostgresContainer {
	env := NewEnv()
	pgContainerName := fmt.Sprintf("postgres:%v", env.POSTGRES_VERSION)
	container, err := postgres.Run(
		ctx,
		pgContainerName,
		postgres.WithDatabase(d.DbName),
		postgres.WithUsername(d.DbUserName),
		postgres.WithPassword(d.DbPassword),
		// postgres.WithWaitStrategy(wait),
	)

	if err != nil {
		log.Fatalf("Failed to start PostgreSQL container: %v", err)
	}

	return container
}

// DbSetup establishes a database connection and applies schema migrations.
func (d *DbAttributes) DbSetup(ctx context.Context, ctr *postgres.PostgresContainer, schemaFileDir string) *sql.DB {
	host, err := ctr.Host(ctx)
	if err != nil {
		log.Fatalf("Failed to get container host: %v", err)
	}

	port, err := ctr.MappedPort(ctx, nat.Port(d.MappedPort))
	if err != nil {
		log.Fatalf("Failed to get container port: %v", err)
	}

	// postgres://user:password@localhost:5432/mydatabase?sslmode=disable
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		d.DbUserName,
		d.DbPassword,
		host,
		port.Port(),
		d.DbName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect db: %v", err)
	}

	time.Sleep(5 * time.Second)
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping the database: %v", err)
	}

	log.Println("Db setup completed.")

	query, err := ReadFile(schemaFileDir)
	if err != nil {
		log.Fatalf("Failed to read sql file: %v", err)
	}

	if err := RunMigration(db, query); err != nil {
		log.Fatalf("Failed to run migration: %v", err)
	}

	log.Println("Db migration completed.")
	return db
}

// TearDown gracefully shuts down the test database and closes connections.
func TearDown(db *sql.DB, ctr *postgres.PostgresContainer) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recover from panic TearDown: %v", r)
		}
	}()

	if ctr != nil {
		if err := ctr.Terminate(context.Background()); err != nil {
			log.Fatalf("Failed to terminate container: %v", err)
		}
		log.Println("Container terminated successfully")
	}

	if db != nil {
		if err := db.Close(); err != nil {
			log.Fatalf("Failed to close db connection: %v", err)
		}

		log.Println("Db connection closed successfully")
	}
}

func TestGetMd5Hash(t *testing.T) {
	got := getMd5("pagy", "myverystrongpassword")
	expected := "md5ac77bfe847b783150cc181043bd7d2d7"
	if got != expected {
		t.Errorf("Expected hash to be %s, got %s", expected, got)
	}
}
