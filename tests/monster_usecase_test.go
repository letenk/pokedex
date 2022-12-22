package tests

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/letenk/pokedex/models/domain"
	"github.com/letenk/pokedex/models/web"
	"github.com/letenk/pokedex/repository"
	"github.com/letenk/pokedex/usecase"
	"github.com/letenk/pokedex/util"
	"github.com/stretchr/testify/require"
)

func RandomCreateMonsterUsecase(t *testing.T) domain.Monster {
	t.Parallel()

	// Read file from local
	file, _ := os.Open("file_sample/image.png")
	defer file.Close()

	// File name formate
	now := time.Now()
	nowRFC3339 := now.Format(time.RFC3339)
	fileName := fmt.Sprintf(`%s_%v_%s`, "usecase_create_test", nowRFC3339, "image.png")

	repositoryMonster := repository.NewMonsterRespository(ConnTest)
	usecaseMonster := usecase.NewUsecaseMonster(repositoryMonster)

	var randTypes []string
	var randCategories []string
	// Get data random category and type
	for i := 0; i < 3; i++ {
		randCategory, randType := RandomCategoryAndType()
		randTypes = append(randTypes, randType)
		randCategories = append(randCategories, randCategory)
	}

	// Test cases
	testCases := []struct {
		name     string
		file     multipart.File
		fileName string
		data     web.MonsterCreateRequest
	}{
		{
			name:     "success_create_monster",
			file:     file,
			fileName: fileName,
			data: web.MonsterCreateRequest{
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
			name:     "failed_create_monster_invalid_category_id",
			file:     file,
			fileName: fileName,
			data: web.MonsterCreateRequest{
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
			name:     "failed_create_monster_invalid_types_id",
			file:     file,
			fileName: fileName,
			data: web.MonsterCreateRequest{
				Name:        util.RandomString(10),
				CategoryID:  randCategories[0],
				Description: util.RandomString(20),
				Length:      54.3,
				Weight:      uint16(util.RandomInt(50, 500)),
				Hp:          uint16(util.RandomInt(50, 500)),
				Attack:      uint16(util.RandomInt(50, 500)),
				Defends:     uint16(util.RandomInt(50, 500)),
				Speed:       uint16(util.RandomInt(50, 500)),
				TypeID:      []string{"558160ef-e8f5-4951-b5f4-feeb0815b510", "d5a8d4bb-eb0a-44a4-ae46-eb2af2b2002d"},
			},
		},
	}

	var monster domain.Monster
	// Test
	for i := range testCases {
		tc := testCases[i]

		// Create
		newMonster, err := usecaseMonster.Create(context.Background(), tc.data, tc.file, tc.fileName)
		if tc.name != "success_create_monster" {
			require.Error(t, err)
			require.Equal(t, "invalid category id or type id, please check valid id in each of their list", err.Error())
		} else {
			require.NoError(t, err)
			require.NotEmpty(t, newMonster.ID)
			require.NotEmpty(t, newMonster.CreatedAt)
			require.NotEmpty(t, newMonster.UpdatedAt)
			require.NotEmpty(t, newMonster.ImageName)
			require.NotEmpty(t, newMonster.ImageURL)

			require.Equal(t, tc.data.Name, newMonster.Name)
			require.Equal(t, tc.data.CategoryID, newMonster.CategoryID)
			require.Equal(t, tc.data.Description, newMonster.Description)
			require.Equal(t, tc.data.Length, newMonster.Length)
			require.Equal(t, tc.data.Weight, newMonster.Weight)
			require.Equal(t, tc.data.Hp, newMonster.Hp)
			require.Equal(t, tc.data.Attack, newMonster.Attack)
			require.Equal(t, tc.data.Defends, newMonster.Defends)
			require.Equal(t, tc.data.Speed, newMonster.Speed)

			monster = newMonster
		}
	}

	return monster
}

func TestCreateMonsterUsecase(t *testing.T) {
	RandomCreateMonsterUsecase(t)
}

func TestFindAllMonsterUsecase(t *testing.T) {
	// Create random monsters
	newMonster, randTypes := RandomCreateMonster(t)

	repositoryMonster := repository.NewMonsterRespository(ConnTest)
	usecaseMonster := usecase.NewUsecaseMonster(repositoryMonster)
	ctx := context.Background()

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

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			// Find All
			monsters, err := usecaseMonster.FindAll(ctx, tc.queryParameter)
			require.NoError(t, err)

			for _, monster := range monsters {
				require.NotEmpty(t, monster.ID)
				require.NotEmpty(t, monster.Name)
				require.NotEmpty(t, monster.Category.Name)
				require.NotEmpty(t, strconv.FormatBool(monster.Catched))
				require.NotEmpty(t, monster.ImageName)
				require.NotEmpty(t, monster.ImageURL)

				require.NotEqual(t, 0, len(monster.Types))
				for i := 0; i < len(monster.Types); i++ {
					require.NotEmpty(t, monster.Types[i].Name)
				}
			}
		})
	}
}

