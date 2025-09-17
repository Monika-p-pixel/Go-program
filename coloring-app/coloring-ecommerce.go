// uploadImageHandler handles image uploads from the user's system
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	"io"
	"bytes"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"coloring-app/models"
)





func (ce *ColoringEcommerce) uploadImageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := r.ParseMultipartForm(10 << 20) // 10MB max
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid form data"})
		return
	}
	file, handler, err := r.FormFile("image")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Image not found in request"})
		return
	}
	defer file.Close()
	// Save file to uploads directory
	savePath := fmt.Sprintf("static/uploads/%d_%s", time.Now().UnixNano(), handler.Filename)
	out, err := os.Create(savePath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to save image"})
		return
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to save image"})
		return
	}
	// Return public URL
	imageURL := "/uploads/" + savePath[len("static/uploads/"):]
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"url": imageURL,
	})
}
// checkoutCartHandler clears the user's cart and returns a success message


// addWorksheetHandler handles worksheet creation

func (ce *ColoringEcommerce) checkoutCartHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	tokenStr := r.Header.Get("Authorization")
	if tokenStr == "" || len(tokenStr) < 7 || tokenStr[:7] != "Bearer " {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Missing or invalid authorization header",
		})
		return
	}
	tokenStr = tokenStr[7:]
	claims, err := ce.validateJWT(tokenStr)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid token",
		})
		return
	}
	ce.Cart[claims.UserID] = []int{}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Checkout successful! Thank you for your purchase.",
	})
}
// removeFromCartHandler removes a worksheet from the user's cart
func (ce *ColoringEcommerce) removeFromCartHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	tokenStr := r.Header.Get("Authorization")
	if tokenStr == "" || len(tokenStr) < 7 || tokenStr[:7] != "Bearer " {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Missing or invalid authorization header",
		})
		return
	}
	tokenStr = tokenStr[7:]
	claims, err := ce.validateJWT(tokenStr)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid token",
		})
		return
	}
	var req struct {
		WorksheetID int `json:"worksheet_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid request format",
		})
		return
	}
	worksheetIDs := ce.Cart[claims.UserID]
	newCart := []int{}
	removed := false
	for _, wid := range worksheetIDs {
		if wid == req.WorksheetID && !removed {
			removed = true
			continue
		}
		newCart = append(newCart, wid)
	}
	ce.Cart[claims.UserID] = newCart
	if removed {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "Worksheet removed from cart",
			"worksheet_id": req.WorksheetID,
		})
	} else {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Worksheet not found in cart",
		})
	}
}
// getCartHandler returns the current user's cart with worksheet details
func (ce *ColoringEcommerce) getCartHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	tokenStr := r.Header.Get("Authorization")
	if tokenStr == "" || len(tokenStr) < 7 || tokenStr[:7] != "Bearer " {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Missing or invalid authorization header",
		})
		return
	}
	tokenStr = tokenStr[7:]
	claims, err := ce.validateJWT(tokenStr)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid token",
		})
		return
	}
	worksheetIDs := ce.Cart[claims.UserID]
	var cartWorksheets []models.Worksheet
	for _, wid := range worksheetIDs {
		for _, ws := range ce.Worksheets {
			if ws.ID == wid {
				cartWorksheets = append(cartWorksheets, ws)
				break
			}
		}
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"cart": cartWorksheets,
	})
}


// addWorksheetHandler handles worksheet creation
func (ce *ColoringEcommerce) addWorksheetHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var ws models.Worksheet
	if err := json.NewDecoder(r.Body).Decode(&ws); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Invalid request format",
		})
		return
	}
	ws.ID = len(ce.Worksheets) + 1
	ws.CreatedAt = time.Now()
	ce.Worksheets = append(ce.Worksheets, ws)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"worksheet": ws,
	})
}

// getWorksheetsHandler returns all worksheets
func (ce *ColoringEcommerce) getWorksheetsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ce.Worksheets)
}

// getWorksheetByIDHandler returns worksheet details by ID
func (ce *ColoringEcommerce) getWorksheetByIDHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	idStr := vars["id"]
	for _, ws := range ce.Worksheets {
		if fmt.Sprintf("%d", ws.ID) == idStr {
			json.NewEncoder(w).Encode(ws)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{"error": "Worksheet not found"})
}




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

// LoginRequest represents the login request payload
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse represents the login response
type LoginResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Token   string `json:"token,omitempty"`
	User    User   `json:"user,omitempty"`
}

// RegisterRequest represents the registration request payload
type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
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
// addToCartHandler adds a worksheet to the user's cart
func (ce *ColoringEcommerce) addToCartHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	tokenStr := r.Header.Get("Authorization")
	if tokenStr == "" || len(tokenStr) < 7 || tokenStr[:7] != "Bearer " {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Missing or invalid authorization header",
		})
		return
	}
	tokenStr = tokenStr[7:]
	claims, err := ce.validateJWT(tokenStr)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid token",
		})
		return
	}
	var req struct {
		WorksheetID int `json:"worksheet_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid request format",
		})
		return
	}
	// Check if worksheet exists
	found := false
	for _, ws := range ce.Worksheets {
		if ws.ID == req.WorksheetID {
			found = true
			break
		}
	}
	if !found {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Worksheet not found",
		})
		return
	}
	// Add to cart
	ce.Cart[claims.UserID] = append(ce.Cart[claims.UserID], req.WorksheetID)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Worksheet added to cart",
		"worksheet_id": req.WorksheetID,
	})
}

