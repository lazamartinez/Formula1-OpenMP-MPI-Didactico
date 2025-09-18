FROM golang:1.23-alpine

WORKDIR /app

# Instalar dependencias del sistema
RUN apk add --no-cache gcc musl-dev bash postgresql-client

# Copiar archivos de dependencias
COPY go.mod go.sum ./
RUN go mod download

# Copiar código fuente
COPY . .

# Compilar la aplicación
RUN go build -o formula1-crud .

# Crear el script de espera directamente en el contenedor (evita problemas de Windows)
RUN echo '#!/bin/sh' > /app/wait-for-postgres.sh && \
    echo 'set -e' >> /app/wait-for-postgres.sh && \
    echo 'echo "⏳ Esperando a PostgreSQL en postgres-formula1:5432..."' >> /app/wait-for-postgres.sh && \
    echo 'until PGPASSWORD="formula1_password" psql -h "postgres-formula1" -p "5432" -U "formula1_user" -d "formula1_db" -c "SELECT 1;" >/dev/null 2>&1; do' >> /app/wait-for-postgres.sh && \
    echo '  echo "⏳ PostgreSQL no disponible - reintentando en 2 segundos..."' >> /app/wait-for-postgres.sh && \
    echo '  sleep 2' >> /app/wait-for-postgres.sh && \
    echo 'done' >> /app/wait-for-postgres.sh && \
    echo 'echo "✅ PostgreSQL disponible - ejecutando aplicación..."' >> /app/wait-for-postgres.sh && \
    echo 'exec ./formula1-crud' >> /app/wait-for-postgres.sh && \
    chmod +x /app/wait-for-postgres.sh

# Exponer puerto
EXPOSE 8080

# Comando para ejecutar
CMD ["/bin/sh", "-c", "/app/wait-for-postgres.sh"]