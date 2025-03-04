package db

import (
	"database/sql"
	"fmt"

	"github.com/gabrigabs/campaign-message-consumer/config"
	"github.com/gabrigabs/campaign-message-consumer/logger"
	_ "github.com/lib/pq"
)

type Postgres struct {
	DB     *sql.DB
	logger logger.Logger
}

func NewPostgresConnection(cfg config.PostgresConfig, log logger.Logger) (*Postgres, error) {
	log.Info("Connecting to PostgreSQL", map[string]any{
		"host":     cfg.Host,
		"port":     cfg.Port,
		"database": cfg.Database,
		"user":     cfg.User,
		"ssl":      cfg.SslMode,
	})

	sslMode := "disable"

	if cfg.SslMode {
		sslMode = "require"
	}

	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database, sslMode,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Error("Failed to connect to PostgreSQL", map[string]any{
			"error": err.Error(),
		})
		return nil, err
	}

	log.Info("Connected to PostgreSQL", map[string]any{
		"database": cfg.Database,
	})

	return &Postgres{
		DB:     db,
		logger: log,
	}, nil
}
