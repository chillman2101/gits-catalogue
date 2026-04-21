package dto

import (
	"time"

	"github.com/chillman2101/gits-catalogue/internal/model"
)

type CreatePublisherRequest struct {
	Name    string `json:"name" binding:"required,min=2,max=100"`
	Address string `json:"address" binding:"omitempty,max=255"`
}

type UpdatePublisherRequest struct {
	Name    string `json:"name" binding:"required,min=2,max=100"`
	Address string `json:"address" binding:"omitempty,max=255"`
}

func (r *CreatePublisherRequest) ToModel() model.Publisher {
	return model.Publisher{Name: r.Name, Address: r.Address}
}

func (r *UpdatePublisherRequest) ToModel() model.Publisher {
	return model.Publisher{Name: r.Name, Address: r.Address}
}

type BookInPublisher struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	ISBN      string    `json:"isbn"`
	Year      int       `json:"year"`
	AuthorID  uint      `json:"author_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PublisherResponse struct {
	ID        uint              `json:"id"`
	Name      string            `json:"name"`
	Address   string            `json:"address"`
	Books     []BookInPublisher `json:"books"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}

func ToPublisherResponse(p model.Publisher) PublisherResponse {
	books := make([]BookInPublisher, 0, len(p.Books))
	for _, b := range p.Books {
		books = append(books, BookInPublisher{
			ID:        b.ID,
			Title:     b.Title,
			ISBN:      b.ISBN,
			Year:      b.Year,
			AuthorID:  b.AuthorID,
			CreatedAt: b.CreatedAt,
			UpdatedAt: b.UpdatedAt,
		})
	}
	return PublisherResponse{
		ID:        p.ID,
		Name:      p.Name,
		Address:   p.Address,
		Books:     books,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}

func ToPublisherResponseList(publishers []model.Publisher) []PublisherResponse {
	result := make([]PublisherResponse, 0, len(publishers))
	for _, p := range publishers {
		result = append(result, ToPublisherResponse(p))
	}
	return result
}
