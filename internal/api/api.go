package api

import (
	"github.com/MPoline/Gophermart/internal/handlers"
	"github.com/MPoline/Gophermart/internal/middleware"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	public := router.Group("/api")
	{
		public.POST("/user/register", handlers.RegisterUser)
		public.POST("/user/login", handlers.LoginUser)
	}

	private := router.Group("/api")
	private.Use(middleware.AuthMiddleware())
	{
		private.POST("/user/orders", handlers.DownloadOrders)
		private.GET("/user/orders", handlers.GetOrders)
		private.GET("/user/balance", handlers.GetBalance)
		private.POST("/user/balance/withdraw", handlers.WithdrawBalance)
		private.GET("/user/withdrawals", handlers.GetWithdrawals)
		private.GET("/orders/:number", handlers.GetOrderAccrual)
	}

	return router
}
