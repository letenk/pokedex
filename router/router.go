package router

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/letenk/pokedex/handlers"
	"github.com/letenk/pokedex/middleware"
	"github.com/letenk/pokedex/repository"
	"github.com/letenk/pokedex/usecase"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://*", "http://*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Use layers users
	repositoryUser := repository.NewUserRepository(db)
	usecaseUser := usecase.NewUsecaseUser(repositoryUser)
	handlerUser := handlers.NewHandlerUser(usecaseUser)

	// Use layers category
	repositoryCategory := repository.NewCategoryRepository(db)
	usecaseCategory := usecase.NewUsecaseCategory(repositoryCategory)
	handlerCategory := handlers.NewHandlerCategory(usecaseCategory)

	// Use layers type
	repositoryType := repository.NewTypeRespository(db)
	usecaseType := usecase.NewUsecaseType(repositoryType)
	handlerType := handlers.NewHandlerType(usecaseType)

	// Use layers montser
	repositoryMonster := repository.NewMonsterRespository(db)
	usecaseMonster := usecase.NewUsecaseMonster(repositoryMonster)
	handlerMonster := handlers.NewHandlerMonster(usecaseMonster)

	// Route home
	router.GET("/", func(c *gin.Context) {
		resp := gin.H{"say": "Server is healthy 💪"}
		c.JSON(http.StatusOK, resp)
	})

	// Group api version 1
	v1 := router.Group("/api/v1")

	// Login
	v1.POST("/login", handlerUser.Login)
	// Categories
	v1.GET("/category", middleware.AuthMiddleware(usecaseUser), handlerCategory.FindAll)
	// Types
	v1.GET("/type", middleware.AuthMiddleware(usecaseUser), handlerType.FindAll)

	// Group endpoint monster
	monster := v1.Group("/monster")
	// Find all monster
	monster.GET("", handlerMonster.FindAll)
	// Create monster
	monster.POST("", middleware.AuthMiddleware(usecaseUser), handlerMonster.Create)

	return router
}
