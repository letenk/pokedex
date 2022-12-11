package tests

import (
	"testing"

	"github.com/letenk/pokedex/models/domain"
	"github.com/letenk/pokedex/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
)

func TestFindByID(t *testing.T) {
	repository := repository.NewUserRepository(ConnTest)

	testCases := []struct {
		name string
		user domain.User
	}{
		{
			name: "success_find_admin",
			user: domain.User{
				ID:       "34a79b77-5201-4fc8-8b8d-7e43350badd4",
				Username: "admin",
				Fullname: "ADMIN",
				Password: "password",
				Role:     "admin",
			},
		},
		{
			name: "success_find_user",
			user: domain.User{
				ID:       "692cf812-e480-443e-b6d5-11b280047dcf",
				Username: "user",
				Fullname: "USER",
				Password: "password",
				Role:     "user",
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Find by id
			ctx := context.Background()
			user, err := repository.FindByUsername(ctx, tc.user.Username)

			require.NoError(t, err)

			require.Equal(t, tc.user.ID, user.ID)
			require.Equal(t, tc.user.Username, user.Username)
			require.Equal(t, tc.user.Fullname, user.Fullname)
			require.Equal(t, tc.user.Role, user.Role)

			// Compare password
			err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(tc.user.Password))
			assert.NoError(t, err)
		})
	}
}
