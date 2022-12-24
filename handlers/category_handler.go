package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/letenk/pokedex/models/domain"
	"github.com/letenk/pokedex/models/web"
	"github.com/letenk/pokedex/usecase"
)

type categoryHandler struct {
	usecase usecase.CategoryUsecase
}

func NewHandlerCategory(usecase usecase.CategoryUsecase) *categoryHandler {
	return &categoryHandler{usecase}
}

func (h *categoryHandler) FindAll(c *gin.Context) {
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

	key := "categories"
	categories, err := cache.Get(key)
	if err == notFound {
		// Get all
		categories, err := h.usecase.FindAll(ctx)
		if err != nil {
			response := web.JSONResponseWithoutData(
				http.StatusBadRequest,
				"error",
				"bad request",
			)
			c.JSON(http.StatusBadRequest, response)
			return
		}

		formatResponseJSON := web.FormatCategoriesResponse(categories)

		// Cache data
		go cache.SetWithTTL(key, formatResponseJSON, time.Hour)

		jsonResponse := web.JSONResponseWithData(
			http.StatusOK,
			"success",
			"list of category",
			formatResponseJSON,
		)
		c.JSON(http.StatusOK, jsonResponse)
		return
	}

	jsonResponse := web.JSONResponseWithData(
		http.StatusOK,
		"success",
		"list of types",
		categories,
	)
	c.JSON(http.StatusOK, jsonResponse)
}
