package main

import (
	"backend/aggregator"
	"backend/db"
	"backend/ingest"
	"backend/storage/postgres"
	"context"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/joho/godotenv"
)

// Entry point for the ETL/ingestion pipeline service.
// Currently empty; wire up OpenFEC → Postgres → Mongo stages here later.
func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	// Log to a dated file under ./logs and also to stdout so
	// this ETL service can run in the background and be inspected later.
	if err := setupLogging(); err != nil {
		log.Printf("Warning: failed to set up file logging: %v", err)
	}

	ctx := context.Background()
	_ = ctx

	client, err := db.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	pgConn, err := postgres.Connect(ctx)
	if err != nil {
		log.Fatalf("Failed to connect to Postgres: %v", err)
	}
	defer pgConn.Close(ctx)

	var one int
	err = pgConn.QueryRow(ctx, "SELECT 1").Scan(&one)
	if err != nil {
		log.Printf("Postgres test query failed: %v", err)
	} else {
		log.Println("Postgres connection verified")
	}

	repo := postgres.NewRepository(pgConn)

	err = ingest.IngestCandidateInfo(ctx, repo, 2024)
	if err != nil {
		log.Printf("Ingestion failed: %v", err)
	}

	err = ingest.IngestCommitteeInfo(ctx, repo)
	if err != nil {
		log.Printf("Committee ingestion failed: %v", err)
	}

	err = ingest.IngestCandidateCommitteeRelations(ctx, repo)
	if err != nil {
		log.Printf("Candidate-Committee relation ingestion failed: %v", err)
	}

	err = ingest.RunScheduleAIngestion(ctx, repo, 2024, 0)
	if err != nil {
		log.Printf("Ingestion failed: %v", err)
	}

	if true {
		aggregator := aggregator.NewScheduleAAggregator(repo, client, "elections")
		err = aggregator.RunAggregation(ctx, 2024)
		if err != nil {
			log.Printf("Aggregation failed: %v", err)
		}
	}

	log.Println("ETL service starting (no jobs configured yet)...")
}

// setupLogging configures log output to both stdout and a dated file under ./logs.
func setupLogging() error {
	if err := os.MkdirAll("logs", 0o755); err != nil {
		return err
	}

	filename := time.Now().Format("2006-01-02") + "-etl.log"
	path := filepath.Join("logs", filename)
	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}

	mw := io.MultiWriter(os.Stdout, file)
	log.SetOutput(mw)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	log.Printf("logging ETL output to %s", path)
	return nil
}
