package model

import "time"

type Publisher struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	Address   string    `json:"address"`
	Books     []Book    `gorm:"foreignKey:PublisherID" json:"books,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PublisherCacheData struct {
	Publishers []Publisher `json:"publishers"`
	Publisher  *Publisher  `json:"publisher"`
	Total      int64       `json:"total"`
}
