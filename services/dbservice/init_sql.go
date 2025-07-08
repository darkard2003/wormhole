package dbservice

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
)

func (db *DBService) InitializeMySql() error {
	if db.Initialized {
		return nil
	}

	cfg := mysql.NewConfig()
	cfg.User = os.Getenv("DB_USER")
	cfg.Net = "tcp"
	cfg.Addr = os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT")
	cfg.DBName = os.Getenv("DB_NAME")
	cfg.Passwd = os.Getenv("DB_PASSWORD")

	var err error
	db.DB, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Println("Error connecting to the database:", err)
		return fmt.Errorf("error connecting to the database: %w", err)
	}

	pingErr := db.DB.Ping()
	if pingErr != nil {
		log.Println("Error pinging the database:", pingErr)
		return fmt.Errorf("error pinging the database: %w", pingErr)
	}
	log.Println("Database connection established successfully")
	err = db.Migrate()
	if err != nil {
		log.Println("Error during migration:", err)
		return fmt.Errorf("error during migration: %w", err)
	}
	log.Println("Database migration completed successfully")
	db.Initialized = true
	return nil
}
