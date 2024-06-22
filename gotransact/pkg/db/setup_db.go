package db

import (
	"fmt"
	config "gotransact/config"
	"strings"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(typ string) {
	var dbURI string
	if strings.EqualFold(typ, "test") {
		config.Loadenv()

		dbURI = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s TimeZone=%s", config.DbHost, config.DbUser, config.DbPassword, "mydatabase", config.DbPort, config.DbTimezone)
	} else {
		dbURI = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s TimeZone=%s", config.DbHost, config.DbUser, config.DbPassword, config.DbName, config.DbPort, config.DbTimezone)
	}
	db, err := gorm.Open(postgres.Open(dbURI), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	DB = db
}
