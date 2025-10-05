// File: models/user.go
package models

import (
	"time"
	
)

type User struct {
	ID        int       `json:"id" db:"id"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"-" db:"password"` // Don't include in JSON responses
	Name      string    `json:"name" db:"name"`
	Role      string    `json:"role" db:"role"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type Worksheet struct {
	ID          int       `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	Difficulty  string    `json:"difficulty" db:"difficulty"`
	Pages       int       `json:"pages" db:"pages"`
	Price       float64   `json:"price" db:"price"`
	ImageURL    string    `json:"image_url" db:"image_url"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	IsActive    bool      `json:"is_active" db:"is_active"`
}

type Order struct {
	ID          int       `json:"id" db:"id"`
	UserID      int       `json:"user_id" db:"user_id"`
	WorksheetID int       `json:"worksheet_id" db:"worksheet_id"`
	Amount      float64   `json:"amount" db:"amount"`
	Status      string    `json:"status" db:"status"` // pending, completed, failed
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}
// File: main.go
package main

import (
	"fmt"
	"net/http"
	"os"
	"github.com/gorilla/handlers"
)

func main() {
	// Load environment variables or use defaults
	port := getEnvOrDefault("PORT", "8080")
	dbURL := getEnvOrDefault("DATABASE_URL", "postgres://user:password@localhost:5432/coloringdb?sslmode=disable")
	jwtSecret := getEnvOrDefault("JWT_SECRET", "supersecretkey")
	
	// Database connection	
	db, err := initDB(dbURL)
	if err != nil {
		fmt.Println("❌ Failed to connect to database:", err)
		return
	}
	defer db.Close()
	fmt.Println("✅ Database connected successfully")
	
	// Initialize auth service
	authService := NewAuthService(jwtSecret)


	// Initialize ColoringEcommerce app
	app := NewColoringEcommerce(db, authService)
	
	// Setup routes
	router := app.SetupRoutes()
	// CORS middleware
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}), // In production, specify your domain
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
		handlers.AllowCredentials(),
	)
	handler := corsHandler(router)
	
	fmt.Println("🚀 Server running on port", port)
	fmt.Println("🔑 JWT Secret:", jwtSecret)
	fmt.Println("🗄️  Database URL:", dbURL
	
	)
	fmt.Println("")
	fmt.Println("Use Ctrl+C to stop the server")

	fmt.Println("")
	fmt.Println("Available endpoints:")
	fmt.Println("POST   /api/register        - Register a new user")
	fmt.Println("POST   /api/login           - Login and receive a JWT")
	fmt.Println("GET    /api/worksheets      - List all worksheets")
	fmt.Println("GET    /api/worksheets/{id} - Get worksheet details")
	fmt.Println("POST   /api/orders          - Create a new order (protected)")
	fmt.Println("GET    /api/orders          - List user orders (protected)")
	fmt.Println("GET    /api/profile         - Get user profile (protected)")
	fmt.Println("PUT    /api/profile         - Update user profile (protected)")
	fmt.Println("Admin-specific endpoints:")
	fmt.Println("POST   /api/worksheets      - Create a new worksheet (admin only)")
	fmt.Println("PUT    /api/worksheets/{id} - Update a worksheet (admin only)")
	fmt.Println("DELETE /api/worksheets/{id} - Delete a worksheet (admin only)")
	fmt.Println("GET    /api/users           - List all users (admin only)")
	fmt.Println("PUT    /api/users/{id}      - Update a user (admin only)")
	fmt.Println("DELETE /api/users/{id}      - Delete a user (admin only)")
	fmt.Println("")
	fmt.Println("Press Ctrl+C to stop the server")
	fmt.Println("================================")

	fmt.Println("🎨 Color Fun Ecommerce Server")
	fmt.Println("================================")
	fmt.Printf("🚀 Server starting on port %s\n", port)
	fmt.Println("📝 Available endpoints:")
	fmt.Println("")
	fmt.Println("📄 Pages:")
	fmt.Println("  GET  /                    - Login page")
	fmt.Println("  GET  /login              - Login page")  
	fmt.Println("  GET  /workshops          - Workshops page")
	fmt.Println("")
	fmt.Println("🔐 Authentication API:")
	fmt.Println("  POST /api/login          - User login")
	fmt.Println("  POST /api/register       - User registration")
	fmt.Println("")
	fmt.Println("🛒 Ecommerce API:"	)
	fmt.Println("  GET  /api/worksheets     - List all worksheets")
	fmt.Println("  GET  /api/worksheets/{id} - Get worksheet details")
	fmt.Println("  GET  /api/dashboard      - User dashboard (protected)")
	fmt.Println("  POST /api/purchase       - Purchase worksheet (protected)")
	fmt.Println("")
	fmt.Println("🔐 Demo credentials:")
	fmt.Println("  Email: mp3593610@gmail.com" )
	fmt.Println("  Password: coloring123")
	fmt.Println("")

}












	














