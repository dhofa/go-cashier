package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/spf13/viper"

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

func main() {
	// =====================
	// Load config (.env)
	// =====================
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Error reading config file:", err)
	}

	config := Config{
		Port:   viper.GetString("PORT"),
		DBConn: viper.GetString("DB_CONN"),
	}

	if config.Port == "" {
		config.Port = "8080"
	}

	if config.DBConn == "" {
		log.Fatal("DB_CONN is required")
	}

	// =====================
	// Init Database
	// =====================
	db, err := database.InitDB(config.DBConn)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// =====================
	// Dependency Injection
	// =====================
	productRepo := repositories.NewProductRepository(db)
	productService := services.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productService)

	// =====================
	// Router (ServeMux)
	// =====================
	mux := http.NewServeMux()

	// Product routes
	mux.HandleFunc("/api/products", productHandler.HandleProducts)
	mux.HandleFunc("/api/products/", productHandler.HandleProductByID)

	// Default route
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  http.StatusOK,
			"message": "Welcome to Cashier API!",
			"data": models.Response{
				Status:  http.StatusOK,
				Message: "Welcome to Cashier API!",
				Data: map[string]interface{}{
					"app":     "Cashier API",
					"version": 1,
				},
			},
		})
	})

	// =====================
	// HTTP Server
	// =====================
	addr := "0.0.0.0:" + config.Port
	server := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// =====================
	// Run server
	// =====================
	go func() {
		fmt.Println("Server running on", addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server error:", err)
		}
	}()

	// =====================
	// Graceful Shutdown
	// =====================
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	fmt.Println("Server exited gracefully")
}
