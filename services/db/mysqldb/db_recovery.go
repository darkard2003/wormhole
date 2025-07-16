package mysqldb

import "database/sql"

func RecoverDB(tx *sql.Tx, err *error) {
	if r := recover(); r != nil {
		tx.Rollback()
		panic(r)
	} else if err != nil {
		tx.Rollback()
	} else {
		tx.Commit()
	}
}
