package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"webhook_receiver/internal/dto"

	"github.com/gin-gonic/gin"
)

// WebhookHandler maneja las peticiones de webhooks
type WebhookHandler struct{}

// NewWebhookHandler crea una nueva instancia del handler
func NewWebhookHandler() *WebhookHandler {
	return &WebhookHandler{}
}

// ReceiveWebhook maneja la recepción de webhooks (consumo y facturas)
// @Summary Recibe un webhook
// @Description Endpoint para recibir webhooks con verificación de firma
// @Tags webhooks
// @Accept json
// @Produce json
// @Param X-Webhook-Signature header string true "Firma HMAC del webhook"
// @Param X-Webhook-Timestamp header string true "Timestamp del webhook"
// @Param X-Webhook-ID header string false "ID del webhook"
// @Param X-Idempotency-Key header string false "Clave de idempotencia"
// @Param payload body dto.WebhookPayload true "Payload del webhook"
// @Success 200 {object} dto.WebhookResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /webhook [post]
func (h *WebhookHandler) ReceiveWebhook(c *gin.Context) {
	// Obtener headers para logging
	headers := dto.WebhookHeaders{
		Signature: c.GetHeader("X-Webhook-Signature"),
		WebhookID: c.GetHeader("X-Webhook-ID"),
		Timestamp: c.GetHeader("X-Webhook-Timestamp"),
		IDKey:     c.GetHeader("X-Idempotency-Key"),
	}

	// Primero, detectar el tipo de webhook leyendo solo el campo data_type
	var basePayload struct {
		DataType string `json:"data_type"`
	}

	// Leer el body completo
	bodyBytes, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":   false,
			"message":   "Failed to read request body: " + err.Error(),
			"timestamp": time.Now(),
		})
		return
	}

	// Parsear para obtener el data_type
	if err := json.Unmarshal(bodyBytes, &basePayload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":   false,
			"message":   "Invalid JSON payload: " + err.Error(),
			"timestamp": time.Now(),
		})
		return
	}

	// Procesar según el tipo de webhook
	var processed bool
	var message string

	switch basePayload.DataType {
	case "consumption":
		processed, message = h.processConsumptionWebhookType(bodyBytes, headers)
	case "bills":
		processed, message = h.processBillsWebhookType(bodyBytes, headers)
	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"success":   false,
			"message":   "Unknown data_type: " + basePayload.DataType,
			"timestamp": time.Now(),
		})
		return
	}

	// Log de la recepción del webhook
	c.Header("X-Webhook-Received", "true")

	response := dto.WebhookResponse{
		Success:   true,
		Message:   message,
		Processed: processed,
		Timestamp: time.Now(),
	}

	c.JSON(http.StatusOK, response)
}

// processConsumptionWebhookType procesa webhooks de tipo CONSUMO
func (h *WebhookHandler) processConsumptionWebhookType(bodyBytes []byte, headers dto.WebhookHeaders) (bool, string) {
	var payload dto.WebhookPayload

	// Parsear el payload específico de consumo
	if err := json.Unmarshal(bodyBytes, &payload); err != nil {
		return false, "Failed to parse consumption payload: " + err.Error()
	}

	// Validar campos requeridos
	if payload.GroupBy == "" || payload.SendInterval == "" {
		return false, "Missing required fields for consumption webhook"
	}

	// Aquí implementarías la lógica específica para datos de consumo
	// Por ejemplo:
	// - Guardar datos de consumo en base de datos
	// - Enviar notificaciones a usuarios
	// - Procesar métricas de energía según el tipo de agrupación

	// Log del procesamiento (descomentado para producción)
	// log.Printf("Processing consumption webhook: ID=%d, Contract=%d, GroupBy=%s, Interval=%s",
	//     payload.WebhookID, payload.Data.ContractID, payload.GroupBy, payload.SendInterval)

	message := fmt.Sprintf("Consumption webhook processed successfully for contract %d (%s)",
		payload.Data.ContractID, payload.Data.ContractName)

	return true, message
}

// processBillsWebhookType procesa webhooks de tipo FACTURAS
func (h *WebhookHandler) processBillsWebhookType(bodyBytes []byte, headers dto.WebhookHeaders) (bool, string) {
	var payload dto.BillWebhookPayload

	// Parsear el payload específico de facturas
	if err := json.Unmarshal(bodyBytes, &payload); err != nil {
		return false, "Failed to parse bills payload: " + err.Error()
	}

	// Validar campos requeridos
	if payload.TriggerType == "" {
		return false, "Missing required fields for bills webhook"
	}

	// Aquí implementarías la lógica específica para eventos de facturas
	// Por ejemplo:
	// - Notificar a usuarios sobre nuevas facturas (trigger_type="available")
	// - Procesar confirmación de pagos (trigger_type="paid")
	// - Actualizar estados de facturas en tu sistema

	// Log del procesamiento (descomentado para producción)
	// log.Printf("Processing bills webhook: ID=%d, TriggerType=%s, BillID=%d",
	//     payload.WebhookID, payload.TriggerType, payload.Bill.BillID)

	message := fmt.Sprintf("Bills webhook processed successfully: %s event for bill %d",
		payload.TriggerType, payload.Bill.BillID)

	return true, message
}

// HealthCheck endpoint de salud
// @Summary Health check
// @Description Verifica el estado del servicio
// @Tags health
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /health [get]
func (h *WebhookHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"timestamp": time.Now(),
		"service":   "webhook-receiver",
	})
}
