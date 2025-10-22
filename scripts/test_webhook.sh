#!/bin/bash

# Script para probar el webhook receiver
# Uso: ./scripts/test_webhook.sh [SECRET_KEY] [URL]

# Configuración por defecto
SECRET_KEY=${1:-"default-secret-key"}
URL=${2:-"http://localhost:8080/webhook"}
TIMESTAMP=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

# Payload de ejemplo (compatible con bia-consumptions)
PAYLOAD='{
  "webhook_id": 12345,
  "data_type": "consumption",
  "group_by": "hour",
  "send_interval": "daily",
  "period": {
    "start_date": "2024-01-15",
    "end_date": "2024-01-16"
  },
  "data": [
    {
      "contract_id": 1001,
      "contract_name": "Contrato Demo",
      "sic": "123456789",
      "consumption": [
        {
          "hour": 0,
          "active_energy": 150.5,
          "active_export": 0.0,
          "inductive_penalized": 10.2,
          "reactive_capacitive": 5.1
        },
        {
          "hour": 1,
          "active_energy": 145.3,
          "active_export": 0.0,
          "inductive_penalized": 9.8,
          "reactive_capacitive": 4.9
        }
      ]
    }
  ],
  "timestamp": "'$TIMESTAMP'"
}'

echo "🧪 Testing Webhook Receiver"
echo "================================"
echo "URL: $URL"
echo "Secret Key: $SECRET_KEY"
echo "Timestamp: $TIMESTAMP"
echo ""

# Generar firma HMAC-SHA256
echo "🔐 Generating signature..."
SIGNATURE=$(echo -n "$PAYLOAD" | openssl dgst -sha256 -hmac "$SECRET_KEY" -binary | xxd -p)
echo "Signature: $SIGNATURE"
echo ""

# Enviar petición
echo "📤 Sending webhook request..."
echo "Payload:"
echo "$PAYLOAD" | jq .
echo ""

# Realizar petición HTTP
RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "$URL" \
  -H "Content-Type: application/json" \
  -H "X-Webhook-Signature: $SIGNATURE" \
  -H "X-Webhook-Timestamp: $TIMESTAMP" \
  -H "X-Webhook-ID: 12345" \
  -H "X-Idempotency-Key: test-$(date +%s)" \
  -d "$PAYLOAD")

# Separar respuesta y código HTTP
HTTP_CODE=$(echo "$RESPONSE" | tail -n 1)
RESPONSE_BODY=$(echo "$RESPONSE" | sed '$d')

echo "📥 Response:"
echo "HTTP Code: $HTTP_CODE"
echo "Response Body:"
echo "$RESPONSE_BODY" | jq . 2>/dev/null || echo "$RESPONSE_BODY"
echo ""

# Verificar resultado
if [ "$HTTP_CODE" = "200" ]; then
    echo "✅ Webhook received successfully!"
else
    echo "❌ Webhook failed with HTTP $HTTP_CODE"
    exit 1
fi
