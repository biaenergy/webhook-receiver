package middleware

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// WebhookSignatureMiddleware middleware para verificar la firma de webhooks
type WebhookSignatureMiddleware struct {
	secretKey string
}

// NewWebhookSignatureMiddleware crea una nueva instancia del middleware
func NewWebhookSignatureMiddleware(secretKey string) *WebhookSignatureMiddleware {
	return &WebhookSignatureMiddleware{
		secretKey: secretKey,
	}
}

// VerifySignature verifica la firma del webhook
func (m *WebhookSignatureMiddleware) VerifySignature() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Obtener headers necesarios
		signature := c.GetHeader("X-Webhook-Signature")
		if signature == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "UNAUTHORIZED",
				"message": "Missing X-Webhook-Signature header",
			})
			c.Abort()
			return
		}

		timestamp := c.GetHeader("X-Webhook-Timestamp")
		if timestamp == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "UNAUTHORIZED",
				"message": "Missing X-Webhook-Timestamp header",
			})
			c.Abort()
			return
		}

		// 2. Validar el timestamp (no más de 5 minutos de antigüedad)
		webhookTime, err := time.Parse(time.RFC3339, timestamp)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "UNAUTHORIZED",
				"message": "Invalid timestamp format",
			})
			c.Abort()
			return
		}

		// Verificar que el webhook no sea muy antiguo (5 minutos)
		if time.Since(webhookTime) > 50*time.Hour {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "UNAUTHORIZED",
				"message": "Webhook timestamp too old",
			})
			c.Abort()
			return
		}

		// 3. Leer el payload completo
		payload, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "BAD_REQUEST",
				"message": "Failed to read request body",
			})
			c.Abort()
			return
		}

		// 4. Verificar la firma
		isValid, err := m.verifySignature(m.secretKey, payload, signature)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "INTERNAL_ERROR",
				"message": "Failed to verify signature",
			})
			c.Abort()
			return
		}

		if !isValid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "UNAUTHORIZED",
				"message": "Invalid signature",
			})
			c.Abort()
			return
		}

		// 5. Restaurar el body para que el handler pueda leerlo
		c.Request.Body = io.NopCloser(bytes.NewReader(payload))

		// Continuar con el siguiente handler
		c.Next()
	}
}

// verifySignature verifica la firma HMAC del payload
func (m *WebhookSignatureMiddleware) verifySignature(secretKey string, payload []byte, receivedSignature string) (bool, error) {
	// Generar la firma esperada
	mac := hmac.New(sha256.New, []byte(secretKey))
	mac.Write(payload)
	expectedSignature := hex.EncodeToString(mac.Sum(nil))

	// Comparar las firmas usando una función de tiempo constante
	return hmac.Equal(
		[]byte(receivedSignature),
		[]byte(expectedSignature),
	), nil
}
