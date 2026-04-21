package mocks

import (
	"github.com/chillman2101/gits-catalogue/internal/model"
	"github.com/chillman2101/gits-catalogue/internal/query"
	"github.com/stretchr/testify/mock"
)

type MockBookService struct {
	mock.Mock
}

func (m *MockBookService) GetAll(p query.Params) ([]model.Book, int64, error) {
	args := m.Called(p)
	return args.Get(0).([]model.Book), args.Get(1).(int64), args.Error(2)
}

func (m *MockBookService) GetByID(id uint) (*model.Book, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Book), args.Error(1)
}

func (m *MockBookService) Create(book *model.Book) error {
	args := m.Called(book)
	return args.Error(0)
}

func (m *MockBookService) Update(id uint, book *model.Book) error {
	args := m.Called(id, book)
	return args.Error(0)
}

func (m *MockBookService) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}
