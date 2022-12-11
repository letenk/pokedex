package tests

import (
	"testing"

	"github.com/letenk/pokedex/models/web"
	"github.com/letenk/pokedex/repository"
	"github.com/letenk/pokedex/usecase"
	"github.com/stretchr/testify/assert"
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
				assert.NotEmpty(t, token)
				assert.NoError(t, err)
			} else {
				assert.Empty(t, token)
				assert.Error(t, err)
				assert.Equal(t, "username or password incorrect", err.Error())
			}
		})
	}
}
