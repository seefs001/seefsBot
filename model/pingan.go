package model

type PingAn struct {
	Name        string `json:"name"`
	IDCard      string `json:"id_card"`
	Gender      string `json:"gender"`
	Phone       string `json:"phone"`
	Email       string `json:"email"`
	Province    string `json:"province"`
	City        string `json:"city"`
	MonthInCome string `gorm:"column:month_income;" json:"in_come"`
	IsMarried   string `json:"is_married"`
}

func (PingAn) TableName() string {
	return "pingan"
}
