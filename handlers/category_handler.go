package handlers

import (
	"net/http"

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

	// Get all
	categories, err := h.usecase.FindAll(c.Request.Context())
	if err != nil {
		response := web.JSONResponseWithoutData(
			http.StatusBadRequest,
			"error",
			"bad request",
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	jsonResponse := web.JSONResponseWithData(
		http.StatusOK,
		"success",
		"list of category",
		web.FormatCategoriesResponse(categories),
	)
	c.JSON(http.StatusOK, jsonResponse)
}
