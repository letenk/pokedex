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
	"strings"
	"testing"

	"github.com/letenk/pokedex/models/web"
	"github.com/letenk/pokedex/util"
	"github.com/stretchr/testify/require"
)

func TestCreateMonsterHandler(t *testing.T) {
	var randTypes []string
	var randCategories []string
	// Get data random category and type
	for i := 0; i < 3; i++ {
		randCategory, randType := RandomCategoryAndType()
		randTypes = append(randTypes, randType)
		randCategories = append(randCategories, randCategory)
	}

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
				CategoryID:  randCategories[0],
				Description: util.RandomString(20),
				Length:      54.3,
				Weight:      uint16(util.RandomInt(50, 500)),
				Hp:          uint16(util.RandomInt(50, 500)),
				Attack:      uint16(util.RandomInt(50, 500)),
				Defends:     uint16(util.RandomInt(50, 500)),
				Speed:       uint16(util.RandomInt(50, 500)),
				TypeID:      randTypes,
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
				CategoryID:  randCategories[0],
				Description: util.RandomString(20),
				Length:      54.3,
				Weight:      uint16(util.RandomInt(50, 500)),
				Hp:          uint16(util.RandomInt(50, 500)),
				Attack:      uint16(util.RandomInt(50, 500)),
				Defends:     uint16(util.RandomInt(50, 500)),
				Speed:       uint16(util.RandomInt(50, 500)),
				TypeID:      randTypes,
			},
		},
		{
			name:     "failed_unauthorized_as_guest",
			reqLogin: web.UserLoginRequest{},

			reqCreateMonster: web.MonsterCreateRequest{
				Name:        util.RandomString(10),
				CategoryID:  randCategories[0],
				Description: util.RandomString(20),
				Length:      54.3,
				Weight:      uint16(util.RandomInt(50, 500)),
				Hp:          uint16(util.RandomInt(50, 500)),
				Attack:      uint16(util.RandomInt(50, 500)),
				Defends:     uint16(util.RandomInt(50, 500)),
				Speed:       uint16(util.RandomInt(50, 500)),
				TypeID:      randTypes,
			},
		},
		{
			name: "failed_invalid_category_id",
			reqLogin: web.UserLoginRequest{
				Username: "admin",
				Password: "password",
			},

			reqCreateMonster: web.MonsterCreateRequest{
				Name:        util.RandomString(10),
				CategoryID:  "4562482c-7acd-4daf-901f-d95c7a7afd65",
				Description: util.RandomString(20),
				Length:      54.3,
				Weight:      uint16(util.RandomInt(50, 500)),
				Hp:          uint16(util.RandomInt(50, 500)),
				Attack:      uint16(util.RandomInt(50, 500)),
				Defends:     uint16(util.RandomInt(50, 500)),
				Speed:       uint16(util.RandomInt(50, 500)),
				TypeID:      randTypes,
			},
		},
		{
			name: "failed_invalid_type_id",
			reqLogin: web.UserLoginRequest{
				Username: "admin",
				Password: "password",
			},

			reqCreateMonster: web.MonsterCreateRequest{
				Name:        util.RandomString(10),
				CategoryID:  randCategories[0],
				Description: util.RandomString(20),
				Length:      54.3,
				Weight:      uint16(util.RandomInt(50, 500)),
				Hp:          uint16(util.RandomInt(50, 500)),
				Attack:      uint16(util.RandomInt(50, 500)),
				Defends:     uint16(util.RandomInt(50, 500)),
				Speed:       uint16(util.RandomInt(50, 500)),
				TypeID:      []string{"558160ef-e8f5-4951-b5f4-feeb0815b511", "d5a8d4bb-eb0a-44a4-ae46-eb2af2b2002c"},
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
				writer.WriteField("type_id", tc.reqCreateMonster.TypeID[0])
				writer.WriteField("type_id", tc.reqCreateMonster.TypeID[1])

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
			} else if tc.name == "failed_invalid_category_id" {
				require.Equal(t, 400, response.StatusCode)
				require.Equal(t, 400, int(responseBody["code"].(float64)))
				require.Equal(t, "error", responseBody["status"])
				require.Equal(t, "create monster failed", responseBody["message"])
				require.Equal(t, "invalid category id or type id, please check valid id in each of their list", responseBody["data"].(map[string]interface{})["errors"])
			} else if tc.name == "failed_invalid_type_id" {
				require.Equal(t, 400, response.StatusCode)
				require.Equal(t, 400, int(responseBody["code"].(float64)))
				require.Equal(t, "error", responseBody["status"])
				require.Equal(t, "create monster failed", responseBody["message"])
				require.Equal(t, "invalid category id or type id, please check valid id in each of their list", responseBody["data"].(map[string]interface{})["errors"])
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
				require.NotEmpty(t, list["image_url"])

				strCatched := strconv.FormatBool(list["catched"].(bool))
				require.NotEmpty(t, strCatched)

				// require.NotEqual(t, 0, len(list["types"].([]any)))
				for _, ty := range list["types"].([]any) {
					listType := ty.(map[string]any)
					require.NotEmpty(t, listType["name"])
				}
			}
		})
	}
}

