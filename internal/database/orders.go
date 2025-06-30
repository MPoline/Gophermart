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
		number BIGINT UNIQUE PRIMARY KEY, 
		login TEXT NOT NULL,
		status TEXT NOT NULL DEFAULT 'NEW' 
            CHECK (status IN ('NEW', 'PROCESSING', 'INVALID', 'PROCESSED', 'WITHDRAWN')),
		accrual FLOAT,
		amount FLOAT,
		uploaded_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	); `

	_, err := db.DBConn.ExecContext(ctx, createQuery)
	if err != nil {
		handlePGError(err)
		return err
	}
	zap.L().Info("Таблица orders создана")
	return nil
}

func (db *Database) AddOrder(ctx context.Context, login string, orderNumber string) error {
	insertQuery := `INSERT INTO orders (login, number) VALUES ($1, $2)`
	_, err := db.DBConn.ExecContext(ctx, insertQuery, login, orderNumber)
	if err != nil {
		handlePGError(err)
		return err
	}

	zap.L().Info("Номер заказа успешно добавлен", zap.String("login, number", login+orderNumber))
	return nil
}

func (db *Database) SelectOrder(ctx context.Context, orderNumber string) (models.Order, error) {
	var order models.Order
	row := db.DBConn.QueryRowContext(ctx, `
        SELECT number, login, status, accrual, amount, uploaded_at 
        FROM orders 
        WHERE number = $1`,
		orderNumber)

	err := row.Scan(
		&order.Number,
		&order.Login,
		&order.Status,
		&order.Accrual,
		&order.Amount,
		&order.UploadedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return order, errors.New("OrderNotFound")
		}
		handlePGError(err)
		return order, err
	}
	return order, nil
}

func (db *Database) GetUserOrders(ctx context.Context, login string) ([]models.Order, error) {
	query := `SELECT number, status, accrual, uploaded_at 
              FROM orders 
              WHERE login = $1 
              ORDER BY uploaded_at DESC`

	rows, err := db.DBConn.QueryContext(ctx, query, login)
	if err != nil {
		handlePGError(err)
		return nil, err
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var order models.Order
		var accrual sql.NullFloat64
		err := rows.Scan(
			&order.Number,
			&order.Status,
			&accrual,
			&order.UploadedAt,
		)
		if err != nil {
			handlePGError(err)
			return nil, err
		}

		if accrual.Valid {
			order.Accrual = &accrual.Float64
		}

		orders = append(orders, order)
	}

	if err = rows.Err(); err != nil {
		handlePGError(err)
		return nil, err
	}

	return orders, nil
}

func (db *Database) GetUserBalance(ctx context.Context, login string) (*models.Balance, error) {
	query := `
        SELECT 
            COALESCE(SUM(CASE 
                WHEN status = 'PROCESSED' THEN accrual 
                WHEN status = 'WITHDRAWN' THEN -amount 
                ELSE 0 
            END), 0) AS current,
            
            COALESCE(SUM(CASE 
                WHEN status = 'WITHDRAWN' THEN amount 
                ELSE 0 
            END), 0) AS withdrawn
        FROM orders
        WHERE login = $1
    `

	var balance models.Balance
	err := db.DBConn.QueryRowContext(ctx, query, login).Scan(
		&balance.Current,
		&balance.Withdrawn,
	)

	if err != nil {
		handlePGError(err)
		return nil, err
	}

	return &balance, nil
}

func (db *Database) WithdrawBalance(ctx context.Context, login, order string, sum float64) error {
	tx, err := db.DBConn.BeginTx(ctx, nil)
	if err != nil {
		handlePGError(err)
		return err
	}
	defer tx.Rollback()

	balance, err := db.GetUserBalance(ctx, login)
	if err != nil {
		handlePGError(err)
		return err
	}

	if balance.Current < sum {
		return errors.New("Недостаточно баллов для списания")
	}

	_, err = tx.ExecContext(ctx, `
        UPDATE orders 
        SET 
            status = 'WITHDRAWN',
            accrual = 0,
            amount = $1,
            uploaded_at = NOW()
        WHERE number = $2 AND login = $3
    `, sum, order, login)

	if err != nil {
		handlePGError(err)
		return err
	}

	if err = tx.Commit(); err != nil {
		handlePGError(err)
		return err
	}

	return nil
}

func (db *Database) GetWithdrawals(ctx context.Context, login string) ([]models.Withdrawal, error) {
	rows, err := db.DBConn.QueryContext(ctx, `
        SELECT 
            number AS order,
            amount AS sum,
            uploaded_at AS processed_at
        FROM orders
        WHERE login = $1 AND status = 'WITHDRAWN'
        ORDER BY uploaded_at DESC
    `, login)

	if err != nil {
		handlePGError(err)
		return nil, err
	}
	defer rows.Close()

	var withdrawals []models.Withdrawal
	for rows.Next() {
		var w models.Withdrawal
		err := rows.Scan(&w.Order, &w.Sum, &w.ProcessedAt)
		if err != nil {
			handlePGError(err)
			return nil, err
		}
		withdrawals = append(withdrawals, w)
	}

	if err = rows.Err(); err != nil {
		handlePGError(err)
		return nil, err
	}

	return withdrawals, nil
}
