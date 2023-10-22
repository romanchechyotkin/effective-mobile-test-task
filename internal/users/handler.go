package users

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"

	"github.com/romanchechyotkin/effective-mobile-test-task/internal/httpserver"
	"github.com/romanchechyotkin/effective-mobile-test-task/pkg/logger"
)

const (
	AgeApi         = "https://api.agify.io/?name="
	GenderApi      = "https://api.genderize.io/?name="
	NationalityApi = "https://api.nationalize.io/?name="
)

const (
	SORT_BY_ASC_AGE  = "age.a"
	SORT_BY_DESC_AGE = "age.d"
)

type storage interface {
	saveUser(ctx context.Context, dto *UserResponseDto) (string, error)
	getAllUsers(ctx context.Context, opt ...string) ([]*UserResponseDto, error)
	getUser(ctx context.Context, id string) (*UserResponseDto, error)
	updateUser(ctx context.Context, id, col string, val any) error
	deleteUser(ctx context.Context, id string) error
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

	group.POST("/", h.createUser)
	group.GET("/", h.getAllUsers)
	group.GET("/:id", h.getUser)
	group.PATCH("/:id", h.updateUser)
	group.DELETE("/:id", h.deleteUser)
	group.GET("/health", h.index)
}

// @Summary Create user
// @Description Endpoint for creating and saving user to database
// @Produce application/json
// @Success 201 {object} UserResponseDto
// @Router /users [post]
func (h *handler) createUser(ctx *gin.Context) {
	var userDto UserRequestDto
	var wg sync.WaitGroup

	err := ctx.ShouldBindJSON(&userDto)
	if err != nil {
		logger.Error(h.log, "error during decoding user", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	h.log.Debug("decoded user dto", slog.Any("dto", userDto))

	var ageDto AgeRequestDto
	var genderDto GenderRequestDto
	var nationalityDto NationalityRequestDto

	wg.Add(1)
	go func() {
		defer wg.Done()
		resp, err := http.Get(fmt.Sprintf("%s%s", AgeApi, userDto.FirstName))
		if err != nil {
			logger.Error(h.log, "error during request to age api", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		defer resp.Body.Close()

		err = json.NewDecoder(resp.Body).Decode(&ageDto)
		if err != nil {
			logger.Error(h.log, "error during decoding age info", err)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		h.log.Debug("decoded age dto", slog.Any("dto", ageDto))
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		resp, err := http.Get(fmt.Sprintf("%s%s", GenderApi, userDto.FirstName))
		if err != nil {
			logger.Error(h.log, "error during request to gender api", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		defer resp.Body.Close()

		err = json.NewDecoder(resp.Body).Decode(&genderDto)
		if err != nil {
			logger.Error(h.log, "error during decoding gender info", err)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		h.log.Debug("decoded gender dto", slog.Any("dto", genderDto))
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		resp, err := http.Get(fmt.Sprintf("%s%s", NationalityApi, userDto.FirstName))
		if err != nil {
			logger.Error(h.log, "error during request to nationality api", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		defer resp.Body.Close()

		err = json.NewDecoder(resp.Body).Decode(&nationalityDto)
		if err != nil {
			logger.Error(h.log, "error during decoding nationality info", err)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		h.log.Debug("decoded nationality dto", slog.Any("dto", nationalityDto))
	}()

	wg.Wait()

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
	sort := ctx.Query("sort")
	limit := ctx.Query("limit")
	h.log.Debug("got sort query value", slog.String("sort", sort))
	h.log.Debug("got limit query value", slog.String("limit", limit))

	var users []*UserResponseDto
	var err error

	switch sort {
	case SORT_BY_ASC_AGE:
		users, err = h.repository.getAllUsers(ctx, sort, limit)
	case SORT_BY_DESC_AGE:
		users, err = h.repository.getAllUsers(ctx, sort, limit)
	default:
		users, err = h.repository.getAllUsers(ctx, sort, limit)
	}

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

// @Summary Users Endpoint Health Check
// @Description Checking health of users endpoint
// @Produce application/json
// @Success 200 {string} nasa
// @Router /users/health [get]
func (h *handler) index(ctx *gin.Context) {
	ctx.String(http.StatusOK, "users")
}

// @Summary Get exact user
// @Description Endpoint for getting user with exact id
// @Produce application/json
// @Success 200 {object} UserResponseDto
// @Param id path string true "id"
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

// @Summary Update exact user
// @Description Endpoint for updating user with exact id
// @Produce application/json
// @Success 204 {object} UserResponseDto
// @Param id path string true "id"
// @Router /users/{id} [patch]
func (h *handler) updateUser(ctx *gin.Context) {
	id := ctx.Param("id")
	h.log.Debug("got id param", slog.String("id", id))

	var dto map[string]any
	err := ctx.ShouldBindJSON(&dto)
	if err != nil {
		logger.Error(h.log, "error during decoding", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	h.log.Debug("decoded update user dto", slog.Any("dto", dto))

	for k, v := range dto {
		err := h.repository.updateUser(ctx, id, k, v)
		if err != nil {
			logger.Error(h.log, "error during updating in database", err)
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
	}

	ctx.JSON(http.StatusNoContent, gin.H{
		"message": "updated successfully",
	})
}

// @Summary Delete exact user
// @Description Endpoint for deleting user with exact id
// @Produce application/json
// @Success 204 {object} UserResponseDto
// @Param id path string true "id"
// @Router /users/{id} [delete]
func (h *handler) deleteUser(ctx *gin.Context) {
	id := ctx.Param("id")
	h.log.Debug("got id param", slog.String("id", id))

	err := h.repository.deleteUser(ctx, id)
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

	ctx.JSON(http.StatusNoContent, gin.H{
		"message": "deleted successfully",
	})
}
