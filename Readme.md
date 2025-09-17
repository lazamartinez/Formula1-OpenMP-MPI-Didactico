# 🏎️ Formula 1 CRUD - Go + PostgreSQL + Docker

Una aplicación completa de gestión de pilotos de Fórmula 1 con API RESTful construida en Go, base de datos PostgreSQL y Docker para containerización, desarrollada para la cátedra de Paradigmas y Lenguajes de Programación de la carrera Lic. en Sistemas de Información (UNAM).

## ✨ Características

- **API RESTful** completa con operaciones CRUD para pilotos de F1
- **Frontend HTML/JS** integrado para interactuar con la API
- **Base de datos PostgreSQL** con datos de ejemplo de pilotos
- **Containerización con Docker** para fácil despliegue
- **Sistema de espera** para PostgreSQL con script personalizado
- **Configuración CORS** para permitir peticiones desde cualquier origen
- **PgAdmin** incluido para administración de la base de datos
- **Script de inicialización** automatizado

## 🚀 Comenzando

### Prerrequisitos

- Docker y Docker Compose instalados en tu sistema

### Instalación y ejecución

1. Clona el repositorio:
```bash
git clone <url-del-repositorio>
cd formula1-crud-go
```

2. Ejecuta el script de inicio (para Linux/macOS):
```bash
chmod +x iniciar.sh
./iniciar.sh
```

Para Windows, ejecuta manualmente:
```bash
docker-compose build
docker-compose up -d
```

3. La aplicación estará disponible en:
   - Frontend: http://localhost:8080
   - API: http://localhost:8080/api/pilotos
   - PgAdmin: http://localhost:5050

### Credenciales de PgAdmin
- Email: admin@formula1.com
- Contraseña: admin123

## 📖 Uso de la API

### Endpoints disponibles

| Método | Endpoint | Descripción |
|--------|----------|-------------|
| GET | `/api/pilotos` | Obtener todos los pilotos |
| GET | `/api/pilotos/:id` | Obtener un piloto por ID |
| POST | `/api/pilotos` | Crear un nuevo piloto |
| PUT | `/api/pilotos/:id` | Actualizar un piloto existente |
| DELETE | `/api/pilotos/:id` | Eliminar un piloto |
| GET | `/api/estadisticas` | Obtener estadísticas de pilotos |
| GET | `/api/buscar?equipo=nombre` | Buscar pilotos por equipo |

### Ejemplos de uso con cURL

**Obtener todos los pilotos:**
```bash
curl http://localhost:8080/api/pilotos
```

**Crear un nuevo piloto:**
```bash
curl -X POST http://localhost:8080/api/pilotos \
  -H "Content-Type: application/json" \
  -d '{
    "nombre": "Lando Norris",
    "equipo": "McLaren",
    "nacionalidad": "Británico",
    "numero": 4,
    "victorias": 1,
    "puntos": 350.5,
    "podios": 15,
    "poles": 2,
    "vueltas_rapidas": 6
  }'
```

**Buscar pilotos por equipo:**
```bash
curl http://localhost:8080/api/buscar?equipo=Ferrari
```

## 🏗️ Estructura del proyecto

```
formula1-crud-go/
├── Dockerfile                 # Configuración de Docker para el backend
├── docker-compose.yml         # Orquestación de contenedores
├── iniciar.sh                 # Script de inicio automatizado
├── wait-for-postgres.sh       # Script de espera para PostgreSQL
├── go.mod                     # Dependencias de Go
├── .env                       # Variables de entorno
├── init-db.sql               # Script de inicialización de la BD
├── main.go                   # Punto de entrada de la aplicación
├── database/
│   └── database.go           # Configuración de conexión a BD
├── handlers/
│   └── pilotos.go           # Manejadores de endpoints API
└── frontend/                 # Frontend HTML/JS (montado desde volumen)
```

## 🔧 Configuración

### Variables de entorno

El proyecto utiliza las siguientes variables de entorno (configuradas en `.env`):

| Variable | Valor por defecto | Descripción |
|----------|-------------------|-------------|
| DB_HOST | postgres-formula1 | Host de PostgreSQL |
| DB_PORT | 5432 | Puerto de PostgreSQL |
| DB_USER | formula1_user | Usuario de PostgreSQL |
| DB_PASSWORD | formula1_password | Contraseña de PostgreSQL |
| DB_NAME | formula1_db | Nombre de la base de datos |
| DB_SSLMODE | disable | Modo SSL para PostgreSQL |
| PORT | 8080 | Puerto del servidor Go |

### Personalización

