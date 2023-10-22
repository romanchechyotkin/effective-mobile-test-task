package users

import (
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/romanchechyotkin/effective-mobile-test-task/internal/httpserver"
)

func RegisterDomain(logger *slog.Logger, pool *pgxpool.Pool) httpserver.Handler {
	repo := newRepository(logger, pool)
	h := newHandler(logger, repo)
	return h
}
