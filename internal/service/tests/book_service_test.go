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

func newBookService(repo *mocks.MockBookRepository) service.BookService {
	return service.NewBookService(repo, cache.NewCacheHelper(nil))
}

func TestBookService_GetAll(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := new(mocks.MockBookRepository)
		svc := newBookService(repo)
		p := query.Params{Page: 1, Limit: 10}
		books := []model.Book{{ID: 1, Title: "Go Programming", ISBN: "978-0-13-468599-1"}}

		repo.On("FindAll", p).Return(books, int64(1), nil)

		result, total, err := svc.GetAll(p)

		assert.NoError(t, err)
		assert.Equal(t, int64(1), total)
		assert.Len(t, result, 1)
		assert.Equal(t, "Go Programming", result[0].Title)
		repo.AssertExpectations(t)
	})

	t.Run("repository error", func(t *testing.T) {
		repo := new(mocks.MockBookRepository)
		svc := newBookService(repo)
		p := query.Params{Page: 1, Limit: 10}

		repo.On("FindAll", p).Return([]model.Book{}, int64(0), errors.New("db error"))

		_, _, err := svc.GetAll(p)

		assert.Error(t, err)
		repo.AssertExpectations(t)
	})
}

func TestBookService_GetByID(t *testing.T) {
	t.Run("found", func(t *testing.T) {
		repo := new(mocks.MockBookRepository)
		svc := newBookService(repo)
		book := &model.Book{ID: 1, Title: "Go Programming"}

		repo.On("FindByID", uint(1)).Return(book, nil)

		result, err := svc.GetByID(1)

		assert.NoError(t, err)
		assert.Equal(t, "Go Programming", result.Title)
		repo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		repo := new(mocks.MockBookRepository)
		svc := newBookService(repo)

		repo.On("FindByID", uint(99)).Return(nil, res_err.ErrBookNotFound)

		result, err := svc.GetByID(99)

		assert.Error(t, err)
		assert.Nil(t, result)
		repo.AssertExpectations(t)
	})
}

func TestBookService_Create(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := new(mocks.MockBookRepository)
		svc := newBookService(repo)
		book := &model.Book{Title: "New Book", ISBN: "978-0-00-000000-0", AuthorID: 1, PublisherID: 1}

		repo.On("FindByTitle", "New Book").Return(nil, nil)
		repo.On("Create", book).Return(nil)

		err := svc.Create(book)

		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("duplicate title", func(t *testing.T) {
		repo := new(mocks.MockBookRepository)
		svc := newBookService(repo)
		book := &model.Book{Title: "Existing Book"}
		existing := &model.Book{ID: 1, Title: "Existing Book"}

		repo.On("FindByTitle", "Existing Book").Return(existing, nil)

		err := svc.Create(book)

		assert.ErrorIs(t, err, res_err.ErrBookConflict)
		repo.AssertNotCalled(t, "Create")
		repo.AssertExpectations(t)
	})

	t.Run("repository error on create", func(t *testing.T) {
		repo := new(mocks.MockBookRepository)
		svc := newBookService(repo)
		book := &model.Book{Title: "New Book"}

		repo.On("FindByTitle", "New Book").Return(nil, nil)
		repo.On("Create", book).Return(errors.New("db error"))

		err := svc.Create(book)

		assert.Error(t, err)
		repo.AssertExpectations(t)
	})
}

func TestBookService_Update(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := new(mocks.MockBookRepository)
		svc := newBookService(repo)
		existing := &model.Book{ID: 1, Title: "Old Title"}
		updated := &model.Book{Title: "New Title"}

		repo.On("FindByID", uint(1)).Return(existing, nil)
		repo.On("Update", updated).Return(nil)

		err := svc.Update(1, updated)

		assert.NoError(t, err)
		assert.Equal(t, uint(1), updated.ID)
		repo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		repo := new(mocks.MockBookRepository)
		svc := newBookService(repo)
		updated := &model.Book{Title: "New Title"}

		repo.On("FindByID", uint(99)).Return(nil, res_err.ErrBookNotFound)

		err := svc.Update(99, updated)

		assert.Error(t, err)
		repo.AssertNotCalled(t, "Update")
		repo.AssertExpectations(t)
	})
}

func TestBookService_Delete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := new(mocks.MockBookRepository)
		svc := newBookService(repo)

		repo.On("Delete", uint(1)).Return(nil)

		err := svc.Delete(1)

		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("repository error", func(t *testing.T) {
		repo := new(mocks.MockBookRepository)
		svc := newBookService(repo)

		repo.On("Delete", uint(99)).Return(errors.New("db error"))

		err := svc.Delete(99)

		assert.Error(t, err)
		repo.AssertExpectations(t)
	})
}
