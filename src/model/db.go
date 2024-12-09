package model

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB // Package-level variable for database connection

func init() {
	var err error
	db, err = gorm.Open(sqlite.Open("db/test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// AutoMigrate database models
	err = db.AutoMigrate(&User{}, &Space{}, &Member{}, &Good{}, &ResponsibleUid{})
	if err != nil {
		panic("failed to migrate database")
	}
}