Puedes modificar los valores por defecto editando el archivo `.env` o pasando las variables de entorno directamente al contenedor.

## 🐛 Solución de problemas

### Ver logs de los contenedores

```bash
docker-compose logs -f
```

### Reiniciar la aplicación

```bash
docker-compose restart
```

### Detener la aplicación

```bash
docker-compose down
```

### Conectar directamente a PostgreSQL

```bash
docker exec -it postgres-formula1 psql -U formula1_user -d formula1_db
```

### 4. Si el problema persiste, prueba ejecutar el backend manualmente:
```bash
# Ejecutar el backend en modo interactivo para ver los errores
docker-compose run --rm backend-formula1 sh

# Dentro del contenedor, prueba:
./wait-for-postgres.sh postgres-formula1 "./formula1-crud"
```

### 5. Solución alternativa - Recrear todo desde cero:
```bash
# Parar todo
docker-compose down -v

# Eliminar imágenes y contenedores huérfanos
docker system prune -a -f

# Reconstruir con más verbosidad
docker-compose build --no-cache --progress=plain

# Iniciar solo PostgreSQL primero
docker-compose up -d postgres-formula1

# Esperar a que PostgreSQL esté listo
sleep 10

# Verificar que PostgreSQL funciona
docker-compose exec postgres-formula1 psql -U formula1_user -d formula1_db -c "SELECT 1;"

# Ahora iniciar el backend
docker-compose up -d backend-formula1

# Ver logs del backend
docker-compose logs -f backend-formula1
```

### 6. Si sigue fallando, modifica el Dockerfile para debugging:
Edita el `Dockerfile` y añade esto al final:
```dockerfile
# Añade esto para debugging
CMD ["sh", "-c", "echo 'Esperando PostgreSQL...' && ./wait-for-postgres.sh postgres-formula1 && echo 'Iniciando aplicación...' && ./formula1-crud"]
```

### 7. Prueba con un comando más simple:
```bash
# Ejecutar manualmente para ver el error real
docker-compose run --rm backend-formula1 ./formula1-crud
```

### 8. Verifica que el binario se creó correctamente:
```bash
docker-compose run --rm backend-formula1 ls -la /app

# Deberías ver formula1-crud en la lista
```
### Problemas comunes

1. **Puerto ya en uso**: Asegúrate de que los puertos 8080, 5432 y 5050 estén libres
2. **Error de conexión a la BD**: Espera unos segundos tras iniciar para que PostgreSQL esté completamente disponible
3. **Permisos denegados en scripts**: Ejecuta `chmod +x wait-for-postgres.sh iniciar.sh`

## 📦 Desarrollo

### Construir manualmente

```bash
# Construir la imagen de Docker
docker-compose build

# Ejecutar en primer plano (para debugging)
docker-compose up

# Ejecutar en segundo plano
docker-compose up -d
```

### Modificar el frontend

El frontend está montado como volumen desde `../frontend-html`. Puedes modificar los archivos en esa carpeta y los cambios se reflejarán inmediatamente.

### Agregar nuevas dependencias de Go

```bash
go get <paquete>
go mod tidy
```

## 🧪 Testing

Puedes probar la API usando las herramientas siguientes:

1. **cURL**: Como se muestra en los ejemplos anteriores
2. **Postman**: Importa la colección de endpoints
3. **Navegador**: Visita http://localhost:8080 para usar la interfaz web

## 📝 Licencia

Este proyecto está bajo la Licencia MIT. Ver el archivo `LICENSE` para más detalles.

## 🤝 Contribuciones

Las contribuciones son bienvenidas. Por favor:

1. Haz fork del proyecto
2. Crea una rama para tu feature (`git checkout -b feature/AmazingFeature`)
3. Commit tus cambios (`git commit -m 'Add some AmazingFeature'`)
4. Push a la rama (`git push origin feature/AmazingFeature`)
5. Abre un Pull Request

## 🧑‍🎓 Alumnos participantes del proyecto

1. Díaz Exequiel Andres - [@exequieldev](https://github.com/exequieldev)
2. Küster Joaquín - [@joaquinkuster](https://github.com/joaquinkuster)
3. Da Silva Marcos - [@Marcos2497](https://github.com/Marcos2497)
4. Martinez Lázaro Ezequiel - [@lazamartinez](https://github.com/lazamartinez)

## 📞 Soporte

Si tienes problemas o preguntas:

1. Revisa la sección de solución de problemas arriba
2. Abre un issue en el repositorio de GitHub
3. Contacta al mantenedor del proyecto

---


¡Disfruta explorando y gestionando los datos de pilotos de Fórmula 1! 🏁
