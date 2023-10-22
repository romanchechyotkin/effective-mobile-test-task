package users

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

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

func (r *repository) saveUser(ctx context.Context, dto *UserResponseDto) error {
	query := `
		INSERT INTO users (last_name, first_name, second_name, age, gender, nationality)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	r.log.Info("database query", slog.String("query", postgresql.FormatQuery(query)))
	exec, err := r.pool.Exec(ctx, query, dto.LastName, dto.FirstName, dto.SecondName, dto.Age, dto.Gender, dto.Nationality)
	if err != nil {
		logger.Error(r.log, "error during execution", err)
		return err
	}
	r.log.Info("result of execution", slog.String("result", fmt.Sprintf("rows affected %d", exec.RowsAffected())))

	return nil
}

func (r *repository) getAllAPODs(ctx context.Context) ([]*UserResponseDto, error) {
	//query := `
	//	SELECT title, explanation, image, media_type, service_version, date FROM apods
	//`
	//
	//r.log.Info("database query", slog.String("query", postgresql.FormatQuery(query)))
	//rows, err := r.pool.Query(ctx, query)
	//if err != nil {
	//	logger.Error(r.log, "error during query", err)
	//	return nil, err
	//}
	//defer rows.Close()
	//
	//var res []*Metadata
	//for rows.Next() {
	//	var m Metadata
	//
	//	var timestamp time.Time
	//	err = rows.Scan(&m.Title, &m.Explanation, &m.URL, &m.MediaType, &m.ServiceVersion, &timestamp)
	//	if err != nil {
	//		logger.Error(r.log, "error during scanning", err)
	//		return nil, err
	//	}
	//
	//	m.URL = "http://localhost:9000/betera/" + m.URL
	//	m.Date = timestamp.String()
	//	res = append(res, &m)
	//}

	return nil, nil
}

func (r *repository) getAPOD(ctx context.Context, date string) (*UserResponseDto, error) {
	//query := `
	//	SELECT title, explanation, image, media_type, service_version, date
	//	FROM apods
	//	WHERE date = $1
	//`
	//
	//var m Metadata
	//var timestamp time.Time
	//
	//r.log.Info("database query", slog.String("query", postgresql.FormatQuery(query)))
	//err := r.pool.QueryRow(ctx, query, date).Scan(&m.Title, &m.Explanation, &m.URL, &m.MediaType, &m.ServiceVersion, &timestamp)
	//if err != nil {
	//	logger.Error(r.log, "error during scanning", err)
	//	if errors.Is(err, pgx.ErrNoRows) {
	//		return nil, ErrNotFound
	//	} else {
	//		return nil, err
	//	}
	//}
	//
	//m.URL = "http://localhost:9000/betera/" + m.URL
	//m.Date = timestamp.String()

	return nil, nil
}
