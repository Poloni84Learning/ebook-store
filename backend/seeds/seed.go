package seeds

import (
	"log"
	"math/rand"
	"time"

	"github.com/Poloni84Learning/ebook-store/config"
	"github.com/Poloni84Learning/ebook-store/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func SeedAll(db *gorm.DB, cfg *config.Config) {
	seedUsers(db)
	seedBooks(db)
	seedOrders(db)
	seedCombos(db)
	seedReviews(db)
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
		{
			Username:     "customer2",
			Email:        "customer2@example.com",
			PasswordHash: hashPassword("Customer@123"),
			Role:         models.RoleCustomer,
			IsActive:     true,
		},
		{
			Username:     "customer3",
			Email:        "customer3@example.com",
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
			CoverImage:  "/uploads/covers/clean-code.jpg",
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
			CoverImage:  "/uploads/covers/design-patterns.jpg",
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
			CoverImage:  "/uploads/covers/the-pragmatic-programer.jpg",
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

func seedCombos(db *gorm.DB) {
	// Lấy dữ liệu sách và user đã tồn tại
	var books []models.Book
	var adminUser models.User
	db.Find(&books)
	db.First(&adminUser, "username = ?", "admin")

	if len(books) < 3 {
		log.Println("Not enough books to seed combos")
		return
	}

	combo := models.BookCombo{
		Title:       "Programming Starter Pack",
		Description: "Essential books for every programmer",
		CreatedBy:   adminUser.ID,
		ComboItems: []models.ComboItem{
			{
				BookID:   books[0].ID, // Clean Code
				IsHidden: false,
			},
			{
				BookID:   books[1].ID, // Design Patterns
				IsHidden: false,
			},
			{
				BookID:   books[2].ID, // The Pragmatic Programmer
				IsHidden: true,        // Ẩn cuốn này trong combo
			},
		},
	}

	// Sử dụng Transaction để đảm bảo toàn vẹn dữ liệu
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&combo).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		log.Printf("Error seeding combo: %v", err)
	}
}

func seedReviews(db *gorm.DB) {
	// Lấy danh sách sách và customer
	var books []models.Book
	var customers []models.User
	db.Find(&books)
	db.Where("role = ?", models.RoleCustomer).Find(&customers)

	if len(books) == 0 || len(customers) < 3 {
		log.Println("Not enough books or customers to seed reviews")
		return
	}

	// Tạo 3 review cho mỗi quyển sách
	for _, book := range books {
		for i := 0; i < 3 && i < len(customers); i++ {
			review := models.Review{
				UserID:    customers[i].ID,
				BookID:    book.ID,
				Rating:    getRandomRating(i),
				Comment:   getRandomComment(book.Title, customers[i].Username),
				ViewCount: i * 10, // ViewCount tăng dần
			}

			if err := db.Create(&review).Error; err != nil {
				log.Printf("Error seeding review for book %s: %v", book.Title, err)
			}
		}
	}
}

// Helper functions
func getRandomRating(index int) int {
	// Rating từ 1-5, nhưng có logic để phân bổ khác nhau
	switch index {
	case 0:
		return 5 // customer1 luôn rating 5 sao
	case 1:
		return 4 // customer2 rating 4 sao
	default:
		return 3 + (index % 3) // customer3 random 3-5 sao
	}
}

func getRandomComment(bookTitle, username string) string {
	comments := []string{
		"Tuyệt vời! " + bookTitle + " đã thay đổi cách tôi viết code",
		username + " recommend cuốn này cho mọi developer",
		"Nội dung hay nhưng giá hơi cao",
		"Giao hàng chậm nhưng sách chất lượng tốt",
		"Đã đọc đi đọc lại nhiều lần",
	}
	return comments[rand.Intn(len(comments))]
}

func hashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("Failed to hash seed password")
	}
	return string(bytes)
}