func TestFindByIDMonsterUsecase(t *testing.T) {
	// Create random monsters
	newMonster, _ := RandomCreateMonster(t)
	repositoryMonster := repository.NewMonsterRespository(ConnTest)
	usecaseMonster := usecase.NewUsecaseMonster(repositoryMonster)
	ctx := context.Background()

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
			// Find by id
			monster, err := usecaseMonster.FindByID(ctx, tc.idMonster)

			if tc.name == "find_by_id_monster_success" {
				require.NoError(t, err)

				require.NotEmpty(t, monster.Category.Name)

				require.Equal(t, newMonster.ID, monster.ID)
				require.Equal(t, newMonster.Name, monster.Name)
				require.Equal(t, newMonster.Description, monster.Description)
				require.Equal(t, newMonster.Length, monster.Length)
				require.Equal(t, newMonster.Weight, monster.Weight)
				require.Equal(t, newMonster.Hp, monster.Hp)
				require.Equal(t, newMonster.Attack, monster.Attack)
				require.Equal(t, newMonster.Defends, monster.Defends)
				require.Equal(t, newMonster.Speed, monster.Speed)
				require.Equal(t, newMonster.Catched, monster.Catched)
				require.Equal(t, newMonster.ImageName, monster.ImageName)
				require.Equal(t, newMonster.ImageURL, monster.ImageURL)

				for _, ty := range monster.Types {
					require.NotEmpty(t, ty.Name)
				}
			} else {
				require.Error(t, err)
				var errTest error
				errMessage := fmt.Sprintf("monster with id %s not found", tc.idMonster)
				errTest = errors.New(errMessage)
				require.Equal(t, errTest, err)
			}
		})
	}
}