func TestFindByIDMonsterHandler(t *testing.T) {
	// Create random monsters
	newMonster, _ := RandomCreateMonster(t)

	testCases := []struct {
		name      string
		idMonster string
	}{
		{
			name:      "find_by_id_monster_success",
			idMonster: newMonster.ID,
		},
		{
			name:      "find_by_id_monster_failed",
			idMonster: "368bd987-dec6-4405-a036-bc1232db21b2",
		},
	}

	// Test
	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			request := httptest.NewRequest(http.MethodGet, "http://localhost:3000/api/v1/monster/"+tc.idMonster, nil)
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

			if tc.name == "find_by_id_monster_success" {
				require.Equal(t, 200, response.StatusCode)
				require.Equal(t, 200, int(responseBody["code"].(float64)))
				require.Equal(t, "success", responseBody["status"])
				require.Equal(t, "profile detail of monsters", responseBody["message"])
				require.NotEmpty(t, responseBody["data"])

				var contextData = responseBody["data"].(map[string]any)

				require.NotEmpty(t, contextData["category_name"])

				require.Equal(t, newMonster.ID, contextData["id"])
				require.Equal(t, newMonster.Name, contextData["name"])
				require.Equal(t, newMonster.Description, contextData["description"])
				require.Equal(t, newMonster.Length, float32(contextData["length"].(float64)))
				require.Equal(t, newMonster.Weight, uint16(contextData["weight"].(float64)))
				require.Equal(t, newMonster.Hp, uint16(contextData["hp"].(float64)))
				require.Equal(t, newMonster.Attack, uint16(contextData["attack"].(float64)))
				require.Equal(t, newMonster.Defends, uint16(contextData["defends"].(float64)))
				require.Equal(t, newMonster.Speed, uint16(contextData["speed"].(float64)))
				require.Equal(t, newMonster.Catched, contextData["catched"])

				require.NotEqual(t, 0, len(contextData["types"].([]any)))
				for _, ty := range contextData["types"].([]any) {
					listType := ty.(map[string]any)
					require.NotEmpty(t, listType["name"])
				}
			} else {
				require.Equal(t, 400, response.StatusCode)
				require.Equal(t, 400, int(responseBody["code"].(float64)))
				require.Equal(t, "error", responseBody["status"])
				require.Equal(t, "bad request", responseBody["message"])

				errMessage := fmt.Sprintf("monster with id %s not found", tc.idMonster)
				require.Equal(t, errMessage, responseBody["data"].(map[string]interface{})["errors"])
			}

		})
	}
}

