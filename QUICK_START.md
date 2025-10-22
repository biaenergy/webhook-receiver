# 🚀 Quick Start - Webhook Receiver

## Inicio Rápido

### 1. Ejecutar el servidor
```bash
cd /Users/jeremy/go/src/webhook_receiver
go run main.go
```

### 2. Verificar que funciona
```bash
curl http://localhost:8080/health
```

### 3. Probar webhook
```bash
# Opción 1: Usar el script de prueba
./scripts/test_webhook.sh

# Opción 2: Usar el ejemplo de cliente
go run examples/client_example.go
```

## 📡 Endpoints Disponibles

- `GET /health` - Health check
- `POST /webhook` - Recibir webhooks (requiere verificación de firma)

## 🔐 Verificación de Firma

El webhook requiere estos headers:
- `X-Webhook-Signature`: Firma HMAC-SHA256 del payload
- `X-Webhook-Timestamp`: Timestamp en formato RFC3339

## 🧪 Testing

### Script de prueba automático:
```bash
./scripts/test_webhook.sh [SECRET_KEY] [URL]
```

### Ejemplo de cliente Go:
```bash
go run examples/client_example.go
```

### Prueba manual con curl:
```bash
# 1. Generar firma
SECRET_KEY="default-secret-key"
PAYLOAD='{"event_type":"test","data":{"test":true},"timestamp":"2024-01-15T10:30:00Z"}'
SIGNATURE=$(echo -n "$PAYLOAD" | openssl dgst -sha256 -hmac "$SECRET_KEY" -binary | xxd -p)

# 2. Enviar webhook
curl -X POST http://localhost:8080/webhook \
  -H "Content-Type: application/json" \
  -H "X-Webhook-Signature: $SIGNATURE" \
  -H "X-Webhook-Timestamp: 2024-01-15T10:30:00Z" \
  -d "$PAYLOAD"
```

## 🐳 Docker

### Construir y ejecutar:
```bash
docker build -t webhook-receiver .
docker run -p 8080:8080 webhook-receiver
```

### Con Docker Compose:
```bash
docker-compose up -d
```

## ⚙️ Configuración

Variables de entorno:
- `PORT=8080` - Puerto del servidor
- `WEBHOOK_SECRET_KEY=default-secret-key` - Clave para verificación
- `GIN_MODE=debug` - Modo de Gin (debug/release/test)

## 📝 Logs

El servidor muestra logs automáticamente:
- ✅ Peticiones HTTP entrantes
- ✅ Errores de verificación de firma
- ✅ Timestamps de webhooks
- ✅ Respuestas del servidor

## 🔧 Comandos Útiles

```bash
# Desarrollo
make dev

# Compilar
make build

# Ejecutar tests
make test

# Limpiar
make clean

# Ver ayuda
make help
```

## ✅ Estado del Proyecto

- ✅ Servidor funcionando en puerto 8080
- ✅ Verificación de firma HMAC-SHA256
- ✅ Health check endpoint
- ✅ Scripts de prueba
- ✅ Ejemplo de cliente
- ✅ Documentación completa
- ✅ Docker support
- ✅ Makefile con comandos útiles
