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

func TestFindByUsername(t *testing.T) {
	repository := repository.NewUserRepository(ConnTest)
	userRoleAdmin, _ := repository.FindByUsername(context.Background(), "admin")
	userRoleUser, _ := repository.FindByUsername(context.Background(), "user")

	testCases := []struct {
		name string
		user domain.User
	}{
		{
			name: "success_find_admin",
			user: domain.User{
				ID:       userRoleAdmin.ID,
				Username: userRoleAdmin.Username,
				Fullname: userRoleAdmin.Fullname,
				Password: "password",
				Role:     userRoleAdmin.Role,
			},
		},
		{
			name: "success_find_user",
			user: domain.User{
				ID:       userRoleUser.ID,
				Username: userRoleUser.Username,
				Fullname: userRoleUser.Fullname,
				Password: "password",
				Role:     userRoleUser.Role,
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

func TestFindByID(t *testing.T) {
	repository := repository.NewUserRepository(ConnTest)
	userRoleAdmin, _ := repository.FindByUsername(context.Background(), "admin")
	userRoleUser, _ := repository.FindByUsername(context.Background(), "user")

	testCases := []struct {
		name string
		user domain.User
	}{
		{
			name: "success_find_admin",
			user: domain.User{
				ID:       userRoleAdmin.ID,
				Username: userRoleAdmin.Username,
				Fullname: userRoleAdmin.Fullname,
				Password: "password",
				Role:     userRoleAdmin.Role,
			},
		},
		{
			name: "success_find_user",
			user: domain.User{
				ID:       userRoleUser.ID,
				Username: userRoleUser.Username,
				Fullname: userRoleUser.Fullname,
				Password: "password",
				Role:     userRoleUser.Role,
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Find by id
			ctx := context.Background()
			user, err := repository.FindByID(ctx, tc.user.ID)

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
