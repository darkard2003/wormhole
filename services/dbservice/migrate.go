package dbservice

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

type MigrationTable struct {
	ID          int    `json:"id"`
	Migration   string `json:"migration"`
	Description string `json:"description"`
	AppliedAt   string `json:"applied_at"`
}

func (db *DBService) InitializeMigrationTable() error {
	if db.DB == nil {
		return fmt.Errorf("database connection is not initialized")
	}
	query := `CREATE TABLE IF NOT EXISTS migrations (
		id INT AUTO_INCREMENT PRIMARY KEY,
		migration VARCHAR(255) NOT NULL,
		description TEXT,
		applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`
	_, err := db.DB.Exec(query)
	if err != nil {
		log.Println("Error creating migrations table:", err)
		return fmt.Errorf("error creating migrations table: %w", err)
	}
	log.Println("Migrations table initialized successfully")
	return nil
}

func (db *DBService) GetAppliedMigrations() ([]MigrationTable, error) {
	if db.DB == nil {
		return nil, fmt.Errorf("database connection is not initialized")
	}
	query := `SELECT id, migration, description, applied_at FROM migrations;`
	rows, err := db.DB.Query(query)
	if err != nil {
		log.Println("Error fetching applied migrations:", err)
		return nil, fmt.Errorf("error fetching applied migrations: %w", err)
	}
	defer rows.Close()

	var migrations []MigrationTable
	for rows.Next() {
		var migration MigrationTable
		if err := rows.Scan(&migration.ID, &migration.Migration, &migration.Description, &migration.AppliedAt); err != nil {
			log.Println("Error scanning row:", err)
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		migrations = append(migrations, migration)
	}
	return migrations, nil
}

func (db *DBService) MigrationExists(migrationFile string) (bool, error) {
	if db.DB == nil {
		return false, fmt.Errorf("database connection is not initialized")
	}
	query := `SELECT COUNT(*) FROM migrations WHERE migration = ?;`
	var count int
	err := db.DB.QueryRow(query, migrationFile).Scan(&count)
	if err != nil {
		log.Println("Error checking migration existence:", err)
		return false, fmt.Errorf("error checking migration existence: %w", err)
	}
	return count > 0, nil
}

func GetMigrationFiles() ([]string, error) {
	os.MkdirAll("migrations", os.ModePerm)
	files, err := os.ReadDir("migrations")
	if err != nil {
		log.Println("Error reading migration files:", err)
		return nil, fmt.Errorf("error reading migration files: %w", err)
	}
	var migrationFiles []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".sql") {
			migrationFiles = append(migrationFiles, file.Name())
		}
	}
	if len(migrationFiles) == 0 {
		log.Println("No migration files found")
		return nil, fmt.Errorf("no migration files found")
	}
	log.Println("Migration files found:", migrationFiles)
	sort.Strings(migrationFiles)
	log.Println("Sorted migration files:", migrationFiles)
	return migrationFiles, nil
}

func (db *DBService) Migrate() error {
	if db.DB == nil {
		return fmt.Errorf("database connection is not initialized")
	}
	er := db.InitializeMigrationTable()
	if er != nil {
		return fmt.Errorf("error initializing migration table: %w", er)
	}
	migrationFiles, err := GetMigrationFiles()
	if err != nil {
		return fmt.Errorf("error getting migration files: %w", err)
	}
	for _, migrationFile := range migrationFiles {
		exists, err := db.MigrationExists(migrationFile)
		if err != nil {
			return fmt.Errorf("error checking migration existence: %w", err)
		}
		if exists {
			fmt.Printf("Migration %s already applied, skipping\n", migrationFile)
			continue
		}

		tx, err := db.DB.Begin()
		if err != nil {
			return fmt.Errorf("error starting transaction: %w", err)
		}
		content, err := os.ReadFile("migrations/" + migrationFile)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("error reading migration file %s: %w", migrationFile, err)
		}
		query := string(content)
		if query == "" {
			tx.Rollback()
			return fmt.Errorf("migration file %s is empty", migrationFile)
		}
		_, err = tx.Exec(query)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("error applying migration %s: %w", migrationFile, err)
		}
		if _, err := tx.Exec("INSERT INTO migrations (migration, description, applied_at) VALUES (?, ?, NOW())", migrationFile, "Migration applied successfully"); err != nil {
			tx.Rollback()
			return fmt.Errorf("error recording migration %s: %w", migrationFile, err)
		}
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("error committing transaction for migration %s: %w", migrationFile, err)
		}
	}
	log.Println("All migrations applied successfully")
	return nil
}
