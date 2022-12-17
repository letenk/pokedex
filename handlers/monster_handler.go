package handlers

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

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

const (
	// Max file size : 2MB
	maxPartSize = int64(5 * 1024 * 1024)
)

func validateUploadFiles(fileHeader *multipart.FileHeader) (bool, string) {
	size := fileHeader.Size
	extension := strings.Split(fileHeader.Filename, ".")

	if size > maxPartSize {
		return false, "Files cannot exceed 2MB"
	}

	if extension[1] != "jpeg" && extension[1] != "jpg" && extension[1] != "png" {
		return false, "File must be format jpeg or png"
	}
	return true, "ok"
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

	// Get file image
	file, fileHeader, err := c.Request.FormFile("image")
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

	// Validate
	valid, message := validateUploadFiles(fileHeader)
	if !valid {
		errorMessage := gin.H{"errors": message}
		response := web.JSONResponseWithData(
			http.StatusBadRequest,
			"error",
			"create monster failed",
			errorMessage,
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// File name formate
	now := time.Now()
	nowRFC3339 := now.Format(time.RFC3339)
	fileName := fmt.Sprintf(`%s-%v-%s`, currentUser.ID, nowRFC3339, fileHeader.Filename)

	// Create new user
	_, err = h.usecase.Create(c.Request.Context(), req, file, fileName)
	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}
		response := web.JSONResponseWithData(
			http.StatusBadRequest,
			"error",
			"create monster failed",
			errorMessage,
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// Create format response
	response := web.JSONResponseWithoutData(
		http.StatusCreated,
		"success",
		"Monster has been created",
	)
	c.JSON(http.StatusCreated, response)
}

func (h *monsterHandler) FindAll(c *gin.Context) {
	// Get query
	var queryParameter web.MonsterQueryRequest
	err := c.Bind(&queryParameter)
	if err != nil {
		response := web.JSONResponseWithoutData(
			http.StatusInternalServerError,
			"error",
			"internal server error",
		)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// Find all montser
	monsters, err := h.usecase.FindAll(c.Request.Context(), queryParameter)
	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}
		response := web.JSONResponseWithData(
			http.StatusBadRequest,
			"error",
			"bad request",
			errorMessage,
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// Create format response
	response := web.JSONResponseWithData(
		http.StatusOK,
		"success",
		"List of monsters",
		web.FormatMonsterResponseList(monsters),
	)

	c.JSON(http.StatusOK, response)
}

func (h *monsterHandler) FindByID(c *gin.Context) {
	var monsterID web.MosterURI
	err := c.ShouldBindUri(&monsterID)
	if err != nil {
		response := web.JSONResponseWithoutData(
			http.StatusInternalServerError,
			"error",
			"internal server error",
		)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// Find by id monster
	monsters, err := h.usecase.FindByID(c.Request.Context(), monsterID.ID)
	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}
		response := web.JSONResponseWithData(
			http.StatusBadRequest,
			"error",
			"bad request",
			errorMessage,
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// Create format response
	response := web.JSONResponseWithData(
		http.StatusOK,
		"success",
		"profile detail of monsters",
		web.FormatMonsterResponseDetail(monsters),
	)

	c.JSON(http.StatusOK, response)
}
