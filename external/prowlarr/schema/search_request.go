package schema

type SearchRequest struct {
	Query string `json:"query"`
}

type SearchByTypeRequest struct {
	Query string `json:"query"`
	Type  string `json:"type"`
}
