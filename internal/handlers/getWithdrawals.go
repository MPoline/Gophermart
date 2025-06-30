package handlers

import (
	"net/http"

	"github.com/MPoline/Gophermart/internal/services"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func GetWithdrawals(c *gin.Context) {
	ctx := c.Request.Context()
	login := c.GetString("userLogin")

	if login == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"mesage": "Требуется авторизация",
		})
		return
	}

	db := services.New()
	defer services.Close(db)

	withdrawals, err := db.GetWithdrawals(ctx, login)
	if err != nil {
		zap.L().Error("Ошибка при получении информации о выводе средств",
			zap.String("login", login),
			zap.Error(err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Внутренняя ошибка сервера",
		})
		return
	}

	if len(withdrawals) == 0 {
		c.Status(http.StatusNoContent)
		return
	}

	c.JSON(http.StatusOK, withdrawals)

}