func TestUpdateMonsterHandler(t *testing.T) {
	newMonster := RandomCreateMonsterUsecase(t)
	var randTypes []string
	var randCategories []string
	// Get data random category and type
	for i := 0; i < 3; i++ {
		randCategory, randType := RandomCategoryAndType()
		randTypes = append(randTypes, randType)
		randCategories = append(randCategories, randCategory)
	}

	// Test Cases
	testCases := []struct {
		name             string
		reqLogin         web.UserLoginRequest
		reqCreateMonster web.MonsterUpdateRequest
	}{
		{
			name: "update_success_with_field_empty",
			reqLogin: web.UserLoginRequest{
				Username: "admin",
				Password: "password",
			},
			reqCreateMonster: web.MonsterUpdateRequest{
				Name:        "",
				CategoryID:  "",
				Description: "",
				Length:      "",
				Weight:      "",
				Hp:          "",
				Attack:      "",
				Defends:     "",
				Speed:       "",
				Catched:     "",
			},
		},
		{
			name: "success_with_role_admin_with_image",

			reqLogin: web.UserLoginRequest{
				Username: "admin",
				Password: "password",
			},

			reqCreateMonster: web.MonsterUpdateRequest{
				Name:        util.RandomString(10),
				CategoryID:  randCategories[0],
				Description: util.RandomString(20),
				Length:      "54.2",
				Weight:      strconv.Itoa(int(util.RandomInt(50, 500))),
				Hp:          strconv.Itoa(int(util.RandomInt(50, 500))),
				Attack:      strconv.Itoa(int(util.RandomInt(50, 500))),
				Defends:     strconv.Itoa(int(util.RandomInt(50, 500))),
				Speed:       strconv.Itoa(int(util.RandomInt(50, 500))),
				Catched:     "true",
				TypeID:      randTypes,
			},
		},
		{
			name: "failed_forbidden_with_role_user",
			reqLogin: web.UserLoginRequest{
				Username: "user",
				Password: "password",
			},

			reqCreateMonster: web.MonsterUpdateRequest{
				Name:        util.RandomString(10),
				CategoryID:  randCategories[0],
				Description: util.RandomString(20),
				Length:      "54.3",
				Weight:      strconv.Itoa(int(util.RandomInt(50, 500))),
				Hp:          strconv.Itoa(int(util.RandomInt(50, 500))),
				Attack:      strconv.Itoa(int(util.RandomInt(50, 500))),
				Defends:     strconv.Itoa(int(util.RandomInt(50, 500))),
				Speed:       strconv.Itoa(int(util.RandomInt(50, 500))),
				Catched:     "true",
				TypeID:      randTypes,
			},
		},
		{
			name:     "failed_unauthorized_as_guest",
			reqLogin: web.UserLoginRequest{},

			reqCreateMonster: web.MonsterUpdateRequest{
				Name:        util.RandomString(10),
				CategoryID:  randCategories[0],
				Description: util.RandomString(20),
				Length:      "54.3",
				Weight:      strconv.Itoa(int(util.RandomInt(50, 500))),
				Hp:          strconv.Itoa(int(util.RandomInt(50, 500))),
				Attack:      strconv.Itoa(int(util.RandomInt(50, 500))),
				Defends:     strconv.Itoa(int(util.RandomInt(50, 500))),
				Speed:       strconv.Itoa(int(util.RandomInt(50, 500))),
				Catched:     "true",
				TypeID:      randTypes,
			},
		},
		{
			name: "failed_invalid_category_id",
			reqLogin: web.UserLoginRequest{
				Username: "admin",
				Password: "password",
			},

			reqCreateMonster: web.MonsterUpdateRequest{
				Name:        util.RandomString(10),
				CategoryID:  "4562482c-7acd-4daf-901f-d95c7a7afd65",
				Description: util.RandomString(20),
				Length:      "54.3",
				Weight:      strconv.Itoa(int(util.RandomInt(50, 500))),
				Hp:          strconv.Itoa(int(util.RandomInt(50, 500))),
				Attack:      strconv.Itoa(int(util.RandomInt(50, 500))),
				Defends:     strconv.Itoa(int(util.RandomInt(50, 500))),
				Speed:       strconv.Itoa(int(util.RandomInt(50, 500))),
				Catched:     "true",
				TypeID:      randTypes,
			},
		},
		{
			name: "failed_invalid_type_id",
			reqLogin: web.UserLoginRequest{
				Username: "admin",
				Password: "password",
			},

			reqCreateMonster: web.MonsterUpdateRequest{
				Name:        util.RandomString(10),
				CategoryID:  "4562482c-7acd-4daf-901f-d95c7a7afd65",
				Description: util.RandomString(20),
				Length:      "54.3",
				Weight:      strconv.Itoa(int(util.RandomInt(50, 500))),
				Hp:          strconv.Itoa(int(util.RandomInt(50, 500))),
				Attack:      strconv.Itoa(int(util.RandomInt(50, 500))),
				Defends:     strconv.Itoa(int(util.RandomInt(50, 500))),
				Speed:       strconv.Itoa(int(util.RandomInt(50, 500))),
				Catched:     "true",
				TypeID:      []string{"558160ef-e8f5-4951-b5f4-feeb0815b511", "d5a8d4bb-eb0a-44a4-ae46-eb2af2b2002c"},
			},
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

			if tc.name == "update_success_with_field_empty" {
				writer.WriteField("name", tc.reqCreateMonster.Name)
				writer.WriteField("category_id", tc.reqCreateMonster.CategoryID)
				writer.WriteField("description", tc.reqCreateMonster.Description)
				writer.WriteField("length", tc.reqCreateMonster.Length)
				writer.WriteField("weight", tc.reqCreateMonster.Weight)
				writer.WriteField("hp", tc.reqCreateMonster.Hp)
				writer.WriteField("attack", tc.reqCreateMonster.Attack)
				writer.WriteField("defends", tc.reqCreateMonster.Defends)
				writer.WriteField("speed", tc.reqCreateMonster.Speed)
				writer.WriteField("catched", tc.reqCreateMonster.Catched)

				writer.Close()
			} else if tc.name != "success_with_role_admin_with_image" {
				// if tc.name is "failed_forbidden_with_role_user", "failed_unauthorized_as_guest", "failed_invalid_category_id", "failed_invalid_type_id",
				writer.WriteField("name", tc.reqCreateMonster.Name)
				writer.WriteField("category_id", tc.reqCreateMonster.CategoryID)
				writer.WriteField("description", tc.reqCreateMonster.Description)
				writer.WriteField("length", tc.reqCreateMonster.Length)
				writer.WriteField("weight", tc.reqCreateMonster.Weight)
				writer.WriteField("hp", tc.reqCreateMonster.Hp)
				writer.WriteField("attack", tc.reqCreateMonster.Attack)
				writer.WriteField("defends", tc.reqCreateMonster.Defends)
				writer.WriteField("speed", tc.reqCreateMonster.Speed)
				writer.WriteField("catched", tc.reqCreateMonster.Catched)
				writer.WriteField("type_id", tc.reqCreateMonster.TypeID[0])
				writer.WriteField("type_id", tc.reqCreateMonster.TypeID[1])

				writer.Close()
			} else {
				// if tc.name == "success_with_role_admin_with_image"
				writer.WriteField("name", tc.reqCreateMonster.Name)
				writer.WriteField("category_id", tc.reqCreateMonster.CategoryID)
				writer.WriteField("description", tc.reqCreateMonster.Description)
				writer.WriteField("length", tc.reqCreateMonster.Length)
				writer.WriteField("weight", tc.reqCreateMonster.Weight)
				writer.WriteField("hp", tc.reqCreateMonster.Hp)
				writer.WriteField("attack", tc.reqCreateMonster.Attack)
				writer.WriteField("defends", tc.reqCreateMonster.Defends)
				writer.WriteField("speed", tc.reqCreateMonster.Speed)
				writer.WriteField("catched", tc.reqCreateMonster.Catched)
				writer.WriteField("type_id", tc.reqCreateMonster.TypeID[0])
				writer.WriteField("type_id", tc.reqCreateMonster.TypeID[1])

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
			}

			// Test access categories
			request := httptest.NewRequest(http.MethodPatch, "http://localhost:3000/api/v1/monster/"+newMonster.ID, bodyRequest)
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

			if tc.name == "failed_forbidden_with_role_user" {
				require.Equal(t, 403, response.StatusCode)
				require.Equal(t, 403, int(responseBody["code"].(float64)))
				require.Equal(t, "error", responseBody["status"])
				require.Equal(t, "forbidden", responseBody["message"])
			} else if tc.name == "failed_unauthorized_as_guest" {
				require.Equal(t, 401, response.StatusCode)
				require.Equal(t, 401, int(responseBody["code"].(float64)))
				require.Equal(t, "error", responseBody["status"])
				require.Equal(t, "unauthorized", responseBody["message"])
			} else if tc.name == "failed_invalid_category_id" || tc.name == "failed_invalid_type_id" {
				require.Equal(t, 400, response.StatusCode)
				require.Equal(t, 400, int(responseBody["code"].(float64)))
				require.Equal(t, "error", responseBody["status"])
				require.Equal(t, "update monster failed", responseBody["message"])
				require.Equal(t, "invalid category id or type id, please check valid id in each of their list", responseBody["data"].(map[string]interface{})["errors"])
			} else if tc.name == "success_with_role_admin_with_image" {
				require.Equal(t, 200, response.StatusCode)
				require.Equal(t, 200, int(responseBody["code"].(float64)))
				require.Equal(t, "success", responseBody["status"])
				require.Equal(t, "Update monster success", responseBody["message"])
				require.NotEmpty(t, responseBody["data"])

				var contextData = responseBody["data"].(map[string]any)

				require.Equal(t, newMonster.ID, contextData["id"])

				require.NotEqual(t, newMonster.Name, contextData["name"])
				require.NotEqual(t, newMonster.Description, contextData["description"])
				require.NotEqual(t, newMonster.Length, float32(contextData["length"].(float64)))
				require.NotEqual(t, newMonster.Weight, uint16(contextData["weight"].(float64)))
				require.NotEqual(t, newMonster.Hp, uint16(contextData["hp"].(float64)))
				require.NotEqual(t, newMonster.Attack, uint16(contextData["attack"].(float64)))
				require.NotEqual(t, newMonster.Defends, uint16(contextData["defends"].(float64)))
				require.NotEqual(t, newMonster.Speed, uint16(contextData["speed"].(float64)))
				require.NotEqual(t, newMonster.Catched, contextData["catched"].(bool))

				require.NotEqual(t, newMonster.ImageURL, contextData["image_url"])

				require.NotEmpty(t, contextData["category_name"])

				require.NotEqual(t, 0, len(contextData["types"].([]any)))
				for _, ty := range contextData["types"].([]any) {
					listType := ty.(map[string]any)
					require.NotEmpty(t, listType["name"])
				}
			} else {
				// if tc.name == "update_success_with_field_empty"
				require.Equal(t, 200, response.StatusCode)
				require.Equal(t, 200, int(responseBody["code"].(float64)))
				require.Equal(t, "success", responseBody["status"])
				require.Equal(t, "Update monster success", responseBody["message"])
				require.NotEmpty(t, responseBody["data"])

				var contextData = responseBody["data"].(map[string]any)

				require.Equal(t, newMonster.ID, contextData["id"])
				require.NotEmpty(t, contextData["category_name"])

				require.Equal(t, newMonster.Name, contextData["name"])
				require.Equal(t, newMonster.Description, contextData["description"])
				require.Equal(t, newMonster.Length, float32(contextData["length"].(float64)))
				require.Equal(t, newMonster.Weight, uint16(contextData["weight"].(float64)))
				require.Equal(t, newMonster.Hp, uint16(contextData["hp"].(float64)))
				require.Equal(t, newMonster.Attack, uint16(contextData["attack"].(float64)))
				require.Equal(t, newMonster.Defends, uint16(contextData["defends"].(float64)))
				require.Equal(t, newMonster.Speed, uint16(contextData["speed"].(float64)))
				require.Equal(t, newMonster.Catched, contextData["catched"])
				require.Equal(t, newMonster.ImageURL, contextData["image_url"])

				require.NotEqual(t, 0, len(contextData["types"].([]any)))
				for _, ty := range contextData["types"].([]any) {
					listType := ty.(map[string]any)
					require.NotEmpty(t, listType["name"])
				}
			}
		})
	}
}

