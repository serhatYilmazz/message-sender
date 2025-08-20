package main

import (
	"github.com/serhatYilmazz/message-sender/internal/config"
	"github.com/serhatYilmazz/message-sender/internal/message"
	"github.com/serhatYilmazz/message-sender/pkg/db"
	"github.com/serhatYilmazz/message-sender/pkg/log"
	"os"
)

func main() {
	env := os.Getenv("environment")
	logger := log.NewLogger(env)

	cfg, err := config.Load(logger, "configs")
	if err != nil {
		return
	}

	postgresDb, err := db.NewPostgresDb(cfg.DbConfig)
	if err != nil {
		logger.Fatal("db connection is failed:", err)
	}

	pgRepository := message.PgRepository{
		Db:     postgresDb,
		Logger: logger,
	}
}
