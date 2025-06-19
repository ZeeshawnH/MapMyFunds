package main

import (
	"log"
	"net/http"

	"backend/aggregator"
	"backend/db"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	router := gin.Default()

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

	client, err := db.ConnectDB()
	// Get request for contributions by state and year
	router.GET("/api/contributions/byState", func(c *gin.Context) {
		resp, err := db.GetContributionsByStateAndYear(client, 2024, "NC")
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
		c.JSON(http.StatusOK, resp)
	})

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	db.PopulateDatabase(client)

	router.Run(":8080")

}
