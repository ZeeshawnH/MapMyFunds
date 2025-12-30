package openfec

import (
	"backend/types"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

const (
	committeePath        = "committees/"
	committeeQueryParams = "sort=name&sort_hide_null=false&sort_null_only=false&sort_nulls_last=false"
)

func FetchCommitteeDataFromFEC(page int) (FECResponse[types.Committee], error) {
	url := fmt.Sprintf("%s%s?per_page=100&page=%d&%s&api_key=%s",
		baseURL,
		committeePath,
		page,
		committeeQueryParams,
		os.Getenv("FEC_API_KEY"),
	)

	resp, err := http.Get(url)
	if err != nil {
		return FECResponse[types.Committee]{}, fmt.Errorf("FEC API request failed: %w", err)
	}
	defer resp.Body.Close()

	var data FECResponse[types.Committee]
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return FECResponse[types.Committee]{}, fmt.Errorf("failed to parse response: %w", err)
	}
	return data, nil
}
