package database

import (
	"context"
	"database/sql"
	"errors"

	"github.com/MPoline/Gophermart/internal/models"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func (db *Database) CreateUsersTable(ctx context.Context) error {
	createQuery := ` CREATE TABLE IF NOT EXISTS users ( 
		id SERIAL PRIMARY KEY, 
		login TEXT UNIQUE NOT NULL, 
		hashedPassword BYTEA NOT NULL
	); `

	_, err := db.DBConn.ExecContext(ctx, createQuery)
	if err != nil {
		handlePGError(err)
		return err
	}
	zap.L().Info("Таблица users создана")
	return nil
}

func (db *Database) CheckLogin(ctx context.Context, login string) (bool, error) {
	var found bool
	err := db.DBConn.QueryRow(`SELECT EXISTS(SELECT 1 FROM users WHERE login = $1)`, login).Scan(&found)
	if err != nil {
		handlePGError(err)
		return false, err
	}
	if found {
		return true, nil
	}
	return false, nil
}

func (db *Database) AddUser(ctx context.Context, login string, hashedPassword string) error {
	insertQuery := `INSERT INTO users (login, hashedPassword) VALUES ($1, $2)`
	_, err := db.DBConn.ExecContext(ctx, insertQuery, login, hashedPassword)
	if err != nil {
		handlePGError(err)
		return err
	}

	zap.L().Info("Пользователь успешно добавлен", zap.String("login", login))
	return nil
}

func (db *Database) FindUser(ctx context.Context, login string) (models.User, error) {
	var user models.User
	row := db.DBConn.QueryRowContext(ctx, ` SELECT id, login, hashedPassword FROM users WHERE login = $1 LIMIT 1 `, login)
	err := row.Scan(&user.ID, &user.Login, &user.HashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			err = errors.New("UserNotFound")
			zap.L().Info("User not found")
			return user, err
		}
		handlePGError(err)
		return user, err
	}
	return user, nil
}
