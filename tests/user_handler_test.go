package tests

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/letenk/pokedex/models/web"
	"github.com/stretchr/testify/require"
)

func TestLoginUserHandler(t *testing.T) {
	// Test Cases
	testCases := []struct {
		name string
		req  web.UserLoginRequest
	}{
		{
			name: "success_login",
			req: web.UserLoginRequest{
				Username: "admin",
				Password: "password",
			},
		},
		{
			name: "failed_login_wrong_username",
			req: web.UserLoginRequest{
				Username: "wrong",
				Password: "password",
			},
		},
		{
			name: "failed_login_wrong_password",
			req: web.UserLoginRequest{
				Username: "admin",
				Password: "wrong",
			},
		},
		{
			name: "failed_login_validation_error",
			req: web.UserLoginRequest{
				Username: "",
				Password: "",
			},
		},
	}
	// Test
	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Data body
			dataBody := fmt.Sprintf(`{"username": "%s", "password": "%s"}`, tc.req.Username, tc.req.Password)
			// New reader
			requestBody := strings.NewReader(dataBody)
			// Create new request
			request := httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/v1/login", requestBody)
			// Added header content type
			request.Header.Add("Content-Type", "application/json")
			// Create new recorder
			recorder := httptest.NewRecorder()

			// Run http test
			RouteTest.ServeHTTP(recorder, request)

			// Get response
			response := recorder.Result()

			// Read all response
			body, _ := io.ReadAll(response.Body)
			var responseBody map[string]interface{}
			json.Unmarshal(body, &responseBody)

			// Check pass test
			if tc.name == "success_login" {
				require.Equal(t, 200, response.StatusCode)
				require.Equal(t, 200, int(responseBody["code"].(float64)))
				require.Equal(t, "success", responseBody["status"])
				require.Equal(t, "login success", responseBody["message"])
				dataToken := responseBody["data"].(map[string]interface{})["token"]
				require.NotEmpty(t, dataToken)
			} else if tc.name == "failed_login_wrong_username" || tc.name == "failed_login_wrong_password" {
				require.Equal(t, 400, response.StatusCode)
				require.Equal(t, 400, int(responseBody["code"].(float64)))
				require.Equal(t, "error", responseBody["status"])
				require.Equal(t, "login failed", responseBody["message"])
				require.Equal(t, "username or password incorrect", responseBody["data"].(map[string]interface{})["errors"])
			} else {
				// Failed login validation error
				require.Equal(t, 400, response.StatusCode)
				require.Equal(t, 400, int(responseBody["code"].(float64)))
				require.Equal(t, "error", responseBody["status"])
				require.Equal(t, "login failed", responseBody["message"])
				require.NotEqual(t, 0, len((responseBody["data"].(map[string]interface{})["errors"].([]interface{}))))
			}
		})
	}
}
