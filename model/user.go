package model

type User struct {
	ID            int     `gorm:"primaryKey" json:"id"`
	UserName      string  `gorm:"column:username;index:user_username,unique" json:"username"`
	Score         int64   `json:"score"`
	FreeScore     int64   `json:"free_score"`
	InviteCode    string  `gorm:"index:user_invite_code,unique" json:"invite_code"`
	BeInvitedCode *string `json:"be_invited_code"`
}
