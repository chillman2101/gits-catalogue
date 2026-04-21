package repository

import (
	"github.com/chillman2101/gits-catalogue/internal/model"
	"github.com/chillman2101/gits-catalogue/internal/query"
	res_err "github.com/chillman2101/gits-catalogue/pkg/response"
	"gorm.io/gorm"
)

type BookRepository interface {
	FindAll(p query.Params) ([]model.Book, int64, error)
	FindByID(id uint) (*model.Book, error)
	FindByTitle(title string) (*model.Book, error)
	Create(book *model.Book) error
	Update(book *model.Book) error
	Delete(id uint) error
}

type bookRepository struct {
	db *gorm.DB
}

func NewBookRepository(db *gorm.DB) BookRepository {
	return &bookRepository{db}
}

func (r *bookRepository) FindAll(p query.Params) ([]model.Book, int64, error) {
	var books []model.Book
	var total int64

	q := r.db.Model(&model.Book{})
	if p.Search != "" {
		q = q.Where("title ILIKE ? OR isbn ILIKE ?", "%"+p.Search+"%", "%"+p.Search+"%")
	}

	q.Count(&total)
	err := q.Preload("Author").Preload("Publisher").
		Order(p.OrderClause()).
		Limit(p.Limit).
		Offset(p.Offset()).
		Find(&books).Error
	if err != nil {
		return nil, 0, res_err.ErrBookNotFound
	}

	return books, total, err
}

func (r *bookRepository) FindByID(id uint) (*model.Book, error) {
	var book model.Book
	err := r.db.Preload("Author").Preload("Publisher").First(&book, id).Error
	if err != nil {
		return nil, res_err.ErrBookNotFound
	}
	return &book, nil
}

func (r *bookRepository) FindByTitle(title string) (*model.Book, error) {
	var book model.Book
	err := r.db.Where("title = ?", title).First(&book).Error
	if err != gorm.ErrRecordNotFound {
		return nil, res_err.ErrBookNotFound
	}
	return &book, nil
}

func (r *bookRepository) Create(book *model.Book) error {
	err := r.db.Create(book).Error
	if err != nil {
		return res_err.ErrBookNotFound
	}
	return nil
}

func (r *bookRepository) Update(book *model.Book) error {
	err := r.db.Save(book).Error
	if err != nil {
		return res_err.ErrBookNotFound
	}
	return nil
}

func (r *bookRepository) Delete(id uint) error {
	err := r.db.Delete(&model.Book{}, id).Error
	if err != nil {
		return res_err.ErrBookNotFound
	}
	return nil
}
