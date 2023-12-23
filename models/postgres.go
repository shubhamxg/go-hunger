package models

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	SSLMode  string
}

func (cfg PostgresConfig) New() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.Database,
		cfg.SSLMode,
	)
}

func (cfg PostgresConfig) Default() string {
	cfg = PostgresConfig{
		Host:     Env(DB_HOST),
		Port:     Env(DB_PORT),
		User:     Env(DB_USER),
		Password: Env(DB_PASSWORD),
		Database: Env(DB_DATABASE),
		SSLMode:  "disable",
	}
	return cfg.New()
}

func (cfg PostgresConfig) Open() (*sqlx.DB, error) {
	// db, err := sql.Open("pgx", cfg.Default())
	db, err := sqlx.Connect("postgres", cfg.Default())
	if err != nil {
		return nil, fmt.Errorf("open: %w", err)
	}
	return db, nil
}

func Start() *sqlx.DB {
	postgres := PostgresConfig{}
	postgres.Default()
	db, err := postgres.Open()
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected to PostgresDB")
	return db
}
