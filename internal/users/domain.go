package users

import (
	"github.com/romanchechyotkin/betera-test-task/internal/httpserver"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterDomain(logger *slog.Logger, pool *pgxpool.Pool) httpserver.Handler {
	repo := newRepository(logger, pool)
	h := newHandler(logger, repo)
	return h
}
