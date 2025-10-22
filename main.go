package main

import (
	"log"
	"os"

	"webhook_receiver/internal/router"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Cargar variables de entorno
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Obtener puerto de las variables de entorno
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Crear router
	router := router.NewRouter()

	// Configurar modo de Gin
	gin.SetMode(getGinMode())

	// Iniciar servidor
	log.Printf("🚀 Webhook Receiver starting on port %s", port)
	log.Printf("📋 Available endpoints:")
	log.Printf("   GET  /health - Health check")
	log.Printf("   POST /webhook - Receive webhooks (requires signature verification)")

	if gin.Mode() != gin.ReleaseMode {
		log.Printf("   GET  / - Service information")
	}

	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

// getGinMode retorna el modo de Gin basado en variables de entorno
func getGinMode() string {
	env := os.Getenv("GIN_MODE")
	if env == "" {
		env = os.Getenv("GO_ENV")
	}

	switch env {
	case "production", "prod":
		return gin.ReleaseMode
	case "test":
		return gin.TestMode
	default:
		return gin.DebugMode
	}
}
