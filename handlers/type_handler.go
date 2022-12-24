package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jellydator/ttlcache/v2"
	"github.com/letenk/pokedex/models/domain"
	"github.com/letenk/pokedex/models/web"
	"github.com/letenk/pokedex/usecase"
)

type typeHandler struct {
	usecase usecase.TypeUsecase
}

func NewHandlerType(usecase usecase.TypeUsecase) *typeHandler {
	return &typeHandler{usecase}
}

var (
	cache    ttlcache.SimpleCache = ttlcache.NewCache()
	notFound                      = ttlcache.ErrNotFound
)

func (h *typeHandler) FindAll(c *gin.Context) {
	// Check Authorization
	// Get current user login
	currentUser := c.MustGet("currentUser").(domain.User)
	if currentUser.Role != "admin" {
		response := web.JSONResponseWithoutData(
			http.StatusForbidden,
			"error",
			"forbidden",
		)
		c.JSON(http.StatusForbidden, response)
		return
	}

	// Create context
	ctx := c.Request.Context()

	key := "types"
	types, err := cache.Get(key)
	if err == notFound {
		// Get all data from db
		types, err := h.usecase.FindAll(ctx)
		if err != nil {
			response := web.JSONResponseWithoutData(
				http.StatusBadRequest,
				"error",
				"bad request",
			)
			c.JSON(http.StatusBadRequest, response)
			return
		}

		formatResponseJSON := web.FormatTypesResponse(types)

		// Cache data
		go cache.SetWithTTL(key, formatResponseJSON, time.Hour)

		jsonResponse := web.JSONResponseWithData(
			http.StatusOK,
			"success",
			"list of types",
			formatResponseJSON,
		)
		c.JSON(http.StatusOK, jsonResponse)
		return
	}

	jsonResponse := web.JSONResponseWithData(
		http.StatusOK,
		"success",
		"list of types",
		types,
	)
	c.JSON(http.StatusOK, jsonResponse)
}
