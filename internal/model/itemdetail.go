package model

type ItemDetail struct {
	ID         int     `json:"id"`
	ItemID     int     `json:"item_id"`
	CategoryID int     `json:"category_id"`
	GroupID    int     `json:"group_id"`
	Cost       float64 `json:"cost"`
	Price      float64 `json:"price"`
	Sort       int     `json:"sort"`
	CreatedAt  string  `json:"created_at"`
	UpdatedAt  string  `json:"updated_at"`
	DeletedAt  string  `json:"deleted_at"`
}

type ItemDetailView struct {
	ID           int     `json:"id"`
	ItemID       int     `json:"item_id"`
	ItemName     string  `json:"item_name"`
	CategoryID   int     `json:"category_id"`
	CategoryName string  `json:"category_name"`
	GroupID      int     `json:"group_id"`
	GroupName    string  `json:"group_name"`
	Cost         float64 `json:"cost"`
	Price        float64 `json:"price"`
	Sort         int     `json:"sort"`
	CreatedAt    string  `json:"created_at"`
	UpdatedAt    string  `json:"updated_at"`
	DeletedAt    string  `json:"deleted_at"`
}

type ItemDetailFilter struct {
	ID           *int    `json:"id"`
	ItemName     *string `json:"item_name"`
	CategoryName *string `json:"category_name"`
	GroupName    *string `json:"group_name"`
}
