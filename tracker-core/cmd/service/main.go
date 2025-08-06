package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"
	"tracker-core/internal/handlers"
	"tracker-core/internal/services"
	"tracker-core/internal/watcher"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found, using default environment variables")
	}

	dbService, err := services.NewDatabaseService()
	if err != nil {
		log.Fatal("Database initialization error:", err)
	}
	defer dbService.Close()

	db := dbService.GetDB()

	ctx := context.Background()
	if err := dbService.CreateTables(ctx); err != nil {
		log.Fatal("Table creation error:", err)
	}

	currencyService := services.NewCurrencyService(db)
	priceService := services.NewPriceService(db)

	watcher := watcher.NewWatcher(db, currencyService, priceService)
	go watcher.Start()

	currencyHandler := handlers.NewCurrencyHandler(currencyService, priceService)
	
	router := gin.Default()

	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})


	api := router.Group("/currency")
	{
		api.POST("/add", currencyHandler.AddCurrency)
		api.DELETE("/remove", currencyHandler.RemoveCurrency)
		api.GET("/price", currencyHandler.GetPrice)
		api.GET("/list", currencyHandler.GetActiveCurrencies)
	}

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "timestamp": time.Now()})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server running on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Server startup error:", err)
	}
}
