package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
		"github.com/MPoline/graduation_project_1/internal/models"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("auth_token")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Требуется авторизация",
			})
			return
		}

		login, err := models.ValidateAuthToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Недействительный токен",
			})
			return
		}

		c.Set("userLogin", login)
		c.Next()
	}
}
