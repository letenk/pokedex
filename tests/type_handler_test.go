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

func TestFindAllTypeHandlerTest(t *testing.T) {
	// Test Cases
	testCases := []struct {
		name string
		req  web.UserLoginRequest
	}{
		{
			name: "success_with_role_admin",
			req: web.UserLoginRequest{
				Username: "admin",
				Password: "password",
			},
		},
		{
			name: "failed_forbidden_with_role_user",
			req: web.UserLoginRequest{
				Username: "user",
				Password: "password",
			},
		},
		{
			name: "failed_unauthorized_as_guest",
			req:  web.UserLoginRequest{},
		},
	}

	// Test
	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Login to get token
			var token string
			// if tc.name same failed_unauthorized_as_guest dont try login
			if tc.name != "failed_unauthorized_as_guest" {
				// Data body
				dataBody := fmt.Sprintf(`{"username": "%s", "password": "%s"}`, tc.req.Username, tc.req.Password)
				requestBody := strings.NewReader(dataBody)
				request := httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/v1/login", requestBody)
				request.Header.Add("Content-Type", "application/json")
				recorder := httptest.NewRecorder()

				RouteTest.ServeHTTP(recorder, request)

				response := recorder.Result()

				// Read all response
				body, _ := io.ReadAll(response.Body)
				var responseBody map[string]interface{}
				json.Unmarshal(body, &responseBody)
				token = responseBody["data"].(map[string]interface{})["token"].(string)
			}

			// Test access categories
			request := httptest.NewRequest(http.MethodGet, "http://localhost:3000/api/v1/type", nil)
			// Added header content type
			request.Header.Add("Content-Type", "application/json")

			// if tc.name same failed_unauthorized_as_guest dont set header
			if tc.name != "failed_unauthorized_as_guest" {
				// Added token in header Authorization
				strToken := fmt.Sprintf("Bearer %s", token)
				request.Header.Add("Authorization", strToken)
			}

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

			if tc.name == "success_with_role_admin" {
				require.Equal(t, 200, response.StatusCode)
				require.Equal(t, 200, int(responseBody["code"].(float64)))
				require.Equal(t, "success", responseBody["status"])
				require.Equal(t, "list of types", responseBody["message"])
				require.NotEmpty(t, responseBody["data"])

				var contextData = responseBody["data"].([]interface{})
				require.NotEqual(t, 0, len(contextData))
				// Check each field is not empty
				for _, data := range contextData {
					list := data.(map[string]interface{})
					require.NotEmpty(t, list["id"])
					require.NotEmpty(t, list["name"])
				}
			} else if tc.name == "failed_forbidden_with_role_user" {
				require.Equal(t, 403, response.StatusCode)
				require.Equal(t, 403, int(responseBody["code"].(float64)))
				require.Equal(t, "error", responseBody["status"])
				require.Equal(t, "forbidden", responseBody["message"])
			} else {
				require.Equal(t, 401, response.StatusCode)
				require.Equal(t, 401, int(responseBody["code"].(float64)))
				require.Equal(t, "error", responseBody["status"])
				require.Equal(t, "unauthorized", responseBody["message"])
			}
		})
	}
}
