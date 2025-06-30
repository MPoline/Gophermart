package handlers

import (
	"net/http"
	"time"

	"github.com/MPoline/Gophermart/internal/services"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func GetOrders(c *gin.Context) {
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

	orders, err := db.GetUserOrders(ctx, login)
	if err != nil {
		zap.L().Info("Ошибка при получении заказов",
			zap.Error(err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Ошибка при получении заказов",
		})
		return
	}

	if len(orders) == 0 {
		c.AbortWithStatus(http.StatusNoContent)
		return
	}

	var response []gin.H
	for _, order := range orders {
		orderData := gin.H{
			"number":      order.Number,
			"status":      order.Status,
			"uploaded_at": order.UploadedAt.Format(time.RFC3339),
		}

		if order.Status == "PROCESSED" && order.Accrual != nil {
			orderData["accrual"] = *order.Accrual
		}

		response = append(response, orderData)
	}

	c.JSON(http.StatusOK, response)

}
