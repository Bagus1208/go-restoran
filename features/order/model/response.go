package model

type OrderResponse struct {
	ID      uint                  `json:"id"`
	NoTable int                   `json:"no_table"`
	Orders  []OrderDetailResponse `json:"orders"`
	Total   int                   `json:"total"`
	Status  string                `json:"status"`
}

type OrderDetailResponse struct {
	ID       uint `json:"id"`
	MenuID   int  `json:"menu_id"`
	Quantity int  `json:"quantity"`
}
