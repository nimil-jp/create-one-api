package response

type SearchResponse struct {
	Count   uint        `json:"count"`
	Records interface{} `json:"records"`
}

func NewSearchResponse(records interface{}, count uint) *SearchResponse {
	return &SearchResponse{
		Count:   count,
		Records: records,
	}
}
