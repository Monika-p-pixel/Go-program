package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
    "time"

    "github.com/golang-jwt/jwt/v5"
    "github.com/gorilla/handlers"
    "github.com/gorilla/mux"
    "golang.org/x/crypto/bcrypt"
    "coloring-app/models"
)

// User represents a user in the system
type User struct {
    ID       int    `json:"id"`
    Email    string `json:"email"`
    Password string `json:"-"` // Don't include in JSON responses
    Name     string `json:"name"`
    Role     string `json:"role"`
}

// ColoringEcommerce is the main application struct
type ColoringEcommerce struct {
    Users      map[string]User      `json:"-"`
    JWTSecret  []byte               `json:"-"`
    Worksheets []models.Worksheet   `json:"-"`
    Cart       map[int][]int        `json:"-"` // userID -> worksheetIDs
}

// JWT Claims structure
type Claims struct {
    UserID int    `json:"user_id"`
    Email  string `json:"email"`
    jwt.RegisteredClaims
}

// NewColoringEcommerce creates a new instance of the application
func NewColoringEcommerce() *ColoringEcommerce {
    // Create demo users with hashed passwords
    hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("coloring123"), bcrypt.DefaultCost)
    users := map[string]User{
        "demo@colorfun.com": {
            ID:       1,
            Email:    "demo@colorfun.com",
            Password: string(hashedPassword),
            Name:     "Demo User",
            Role:     "user",
        },
        "admin@colorfun.com": {
            ID:       2,
            Email:    "admin@colorfun.com",
            Password: string(hashedPassword),
            Name:     "Admin User",
            Role:     "admin",
        },
    }

    return &ColoringEcommerce{
        Users:     users,
        JWTSecret: []byte("your-secret-key-change-in-production"),
        Cart:      make(map[int][]int),
    }
}

// setupRoutes configures all the HTTP routes
func (ce *ColoringEcommerce) setupRoutes() http.Handler {
    r := mux.NewRouter()

    // Serve static pages
    r.HandleFunc("/add-worksheet", func(w http.ResponseWriter, req *http.Request) {
        http.ServeFile(w, req, "static/add-worksheet.html")
    }).Methods("GET")
    r.HandleFunc("/worksheets", func(w http.ResponseWriter, req *http.Request) {
        http.ServeFile(w, req, "static/worksheets.html")
    }).Methods("GET")
    r.HandleFunc("/cart", func(w http.ResponseWriter, req *http.Request) {
        http.ServeFile(w, req, "static/cart.html")
    }).Methods("GET")
    r.HandleFunc("/login", func(w http.ResponseWriter, req *http.Request) {
        http.ServeFile(w, req, "static/login.html")
    }).Methods("GET")
    r.HandleFunc("/register", func(w http.ResponseWriter, req *http.Request) {
        http.ServeFile(w, req, "static/register.html")
    }).Methods("GET")

    // API routes
    api := r.PathPrefix("/api").Subrouter()
    api.HandleFunc("/login", ce.loginHandler).Methods("POST")
    api.HandleFunc("/register", ce.registerHandler).Methods("POST")
    api.HandleFunc("/worksheets", ce.addWorksheetHandler).Methods("POST")
    api.HandleFunc("/worksheets", ce.getWorksheetsHandler).Methods("GET")
    api.HandleFunc("/cart", ce.addToCartHandler).Methods("POST")
    api.HandleFunc("/cart", ce.getCartHandler).Methods("GET")
    api.HandleFunc("/cart/checkout", ce.checkoutCartHandler).Methods("POST")

    // CORS middleware
    corsHandler := handlers.CORS(
        handlers.AllowedOrigins([]string{"*"}), // In production, specify your frontend domain
        handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
        handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
    )

    return corsHandler(r)
}

// StartColoringEcommerceServer starts the server
func StartColoringEcommerceServer(port string) {
    ce := NewColoringEcommerce()
    handler := ce.setupRoutes()
    fmt.Println("🌐 Server running on http://localhost" + port)
    log.Fatal(http.ListenAndServe(port, handler))
}





	