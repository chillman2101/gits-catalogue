package mocks

import (
	"github.com/chillman2101/gits-catalogue/internal/service"
	"github.com/stretchr/testify/mock"
)

type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Register(email, password string) error {
	args := m.Called(email, password)
	return args.Error(0)
}

func (m *MockAuthService) Login(email, password string) (*service.TokenPair, error) {
	args := m.Called(email, password)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*service.TokenPair), args.Error(1)
}

func (m *MockAuthService) Refresh(refreshToken string) (*service.TokenPair, error) {
	args := m.Called(refreshToken)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*service.TokenPair), args.Error(1)
}

func (m *MockAuthService) Logout(userID uint) error {
	args := m.Called(userID)
	return args.Error(0)
}
