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

func newAuthorService(repo *mocks.MockAuthorRepository) service.AuthorService {
	return service.NewAuthorService(repo, cache.NewCacheHelper(nil))
}

func TestAuthorService_GetAll(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := new(mocks.MockAuthorRepository)
		svc := newAuthorService(repo)
		p := query.Params{Page: 1, Limit: 10}
		authors := []model.Author{{ID: 1, Name: "Author One"}}

		repo.On("FindAll", p).Return(authors, int64(1), nil)

		result, total, err := svc.GetAll(p)

		assert.NoError(t, err)
		assert.Equal(t, int64(1), total)
		assert.Len(t, result, 1)
		assert.Equal(t, "Author One", result[0].Name)
		repo.AssertExpectations(t)
	})

	t.Run("repository error", func(t *testing.T) {
		repo := new(mocks.MockAuthorRepository)
		svc := newAuthorService(repo)
		p := query.Params{Page: 1, Limit: 10}

		repo.On("FindAll", p).Return([]model.Author{}, int64(0), errors.New("db error"))

		_, _, err := svc.GetAll(p)

		assert.Error(t, err)
		repo.AssertExpectations(t)
	})
}

func TestAuthorService_GetByID(t *testing.T) {
	t.Run("found", func(t *testing.T) {
		repo := new(mocks.MockAuthorRepository)
		svc := newAuthorService(repo)
		author := &model.Author{ID: 1, Name: "Author One"}

		repo.On("FindByID", uint(1)).Return(author, nil)

		result, err := svc.GetByID(1)

		assert.NoError(t, err)
		assert.Equal(t, "Author One", result.Name)
		repo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		repo := new(mocks.MockAuthorRepository)
		svc := newAuthorService(repo)

		repo.On("FindByID", uint(99)).Return(nil, res_err.ErrAuthorNotFound)

		result, err := svc.GetByID(99)

		assert.Error(t, err)
		assert.Nil(t, result)
		repo.AssertExpectations(t)
	})
}

func TestAuthorService_Create(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := new(mocks.MockAuthorRepository)
		svc := newAuthorService(repo)
		author := &model.Author{Name: "New Author"}

		repo.On("FindByName", "New Author").Return(nil, nil)
		repo.On("Create", author).Return(nil)

		err := svc.Create(author)

		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("duplicate name", func(t *testing.T) {
		repo := new(mocks.MockAuthorRepository)
		svc := newAuthorService(repo)
		author := &model.Author{Name: "Existing Author"}
		existing := &model.Author{ID: 1, Name: "Existing Author"}

		repo.On("FindByName", "Existing Author").Return(existing, nil)

		err := svc.Create(author)

		assert.ErrorIs(t, err, res_err.ErrAuthorConflict)
		repo.AssertNotCalled(t, "Create")
		repo.AssertExpectations(t)
	})

	t.Run("repository error on create", func(t *testing.T) {
		repo := new(mocks.MockAuthorRepository)
		svc := newAuthorService(repo)
		author := &model.Author{Name: "New Author"}

		repo.On("FindByName", "New Author").Return(nil, nil)
		repo.On("Create", author).Return(errors.New("db error"))

		err := svc.Create(author)

		assert.Error(t, err)
		repo.AssertExpectations(t)
	})
}

func TestAuthorService_Update(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := new(mocks.MockAuthorRepository)
		svc := newAuthorService(repo)
		existing := &model.Author{ID: 1, Name: "Old Name"}
		updated := &model.Author{Name: "New Name"}

		repo.On("FindByID", uint(1)).Return(existing, nil)
		repo.On("Update", updated).Return(nil)

		err := svc.Update(1, updated)

		assert.NoError(t, err)
		assert.Equal(t, uint(1), updated.ID)
		repo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		repo := new(mocks.MockAuthorRepository)
		svc := newAuthorService(repo)
		updated := &model.Author{Name: "New Name"}

		repo.On("FindByID", uint(99)).Return(nil, res_err.ErrAuthorNotFound)

		err := svc.Update(99, updated)

		assert.Error(t, err)
		repo.AssertNotCalled(t, "Update")
		repo.AssertExpectations(t)
	})
}

func TestAuthorService_Delete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := new(mocks.MockAuthorRepository)
		svc := newAuthorService(repo)

		repo.On("Delete", uint(1)).Return(nil)

		err := svc.Delete(1)

		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("repository error", func(t *testing.T) {
		repo := new(mocks.MockAuthorRepository)
		svc := newAuthorService(repo)

		repo.On("Delete", uint(99)).Return(errors.New("db error"))

		err := svc.Delete(99)

		assert.Error(t, err)
		repo.AssertExpectations(t)
	})
}
