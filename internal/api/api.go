package api

import (
	"github.com/MPoline/Gophermart/internal/middleware"
	"github.com/MPoline/Gophermart/internal/services"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	public := router.Group("/api")
	{
		public.POST("/user/register", services.RegisterUser)
		public.POST("/user/login", services.LoginUser)
	}

	private := router.Group("/api")
	private.Use(middleware.AuthMiddleware())
	{
		private.POST("/user/orders", services.DownloadOrders)
		// router.GET("/api/user/orders", services.GetOrders)
		// router.GET("/api/user/balance", services.GetBalance)
		// router.POST("/api/user/balance/withdraw", services.WithdrawBalance)
		// router.GET("/api/user/withdrawals", services.GetWithdrawals)
	}

	return router
}
