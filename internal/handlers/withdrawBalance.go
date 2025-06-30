package handlers

import (
	"net/http"

	"github.com/MPoline/Gophermart/internal/services"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func WithdrawBalance(c *gin.Context) {
	ctx := c.Request.Context()
	login := c.GetString("userLogin")

	if login == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"mesage": "Требуется авторизация",
		})
		return
	}

	var request struct {
		Order string  `json:"order"`
		Sum   float64 `json:"sum"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		zap.L().Error("Ошибка парсинга запроса",
			zap.String("login", login),
			zap.Error(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Неверный формат запроса",
		})
		return
	}

	if !isValidNumber(request.Order) {
		zap.L().Info("Номер заказа должен содержать только цифры")
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"message": "Номер заказа должен содержать только цифры",
		})
		return
	}

	if !luhnCheck(request.Order) {
		zap.L().Info("Неверный формат номера заказа")
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"message": "Неверный формат номера заказа",
		})
		return
	}

	db := services.New()
	defer services.Close(db)

	order, err := db.SelectOrder(ctx, request.Order)
	if err != nil && err.Error() != "OrderNotFound" {
		zap.L().Info("Ошибка при получении заказов",
			zap.Error(err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	switch {
	case order.Login == login:
		c.JSON(http.StatusOK, gin.H{
			"message": "Этот номер заказа уже был загружен этим пользователем",
		})
	case order.Login != "":
		c.AbortWithStatusJSON(http.StatusConflict, gin.H{
			"message": "Этот номер заказа уже был загружен другим пользователем",
		})
	default:
		err = db.AddOrder(ctx, login, request.Order)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}
	}

	err = db.WithdrawBalance(ctx, login, request.Order, request.Sum)
	if err != nil {
		zap.L().Error("Ошибка списания средств",
			zap.String("login", login),
			zap.String("order", request.Order),
			zap.Float64("sum", request.Sum),
			zap.Error(err))
		if err.Error() == "Недостаточно баллов для списания" {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "Недостаточно баллов для списания",
			})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Внутренняя ошибка сервера",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Баллы успешно списаны",
	})
}
