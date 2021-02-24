package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Article struct {
	Id                     uint16
	Title, Anons, FullText string
}

func initDB() {
	db, _ := gorm.Open("sqlite3", "./gorm.db")
	defer db.Close()

	db.AutoMigrate(&Article{})

	db.Create(&p2)
}

func createPost() {

}
