package model

type FindRequest struct {
	Category string `json:"category,omitempty"`
	Provider string `json:"provider,omitempty"`
	Limit    int    `json:"limit,omitempty"`
	Page     int    `json:"page,omitempty"`
	Sort     string `json:"sort,omitempty"`
	Order    string `json:"order,omitempty"`
}

type FindResponse struct {
	Criteria FindRequest `json:"criteria,omitempty"`
	Articles []Article   `json:"articles,omitempty"`
	Total    int         `json:"total,omitempty"`
}
