package database

import (
	"context"
	"database/sql"

	"github.com/MPoline/graduation_project_1/internal/flags"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"go.uber.org/zap"
)

type Database struct {
	dbConn *sql.DB
}

func New() *Database {
	dbConn, err := OpenDBConnection()
	if err != nil {
		zap.L().Error("Error opening database: ", zap.Error(err))
	}
	return &Database{
		dbConn: dbConn,
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

func (db *Database) Close() {
	if err := db.dbConn.Close(); err != nil {
		zap.L().Error("Error closing database: ", zap.Error(err))
	} else {
		zap.L().Info("The database connection was closed")
	}
}

func DBInit() error {
	db := New()
	defer db.Close()

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

func handlePGError(err error) {
	if pgErr, ok := err.(*pgconn.PgError); ok {
		switch pgErr.Code {
		case pgerrcode.UniqueViolation:
			zap.L().Warn("PostgreSQL unique constraint violation",
				zap.String("constraint_name", pgErr.ConstraintName),
				zap.String("detail", pgErr.Detail))
		case pgerrcode.ForeignKeyViolation:
			zap.L().Warn("PostgreSQL foreign key constraint violation",
				zap.String("constraint_name", pgErr.ConstraintName),
				zap.String("detail", pgErr.Detail))
		case pgerrcode.CheckViolation:
			zap.L().Warn("PostgreSQL check constraint violation",
				zap.String("constraint_name", pgErr.ConstraintName),
				zap.String("detail", pgErr.Detail))
		case pgerrcode.NotNullViolation:
			zap.L().Warn("PostgreSQL not-null constraint violation",
				zap.String("column_name", pgErr.ColumnName),
				zap.String("table_name", pgErr.TableName))
		default:
			zap.L().Error("Unhandled PostgreSQL error:", zap.Error(pgErr))
		}
	} else {
		zap.L().Error("Unknown error:", zap.Error(err))
	}
}
