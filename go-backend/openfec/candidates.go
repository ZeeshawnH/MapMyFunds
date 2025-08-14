package openfec

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

const (
	candidatePath        = "candidates/"
	candidateQueryParams = "per_page=100&sort=name&sort_hide_null=true&sort_null_only=false&sort_nulls_last=false&api_key=mpes9XAfrLNioHVlF4mMflhFi1Kd8kfuZAiI4CFC"
)

type Candidate struct {
	CandidateID         string `json:"candidate_id"`
	CandidateName       string `json:"name"`
	CandidateOffice     string `json:"office"`
	CandidateOfficeFull string `json:"office_full"`
	CandidatePartyAbbr  string `json:"party"`
	CandidatePartyFull  string `json:"party_full"`
}

func FetchCandidateDataFromFEC(page int, year int) (FECResponse[Candidate], error) {
	url := fmt.Sprintf("%s%s?page=%d&election_year=%d&office=%s&%s&api_key=%s",
		baseURL,
		candidatePath,
		page,
		year,
		"P",
		candidateQueryParams,
		os.Getenv("FEC_API_KEY"),
	)

	resp, err := http.Get(url)
	if err != nil {
		return FECResponse[Candidate]{}, fmt.Errorf("FEC API request failed: %w", err)
	}
	defer resp.Body.Close()

	var data FECResponse[Candidate]
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return FECResponse[Candidate]{}, fmt.Errorf("failed to parse response: %w", err)
	}
	return data, nil
}

func GetCandidateData(year int) ([]Candidate, error) {
	data, err := FetchCandidateDataFromFEC(1, year)
	if err != nil {
		return nil, err
	}

	results := data.Results
	pages := data.Pagination.Pages

	for page := 2; page <= pages; page++ {
		pageResults, err := getDataPaginated(page, year, FetchCandidateDataFromFEC)
		if err != nil {
			return nil, err
		}
		results = append(results, pageResults...)
	}
	return results, nil
}
