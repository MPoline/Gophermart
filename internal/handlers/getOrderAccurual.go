package handlers

import (
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func GetOrderAccrual(c *gin.Context) {
	orderNumber := c.Param("number")

	if !isValidNumber(orderNumber) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Номер заказа должен содержать только цифры",
		})
		return
	}

	if !luhnCheck(orderNumber) {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"message": "Неверный формат номера заказа",
		})
		return
	}

	accrualInfo, err := getAccrualInfoFromService(orderNumber)
	if err != nil {
		if isRateLimitError(err) {
			retryAfter := 60
			c.Header("Retry-After", strconv.Itoa(retryAfter))
			c.String(http.StatusTooManyRequests, "No more than N requests per minute allowed")
			return
		}

		zap.L().Error("Ошибка при запросе информации о начислении",
			zap.String("order", orderNumber),
			zap.Error(err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Внутренняя ошибка сервера",
		})
		return
	}

	if accrualInfo == nil {
		c.Status(http.StatusNoContent)
		return
	}

	response := gin.H{
		"order":  accrualInfo.Order,
		"status": accrualInfo.Status,
	}

	if accrualInfo.Accrual > 0 {
		response["accrual"] = accrualInfo.Accrual
	}

	c.JSON(http.StatusOK, response)
}

type AccrualInfo struct {
	Order   string
	Status  string
	Accrual float64
}

func getAccrualInfoFromService(orderNumber string) (*AccrualInfo, error) {
	
	if rand.Intn(10) == 0 { // 10% chance for rate limit
		return nil, &RateLimitError{}
	}

	if rand.Intn(5) == 0 { // 20% chance for no content
		return nil, nil
	}

	statuses := []string{"REGISTERED", "INVALID", "PROCESSING", "PROCESSED"}
	status := statuses[rand.Intn(len(statuses))]

	var accrual float64
	if status == "PROCESSED" {
		accrual = float64(rand.Intn(10000)) / 100
	}

	return &AccrualInfo{
		Order:   orderNumber,
		Status:  status,
		Accrual: accrual,
	}, nil
}

type RateLimitError struct{}

func (e *RateLimitError) Error() string {
	return "rate limit exceeded"
}

func isRateLimitError(err error) bool {
	_, ok := err.(*RateLimitError)
	return ok
}
