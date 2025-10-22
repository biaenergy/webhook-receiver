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

// WebhookPayload representa el payload que se enviará (compatible con bia-consumptions)
type WebhookPayload struct {
	WebhookID    int                   `json:"webhook_id"`
	DataType     string                `json:"data_type"`
	GroupBy      string                `json:"group_by"`
	SendInterval string                `json:"send_interval"`
	Period       WebhookPeriod         `json:"period"`
	Data         []WebhookContractData `json:"data"`
	Timestamp    time.Time             `json:"timestamp"`
}

// WebhookPeriod representa el período de tiempo
type WebhookPeriod struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

// WebhookContractData representa los datos de un contrato
type WebhookContractData struct {
	ContractID   int         `json:"contract_id"`
	ContractName string      `json:"contract_name"`
	SIC          string      `json:"sic"`
	Consumption  interface{} `json:"consumption"`
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

	// Crear payload (compatible con bia-consumptions)
	payload := WebhookPayload{
		WebhookID:    12345,
		DataType:     "consumption",
		GroupBy:      "hour",
		SendInterval: "daily",
		Period: WebhookPeriod{
			StartDate: "2024-01-15",
			EndDate:   "2024-01-16",
		},
		Data: []WebhookContractData{
			{
				ContractID:   1001,
				ContractName: "Contrato Demo",
				SIC:          "123456789",
				Consumption: []map[string]interface{}{
					{
						"hour":                0,
						"active_energy":       150.5,
						"active_export":       0.0,
						"inductive_penalized": 10.2,
						"reactive_capacitive": 5.1,
					},
					{
						"hour":                1,
						"active_energy":       145.3,
						"active_export":       0.0,
						"inductive_penalized": 9.8,
						"reactive_capacitive": 4.9,
					},
				},
			},
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
	req.Header.Set("X-Webhook-ID", "12345")
	req.Header.Set("X-Idempotency-Key", fmt.Sprintf("test-%d", time.Now().Unix()))

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
		fmt.Println("✅ Webhook enviado exitosamente!")
	} else {
		fmt.Println("❌ Error al enviar webhook")
	}
}

// generateSignature genera la firma HMAC-SHA256 del payload
func generateSignature(secretKey string, payload []byte) string {
	mac := hmac.New(sha256.New, []byte(secretKey))
	mac.Write(payload)
	return hex.EncodeToString(mac.Sum(nil))
}
