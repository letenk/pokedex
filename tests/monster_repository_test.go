package tests

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/letenk/pokedex/models/domain"
	"github.com/letenk/pokedex/models/web"
	"github.com/letenk/pokedex/repository"
	"github.com/letenk/pokedex/util"
	"github.com/stretchr/testify/require"
)

func RandomCategoryAndType() (string, string) {
	// Random data id categories
	repositoryCategory := repository.NewCategoryRepository(ConnTest)
	categories, _ := repositoryCategory.FindAll(context.Background())

	var category []string
	for _, data := range categories {
		category = append(category, data.ID)
	}
	randCategory := util.RandomStringFromSet(category...)

	// Random data id types
	repositoryType := repository.NewTypeRespository(ConnTest)
	dataTypes, _ := repositoryType.FindAll(context.Background())

	var types []string
	for _, data := range dataTypes {
		types = append(types, data.ID)
	}
	randType := util.RandomStringFromSet(types...)

	return randCategory, randType
}

func RandomCreateMonster(t *testing.T) (domain.Monster, []string) {
	t.Parallel()
	repositoryMonster := repository.NewMonsterRespository(ConnTest)

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
		name string
		data domain.Monster
	}{
		{
			name: "success_create_monster",
			data: domain.Monster{
				Name:        util.RandomString(10),
				CategoryID:  randCategories[0],
				Description: util.RandomString(20),
				Length:      54.3,
				Weight:      uint16(util.RandomInt(50, 500)),
				Hp:          uint16(util.RandomInt(50, 500)),
				Attack:      uint16(util.RandomInt(50, 500)),
				Defends:     uint16(util.RandomInt(50, 500)),
				Speed:       uint16(util.RandomInt(50, 500)),
				ImageName:   util.RandomString(10),
				ImageURL:    util.RandomString(10),
				TypeID:      randTypes,
			},
		},
		{
			name: "failed_create_monster_invalid_category_id",
			data: domain.Monster{
				Name:        util.RandomString(10),
				CategoryID:  "4562482c-7acd-4daf-901f-d95c7a7afd65",
				Description: util.RandomString(20),
				Length:      54.3,
				Weight:      uint16(util.RandomInt(50, 500)),
				Hp:          uint16(util.RandomInt(50, 500)),
				Attack:      uint16(util.RandomInt(50, 500)),
				Defends:     uint16(util.RandomInt(50, 500)),
				Speed:       uint16(util.RandomInt(50, 500)),
				ImageName:   util.RandomString(10),
				ImageURL:    util.RandomString(10),
				TypeID:      randTypes,
			},
		},
		{
			name: "failed_create_monster_invalid_types_id",
			data: domain.Monster{
				Name:        util.RandomString(10),
				CategoryID:  randCategories[0],
				Description: util.RandomString(20),
				Length:      54.3,
				Weight:      uint16(util.RandomInt(50, 500)),
				Hp:          uint16(util.RandomInt(50, 500)),
				Attack:      uint16(util.RandomInt(50, 500)),
				Defends:     uint16(util.RandomInt(50, 500)),
				Speed:       uint16(util.RandomInt(50, 500)),
				ImageName:   util.RandomString(10),
				ImageURL:    util.RandomString(10),
				TypeID:      []string{"558160ef-e8f5-4951-b5f4-feeb0815b510", "d5a8d4bb-eb0a-44a4-ae46-eb2af2b2002d"},
			},
		},
	}

	var monster domain.Monster
	// Test
	for i := range testCases {
		tc := testCases[i]

		// Create
		if tc.name != "success_create_monster" {
			_, err := repositoryMonster.Create(context.Background(), tc.data)
			require.Error(t, err)
			errMessage := fmt.Sprint("invalid category id or type id, please check valid id in each of their list")
			require.Equal(t, errMessage, err.Error())
		} else {
			newMonster, err := repositoryMonster.Create(context.Background(), tc.data)
			require.NoError(t, err)
			require.NotEmpty(t, newMonster.ID)
			require.NotEmpty(t, newMonster.CreatedAt)
			require.NotEmpty(t, newMonster.UpdatedAt)

			require.Equal(t, tc.data.Name, newMonster.Name)
			require.Equal(t, tc.data.CategoryID, newMonster.CategoryID)
			require.Equal(t, tc.data.Description, newMonster.Description)
			require.Equal(t, tc.data.Length, newMonster.Length)
			require.Equal(t, tc.data.Weight, newMonster.Weight)
			require.Equal(t, tc.data.Hp, newMonster.Hp)
			require.Equal(t, tc.data.Attack, newMonster.Attack)
			require.Equal(t, tc.data.Defends, newMonster.Defends)
			require.Equal(t, tc.data.Speed, newMonster.Speed)
			require.Equal(t, tc.data.ImageName, newMonster.ImageName)
			require.Equal(t, tc.data.ImageURL, newMonster.ImageURL)

			monster = newMonster
		}
	}
	// Return result to use other test
	return monster, randTypes
}

func TestCreateMonsterRepository(t *testing.T) {
	RandomCreateMonster(t)
}

