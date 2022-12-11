package middleware

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/letenk/pokedex/models/web"
	"github.com/letenk/pokedex/usecase"
	"github.com/letenk/pokedex/util"
)

// Function for auth middleware
func AuthMiddleware(userUsecase usecase.UserUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get header with name `Authorization`
		authHeader := c.GetHeader("Authorization")

		// If inside authHeader doesn't have `Bearer`
		if !strings.Contains(authHeader, "Bearer") {
			// Create format response
			response := web.JSONResponseWithoutData(http.StatusUnauthorized, "error", "unauthorized")
			// Stop process and return response
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		// If there is, create new variable with empty string value
		tokenString := ""
		// Split authHeader with white space
		arrayToken := strings.Split(authHeader, " ")
		// If length arrayToken is same the 2
		if len(arrayToken) == 2 {
			// Get arrayToken with index 1 / only token jwt
			tokenString = arrayToken[1]
		}

		// Parse token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)

			if !ok {
				return nil, errors.New("invalid token")
			}

			// Load Config
			config, err := util.LoadConfig("../.")
			if err != nil {
				log.Fatal("cannot load config:", err)
			}

			return []byte(config.JWT_SECRET_KEY), nil
		})

		// If error
		if err != nil {
			// Create format response
			response := web.JSONResponseWithoutData(http.StatusUnauthorized, "error", "unauthorized")
			// Stop process and return response
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		// Get payload token
		claim, ok := token.Claims.(jwt.MapClaims)
		// If not `ok` and token invalid
		if !ok || !token.Valid {
			// Create format response
			response := web.JSONResponseWithoutData(http.StatusUnauthorized, "error", "unauthorized")
			// Stop process and return response
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		// Get payload `user_id` and convert to `string`
		userId := claim["user_id"].(string)

		// Find user on db with service
		user, err := userUsecase.FindOneByID(context.Background(), userId)
		// If error
		if err != nil {
			// Create format response
			response := web.JSONResponseWithoutData(http.StatusUnauthorized, "error", "unauthorized")
			// Stop process and return response
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		// Set user to context with name `currentUser`
		c.Set("currentUser", user)
	}
}
