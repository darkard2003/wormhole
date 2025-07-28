package mysqldb

import (
	"errors"

	errorNumbers "github.com/bombsimon/mysql-error-numbers"
	"github.com/darkard2003/wormhole/internals/services/db"
	"github.com/go-sql-driver/mysql"
)

func ToDBError(err error, table, item string) error {
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		switch mysqlErr.Number {
		case errorNumbers.ER_DUP_ENTRY:
			return db.NewAlreadyExistsError(table, item, mysqlErr.Message)

		case errorNumbers.ER_NO_SUCH_TABLE:
			return db.NewNotFoundError(table, item)
		case errorNumbers.ER_BAD_FIELD_ERROR:
			return db.NewValidationError(item, mysqlErr.Message)
		default:
			return db.NewDBError(table, item, mysqlErr.Message)
		}
	}

	return db.NewInternalError(err.Error())
}
