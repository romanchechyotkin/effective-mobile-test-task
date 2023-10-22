package users

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

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

func (r *repository) getAllUsers(ctx context.Context, opt ...string) ([]*UserResponseDto, error) {
	var orderBy, query string
	limit := 3 // default limit

	if len(opt) != 0 {
		orderBy = opt[0]
	}

	if opt[1] != "" {
		i, _ := strconv.ParseInt(opt[1], 10, 64)
		limit = int(i)
	}

	switch orderBy {
	case SORT_BY_ASC_AGE:
		query = `
		SELECT id, last_name, first_name, second_name, age, gender, nationality 
		FROM effective.public.users
		ORDER BY age
		LIMIT $1
	`
	case SORT_BY_DESC_AGE:
		query = `
		SELECT id, last_name, first_name, second_name, age, gender, nationality 
		FROM effective.public.users
		ORDER BY age DESC 
		LIMIT $1
	`
	default:
		query = `
		SELECT id, last_name, first_name, second_name, age, gender, nationality 
		FROM effective.public.users	
		ORDER BY created_at 
		LIMIT $1
	`
	}

	r.log.Info("database query", slog.String("query", postgresql.FormatQuery(query)))
	rows, err := r.pool.Query(ctx, query, limit)
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

func (r *repository) updateUser(ctx context.Context, id, col string, val any) error {
	query := fmt.Sprintf(`
		UPDATE users
		SET %s = $1
		WHERE id = $2
`, col)

	r.log.Info("database query", slog.String("query", postgresql.FormatQuery(query)))
	exec, err := r.pool.Exec(ctx, query, val, id)
	if err != nil {
		logger.Error(r.log, "error during execution", err)
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		} else {
			return err
		}
	}
	r.log.Info("result of execution", slog.Int("rows affected", int(exec.RowsAffected())))

	return nil
}

func (r *repository) deleteUser(ctx context.Context, id string) error {
	query := `
		DELETE FROM effective.public.users
		WHERE id = $1
	`

	r.log.Info("database query", slog.String("query", postgresql.FormatQuery(query)))
	exec, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		logger.Error(r.log, "error during scanning", err)
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		} else {
			return err
		}
	}
	r.log.Info("result of execution", slog.Int("rows affected", int(exec.RowsAffected())))

	return nil
}
