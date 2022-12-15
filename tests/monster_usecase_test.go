package tests

import (
	"context"
	"fmt"
	"os"
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
	// Create image
	_ = CreateImage()

	// Read file from local
	file, _ := os.Open("image.png")
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
