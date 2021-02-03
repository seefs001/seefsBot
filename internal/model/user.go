package model

import (
	"github.com/seefs001/seefsBot/pkg/orm"
	"time"
)

type User struct {
	ID        int64     `gorm:"primarykey" json:"id,omitempty"`
	Role      int       `json:"role"`
	SecretKey string    `json:"secret_key"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

const (
	NormalRole = iota
	AdminRole
	BlackListRole
)

func GetUserIDBySecretKey(secretKey string) (int64, error) {
	user := User{}
	err := orm.DB().Model(&User{}).
		Where("secret_key = ?", secretKey).
		First(&user).Error

	if err != nil {
		return 0, err
	}
	return user.ID, nil
}
