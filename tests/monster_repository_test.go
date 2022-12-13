package tests

import (
	"context"
	"testing"

	"github.com/letenk/pokedex/models/domain"
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

func TestCreateMonsterRepository(t *testing.T) {
	t.Parallel()
	repositoryMonster := repository.NewMonsterRespository(ConnTest)

	// Get data random category and type
	randCategory, randType := RandomCategoryAndType()

	// Sample data
	monster := domain.Monster{
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
	}

	// Create
	newMonster, err := repositoryMonster.Create(context.Background(), monster)
	require.NoError(t, err)

	require.NotEmpty(t, newMonster.ID)
	require.NotEmpty(t, newMonster.CreatedAt)
	require.NotEmpty(t, newMonster.UpdatedAt)

	require.Equal(t, monster.Name, newMonster.Name)
	require.Equal(t, monster.CategoryID, newMonster.CategoryID)
	require.Equal(t, monster.Description, newMonster.Description)
	require.Equal(t, monster.Length, newMonster.Length)
	require.Equal(t, monster.Weight, newMonster.Weight)
	require.Equal(t, monster.Hp, newMonster.Hp)
	require.Equal(t, monster.Attack, newMonster.Attack)
	require.Equal(t, monster.Defends, newMonster.Defends)
	require.Equal(t, monster.Speed, newMonster.Speed)
	require.Equal(t, monster.Image, newMonster.Image)
}
