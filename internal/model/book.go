package model

import "time"

type Book struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Title       string    `gorm:"not null" json:"title"`
	ISBN        string    `gorm:"uniqueIndex;not null" json:"isbn"`
	Year        int       `json:"year"`
	AuthorID    uint      `gorm:"not null" json:"author_id"`
	Author      Author    `gorm:"foreignKey:AuthorID" json:"author,omitempty"`
	PublisherID uint      `gorm:"not null" json:"publisher_id"`
	Publisher   Publisher `gorm:"foreignKey:PublisherID" json:"publisher,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type BookCacheData struct {
	Books []Book `json:"books"`
	Book  *Book  `json:"book"`
	Total int64  `json:"total"`
}
