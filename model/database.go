package model

import (
	"carbonfootprint/config"
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var Db *gorm.DB

func SetDB() {
	config := config.GetConfig()
	dbUser := config.Database.User
	dbPassword := config.Database.Password
	dbHost := config.Database.Host
	dbPort := config.Database.DbPort
	dbName := config.Database.Name
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPassword, dbHost, dbPort, dbName)
	var err error
	Db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		log.Fatal("Veritabanına bağlanırken hata oluştu:", err)
	}
}

func GetDB() *gorm.DB {
	return Db
}
