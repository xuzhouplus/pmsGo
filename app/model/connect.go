package model

import "time"

type Connect struct {
	ID        int       `gorm:"private_key" json:"id"`
	AdminId   int       `json:"admin_id"`
	Type      string    `json:"type"`
	Avatar    string    `json:"avatar"`
	Account   string    `json:"account"`
	UnionId   string    `json:"union_id"`
	Status    int       `json:"status"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
