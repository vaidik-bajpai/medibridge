package models

// Paginate represents the parameters used for pagination in the API.
// swagger:parameters paginateRequest
type Paginate struct {
	// PageSize defines the number of items to be returned per page.
	// required: true
	// example: 10
	PageSize int64 `json:"pageSize" validate:"required,numeric,gt=0"`

	// LastID is the ID of the last item from the previous page, used for pagination.
	// It helps in retrieving the next set of results.
	// example: "123abc"
	LastID string `json:"lastID" validate:"omitempty,uuid"`
}