func TestUpdateMarkMonsterCapturedHandler(t *testing.T) {
	newMonster, _ := RandomCreateMonster(t)

	// Test Cases
	testCases := []struct {
		name      string
		idMonster string
		reqLogin  web.UserLoginRequest
		reqUpdate web.MonsterUpdateRequestMonsterCapture
	}{
		{
			name:      "update_mark_monster_captured_monster_success",
			idMonster: newMonster.ID,
			reqLogin: web.UserLoginRequest{
				Username: "user",
				Password: "password",
			},
			reqUpdate: web.MonsterUpdateRequestMonsterCapture{
				Catched: true,
			},
		},
		{
			name:      "update_mark_monster_captured_monster_failed_monster_not_found",
			idMonster: "368bd987-dec6-4405-a036-bc1232db21b2",
			reqLogin: web.UserLoginRequest{
				Username: "user",
				Password: "password",
			},
			reqUpdate: web.MonsterUpdateRequestMonsterCapture{
				Catched: true,
			},
		},
		{
			name:      "failed_forbidden_with_role_admin",
			idMonster: newMonster.ID,
			reqLogin: web.UserLoginRequest{
				Username: "admin",
				Password: "password",
			},
			reqUpdate: web.MonsterUpdateRequestMonsterCapture{
				Catched: true,
			},
		},
		{
			name:      "failed_unauthorized_as_guest",
			idMonster: newMonster.ID,
			reqLogin:  web.UserLoginRequest{},
			reqUpdate: web.MonsterUpdateRequestMonsterCapture{
				Catched: true,
			},
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
				token = GetToken(tc.reqLogin)
			}

			// Request body
			dataBody := fmt.Sprintf(`{"catched": "%v"}`, tc.reqUpdate.Catched)
			requestBody := strings.NewReader(dataBody)

			// Url
			url := fmt.Sprintf("http://localhost:3000/api/v1/monster/%s/captured", tc.idMonster)
			// Test access categories
			request := httptest.NewRequest(http.MethodPatch, url, requestBody)

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

			if tc.name == "failed_forbidden_with_role_admin" {
				require.Equal(t, 403, response.StatusCode)
				require.Equal(t, 403, int(responseBody["code"].(float64)))
				require.Equal(t, "error", responseBody["status"])
				require.Equal(t, "forbidden", responseBody["message"])
			} else if tc.name == "failed_unauthorized_as_guest" {
				require.Equal(t, 401, response.StatusCode)
				require.Equal(t, 401, int(responseBody["code"].(float64)))
				require.Equal(t, "error", responseBody["status"])
				require.Equal(t, "unauthorized", responseBody["message"])
			} else if tc.name == "update_mark_monster_captured_monster_failed_monster_not_found" {
				require.Equal(t, 400, response.StatusCode)
				require.Equal(t, 400, int(responseBody["code"].(float64)))
				require.Equal(t, "error", responseBody["status"])
				require.Equal(t, "update monster captured failed", responseBody["message"])

				errMessage := fmt.Sprintf("monster with id %s not found", tc.idMonster)
				require.Equal(t, errMessage, responseBody["data"].(map[string]interface{})["errors"])
			} else {
				// if tc.name == "update_mark_monster_captured_monster_success",
				require.Equal(t, 200, response.StatusCode)
				require.Equal(t, 200, int(responseBody["code"].(float64)))
				require.Equal(t, "success", responseBody["status"])
				msgSuccess := fmt.Sprintf("monster with id %s updated", tc.idMonster)
				require.Equal(t, msgSuccess, responseBody["message"])
			}
		})
	}
}

