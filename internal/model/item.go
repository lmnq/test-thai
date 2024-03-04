package model

type Item struct {
	ID        int    `json:"id"`
	ItemName  string `json:"item_name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	DeletedAt string `json:"deleted_at"`
}
