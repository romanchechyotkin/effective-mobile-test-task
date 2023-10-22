package users

import (
	"context"
	"errors"
	"github.com/romanchechyotkin/betera-test-task/internal/httpserver"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/romanchechyotkin/betera-test-task/pkg/logger"
)

type storage interface {
	saveAPOD(ctx context.Context, dto *User) error
	getAllAPODs(ctx context.Context) ([]*User, error)
	getAPOD(ctx context.Context, date string) (*User, error)
}

type handler struct {
	log        *slog.Logger
	repository storage
}

func newHandler(logger *slog.Logger, repo storage) httpserver.Handler {
	h := &handler{
		log:        logger,
		repository: repo,
	}

	return h
}

func (h *handler) RegisterRoutes(engine *gin.Engine) {
	group := engine.Group("/nasa")

	group.GET("/", h.getAllAPODs)
	group.GET("/:date", h.getAPOD)
	group.GET("/health", h.index)
}

// @Summary The whole album
// @Description Endpoint for getting the whole album
// @Produce application/json
// @Success 200 {object} []Metadata{}
// @Router /nasa [get]
func (h *handler) getAllAPODs(ctx *gin.Context) {
	apods, err := h.repository.getAllAPODs(ctx)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}
		return
	}

	ctx.JSON(http.StatusOK, apods)
}

//func (h *handler) parseMetadata() {
//	ticker := time.NewTicker(time.Hour * 24)
//	for {
//		select {
//		case <-ticker.C:
//			resp, err := http.Get("https://api.nasa.gov/planetary/apod?api_key=DEMO_KEY")
//			if err != nil {
//				logger.Error(h.log, "error during request to NASA API", err)
//				return
//			}
//
//			var dto Metadata
//			err = json.NewDecoder(resp.Body).Decode(&dto)
//			if err != nil {
//				logger.Error(h.log, "error during json decoding", err)
//				return
//			}
//
//			h.log.Info("got response", slog.Any("metadata", dto))
//			if dto.URL == "" {
//				logger.Error(h.log, "empty url", nil)
//				return
//			}
//
//			resp, err = http.Get(dto.URL)
//			if err != nil {
//				logger.Error(h.log, "error during request to get image", err)
//				return
//			}
//
//			file, err := os.Create("tmp.jpg")
//			if err != nil {
//				logger.Error(h.log, "failed to create image file", err)
//				return
//			}
//			defer file.Close()
//
//			_, err = io.Copy(file, resp.Body)
//			if err != nil {
//				logger.Error(h.log, "failed to create image file", err)
//				return
//			}
//
//			fileName := uuid.New().String() + ".jpg"
//
//			err = h.saveToMinio(context.Background(), h.minioClient, fileName, "tmp.jpg")
//			if err != nil {
//				logger.Error(h.log, "error during saving to minio", err)
//				return
//			}
//
//			dto.URL = fileName
//
//			err = h.repository.saveAPOD(context.Background(), &dto)
//			if err != nil {
//				logger.Error(h.log, "error during saving apod to database", err)
//				return
//			}
//		}
//	}
//}

// @Summary Nasa Endpoint Health Check
// @Description Checking health of nasa endpoint
// @Produce application/json
// @Success 200 {string} nasa
// @Router /nasa/health [get]
func (h *handler) index(ctx *gin.Context) {
	ctx.String(http.StatusOK, "nasa")
}

// @Summary The exact APOD
// @Description Endpoint for getting the APOD with exact date
// @Produce application/json
// @Success 200 {object} Metadata
// @Param date path string true "Date"
// @Router /nasa/{date} [get]
func (h *handler) getAPOD(ctx *gin.Context) {
	date := ctx.Param("date")
	h.log.Info("got date param", slog.String("date", date))

	apod, err := h.repository.getAPOD(ctx, date)
	if err != nil {
		logger.Error(h.log, "error during db query", err)
		if errors.Is(err, ErrNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}
		return
	}

	ctx.JSON(http.StatusOK, apod)
}

//func (h *handler) saveToMinio(ctx context.Context, client *minio.Client, fileName, filePath string) error {
//	info, err := client.FPutObject(ctx, "betera", fileName, filePath, minio.PutObjectOptions{ContentType: "image/jpg"})
//	if err != nil {
//		return err
//	}
//
//	h.log.Info("image successfully uploaded", slog.String("info", fmt.Sprintf("%s of size %d", fileName, info.Size)))
//
//	return nil
//}
