package model

type UserProperty struct {
	Id          int64   `gorm:"column:id; primaryKey" json:"id"`
	UserId      int64   `gorm:"column:user_id; not null; default:0;" json:"userId"`               //用户ID
	TokenName   string  `gorm:"column:token_name; type:varchar(255); default:0;"json:"tokenName"` //代币名称
	TokenNumber float64 `gorm:"column:token_number;  not null; default:0;"json:"tokenNumber"`     //代币金额
	Time        string  `gorm:"column:time;  type:varchar(255); default:0;"json:"createTA"`       //时间
}

func (*UserProperty) TableName() string {
	return "userProperty"
}
