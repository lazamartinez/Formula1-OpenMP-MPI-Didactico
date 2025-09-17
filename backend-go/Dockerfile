FROM golang:1.23-alpine

WORKDIR /app

# Instalar dependencias del sistema
RUN apk add --no-cache gcc musl-dev bash postgresql-client

# Copiar archivos de dependencias
COPY go.mod go.sum ./
RUN go mod download

# Copiar script de espera
COPY wait-for-postgres.sh ./
RUN chmod +x wait-for-postgres.sh

# Copiar código fuente
COPY . .

# Compilar la aplicación
RUN go build -o formula1-crud .

# Exponer puerto
EXPOSE 8080

# Comando para ejecutar
CMD ["./wait-for-postgres.sh", "postgres-formula1", "./formula1-crud"]
