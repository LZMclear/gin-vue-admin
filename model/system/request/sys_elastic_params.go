package request

type SearchParams struct {
	Keyword string `json:"keyword"`
	Page    int    `json:"page"`
	Size    int    `json:"size"`
	Sort    string `json:"sort"`
}
