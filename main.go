package main

import (
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"

	"wallet-service/config"

	"wallet-service/handlers"
	"wallet-service/storage"
)

func main() {
	config.LoadConfig()

	db, err := storage.InitDB()
	if err != nil {
		log.Fatal(err)
	}

	storage.CreateTables(db)

	handlers.SetDB(db)
	handlers.InitWorkerPool(10)

	router := gin.Default()

	apiV1 := router.Group("/api/v1")
	apiV1.POST("/wallet", handlers.Transaction)
	apiV1.GET("/wallets/:id", handlers.WalletBalance)

	port := config.ServerConfig()

	log.Printf("Server starting on port %s", port)
	router.Run(":" + port)
}