func TestUpdateMonsterUsecase(t *testing.T) {
	// File name format
	now := time.Now()
	nowRFC3339 := now.Format(time.RFC3339)
	fileName := fmt.Sprintf(`%s_%v_%s`, "usecase_update_test", nowRFC3339, "image.png")

	// Create random monsters
	newMonster, _ := RandomCreateMonster(t)
	repositoryMonster := repository.NewMonsterRespository(ConnTest)
	usecaseMonster := usecase.NewUsecaseMonster(repositoryMonster)

	var randTypes []string
	var randCategories []string
	// Get data random category and type
	for i := 0; i < 3; i++ {
		randCategory, randType := RandomCategoryAndType()
		randTypes = append(randTypes, randType)
		randCategories = append(randCategories, randCategory)
	}

	testCases := []struct {
		id       string
		name     string
		fileName string
		req      web.MonsterUpdateRequest
	}{
		{
			id:       newMonster.ID,
			name:     "update_success_with_field_empty",
			fileName: "",
			req: web.MonsterUpdateRequest{
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
			id:       newMonster.ID,
			name:     "update_success_with_image",
			fileName: fileName,
			req: web.MonsterUpdateRequest{
				Name:        "UPDATED",
				CategoryID:  randCategories[0],
				Description: "UPDATED",
				Length:      "1.1",
				Weight:      "1",
				Hp:          "1",
				Attack:      "1",
				Defends:     "1",
				Speed:       "1",
				Catched:     "true",
				TypeID:      randTypes,
			},
		},
		{
			id:       newMonster.ID,
			name:     "update_failed_category_id_invalid",
			fileName: "",
			req: web.MonsterUpdateRequest{
				Name:        "UPDATED",
				CategoryID:  "719d94b8-81a8-48ba-8052-0bdeda9643ad",
				Description: "UPDATED",
				Length:      "1.1",
				Weight:      "1",
				Hp:          "1",
				Attack:      "1",
				Defends:     "1",
				Speed:       "1",
				Catched:     "true",
				TypeID:      randTypes,
			},
		},
		{
			id:       newMonster.ID,
			name:     "update_failed_types_id_invalid",
			fileName: "",
			req: web.MonsterUpdateRequest{
				Name:        "UPDATED",
				CategoryID:  randCategories[0],
				Description: "UPDATED",
				Length:      "1.1",
				Weight:      "1",
				Hp:          "1",
				Attack:      "1",
				Defends:     "1",
				Speed:       "1",
				Catched:     "true",
				TypeID:      []string{"558160ef-e8f5-4951-b5f4-feeb0815b510", "d5a8d4bb-eb0a-44a4-ae46-eb2af2b2002d"},
			},
		},
	}

	// Test
	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			// Read file from local
			file, _ := os.Open("file_sample/image.png")
			defer file.Close()

			// Update
			updatedMonster, err := usecaseMonster.Update(context.Background(), tc.id, tc.req, file, tc.fileName)

			if tc.name == "update_failed_category_id_invalid" || tc.name == "update_failed_types_id_invalid" {
				var errTest error
				errMessage := fmt.Sprint("invalid category id or type id, please check valid id in each of their list")
				errTest = errors.New(errMessage)
				require.Equal(t, errTest, err)
			} else if tc.name == "update_success_with_field_empty" {
				require.NoError(t, err)

				require.Equal(t, newMonster.ID, updatedMonster.ID)
				require.Equal(t, newMonster.Name, updatedMonster.Name)
				require.Equal(t, newMonster.CategoryID, updatedMonster.CategoryID)
				require.Equal(t, newMonster.Description, updatedMonster.Description)
				require.Equal(t, newMonster.Length, updatedMonster.Length)
				require.Equal(t, newMonster.Weight, updatedMonster.Weight)
				require.Equal(t, newMonster.Hp, updatedMonster.Hp)
				require.Equal(t, newMonster.Attack, updatedMonster.Attack)
				require.Equal(t, newMonster.Defends, updatedMonster.Defends)
				require.Equal(t, newMonster.Speed, updatedMonster.Speed)
				require.Equal(t, newMonster.Catched, updatedMonster.Catched)
				require.Equal(t, newMonster.ImageName, updatedMonster.ImageName)
				require.Equal(t, newMonster.ImageURL, updatedMonster.ImageURL)

				for i := 0; i < len(updatedMonster.TypeID); i++ {
					require.Equal(t, newMonster.TypeID[i], updatedMonster.TypeID[i])
				}
			} else if tc.name == "update_success_with_image" {
				require.NoError(t, err)

				require.Equal(t, newMonster.ID, updatedMonster.ID)

				require.NotEqual(t, newMonster.Name, updatedMonster.Name)
				require.NotEqual(t, newMonster.Category.Name, updatedMonster.Category.Name)
				require.NotEqual(t, newMonster.Description, updatedMonster.Description)
				require.NotEqual(t, newMonster.Length, updatedMonster.Length)
				require.NotEqual(t, newMonster.Weight, updatedMonster.Weight)
				require.NotEqual(t, newMonster.Hp, updatedMonster.Hp)
				require.NotEqual(t, newMonster.Attack, updatedMonster.Attack)
				require.NotEqual(t, newMonster.Defends, updatedMonster.Defends)
				require.NotEqual(t, newMonster.Speed, updatedMonster.Speed)
				require.NotEqual(t, newMonster.Catched, updatedMonster.Catched)

				require.NotEqual(t, newMonster.ImageName, updatedMonster.ImageName)
				require.NotEqual(t, newMonster.ImageURL, updatedMonster.ImageURL)

				for i := 0; i < len(updatedMonster.TypeID); i++ {
					require.NotEqual(t, newMonster.TypeID[i], updatedMonster.TypeID[i])
				}
			}
			// else {
			// 	// If tc.name == "update_success_without_image"
			// 	require.NoError(t, err)

			// 	require.Equal(t, newMonster.ID, updatedMonster.ID)
			// 	// require.Equal(t, newMonster.ImageName, updatedMonster.ImageName)
			// 	// require.Equal(t, newMonster.ImageURL, updatedMonster.ImageURL)

			// 	require.NotEqual(t, newMonster.Name, updatedMonster.Name)
			// 	require.NotEqual(t, newMonster.CategoryID, updatedMonster.CategoryID)
			// 	require.NotEqual(t, newMonster.Description, updatedMonster.Description)
			// 	require.NotEqual(t, newMonster.Length, updatedMonster.Length)
			// 	require.NotEqual(t, newMonster.Weight, updatedMonster.Weight)
			// 	require.NotEqual(t, newMonster.Hp, updatedMonster.Hp)
			// 	require.NotEqual(t, newMonster.Attack, updatedMonster.Attack)
			// 	require.NotEqual(t, newMonster.Defends, updatedMonster.Defends)
			// 	require.NotEqual(t, newMonster.Speed, updatedMonster.Speed)
			// 	require.NotEqual(t, newMonster.Catched, updatedMonster.Catched)

			// 	for i := 0; i < len(updatedMonster.TypeID); i++ {
			// 		require.NotEqual(t, newMonster.TypeID[i], updatedMonster.TypeID[i])
			// 	}
			// }

		})
	}
}

