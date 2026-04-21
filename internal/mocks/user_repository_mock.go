package mocks

import (
	"github.com/chillman2101/gits-catalogue/internal/model"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) FindByEmail(email string) (*model.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) FindByID(id uint) (*model.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) FindByIDAndTokenHash(id uint, tokenHash string) (*model.User, error) {
	args := m.Called(id, tokenHash)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) FindByIDAndRefreshTokenHash(id uint, refreshTokenHash string) (*model.User, error) {
	args := m.Called(id, refreshTokenHash)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) Create(user *model.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) UpdateTokenHash(id uint, hash string) error {
	args := m.Called(id, hash)
	return args.Error(0)
}

func (m *MockUserRepository) UpdateRefreshTokenHash(id uint, hash string) error {
	args := m.Called(id, hash)
	return args.Error(0)
}
