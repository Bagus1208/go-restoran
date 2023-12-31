package model

type MenuResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Category    string `json:"category"`
	Description string `json:"description"`
	Price       int    `json:"price"`
	Image       string `json:"image"`
}

type Favorite struct {
	MenuName   string `json:"menu_name"`
	TotalOrder int    `json:"total_order"`
}

type Pagination struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalItems int64 `json:"total_items"`
}
