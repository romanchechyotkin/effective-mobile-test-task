package main

import (
	"context"
	"os"

	"github.com/romanchechyotkin/betera-test-task/internal/httpserver"
	"github.com/romanchechyotkin/betera-test-task/internal/users"
	"github.com/romanchechyotkin/betera-test-task/pkg/logger"
	"github.com/romanchechyotkin/betera-test-task/pkg/postgresql"
)

// @title Swagger Documentation
// @version 1.0
// @description Effective Mobile test task in Gin Framework
// @host localhost:8080
func main() {
	log := logger.New(os.Stdout)
	log.Debug("app running")

	pgConfig := postgresql.NewPgConfig(
		"chechyotka",
		"5432",
		"localhost",
		"5432",
		"effective",
	)

	//pgConfig := postgresql.NewPgConfig(
	//	os.Getenv("POSTGRES_USER"),
	//	os.Getenv("POSTGRES_PASSWORD"),
	//	os.Getenv("POSTGRES_HOST"),
	//	os.Getenv("POSTGRES_PORT"),
	//	os.Getenv("POSTGRES_DB"),
	//)

	pgClient := postgresql.NewClient(context.Background(), log, pgConfig)
	usersDomain := users.RegisterDomain(log, pgClient)

	httpserver.Run(log, usersDomain)
}
