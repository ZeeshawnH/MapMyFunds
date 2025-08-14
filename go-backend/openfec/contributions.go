package openfec

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

const (
	candidateContributionsPath = "presidential/contributions/by_candidate/"
	contributionsQueryParams   = "per_page=100&sort=-net_receipts&sort_hide_null=false&sort_null_only=false&sort_nulls_last=false"
)

type CandidateContribution struct {
	CandidateID               string  `json:"candidate_id"`
	CandidateLastName         string  `json:"candidate_last_name"`
	CandidatePartyAffiliation string  `json:"candidate_party_affiliation"`
	ContributorState          string  `json:"contributor_state"`
	ElectionYear              int     `json:"election_year"`
	NetReceipts               float64 `json:"net_receipts"`
}

func FetchContributionDataFromFEC(page int, year int) (FECResponse[CandidateContribution], error) {
	url := fmt.Sprintf("%s%s?page=%d&election_year=%d&%s&api_key=%s",
		baseURL,
		candidateContributionsPath,
		page,
		year,
		contributionsQueryParams,
		os.Getenv("FEC_API_KEY"),
	)

	resp, err := http.Get(url)
	if err != nil {
		return FECResponse[CandidateContribution]{}, fmt.Errorf("FEC API request failed: %w", err)
	}
	defer resp.Body.Close()

	var data FECResponse[CandidateContribution]
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return FECResponse[CandidateContribution]{}, fmt.Errorf("failed to parse response: %w", err)
	}
	return data, nil
}

func GetContributions(year int) ([]CandidateContribution, error) {
	data, err := FetchContributionDataFromFEC(1, year)
	if err != nil {
		return nil, err
	}

	results := data.Results
	pages := data.Pagination.Pages

	for page := 2; page <= pages; page++ {
		pageResults, err := getDataPaginated(page, year, FetchContributionDataFromFEC)
		if err != nil {
			return nil, err
		}
		results = append(results, pageResults...)
	}
	return results, nil
}
