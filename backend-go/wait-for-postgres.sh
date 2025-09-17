#!/bin/sh
# wait-for-postgres.sh
# Espera a que PostgreSQL esté disponible antes de ejecutar el backend

set -e

# Usar variables de entorno o valores por defecto
host="${DB_HOST:-postgres-formula1}"
user="${DB_USER:-formula1_user}"
password="${DB_PASSWORD:-formula1_password}"
dbname="${DB_NAME:-formula1_db}"
port="${DB_PORT:-5432}"

cmd="$@"

echo "⏳ Esperando a que PostgreSQL en $host:$port esté disponible..."

until PGPASSWORD="$password" psql -h "$host" -p "$port" -U "$user" -d "$dbname" -c '\q' >/dev/null 2>&1; do
  >&2 echo "⏳ PostgreSQL no está disponible - esperando 2 segundos..."
  sleep 2
done

>&2 echo "✅ PostgreSQL está disponible - ejecutando comando"
exec $cmd