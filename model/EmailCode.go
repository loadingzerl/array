package model

type EmailCode struct {
	Id        int64  `gorm:"column:id; primaryKey" json:"id"`
	Email     string `gorm:"column:email; type:varchar(255); default:0;"json:"email"`          //邮箱
	EmailCode string `gorm:"column:email_code; type:varchar(255); default:0;"json:"emailCode"` //验证码
	Time      string `gorm:"column:time;  type:varchar(255); default:0;"json:"createTA"`       //时间
}

func (*EmailCode) TableName() string {
	return "emailCode"
}
