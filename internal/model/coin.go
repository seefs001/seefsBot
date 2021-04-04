package model

import "time"

type Coin struct {
	ID        int64     `gorm:"primarykey,autoIncrement" json:"id,omitempty"`
	Type      string    `json:"type"`
	Price     string    `json:"price"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
