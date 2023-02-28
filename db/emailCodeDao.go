package db

import "array/model"

type emailCodeDao struct{}

var EmailCodeDao *emailCodeDao

func init() {
	EmailCodeDao = &emailCodeDao{}
}

func (*emailCodeDao) CreateUser(email *model.EmailCode) {
	mainDB.Create(&email)
}

func (*emailCodeDao) DeleteCode(email *model.EmailCode) {
	mainDB.Delete(&email, email.Id)
}

func (*emailCodeDao) GetCode(email string) (res *model.EmailCode, err error) {
	err = mainDB.Find(&res, "email=?", email).Error
	return
}
