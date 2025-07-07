package dbservice

import "database/sql"

type DBService struct {
	DB          *sql.DB
	Initialized bool
} 
