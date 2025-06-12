package main

import (
	"elections-backend/aggregator"
	"elections-backend/openfec"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		contributions, err := openfec.FetchContributionsByState()
		if err != nil {
			http.Error(w, "Failed to fetch contributions", http.StatusInternalServerError)
			return
		}

		topCandidates, err := aggregator.TopCandidatesByState(contributions)
		if err != nil {
			http.Error(w, "Failed to aggregate contributions", http.StatusInternalServerError)
			return
		}

		fmt.Fprintln(w, topCandidates)
		fmt.Println("Received a ping request")
	})

	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
