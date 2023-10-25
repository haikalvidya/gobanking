package payload

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Meta    *Pagination `json:"meta,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

type Pagination struct {
	PerPage     int64 `json:"per_page"`
	CurrentPage int64 `json:"current_page"`
	LastPage    int64 `json:"last_page"`
	IsLoadMore  bool  `json:"is_load_more"`
}

type PaginationRequest struct {
	Page    int64 `json:"page"`
	PerPage int64 `json:"perpage"`
	Limit   int64
	Offset  int64
}
