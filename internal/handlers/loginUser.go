package handlers

import (
	"net/http"

	"github.com/MPoline/Gophermart/internal/models"
	"github.com/MPoline/Gophermart/internal/services"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func LoginUser(c *gin.Context) {
	ctx := c.Request.Context()

	var input struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	err := c.ShouldBindJSON(&input)
	if err != nil {
		zap.L().Info("Неверный формат JSON LoginUser",
			zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"message": "Неверный формат JSON LoginUser"})
		return
	}

	db := services.New()
	defer services.Close(db)

	user, err := db.FindUser(ctx, input.Login)
	if err != nil {
		zap.L().Info("Ошибка при проверке пользователя",
			zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	err = bcrypt.CompareHashAndPassword(user.HashedPassword, []byte(input.Password))
	if err != nil {
		zap.L().Info("Пароль неверный",
			zap.Error(err))
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Пароль неверный"})
		return
	}

	token, err := models.GenerateAuthToken(input.Login)
	if err != nil {
		zap.L().Info("Ошибка при создании токена авторизации",
			zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка при создании токена авторизации"})
		return
	}

	c.SetCookie(
		"auth_token",
		token,
		int(models.TokenTTL.Seconds()),
		"/",
		"",
		false, // secure (false для HTTP, true для HTTPS)
		true,  // httpOnly
	)

	c.JSON(http.StatusOK, gin.H{"message": "Успешная авторизация"})
}
