# Webhook Receiver

Un servidor sencillo para recibir webhooks con verificación de firma HMAC.

## 🚀 Características

- ✅ **Verificación de firma HMAC-SHA256** para autenticar webhooks
- ✅ **Validación de timestamp** para prevenir replay attacks
- ✅ **Middleware de seguridad** para headers requeridos
- ✅ **Estructura limpia** siguiendo patrones de Go
- ✅ **Dos tipos de webhooks**: Consumo y Facturas
- ✅ **Sincronizado con bia-consumptions** (estructura actualizada)
- ✅ **Documentación automática** con Swagger
- ✅ **Health check** endpoint

## ⚠️ IMPORTANTE: Estructura de Datos

Este proyecto está **sincronizado con la estructura real** que envía `bia-consumptions`:

### Webhooks de Consumo (`data_type: "consumption"`)
- ✅ El campo `data` es **UN SOLO contrato** (no un array)
- ✅ Campos de energía: `active_energy`, `active_export`, `inductive_penalized`, `reactive_capacitive`
- ✅ Soporta agrupación por: `hour`, `day`, `month`
- ✅ Intervalos de envío: `hourly`, `daily`, `monthly`

### Webhooks de Facturas (`data_type: "bills"`)
- ✅ Dos tipos de trigger: `available` (factura disponible) y `paid` (factura pagada)
- ✅ Campo `payment` solo presente cuando `trigger_type="paid"`

## 📁 Estructura del Proyecto

```
webhook_receiver/
├── internal/
│   ├── dto/                    # Data Transfer Objects
│   │   └── webhook_dto.go
│   ├── handlers/               # HTTP Handlers
│   │   └── webhook_handler.go
│   ├── middleware/             # Middleware
│   │   └── signature_middleware.go
│   └── router/                 # Router configuration
│       └── router.go
├── main.go                     # Punto de entrada
├── go.mod                      # Dependencias
├── config.env.example         # Variables de entorno de ejemplo
└── README.md                  # Este archivo
```

## 🛠️ Instalación y Uso

### 1. Clonar y configurar

```bash
cd /Users/jeremy/go/src/webhook_receiver
go mod tidy
```

### 2. Configurar variables de entorno

```bash
# Copiar archivo de ejemplo
cp config.env.example .env

# Editar variables
nano .env
```

Variables de entorno:
```env
PORT=8080
WEBHOOK_SECRET_KEY=tu-clave-secreta-aqui
LOG_LEVEL=info
```

### 3. Ejecutar el servidor

```bash
# Desarrollo
go run main.go

# O con variables de entorno
WEBHOOK_SECRET_KEY=mi-clave-secreta PORT=8080 go run main.go
```

## 📡 Endpoints

### Health Check
```http
GET /health
```

**Respuesta:**
```json
{
  "status": "healthy",
  "timestamp": "2024-01-15T10:30:00Z",
  "service": "webhook-receiver"
}
```

### Recibir Webhook
```http
POST /webhook
```

**Headers requeridos:**
- `X-Webhook-Signature`: Firma HMAC-SHA256 del payload
- `X-Webhook-Timestamp`: Timestamp en formato RFC3339

**Headers opcionales:**
- `X-Webhook-ID`: ID del webhook
- `X-Idempotency-Key`: Clave de idempotencia

**Payload de ejemplo para CONSUMO:**
```json
{
  "webhook_id": 12345,
  "data_type": "consumption",
  "group_by": "hour",
  "send_interval": "daily",
  "period": {
    "start_date": "2024-01-15",
    "end_date": "2024-01-16"
  },
  "data": {
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
      }
    ]
  },
  "timestamp": "2024-01-15T10:30:00Z"
}
```

**Payload de ejemplo para FACTURAS:**
```json
{
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
  "timestamp": "2024-01-15T10:30:00Z"
}
```

**Respuesta exitosa:**
```json
{
  "success": true,
  "message": "Webhook received and processed successfully",
  "processed": true,
  "timestamp": "2024-01-15T10:30:00Z"
}
```

## 🔐 Verificación de Firma

El servidor verifica automáticamente la firma de cada webhook usando HMAC-SHA256:

### Cómo generar la firma (lado del cliente):

