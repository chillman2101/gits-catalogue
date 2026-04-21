package dto

import (
	"time"

	"github.com/chillman2101/gits-catalogue/internal/model"
)

type CreateAuthorRequest struct {
	Name string `json:"name" binding:"required,min=2,max=100"`
	Bio  string `json:"bio" binding:"omitempty,max=500"`
}

type UpdateAuthorRequest struct {
	Name string `json:"name" binding:"required,min=2,max=100"`
	Bio  string `json:"bio" binding:"omitempty,max=500"`
}

func (r *CreateAuthorRequest) ToModel() model.Author {
	return model.Author{Name: r.Name, Bio: r.Bio}
}

func (r *UpdateAuthorRequest) ToModel() model.Author {
	return model.Author{Name: r.Name, Bio: r.Bio}
}

type BookInAuthor struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	ISBN        string    `json:"isbn"`
	Year        int       `json:"year"`
	PublisherID uint      `json:"publisher_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type AuthorResponse struct {
	ID        uint           `json:"id"`
	Name      string         `json:"name"`
	Bio       string         `json:"bio"`
	Books     []BookInAuthor `json:"books"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

func ToAuthorResponse(a model.Author) AuthorResponse {
	books := make([]BookInAuthor, 0, len(a.Books))
	for _, b := range a.Books {
		books = append(books, BookInAuthor{
			ID:          b.ID,
			Title:       b.Title,
			ISBN:        b.ISBN,
			Year:        b.Year,
			PublisherID: b.PublisherID,
			CreatedAt:   b.CreatedAt,
			UpdatedAt:   b.UpdatedAt,
		})
	}
	return AuthorResponse{
		ID:        a.ID,
		Name:      a.Name,
		Bio:       a.Bio,
		Books:     books,
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
	}
}

func ToAuthorResponseList(authors []model.Author) []AuthorResponse {
	result := make([]AuthorResponse, 0, len(authors))
	for _, a := range authors {
		result = append(result, ToAuthorResponse(a))
	}
	return result
}
