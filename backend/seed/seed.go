package seed

import (
	"log"

	"github.com/Poloni84Learning/ebook-store/config"
	"github.com/Poloni84Learning/ebook-store/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeedAll(db *gorm.DB, cfg *config.Config) {
	seedUsers(db)
	seedBooks(db)
	seedOrders(db)
	log.Println("🌱 Database seeding completed!")
}

func seedUsers(db *gorm.DB) {
	users := []models.User{
		{
			Username:     "admin",
			Email:        "admin@ebookstore.com",
			PasswordHash: hashPassword("Admin@123"),
			Role:         models.RoleAdmin,
			IsActive:     true,
		},
		{
			Username:     "staff1",
			Email:        "staff1@ebookstore.com",
			PasswordHash: hashPassword("Staff@123"),
			Role:         models.RoleStaff,
			IsActive:     true,
		},
		{
			Username:     "customer1",
			Email:        "customer1@example.com",
			PasswordHash: hashPassword("Customer@123"),
			Role:         models.RoleCustomer,
			IsActive:     true,
		},
	}

	for _, user := range users {
		if err := db.FirstOrCreate(&user, "email = ?", user.Email).Error; err != nil {
			log.Printf("Error seeding user %s: %v", user.Email, err)
		}
	}
}

func seedBooks(db *gorm.DB) {
	books := []models.Book{
		{
			Title:       "Clean Code",
			Author:      "Robert Martin",
			Description: "A handbook of agile software craftsmanship",
			Price:       29.99,
			Stock:       100,
			ISBN:        "9780132350884",
			Pages:       464,
			Language:    "English",
			Category:    "Programming",
		},
		{
			Title:       "Design Patterns",
			Author:      "Erich Gamma",
			Description: "Elements of Reusable Object-Oriented Software",
			Price:       49.99,
			Stock:       50,
			ISBN:        "9780201633610",
			Pages:       395,
			Language:    "English",
			Category:    "Programming",
		},
		{
			Title:       "The Pragmatic Programmer",
			Author:      "Andrew Hunt",
			Description: "Your journey to mastery",
			Price:       39.99,
			Stock:       75,
			ISBN:        "9780135957059",
			Pages:       352,
			Language:    "English",
			Category:    "Programming",
		},
	}

	for _, book := range books {
		if err := db.FirstOrCreate(&book, "isbn = ?", book.ISBN).Error; err != nil {
			log.Printf("Error seeding book %s: %v", book.Title, err)
		}
	}
}

func seedOrders(db *gorm.DB) {
	var users []models.User
	var books []models.Book
	db.Find(&users)
	db.Find(&books)

	if len(users) < 1 || len(books) < 2 {
		log.Println("Not enough users or books to seed orders")
		return
	}

	orders := []models.Order{
		{
			UserID:      users[2].ID, // customer1
			TotalAmount: books[0].Price*2 + books[1].Price,
			Status:      "completed",
			OrderItems: []models.OrderItem{
				{
					BookID:   books[0].ID,
					Quantity: 2,
					Price:    books[0].Price,
				},
				{
					BookID:   books[1].ID,
					Quantity: 1,
					Price:    books[1].Price,
				},
			},
		},
		{
			UserID:      users[2].ID, // customer1
			TotalAmount: books[2].Price * 3,
			Status:      "pending",
			OrderItems: []models.OrderItem{
				{
					BookID:   books[2].ID,
					Quantity: 3,
					Price:    books[2].Price,
				},
			},
		},
	}

	for _, order := range orders {
		if err := db.Create(&order).Error; err != nil {
			log.Printf("Error seeding order: %v", err)
		}
	}
}

func hashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("Failed to hash seed password")
	}
	return string(bytes)
}
