package router

import (
	"os"

	"webhook_receiver/internal/handlers"
	"webhook_receiver/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// NewRouter crea y configura el router principal
func NewRouter() *gin.Engine {
	// Cargar variables de entorno
	if err := godotenv.Load(); err != nil {
		// Si no existe el archivo .env, continuar sin error
	}

	// Configurar Gin
	gin.SetMode(getGinMode())

	// Crear router
	router := gin.New()

	// Middleware global
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(corsMiddleware())

	// Obtener secret key de las variables de entorno
	secretKey := os.Getenv("WEBHOOK_SECRET_KEY")
	if secretKey == "" {
		secretKey = "secret_key" // Solo para desarrollo
	}

	// Crear middleware de verificación de firma
	signatureMiddleware := middleware.NewWebhookSignatureMiddleware(secretKey)

	// Crear handlers
	webhookHandler := handlers.NewWebhookHandler()

	// Configurar rutas
	configureRoutes(router, webhookHandler, signatureMiddleware)

	return router
}

// configureRoutes configura todas las rutas de la aplicación
func configureRoutes(router *gin.Engine, webhookHandler *handlers.WebhookHandler, signatureMiddleware *middleware.WebhookSignatureMiddleware) {
	// Grupo de rutas públicas (sin autenticación)
	public := router.Group("/")
	{
		public.GET("/health", webhookHandler.HealthCheck)
	}

	// Grupo de rutas protegidas (con verificación de firma)
	protected := router.Group("/")
	protected.Use(signatureMiddleware.VerifySignature())
	{
		protected.POST("/webhook", webhookHandler.ReceiveWebhook)
	}

	// Ruta de documentación (solo en desarrollo)
	if gin.Mode() != gin.ReleaseMode {
		router.GET("/", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"service": "Webhook Receiver",
				"version": "1.0.0",
				"endpoints": gin.H{
					"health":  "GET /health",
					"webhook": "POST /webhook (requires signature verification)",
				},
			})
		})
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

// corsMiddleware configura CORS
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Webhook-Signature, X-Webhook-Timestamp, X-Webhook-ID, X-Idempotency-Key")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
