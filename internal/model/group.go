package model

type Group struct {
	ID        int    `json:"id"`
	GroupName string `json:"group_name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	DeletedAt string `json:"deleted_at"`
}
