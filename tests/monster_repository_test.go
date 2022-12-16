package tests

import (
	"context"
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
				Image:       util.RandomString(10),
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
				Image:       util.RandomString(10),
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
				Image:       util.RandomString(10),
				TypeID:      []string{"558160ef-e8f5-4951-b5f4-feeb0815b510", "d5a8d4bb-eb0a-44a4-ae46-eb2af2b2002d"},
			},
		},
	}

	var monster domain.Monster
	// Test
	for i := range testCases {
		tc := testCases[i]

		// Create
		newMonster, err := repositoryMonster.Create(context.Background(), tc.data)
		if tc.name != "success_create_monster" {
			require.Error(t, err)
			require.Equal(t, "invalid category id or type id, please check valid id in each of their list", err.Error())
		} else {
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
			require.Equal(t, tc.data.Image, newMonster.Image)

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
				require.NotEmpty(t, monster.Image)

				require.NotEqual(t, 0, len(monster.Types))
				for i := 0; i < len(monster.Types); i++ {
					require.NotEmpty(t, monster.Types[i].Name)
				}
			}
		})
	}
}
