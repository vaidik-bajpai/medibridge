package dto

// Paginate represents the parameters used for pagination in the API.
// swagger:parameters paginateRequest
type Paginate struct {
	// SearchTerm is the term used to filter the results.
	// It is optional, and if provided, it will be used to search within the data.
	// example: "john doe"
	SearchTerm string `json:"searchTerm"`

	// PageSize defines the number of items to be returned per page.
	// required: true
	// example: 10
	PageSize int `json:"pageSize"`

	// LastID is the ID of the last item from the previous page, used for pagination.
	// It helps in retrieving the next set of results.
	// example: "123abc"
	LastID string `json:"lastID"`
}
