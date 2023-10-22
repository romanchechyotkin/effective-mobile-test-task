package users

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/romanchechyotkin/effective-mobile-test-task/internal/httpserver"
	"github.com/romanchechyotkin/effective-mobile-test-task/pkg/logger"
)

const (
	AgeApi         = "https://api.agify.io/?name="
	GenderApi      = "https://api.genderize.io/?name="
	NationalityApi = "https://api.nationalize.io/?name="
)

type storage interface {
	saveUser(ctx context.Context, dto *UserResponseDto) (string, error)
	getAllUsers(ctx context.Context) ([]*UserResponseDto, error)
	getUser(ctx context.Context, id string) (*UserResponseDto, error)
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
	group := engine.Group("/users")

	group.POST("/", h.CreateUser)
	group.GET("/", h.getAllUsers)
	group.GET("/:id", h.getUser)
	group.GET("/health", h.index)
}

// @Summary Create user
// @Description Endpoint for creating and saving user to database
// @Produce application/json
// @Success 201 {object} User
// @Router /users [post]
func (h *handler) CreateUser(ctx *gin.Context) {
	var userDto UserRequestDto

	err := ctx.ShouldBindJSON(&userDto)
	if err != nil {
		logger.Error(h.log, "error during decoding user", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	h.log.Debug("decoded user dto", slog.Any("dto", userDto))

	resp, err := http.Get(fmt.Sprintf("%s%s", AgeApi, userDto.FirstName))
	if err != nil {
		logger.Error(h.log, "error during request to age api", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	defer resp.Body.Close()

	var ageDto AgeRequestDto
	err = json.NewDecoder(resp.Body).Decode(&ageDto)
	if err != nil {
		logger.Error(h.log, "error during decoding age info", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	h.log.Debug("decoded age dto", slog.Any("dto", ageDto))

	resp, err = http.Get(fmt.Sprintf("%s%s", GenderApi, userDto.FirstName))
	if err != nil {
		logger.Error(h.log, "error during request to gender api", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	defer resp.Body.Close()

	var genderDto GenderRequestDto
	err = json.NewDecoder(resp.Body).Decode(&genderDto)
	if err != nil {
		logger.Error(h.log, "error during decoding gender info", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	h.log.Debug("decoded gender dto", slog.Any("dto", genderDto))

	resp, err = http.Get(fmt.Sprintf("%s%s", NationalityApi, userDto.FirstName))
	if err != nil {
		logger.Error(h.log, "error during request to nationality api", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	defer resp.Body.Close()

	var nationalityDto NationalityRequestDto
	err = json.NewDecoder(resp.Body).Decode(&nationalityDto)
	if err != nil {
		logger.Error(h.log, "error during decoding nationality info", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	h.log.Debug("decoded nationality dto", slog.Any("dto", nationalityDto))

	response := &UserResponseDto{
		LastName:    userDto.LastName,
		FirstName:   userDto.FirstName,
		SecondName:  userDto.SecondName,
		Age:         ageDto.Age,
		Gender:      genderDto.Gender,
		Nationality: nationalityDto.Country[0].CountryID,
	}

	id, err := h.repository.saveUser(ctx, response)
	if err != nil {
		logger.Error(h.log, "error during saving user to database", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	response.ID = id

	h.log.Info("user created", slog.Any("user", response))
	ctx.JSON(http.StatusCreated, response)
}

// @Summary All users
// @Description Endpoint for getting all users
// @Produce application/json
// @Success 200 {object} []UserResponseDto{}
// @Router /users [get]
func (h *handler) getAllUsers(ctx *gin.Context) {
	users, err := h.repository.getAllUsers(ctx)
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

	ctx.JSON(http.StatusOK, users)
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

// @Summary Users Endpoint Health Check
// @Description Checking health of users endpoint
// @Produce application/json
// @Success 200 {string} nasa
// @Router /users/health [get]
func (h *handler) index(ctx *gin.Context) {
	ctx.String(http.StatusOK, "nasa")
}

// @Summary The exact User
// @Description Endpoint for getting user with exact id
// @Produce application/json
// @Success 200 {object} UserResponseDto
// @Param date path string true "Date"
// @Router /users/{id} [get]
func (h *handler) getUser(ctx *gin.Context) {
	id := ctx.Param("id")
	h.log.Debug("got id param", slog.String("id", id))

	user, err := h.repository.getUser(ctx, id)
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

	ctx.JSON(http.StatusOK, user)
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
