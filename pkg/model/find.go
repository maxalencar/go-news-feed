package model

type FindRequest struct {
	Category string `json:"category,omitempty" query:"category"`
	Provider string `json:"provider,omitempty" query:"provider"`
	Limit    int    `json:"limit,omitempty" query:"limit"`
	Page     int    `json:"page,omitempty" query:"page"`
	Sort     string `json:"sort,omitempty" query:"sort"`
	Order    string `json:"order,omitempty" query:"order"`
}

type FindResponse struct {
	Criteria FindRequest `json:"criteria,omitempty"`
	Articles []Article   `json:"articles,omitempty"`
	Total    int         `json:"total,omitempty"`
}
