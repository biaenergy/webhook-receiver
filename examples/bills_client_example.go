package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// BillsWebhookPayload representa el payload para facturas (compatible con bia-consumptions)
type BillsWebhookPayload struct {
	WebhookID   int                 `json:"webhook_id"`
	DataType    string              `json:"data_type"`    // "bills"
	TriggerType string              `json:"trigger_type"` // "available" o "paid"
	Bill        BillWebhookData     `json:"bill"`
	Payment     *PaymentWebhookData `json:"payment,omitempty"` // Solo para trigger "paid"
	Timestamp   time.Time           `json:"timestamp"`
}

// BillWebhookData contiene los datos de la factura
type BillWebhookData struct {
	BillID     int     `json:"bill_id"`
	ContractID int     `json:"contract_id"`
	Period     string  `json:"period"`
	Total      float64 `json:"total"`
	Status     string  `json:"status"`
	XmlUrl     string  `json:"xml_url"`
}

// PaymentWebhookData contiene los datos del pago
type PaymentWebhookData struct {
	PaymentDate   time.Time `json:"payment_date"`
	TransactionID int       `json:"transaction_id,omitempty"`
	PaymentMethod string    `json:"payment_method,omitempty"`
}

// WebhookResponse representa la respuesta del servidor
type WebhookResponse struct {
	Success   bool      `json:"success"`
	Message   string    `json:"message,omitempty"`
	Processed bool      `json:"processed"`
	Timestamp time.Time `json:"timestamp"`
}

func main() {
	// Configuración
	secretKey := "default-secret-key"
	webhookURL := "http://localhost:8080/webhook"

	// Crear payload para factura disponible
	payload := BillsWebhookPayload{
		WebhookID:   67890,
		DataType:    "bills",
		TriggerType: "available",
		Bill: BillWebhookData{
			BillID:     1001,
			ContractID: 2001,
			Period:     "2024-01",
			Total:      1250.75,
			Status:     "pending",
			XmlUrl:     "https://example.com/bill_1001.xml",
		},
		Timestamp: time.Now(),
	}

	// Convertir a JSON
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("Error marshaling payload: %v\n", err)
		return
	}

	// Generar firma HMAC-SHA256
	signature := generateSignature(secretKey, payloadBytes)

	// Crear petición HTTP
	req, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return
	}

	// Agregar headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Webhook-Signature", signature)
	req.Header.Set("X-Webhook-Timestamp", time.Now().UTC().Format(time.RFC3339))
	req.Header.Set("X-Webhook-ID", "67890")
	req.Header.Set("X-Idempotency-Key", fmt.Sprintf("bills-test-%d", time.Now().Unix()))

	// Enviar petición
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error sending request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// Leer respuesta
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response: %v\n", err)
		return
	}

	// Parsear respuesta
	var webhookResp WebhookResponse
	if err := json.Unmarshal(body, &webhookResp); err != nil {
		fmt.Printf("Error parsing response: %v\n", err)
		return
	}

	// Mostrar resultado
	fmt.Printf("Status Code: %d\n", resp.StatusCode)
	fmt.Printf("Response: %+v\n", webhookResp)

	if resp.StatusCode == 200 && webhookResp.Success {
		fmt.Println("✅ Bills webhook enviado exitosamente!")
	} else {
		fmt.Println("❌ Error al enviar bills webhook")
	}
}

// generateSignature genera la firma HMAC-SHA256 del payload
func generateSignature(secretKey string, payload []byte) string {
	mac := hmac.New(sha256.New, []byte(secretKey))
	mac.Write(payload)
	return hex.EncodeToString(mac.Sum(nil))
}
