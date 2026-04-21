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

type PublisherService interface {
	GetAll(p query.Params) ([]model.Publisher, int64, error)
	GetByID(id uint) (*model.Publisher, error)
	Create(publisher *model.Publisher) error
	Update(id uint, publisher *model.Publisher) error
	Delete(id uint) error
}

type publisherService struct {
	repo  repository.PublisherRepository
	cache *cache.CacheHelper
}

func NewPublisherService(repo repository.PublisherRepository, cache *cache.CacheHelper) PublisherService {
	return &publisherService{repo, cache}
}

func (s *publisherService) GetAll(p query.Params) ([]model.Publisher, int64, error) {
	key := map[string]string{
		"search": p.Search,
		"page":   strconv.Itoa(p.Page),
		"limit":  strconv.Itoa(p.Limit),
		"order":  p.Order,
		"sort":   p.Sort,
		"module": "publisher",
	}

	cacheData, err := s.cache.GetPublisherTypedCache(key)
	if err == nil {
		return cacheData.Publishers, cacheData.Total, nil
	}

	publishers, total, err := s.repo.FindAll(p)
	if err != nil {
		return nil, 0, err
	}

	if err := s.cache.SetPublisherCache(key, &model.PublisherCacheData{
		Publishers: publishers,
		Total:      total,
	}); err != nil {
		logrus.WithError(err).Error("failed to set publisher cache")
	}

	return publishers, total, nil
}

func (s *publisherService) GetByID(id uint) (*model.Publisher, error) {
	key := map[string]string{
		"id": strconv.Itoa(int(id)),
	}

	cacheData, err := s.cache.GetPublisherTypedCache(key)
	if err == nil {
		return cacheData.Publisher, nil
	}

	publisher, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if err := s.cache.SetPublisherCache(key, &model.PublisherCacheData{
		Publisher: publisher,
	}); err != nil {
		logrus.WithError(err).Error("failed to set publisher cache")
	}

	return publisher, nil
}

func (s *publisherService) Create(publisher *model.Publisher) error {
	existing, err := s.repo.FindByName(publisher.Name)
	if err != nil {
		logrus.WithError(err).WithField("publisher_name", publisher.Name).Error("failed to find publisher")
		return err
	}
	if existing != nil {
		return res_err.ErrPublisherConflict
	}

	if err := s.repo.Create(publisher); err != nil {
		return err
	}

	if err := s.cache.InvalidatePublisherCache(); err != nil {
		logrus.WithError(err).Error("failed to invalidate publisher cache")
	}

	return nil
}

func (s *publisherService) Update(id uint, publisher *model.Publisher) error {
	existing, err := s.repo.FindByID(id)
	if err != nil {
		logrus.WithError(err).WithField("publisher_id", id).Error("failed to find publisher")
		return err
	}
	publisher.ID = existing.ID

	if err := s.repo.Update(publisher); err != nil {
		return err
	}

	if err := s.cache.InvalidatePublisherCache(); err != nil {
		logrus.WithError(err).Error("failed to invalidate publisher cache")
	}
	if err := s.cache.InvalidateBookCache(); err != nil {
		logrus.WithError(err).Error("failed to invalidate book cache after publisher update")
	}

	return nil
}

func (s *publisherService) Delete(id uint) error {
	if err := s.repo.Delete(id); err != nil {
		return err
	}

	if err := s.cache.InvalidatePublisherCache(); err != nil {
		logrus.WithError(err).Error("failed to invalidate publisher cache")
	}
	if err := s.cache.InvalidateBookCache(); err != nil {
		logrus.WithError(err).Error("failed to invalidate book cache after publisher delete")
	}

	return nil
}
