package tests

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/letenk/pokedex/models/web"
	"github.com/letenk/pokedex/repository"
	"github.com/letenk/pokedex/usecase"
	"github.com/letenk/pokedex/util"
	"github.com/stretchr/testify/require"
)

func TestCreateMonsterUsecase(t *testing.T) {
	t.Parallel()

	// Read file from local
	file, _ := os.Open("file_sample/image.png")
	defer file.Close()

	// File name formate
	now := time.Now()
	nowRFC3339 := now.Format(time.RFC3339)
	fileName := fmt.Sprintf(`%s-%v-%s`, "usecasetest", nowRFC3339, "image.png")

	repositoryMonster := repository.NewMonsterRespository(ConnTest)
	usecaseMonster := usecase.NewUsecaseMonster(repositoryMonster)

	// Get data random category and type
	randCategory, randType := RandomCategoryAndType()

	// Sample data
	monster := web.MonsterCreateRequest{
		Name:        util.RandomString(10),
		CategoryID:  randCategory,
		Description: util.RandomString(20),
		Length:      54.3,
		Weight:      uint16(util.RandomInt(50, 500)),
		Hp:          uint16(util.RandomInt(50, 500)),
		Attack:      uint16(util.RandomInt(50, 500)),
		Defends:     uint16(util.RandomInt(50, 500)),
		Speed:       uint16(util.RandomInt(50, 500)),
		// Image:       file,
		TypeID: []string{randType, randType, randType},
	}

	// Create
	newMonster, err := usecaseMonster.Create(context.Background(), monster, file, fileName)
	require.NoError(t, err)

	require.NotEmpty(t, newMonster.ID)
	require.NotEmpty(t, newMonster.CreatedAt)
	require.NotEmpty(t, newMonster.UpdatedAt)
	require.NotEmpty(t, newMonster.Image)

	require.Equal(t, monster.Name, newMonster.Name)
	require.Equal(t, monster.CategoryID, newMonster.CategoryID)
	require.Equal(t, monster.Description, newMonster.Description)
	require.Equal(t, monster.Length, newMonster.Length)
	require.Equal(t, monster.Weight, newMonster.Weight)
	require.Equal(t, monster.Hp, newMonster.Hp)
	require.Equal(t, monster.Attack, newMonster.Attack)
	require.Equal(t, monster.Defends, newMonster.Defends)
	require.Equal(t, monster.Speed, newMonster.Speed)
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
				require.NotEmpty(t, monster.Image)

				require.NotEqual(t, 0, len(monster.Types))
				for i := 0; i < len(monster.Types); i++ {
					require.NotEmpty(t, monster.Types[i].Name)
				}
			}
		})
	}
}
