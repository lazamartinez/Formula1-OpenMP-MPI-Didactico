#!/bin/bash

echo "�️ INICIANDO PROYECTO FORMULA 1 CRUD - GO + POSTGRESQL"
echo "======================================================="

# Verificar que Docker está instalado
if ! command -v docker &> /dev/null; then
    echo "❌ Docker no está instalado. Por favor instala Docker primero."
    exit 1
fi

# Verificar que Docker Compose está instalado
if ! command -v docker-compose &> /dev/null; then
    echo "❌ Docker Compose no está instalado. Instalando..."
    sudo curl -L "https://github.com/docker/compose/releases/download/v2.23.0/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
    sudo chmod +x /usr/local/bin/docker-compose
fi

echo "� Construyendo contenedores Docker..."
docker-compose build

echo "� Iniciando la aplicación con PostgreSQL..."
docker-compose up -d

echo "⏳ Esperando que los servicios inicien..."
sleep 10

echo "✅ ¡Aplicación iniciada correctamente!"
echo "� Frontend: http://localhost:8080"
echo "� API: http://localhost:8080/api/pilotos"
echo "� PostgreSQL: localhost:5432"
echo "� pgAdmin: http://localhost:5050 (admin@formula1.com / admin123)"
echo " "
echo "� Comandos útiles:"
echo "   docker-compose logs -f              # Ver logs en tiempo real"
echo "   docker-compose down                 # Detener aplicación"
echo "   docker-compose restart              # Reiniciar aplicación"
echo "   docker exec -it postgres-formula1 psql -U formula1_user -d formula1_db # Conectar a PostgreSQL"

# Esperar un poco más y probar la API
sleep 5
echo " "
echo "� Probando conexión a la API..."
curl -s http://localhost:8080/api/pilotos | python3 -m json.tool || echo "⚠️  La API aún no está lista, espera unos segundos más..."

echo " "
echo "�� ¡Configuración completada! La aplicación está corriendo con PostgreSQL."
