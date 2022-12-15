package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/png"
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

func CreateImage() image.Image {
	width := 200
	height := 100

	upLeft := image.Point{0, 0}
	lowRight := image.Point{width, height}

	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	// Colors are defined by Red, Green, Blue, Alpha uint8 values.
	cyan := color.RGBA{100, 200, 200, 0xff}

	// Set color for each pixel.
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			switch {
			case x < width/2 && y < height/2: // upper left quadrant
				img.Set(x, y, cyan)
			case x >= width/2 && y >= height/2: // lower right quadrant
				img.Set(x, y, color.White)
			default:
				// Use zero value.
			}
		}
	}

	// Encode as PNG.
	f, _ := os.Create("image.png")
	png.Encode(f, img)

	return img
}

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
				file, _ := os.Open("image.png")
				defer file.Close()

				part, err := writer.CreateFormFile("image", "image.png")

				if err != nil {
					log.Fatal(err)
				}

				_ = CreateImage()
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
