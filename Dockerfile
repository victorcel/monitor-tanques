FROM golang:1.24-alpine AS builder

# Instalamos dependencias necesarias
RUN apk update && apk add --no-cache git ca-certificates tzdata && update-ca-certificates

# Configuramos el directorio de trabajo
WORKDIR /app

# Copiamos los archivos go.mod y go.sum y descargamos las dependencias
# Esto aprovecha la caché de Docker si las dependencias no han cambiado
COPY go.mod go.sum ./
RUN go mod download

# Copiamos el resto del código fuente
COPY . .

# Compilamos la aplicación con optimizaciones
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o /go/bin/monitor-tanques

# Imagen final más pequeña
FROM alpine:latest

# Añadimos certificados y zona horaria
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/bin/monitor-tanques /app/monitor-tanques

# Directorio para configuraciones y datos persistentes
VOLUME ["/app/data", "/app/configs"]

# Puerto para la API REST
EXPOSE 8080

# Ejecutamos la aplicación
CMD ["/app/monitor-tanques"]
