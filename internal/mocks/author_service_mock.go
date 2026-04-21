package mocks

import (
	"github.com/chillman2101/gits-catalogue/internal/model"
	"github.com/chillman2101/gits-catalogue/internal/query"
	"github.com/stretchr/testify/mock"
)

type MockAuthorService struct {
	mock.Mock
}

func (m *MockAuthorService) GetAll(p query.Params) ([]model.Author, int64, error) {
	args := m.Called(p)
	return args.Get(0).([]model.Author), args.Get(1).(int64), args.Error(2)
}

func (m *MockAuthorService) GetByID(id uint) (*model.Author, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Author), args.Error(1)
}

func (m *MockAuthorService) Create(author *model.Author) error {
	args := m.Called(author)
	return args.Error(0)
}

func (m *MockAuthorService) Update(id uint, author *model.Author) error {
	args := m.Called(id, author)
	return args.Error(0)
}

func (m *MockAuthorService) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}
