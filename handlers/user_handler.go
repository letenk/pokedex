package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/letenk/pokedex/models/web"
	"github.com/letenk/pokedex/usecase"
)

type userHandler struct {
	usecase usecase.UserUsecase
}

func NewHandlerUser(usecase usecase.UserUsecase) *userHandler {
	return &userHandler{usecase}
}

func (h *userHandler) Login(c *gin.Context) {
	var req web.UserLoginRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		errors := web.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}
		response := web.JSONResponseWithData(
			http.StatusBadRequest,
			"error",
			"login failed",
			errorMessage,
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// Login
	token, err := h.usecase.Login(c.Request.Context(), req)
	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}
		response := web.JSONResponseWithData(
			http.StatusBadRequest,
			"error",
			"login failed",
			errorMessage,
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	fieldToken := gin.H{"token": token}
	// Create format response
	response := web.JSONResponseWithData(
		http.StatusOK,
		"success",
		"login success",
		fieldToken,
	)
	c.JSON(http.StatusOK, response)
}
