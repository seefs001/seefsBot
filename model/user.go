package model

// User subscriber
//
// TelegramID 用作外键
type User struct {
	ID         int `gorm:"primary_key"`
	UserName   string
	TelegramID int
	State      int `gorm:"DEFAULT:0;"`
}

// FindOrCreateUserByTelegramID find subscriber or init a subscriber by telegram ID
func FindOrCreateUserByTelegramIDAndUserName(telegramID int, username string) (*User, error) {
	var user User
	DB.Where(User{TelegramID: telegramID, UserName: username}).FirstOrCreate(&user)

	return &user, nil
}
