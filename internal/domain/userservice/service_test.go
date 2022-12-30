package userservice

import (
	"fmt"
	"github.com/kamkali/go-timeline/internal/domain"
	"github.com/kamkali/go-timeline/internal/domain/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
	"testing"
)

func TestChangePassword(t *testing.T) {
	ctx := context.Background()
	repo := mocks.NewUserRepository(t)
	service := New(nil, repo)
	email := "test@example.com"
	password := "newpassword"

	t.Run("Test ChangePassword with valid input", func(t *testing.T) {
		repo.On("ChangePassword", ctx, email, password).
			Return(nil)

		err := service.ChangePassword(ctx, email, password)
		require.NoError(t, err)
	})

	t.Run("Test ChangePassword with empty password", func(t *testing.T) {
		password = ""
		err := service.ChangePassword(ctx, email, password)
		require.Error(t, err)
	})
}

func TestCreateUser(t *testing.T) {
	tests := map[string]struct {
		user        domain.User
		setMockFunc func(*mocks.UserRepository, domain.User)
		wantErr     assert.ErrorAssertionFunc
	}{
		"valid user": {
			user: domain.User{
				Email:    "test@example.com",
				Password: "password",
			},
			setMockFunc: func(repo *mocks.UserRepository, user domain.User) {
				repo.On("CreateUser", mock.Anything, user).Return(nil)

			},
			wantErr: assert.NoError,
		},
		"empty email": {
			user: domain.User{
				Email:    "",
				Password: "password",
			},
			wantErr: assert.Error,
		},
		"empty password": {
			user: domain.User{
				Email:    "test@example.com",
				Password: "",
			},
			wantErr: assert.Error,
		},
		"error from repository": {
			user: domain.User{
				Email:    "test@example.com",
				Password: "password",
			},
			setMockFunc: func(repo *mocks.UserRepository, user domain.User) {
				repo.On("CreateUser", mock.Anything, user).Return(fmt.Errorf("error creating user"))
			},
			wantErr: assert.Error,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t1 *testing.T) {
			ctx := context.Background()
			repo := mocks.NewUserRepository(t)
			if tt.setMockFunc != nil {
				tt.setMockFunc(repo, tt.user)
			}
			service := New(nil, repo)
			tt.wantErr(t1, service.CreateUser(ctx, tt.user), fmt.Sprintf("CreateUser: %v", tt.user))
		})
	}
}

func TestLoginUser(t *testing.T) {
	mockErr := fmt.Errorf("error creating user")
	tests := map[string]struct {
		user        domain.User
		setMockFunc func(*mocks.UserRepository, domain.User)
		want        domain.User
		wantErr     func(t *testing.T, err error)
	}{
		"valid user": {
			user: domain.User{
				Email:    "test@example.com",
				Password: "password",
			},
			setMockFunc: func(repo *mocks.UserRepository, user domain.User) {
				passHash, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
				require.NoError(t, err)
				repo.On("GetUserByEmail", mock.Anything, user.Email).
					Return(domain.User{
						Email:    "test@example.com",
						Password: string(passHash),
					}, nil)
			},
			want: domain.User{
				Email: "test@example.com",
			},
			wantErr: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		"empty email": {
			user: domain.User{
				Email:    "",
				Password: "password",
			},
			wantErr: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
			want: domain.User{},
		},
		"empty password": {
			user: domain.User{
				Email:    "test@example.com",
				Password: "",
			},
			wantErr: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
			want: domain.User{},
		},
		"error from repository": {
			user: domain.User{
				Email:    "test@example.com",
				Password: "password",
			},
			setMockFunc: func(repo *mocks.UserRepository, user domain.User) {
				repo.On("GetUserByEmail", mock.Anything, user.Email).
					Return(domain.User{}, mockErr)
			},
			wantErr: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, mockErr)
			},
			want: domain.User{},
		},
		"invalid password": {
			user: domain.User{
				Email:    "test@example.com",
				Password: "password",
			},
			setMockFunc: func(repo *mocks.UserRepository, user domain.User) {
				passHash, err := bcrypt.GenerateFromPassword([]byte("something-different"), bcrypt.DefaultCost)
				require.NoError(t, err)
				repo.On("GetUserByEmail", mock.Anything, user.Email).
					Return(domain.User{
						Email:    "test@example.com",
						Password: string(passHash),
					}, nil)
			},
			wantErr: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, domain.ErrUnauthorized)
			},
			want: domain.User{},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t1 *testing.T) {
			ctx := context.Background()
			repo := mocks.NewUserRepository(t)
			if tt.setMockFunc != nil {
				tt.setMockFunc(repo, tt.user)
			}
			service := New(nil, repo)
			got, err := service.LoginUser(ctx, &tt.user)
			tt.wantErr(t, err)
			assert.Equal(t1, tt.want.Email, got.Email)
		})
	}
}
