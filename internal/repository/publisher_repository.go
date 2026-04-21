package repository

import (
	"github.com/chillman2101/gits-catalogue/internal/model"
	"github.com/chillman2101/gits-catalogue/internal/query"
	res_err "github.com/chillman2101/gits-catalogue/pkg/response"
	"gorm.io/gorm"
)

type PublisherRepository interface {
	FindAll(p query.Params) ([]model.Publisher, int64, error)
	FindByID(id uint) (*model.Publisher, error)
	FindByName(name string) (*model.Publisher, error)
	Create(publisher *model.Publisher) error
	Update(publisher *model.Publisher) error
	Delete(id uint) error
}

type publisherRepository struct {
	db *gorm.DB
}

func NewPublisherRepository(db *gorm.DB) PublisherRepository {
	return &publisherRepository{db}
}

func (r *publisherRepository) FindAll(p query.Params) ([]model.Publisher, int64, error) {
	var publishers []model.Publisher
	var total int64

	q := r.db.Model(&model.Publisher{})
	if p.Search != "" {
		q = q.Where("name ILIKE ?", "%"+p.Search+"%")
	}

	q.Count(&total)
	err := q.Preload("Books").
		Order(p.OrderClause()).
		Limit(p.Limit).
		Offset(p.Offset()).
		Find(&publishers).Error
	if err != nil {
		return nil, 0, res_err.ErrPublisherNotFound
	}

	return publishers, total, nil
}

func (r *publisherRepository) FindByID(id uint) (*model.Publisher, error) {
	var publisher model.Publisher
	err := r.db.Preload("Books").First(&publisher, id).Error
	if err != nil {
		return nil, res_err.ErrPublisherNotFound
	}
	return &publisher, nil
}

func (r *publisherRepository) FindByName(name string) (*model.Publisher, error) {
	var publisher model.Publisher
	err := r.db.Where("name = ?", name).First(&publisher).Error
	if err != gorm.ErrRecordNotFound {
		return nil, res_err.ErrPublisherNotFound
	}
	return &publisher, nil
}

func (r *publisherRepository) Create(publisher *model.Publisher) error {
	err := r.db.Create(publisher).Error
	if err != nil {
		return res_err.ErrPublisherNotFound
	}
	return nil
}

func (r *publisherRepository) Update(publisher *model.Publisher) error {
	err := r.db.Save(publisher).Error
	if err != nil {
		return res_err.ErrPublisherNotFound
	}
	return nil
}

func (r *publisherRepository) Delete(id uint) error {
	err := r.db.Delete(&model.Publisher{}, id).Error
	if err != nil {
		return res_err.ErrPublisherNotFound
	}
	return nil
}
