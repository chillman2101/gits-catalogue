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

type AuthorService interface {
	GetAll(p query.Params) ([]model.Author, int64, error)
	GetByID(id uint) (*model.Author, error)
	Create(author *model.Author) error
	Update(id uint, author *model.Author) error
	Delete(id uint) error
}

type authorService struct {
	repo  repository.AuthorRepository
	cache *cache.CacheHelper
}

func NewAuthorService(repo repository.AuthorRepository, cache *cache.CacheHelper) AuthorService {
	return &authorService{repo, cache}
}

func (s *authorService) GetAll(p query.Params) ([]model.Author, int64, error) {
	key := map[string]string{
		"search": p.Search,
		"page":   strconv.Itoa(p.Page),
		"limit":  strconv.Itoa(p.Limit),
		"order":  p.Order,
		"sort":   p.Sort,
		"module": "author",
	}

	// cache miss (redis.Nil) is normal — fall through to DB
	if data, err := s.cache.GetAuthorTypedCache(key); err == nil {
		return data.Authors, data.Total, nil
	}

	authors, total, err := s.repo.FindAll(p)
	if err != nil {
		logrus.WithError(err).Error("failed to find authors")
		return nil, 0, err
	}

	if err := s.cache.SetAuthorCache(key, &model.AuthorCacheData{
		Authors: authors,
		Total:   total,
	}); err != nil {
		logrus.WithError(err).Error("failed to set author cache")
	}

	return authors, total, nil
}

func (s *authorService) GetByID(id uint) (*model.Author, error) {
	key := map[string]string{
		"id":     strconv.Itoa(int(id)),
		"module": "author",
	}

	if data, err := s.cache.GetAuthorTypedCache(key); err == nil {
		return data.Author, nil
	}

	author, err := s.repo.FindByID(id)
	if err != nil {
		logrus.WithError(err).Error("failed to find author")
		return nil, err
	}

	if err := s.cache.SetAuthorCache(key, &model.AuthorCacheData{
		Author: author,
	}); err != nil {
		logrus.WithError(err).Error("failed to set author cache")
	}

	return author, nil
}

func (s *authorService) Create(author *model.Author) error {
	existing, err := s.repo.FindByName(author.Name)
	if err != nil {
		logrus.WithError(err).WithField("author_name", author.Name).Error("failed to find author")
		return err
	}
	if existing != nil {
		return res_err.ErrAuthorConflict
	}

	if err := s.repo.Create(author); err != nil {
		logrus.WithError(err).Error("failed to create author")
		return err
	}

	if err := s.cache.InvalidateAuthorCache(); err != nil {
		logrus.WithError(err).Error("failed to invalidate author cache")
	}

	return nil
}

func (s *authorService) Update(id uint, author *model.Author) error {
	existing, err := s.repo.FindByID(id)
	if err != nil {
		logrus.WithError(err).WithField("author_id", id).Error("failed to find author")
		return err
	}
	author.ID = existing.ID

	if err := s.repo.Update(author); err != nil {
		logrus.WithError(err).Error("failed to update author")
		return err
	}

	// invalidate book cache juga karena book preload author
	if err := s.cache.InvalidateAuthorCache(); err != nil {
		logrus.WithError(err).Error("failed to invalidate author cache")
	}
	if err := s.cache.InvalidateBookCache(); err != nil {
		logrus.WithError(err).Error("failed to invalidate book cache after author update")
	}

	return nil
}

func (s *authorService) Delete(id uint) error {
	if err := s.repo.Delete(id); err != nil {
		logrus.WithError(err).Error("failed to delete author")
		return err
	}

	if err := s.cache.InvalidateAuthorCache(); err != nil {
		logrus.WithError(err).Error("failed to invalidate author cache")
	}
	if err := s.cache.InvalidateBookCache(); err != nil {
		logrus.WithError(err).Error("failed to invalidate book cache after author delete")
	}

	return nil
}
