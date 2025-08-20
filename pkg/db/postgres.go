package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/serhatYilmazz/message-sender/internal/config"
)

func NewPostgresDb(cfg config.DbConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DbName, cfg.SslMode,
	)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
