package models

// Paginate represents the parameters used for pagination in the API.
// swagger:parameters paginateRequest
type Paginate struct {
	// PageSize defines the number of items to be returned per page.
	// required: true
	// example: 10
	PageSize int64 `json:"pageSize" validate:"required,gt=0"`

	// Offset is the number of items to skip before starting to collect the result set.
	// required: true
	// example: 20
	Page int64 `json:"page" validate:"required,gte=1"`

	// SearchTerm is a keyword or phrase used to filter the results.
	// This is optional.
	// example: "john doe"
	SearchTerm string `json:"searchTerm"`
}
