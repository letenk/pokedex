package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/letenk/pokedex/models/domain"
	"github.com/letenk/pokedex/models/web"
	"github.com/letenk/pokedex/usecase"
)

type monsterHandler struct {
	usecase usecase.MonsterUsecase
}

func NewHandlerMonster(usecase usecase.MonsterUsecase) *monsterHandler {
	return &monsterHandler{usecase}
}

func (h *monsterHandler) Create(c *gin.Context) {
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

	// Get payload body
	var req web.MonsterCreateRequest
	err := c.ShouldBind(&req)
	if err != nil {
		errors := web.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}
		response := web.JSONResponseWithData(
			http.StatusBadRequest,
			"error",
			"create monster failed",
			errorMessage,
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	fmt.Println("DEBUG 1")
	fmt.Println("DEBUG 1")
	fmt.Println("DEBUG 1")
	fmt.Println("DEBUG 1")
	// Get image
	file, err := c.FormFile("image")
	if err != nil {
		errorMessage := gin.H{"errors": err}
		response := web.JSONResponseWithData(
			http.StatusBadRequest,
			"error",
			"create monster failed",
			errorMessage,
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	fmt.Println("DEBUG 2")
	fmt.Println("DEBUG 2")
	fmt.Println("DEBUG 2")
	fmt.Println("DEBUG 2")
	// Create path image name
	path := fmt.Sprintf("images/%s-%s", currentUser.ID, file.Filename)
	// Move file image to folder images
	err = c.SaveUploadedFile(file, path)
	if err != nil {
		// errorMessage := gin.H{"errors": err}
		response := web.JSONResponseWithoutData(
			http.StatusBadRequest,
			"error",
			"create monster failed",
			// errorMessage,
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	fmt.Println("DEBUG 3")
	fmt.Println("DEBUG 3")
	fmt.Println("DEBUG 3")
	fmt.Println("DEBUG 3")
	// Create new user
	_, err = h.usecase.Create(c.Request.Context(), req, path)
	if err != nil {
		errorMessage := gin.H{"errors": err}
		response := web.JSONResponseWithData(
			http.StatusBadRequest,
			"error",
			"create monster failed",
			errorMessage,
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	fmt.Println("DEBUG 4")
	fmt.Println("DEBUG 4")
	fmt.Println("DEBUG 4")
	fmt.Println("DEBUG 4")
	// Create format response
	response := web.JSONResponseWithoutData(
		http.StatusCreated,
		"success",
		"Monster has been created",
	)
	c.JSON(http.StatusCreated, response)
}
