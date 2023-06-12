package database

import (
	"database/sql"
	"fmt"

	"github.com/bhanupbalusu/gocomboums_v4/pkg/config"
	"github.com/bhanupbalusu/gocomboums_v4/pkg/logs"
	_ "github.com/lib/pq" // Here we import the pq package for its side-effects, initializing its driver.
)

func NewDB(cfg *config.Config) (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	logs.Info("Successfully connected to database")

	return db, nil
}
