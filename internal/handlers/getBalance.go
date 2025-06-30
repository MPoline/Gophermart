package handlers

import (
	"net/http"

	"github.com/MPoline/Gophermart/internal/services"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func GetBalance(c *gin.Context) {
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

	balance, err := db.GetUserBalance(ctx, login)
	if err != nil {
		zap.L().Error("Ошибка получения баланса",
			zap.String("login", login),
			zap.Error(err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Внутренняя ошибка сервера",
		})
		return
	}

	response := gin.H{
		"current":   balance.Current,
		"withdrawn": balance.Withdrawn,
	}

	c.JSON(http.StatusOK, response)
}
