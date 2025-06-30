package database

import (
	"database/sql"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"go.uber.org/zap"
)

type Database struct {
	DBConn *sql.DB
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
