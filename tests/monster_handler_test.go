package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
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
				// Image:       util.RandomString(10),
				TypeID: []string{randType, randType, randType},
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
				// Image:       util.RandomString(10),
				TypeID: []string{randType, randType, randType},
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
				// Image:       util.RandomString(10),
				TypeID: []string{randType, randType, randType},
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
			bodyRequest := new(bytes.Buffer)
			writer := multipart.NewWriter(bodyRequest)

			if tc.name != "failed_validation_error" {
				writer.WriteField("name", tc.reqCreateMonster.Name)
				writer.WriteField("category_id", tc.reqCreateMonster.CategoryID)
				writer.WriteField("description", tc.reqCreateMonster.Description)
				writer.WriteField("length", strconv.FormatFloat(float64(tc.reqCreateMonster.Length), 'f', 6, 64))
				writer.WriteField("weight", strconv.Itoa(int(tc.reqCreateMonster.Weight)))
				writer.WriteField("hp", strconv.Itoa(int(tc.reqCreateMonster.Hp)))
				writer.WriteField("attack", strconv.Itoa(int(tc.reqCreateMonster.Attack)))
				writer.WriteField("defends", strconv.Itoa(int(tc.reqCreateMonster.Defends)))
				writer.WriteField("speed", strconv.Itoa(int(tc.reqCreateMonster.Speed)))
				writer.WriteField("type_id", randType)

				// Read file from local
				file, _ := os.Open("file_sample/image.png")
				defer file.Close()

				part, err := writer.CreateFormFile("image", "image.png")

				if err != nil {
					log.Fatal(err)
				}

				_, err = io.Copy(part, file)
				if err != nil {
					log.Fatal(err)
				}

				writer.Close()
			} else {
				// Validation error field is empty
				writer.WriteField("name", "")
				writer.WriteField("category_id", "")
				writer.WriteField("description", "")
				writer.WriteField("length", "")
				writer.WriteField("weight", "")
				writer.WriteField("hp", "")
				writer.WriteField("attack", "")
				writer.WriteField("defends", "")
				writer.WriteField("speed", "")
				writer.Close()
			}

			// Test access categories
			request := httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/v1/monster", bodyRequest)
			// Added header content type
			request.Header.Set("Content-Type", writer.FormDataContentType())

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

func TestFindAllMonsterHandler(t *testing.T) {
	// Create random monsters
	newMonster, randTypes := RandomCreateMonster(t)

	testCases := []struct {
		name           string
		queryParameter web.MonsterQueryRequest
	}{
		{
			name:           "find_all_monsters_without_query_parameter",
			queryParameter: web.MonsterQueryRequest{},
		},
		{
			name: "find_all_monsters_with_query_parameter_name",
			queryParameter: web.MonsterQueryRequest{
				Name: newMonster.Name,
			},
		},
		{
			name: "find_all_monsters_with_query_parameter_catched_false",
			queryParameter: web.MonsterQueryRequest{
				Catched: "false",
			},
		},
		{
			name: "find_all_monsters_with_query_parameter_sort_by_name",
			queryParameter: web.MonsterQueryRequest{
				Sort: "name",
			},
		},
		{
			name: "find_all_monsters_with_query_parameter_sort_by_id",
			queryParameter: web.MonsterQueryRequest{
				Sort: "id",
			},
		},
		{
			name: "find_all_monsters_with_query_parameter_order_by_desc",
			queryParameter: web.MonsterQueryRequest{
				Sort:  "name",
				Order: "desc",
			},
		},
		{
			name: "find_all_monsters_with_query_parameter_types",
			queryParameter: web.MonsterQueryRequest{
				Types: randTypes,
			},
		},
		{
			name: "find_all_monsters_with_query_parameter_types",
			queryParameter: web.MonsterQueryRequest{
				Types: randTypes,
			},
		},
	}

	// Test
	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Test access categories
			var urlTarget string
			if tc.name == "find_all_monsters_with_query_parameter_name" {
				urlTarget = fmt.Sprintf("http://localhost:3000/api/v1/monster?name=%s", tc.queryParameter.Name)
			} else if tc.name == "find_all_monsters_with_query_parameter_catched_false" {
				urlTarget = fmt.Sprintf("http://localhost:3000/api/v1/monster?catched=%s", tc.queryParameter.Catched)
			} else if tc.name == "find_all_monsters_with_query_parameter_sort_by_name" || tc.name == "find_all_monsters_with_query_parameter_sort_by_id" {
				urlTarget = fmt.Sprintf("http://localhost:3000/api/v1/monster?sort=%s", tc.queryParameter.Sort)
			} else if tc.name == "find_all_monsters_with_query_parameter_order_by_asc" || tc.name == "find_all_monsters_with_query_parameter_order_by_desc" {
				urlTarget = fmt.Sprintf("http://localhost:3000/api/v1/monster?sort=%s&order=%s", tc.queryParameter.Sort, tc.queryParameter.Order)
			} else if tc.name == "find_all_monsters_with_query_parameter_types" {
				urlTarget = fmt.Sprintf("http://localhost:3000/api/v1/monster?types=%s&types=%s&types=%s", tc.queryParameter.Types[0], tc.queryParameter.Types[1], tc.queryParameter.Types[2])
			} else {
				urlTarget = fmt.Sprint("http://localhost:3000/api/v1/monster")
			}

			request := httptest.NewRequest(http.MethodGet, urlTarget, nil)
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

			require.Equal(t, 200, response.StatusCode)
			require.Equal(t, 200, int(responseBody["code"].(float64)))
			require.Equal(t, "success", responseBody["status"])
			require.Equal(t, "List of monsters", responseBody["message"])
			require.NotEmpty(t, responseBody["data"])

			var contextData = responseBody["data"].([]any)
			require.NotEqual(t, 0, len(contextData))

			// Check each field is not empty
			for _, data := range contextData {
				list := data.(map[string]any)
				require.NotEmpty(t, list["id"])
				require.NotEmpty(t, list["name"])
				require.NotEmpty(t, list["category_name"])
				require.NotEmpty(t, list["image"])

				strCatched := strconv.FormatBool(list["catched"].(bool))
				require.NotEmpty(t, strCatched)

				require.NotEqual(t, 0, len(list["types"].([]any)))
				for _, ty := range list["types"].([]any) {
					listType := ty.(map[string]any)
					require.NotEmpty(t, listType["name"])
				}
			}
		})
	}
}
