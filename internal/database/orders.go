package database

import (
	"context"
	"database/sql"
	"errors"

	"github.com/MPoline/Gophermart/internal/models"
	"go.uber.org/zap"
)

func (db *Database) CreateOrdersTable(ctx context.Context) error {
	createQuery := ` CREATE TABLE IF NOT EXISTS orders ( 
		number BIGINT PRIMARY KEY, 
		login TEXT UNIQUE NOT NULL 
	); `

	_, err := db.dbConn.ExecContext(ctx, createQuery)
	if err != nil {
		handlePGError(err)
		return err
	}
	zap.L().Info("Таблица orders создана")
	return nil
}

func (db *Database) AddOrder(ctx context.Context, login string, orderNumber string) error {
	insertQuery := `INSERT INTO orders (login, number) VALUES ($1, $2)`
	_, err := db.dbConn.ExecContext(ctx, insertQuery, login, orderNumber)
	if err != nil {
		handlePGError(err)
		return err
	}

	zap.L().Info("Номер заказа успешно добавлен", zap.String("login, number", login+orderNumber))
	return nil
}

func (db *Database) SelectOrder(ctx context.Context, orderNumber string) (models.Order, error) {
	var order models.Order
	row := db.dbConn.QueryRowContext(ctx, `SELECT * FROM orders WHERE number = $1`, orderNumber)
	err := row.Scan(&order.OrderNumber, &order.Login)
	if err != nil {
		if err == sql.ErrNoRows {
			err = errors.New("OrderNotFound")
			zap.L().Info("Order not found")
			return order, err
		}
		handlePGError(err)
		return order, err
	}
	return order, nil
}
