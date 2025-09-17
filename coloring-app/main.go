// File: main.go
package main



import (
	"database/sql"
	"fmt"
	"log"
	
	

	_ "github.com/lib/pq"
	
	
	
	
	
)

func main() {
	// Load environment variables or use defaults
	dbURL := getEnvOrDefault("DATABASE_URL", "")
	port := getEnvOrDefault("PORT", "8080")

	// Initialize database (optional - remove if not using DB yet)
	var db *sql.DB
	var err error
	if dbURL != "" {
		db, err = sql.Open("postgres", dbURL)
		if err != nil {
			log.Fatal("Failed to connect to database:", err)
		}
		defer db.Close()
		
		if err = db.Ping(); err != nil {
			log.Fatal("Failed to ping database:", err)
		}
		fmt.Println("✅ Database connected successfully")
	}

	StartColoringEcommerceServer(":" + port)
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getenv(key string) string {
	return "" // stub for environment variable, replace with os.Getenv(key) if needed
}