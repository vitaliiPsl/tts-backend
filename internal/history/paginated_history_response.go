package history

type PaginatedHistoryResponse struct {
	Records      []HistoryRecordDto `json:"records"`
	TotalRecords int              `json:"totalRecords"`
	TotalPages   int                `json:"totalPages"`
	CurrentPage  int                `json:"currentPage"`
	HasMore      bool               `json:"hasMore"`
}
