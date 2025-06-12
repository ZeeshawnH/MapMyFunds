package openfec

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type FECResult struct {
	CandidateID        string  `json:"candidate_id"`
	ContributionAmount float64 `json:"contribution_receipt_amount"`
	ContributionState  string  `json:"contribution_state"`
	ElectionYear       int     `json:"election_year"`
}

type FECResponse struct {
	Results []FECResult `json:"results"`
}

const apiKey = "mpes9XAfrLNioHVlF4mMflhFi1Kd8kfuZAiI4CFC"

func FetchContributionsByState() ([]FECResult, error) {
	// Make API request to FEC API for contributions by state
	url := fmt.Sprintf("https://api.open.fec.gov/v1/presidential/contributions/by_state/?page=1&per_page=100&election_year=2024&sort=-contribution_receipt_amount&sort_hide_null=false&sort_null_only=false&sort_nulls_last=false&api_key=%s", apiKey)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	fmt.Println(string(body))
	if err != nil {
		return nil, err
	}

	var response FECResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return response.Results, err
}
