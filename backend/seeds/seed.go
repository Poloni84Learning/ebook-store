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

func SeedAll(db *gorm.DB, cfg *config.Config) {
	// Kiểm tra xem đã có dữ liệu chưa
	var count int64
	db.Model(&models.User{}).Count(&count)

	if count == 0 {
		seedUsers(db)
		seedBooks(db)
		seedOrders(db)
		seedCombos(db)
		seedReviews(db)
		log.Println("🌱 Database seeding completed!")
	}
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
		{
			Title:       "You Don't Know JS Yet",
			Author:      "Kyle Simpson",
			Description: "A deep dive into JavaScript core concepts",
			Price:       29.99,
			Stock:       50,
			CoverImage:  "/uploads/covers/you-dont-know-js.jpg",
			ISBN:        "9781091210092",
			Pages:       278,
			Language:    "English",
			Category:    "Programming",
		},
		{
			Title:       "Introduction to Probability",
			Author:      "Dimitri P. Bertsekas",
			Description: "A comprehensive guide to probability theory",
			Price:       45.00,
			Stock:       40,
			CoverImage:  "/uploads/covers/intro-to-probability.jpg",
			ISBN:        "9781886529236",
			Pages:       544,
			Language:    "English",
			Category:    "Mathematics",
		},
		{
			Title:       "How Not to Be Wrong",
			Author:      "Jordan Ellenberg",
			Description: "The power of mathematical thinking in everyday life",
			Price:       18.99,
			Stock:       60,
			CoverImage:  "/uploads/covers/how-not-to-be-wrong.jpg",
			ISBN:        "9780143127536",
			Pages:       480,
			Language:    "English",
			Category:    "Mathematics",
		},
		{
			Title:       "Calculus, 10th Edition",
			Author:      "Ron Larson",
			Description: "An in-depth textbook for learning calculus",
			Price:       79.95,
			Stock:       25,
			CoverImage:  "/uploads/covers/calculus-larson.jpg",
			ISBN:        "9781337624183",
			Pages:       1280,
			Language:    "English",
			Category:    "Mathematics",
		},
		{
			Title:       "Harry Potter and the Sorcerer's Stone",
			Author:      "J.K. Rowling",
			Description: "The beginning of the magical journey",
			Price:       24.99,
			Stock:       120,
			CoverImage:  "/uploads/covers/hp1.jpg",
			ISBN:        "9780590353427",
			Pages:       320,
			Language:    "English",
			Category:    "Fantasy",
		},
		{
			Title:       "Harry Potter and the Chamber of Secrets",
			Author:      "J.K. Rowling",
			Description: "The second year at Hogwarts",
			Price:       24.99,
			Stock:       100,
			CoverImage:  "/uploads/covers/hp2.jpg",
			ISBN:        "9780439064873",
			Pages:       341,
			Language:    "English",
			Category:    "Fantasy",
		},
		{
			Title:       "Harry Potter and the Prisoner of Azkaban",
			Author:      "J.K. Rowling",
			Description: "A dark past returns to haunt Harry",
			Price:       26.99,
			Stock:       95,
			CoverImage:  "/uploads/covers/hp3.jpg",
			ISBN:        "9780439136365",
			Pages:       435,
			Language:    "English",
			Category:    "Fantasy",
		},
		{
			Title:       "Harry Potter and the Goblet of Fire",
			Author:      "J.K. Rowling",
			Description: "The Triwizard Tournament begins",
			Price:       28.99,
			Stock:       90,
			CoverImage:  "/uploads/covers/hp4.jpg",
			ISBN:        "9780439139601",
			Pages:       734,
			Language:    "English",
			Category:    "Fantasy",
		},
		{
			Title:       "Harry Potter and the Order of the Phoenix",
			Author:      "J.K. Rowling",
			Description: "Rebellion brews as Voldemort returns",
			Price:       29.99,
			Stock:       85,
			CoverImage:  "/uploads/covers/hp5.jpg",
			ISBN:        "9780439358071",
			Pages:       870,
			Language:    "English",
			Category:    "Fantasy",
		},
		{
			Title:       "Harry Potter and the Half-Blood Prince",
			Author:      "J.K. Rowling",
			Description: "Secrets of Voldemort's past are revealed",
			Price:       29.99,
			Stock:       80,
			CoverImage:  "/uploads/covers/hp6.jpg",
			ISBN:        "9780439785969",
			Pages:       652,
			Language:    "English",
			Category:    "Fantasy",
		},
		{
			Title:       "Harry Potter and the Deathly Hallows",
			Author:      "J.K. Rowling",
			Description: "The final battle begins",
			Price:       32.99,
			Stock:       100,
			CoverImage:  "/uploads/covers/hp7.jpg",
			ISBN:        "9780545010221",
			Pages:       759,
			Language:    "English",
			Category:    "Fantasy",
		},
		{
			Title:       "Twilight",
			Author:      "Stephenie Meyer",
			Description: "A love story between a human and a vampire",
			Price:       22.99,
			Stock:       80,
			CoverImage:  "/uploads/covers/twilight.jpg",
			ISBN:        "9780316015844",
			Pages:       544,
			Language:    "English",
			Category:    "Fantasy",
		},
		{
			Title:       "New Moon",
			Author:      "Stephenie Meyer",
			Description: "Bella faces new heartbreaks and discoveries",
			Price:       22.99,
			Stock:       75,
			CoverImage:  "/uploads/covers/new-moon.jpg",
			ISBN:        "9780316024969",
			Pages:       608,
			Language:    "English",
			Category:    "Fantasy",
		},
		{
			Title:       "Eclipse",
			Author:      "Stephenie Meyer",
			Description: "A choice between love and friendship",
			Price:       23.99,
			Stock:       70,
			CoverImage:  "/uploads/covers/eclipse.jpg",
			ISBN:        "9780316160209",
			Pages:       640,
			Language:    "English",
			Category:    "Fantasy",
		},
		{
			Title:       "Breaking Dawn",
			Author:      "Stephenie Meyer",
			Description: "Love, sacrifice, and transformation",
			Price:       25.99,
			Stock:       60,
			CoverImage:  "/uploads/covers/breaking-dawn.jpg",
			ISBN:        "9780316067928",
			Pages:       768,
			Language:    "English",
			Category:    "Fantasy",
		},
		{
			Title:       "Midnight Sun",
			Author:      "Stephenie Meyer",
			Description: "Edward's perspective of Twilight",
			Price:       27.99,
			Stock:       50,
			CoverImage:  "/uploads/covers/midnight-sun.jpg",
			ISBN:        "9780316707046",
			Pages:       672,
			Language:    "English",
			Category:    "Fantasy",
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
			Status:      "completed",
			OrderItems: []models.OrderItem{
				{
					BookID:   books[2].ID,
					Quantity: 3,
					Price:    books[2].Price,
				},
			},
		},
		{
			UserID:      users[3].ID, // customer2
			TotalAmount: books[7].Price + books[8].Price + books[9].Price + books[10].Price,
			Status:      "completed",
			OrderItems: []models.OrderItem{
				{
					BookID:   books[7].ID, // HP 1
					Quantity: 1,
					Price:    books[7].Price,
				},
				{
					BookID:   books[8].ID, // HP 2
					Quantity: 1,
					Price:    books[8].Price,
				},
				{
					BookID:   books[9].ID, // HP 3
					Quantity: 1,
					Price:    books[9].Price,
				},
				{
					BookID:   books[10].ID, // HP 4
					Quantity: 1,
					Price:    books[10].Price,
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
	rand.Seed(time.Now().UnixNano())

	booksToUpdate := make(map[uint]bool)

	// Tạo 3 review cho mỗi quyển sách
	for _, book := range books {

		numReviewers := rand.Intn(min(3, len(customers))) + 1 // ít nhất 1 người review
		randomIndexes := rand.Perm(len(customers))[:numReviewers]

		for idx := range randomIndexes {
			customer := customers[idx]

			review := models.Review{
				UserID:  customer.ID,
				BookID:  book.ID,
				Rating:  getRandomRating(),
				Comment: getRandomComment(book.Title, customer.Username),
			}

			if err := db.Create(&review).Error; err != nil {
				log.Printf("Error seeding review for book %s: %v", book.Title, err)
			} else {
				// Đánh dấu sách cần cập nhật average_rating
				booksToUpdate[book.ID] = true
			}
		}
	}
	for bookID := range booksToUpdate {
		updateBookAverageRating(db, bookID)
	}
}

func updateBookAverageRating(db *gorm.DB, bookID uint) {
	var avgRating float64
	err := db.Model(&models.Review{}).
		Where("book_id = ?", bookID).
		Select("COALESCE(AVG(rating), 0)").
		Row().
		Scan(&avgRating)

	if err != nil {
		log.Printf("Error calculating average rating for book %d: %v", bookID, err)
		return
	}

	err = db.Model(&models.Book{}).
		Where("id = ?", bookID).
		Update("average_rating", avgRating).Error

	if err != nil {
		log.Printf("Error updating average rating for book %d: %v", bookID, err)
	} else {
		log.Printf("Updated average rating for book %d to %.2f", bookID, avgRating)
	}
}

// Helper functions
func getRandomRating() int {
	// Khởi tạo nguồn ngẫu nhiên
	return rand.Intn(3) + 3 // Giá trị ngẫu nhiên trong phạm vi [3, 5]
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