```go
import (
    "crypto/hmac"
    "crypto/sha256"
    "encoding/hex"
)

func generateSignature(secretKey string, payload []byte) string {
    mac := hmac.New(sha256.New, []byte(secretKey))
    mac.Write(payload)
    return hex.EncodeToString(mac.Sum(nil))
}
```

### Headers requeridos:

1. **X-Webhook-Signature**: Firma HMAC del payload
2. **X-Webhook-Timestamp**: Timestamp en formato RFC3339

### Validaciones de seguridad:

- ✅ Verificación de firma HMAC-SHA256
- ✅ Validación de timestamp (máximo 5 minutos de antigüedad)
- ✅ Verificación de headers requeridos
- ✅ Prevención de replay attacks

## 🧪 Testing

### Ejemplo de curl para testing:

```bash
# 1. Definir payload de consumo
SECRET_KEY="default-secret-key"
PAYLOAD='{
  "webhook_id": 12345,
  "data_type": "consumption",
  "group_by": "hour",
  "send_interval": "daily",
  "period": {
    "start_date": "2024-01-15",
    "end_date": "2024-01-16"
  },
  "data": {
    "contract_id": 1001,
    "contract_name": "Contrato Demo",
    "sic": "123456789",
    "consumption": [{"hour": 0, "active_energy": 150.5}]
  },
  "timestamp": "2024-01-15T10:30:00Z"
}'

# 2. Generar firma HMAC-SHA256
SIGNATURE=$(echo -n "$PAYLOAD" | openssl dgst -sha256 -hmac "$SECRET_KEY" -binary | xxd -p)

# 3. Enviar webhook
curl -X POST http://localhost:8080/webhook \
  -H "Content-Type: application/json" \
  -H "X-Webhook-Signature: $SIGNATURE" \
  -H "X-Webhook-Timestamp: 2024-01-15T10:30:00Z" \
  -H "X-Webhook-ID: 12345" \
  -d "$PAYLOAD"
```

## 🔧 Configuración Avanzada

### Variables de entorno:

| Variable | Descripción | Valor por defecto |
|----------|-------------|-------------------|
| `PORT` | Puerto del servidor | `8080` |
| `WEBHOOK_SECRET_KEY` | Clave secreta para verificación | `default-secret-key` |
| `GIN_MODE` | Modo de Gin (debug/release/test) | `debug` |
| `LOG_LEVEL` | Nivel de logging | `info` |

### Modos de ejecución:

```bash
# Desarrollo (con logs detallados)
GIN_MODE=debug go run main.go

# Producción (logs mínimos)
GIN_MODE=release go run main.go

# Testing
GIN_MODE=test go run main.go
```

## 🚀 Despliegue

### Docker (opcional):

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o webhook-receiver main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/webhook-receiver .
EXPOSE 8080
CMD ["./webhook-receiver"]
```

### Docker Compose:

```yaml
version: '3.8'
services:
  webhook-receiver:
    build: .
    ports:
      - "8080:8080"
    environment:
      - PORT=8080
      - WEBHOOK_SECRET_KEY=tu-clave-secreta
      - GIN_MODE=release
```

## 📝 Logs

El servidor registra automáticamente:

- ✅ Peticiones HTTP entrantes
- ✅ Errores de verificación de firma
- ✅ Timestamps de webhooks
- ✅ Respuestas del servidor

## 🔍 Troubleshooting

### Error: "Missing X-Webhook-Signature header"
- **Causa**: El cliente no está enviando el header de firma
- **Solución**: Asegúrate de incluir `X-Webhook-Signature` en la petición

### Error: "Invalid signature"
- **Causa**: La firma no coincide con el payload
- **Solución**: Verifica que estés usando la misma clave secreta y el payload correcto

### Error: "Webhook timestamp too old"
- **Causa**: El timestamp es mayor a 5 minutos
- **Solución**: Asegúrate de que el timestamp esté en formato RFC3339 y sea reciente

## 🤝 Contribución

1. Fork el proyecto
2. Crea una rama para tu feature (`git checkout -b feature/nueva-funcionalidad`)
3. Commit tus cambios (`git commit -m 'feat: agregar nueva funcionalidad'`)
4. Push a la rama (`git push origin feature/nueva-funcionalidad`)
5. Abre un Pull Request

## 📄 Licencia

Este proyecto es de uso interno para BIA Energy.
