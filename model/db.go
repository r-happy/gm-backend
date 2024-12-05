package model

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func init() {
	var err error
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	// db, err = gorm.Open("sqlite3", "db/sample.db")
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Space{})
	db.AutoMigrate(&Member{})
	db.AutoMigrate(&Good{})
	db.AutoMigrate(&ResponsibleUid{})
}
