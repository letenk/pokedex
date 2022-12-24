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

	// File name format
	now := time.Now()
	nowRFC3339 := now.Format(time.RFC3339)
	fileName := fmt.Sprintf(`%s_%v_%s`, currentUser.ID, nowRFC3339, fileHeader.Filename)

	// Create new user
	newMonster, err := h.usecase.Create(c.Request.Context(), req, file, fileName)
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

	// Cache for get detail
	formatResponseJSON := web.FormatMonsterResponseDetail(newMonster)
	key := fmt.Sprintf("monster_id_%s", newMonster.ID)
	go cache.SetWithTTL(key, formatResponseJSON, time.Hour)
	// Remove cache
	go cache.Remove("monsters")

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

	if queryParameter.Name != "" || queryParameter.Catched != "" || queryParameter.Sort != "" || queryParameter.Order != "" || len(queryParameter.Types) != 0 {
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

		formatResponseJSON := web.FormatMonsterResponseList(monsters)

		// Create format response
		response := web.JSONResponseWithData(
			http.StatusOK,
			"success",
			"List of monsters",
			formatResponseJSON,
		)

		c.JSON(http.StatusOK, response)
		return
	}

	key := "monsters"
	monsters, err := cache.Get(key)
	if err == notFound {
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

		formatResponseJSON := web.FormatMonsterResponseList(monsters)

		// Cache data
		go cache.SetWithTTL(key, formatResponseJSON, time.Hour)

		// Create format response
		response := web.JSONResponseWithData(
			http.StatusOK,
			"success",
			"List of monsters",
			formatResponseJSON,
		)

		c.JSON(http.StatusOK, response)
		return
	}

	jsonResponse := web.JSONResponseWithData(
		http.StatusOK,
		"success",
		"List of monsters",
		monsters,
	)
	c.JSON(http.StatusOK, jsonResponse)
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

	// Get from cache
	key := fmt.Sprintf("monster_id_%s", monsterID.ID)
	monsters, err := cache.Get(key)
	if err == notFound {
		// Find by id monster from database
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

		// Cache for get detail
		formatResponseJSON := web.FormatMonsterResponseDetail(monsters)
		key := fmt.Sprintf("monster_id_%s", monsterID.ID)
		go cache.SetWithTTL(key, formatResponseJSON, time.Hour)

		// Create format response
		response := web.JSONResponseWithData(
			http.StatusOK,
			"success",
			"profile detail of monsters",
			formatResponseJSON,
		)

		c.JSON(http.StatusOK, response)
		return
	}

	// Create format response
	response := web.JSONResponseWithData(
		http.StatusOK,
		"success",
		"profile detail of monsters",
		monsters,
	)

	c.JSON(http.StatusOK, response)
}

func (h *monsterHandler) Update(c *gin.Context) {
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

	// Get id monster from path
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

	// Get payload body
	var reqUpdate web.MonsterUpdateRequest
	err = c.ShouldBind(&reqUpdate)
	if err != nil {
		errors := web.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}
		response := web.JSONResponseWithData(
			http.StatusBadRequest,
			"error",
			"update monster failed",
			errorMessage,
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// Get file image
	file, fileHeader, err := c.Request.FormFile("image")
	if err != nil && err.Error() != "http: no such file" {
		errorMessage := gin.H{"errors": err}
		response := web.JSONResponseWithData(
			http.StatusBadRequest,
			"error",
			"update monster failed",
			errorMessage,
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	var fileName string
	if file != nil {
		// Validate
		valid, message := validateUploadFiles(fileHeader)
		if !valid {
			errorMessage := gin.H{"errors": message}
			response := web.JSONResponseWithData(
				http.StatusBadRequest,
				"error",
				"updated monster failed",
				errorMessage,
			)
			c.JSON(http.StatusBadRequest, response)
			return
		}

		// File name format
		now := time.Now()
		nowRFC3339 := now.Format(time.RFC3339)
		fileName = fmt.Sprintf(`%s_%v_%s`, currentUser.ID, nowRFC3339, fileHeader.Filename)
	}

	// Update
	monsterUpdated, err := h.usecase.Update(c.Request.Context(), monsterID.ID, reqUpdate, file, fileName)
	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}
		response := web.JSONResponseWithData(
			http.StatusBadRequest,
			"error",
			"update monster failed",
			errorMessage,
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// Cache for get detail
	formatResponseJSON := web.FormatMonsterResponseDetail(monsterUpdated)
	key := fmt.Sprintf("monster_id_%s", monsterID.ID)
	// Remove cache
	go cache.Remove(key)
	go cache.SetWithTTL(key, formatResponseJSON, time.Hour)
	// Remove cache
	go cache.Remove("monsters")

	// Create format response
	response := web.JSONResponseWithData(
		http.StatusOK,
		"success",
		"Update monster success",
		formatResponseJSON,
	)

	c.JSON(http.StatusOK, response)
}

func (h *monsterHandler) UpdateMarkMonsterCaptured(c *gin.Context) {
	// Check Authorization
	// Get current user login
	currentUser := c.MustGet("currentUser").(domain.User)
	if currentUser.Role != "user" {
		response := web.JSONResponseWithoutData(
			http.StatusForbidden,
			"error",
			"forbidden",
		)
		c.JSON(http.StatusForbidden, response)
		return
	}

	// Get id monster from path
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

	// Get payload body
	var reqUpdate web.MonsterUpdateRequestMonsterCapture
	err = c.ShouldBind(&reqUpdate)
	if err != nil {
		errors := web.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}
		response := web.JSONResponseWithData(
			http.StatusBadRequest,
			"error",
			"update monster captured failed",
			errorMessage,
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// Update
	_, err = h.usecase.UpdateMarkMonsterCaptured(c.Request.Context(), monsterID.ID, reqUpdate)
	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}
		response := web.JSONResponseWithData(
			http.StatusBadRequest,
			"error",
			"update monster captured failed",
			errorMessage,
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	key := fmt.Sprintf("monster_id_%s", monsterID.ID)
	// Remove cache
	go cache.Remove(key)
	// Remove cache
	go cache.Remove("monsters")

	msgSuccess := fmt.Sprintf("monster with id %s updated", monsterID.ID)
	// Create format response
	response := web.JSONResponseWithoutData(
		http.StatusOK,
		"success",
		msgSuccess,
	)

	c.JSON(http.StatusOK, response)
}

func (h *monsterHandler) Delete(c *gin.Context) {
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

	// Get id monster from path
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

	// Delete
	_, err = h.usecase.Delete(c.Request.Context(), monsterID.ID)
	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}
		response := web.JSONResponseWithData(
			http.StatusBadRequest,
			"error",
			"delete monster failed",
			errorMessage,
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	key := fmt.Sprintf("monster_id_%s", monsterID.ID)
	// Remove cache
	go cache.Remove(key)
	// Remove cache
	go cache.Remove("monsters")

	// Create format response
	response := web.JSONResponseWithoutData(
		http.StatusOK,
		"success",
		"monster deleted",
	)

	c.JSON(http.StatusOK, response)
}
