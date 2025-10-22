package dto

import "time"

// WebhookReceivedPayload representa el payload recibido en un webhook (dinámico para consumo y facturas)
type WebhookReceivedPayload struct {
	WebhookID int       `json:"webhook_id" binding:"required"`
	DataType  string    `json:"data_type" binding:"required"`
	Timestamp time.Time `json:"timestamp" binding:"required"`

	// Campos específicos para consumo (opcionales)
	GroupBy      *string               `json:"group_by,omitempty"`
	SendInterval *string               `json:"send_interval,omitempty"`
	Period       *WebhookPeriod        `json:"period,omitempty"`
	Data         []WebhookContractData `json:"data,omitempty"`

	// Campos específicos para facturas (opcionales)
	TriggerType *string             `json:"trigger_type,omitempty"`
	Bill        *BillWebhookData    `json:"bill,omitempty"`
	Payment     *PaymentWebhookData `json:"payment,omitempty"`
}

// WebhookPeriod representa el período de tiempo de los datos
type WebhookPeriod struct {
	StartDate string `json:"start_date"` // "2025-10-08"
	EndDate   string `json:"end_date"`   // "2025-10-09"
}

// WebhookContractData representa los datos de consumo de un contrato específico
type WebhookContractData struct {
	ContractID   int         `json:"contract_id"`
	ContractName string      `json:"contract_name"`
	SIC          string      `json:"sic"`
	Consumption  interface{} `json:"consumption"` // Puede ser WebhookHourlyConsumptionSummary, WebhookDailyConsumptionSummary, etc.
}

// WebhookEnergyMetrics contiene las métricas de consumo de energía
type WebhookEnergyMetrics struct {
	ActiveEnergy       *float64 `json:"active_energy,omitempty"`
	ActiveExport       *float64 `json:"active_export,omitempty"`
	InductivePenalized *float64 `json:"inductive_penalized,omitempty"`
	ReactiveCapacitive *float64 `json:"reactive_capacitive,omitempty"`
}

// WebhookMonthlyConsumptionSummary representa datos de consumo agregados por mes
type WebhookMonthlyConsumptionSummary struct {
	Month string `json:"month"` // "2025-10"
	WebhookEnergyMetrics
}

// WebhookDailyConsumptionSummary representa datos de consumo agregados por día
type WebhookDailyConsumptionSummary struct {
	Date string `json:"date"` // "2025-10-08"
	WebhookEnergyMetrics
}

// WebhookHourlyConsumptionSummary representa datos de consumo agregados por hora
type WebhookHourlyConsumptionSummary struct {
	Hour int `json:"hour"` // 0-23
	WebhookEnergyMetrics
}

// WebhookDateAndHourlyConsumptionSummary representa consumo agrupado por fecha con todas las horas del día
type WebhookDateAndHourlyConsumptionSummary struct {
	Date  string                            `json:"date"`  // "2025-10-01"
	Hours []WebhookHourlyConsumptionSummary `json:"hours"` // Array de 24 horas (0-23)
}

// BillWebhookPayload representa el payload para eventos de facturas
type BillWebhookPayload struct {
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

// WebhookReceivedResponse representa la respuesta del webhook receiver
type WebhookReceivedResponse struct {
	Success   bool      `json:"success"`
	Message   string    `json:"message,omitempty"`
	Processed bool      `json:"processed"`
	Timestamp time.Time `json:"timestamp"`
}

// WebhookReceivedHeaders representa los headers importantes del webhook
type WebhookReceivedHeaders struct {
	Signature string `json:"signature"`
	WebhookID string `json:"webhook_id"`
	Timestamp string `json:"timestamp"`
	IDKey     string `json:"idempotency_key"`
}