func TestDeleteMonsterHandler(t *testing.T) {
	newMonster := RandomCreateMonsterUsecase(t)

	testCases := []struct {
		name      string
		idMonster string
		reqLogin  web.UserLoginRequest
	}{
		{
			name:      "delete_monster_success",
			idMonster: newMonster.ID,
			reqLogin: web.UserLoginRequest{
				Username: "admin",
				Password: "password",
			},
		},
		{
			name:      "delete_monster_failed_monster_not_found",
			idMonster: "368bd987-dec6-4405-a036-bc1232db21b2",
			reqLogin: web.UserLoginRequest{
				Username: "admin",
				Password: "password",
			},
		},
		{
			name:      "failed_forbidden_with_role_user",
			idMonster: newMonster.ID,
			reqLogin: web.UserLoginRequest{
				Username: "user",
				Password: "password",
			},
		},
		{
			name:      "failed_unauthorized_as_guest",
			idMonster: newMonster.ID,
			reqLogin:  web.UserLoginRequest{},
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
				token = GetToken(tc.reqLogin)
			}

			// Url
			url := fmt.Sprintf("http://localhost:3000/api/v1/monster/%s", tc.idMonster)
			// Test access categories
			request := httptest.NewRequest(http.MethodDelete, url, nil)

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

			if tc.name == "failed_forbidden_with_role_user" {
				require.Equal(t, 403, response.StatusCode)
				require.Equal(t, 403, int(responseBody["code"].(float64)))
				require.Equal(t, "error", responseBody["status"])
				require.Equal(t, "forbidden", responseBody["message"])
			} else if tc.name == "failed_unauthorized_as_guest" {
				require.Equal(t, 401, response.StatusCode)
				require.Equal(t, 401, int(responseBody["code"].(float64)))
				require.Equal(t, "error", responseBody["status"])
				require.Equal(t, "unauthorized", responseBody["message"])
			} else if tc.name == "delete_monster_failed_monster_not_found" {
				require.Equal(t, 400, response.StatusCode)
				require.Equal(t, 400, int(responseBody["code"].(float64)))
				require.Equal(t, "error", responseBody["status"])
				require.Equal(t, "delete monster failed", responseBody["message"])

				errMessage := fmt.Sprintf("monster with id %s not found", tc.idMonster)
				require.Equal(t, errMessage, responseBody["data"].(map[string]interface{})["errors"])
			} else {
				// if tc.name == "delete_monster_success"",
				require.Equal(t, 200, response.StatusCode)
				require.Equal(t, 200, int(responseBody["code"].(float64)))
				require.Equal(t, "success", responseBody["status"])
				require.Equal(t, "monster deleted", responseBody["message"])
			}
		})
	}
}
