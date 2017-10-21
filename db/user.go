package db

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Pwd   string
	Name  string `gorm:"unique_index"`
	Email string `gorm:"unique_index"`
}

func QueryUser(u interface{}) *User {
	user := new(User)
	db.Where(u).First(&user)
	return user
}

func FindUser(u interface{}, out interface{}) {
	db.Model(&User{}).Where(u).Select([]string{"id", "name", "email"}).Limit(1).Scan(out)
}

func ChangePwd(id uint, pwd string) error {
	return db.Model(&User{}).Where("id = ?", id).Update("pwd", pwd).Error
}
