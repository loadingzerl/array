package model

type User struct {
	Id                  int64   `gorm:"column:id; primaryKey" json:"id"`
	Email               string  `gorm:"column:email; type:varchar(255); default:0;"json:"email"`                        //邮箱
	PassWord            string  `gorm:"column:pass_word; type:varchar(255); default:0;"json:"passWord"`                 //密码
	DataBarth           string  `gorm:"column:data_barth; type:varchar(255); default:0;"json:"dataBarth"`               //出生日期
	UUid                string  `gorm:"column:uuid; type:varchar(255); default:0;"json:"uuid"`                          //uid
	Name                string  `gorm:"column:name; type:varchar(255); default:0;"json:"name"`                          // 昵称
	HeadPhoto           string  `gorm:"column:head_photo; type:varchar(255); default:0;"json:"headPhoto"`               // 头像
	PersonaSignature    string  `gorm:"column:persona_signature; type:varchar(255); default:0;"json:"personaSignature"` //个性签名
	ETHAddress          string  `gorm:"column:eth_address; type:varchar(255); default:0;"json:"ethAddress"`             //绑定的ETH地址
	Siberian            string  `gorm:"column:siberian; type:varchar(255); default:0;"json:"siberian"`                  // 卡密
	LimitNumber         float64 `gorm:"column:limit_number;  not null; default:0;"json:"limitNumber"`                   // 额度
	FormerlyLimitNumber float64 `gorm:"column:former_limitNumber;  not null; default:0;"json:"formerlyLimitNumber"`     // 已使用额度
	Time                string  `gorm:"column:time;  type:varchar(255); default:0;"json:"createTA"`                     //时间
}

type UserMess struct {
	UUID      string `json:"uuid"`
	Name      string `json:"name"`
	HeadPhoto string `json:"headPhoto"`
	Email     string `json:"email"`
}

func (*User) TableName() string {
	return "user"
}
