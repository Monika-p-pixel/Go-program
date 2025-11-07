package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type ColoringSheet struct {
	ID         int     `json:"id"`
	Name       string  `json:"name"`
	Image      string  `json:"image"`
	Quantity   int     `json:"quantity"`
	Price      float64 `json:"price"`
	UploadedAt string  `json:"uploaded_at"`
}

type CartItem struct {
	SheetID  int `json:"sheet_id"`
	Quantity int `json:"quantity"`
}

type Order struct {
	Items        []CartItem `json:"items"`
	CustomerName string     `json:"customer_name"`
	Email        string     `json:"email"`
	Phone        string     `json:"phone"`
	Address      string     `json:"address"`
	City         string     `json:"city"`
	State        string     `json:"state"`
	Pincode      string     `json:"pincode"`
	TotalAmount  float64    `json:"total_amount"`
	OrderDate    string     `json:"order_date"`
}

var (
	sheets    = make(map[int]*ColoringSheet)
	sheetID   = 1
	sheetLock sync.RWMutex
)

func main() {
	// Serve static files
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	// API endpoints
	http.HandleFunc("/api/sheets", sheetsHandler)
	http.HandleFunc("/api/sheets/add", addSheetHandler)
	http.HandleFunc("/api/sheets/download/", downloadHandler)
	http.HandleFunc("/api/order", orderHandler)

	fmt.Println("✓ Server starting on http://localhost:8080")
	fmt.Println("✓ Make sure index.html is in the same directory")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func sheetsHandler(w http.ResponseWriter, r *http.Request) {
	// Enable CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		return
	}

	sheetLock.RLock()
	defer sheetLock.RUnlock()

	sheetList := make([]*ColoringSheet, 0, len(sheets))
	for _, sheet := range sheets {
		sheetList = append(sheetList, sheet)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sheetList)
}

func addSheetHandler(w http.ResponseWriter, r *http.Request) {
	// Enable CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	quantity, _ := strconv.Atoi(r.FormValue("quantity"))
	price, _ := strconv.ParseFloat(r.FormValue("price"), 64)

	file, _, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Unable to read image", http.StatusBadRequest)
		return
	}
	defer file.Close()

	imageData, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Unable to read image data", http.StatusInternalServerError)
		return
	}

	base64Image := base64.StdEncoding.EncodeToString(imageData)

	sheetLock.Lock()
	sheet := &ColoringSheet{
		ID:         sheetID,
		Name:       name,
		Image:      base64Image,
		Quantity:   quantity,
		Price:      price,
		UploadedAt: time.Now().Format("2006-01-02 15:04:05"),
	}
	sheets[sheetID] = sheet
	sheetID++
	sheetLock.Unlock()

	log.Printf("✓ Added new sheet: %s (ID: %d, Price: ₹%.2f)", name, sheet.ID, price)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sheet)
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	// Enable CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")

	idStr := r.URL.Path[len("/api/sheets/download/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	sheetLock.RLock()
	sheet, exists := sheets[id]
	sheetLock.RUnlock()

	if !exists {
		http.Error(w, "Sheet not found", http.StatusNotFound)
		return
	}

	imageData, err := base64.StdEncoding.DecodeString(sheet.Image)
	if err != nil {
		http.Error(w, "Unable to decode image", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.png", sheet.Name))
	w.Write(imageData)

	log.Printf("✓ Downloaded sheet: %s (ID: %d)", sheet.Name, id)
}

func orderHandler(w http.ResponseWriter, r *http.Request) {
	// Enable CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var order Order
	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	order.OrderDate = time.Now().Format("2006-01-02 15:04:05")
	orderID := fmt.Sprintf("ORD-%d", time.Now().Unix())

	// Log order details
	log.Printf("✓ New Order Received: %s", orderID)
	log.Printf("  Customer: %s", order.CustomerName)
	log.Printf("  Email: %s", order.Email)
	log.Printf("  Phone: %s", order.Phone)
	log.Printf("  Address: %s, %s, %s - %s", order.Address, order.City, order.State, order.Pincode)
	log.Printf("  Total Amount: ₹%.2f", order.TotalAmount)
	log.Printf("  Items: %d", len(order.Items))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":   "success",
		"message":  "Order placed successfully!",
		"order_id": orderID,
	})
}
