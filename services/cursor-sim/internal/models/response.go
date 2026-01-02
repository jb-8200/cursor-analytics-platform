package models

// PaginatedResponse wraps API responses with pagination metadata.
// This matches the Cursor API response format.
type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Pagination Pagination  `json:"pagination,omitempty"`
	Params     Params      `json:"params,omitempty"`
}

// Pagination contains pagination metadata for list responses.
type Pagination struct {
	Page            int  `json:"page"`
	PageSize        int  `json:"pageSize"`
	TotalUsers      int  `json:"totalUsers,omitempty"`
	TotalPages      int  `json:"totalPages"`
	HasNextPage     bool `json:"hasNextPage"`
	HasPreviousPage bool `json:"hasPreviousPage"`
}

// Params contains the request parameters echoed back in the response.
type Params struct {
	From     string `json:"from,omitempty"`
	To       string `json:"to,omitempty"`
	Page     int    `json:"page,omitempty"`
	PageSize int    `json:"pageSize,omitempty"`
	UserID   string `json:"userId,omitempty"`
	RepoName string `json:"repoName,omitempty"`
}
