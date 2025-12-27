package openfec

const (
	baseURL = "https://api.open.fec.gov/v1/"
)

type FECPagination struct {
	Count        int  `json:"count"`
	IsCountExact bool `json:"is_count_exact"`
	Page         int  `json:"page"`
	Pages        int  `json:"pages"`
	PerPage      int  `json:"per_page"`
}

// Used for indexed pagination
type FECLastIndexes struct {
	LastIndex                   string `json:"last_index"`
	LastContributionReceiptDate string `json:"last_contribution_receipt_date,omitempty"`
	SortNullOnly                bool   `json:"sort_null_only,omitempty"`
}

// Used for indexed pagination
type FECIndexedPagination struct {
	Count        int  `json:"count"`
	IsCountExact bool `json:"is_count_exact"`
	PerPage      int  `json:"per_page"`
	Pages        int  `json:"pages"`

	LastIndexes FECLastIndexes `json:"last_indexes"`
}

type FECResponse[T any] struct {
	Pagination FECPagination `json:"pagination"`
	Results    []T           `json:"results"`
}

// Used for indexed pagination
type FECIndexedResponse[T any] struct {
	Pagination FECIndexedPagination `json:"pagination"`
	Results    []T                  `json:"results"`
}


type FECError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func getDataPaginated[T any](page int, year int, fetcher func(page int, year int) (FECResponse[T], error)) ([]T, error) {
	data, err := fetcher(page, year)
	if err != nil {
		return nil, err
	}
	return data.Results, nil
}
