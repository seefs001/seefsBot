package model

import "gorm.io/gorm"

type Card struct {
	ID        uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Content   string         `gorm:"index:card_content,unique" json:"content"`
	Score     int64          `json:"score"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
