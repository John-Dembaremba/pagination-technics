package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/John-Dembaremba/pagination-technics/internal/api"
	"github.com/John-Dembaremba/pagination-technics/internal/domain"
	"github.com/John-Dembaremba/pagination-technics/internal/domain/pagination"
	"github.com/John-Dembaremba/pagination-technics/internal/repo"
	"github.com/John-Dembaremba/pagination-technics/pkg"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func main() {

	log.Println("Setting Env Variables ...")
	env := pkg.NewEnv()

	db, err := pkg.NewPgDb(env.POSTGRES_CONTAINER_NAME, env.POSTGRES_DB, env.POSTGRES_USER, env.POSTGRES_PSW, env.POSTGRES_PORT)
	if err != nil {
		log.Fatalf("failed to init database with error: %v", err)
	}

	migration_query, err := pkg.ReadFile("./schema.sql")
	if err != nil {
		log.Fatalf("failed to read sql schema with error: %v", err)
	}

	if err = pkg.RunMigration(db, migration_query); err != nil {
		log.Fatalf("failed to run migration with error: %v", err)
	}
	log.Println("Migration completed successfully.")

	repo := repo.RepositoryHandler{Db: db}

	numUsers := 1000

	log.Printf("Seeding %v of users", numUsers)

	seedH := domain.SeedHandler{
		Generator: domain.DataGenHandler{},
		Repo:      repo,
	}
	if err = seedH.Seed(numUsers); err != nil {
		log.Fatalf("Seeding failed with error: %v", err)
	}
	log.Println("Seeding completed.")

	log.Printf("Starting Server on port: %v\n", env.ServerPort)
	defer log.Println("--------------------------")

	mux := http.NewServeMux()

	log.Println("Init Prometheus Metrics http handler ....")
	promHttpH := pkg.NewPromMetricsHttpHandler()
	mux.Handle("/metrics", promHttpH)
	log.Println("Prometheus Metrics http handler set")

	// Tracing
	log.Println("Initialize Tracer .....")
	tracerHander := pkg.TracerConfigHandler{}
	tp, err := tracerHander.InitTracer(env.OTLP_HTTP_PORT, env.JAEGER_HOST, "pagination-app")
	log.Fatalf("Tracer init failed with error: %v", err)

	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Fatalf("Error shutting down tracer handler: %v", err)
		}
		log.Println("Gracefully shutdown tracer handler")
	}()

	log.Println("Tracer initialized successfully and running ..")

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello Paginators are ready")
	})

	cursorBsdHandler := pagination.NewCursorBasedHandler(db)
	cursorBsdHttpControler := api.NewCursorBasedHttpController(cursorBsdHandler)
	mux.Handle("GET /users/cursor-based",
		otelhttp.NewHandler(
			http.HandlerFunc(cursorBsdHttpControler.GetUsers),
			"cursor-based-pagination",
		))

	limitOffsetHandler := pagination.NewLimitOffSetHandler(db)
	limitOffsetHttpController := api.NewLimitOffsetHttpControler(limitOffsetHandler)

	mux.Handle("GET /users/limit-offset",
		otelhttp.NewHandler(
			http.HandlerFunc(limitOffsetHttpController.GetUsers),
			"cursor-based-pagination",
		))

	http.ListenAndServe(fmt.Sprintf(":%v", env.ServerPort), mux)
}
