#!/bin/bash

echo "ÌøéÔ∏è INICIANDO PROYECTO FORMULA 1 CRUD - GO + POSTGRESQL"
echo "======================================================="

# Verificar que Docker est√° instalado
if ! command -v docker &> /dev/null; then
    echo "‚ùå Docker no est√° instalado. Por favor instala Docker primero."
    exit 1
fi

# Verificar que Docker Compose est√° instalado
if ! command -v docker-compose &> /dev/null; then
    echo "‚ùå Docker Compose no est√° instalado. Instalando..."
    sudo curl -L "https://github.com/docker/compose/releases/download/v2.23.0/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
    sudo chmod +x /usr/local/bin/docker-compose
fi

echo "Ì∞≥ Construyendo contenedores Docker..."
docker-compose build

echo "Ì∫Ä Iniciando la aplicaci√≥n con PostgreSQL..."
docker-compose up -d

echo "‚è≥ Esperando que los servicios inicien..."
sleep 10

echo "‚úÖ ¬°Aplicaci√≥n iniciada correctamente!"
echo "Ìºê Frontend: http://localhost:8080"
echo "Ì≥ä API: http://localhost:8080/api/pilotos"
echo "Ì∞ò PostgreSQL: localhost:5432"
echo "Ì≥ã pgAdmin: http://localhost:5050 (admin@formula1.com / admin123)"
echo " "
echo "Ì≥ã Comandos √∫tiles:"
echo "   docker-compose logs -f              # Ver logs en tiempo real"
echo "   docker-compose down                 # Detener aplicaci√≥n"
echo "   docker-compose restart              # Reiniciar aplicaci√≥n"
echo "   docker exec -it postgres-formula1 psql -U formula1_user -d formula1_db # Conectar a PostgreSQL"

# Esperar un poco m√°s y probar la API
sleep 5
echo " "
echo "Ì∑™ Probando conexi√≥n a la API..."
curl -s http://localhost:8080/api/pilotos | python3 -m json.tool || echo "‚ö†Ô∏è  La API a√∫n no est√° lista, espera unos segundos m√°s..."

echo " "
echo "ÔøΩÔøΩ ¬°Configuraci√≥n completada! La aplicaci√≥n est√° corriendo con PostgreSQL."
