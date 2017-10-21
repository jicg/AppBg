package db

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	//_ "github.com/jinzhu/gorm/dialects/sqlite"
	"time"

	"github.com/jicg/AppBg/conf"
)

var (
	db *gorm.DB
)

func init() {
	var err error
	db, err = gorm.Open("mysql", conf.GetConf().Db)
	//db, err = gorm.Open("sqlite3", "data/data.db")
	db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&User{}, &Offorder{}, &Cache{})
	//if err != nil {
	//	fmt.Println(err.Error())
	//	return
	//}
	//db.AutoMigrate(&User{}, &Offorder{})
	if err != nil {
		fmt.Println(err.Error())
	}
	//db.LogMode(true)
	db.Create(&User{Name: "admin", Pwd: "admin", Email: "jicg@qq.com"})
}

type Model struct {
	ID        uint       `json:"id" gorm:"primary_key"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-" sql:"index"`
}
