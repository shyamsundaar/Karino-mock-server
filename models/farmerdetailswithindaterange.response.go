package models

// SuccessListResponse represents the paginated response for farmers
type SuccessListResponse struct {
	Success bool               `json:"success"`
	Data    ListFarmerResponse `json:"data"`
}

// ListFarmerResponse contains the items and paging metadata
type ListFarmerResponse struct {
	Items    []FarmerDetails    `json:"items"`
	Metadata PaginationMetadata `json:"metadata"`
}

// PaginationMetadata holds paging info
type PaginationMetadata struct {
	TotalRecord int `json:"totalRecord"`
	TotalPage   int `json:"totalPage"`
	CurrentPage int `json:"currentPage"`
	Limit       int `json:"limit"`
}
