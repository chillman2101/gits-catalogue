package mocks

import (
	"github.com/chillman2101/gits-catalogue/internal/model"
	"github.com/chillman2101/gits-catalogue/internal/query"
	"github.com/stretchr/testify/mock"
)

type MockPublisherRepository struct {
	mock.Mock
}

func (m *MockPublisherRepository) FindAll(p query.Params) ([]model.Publisher, int64, error) {
	args := m.Called(p)
	return args.Get(0).([]model.Publisher), args.Get(1).(int64), args.Error(2)
}

func (m *MockPublisherRepository) FindByID(id uint) (*model.Publisher, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Publisher), args.Error(1)
}

func (m *MockPublisherRepository) FindByName(name string) (*model.Publisher, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Publisher), args.Error(1)
}

func (m *MockPublisherRepository) Create(publisher *model.Publisher) error {
	args := m.Called(publisher)
	return args.Error(0)
}

func (m *MockPublisherRepository) Update(publisher *model.Publisher) error {
	args := m.Called(publisher)
	return args.Error(0)
}

func (m *MockPublisherRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}
