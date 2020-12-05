package model

type QQ struct {
	Phone string `json:"phone"`
	QQ    string `json:"qq"`
}

func (QQ) TableName() string {
	return "qq"
}
