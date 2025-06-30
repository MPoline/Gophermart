package services

import (
	"context"
	"database/sql"

	"github.com/MPoline/Gophermart/internal/database"
	"github.com/MPoline/Gophermart/internal/flags"
	"go.uber.org/zap"
)

func New() *database.Database {
	dbConn, err := OpenDBConnection()
	if err != nil {
		zap.L().Error("Error opening database: ", zap.Error(err))
	}
	return &database.Database{
		DBConn: dbConn,
	}
}

func OpenDBConnection() (*sql.DB, error) {
	db, err := sql.Open("postgres", flags.FlagDatabaseURI)
	if err != nil {
		zap.L().Error("Error opening database: ", zap.Error(err))
		return nil, err
	}
	zap.L().Info("Successful open to the database")
	return db, nil
}

func Close(db *database.Database) {
	if err := db.DBConn.Close(); err != nil {
		zap.L().Error("Error closing database: ", zap.Error(err))
	} else {
		zap.L().Info("The database connection was closed")
	}
}

func DBInit() error {
	db := New()
	defer Close(db)

	err := db.CreateUsersTable(context.Background())
	if err != nil {
		zap.L().Error("Error create users table: ", zap.Error(err))
		return err
	}

	err = db.CreateOrdersTable(context.Background())
	if err != nil {
		zap.L().Error("Error create orders table: ", zap.Error(err))
		return err
	}
	return nil
}
