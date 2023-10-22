package httpserver

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	swagger "github.com/swaggo/gin-swagger"

	_ "github.com/romanchechyotkin/betera-test-task/docs"
	"github.com/romanchechyotkin/betera-test-task/pkg/logger"
)

type Handler interface {
	RegisterRoutes(engine *gin.Engine)
}

func Run(log *slog.Logger, usersHandler Handler) {
	engine := gin.Default()
	engine.Use(CORSMiddleware())

	registerGinRoutes(engine)
	usersHandler.RegisterRoutes(engine)

	err := engine.Run(":8080")
	if err != nil {
		logger.Error(log, "http server init failed", err)
		os.Exit(1)
	}
}

func registerGinRoutes(engine *gin.Engine) {
	engine.GET("/swagger/*any", swagger.WrapHandler(swaggerFiles.Handler))
	engine.GET("/health", health)
}

// @Summary Health Check
// @Description Checking health of backend
// @Produce application/json
// @Success 200 {string} health
// @Router /health [get]
func health(ctx *gin.Context) {
	ctx.String(http.StatusOK, "health")
}
