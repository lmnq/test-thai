package model

import "time"

type Item struct {
	ID        int       `json:"id"`
	ItemName  string    `json:"item_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}
