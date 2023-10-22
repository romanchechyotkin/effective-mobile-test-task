package postgresql

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/romanchechyotkin/betera-test-task/pkg/logger"
)

type pgConfig struct {
	username string
	password string
	host     string
	port     string
	database string
}

func NewPgConfig(username, password, host, port, database string) *pgConfig {
	return &pgConfig{
		username: username,
		password: password,
		host:     host,
		port:     port,
		database: database,
	}
}

func NewClient(ctx context.Context, log *slog.Logger, cfg *pgConfig) *pgxpool.Pool {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", cfg.username, cfg.password, cfg.host, cfg.port, cfg.database)
	log.Debug("got connection string", slog.String("connection string", connString))

	pool, err := pgxpool.New(ctx, connString)
	err = pool.Ping(ctx)
	if err != nil {
		logger.Error(log, "cannot to connect to postgres", err)
		os.Exit(1)
	}
	log.Debug("postgresql client init", slog.String("client", fmt.Sprintf("%#v", pool)))

	return pool
}

func FormatQuery(q string) string {
	return strings.ReplaceAll(strings.ReplaceAll(q, "\t", ""), "\n", "")
}
