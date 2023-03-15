package db

import (
	"encoding/json"
	"fmt"
	"github.com/ZakirAvrora/TechHome/config"
	"github.com/ZakirAvrora/TechHome/internals/models"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"io"
	"log"
	"os"
)

const migrationPath = "file:./migrations"

func DatabaseInit(conf *config.DatabaseConfig) (*sqlx.DB, error) {
	dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", conf.User, conf.Password, conf.Host, conf.Port, conf.DbName)
	log.Println(dbUrl)

	db, err := sqlx.Open("postgres", dbUrl)
	if err != nil {
		return nil, fmt.Errorf("connection to database failed:%w", err)
	}

	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return nil, fmt.Errorf("the migration instance error: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		migrationPath, "postgres", driver)
	if err != nil {
		return nil, fmt.Errorf("migration error: %w", err)
	}

	if err := m.Up(); err != nil {
		return nil, fmt.Errorf("migration up error: %w", err)
	}

	return db, nil
}

func JsonUploadToDb(db *sqlx.DB, filePath string) (err error) {
	var links []models.Link

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error in opening file: %w", err)
	}

	defer func() {
		if cErr := file.Close(); cErr != nil {
			err = cErr
		}
	}()

	content, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("error in reading the file %s: %w", filePath, err)
	}

	if err := json.Unmarshal(content, &links); err != nil {
		return fmt.Errorf("error in marshalling json file: %w", err)
	}

	tx, err := db.Beginx()
	if err != nil {
		fmt.Errorf("error in starting transaction: %w", err)
	}

	_, err = tx.NamedExec(`INSERT INTO redirects (active_link, history_link)
        VALUES (:active_link, :history_link)`, links)
	if err != nil {
		if rErr := tx.Rollback(); rErr != nil {
			return fmt.Errorf("error in rollbacking transaction: %w", err)
		}
		return fmt.Errorf("error in transaction: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error in transaction commit: %w", err)
	}

	return nil
}
