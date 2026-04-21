package repository

import (
	res_err "github.com/chillman2101/gits-catalogue/pkg/response"
	"github.com/sirupsen/logrus"

	"github.com/chillman2101/gits-catalogue/internal/model"
	"github.com/chillman2101/gits-catalogue/internal/query"
	"gorm.io/gorm"
)

type AuthorRepository interface {
	FindAll(p query.Params) ([]model.Author, int64, error)
	FindByID(id uint) (*model.Author, error)
	FindByName(name string) (*model.Author, error)
	Create(author *model.Author) error
	Update(author *model.Author) error
	Delete(id uint) error
}

type authorRepository struct {
	db *gorm.DB
}

func NewAuthorRepository(db *gorm.DB) AuthorRepository {
	return &authorRepository{db}
}

func (r *authorRepository) FindAll(p query.Params) ([]model.Author, int64, error) {
	var authors []model.Author
	var total int64

	q := r.db.Model(&model.Author{})
	if p.Search != "" {
		q = q.Where("name ILIKE ?", "%"+p.Search+"%")
	}

	q.Count(&total)
	err := q.Preload("Books").
		Order(p.OrderClause()).
		Limit(p.Limit).
		Offset(p.Offset()).
		Find(&authors).Error
	if err != nil {
		return nil, 0, res_err.ErrAuthorNotFound
	}

	return authors, total, nil
}

func (r *authorRepository) FindByID(id uint) (*model.Author, error) {
	var author model.Author
	err := r.db.Preload("Books").First(&author, id).Error
	if err != nil {
		return nil, res_err.ErrAuthorNotFound
	}
	return &author, nil
}

func (r *authorRepository) FindByName(name string) (*model.Author, error) {
	var author model.Author
	err := r.db.Where("name = ?", name).First(&author).Error
	if err != gorm.ErrRecordNotFound {
		logrus.Error(err)
		return nil, res_err.ErrAuthorNotFound
	}
	return &author, nil
}

func (r *authorRepository) Create(author *model.Author) error {
	err := r.db.Create(author).Error
	if err != nil {
		return res_err.ErrAuthorNotFound
	}
	return nil
}

func (r *authorRepository) Update(author *model.Author) error {
	err := r.db.Save(author).Error
	if err != nil {
		return res_err.ErrAuthorNotFound
	}
	return nil
}

func (r *authorRepository) Delete(id uint) error {
	err := r.db.Delete(&model.Author{}, id).Error
	if err != nil {
		return res_err.ErrAuthorNotFound
	}
	return nil
}
