package main

import (
	accounts "gotransact/apps/accounts/models"
	base "gotransact/apps/base/models"
	base_utils "gotransact/apps/base/utils"

	// transaction_handler "gotransact/apps/transaction/handlers"
	transaction "gotransact/apps/transaction/models"
	transaction_utils "gotransact/apps/transaction/utils"
	config "gotransact/config"
	logger "gotransact/log"
	database "gotransact/pkg/db"

	// "net/http"
	// "github.com/gin-gonic/gin"
	"github.com/robfig/cron"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gotransact/routers"
)

func main() {
	config.Loadenv()
	database.InitDB("")
	// validators.Init()
	// validators2.InitValidation()
	logger.Init()

	database.DB.AutoMigrate(
		&base.Base{}, &accounts.User{}, &accounts.Company{},
		&transaction.PaymentGateway{}, &transaction.TransactionRequest{},
		&transaction.TransactionHistory{})
	// Define your routes here

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
