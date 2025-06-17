package main

import (
	"log"
	"net/http"

	"backend/aggregator"

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

	router.Run(":8080")

}
