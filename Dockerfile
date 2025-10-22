# Multi-stage build para webhook-receiver
FROM golang:1.21-alpine AS builder

# Instalar dependencias del sistema
RUN apk add --no-cache git ca-certificates tzdata

# Crear directorio de trabajo
WORKDIR /app

# Copiar archivos de dependencias
COPY go.mod go.sum ./

# Descargar dependencias
RUN go mod download

# Copiar código fuente
COPY . .

# Compilar la aplicación
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o webhook-receiver main.go

# Imagen final
FROM alpine:latest

# Instalar certificados y timezone data
RUN apk --no-cache add ca-certificates tzdata

# Crear usuario no-root
RUN adduser -D -s /bin/sh webhook

# Crear directorio de trabajo
WORKDIR /app

# Copiar binario desde builder
COPY --from=builder /app/webhook-receiver .

# Cambiar ownership
RUN chown webhook:webhook /app/webhook-receiver

# Cambiar a usuario no-root
USER webhook

# Exponer puerto
EXPOSE 8080

# Variables de entorno por defecto
ENV PORT=8080
ENV GIN_MODE=release
ENV WEBHOOK_SECRET_KEY=default-secret-key

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Comando por defecto
CMD ["./webhook-receiver"]
