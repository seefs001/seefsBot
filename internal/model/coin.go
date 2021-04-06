package model

import "time"

type Coin struct {
	ID        int64     `gorm:"primaryKey,autoIncrement" json:"id"`
	Type      string    `json:"type"`
	Price     string    `json:"price"`
	Increase  float64   `json:"increase"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
