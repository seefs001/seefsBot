package model

type Car struct {
	Name        string `json:"name"`
	IDCard      string `json:"id_card"`
	Gender      string `json:"gender"`
	Phone       string `json:"phone"`
	Email       string `json:"email"`
	Province    string `json:"province"`
	City        string `json:"city"`
	Address     string `json:"address"`
	Birthday    string `json:"birthday"`
	Industry    string `json:"industry"`
	PostCode    string `json:"post_code"`
	Education   string `json:"education"`
	MonthInCome string `gorm:"column:month_income;" json:"in_come"`
	IsMarried   string `json:"is_married"`
}

func (Car) TableName() string {
	return "car"
}
