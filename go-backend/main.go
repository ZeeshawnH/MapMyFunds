package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"backend/aggregator"
	"backend/db"
	"backend/openfec"
	"backend/types"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
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
		yearParam := c.Query("year")
		year := 2024
		if yearParam != "" {
			if parsed, err := strconv.Atoi(yearParam); err == nil {
				year = parsed
			}
		}

		resp, err := db.GetAllContributionsWithCandidates(client, year)
		fmt.Println(resp)
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
		}
		c.JSON(http.StatusOK, resp)
	})

	// Batch update candidates' image_url. Accepts a list of full candidate objects
	// in the body, but only uses candidate_id and image_url for each update.
	router.PUT("/api/candidates/image", func(c *gin.Context) {
		var payload []types.Candidate
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
			return
		}

		if len(payload) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "empty candidate list"})
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
		defer cancel()

		collection := client.Database("election_data").Collection("candidates")
		var modified int64
		var notFound []string
		var invalid []string

		for _, cand := range payload {
			if cand.CandidateID == "" {
				invalid = append(invalid, cand.CandidateName)
				continue
			}

			filter := bson.M{"candidate_id": cand.CandidateID}
			update := bson.M{"$set": bson.M{"image_url": cand.ImageURL}}

			result, err := collection.UpdateOne(ctx, filter, update)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			if result.MatchedCount == 0 {
				notFound = append(notFound, cand.CandidateID)
				continue
			}

			modified += result.ModifiedCount
		}

		c.JSON(http.StatusOK, gin.H{
			"status":   "ok",
			"modified": modified,
			"notFound": notFound,
			"invalid":  invalid,
		})
	})

	router.Run(":8080")
}
