package model

import "time"

type Author struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	Bio       string    `json:"bio"`
	Books     []Book    `gorm:"foreignKey:AuthorID" json:"books,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AuthorCacheData struct {
	Authors []Author `json:"authors"`
	Author  *Author  `json:"author"`
	Total   int64    `json:"total"`
}
