package model

import "time"

type Category struct {
	ID           int        `json:"id"`
	CategoryName string     `json:"category_name"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at"`
}
