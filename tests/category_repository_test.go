package tests

import (
	"context"
	"testing"

	"github.com/letenk/pokedex/repository"
	"github.com/stretchr/testify/require"
)

func TestFindAllCategories(t *testing.T) {
	t.Parallel()
	repository := repository.NewCategoryRepository(ConnTest)

	// Find all
	todos, err := repository.FindAll(context.Background())

	require.NoError(t, err)
	require.NotEqual(t, 0, len(todos))

	for _, data := range todos {
		require.NotEmpty(t, data.ID)
		require.NotEmpty(t, data.Name)
		require.NotEmpty(t, data.CreatedAt)
		require.NotEmpty(t, data.UpdatedAt)
	}
}
