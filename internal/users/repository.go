package users

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"

	"github.com/romanchechyotkin/effective-mobile-test-task/pkg/logger"
	"github.com/romanchechyotkin/effective-mobile-test-task/pkg/postgresql"
)

var ErrNotFound = errors.New("not found")

type repository struct {
	log  *slog.Logger
	pool *pgxpool.Pool
}

func newRepository(logger *slog.Logger, pool *pgxpool.Pool) storage {
	return &repository{
		log:  logger,
		pool: pool,
	}
}

func (r *repository) saveUser(ctx context.Context, dto *UserResponseDto) (string, error) {
	query := `
		INSERT INTO users (last_name, first_name, second_name, age, gender, nationality)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	var id string
	r.log.Info("database query", slog.String("query", postgresql.FormatQuery(query)))
	err := r.pool.QueryRow(ctx, query, dto.LastName, dto.FirstName, dto.SecondName, dto.Age, dto.Gender, dto.Nationality).Scan(&id)
	if err != nil {
		logger.Error(r.log, "error during execution", err)
		return "", err
	}

	return id, nil
}

func (r *repository) getAllUsers(ctx context.Context) ([]*UserResponseDto, error) {
	query := `
		SELECT id, last_name, first_name, second_name, age, gender, nationality FROM effective.public.users
	`

	r.log.Info("database query", slog.String("query", postgresql.FormatQuery(query)))
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		logger.Error(r.log, "error during query", err)
		return nil, err
	}
	defer rows.Close()

	var res []*UserResponseDto
	for rows.Next() {
		var dto UserResponseDto
		err = rows.Scan(&dto.ID, &dto.LastName, &dto.FirstName, &dto.SecondName, &dto.Age, &dto.Gender, &dto.Nationality)
		if err != nil {
			logger.Error(r.log, "error during scanning", err)
			return nil, err
		}

		res = append(res, &dto)
	}

	return res, nil
}

func (r *repository) getUser(ctx context.Context, id string) (*UserResponseDto, error) {
	query := `
		SELECT id, last_name, first_name, second_name, age, gender, nationality
		FROM effective.public.users
		WHERE id = $1
	`

	var dto UserResponseDto

	r.log.Info("database query", slog.String("query", postgresql.FormatQuery(query)))
	err := r.pool.QueryRow(ctx, query, id).Scan(&dto.ID, &dto.LastName, &dto.FirstName, &dto.SecondName, &dto.Age, &dto.Gender, &dto.Nationality)
	if err != nil {
		logger.Error(r.log, "error during scanning", err)
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		} else {
			return nil, err
		}
	}

	return &dto, nil
}
