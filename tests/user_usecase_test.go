package tests

import (
	"testing"

	"github.com/letenk/pokedex/models/domain"
	"github.com/letenk/pokedex/models/web"
	"github.com/letenk/pokedex/repository"
	"github.com/letenk/pokedex/usecase"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
)

func TestLogin(t *testing.T) {
	repository := repository.NewUserRepository(ConnTest)
	usecase := usecase.NewUsecaseUser(repository)

	testCases := []struct {
		name string
		req  web.UserLoginRequest
	}{
		{
			name: "success_login",
			req: web.UserLoginRequest{
				Username: "admin",
				Password: "password",
			},
		},
		{
			name: "success_failed_username_wrong",
			req: web.UserLoginRequest{
				Username: "wrong",
				Password: "password",
			},
		},
		{
			name: "success_failed_password_wrong",
			req: web.UserLoginRequest{
				Username: "admin",
				Password: "wrong",
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			// Login
			ctx := context.Background()
			token, err := usecase.Login(ctx, tc.req)

			if tc.name == "success_login" {
				require.NotEmpty(t, token)
				require.NoError(t, err)
			} else {
				require.Empty(t, token)
				require.Error(t, err)
				require.Equal(t, "username or password incorrect", err.Error())
			}
		})
	}
}

func TestFindOneByID(t *testing.T) {
	repository := repository.NewUserRepository(ConnTest)
	usecase := usecase.NewUsecaseUser(repository)

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
			user, err := usecase.FindOneByID(ctx, tc.user.ID)

			require.NoError(t, err)

			require.Equal(t, tc.user.ID, user.ID)
			require.Equal(t, tc.user.Username, user.Username)
			require.Equal(t, tc.user.Fullname, user.Fullname)
			require.Equal(t, tc.user.Role, user.Role)

			// Compare password
			err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(tc.user.Password))
			require.NoError(t, err)
		})
	}
}
