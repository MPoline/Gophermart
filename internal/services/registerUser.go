package services

import (
	"errors"
	"net/http"

	"github.com/MPoline/Gophermart/internal/database"
	"github.com/MPoline/Gophermart/internal/models"
	"github.com/gin-gonic/gin"
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
		c.JSON(http.StatusBadRequest, gin.H{"message": "Неверный формат JSON RegisterUser"})
		return
	}

	db := database.New()
	defer db.Close()

	isLoginExist, err := db.CheckLogin(ctx, input.Login)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка сервера при обработке запроса (isLoginExist)"})
		return
	}
	if isLoginExist {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Пользователь с таким логином уже существует"})
		return
	}

	if valid, err := isValidPassword(input.Password); !valid {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	hashedPassword, err := hashPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка сервера при обработке запроса (hashedPassword)"})
		return
	}

	err = db.AddUser(ctx, input.Login, string(hashedPassword))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка сервера при регистрации пользователя"})
		return
	}

	token, err := models.GenerateAuthToken(input.Login)
	if err != nil {
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
