package services

import (
	"net/http"
	"strings"

	"github.com/MPoline/graduation_project_1/internal/database"
	"github.com/MPoline/graduation_project_1/internal/models"

	"github.com/gin-gonic/gin"
)

func isValidNumber(s string) bool {
	if s == "" {
		return false
	}

	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

func luhnCheck(number string) bool {
	sum := 0
	parity := len(number) % 2

	for i, digit := range number {
		d := int(digit - '0')

		if i%2 == parity {
			d *= 2
			if d > 9 {
				d -= 9
			}
		}

		sum += d
	}

	return sum%10 == 0
}

func DownloadOrders(c *gin.Context) {
	ctx := c.Request.Context()

	token, err := c.Cookie("auth_token")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Пользователь не аутентифицирован",
		})
		return
	}

	login, err := models.ValidateAuthToken(token)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Недействительный токен авторизации",
		})
		return
	}

	orderNumber, err := c.GetRawData()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Неверный формат запроса",
		})
		return
	}

	orderStr := strings.TrimSpace(string(orderNumber))

	if !isValidNumber(orderStr) {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"message": "Номер заказа должен содержать только цифры",
		})
		return
	}

	if !luhnCheck(orderStr) {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"message": "Неверный формат номера заказа",
		})
		return
	}

	db := database.New()
	defer db.Close()

	order, err := db.SelectOrder(ctx, orderStr)
	if err != nil && err.Error() != "OrderNotFound" {
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
		err = db.AddOrder(ctx, login, orderStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}
		c.JSON(http.StatusAccepted, gin.H{
			"message": "Новый номер заказа принят в обработку",
		})
	}
}
