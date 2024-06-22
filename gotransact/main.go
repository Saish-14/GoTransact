package main

import (
	accounts "gotransact/apps/accounts/models"
	account_validators "gotransact/apps/accounts/validators"
	base "gotransact/apps/base/models"
	base_utils "gotransact/apps/base/utils"
	transaction "gotransact/apps/transaction/models"
	transaction_utils "gotransact/apps/transaction/utils"
	transaction_validators "gotransact/apps/transaction/validators"
	config "gotransact/config"
	logger "gotransact/log"
	database "gotransact/pkg/db"
	routers "gotransact/routers"
	cron "github.com/robfig/cron"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	config.Loadenv()
	database.InitDB("")
	account_validators.Init()
	transaction_validators.Init()
	logger.Init()

	database.DB.AutoMigrate(
		&base.Base{}, &accounts.User{}, &accounts.Company{},
		&transaction.PaymentGateway{}, &transaction.TransactionRequest{},
		&transaction.TransactionHistory{})

    var c = cron.New()
    c.AddFunc("@every 1h", func() {
        transactions := transaction_utils.FetchTransactionsLast24Hours()
        filePath, err := transaction_utils.GenerateExcel(transactions)
        if err != nil {
            logger.ErrorLogger.Fatalf("failed to generate excel: %v", err)
        }
        base_utils.SendMailWithAttachment(
			"nsaish@trellissoft.ai",
			"Daily Transactions Report",
			"Please find attached the daily transactions report.",
			filePath)
    })
    c.Start()
	r := routers.Routers()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Run(":8080")
}
