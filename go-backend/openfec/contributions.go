package openfec

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

const (
	baseURL     = "https://api.open.fec.gov/v1/presidential/contributions/by_candidate/"
	queryParams = "per_page=100&sort=-net_receipts&sort_hide_null=false&sort_null_only=false&sort_nulls_last=false"
)

type CandidateContribution struct {
	CandidateID               string  `json:"candidate_id"`
	CandidateLastName         string  `json:"candidate_last_name"`
	CandidatePartyAffiliation string  `json:"candidate_party_affiliation"`
	ContributorState          string  `json:"contributor_state"`
	ElectionYear              int     `json:"election_year"`
	NetReceipts               float64 `json:"net_receipts"`
}

type FECPagination struct {
	Count        int  `json:"count"`
	IsCountExact bool `json:"is_count_exact"`
	Page         int  `json:"page"`
	Pages        int  `json:"pages"`
	PerPage      int  `json:"per_page"`
}

type FECResponse struct {
	Pagination FECPagination           `json:"pagination"`
	Results    []CandidateContribution `json:"results"`
}

type FECError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func getFECData(page int, year int) (FECResponse, error) {
	url := fmt.Sprintf("%s?page=%d&election_year=%d&%s&api_key=%s",
		baseURL,
		page,
		year,
		queryParams,
		os.Getenv("FEC_API_KEY"),
	)

	resp, err := http.Get(url)
	if err != nil {
		return FECResponse{}, fmt.Errorf("FEC API request failed: %w", err)
	}
	defer resp.Body.Close()

	var data FECResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return FECResponse{}, fmt.Errorf("failed to parse response: %w", err)
	}
	return data, nil
}

func getContributionsByStatePaginated(page int, year int) ([]CandidateContribution, error) {
	data, err := getFECData(page, year)
	if err != nil {
		return nil, err
	}
	return data.Results, nil
}

func GetContributions(year int) ([]CandidateContribution, error) {
	data, err := getFECData(1, year)
	if err != nil {
		return nil, err
	}

	results := data.Results
	pages := data.Pagination.Pages

	for page := 2; page <= pages; page++ {
		pageResults, err := getContributionsByStatePaginated(page, year)
		if err != nil {
			return nil, err
		}
		results = append(results, pageResults...)
	}
	return results, nil
}
