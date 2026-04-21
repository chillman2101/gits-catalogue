package service_test

import (
	"os"
	"testing"

	"github.com/chillman2101/gits-catalogue/internal/mocks"
	"github.com/chillman2101/gits-catalogue/internal/model"
	"github.com/chillman2101/gits-catalogue/internal/service"
	res_err "github.com/chillman2101/gits-catalogue/pkg/response"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func init() {
	os.Setenv("JWT_SECRET", "test-secret-key")
	os.Setenv("JWT_REFRESH_SECRET", "test-refresh-secret-key")
}

func newAuthService(repo *mocks.MockUserRepository) service.AuthService {
	return service.NewAuthService(repo)
}

func hashPassword(password string) string {
	h, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	return string(h)
}

func TestAuthService_Register(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := new(mocks.MockUserRepository)
		svc := newAuthService(repo)

		// Create is called with any *model.User (bcrypt hash differs each run)
		repo.On("Create", mock.MatchedBy(func(u *model.User) bool {
			return u.Email == "test@example.com" && u.Password != ""
		})).Return(nil)

		err := svc.Register("test@example.com", "password123")

		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("duplicate email", func(t *testing.T) {
		repo := new(mocks.MockUserRepository)
		svc := newAuthService(repo)

		repo.On("Create", mock.MatchedBy(func(u *model.User) bool {
			return u.Email == "existing@example.com"
		})).Return(res_err.ErrUserConflict)

		err := svc.Register("existing@example.com", "password123")

		assert.ErrorIs(t, err, res_err.ErrUserConflict)
		repo.AssertExpectations(t)
	})
}

func TestAuthService_Login(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := new(mocks.MockUserRepository)
		svc := newAuthService(repo)
		user := &model.User{ID: 1, Email: "test@example.com", Password: hashPassword("password123")}

		repo.On("FindByEmail", "test@example.com").Return(user, nil)
		repo.On("UpdateTokenHash", uint(1), mock.AnythingOfType("string")).Return(nil)
		repo.On("UpdateRefreshTokenHash", uint(1), mock.AnythingOfType("string")).Return(nil)

		pair, err := svc.Login("test@example.com", "password123")

		assert.NoError(t, err)
		assert.NotEmpty(t, pair.AccessToken)
		assert.NotEmpty(t, pair.RefreshToken)
		repo.AssertExpectations(t)
	})

	t.Run("user not found", func(t *testing.T) {
		repo := new(mocks.MockUserRepository)
		svc := newAuthService(repo)

		repo.On("FindByEmail", "notfound@example.com").Return(nil, res_err.ErrUserNotFound)

		pair, err := svc.Login("notfound@example.com", "password123")

		assert.Error(t, err)
		assert.Nil(t, pair)
		repo.AssertExpectations(t)
	})

	t.Run("wrong password", func(t *testing.T) {
		repo := new(mocks.MockUserRepository)
		svc := newAuthService(repo)
		user := &model.User{ID: 1, Email: "test@example.com", Password: hashPassword("correctpassword")}

		repo.On("FindByEmail", "test@example.com").Return(user, nil)

		pair, err := svc.Login("test@example.com", "wrongpassword")

		assert.ErrorIs(t, err, res_err.ErrUserInvalidCredentials)
		assert.Nil(t, pair)
		repo.AssertExpectations(t)
	})
}

func TestAuthService_Logout(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := new(mocks.MockUserRepository)
		svc := newAuthService(repo)

		repo.On("UpdateTokenHash", uint(1), "").Return(nil)
		repo.On("UpdateRefreshTokenHash", uint(1), "").Return(nil)

		err := svc.Logout(1)

		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("update token hash fails", func(t *testing.T) {
		repo := new(mocks.MockUserRepository)
		svc := newAuthService(repo)

		repo.On("UpdateTokenHash", uint(1), "").Return(res_err.ErrUserNotFound)

		err := svc.Logout(1)

		assert.Error(t, err)
		repo.AssertExpectations(t)
	})
}

