package handlers

import (
	"errors"
	"net/http"

	"github.com/MPoline/Gophermart/internal/models"
	"github.com/MPoline/Gophermart/internal/services"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func isValidPassword(password string) (bool, error) {
	if len(password) < 8 {
		return false, errors.New("Минимальная длина пароля 8 символов")
	}
	return true, nil
}

func hashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func RegisterUser(c *gin.Context) {
	ctx := c.Request.Context()

	var input struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	err := c.ShouldBindJSON(&input)
	if err != nil {
		zap.L().Info("Неверный формат JSON RegisterUser",
			zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"message": "Неверный формат JSON RegisterUser"})
		return
	}

	db := services.New()
	defer services.Close(db)

	isLoginExist, err := db.CheckLogin(ctx, input.Login)
	if err != nil {
		zap.L().Info("Ошибка сервера при обработке запроса (isLoginExist)",
			zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка сервера при обработке запроса (isLoginExist)"})
		return
	}
	if isLoginExist {
		zap.L().Info("Пользователь с таким логином уже существует",
			zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"message": "Пользователь с таким логином уже существует"})
		return
	}

	if valid, err := isValidPassword(input.Password); !valid {
		zap.L().Info("Пароль не соответствует требованиям",
			zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	hashedPassword, err := hashPassword(input.Password)
	if err != nil {
		zap.L().Info("Пароль не верный",
			zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка сервера при обработке запроса (hashedPassword)"})
		return
	}

	err = db.AddUser(ctx, input.Login, string(hashedPassword))
	if err != nil {
		zap.L().Info("Ошибка сервера при регистрации пользователя",
			zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка сервера при регистрации пользователя"})
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

	c.JSON(http.StatusOK, gin.H{"message": "Пользователь успешно зарегистрирован и авторизован"})

}
