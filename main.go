package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"github.com/joho/godotenv"
	"go-cashier/database"
	"go-cashier/handlers"
	"go-cashier/models"
	"go-cashier/repositories"
	"go-cashier/services"
)

type Config struct {
	Port   string `mapstructure:"PORT"`
	DBConn string `mapstructure:"DB_CONN"`
}


func loadConfig() Config {
	// =====================
	// Load .env ONLY if exists (LOCAL)
	// =====================
	if _, err := os.Stat(".env"); err == nil {
		log.Println("📦 Loading .env (local mode)")
		_ = godotenv.Load()
	} else {
		log.Println("🚄 No .env found (Railway / production mode)")
	}

	port := os.Getenv("PORT")
	dbConn := os.Getenv("DB_CONN")

	log.Println("ENV PORT   =", port)
	log.Println("ENV DB_CONN=", dbConn)

	if port == "" {
		port = "8080"
	}

	if dbConn == "" {
		log.Fatal("❌ DB_CONN is required (Railway ENV not injected)")
	}

	return Config{
		Port:   port,
		DBConn: dbConn,
	}
}


func main() {
	// =====================
	// Load Config
	// =====================
	config := loadConfig()

	// =====================
	// Init Database
	// =====================
	db, err := database.InitDB(config.DBConn)
	if err != nil {
		log.Fatal("❌ Failed to initialize database:", err)
	}
	defer db.Close()

	// =====================
	// Migration (Auto-Schema)
	// =====================
	if err := database.Migrate(db); err != nil {
		log.Fatal("❌ Migration failed:", err)
	}

	// =====================
	// Dependency Injection
	// =====================
	// Product
	productRepo := repositories.NewProductRepository(db)
	productService := services.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productService)

	// Category
	categoryRepo := repositories.NewCategoryRepository(db)
	categoryService := services.NewCategoryService(categoryRepo)
	categoryHandler := handlers.NewCategoryHandler(categoryService)

	// Cart
	cartRepo := repositories.NewCartRepository(db)
	cartService := services.NewCartService(cartRepo)
	cartHandler := handlers.NewCartHandler(cartService)

	// Transaction
	transactionRepo := repositories.NewTransactionRepository(db)
	transactionService := services.NewTransactionService(transactionRepo)
	transactionHandler := handlers.NewTransactionHandler(transactionService)

	// =====================
	// Router
	// =====================
	mux := http.NewServeMux()

	// Products
	mux.HandleFunc("/api/products", productHandler.HandleProducts)
	mux.HandleFunc("/api/products/", productHandler.HandleProductByID)

	// Categories
	mux.HandleFunc("/api/categories", categoryHandler.HandleCategories)
	mux.HandleFunc("/api/categories/", categoryHandler.HandleCategoryByID)

	// Carts
	// POST /api/cart/items handles both:
	// 1. Adding item to existing cart (if cart_id provided)
	// 2. Creating new cart + adding item (if cart_id is 0/empty)
	mux.HandleFunc("POST /api/cart/items", cartHandler.AddToCart)
	
	// Helper to view cart
	mux.HandleFunc("GET /api/carts/{id}", cartHandler.GetCart)
	
	// Remove item
	mux.HandleFunc("DELETE /api/carts/{id}/items/{product_id}", cartHandler.RemoveItem)

	// Update item quantity
	mux.HandleFunc("PUT /api/cart/items", cartHandler.UpdateItem)

	// Checkout - accepts array of items in request body
	mux.HandleFunc("POST /api/checkout", transactionHandler.Checkout)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(models.Response{
			Status:  http.StatusOK,
			Message: "Welcome to Cashier API!",
			Data: map[string]interface{}{
				"app":     "Cashier API",
				"version": 1,
			},
		})
	})

	// =====================
	// HTTP Server
	// =====================
	addr := ":" + config.Port
	server := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// =====================
	// Run Server
	// =====================
	go func() {
		log.Println("🚀 Server running on", addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("❌ Server error:", err)
		}
	}()

	// =====================
	// Graceful Shutdown
	// =====================
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("❌ Server forced to shutdown:", err)
	}

	log.Println("✅ Server exited gracefully")
}
