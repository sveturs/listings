package service

// SearchQuery представляет поисковый запрос в системе
type SearchQuery struct {
	ID              int    `json:"id"`
	Query           string `json:"query"`
	NormalizedQuery string `json:"normalized_query"`
	SearchCount     int    `json:"search_count"`
	LastSearched    string `json:"last_searched"`
	Language        string `json:"language"`
	ResultsCount    int    `json:"results_count"`
}
