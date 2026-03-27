# Dockerfile

# ── Stage 1: Build ──────────────────────────────────────────────
FROM golang:1.25.5-alpine AS builder

# Dependencias del sistema necesarias para cgo (lib/pq)
RUN apk add --no-cache gcc musl-dev

WORKDIR /app

# Copiar dependencias primero (mejor uso del cache de Docker)
COPY go.mod go.sum ./
RUN go mod download

# Copiar el código fuente
COPY . .

# Compilar: binario estático y optimizado
RUN CGO_ENABLED=0 GOOS=linux go build \
  -ldflags="-w -s" \
  -o /app/server \
  ./cmd/api

# ── Stage 2: Runtime ────────────────────────────────────────────
FROM alpine:3.20

# Certificados SSL para conexiones TLS externas
RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

# Solo copiar el binario compilado (imagen final muy pequeña ~15MB)
COPY --from=builder /app/server .

EXPOSE 8080

ENTRYPOINT ["./server"]