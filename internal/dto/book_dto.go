package dto

import (
	"time"

	"github.com/chillman2101/gits-catalogue/internal/model"
)

type CreateBookRequest struct {
	Title       string `json:"title" binding:"required,min=1,max=255"`
	ISBN        string `json:"isbn" binding:"required,len=17"`
	Year        int    `json:"year" binding:"required,min=1000,max=2100"`
	AuthorID    uint   `json:"author_id" binding:"required,min=1"`
	PublisherID uint   `json:"publisher_id" binding:"required,min=1"`
}

type UpdateBookRequest struct {
	Title       string `json:"title" binding:"required,min=1,max=255"`
	ISBN        string `json:"isbn" binding:"required,len=17"`
	Year        int    `json:"year" binding:"required,min=1000,max=2100"`
	AuthorID    uint   `json:"author_id" binding:"required,min=1"`
	PublisherID uint   `json:"publisher_id" binding:"required,min=1"`
}

func (r *CreateBookRequest) ToModel() model.Book {
	return model.Book{
		Title:       r.Title,
		ISBN:        r.ISBN,
		Year:        r.Year,
		AuthorID:    r.AuthorID,
		PublisherID: r.PublisherID,
	}
}

func (r *UpdateBookRequest) ToModel() model.Book {
	return model.Book{
		Title:       r.Title,
		ISBN:        r.ISBN,
		Year:        r.Year,
		AuthorID:    r.AuthorID,
		PublisherID: r.PublisherID,
	}
}

type AuthorInBook struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Bio  string `json:"bio"`
}

type PublisherInBook struct {
	ID      uint   `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
}

type BookResponse struct {
	ID          uint            `json:"id"`
	Title       string          `json:"title"`
	ISBN        string          `json:"isbn"`
	Year        int             `json:"year"`
	AuthorID    uint            `json:"author_id"`
	Author      AuthorInBook    `json:"author"`
	PublisherID uint            `json:"publisher_id"`
	Publisher   PublisherInBook `json:"publisher"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

func ToBookResponse(b model.Book) BookResponse {
	return BookResponse{
		ID:          b.ID,
		Title:       b.Title,
		ISBN:        b.ISBN,
		Year:        b.Year,
		AuthorID:    b.AuthorID,
		Author:      AuthorInBook{ID: b.Author.ID, Name: b.Author.Name, Bio: b.Author.Bio},
		PublisherID: b.PublisherID,
		Publisher:   PublisherInBook{ID: b.Publisher.ID, Name: b.Publisher.Name, Address: b.Publisher.Address},
		CreatedAt:   b.CreatedAt,
		UpdatedAt:   b.UpdatedAt,
	}
}

func ToBookResponseList(books []model.Book) []BookResponse {
	result := make([]BookResponse, 0, len(books))
	for _, b := range books {
		result = append(result, ToBookResponse(b))
	}
	return result
}
