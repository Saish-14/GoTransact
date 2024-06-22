package routers

import (
	account "gotransact/apps/accounts/handlers"
	transaction "gotransact/apps/transaction/handlers"
	docs "gotransact/docs"
	"gotransact/middleware"

	"github.com/gin-gonic/gin"
)

func Routers() *gin.Engine {

	r := gin.Default()

	docs.SwaggerInfo.BasePath = "/api"

	api := r.Group("/api")
	{
		api.POST("/register", account.SignupHandler)
		api.POST("/login", account.LoginHandler)
		api.GET("/confirm-payment", transaction.ConfirmPaymentHandler)

		protected := api.Group("/protected")
		protected.Use(middleware.AuthMiddleware())
		{
			protected.POST("/post-payment", transaction.PaymentRequest)
			protected.POST("/logout", account.LogoutHandler)
		}
	}
	return r
}