# Changelog - Webhook Receiver

## [2.0.0] - 2025-10-28

### ✨ ACTUALIZACIÓN MAYOR: Sincronización con bia-consumptions

Este release actualiza completamente la estructura de datos del proyecto `webhook_receiver` para estar **100% sincronizado** con la estructura real que envía `bia-consumptions`.

### 🔄 Cambios Importantes (BREAKING CHANGES)

#### 1. **DTOs Separados por Tipo de Webhook**
- ✅ **Antes**: Un DTO genérico `WebhookReceivedPayload` con campos opcionales
- ✅ **Ahora**: Dos DTOs específicos:
  - `WebhookPayload` para webhooks de **consumo**
  - `BillWebhookPayload` para webhooks de **facturas**

#### 2. **Campo `Data` en Webhooks de Consumo**
- ⚠️ **BREAKING CHANGE**: El campo `data` ahora es **UN SOLO objeto de contrato**, no un array
- **Antes**: `Data []WebhookContractData`
- **Ahora**: `Data WebhookContractData`
- **Razón**: Cada webhook de bia-consumptions envía datos de un solo contrato por webhook

#### 3. **Campos de Energía Actualizados**
Los campos de energía están correctamente sincronizados:
- ✅ `active_energy` - Energía activa
- ✅ `active_export` - Exportación de energía activa
- ✅ `inductive_penalized` - Reactiva inductiva penalizada
- ✅ `reactive_capacitive` - Reactiva capacitiva

### 📦 Cambios en Archivos

#### DTOs (`internal/dto/webhook_dto.go`)
```go
// ANTES (Estructura genérica)
type WebhookReceivedPayload struct {
    Data []WebhookContractData `json:"data,omitempty"` // Array
    // ... campos mezclados para consumo y facturas
}

// AHORA (Estructuras separadas)
type WebhookPayload struct {
    Data WebhookContractData `json:"data"` // Objeto único
    // ... solo campos de consumo
}

type BillWebhookPayload struct {
    // ... solo campos de facturas
}
```

#### Handler (`internal/handlers/webhook_handler.go`)
- ✅ Detecta automáticamente el tipo de webhook por el campo `data_type`
- ✅ Procesa cada tipo con su DTO específico
- ✅ Mejor manejo de errores y validación
- ✅ Mensajes de respuesta más descriptivos

#### Ejemplos de Cliente
- ✅ `examples/client_example.go` - Actualizado con estructura de consumo correcta
- ✅ `examples/bills_client_example.go` - Actualizado con comentarios explicativos

#### Documentación
- ✅ `README.md` - Ejemplos actualizados con payloads reales
- ✅ `QUICK_START.md` - Comandos curl con estructura correcta
- ✅ Nueva sección "⚠️ IMPORTANTE: Estructura de Datos"

### 🎯 Compatibilidad

Este proyecto ahora refleja **exactamente** la estructura que envía `bia-consumptions`:

#### Webhooks de Consumo
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
    "consumption": [...]
  },
  "timestamp": "2024-01-15T10:30:00Z"
}
```

#### Webhooks de Facturas
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
    "xml_url": "https://..."
  },
  "timestamp": "2024-01-15T10:30:00Z"
}
```

### 🚀 Migración desde v1.x

Si estabas usando la versión anterior:

1. **Actualizar parsing de payload de consumo:**
   ```go
   // ANTES
   var payload dto.WebhookReceivedPayload
   // payload.Data es un array
   
   // AHORA
   var payload dto.WebhookPayload
   // payload.Data es un objeto único
   ```

2. **Actualizar procesamiento:**
   ```go
   // ANTES
   for _, contract := range payload.Data {
       // procesar cada contrato
   }
   
   // AHORA
   contract := payload.Data
   // procesar el contrato único
   ```

### 📝 Notas para Desarrolladores

- Este proyecto es un **ejemplo de referencia** para desarrolladores que necesiten integrar webhooks de bia-consumptions
- La estructura está **sincronizada** con el sistema de producción
- Los campos de energía y la estructura de datos son **idénticos** a los que envía bia-consumptions

### 🔒 Seguridad

No hay cambios en el sistema de seguridad:
- ✅ Verificación HMAC-SHA256 sigue igual
- ✅ Validación de timestamp funciona igual
- ✅ Headers requeridos sin cambios

---

**Fecha de actualización**: 28 de Octubre, 2025
**Sincronizado con**: bia-consumptions v1.0

