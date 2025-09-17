package database

import (
    "formula1-crud-go/models"
    "log"
    "os"
    "time"

    "github.com/joho/godotenv"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConectarBaseDeDatos() {
    // Cargar variables de entorno
    err := godotenv.Load()
    if err != nil {
        log.Println("⚠️  No se encontró archivo .env, usando variables de entorno del sistema")
    }

    // Configurar conexión PostgreSQL
    dsn := obtenerDSN()

    // Configurar logger de GORM
    newLogger := logger.New(
        log.New(os.Stdout, "\r\n", log.LstdFlags),
        logger.Config{
            SlowThreshold: time.Second,
            LogLevel:      logger.Info,
            Colorful:      true,
        },
    )

    // Conectar a PostgreSQL
    DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
        Logger: newLogger,
    })
    if err != nil {
        log.Fatal("❌ Error conectando a PostgreSQL:", err)
    }

    log.Println("✅ Conectado a PostgreSQL exitosamente")

    // Migrar el schema
    err = DB.AutoMigrate(&models.Piloto{})
    if err != nil {
        log.Fatal("❌ Error migrando la base de datos:", err)
    }

    log.Println("✅ Esquema migrado exitosamente")

    // Insertar datos de ejemplo
    insertarDatosEjemplo()
}

func obtenerDSN() string {
    host := getEnv("DB_HOST", "localhost")
    port := getEnv("DB_PORT", "5432")
    user := getEnv("DB_USER", "formula1_user")
    password := getEnv("DB_PASSWORD", "formula1_password")
    dbname := getEnv("DB_NAME", "formula1_db")
    sslmode := getEnv("DB_SSLMODE", "disable")

    return "host=" + host + " port=" + port + " user=" + user + 
           " password=" + password + " dbname=" + dbname + 
           " sslmode=" + sslmode + " TimeZone=UTC"
}

func getEnv(key, defaultValue string) string {
    value := os.Getenv(key)
    if value == "" {
        return defaultValue
    }
    return value
}

func insertarDatosEjemplo() {
    pilotosEjemplo := []models.Piloto{
        {
            Nombre: "Max Verstappen", 
            Equipo: "Red Bull", 
            Nacionalidad: "Holandés", 
            Numero: 1, 
            Victorias: 54, 
            Puntos: 575.5, 
            Podios: 98, 
            Poles: 33, 
            VueltasRapidas: 30,
        },
        {
            Nombre: "Lewis Hamilton", 
            Equipo: "Mercedes", 
            Nacionalidad: "Británico", 
            Numero: 44, 
            Victorias: 103, 
            Puntos: 4637.5, 
            Podios: 197, 
            Poles: 104, 
            VueltasRapidas: 65,
        },
        {
            Nombre: "Charles Leclerc", 
            Equipo: "Ferrari", 
            Nacionalidad: "Monegasco", 
            Numero: 16, 
            Victorias: 5, 
            Puntos: 1074, 
            Podios: 32, 
            Poles: 23, 
            VueltasRapidas: 9,
        },
        {
            Nombre: "Fernando Alonso", 
            Equipo: "Aston Martin", 
            Nacionalidad: "Español", 
            Numero: 14, 
            Victorias: 32, 
            Puntos: 2267, 
            Podios: 106, 
            Poles: 22, 
            VueltasRapidas: 24,
        },
        {
            Nombre: "Carlos Sainz", 
            Equipo: "Ferrari", 
            Nacionalidad: "Español", 
            Numero: 55, 
            Victorias: 2, 
            Puntos: 782.5, 
            Podios: 18, 
            Poles: 5, 
            VueltasRapidas: 3,
        },
    }

    for _, piloto := range pilotosEjemplo {
        var existe models.Piloto
        result := DB.Where("nombre = ?", piloto.Nombre).First(&existe)
        if result.Error != nil {
            DB.Create(&piloto)
            log.Printf("✅ Insertado piloto: %s", piloto.Nombre)
        }
    }
}
