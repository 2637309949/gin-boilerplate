package http

type Pager struct {
	CurPage    int32       `json:"cur_page"`
	TotalPage  int64       `json:"total_page"`
	TotalCount int64       `json:"total_count"`
	Data       interface{} `json:"data"`
}
