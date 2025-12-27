package openfec

import (
	"backend/types"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

const (
	contributorReceiptsPath        = "schedules/schedule_c/"
	contributorReceiptsQueryParams = "per_page=100&sort=-incurred_date&sort_hide_null=true&sort_null_only=false&sort_nulls_last=true"
)

func FetchContributorReceiptDataFromFEC(
	lastIndex *string,
	lastReceiptDate *string,
	sortNullOnly *bool,
	years []int,
) (FECIndexedResponse[types.ContributorReceipt], error) {

	params := url.Values{}

	if lastIndex != nil {
		params.Set("last_index", *lastIndex)
	}
	if lastReceiptDate != nil {
		params.Set("last_contribution_receipt_date", *lastReceiptDate)
	}
	if sortNullOnly != nil {
		params.Set("sort_null_only", strconv.FormatBool(*sortNullOnly))
	}

	for _, y := range years {
		params.Add("years", strconv.Itoa(y))
	}

	var requestUrl string

	if contributorReceiptsQueryParams != "" {
		paramsStr := params.Encode() + "&" + contributorReceiptsQueryParams
		requestUrl = fmt.Sprintf("%s%s?%s&api_key=%s", baseURL, contributorReceiptsPath, paramsStr, os.Getenv("FEC_API_KEY"))
	} else {
		requestUrl = fmt.Sprintf("%s%s?%s&api_key=%s", baseURL, contributorReceiptsPath, params.Encode(), os.Getenv("FEC_API_KEY"))
	}

	fmt.Println(requestUrl)

	resp, err := http.Get(requestUrl)
	if err != nil {
		return FECIndexedResponse[types.ContributorReceipt]{}, fmt.Errorf("FEC API request failed: %w", err)
	}
	defer resp.Body.Close()

	var data FECIndexedResponse[types.ContributorReceipt]
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return FECIndexedResponse[types.ContributorReceipt]{}, fmt.Errorf("failed to parse response: %w", err)
	}
	return data, nil
}

