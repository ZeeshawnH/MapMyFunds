package main

import (
	"log"
	"net/http"

	"backend/openfec"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	router := gin.Default()

	// Simple get request
	router.GET("/ping", func(c *gin.Context) {
		resp, err := openfec.GetContributions(2024)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		}
		c.JSON(http.StatusOK, resp)
	})

	router.Run(":8080")

}