func TestUpdateMarkMonsterCapturedUsecase(t *testing.T) {
	// Create random monsters
	newMonster, _ := RandomCreateMonster(t)
	repositoryMonster := repository.NewMonsterRespository(ConnTest)
	usecaseMonster := usecase.NewUsecaseMonster(repositoryMonster)

	testCases := []struct {
		name      string
		idMonster string
		reqUpdate web.MonsterUpdateRequestMonsterCapture
	}{
		{
			name:      "update_mark_monster_captured_monster_success",
			idMonster: newMonster.ID,
			reqUpdate: web.MonsterUpdateRequestMonsterCapture{
				Catched: true,
			},
		},
		{
			name:      "update_mark_monster_captured_monster_failed_monster_not_found",
			idMonster: "368bd987-dec6-4405-a036-bc1232db21b2",
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
			ok, err := usecaseMonster.UpdateMarkMonsterCaptured(context.Background(), tc.idMonster, tc.reqUpdate)

			monsterUpdated, _ := usecaseMonster.FindByID(context.Background(), newMonster.ID)

			if tc.name == "update_mark_monster_captured_monster_success" {
				require.NoError(t, err)
				require.True(t, ok)
				require.True(t, monsterUpdated.Catched)
			} else {
				require.Error(t, err)
				var errTest error
				errMessage := fmt.Sprintf("monster with id %s not found", tc.idMonster)
				errTest = errors.New(errMessage)
				require.Equal(t, errTest, err)
			}
		})
	}
}

func TestDeleteMonsterUsecase(t *testing.T) {
	newMonster := RandomCreateMonsterUsecase(t)

	repositoryMonster := repository.NewMonsterRespository(ConnTest)
	usecaseMonster := usecase.NewUsecaseMonster(repositoryMonster)

	testCases := []struct {
		name      string
		idMonster string
	}{
		{
			name:      "delete_monster_success",
			idMonster: newMonster.ID,
		},
		{
			name:      "delete_monster_failed_monster_not_found",
			idMonster: "368bd987-dec6-4405-a036-bc1232db21b2",
		},
	}

	// Test
	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			// Delete
			ok, err := usecaseMonster.Delete(context.Background(), tc.idMonster)

			if tc.name == "delete_monster_success" {
				require.NoError(t, err)
				require.True(t, ok)

				// Find by id
				monster, err := usecaseMonster.FindByID(context.Background(), tc.idMonster)
				require.Empty(t, monster.ID)

				msg := fmt.Sprintf("monster with id %s not found", tc.idMonster)
				require.Equal(t, msg, err.Error())
			} else {
				require.Error(t, err)

				msg := fmt.Sprintf("monster with id %s not found", tc.idMonster)
				require.Equal(t, msg, err.Error())
			}
		})
	}
}
