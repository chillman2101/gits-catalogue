package service_test

import (
	"errors"
	"testing"

	"github.com/chillman2101/gits-catalogue/internal/mocks"
	"github.com/chillman2101/gits-catalogue/internal/model"
	"github.com/chillman2101/gits-catalogue/internal/query"
	"github.com/chillman2101/gits-catalogue/internal/service"
	"github.com/chillman2101/gits-catalogue/pkg/redis/cache"
	res_err "github.com/chillman2101/gits-catalogue/pkg/response"
	"github.com/stretchr/testify/assert"
)

func newPublisherService(repo *mocks.MockPublisherRepository) service.PublisherService {
	return service.NewPublisherService(repo, cache.NewCacheHelper(nil))
}

func TestPublisherService_GetAll(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := new(mocks.MockPublisherRepository)
		svc := newPublisherService(repo)
		p := query.Params{Page: 1, Limit: 10}
		publishers := []model.Publisher{{ID: 1, Name: "Gramedia"}}

		repo.On("FindAll", p).Return(publishers, int64(1), nil)

		result, total, err := svc.GetAll(p)

		assert.NoError(t, err)
		assert.Equal(t, int64(1), total)
		assert.Len(t, result, 1)
		assert.Equal(t, "Gramedia", result[0].Name)
		repo.AssertExpectations(t)
	})

	t.Run("repository error", func(t *testing.T) {
		repo := new(mocks.MockPublisherRepository)
		svc := newPublisherService(repo)
		p := query.Params{Page: 1, Limit: 10}

		repo.On("FindAll", p).Return([]model.Publisher{}, int64(0), errors.New("db error"))

		_, _, err := svc.GetAll(p)

		assert.Error(t, err)
		repo.AssertExpectations(t)
	})
}

func TestPublisherService_GetByID(t *testing.T) {
	t.Run("found", func(t *testing.T) {
		repo := new(mocks.MockPublisherRepository)
		svc := newPublisherService(repo)
		publisher := &model.Publisher{ID: 1, Name: "Gramedia"}

		repo.On("FindByID", uint(1)).Return(publisher, nil)

		result, err := svc.GetByID(1)

		assert.NoError(t, err)
		assert.Equal(t, "Gramedia", result.Name)
		repo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		repo := new(mocks.MockPublisherRepository)
		svc := newPublisherService(repo)

		repo.On("FindByID", uint(99)).Return(nil, res_err.ErrPublisherNotFound)

		result, err := svc.GetByID(99)

		assert.Error(t, err)
		assert.Nil(t, result)
		repo.AssertExpectations(t)
	})
}

func TestPublisherService_Create(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := new(mocks.MockPublisherRepository)
		svc := newPublisherService(repo)
		publisher := &model.Publisher{Name: "New Publisher"}

		repo.On("FindByName", "New Publisher").Return(nil, nil)
		repo.On("Create", publisher).Return(nil)

		err := svc.Create(publisher)

		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("duplicate name", func(t *testing.T) {
		repo := new(mocks.MockPublisherRepository)
		svc := newPublisherService(repo)
		publisher := &model.Publisher{Name: "Existing Publisher"}
		existing := &model.Publisher{ID: 1, Name: "Existing Publisher"}

		repo.On("FindByName", "Existing Publisher").Return(existing, nil)

		err := svc.Create(publisher)

		assert.ErrorIs(t, err, res_err.ErrPublisherConflict)
		repo.AssertNotCalled(t, "Create")
		repo.AssertExpectations(t)
	})
}

func TestPublisherService_Update(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := new(mocks.MockPublisherRepository)
		svc := newPublisherService(repo)
		existing := &model.Publisher{ID: 1, Name: "Old Name"}
		updated := &model.Publisher{Name: "New Name"}

		repo.On("FindByID", uint(1)).Return(existing, nil)
		repo.On("Update", updated).Return(nil)

		err := svc.Update(1, updated)

		assert.NoError(t, err)
		assert.Equal(t, uint(1), updated.ID)
		repo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		repo := new(mocks.MockPublisherRepository)
		svc := newPublisherService(repo)
		updated := &model.Publisher{Name: "New Name"}

		repo.On("FindByID", uint(99)).Return(nil, res_err.ErrPublisherNotFound)

		err := svc.Update(99, updated)

		assert.Error(t, err)
		repo.AssertNotCalled(t, "Update")
		repo.AssertExpectations(t)
	})
}

func TestPublisherService_Delete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := new(mocks.MockPublisherRepository)
		svc := newPublisherService(repo)

		repo.On("Delete", uint(1)).Return(nil)

		err := svc.Delete(1)

		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("repository error", func(t *testing.T) {
		repo := new(mocks.MockPublisherRepository)
		svc := newPublisherService(repo)

		repo.On("Delete", uint(99)).Return(errors.New("db error"))

		err := svc.Delete(99)

		assert.Error(t, err)
		repo.AssertExpectations(t)
	})
}
