#!/bin/bash

# Script para probar el webhook receiver con payload de facturas
# Uso: ./scripts/test_bills_webhook.sh [SECRET_KEY] [URL]

# Configuración por defecto
SECRET_KEY=${1:-"default-secret-key"}
URL=${2:-"http://localhost:8080/webhook"}
TIMESTAMP=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

# Payload de ejemplo para facturas (compatible con bia-consumptions)
PAYLOAD='{
  "webhook_id": 67890,
  "data_type": "bills",
  "trigger_type": "available",
  "bill": {
    "bill_id": 1001,
    "contract_id": 2001,
    "period": "2024-01",
    "total": 1250.75,
    "status": "pending",
    "xml_url": "https://example.com/bill_1001.xml"
  },
  "timestamp": "'$TIMESTAMP'"
}'

echo "🧪 Testing Bills Webhook Receiver"
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
echo "📤 Sending bills webhook request..."
echo "Payload:"
echo "$PAYLOAD" | jq .
echo ""

# Realizar petición HTTP
RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "$URL" \
  -H "Content-Type: application/json" \
  -H "X-Webhook-Signature: $SIGNATURE" \
  -H "X-Webhook-Timestamp: $TIMESTAMP" \
  -H "X-Webhook-ID: 67890" \
  -H "X-Idempotency-Key: bills-test-$(date +%s)" \
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
    echo "✅ Bills webhook received successfully!"
else
    echo "❌ Bills webhook failed with HTTP $HTTP_CODE"
    exit 1
fi