func TestFindAllMonsterRepository(t *testing.T) {
	// Create random monsters
	newMonster, randTypes := RandomCreateMonster(t)

	repositoryMonster := repository.NewMonsterRespository(ConnTest)
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
			name: "find_all_monsters_with_query_parameter_order_by_asc",
			queryParameter: web.MonsterQueryRequest{
				Sort:  "name",
				Order: "asc",
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
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			// Find All
			monsters, err := repositoryMonster.FindAll(ctx, tc.queryParameter)
			require.NoError(t, err)

			for _, monster := range monsters {
				require.NotEmpty(t, monster.ID)
				require.NotEmpty(t, monster.Name)
				require.NotEmpty(t, monster.Category.Name)
				require.NotEmpty(t, strconv.FormatBool(monster.Catched))
				require.NotEmpty(t, monster.ImageName)
				require.NotEmpty(t, monster.ImageURL)

				for i := 0; i < len(monster.Types); i++ {
					require.NotEmpty(t, monster.Types[i].Name)
				}
			}
		})
	}
}

func TestFindByIDMonsterRepository(t *testing.T) {
	// Create random monsters
	newMonster, _ := RandomCreateMonster(t)
	repositoryMonster := repository.NewMonsterRespository(ConnTest)
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
			t.Parallel()
			// Find by id
			monster, err := repositoryMonster.FindByID(ctx, tc.idMonster)

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

func TestUpdateMonsterRepository(t *testing.T) {
	// Create random monsters
	newMonster, _ := RandomCreateMonster(t)

	repositoryMonster := repository.NewMonsterRespository(ConnTest)
	ctx := context.Background()

	var randTypes []string
	var randCategories []string
	// Get data random category and type
	for i := 0; i < 3; i++ {
		randCategory, randType := RandomCategoryAndType()
		randTypes = append(randTypes, randType)
		randCategories = append(randCategories, randCategory)
	}

	testCases := []struct {
		name       string
		dataUpdate domain.Monster
	}{
		{
			name: "update_success_with_types",
			dataUpdate: domain.Monster{
				ID:          newMonster.ID,
				Name:        "UPDATED",
				CategoryID:  randCategories[0],
				Description: "UPDATED",
				Length:      1.1,
				Weight:      1,
				Hp:          1,
				Attack:      1,
				Defends:     1,
				Speed:       1,
				Catched:     true,
				ImageName:   "UPDATED",
				ImageURL:    "UPDATED",
				TypeID:      randTypes,
			},
		},
		{
			name: "update_success_without_types",
			dataUpdate: domain.Monster{
				ID:          newMonster.ID,
				Name:        "UPDATED",
				CategoryID:  randCategories[0],
				Description: "UPDATED",
				Length:      1.1,
				Weight:      1,
				Hp:          1,
				Attack:      1,
				Defends:     1,
				Speed:       1,
				Catched:     true,
				ImageName:   "UPDATED",
				ImageURL:    "UPDATED",
				// TypeID:      randTypes,
			},
		},
		{
			name: "update_failed_category_id_invalid",
			dataUpdate: domain.Monster{
				ID:          newMonster.ID,
				Name:        "UPDATED",
				CategoryID:  "3e1a7c6c-4272-42cf-a9d2-d0814b4454e9",
				Description: "UPDATED",
				Length:      1.1,
				Weight:      1,
				Hp:          1,
				Attack:      1,
				Defends:     1,
				Speed:       1,
				Catched:     true,
				ImageName:   "UPDATED",
				ImageURL:    "UPDATED",
				TypeID:      randTypes,
			},
		},
		{
			name: "update_failed_types_id_invalid",
			dataUpdate: domain.Monster{
				ID:          newMonster.ID,
				Name:        "UPDATED",
				CategoryID:  randCategories[0],
				Description: "UPDATED",
				Length:      1.1,
				Weight:      1,
				Hp:          1,
				Attack:      1,
				Defends:     1,
				Speed:       1,
				Catched:     true,
				ImageName:   "UPDATED",
				ImageURL:    "UPDATED",
				TypeID:      []string{"558160ef-e8f5-4951-b5f4-feeb0815b510", "d5a8d4bb-eb0a-44a4-ae46-eb2af2b2002d"},
			},
		},
	}

	// Test
	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			// Update
			_, err := repositoryMonster.Update(ctx, tc.dataUpdate)

			// Find monster by id
			updatedMonster, _ := repositoryMonster.FindByID(ctx, tc.dataUpdate.ID)

			if tc.name == "update_success_without_types" {
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
					require.Equal(t, newMonster.TypeID[i], updatedMonster.TypeID[i])
				}
			} else if tc.name == "update_failed_category_id_invalid" || tc.name == "update_failed_types_id_invalid" {
				require.Error(t, err)
				var errTest error
				errMessage := fmt.Sprintf("invalid category id or type id, please check valid id in each of their list")
				errTest = errors.New(errMessage)
				require.Equal(t, errTest, err)
			} else {
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
		})
	}
}

func TestDeleteMonsterRepository(t *testing.T) {
	newMonster, _ := RandomCreateMonster(t)

	repositoryMonster := repository.NewMonsterRespository(ConnTest)
	ctx := context.Background()

	ok, err := repositoryMonster.Delete(ctx, newMonster)

	require.NoError(t, err)
	require.True(t, ok)

	// Check monster in table
	monster, _ := repositoryMonster.FindByID(ctx, newMonster.ID)
	require.Empty(t, monster.ID)

}
