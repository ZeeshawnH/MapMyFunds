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

type FECResponse[T any] struct {
	Pagination FECPagination `json:"pagination"`
	Results    []T           `json:"results"`
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