// Login handles user authentication
func (ce *ColoringEcommerce) Login(email, password string) (*User, string, error) {
	// Find user
	user, exists := ce.Users[email]
	if !exists {
		return nil, "", fmt.Errorf("invalid credentials")
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, "", fmt.Errorf("invalid credentials")
	}

	// Generate JWT token
	token, err := ce.generateJWT(user)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate token")
	}

	return &user, token, nil
}

// Register handles user registration
func (ce *ColoringEcommerce) Register(email, password, name string) (*User, string, error) {
	// Check if user already exists
	if _, exists := ce.Users[email]; exists {
		return nil, "", fmt.Errorf("user already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", fmt.Errorf("failed to process password")
	}

	// Create new user
	newUser := User{
		ID:       len(ce.Users) + 1,
		Email:    email,
		Password: string(hashedPassword),
		Name:     name,
		Role:     "user",
	}

	// Add to users map
	ce.Users[email] = newUser

	// Generate JWT token
	token, err := ce.generateJWT(newUser)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate token")
	}

	return &newUser, token, nil
}

// generateJWT creates a JWT token for the user
func (ce *ColoringEcommerce) generateJWT(user User) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	
	claims := &Claims{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(ce.JWTSecret)
}

// validateJWT validates and parses a JWT token
func (ce *ColoringEcommerce) validateJWT(tokenStr string) (*Claims, error) {
	claims := &Claims{}
	
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return ce.JWTSecret, nil
	})
	
	if err != nil {
		return nil, err
	}
	
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	
	return claims, nil
}

// GetUserByID retrieves a user by ID
func (ce *ColoringEcommerce) GetUserByID(userID int) (*User, error) {
	for _, user := range ce.Users {
		if user.ID == userID {
			return &user, nil
		}
	}
	return nil, fmt.Errorf("user not found")
}

// HTTP Handlers

// loginHandler handles user login HTTP requests
func (ce *ColoringEcommerce) loginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var loginReq LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(LoginResponse{
			Success: false,
			Message: "Invalid request format",
		})
		return
	}

	// Basic validation
	if loginReq.Email == "" || loginReq.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(LoginResponse{
			Success: false,
			Message: "Email and password are required",
		})
		return
	}

	// Use the Login method
	user, token, err := ce.Login(loginReq.Email, loginReq.Password)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(LoginResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	// Return success response
	json.NewEncoder(w).Encode(LoginResponse{
		Success: true,
		Message: "Login successful",
		Token:   token,
		User:    *user,
	})
}

// registerHandler handles user registration HTTP requests
func (ce *ColoringEcommerce) registerHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Debug: print raw request body
	bodyBytes, _ := io.ReadAll(r.Body)
	log.Printf("Register raw body: %s", string(bodyBytes))
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	var registerReq RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&registerReq); err != nil {
		log.Printf("Register decode error: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(LoginResponse{
			Success: false,
			Message: "Invalid request format",
		})
		return
	}

	// Basic validation
	if registerReq.Email == "" || registerReq.Password == "" || registerReq.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(LoginResponse{
			Success: false,
			Message: "All fields are required",
		})
		return
	}

	// Use the Register method
	user, token, err := ce.Register(registerReq.Email, registerReq.Password, registerReq.Name)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "user already exists" {
			statusCode = http.StatusConflict
		}
		
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(LoginResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(LoginResponse{
		Success: true,
		Message: "Registration successful",
		Token:   token,
		User:    *user,
	})
}

