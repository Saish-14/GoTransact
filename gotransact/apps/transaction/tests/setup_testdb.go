package test

import (
	accounts_model "gotransact/apps/accounts/models"
	transaction_model "gotransact/apps/transaction/models"
	"gotransact/pkg/db"
	"fmt"
)

func SetupTestDb() {
	db.InitDB("test")
	if err := db.DB.AutoMigrate(
		&accounts_model.User{}, &accounts_model.Company{},
		&transaction_model.PaymentGateway{},
		&transaction_model.TransactionRequest{},
		&transaction_model.TransactionHistory{},
		); err != nil {
		fmt.Printf("Error automigrating models : %s", err.Error())
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
	db.DB.Exec("DELETE FROM transaction_histories")
	db.DB.Exec("DELETE FROM transaction_requests")
	db.DB.Exec("DELETE FROM payment_gateways")
}
