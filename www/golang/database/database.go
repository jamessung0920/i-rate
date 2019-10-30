package database

import (
	"fmt"
	"os"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var DB *gorm.DB

func ConnectionDB() (*gorm.DB, error) {

	database := os.Getenv("DB_NAME")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")

	// name := "dbname=" + database + " user=" + user + " password=" + password + " host=postgres sslmode=disable"
	name := user + ":" + password + "@(db)/" + database + "?charset=utf8&parseTime=True&loc=Local"

	var db *gorm.DB
	var err error
	for {
		db, err = gorm.Open("mysql", name)
		if err != nil {
			fmt.Println(err)
		} else {
			err = db.DB().Ping()
			if err == nil {
				break
			}
		}
		time.Sleep(time.Second * 3)
	}

	//defer db.Close()
	db.DB().SetMaxIdleConns(10)
	db.LogMode(true)
	DB = db
	return db, nil
}

func GetDB() *gorm.DB {
	return DB
}
