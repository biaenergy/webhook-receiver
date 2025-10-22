package handlers

import (
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

// ReceiveWebhook maneja la recepción de webhooks
// @Summary Recibe un webhook
// @Description Endpoint para recibir webhooks con verificación de firma
// @Tags webhooks
// @Accept json
// @Produce json
// @Param X-Webhook-Signature header string true "Firma HMAC del webhook"
// @Param X-Webhook-Timestamp header string true "Timestamp del webhook"
// @Param X-Webhook-ID header string false "ID del webhook"
// @Param X-Idempotency-Key header string false "Clave de idempotencia"
// @Param payload body dto.WebhookReceivedPayload true "Payload del webhook"
// @Success 200 {object} dto.WebhookReceivedResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /webhook [post]
func (h *WebhookHandler) ReceiveWebhook(c *gin.Context) {
	var payload dto.WebhookReceivedPayload

	// Parsear el JSON del body
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":   false,
			"message":   "Invalid JSON payload: " + err.Error(),
			"timestamp": time.Now(),
		})
		return
	}

	// Obtener headers para logging
	headers := dto.WebhookReceivedHeaders{
		Signature: c.GetHeader("X-Webhook-Signature"),
		WebhookID: c.GetHeader("X-Webhook-ID"),
		Timestamp: c.GetHeader("X-Webhook-Timestamp"),
		IDKey:     c.GetHeader("X-Idempotency-Key"),
	}

	// Log de la recepción del webhook
	c.Header("X-Webhook-Received", "true")

	// Aquí puedes agregar tu lógica de procesamiento del webhook
	// Por ejemplo: guardar en base de datos, enviar notificaciones, etc.

	// Simular procesamiento
	processed := h.processWebhook(payload, headers)

	response := dto.WebhookReceivedResponse{
		Success:   true,
		Message:   "Webhook received and processed successfully",
		Processed: processed,
		Timestamp: time.Now(),
	}

	c.JSON(http.StatusOK, response)
}

// processWebhook procesa el webhook recibido
func (h *WebhookHandler) processWebhook(payload dto.WebhookReceivedPayload, headers dto.WebhookReceivedHeaders) bool {
	// Log del procesamiento
	// log.Printf("Processing webhook: WebhookID=%d, DataType=%s", payload.WebhookID, payload.DataType)

	// Procesar según el tipo de datos
	switch payload.DataType {
	case "consumption":
		return h.processConsumptionWebhook(payload)
	case "bills":
		return h.processBillsWebhook(payload)
	default:
		// log.Printf("Unknown data type: %s", payload.DataType)
		return true
	}
}

// processConsumptionWebhook procesa webhooks de consumo
func (h *WebhookHandler) processConsumptionWebhook(payload dto.WebhookReceivedPayload) bool {
	// Validar que tenga los campos necesarios para consumo
	if payload.GroupBy == nil || payload.SendInterval == nil || payload.Period == nil {
		// log.Printf("Missing required fields for consumption webhook")
		return false
	}

	// Aquí implementarías la lógica específica para datos de consumo
	// Por ejemplo:
	// - Guardar datos de consumo en base de datos
	// - Enviar notificaciones a usuarios
	// - Procesar métricas de energía

	// log.Printf("Processing consumption data: %d contracts, period %s to %s, group by %s",
	//     len(payload.Data), payload.Period.StartDate, payload.Period.EndDate, *payload.GroupBy)

	return true
}

// processBillsWebhook procesa webhooks de facturas
func (h *WebhookHandler) processBillsWebhook(payload dto.WebhookReceivedPayload) bool {
	// Validar que tenga los campos necesarios para facturas
	if payload.TriggerType == nil || payload.Bill == nil {
		// log.Printf("Missing required fields for bills webhook")
		return false
	}

	// Aquí implementarías la lógica específica para eventos de facturas
	// Por ejemplo:
	// - Notificar a usuarios sobre nuevas facturas
	// - Procesar pagos
	// - Actualizar estados de facturas

	// log.Printf("Processing bills data: webhook_id=%d, trigger_type=%s, bill_id=%d",
	//     payload.WebhookID, *payload.TriggerType, payload.Bill.BillID)

	return true
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
