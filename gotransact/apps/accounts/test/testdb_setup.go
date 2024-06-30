package tests

import (
	"fmt"
	"gotransact/apps/accounts/models"
	"gotransact/pkg/db"
)

func SetupTestDb() {
	db.InitDB("test")
	if err := db.DB.AutoMigrate(&models.User{}, &models.Company{}); err != nil {
		fmt.Printf("Error autoigrating models : %s", err.Error())
	}
}

func CloseTestDb() {
	sqlDB, err := db.DB.DB()
	if err != nil {
		fmt.Printf("Error getting sqlDB from gorm DB: %s", err.Error())
		return
	}
	if err := sqlDB.Close(); err != nil {
		fmt.Printf("Error closing database: %s", err.Error())
	}
}

func ClearDatabase() {
	db.DB.Exec("DELETE FROM companies")
	db.DB.Exec("DELETE FROM users")
}
