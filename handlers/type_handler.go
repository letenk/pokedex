package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
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

	// Get all
	types, err := h.usecase.FindAll(c.Request.Context())
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
		"list of types",
		web.FormatTypesResponse(types),
	)
	c.JSON(http.StatusOK, jsonResponse)
}
