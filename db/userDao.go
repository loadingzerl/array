package db

import "array/model"

type userDao struct{}

var UserDao *userDao

func init() {
	UserDao = &userDao{}
}

func (*userDao) GetUser(email string) (res *model.User, err error) {
	err = mainDB.Find(&res, "email=?", email).Error
	return
}

func (*userDao) GetUserID(id int64) (res *model.User, err error) {
	err = mainDB.Find(&res, "id=?", id).Error
	return
}

// CreateUser 录入用户
func (*userDao) CreateUser(user *model.User) {
	mainDB.Create(&user)
}

// UpdatePassWord  更新密码
func (*userDao) UpdatePassWord(user *model.User, password string) int64 {
	t := mainDB.Model(&model.User{}).Where("email = ?", user.Email).Update("pass_word", password)
	return t.RowsAffected
}

func (*userDao) UpdateUser(user *model.User) int64 {
	rews := mainDB.Model(&model.User{}).Where("id", user.Id).Updates(user)
	return rews.RowsAffected
}

// GetAll 返回所有用户信息
func (*userDao) GetAll() (res []*model.User, err error) {
	err = MainDB().Find(&res).Error
	return
}
