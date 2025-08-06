package services

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"tracker-core/internal/models"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

type DatabaseService struct {
	db *bun.DB
}

func NewDatabaseService() (*DatabaseService, error) {
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "password")
	dbname := getEnv("DB_NAME", "crypto_tracker")

	adminDSN := fmt.Sprintf("postgres://%s:%s@%s:%s/postgres?sslmode=disable",
		user, password, host, port)

	adminDB := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(adminDSN)))
	defer adminDB.Close()

	var exists bool
	err := adminDB.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)",
		dbname,
	).Scan(&exists)

	if err != nil {
		return nil, fmt.Errorf("database existence check failed: %v", err)
	}

	if !exists {
		_, err = adminDB.Exec("CREATE DATABASE " + dbname)
		if err != nil {
			return nil, fmt.Errorf("database creation failed: %v", err)
		}
		log.Printf("Database %s created", dbname)
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		user, password, host, port, dbname)

	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db := bun.NewDB(sqldb, pgdialect.New())

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("database connection failed: %v", err)
	}

	log.Println("Successfully connected to database")
	return &DatabaseService{db: db}, nil
}

func (s *DatabaseService) GetDB() *bun.DB {
	return s.db
}

func (s *DatabaseService) Close() error {
	return s.db.Close()
}

func (s *DatabaseService) CreateTables(ctx context.Context) error {
	_, err := s.db.NewCreateTable().
		Model((*models.Currency)(nil)).
		IfNotExists().
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to create currencies table: %v", err)
	}

	_, err = s.db.NewCreateTable().
		Model((*models.Price)(nil)).
		IfNotExists().
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to create prices table: %v", err)
	}

	_, err = s.db.NewCreateIndex().
		Model((*models.Price)(nil)).
		Index("idx_prices_currency_timestamp").
		Column("currency_id", "timestamp").
		IfNotExists().
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to create idx_prices_currency_timestamp: %v", err)
	}

	_, err = s.db.NewCreateIndex().
		Model((*models.Currency)(nil)).
		Index("idx_currencies_active").
		Column("is_active").
		IfNotExists().
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to create idx_currencies_active: %v", err)
	}

	log.Println("Tables created successfully")
	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
