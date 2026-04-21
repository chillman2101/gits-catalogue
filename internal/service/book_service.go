package service

import (
	"strconv"

	"github.com/chillman2101/gits-catalogue/internal/model"
	"github.com/chillman2101/gits-catalogue/internal/query"
	"github.com/chillman2101/gits-catalogue/internal/repository"
	"github.com/chillman2101/gits-catalogue/pkg/redis/cache"
	res_err "github.com/chillman2101/gits-catalogue/pkg/response"
	"github.com/sirupsen/logrus"
)

type BookService interface {
	GetAll(p query.Params) ([]model.Book, int64, error)
	GetByID(id uint) (*model.Book, error)
	Create(book *model.Book) error
	Update(id uint, book *model.Book) error
	Delete(id uint) error
}

type bookService struct {
	repo  repository.BookRepository
	cache *cache.CacheHelper
}

func NewBookService(repo repository.BookRepository, cache *cache.CacheHelper) BookService {
	return &bookService{repo, cache}
}

func (s *bookService) GetAll(p query.Params) ([]model.Book, int64, error) {
	key := map[string]string{
		"search": p.Search,
		"page":   strconv.Itoa(p.Page),
		"limit":  strconv.Itoa(p.Limit),
		"order":  p.Order,
		"sort":   p.Sort,
		"module": "book",
	}

	if data, err := s.cache.GetBookTypedCache(key); err == nil {
		return data.Books, data.Total, nil
	}

	books, total, err := s.repo.FindAll(p)
	if err != nil {
		return nil, 0, err
	}

	if err := s.cache.SetBookCache(key, &model.BookCacheData{
		Books: books,
		Total: total,
	}); err != nil {
		logrus.WithError(err).Error("failed to set book cache")
	}

	return books, total, nil
}

func (s *bookService) GetByID(id uint) (*model.Book, error) {
	key := map[string]string{
		"id":     strconv.Itoa(int(id)),
		"module": "book",
	}

	if data, err := s.cache.GetBookTypedCache(key); err == nil {
		return data.Book, nil
	}

	book, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if err := s.cache.SetBookCache(key, &model.BookCacheData{
		Book: book,
	}); err != nil {
		logrus.WithError(err).Error("failed to set book cache")
	}

	return book, nil
}

func (s *bookService) Create(book *model.Book) error {
	existing, err := s.repo.FindByTitle(book.Title)
	if err != nil {
		logrus.WithError(err).WithField("book_title", book.Title).Error("failed to find book")
		return err
	}
	if existing != nil {
		return res_err.ErrBookConflict
	}

	if err := s.repo.Create(book); err != nil {
		logrus.WithError(err).Error("failed to create book")
		return err
	}

	if err := s.cache.InvalidateBookCache(); err != nil {
		logrus.WithError(err).Error("failed to invalidate book cache")
	}

	return nil
}

func (s *bookService) Update(id uint, book *model.Book) error {
	existing, err := s.repo.FindByID(id)
	if err != nil {
		logrus.WithError(err).WithField("book_id", id).Error("failed to find book")
		return err
	}
	book.ID = existing.ID

	if err := s.repo.Update(book); err != nil {
		logrus.WithError(err).Error("failed to update book")
		return err
	}

	if err := s.cache.InvalidateBookCache(); err != nil {
		logrus.WithError(err).Error("failed to invalidate book cache")
	}

	return nil
}

func (s *bookService) Delete(id uint) error {
	if err := s.repo.Delete(id); err != nil {
		logrus.WithError(err).Error("failed to delete book")
		return err
	}

	if err := s.cache.InvalidateBookCache(); err != nil {
		logrus.WithError(err).Error("failed to invalidate book cache")
	}

	return nil
}