// protectedHandler handles protected routes that require authentication
func (ce *ColoringEcommerce) protectedHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	tokenStr := r.Header.Get("Authorization")
	if tokenStr == "" || len(tokenStr) < 7 || tokenStr[:7] != "Bearer " {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Missing or invalid authorization header",
		})
		return
	}

	// Extract token (remove "Bearer " prefix)
	tokenStr = tokenStr[7:]
	
	claims, err := ce.validateJWT(tokenStr)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid token",
		})
		return
	}

	// Get user details
	user, err := ce.GetUserByID(claims.UserID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "User not found",
		})
		return
	}

	// Return user info and all worksheets
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Welcome to Color Fun!",
		"user": map[string]interface{}{
			"id":    user.ID,
			"email": user.Email,
			"name":  user.Name,
			"role":  user.Role,
		},
		"worksheets": ce.Worksheets,
	})
}

// forgotPasswordHandler handles password reset requests
func (ce *ColoringEcommerce) forgotPasswordHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req struct {
		Email string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid request format",
		})
		return
	}

	// Check if user exists
	if _, exists := ce.Users[req.Email]; exists {
		log.Printf("Password reset requested for: %s", req.Email)
		// In production, send actual email with reset token
	}

	// Always return the same response to avoid email enumeration
	json.NewEncoder(w).Encode(map[string]string{
		"message": "If the email exists, a password reset link has been sent",
	})
}

// setupRoutes configures all the HTTP routes
func (ce *ColoringEcommerce) setupRoutes() http.Handler {
	r := mux.NewRouter()
	// Serve add worksheet page
	r.HandleFunc("/add-worksheet", func(w http.ResponseWriter, req *http.Request) {
		http.ServeFile(w, req, "static/add-worksheet.html")
	}).Methods("GET")
	// Serve worksheet details page
	r.HandleFunc("/worksheet/{id}", func(w http.ResponseWriter, req *http.Request) {
		http.ServeFile(w, req, "static/worksheet.html")
	}).Methods("GET")
	// Serve worksheets page
	r.HandleFunc("/worksheets", func(w http.ResponseWriter, req *http.Request) {
		http.ServeFile(w, req, "static/worksheets.html")
	}).Methods("GET")
	// Serve cart page
	r.HandleFunc("/cart", func(w http.ResponseWriter, req *http.Request) {
		http.ServeFile(w, req, "static/cart.html")
	}).Methods("GET")
	// Serve dashboard page
	r.HandleFunc("/dashboard", func(w http.ResponseWriter, req *http.Request) {
		http.ServeFile(w, req, "static/dashboard.html")
	}).Methods("GET")
	// Serve login page
	r.HandleFunc("/login", func(w http.ResponseWriter, req *http.Request) {
		http.ServeFile(w, req, "static/login.html")
	}).Methods("GET")

	// CORS middleware
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}), // In production, specify your frontend domain
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	// API routes
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/login", ce.loginHandler).Methods("POST")
	api.HandleFunc("/register", ce.registerHandler).Methods("POST")
	api.HandleFunc("/forgot-password", ce.forgotPasswordHandler).Methods("POST")
	api.HandleFunc("/dashboard", ce.protectedHandler).Methods("GET")

	// Add worksheet API
	api.HandleFunc("/worksheets", ce.addWorksheetHandler).Methods("POST")
	api.HandleFunc("/worksheets", ce.getWorksheetsHandler).Methods("GET")
	api.HandleFunc("/worksheets/{id}", ce.getWorksheetByIDHandler).Methods("GET")
	// Image upload API
	api.HandleFunc("/upload-image", ce.uploadImageHandler).Methods("POST")
	// Cart API
	api.HandleFunc("/cart", ce.addToCartHandler).Methods("POST")
	api.HandleFunc("/cart", ce.getCartHandler).Methods("GET")
	api.HandleFunc("/cart", ce.removeFromCartHandler).Methods("DELETE")
	api.HandleFunc("/cart/checkout", ce.checkoutCartHandler).Methods("POST")

	// Serve register page
	r.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/register.html")
	}).Methods("GET")

	// Health check endpoint
	api.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status": "healthy",
			"service": "coloring-app-api",
		})
	}).Methods("GET")

	// Serve static files (optional)
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	// Apply CORS middleware
	return corsHandler(r)
}

func StartColoringEcommerceServer(port string) {
	ce := NewColoringEcommerce()
	handler := ce.setupRoutes()

	fmt.Println("🎨 Color Fun API Server starting on", port)
	fmt.Println("📝 Available endpoints:")
	fmt.Println("  POST /api/login")
	fmt.Println("  POST /api/register")
	fmt.Println("  POST /api/forgot-password")
	fmt.Println("  GET  /api/dashboard (protected)")
	fmt.Println("  GET  /api/health")
	fmt.Println("")
	fmt.Println("🔐 Demo credentials:")
	fmt.Println("  Email: demo@colorfun.com")
	fmt.Println("  Password: coloring123")

	log.Fatal(http.ListenAndServe(port, handler))
}