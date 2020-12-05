package model

type SF struct {
	Name     string `json:"name"`
	Phone1   string `json:"phone1"`
	Phone2   string `json:"phone2"`
	Address  string `json:"address"`
	Province string `json:"province"`
	City     string `json:"city"`
}

func (SF) TableName() string {
	return "sf"
}
