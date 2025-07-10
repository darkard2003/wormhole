package mysqldb

import "database/sql"

type MySqlRepo struct {
	DB          *sql.DB
	Initialized bool
}
