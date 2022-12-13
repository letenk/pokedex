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
	"github.com/letenk/pokedex/util"
	"github.com/stretchr/testify/require"
)

func TestCreateMonsterHandler(t *testing.T) {
	// Get data random category and type
	randCategory, randType := RandomCategoryAndType()

	// Test Cases
	testCases := []struct {
		name             string
		reqLogin         web.UserLoginRequest
		reqCreateMonster web.MonsterCreateRequest
	}{
		{
			name: "success_with_role_admin",

			reqLogin: web.UserLoginRequest{
				Username: "admin",
				Password: "password",
			},

			reqCreateMonster: web.MonsterCreateRequest{
				Name:        util.RandomString(10),
				CategoryID:  randCategory,
				Description: util.RandomString(20),
				Length:      54.3,
				Weight:      uint16(util.RandomInt(50, 500)),
				Hp:          uint16(util.RandomInt(50, 500)),
				Attack:      uint16(util.RandomInt(50, 500)),
				Defends:     uint16(util.RandomInt(50, 500)),
				Speed:       uint16(util.RandomInt(50, 500)),
				Image:       util.RandomString(10),
				TypeID:      []string{randType, randType, randType},
			},
		},
		{
			name: "failed_forbidden_with_role_user",
			reqLogin: web.UserLoginRequest{
				Username: "user",
				Password: "password",
			},

			reqCreateMonster: web.MonsterCreateRequest{
				Name:        util.RandomString(10),
				CategoryID:  randCategory,
				Description: util.RandomString(20),
				Length:      54.3,
				Weight:      uint16(util.RandomInt(50, 500)),
				Hp:          uint16(util.RandomInt(50, 500)),
				Attack:      uint16(util.RandomInt(50, 500)),
				Defends:     uint16(util.RandomInt(50, 500)),
				Speed:       uint16(util.RandomInt(50, 500)),
				Image:       util.RandomString(10),
				TypeID:      []string{randType, randType, randType},
			},
		},
		{
			name:     "failed_unauthorized_as_guest",
			reqLogin: web.UserLoginRequest{},

			reqCreateMonster: web.MonsterCreateRequest{
				Name:        util.RandomString(10),
				CategoryID:  randCategory,
				Description: util.RandomString(20),
				Length:      54.3,
				Weight:      uint16(util.RandomInt(50, 500)),
				Hp:          uint16(util.RandomInt(50, 500)),
				Attack:      uint16(util.RandomInt(50, 500)),
				Defends:     uint16(util.RandomInt(50, 500)),
				Speed:       uint16(util.RandomInt(50, 500)),
				Image:       util.RandomString(10),
				TypeID:      []string{randType, randType, randType},
			},
		},
		{
			name: "failed_validation_error",
			reqLogin: web.UserLoginRequest{
				Username: "admin",
				Password: "password",
			},
			reqCreateMonster: web.MonsterCreateRequest{},
		},
	}

	// Test
	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			// Login to get token
			var token string
			// if tc.name same failed_unauthorized_as_guest dont try login
			if tc.name != "failed_unauthorized_as_guest" {
				token = GetToken(tc.reqLogin)
			}

			// Create new monster
			// Data body
			var dataBody string
			if tc.name != "failed_validation_error" {
				dataBody = fmt.Sprintf(`{"name": "%s", "category_id": "%s", "description": "%s", "length": %.1f, "weight": %d, "hp": %d, "attack": %d, "defends": %d, "speed": %d, "image": "%s", "type_id": ["%s", "%s"]}`, tc.reqCreateMonster.Name, tc.reqCreateMonster.CategoryID, tc.reqCreateMonster.Description, tc.reqCreateMonster.Length, tc.reqCreateMonster.Weight, tc.reqCreateMonster.Hp, tc.reqCreateMonster.Attack, tc.reqCreateMonster.Defends, tc.reqCreateMonster.Speed, tc.reqCreateMonster.Image, randType, randType)
			} else {
				var emptySlice []string
				dataBody = fmt.Sprintf(`{"name": "%s", "category_id": "%s", "description": "%s", "length": %.1f, "weight": %d, "hp": %d, "attack": %d, "defends": %d, "speed": %d, "image": "%s", "type_id": %v}`, "", tc.reqCreateMonster.CategoryID, tc.reqCreateMonster.Description, tc.reqCreateMonster.Length, tc.reqCreateMonster.Weight, tc.reqCreateMonster.Hp, tc.reqCreateMonster.Attack, tc.reqCreateMonster.Defends, tc.reqCreateMonster.Speed, tc.reqCreateMonster.Image, emptySlice)
			}

			// New reader
			requestBody := strings.NewReader(dataBody)

			// Test access categories
			request := httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/v1/monster", requestBody)
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
				require.Equal(t, 201, response.StatusCode)
				require.Equal(t, 201, int(responseBody["code"].(float64)))
				require.Equal(t, "success", responseBody["status"])
				require.Equal(t, "Monster has been created", responseBody["message"])
			} else if tc.name == "failed_forbidden_with_role_user" {
				require.Equal(t, 403, response.StatusCode)
				require.Equal(t, 403, int(responseBody["code"].(float64)))
				require.Equal(t, "error", responseBody["status"])
				require.Equal(t, "forbidden", responseBody["message"])
			} else if tc.name == "failed_unauthorized_as_guest" {
				require.Equal(t, 401, response.StatusCode)
				require.Equal(t, 401, int(responseBody["code"].(float64)))
				require.Equal(t, "error", responseBody["status"])
				require.Equal(t, "unauthorized", responseBody["message"])
			} else {
				// If validation error
				require.Equal(t, 400, response.StatusCode)
				require.Equal(t, 400, int(responseBody["code"].(float64)))
				require.Equal(t, "error", responseBody["status"])
				require.Equal(t, "create monster failed", responseBody["message"])
				require.NotEqual(t, 0, len((responseBody["data"].(map[string]interface{})["errors"].([]interface{}))))
			}
		})

	}
}
