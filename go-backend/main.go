package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"

	"backend/aggregator"
	"backend/db"
	"backend/ingest"
	"backend/openfec"
	"backend/storage/postgres"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var debugMode bool

func init() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	flag.BoolVar(&debugMode, "debug", false, "Enable debug mode")
	flag.Parse()

	if debugMode {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(os.Getenv("GIN_MODE"))
	}
}

func main() {
	router := gin.Default()
	router.Use(cors.Default())

	client, err := db.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	ctx := context.Background()

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

	err = ingest.IngestCommitteeInfo(ctx, repo)
	if err != nil {
		log.Printf("Committee ingestion failed: %v", err)
	}

	// err = ingest.IngestCandidateInfo(ctx, repo, 2024)
	// if err != nil {
	// 	log.Printf("Ingestion failed: %v", err)
	// }

	// err = ingest.RunScheduleAIngestion(ctx, repo, 2024, 0)
	// if err != nil {
	// 	log.Printf("Ingestion failed: %v", err)
	// }

	if true {
		aggregator := aggregator.NewScheduleAAggregator(repo, client, "elections")
		err = aggregator.RunAggregation(ctx, 2024)
		if err != nil {
			log.Printf("Aggregation failed: %v", err)
		}
	}

	// Get request for contributions by state
	router.GET("/api/contributions", func(c *gin.Context) {
		resp, err := aggregator.FetchContributionAmountByStateAndCandidate()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		}
		c.JSON(http.StatusOK, resp)
	})

	// Get request for contributions by state and year
	router.GET("/api/contributions/byState", func(c *gin.Context) {
		resp, err := db.GetContributionsByStateAndYear(client, 2024, "NC")
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
		c.JSON(http.StatusOK, resp)
	})

	router.GET("api/candidates", func(c *gin.Context) {
		resp, err := openfec.GetCandidateData(2024)
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
		}

		c.JSON(http.StatusOK, resp)
	})

	router.GET("/api/contributions/withCandidates", func(c *gin.Context) {
		resp, err := db.GetAllContributionsWithCandidates(client, 2024)
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
		}
		c.JSON(http.StatusOK, resp)
	})

	router.GET("/api/test", func(c *gin.Context) {
		resp, err := openfec.FetchCommitteeDataFromFEC(1)
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
		}
		c.JSON(http.StatusOK, resp.Results)
	})

	router.Run(":8080")
}
