package main

import (
	"formula1-crud-go/database"
	"formula1-crud-go/handlers"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// Conectar a la base de datos
	database.ConectarBaseDeDatos()

	// Crear router
	router := gin.Default()

	// Configurar CORS
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Inicializar manejador
	manejador := handlers.NuevoManejadorPilotos(database.DB)

	// Routes
	api := router.Group("/api")
	{
		api.GET("/pilotos", manejador.ObtenerPilotos)
		api.GET("/pilotos/:id", manejador.ObtenerPiloto)
		api.POST("/pilotos", manejador.CrearPiloto)
		api.PUT("/pilotos/:id", manejador.ActualizarPiloto)
		api.DELETE("/pilotos/:id", manejador.EliminarPiloto)
		api.GET("/estadisticas", manejador.ObtenerEstadisticas)
		api.GET("/buscar", manejador.BuscarPorEquipo)
	}

	// Servir archivos estáticos
	router.Static("/static", "./frontend")
	router.GET("/", func(c *gin.Context) {
		c.File("./frontend/index.html")
	})

	// Obtener puerto de environment variable
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Servidor Formula 1 CRUD ejecutándose en puerto %s", port)
	log.Printf("API disponible en: http://localhost:%s/api/pilotos", port)
	log.Printf("Frontend disponible en: http://localhost:%s", port)

	if err := router.Run(":" + port); err != nil {
		log.Fatal("Error iniciando el servidor:", err)
	}
}
