package main

import (
	"log"

	"github.com/chillman2101/gits-catalogue/internal/config"
	"github.com/chillman2101/gits-catalogue/internal/model"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, using environment variables")
	}

	cfg := config.Load()
	db, err := cfg.ConnectDB()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	log.Println("seeding database...")

	seedUsers(db)
	seedAuthors(db)
	seedPublishers(db)
	seedBooks(db)

	log.Println("seeding completed successfully")
}

func seedUsers(db *gorm.DB) {
	hashed, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	users := []model.User{
		{Email: "admin@catalogue.com", Password: string(hashed)},
		{Email: "user@catalogue.com", Password: string(hashed)},
	}
	for _, u := range users {
		if err := db.Where("email = ?", u.Email).FirstOrCreate(&u).Error; err != nil {
			log.Printf("failed to seed user %s: %v", u.Email, err)
		}
	}
	log.Println("users seeded")
}

func seedAuthors(db *gorm.DB) {
	authors := []model.Author{
		{Name: "J.K. Rowling", Bio: "British author, best known for the Harry Potter series."},
		{Name: "George R.R. Martin", Bio: "American novelist known for A Song of Ice and Fire."},
		{Name: "Haruki Murakami", Bio: "Japanese author known for surrealist fiction."},
	}
	for _, a := range authors {
		if err := db.Where("name = ?", a.Name).FirstOrCreate(&a).Error; err != nil {
			log.Printf("failed to seed author %s: %v", a.Name, err)
		}
	}
	log.Println("authors seeded")
}

func seedPublishers(db *gorm.DB) {
	publishers := []model.Publisher{
		{Name: "Bloomsbury", Address: "50 Bedford Square, London, UK"},
		{Name: "Bantam Books", Address: "1745 Broadway, New York, USA"},
		{Name: "Kodansha", Address: "2-12-21 Otowa, Bunkyo, Tokyo, Japan"},
	}
	for _, p := range publishers {
		if err := db.Where("name = ?", p.Name).FirstOrCreate(&p).Error; err != nil {
			log.Printf("failed to seed publisher %s: %v", p.Name, err)
		}
	}
	log.Println("publishers seeded")
}

func seedBooks(db *gorm.DB) {
	var (
		rowling, martin, murakami    model.Author
		bloomsbury, bantam, kodansha model.Publisher
	)

	db.Where("name = ?", "J.K. Rowling").First(&rowling)
	db.Where("name = ?", "George R.R. Martin").First(&martin)
	db.Where("name = ?", "Haruki Murakami").First(&murakami)
	db.Where("name = ?", "Bloomsbury").First(&bloomsbury)
	db.Where("name = ?", "Bantam Books").First(&bantam)
	db.Where("name = ?", "Kodansha").First(&kodansha)

	books := []model.Book{
		{
			Title:       "Harry Potter and the Philosopher's Stone",
			ISBN:        "978-0-7475-3269-9",
			Year:        1997,
			AuthorID:    rowling.ID,
			PublisherID: bloomsbury.ID,
		},
		{
			Title:       "Harry Potter and the Chamber of Secrets",
			ISBN:        "978-0-7475-3849-3",
			Year:        1998,
			AuthorID:    rowling.ID,
			PublisherID: bloomsbury.ID,
		},
		{
			Title:       "A Game of Thrones",
			ISBN:        "978-0-553-10354-0",
			Year:        1996,
			AuthorID:    martin.ID,
			PublisherID: bantam.ID,
		},
		{
			Title:       "A Clash of Kings",
			ISBN:        "978-0-553-10803-3",
			Year:        1998,
			AuthorID:    martin.ID,
			PublisherID: bantam.ID,
		},
		{
			Title:       "Norwegian Wood",
			ISBN:        "978-4-06-182462-0",
			Year:        1987,
			AuthorID:    murakami.ID,
			PublisherID: kodansha.ID,
		},
	}

	for _, b := range books {
		if err := db.Where("isbn = ?", b.ISBN).FirstOrCreate(&b).Error; err != nil {
			log.Printf("failed to seed book %s: %v", b.Title, err)
		}
	}
	log.Println("books seeded")
}
