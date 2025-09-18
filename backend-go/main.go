package main

import (
	"fmt"
	"formula1-crud-go/database"
	"formula1-crud-go/handlers"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// -------------------- Configuración WebSocket --------------------
var actualizador = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// -------------------- Tipos de Mensajes --------------------
type MensajeWS struct {
	Tipo   string `json:"tipo"`
	Topico string `json:"topico,omitempty"`
	Texto  string `json:"texto,omitempty"`
	Obj    any    `json:"obj,omitempty"`
}

type ResultadoOpenMP struct {
	AutoID          int     `json:"auto_id"`
	MejorVuelta     float64 `json:"mejor_vuelta"`
	CantidadVueltas int     `json:"cantidad_vueltas"`
}

// -------------------- Funciones de Simulación --------------------
func correrMPI(sectores int, vueltas int, enviar chan MensajeWS) {
	if sectores < 1 {
		enviar <- MensajeWS{Tipo: "registro", Topico: "mpi", Texto: "Error: sectores debe ser >= 1"}
		enviar <- MensajeWS{Tipo: "finalizado", Topico: "mpi"}
		return
	}
	if vueltas < 1 {
		vueltas = 1
	}

	enviar <- MensajeWS{Tipo: "registro", Topico: "mpi", Texto: fmt.Sprintf("Iniciando MPI: %d sectores, %d vueltas", sectores, vueltas)}

	for v := 1; v <= vueltas; v++ {
		enviar <- MensajeWS{Tipo: "registro", Topico: "mpi", Texto: fmt.Sprintf("=== Vuelta %d ===", v)}
		for s := 1; s <= sectores; s++ {
			tiempoSector := float64(rand.Intn(2300)+1200) / 100.0
			time.Sleep(300 * time.Millisecond)
			enviar <- MensajeWS{
				Tipo:   "registro",
				Topico: "mpi",
				Texto:  fmt.Sprintf("Tiempo de sector %d: %.2f s (vuelta %d)", s, tiempoSector, v),
			}
		}
	}

	enviar <- MensajeWS{Tipo: "resumen", Topico: "mpi", Obj: map[string]any{"mensaje": "MPI finalizado"}}
	enviar <- MensajeWS{Tipo: "finalizado", Topico: "mpi"}
}

func correrOpenMP(cantidadAutos int, vueltas int, enviar chan MensajeWS) {
	if cantidadAutos < 1 {
		enviar <- MensajeWS{Tipo: "registro", Topico: "openmp", Texto: "Error: cantidad de autos debe ser >= 1"}
		enviar <- MensajeWS{Tipo: "finalizado", Topico: "openmp"}
		return
	}
	if vueltas < 1 {
		vueltas = 1
	}

	enviar <- MensajeWS{Tipo: "registro", Topico: "openmp", Texto: fmt.Sprintf("Iniciando OpenMP: %d autos, %d vueltas cada uno", cantidadAutos, vueltas)}

	resultados := make([]ResultadoOpenMP, cantidadAutos)
	done := make(chan struct{})

	for auto := 0; auto < cantidadAutos; auto++ {
		go func(autoID int) {
			defer func() { done <- struct{}{} }()

			mejor := 1e9
			for v := 1; v <= vueltas; v++ {
				tiempoVuelta := float64(rand.Intn(2099)+7500) / 100.0
				time.Sleep(200 * time.Millisecond)
				enviar <- MensajeWS{Tipo: "registro", Topico: "openmp", Texto: fmt.Sprintf("Auto %d - Vuelta %d: %.2f s", autoID+1, v, tiempoVuelta)}
				if tiempoVuelta < mejor {
					mejor = tiempoVuelta
					enviar <- MensajeWS{Tipo: "registro", Topico: "openmp", Texto: fmt.Sprintf("Auto %d - Nueva mejor vuelta: %.2f s", autoID+1, mejor)}
				}
			}
			resultados[autoID] = ResultadoOpenMP{AutoID: autoID + 1, MejorVuelta: mejor, CantidadVueltas: vueltas}
		}(auto)
	}

	for i := 0; i < cantidadAutos; i++ {
		<-done
	}

	mejorGeneral := ResultadoOpenMP{AutoID: -1, MejorVuelta: 1e9}
	for _, r := range resultados {
		if r.MejorVuelta < mejorGeneral.MejorVuelta {
			mejorGeneral = r
		}
	}

	enviar <- MensajeWS{Tipo: "resumen", Topico: "openmp", Obj: map[string]any{
		"mejor_por_auto": resultados,
		"mejor_general":  mejorGeneral,
	}}
	enviar <- MensajeWS{Tipo: "finalizado", Topico: "openmp"}
}

// -------------------- Handlers WebSocket --------------------
func wsHandler(c *gin.Context) {
	conn, err := actualizador.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Error al actualizar a websocket:", err)
		return
	}
	defer conn.Close()

	enviar := make(chan MensajeWS, 100)
	defer close(enviar)

	go func() {
		for msg := range enviar {
			if err := conn.WriteJSON(msg); err != nil {
				log.Println("Error escribiendo en websocket:", err)
				return
			}
		}
	}()

	for {
		var comando map[string]any
		if err := conn.ReadJSON(&comando); err != nil {
			log.Println("Conexión cerrada o error de lectura:", err)
			return
		}
		switch comando["action"] {
		case "iniciar_mpi":
			sectores := 1
			vueltas := 1
			if v, ok := comando["sectores"].(float64); ok {
				sectores = int(v)
			}
			if v, ok := comando["vueltas"].(float64); ok {
				vueltas = int(v)
			}
			go correrMPI(sectores, vueltas, enviar)
		case "iniciar_openmp":
			autos := 3
			vueltas := 5
			if v, ok := comando["autos"].(float64); ok {
				autos = int(v)
			}
			if v, ok := comando["vueltas"].(float64); ok {
				vueltas = int(v)
			}
			go correrOpenMP(autos, vueltas, enviar)
		default:
			enviar <- MensajeWS{Tipo: "registro", Texto: fmt.Sprintf("Comando no reconocido: %v", comando["action"])}
		}
	}
}

// -------------------- HTML Template --------------------
var plantillaSimulacion = template.Must(template.New("simulacion").Parse(htmlSimulacion))

func simulacionHandler(c *gin.Context) {
	c.Header("Content-Type", "text/html; charset=utf-8")
	plantillaSimulacion.Execute(c.Writer, nil)
}

// -------------------- Main --------------------
func main() {
	rand.Seed(time.Now().UnixNano())

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

	// Routes API CRUD
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

	// Routes Simulación
	router.GET("/simulacion", simulacionHandler)
	router.GET("/ws", func(c *gin.Context) {
		wsHandler(c)
	})

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

	log.Printf("Servidor Formula 1 ejecutándose en puerto %s", port)
	log.Printf("API CRUD disponible en: http://localhost:%s/api/pilotos", port)
	log.Printf("Frontend disponible en: http://localhost:%s", port)
	log.Printf("Simulación disponible en: http://localhost:%s/simulacion", port)

	if err := router.Run(":" + port); err != nil {
		log.Fatal("Error iniciando el servidor:", err)
	}
}

// -------------------- HTML + JS embebido --------------------
const htmlSimulacion = `
<!DOCTYPE html>
<html lang="es">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Simulación F1 - MPI & OpenMP</title>
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.4.0/css/all.min.css">
    <style>
        :root {
            --f1-red: #e10600;
            --f1-black: #15151e;
            --f1-white: #ffffff;
            --f1-gray: #38383f;
            --f1-blue: #0066cc;
            --omp-yellow: #ffcc00;
            --mpi-blue: #0099ff;
            --grid-green: #00cc66;
            --circuit-gray: #4a4a57;
            --grass-green: #0a7e3a;
            --kerb-red: #ff2800;
            --kerb-yellow: #ffd700;
        }

        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }

        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            background: linear-gradient(135deg, var(--f1-black) 0%, #2a2a35 100%);
            color: var(--f1-white);
            line-height: 1.6;
            min-height: 100vh;
            overflow-x: hidden;
        }

        .container {
            max-width: 1600px;
            margin: 0 auto;
            padding: 20px;
        }

        /* Header con efecto F1 */
        header {
            text-align: center;
            margin-bottom: 30px;
            padding: 30px;
            background: linear-gradient(90deg, var(--f1-red) 0%, #ff2b2b 100%);
            border-radius: 15px;
            position: relative;
            overflow: hidden;
            box-shadow: 0 10px 30px rgba(225, 6, 0, 0.3);
        }

        header::before {
            content: '';
            position: absolute;
            top: 0;
            left: -100%;
            width: 100%;
            height: 100%;
            background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.2), transparent);
            animation: shine 3s infinite;
        }

        @keyframes shine {
            0% { left: -100%; }
            100% { left: 100%; }
        }

        header h1 {
            font-size: 2.8em;
            margin-bottom: 10px;
            text-shadow: 2px 2px 4px rgba(0, 0, 0, 0.5);
        }

        header p {
            font-size: 1.2em;
            opacity: 0.9;
        }

        /* Botones de navegación */
        .nav-buttons {
            display: flex;
            gap: 15px;
            justify-content: center;
            margin-top: 20px;
        }

        .nav-button {
            background: linear-gradient(135deg, var(--omp-yellow) 0%, #ffd700 100%);
            color: black;
            padding: 12px 20px;
            border: none;
            border-radius: 8px;
            cursor: pointer;
            font-weight: bold;
            display: flex;
            align-items: center;
            gap: 8px;
            text-decoration: none;
            transition: transform 0.2s ease, box-shadow 0.2s ease;
        }

        .nav-button:hover {
            transform: translateY(-2px);
            box-shadow: 0 5px 15px rgba(255, 204, 0, 0.4);
        }

        .nav-button.back {
            background: linear-gradient(135deg, var(--mpi-blue) 0%, #00aaff 100%);
            color: white;
        }

        /* Contenido principal */
        .main-content {
            display: flex;
            gap: 30px;
            margin-bottom: 30px;
        }

        /* Panel de control */
        .control-panel {
            flex: 1;
            background: rgba(40, 40, 50, 0.8);
            border-radius: 15px;
            padding: 20px;
            box-shadow: 0 5px 20px rgba(0, 0, 0, 0.3);
        }

        .panel-section {
            margin-bottom: 25px;
            padding-bottom: 15px;
            border-bottom: 1px solid var(--f1-gray);
        }

        .panel-section h3 {
            margin-bottom: 15px;
            color: var(--f1-red);
            display: flex;
            align-items: center;
            gap: 10px;
        }

        .panel-section h3 i {
            font-size: 1.2em;
        }

        .input-group {
            margin-bottom: 12px;
        }

        .input-group label {
            display: block;
            margin-bottom: 5px;
            font-weight: 500;
        }

        .input-group input {
            width: 100%;
            padding: 10px;
            border-radius: 6px;
            border: 1px solid var(--f1-gray);
            background: var(--f1-black);
            color: var(--f1-white);
        }

        .btn {
            width: 100%;
            padding: 12px;
            border: none;
            border-radius: 8px;
            cursor: pointer;
            font-weight: bold;
            transition: all 0.2s ease;
        }

        .btn-mpi {
            background: linear-gradient(135deg, var(--mpi-blue) 0%, #0066cc 100%);
            color: white;
        }

        .btn-mpi:hover {
            background: linear-gradient(135deg, #0066cc 0%, #004d99 100%);
            transform: translateY(-2px);
            box-shadow: 0 5px 15px rgba(0, 102, 204, 0.3);
        }

        .btn-openmp {
            background: linear-gradient(135deg, var(--omp-yellow) 0%, #e6b800 100%);
            color: black;
        }

        .btn-openmp:hover {
            background: linear-gradient(135deg, #e6b800 0%, #cc9900 100%);
            transform: translateY(-2px);
            box-shadow: 0 5px 15px rgba(255, 204, 0, 0.3);
        }

        /* Circuito */
        .circuit-container {
            flex: 2;
            position: relative;
            background: var(--grass-green);
            border-radius: 15px;
            overflow: hidden;
            box-shadow: 0 5px 20px rgba(0, 0, 0, 0.3);
            min-height: 500px;
        }

        .circuit {
            position: relative;
            width: 100%;
            height: 100%;
            padding: 20px;
        }

        .circuit-track {
            position: absolute;
            top: 20px;
            left: 20px;
            right: 20px;
            bottom: 20px;
            border: 10px solid var(--circuit-gray);
            border-radius: 50%;
        }

        .circuit-track::before {
            content: '';
            position: absolute;
            top: -10px;
            left: -10px;
            right: -10px;
            bottom: -10px;
            border: 3px solid var(--kerb-red);
            border-radius: 50%;
            border-style: dashed;
        }

        .sector {
            position: absolute;
            width: 60px;
            height: 60px;
            background: var(--circuit-gray);
            border: 3px solid var(--kerb-yellow);
            border-radius: 50%;
            display: flex;
            align-items: center;
            justify-content: center;
            color: white;
            font-weight: bold;
            z-index: 10;
        }

        .sector-node {
            position: absolute;
            width: 40px;
            height: 40px;
            background: var(--f1-black);
            border: 2px solid var(--mpi-blue);
            border-radius: 50%;
            display: flex;
            align-items: center;
            justify-content: center;
            color: white;
            font-weight: bold;
            z-index: 20;
            box-shadow: 0 0 10px rgba(0, 153, 255, 0.5);
        }

        .f1-car {
            position: absolute;
            width: 40px;
            height: 20px;
            background: var(--f1-red);
            border-radius: 4px;
            z-index: 100;
            transition: all 0.5s linear;
            box-shadow: 0 0 10px rgba(225, 6, 0, 0.7);
        }

        .f1-car::before {
            content: '';
            position: absolute;
            top: -5px;
            left: 5px;
            width: 10px;
            height: 5px;
            background: #333;
            border-radius: 2px 2px 0 0;
        }

        /* Panel de logs */
        .logs-container {
            display: flex;
            gap: 20px;
        }

        .log-panel {
            flex: 1;
            background: rgba(40, 40, 50, 0.8);
            border-radius: 15px;
            padding: 20px;
            box-shadow: 0 5px 20px rgba(0, 0, 0, 0.3);
            max-height: 300px;
            overflow-y: auto;
        }

        .log-panel h3 {
            margin-bottom: 15px;
            color: var(--f1-red);
            display: flex;
            align-items: center;
            gap: 10px;
            position: sticky;
            top: 0;
            background: rgba(40, 40, 50, 0.9);
            padding: 5px 0;
            z-index: 10;
        }

        .log-content {
            font-family: 'Courier New', monospace;
            font-size: 0.9em;
            line-height: 1.5;
        }

        .log-entry {
            margin-bottom: 8px;
            padding: 8px;
            border-radius: 4px;
            background: rgba(255, 255, 255, 0.05);
        }

        .log-entry.mpi {
            border-left: 3px solid var(--mpi-blue);
        }

        .log-entry.openmp {
            border-left: 3px solid var(--omp-yellow);
        }

        .log-entry.info {
            border-left: 3px solid var(--grid-green);
        }

        /* Estadísticas */
        .stats-panel {
            background: rgba(40, 40, 50, 0.8);
            border-radius: 15px;
            padding: 20px;
            margin-top: 20px;
            box-shadow: 0 5px 20px rgba(0, 0, 0, 0.3);
        }

        .stats-panel h3 {
            margin-bottom: 15px;
            color: var(--f1-red);
            display: flex;
            align-items: center;
            gap: 10px;
        }

        .stats-grid {
            display: grid;
            grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
            gap: 15px;
        }

        .stat-card {
            background: rgba(30, 30, 40, 0.8);
            padding: 15px;
            border-radius: 8px;
            box-shadow: 0 3px 10px rgba(0, 0, 0, 0.2);
        }

        .stat-card h4 {
            color: var(--mpi-blue);
            margin-bottom: 8px;
        }

        .stat-value {
            font-size: 1.5em;
            font-weight: bold;
            color: var(--omp-yellow);
        }

        /* Animaciones */
        @keyframes pulse {
            0% { opacity: 1; }
            50% { opacity: 0.7; }
            100% { opacity: 1; }
        }

        .pulse {
            animation: pulse 1.5s infinite;
        }

        @keyframes highlight {
            0% { background-color: rgba(255, 255, 255, 0.05); }
            50% { background-color: rgba(0, 153, 255, 0.2); }
            100% { background-color: rgba(255, 255, 255, 0.05); }
        }

        .highlight {
            animation: highlight 1s ease;
        }

        /* Responsive */
        @media (max-width: 1200px) {
            .main-content {
                flex-direction: column;
            }
            
            .logs-container {
                flex-direction: column;
            }
        }
    </style>
</head>
<body>
    <div class="container">
        <header>
            <h1><i class="fas fa-flag-checkered"></i> Simulación Fórmula 1 - MPI & OpenMP</h1>
            <p>Visualización de procesamiento paralelo con temática de carrera</p>
            <div class="nav-buttons">
                <a href="/" class="nav-button back">
                    <i class="fas fa-arrow-left"></i> Volver al CRUD
                </a>
            </div>
        </header>

        <div class="main-content">
            <div class="control-panel">
                <div class="panel-section">
                    <h3><i class="fas fa-project-diagram"></i> MPI - Sectores</h3>
                    <div class="input-group">
                        <label for="mpi-sectores">Cantidad de Sectores:</label>
                        <input type="number" id="mpi-sectores" min="1" value="5">
                    </div>
                    <div class="input-group">
                        <label for="mpi-vueltas">Vueltas:</label>
                        <input type="number" id="mpi-vueltas" min="1" value="3">
                    </div>
                    <button class="btn btn-mpi" id="start-mpi">
                        <i class="fas fa-play-circle"></i> Iniciar Simulación MPI
                    </button>
                </div>

                <div class="panel-section">
                    <h3><i class="fas fa-tachometer-alt"></i> OpenMP - Vueltas Rápidas</h3>
                    <div class="input-group">
                        <label for="openmp-autos">Cantidad de Autos:</label>
                        <input type="number" id="openmp-autos" min="1" value="4">
                    </div>
                    <div class="input-group">
                        <label for="openmp-vueltas">Vueltas por Auto:</label>
                        <input type="number" id="openmp-vueltas" min="1" value="5">
                    </div>
                    <button class="btn btn-openmp" id="start-openmp">
                        <i class="fas fa-play-circle"></i> Iniciar Simulación OpenMP
                    </button>
                </div>

                <div class="panel-section">
                    <h3><i class="fas fa-info-circle"></i> Información</h3>
                    <p>MPI simula el procesamiento de sectores de una pista en un anillo de nodos.</p>
                    <p>OpenMP simula múltiples autos corriendo en paralelo con vueltas cronometradas.</p>
                </div>
            </div>

            <div class="circuit-container">
                <div class="circuit">
                    <div class="circuit-track"></div>
                    <div id="sectors-container"></div>
                    <div id="nodes-container"></div>
                    <div id="cars-container"></div>
                </div>
            </div>
        </div>

        <div class="logs-container">
            <div class="log-panel">
                <h3><i class="fas fa-list-alt"></i> Logs MPI</h3>
                <div id="mpi-log" class="log-content"></div>
            </div>
            <div class="log-panel">
                <h3><i class="fas fa-list-alt"></i> Logs OpenMP</h3>
                <div id="openmp-log" class="log-content"></div>
            </div>
        </div>

        <div class="stats-panel">
            <h3><i class="fas fa-chart-line"></i> Estadísticas</h3>
            <div class="stats-grid">
                <div class="stat-card">
                    <h4>Mejor Tiempo MPI</h4>
                    <div id="best-mpi-time" class="stat-value">-</div>
                </div>
                <div class="stat-card">
                    <h4>Sectores Procesados</h4>
                    <div id="sectors-processed" class="stat-value">0</div>
                </div>
                <div class="stat-card">
                    <h4>Mejor Tiempo OpenMP</h4>
                    <div id="best-openmp-time" class="stat-value">-</div>
                </div>
                <div class="stat-card">
                    <h4>Vueltas Completadas</h4>
                    <div id="laps-completed" class="stat-value">0</div>
                </div>
            </div>
        </div>
    </div>

    <script>
        document.addEventListener('DOMContentLoaded', function() {
            // Conexión WebSocket
            const ws = new WebSocket("ws://" + location.host + "/ws");
            const mpiLog = document.getElementById("mpi-log");
            const openmpLog = document.getElementById("openmp-log");
            const sectorsContainer = document.getElementById("sectors-container");
            const nodesContainer = document.getElementById("nodes-container");
            const carsContainer = document.getElementById("cars-container");
            
            // Elementos de estadísticas
            const bestMpiTime = document.getElementById("best-mpi-time");
            const sectorsProcessed = document.getElementById("sectors-processed");
            const bestOpenmpTime = document.getElementById("best-openmp-time");
            const lapsCompleted = document.getElementById("laps-completed");
            
            // Variables de estado
            let currentSectors = 5;
            let mpiStats = { bestTime: Infinity, sectorsProcessed: 0 };
            let openmpStats = { bestTime: Infinity, lapsCompleted: 0 };
            let carPositions = {};
            let nodePositions = [];

            ws.onopen = () => appendLog("info", "Conexión WebSocket establecida.");
            ws.onclose = () => appendLog("info", "WebSocket cerrado.");
            ws.onerror = (e) => appendLog("info", "Error WebSocket: " + e);

            ws.onmessage = (evt) => {
                try {
                    const msg = JSON.parse(evt.data);
                    
                    if (msg.tipo === "registro") {
                        if (msg.topico === "mpi") {
                            appendLog("mpi", msg.texto);
                            processMpiMessage(msg.texto);
                        } else if (msg.topico === "openmp") {
                            appendLog("openmp", msg.texto);
                            processOpenmpMessage(msg.texto);
                        } else {
                            appendLog("info", msg.texto);
                        }
                    } else if (msg.tipo === "resumen") {
                        if (msg.topico === "mpi") {
                            appendLog("mpi", "<b>Resumen MPI:</b> " + JSON.stringify(msg.obj));
                        } else if (msg.topico === "openmp") {
                            appendLog("openmp", "<b>Resumen OpenMP:</b> " + JSON.stringify(msg.obj));
                            updateOpenmpStats(msg.obj);
                        }
                    } else if (msg.tipo === "finalizado") {
                        appendLog(msg.topico, "<i>Proceso " + msg.topico + " finalizado</i>");
                    }
                } catch(e) {
                    appendLog("info", "Mensaje no JSON: " + evt.data);
                }
            };

            // Configurar eventos de botones
            document.getElementById("start-mpi").addEventListener("click", () => {
                const sectores = parseInt(document.getElementById("mpi-sectores").value) || 5;
                const vueltas = parseInt(document.getElementById("mpi-vueltas").value) || 3;
                
                // Actualizar circuito
                currentSectors = sectores;
                setupCircuit(sectores);
                
                // Reiniciar estadísticas
                mpiStats = { bestTime: Infinity, sectorsProcessed: 0 };
                updateMpiStats();
                
                ws.send(JSON.stringify({action: "iniciar_mpi", sectores: sectores, vueltas: vueltas}));
                appendLog("mpi", "<b>Comando enviado: iniciar MPI</b>");
            });

            document.getElementById("start-openmp").addEventListener("click", () => {
                const autos = parseInt(document.getElementById("openmp-autos").value) || 4;
                const vueltas = parseInt(document.getElementById("openmp-vueltas").value) || 5;
                
                // Preparar autos
                setupCars(autos);
                
                // Reiniciar estadísticas
                openmpStats = { bestTime: Infinity, lapsCompleted: 0 };
                updateOpenmpStats();
                
                ws.send(JSON.stringify({action: "iniciar_openmp", autos: autos, vueltas: vueltas}));
                appendLog("openmp", "<b>Comando enviado: iniciar OpenMP</b>");
            });

            // Funciones auxiliares
            function appendLog(type, text) {
                const target = type === "mpi" ? mpiLog : 
                              type === "openmp" ? openmpLog : 
                              (mpiLog, openmpLog);
                
                const entry = document.createElement("div");
                entry.className = "log-entry " + type;
                entry.innerHTML = sanitize(text);
                
                if (type === "mpi") {
                    mpiLog.appendChild(entry);
                    mpiLog.scrollTop = mpiLog.scrollHeight;
                } else if (type === "openmp") {
                    openmpLog.appendChild(entry);
                    openmpLog.scrollTop = openmpLog.scrollHeight;
                } else {
                    mpiLog.appendChild(entry.cloneNode(true));
                    openmpLog.appendChild(entry);
                    mpiLog.scrollTop = mpiLog.scrollHeight;
                    openmpLog.scrollTop = openmpLog.scrollHeight;
                }
                
                // Efecto de highlight
                entry.classList.add("highlight");
                setTimeout(() => entry.classList.remove("highlight"), 1000);
            }

            function sanitize(s) {
                if (!s) return "";
                return s.replace(/</g, "&lt;").replace(/>/g, "&gt;");
            }

            function processMpiMessage(text) {
                // Extraer tiempo de sector del mensaje
                const timeMatch = text.match(/Tiempo de sector (\d+): (\d+\.\d+) s/);
                if (timeMatch) {
                    const sector = parseInt(timeMatch[1]);
                    const time = parseFloat(timeMatch[2]);
                    
                    // Actualizar estadísticas
                    if (time < mpiStats.bestTime) {
                        mpiStats.bestTime = time;
                    }
                    mpiStats.sectorsProcessed++;
                    updateMpiStats();
                    
                    // Resaltar el sector procesado
                    highlightSector(sector);
                    
                    // Mover el auto al siguiente sector
                    moveCarToSector(0, sector);
                }
            }

            function processOpenmpMessage(text) {
                // Extraer información de auto y vuelta
                const autoMatch = text.match(/Auto (\d+) - Vuelta (\d+): (\d+\.\d+) s/);
                const bestMatch = text.match(/Auto (\d+) - Nueva mejor vuelta: (\d+\.\d+) s/);
                
                if (autoMatch) {
                    const autoId = parseInt(autoMatch[1]);
                    const lap = parseInt(autoMatch[2]);
                    const time = parseFloat(autoMatch[3]);
                    
                    // Actualizar estadísticas
                    if (time < openmpStats.bestTime) {
                        openmpStats.bestTime = time;
                    }
                    openmpStats.lapsCompleted++;
                    updateOpenmpStats();
                    
                    // Mover el auto
                    moveCar(autoId, lap);
                } else if (bestMatch) {
                    const autoId = parseInt(bestMatch[1]);
                    const time = parseFloat(bestMatch[2]);
                    
                    // Actualizar mejor tiempo
                    openmpStats.bestTime = time;
                    updateOpenmpStats();
                    
                    // Destacar auto con mejor tiempo
                    highlightCar(autoId);
                }
            }

            function updateMpiStats() {
                bestMpiTime.textContent = mpiStats.bestTime !== Infinity ? 
                    mpiStats.bestTime.toFixed(2) + " s" : "-";
                sectorsProcessed.textContent = mpiStats.sectorsProcessed;
            }

            function updateOpenmpStats(data = null) {
                if (data && data.mejor_general) {
                    bestOpenmpTime.textContent = data.mejor_general.MejorVuelta.toFixed(2) + " s";
                } else {
                    bestOpenmpTime.textContent = openmpStats.bestTime !== Infinity ? 
                        openmpStats.bestTime.toFixed(2) + " s" : "-";
                }
                lapsCompleted.textContent = openmpStats.lapsCompleted;
            }

            function setupCircuit(sectors) {
                // Limpiar contenedores
                sectorsContainer.innerHTML = '';
                nodesContainer.innerHTML = '';
                
                // Calcular posiciones de sectores y nodos
                const centerX = document.querySelector('.circuit-track').offsetWidth / 2;
                const centerY = document.querySelector('.circuit-track').offsetHeight / 2;
                const radius = Math.min(centerX, centerY) - 50;
                
                nodePositions = [];
                
                for (let i = 0; i < sectors; i++) {
                    const angle = (2 * Math.PI * i) / sectors - Math.PI / 2;
                    const x = centerX + radius * Math.cos(angle);
                    const y = centerY + radius * Math.sin(angle);
                    
                    // Crear sector
                    const sector = document.createElement('div');
                    sector.className = 'sector';
                    sector.style.left = (x - 30) + "px";
                    sector.style.top = (y - 30) + "px";
                    sector.textContent = i + 1;
                    sectorsContainer.appendChild(sector);
                    
                    // Crear nodo (ligeramente más adentro)
                    const nodeRadius = radius - 40;
                    const nodeX = centerX + nodeRadius * Math.cos(angle);
                    const nodeY = centerY + nodeRadius * Math.sin(angle);
                    
                    const node = document.createElement('div');
                    node.className = 'sector-node';
                    node.style.left = (nodeX - 20) + "px";
                    node.style.top = (nodeY - 20) + "px";
                    node.textContent = i + 1;
                    nodesContainer.appendChild(node);
                    
                    nodePositions.push({ x: nodeX - 20, y: nodeY - 20 });
                }
                
                // Crear auto de F1 para MPI
                carsContainer.innerHTML = '';
                const car = document.createElement('div');
                car.className = 'f1-car';
                car.id = 'car-0';
                car.style.left = nodePositions[0].x + "px";
                car.style.top = nodePositions[0].y + "px";
                carsContainer.appendChild(car);
            }

            function setupCars(carCount) {
                // Limpiar autos existentes (excepto el de MPI)
                const mpiCar = document.getElementById('car-0');
                carsContainer.innerHTML = '';
                if (mpiCar) carsContainer.appendChild(mpiCar);
                
                carPositions = {};
                
                // Crear nuevos autos para OpenMP
                for (let i = 1; i <= carCount; i++) {
                    const car = document.createElement('div');
                    car.className = 'f1-car';
                    car.id = "car-" + i;
                    
                    // Posición inicial aleatoria en la pista
                    const randomPos = Math.floor(Math.random() * nodePositions.length);
                    car.style.left = nodePositions[randomPos].x + "px";
                    car.style.top = nodePositions[randomPos].y + "px";
                    
                    // Color diferente para cada auto
                    const hue = (i * 360 / carCount) % 360;
                    car.style.background = "hsl(" + hue + ", 100%, 50%)";
                    
                    carsContainer.appendChild(car);
                    carPositions[i] = randomPos;
                }
            }

            function moveCarToSector(carId, sector) {
                const car = document.getElementById("car-" + carId);
                if (!car || !nodePositions[sector - 1]) return;
                
                car.style.transition = 'all 0.5s ease-in-out';
                car.style.left = nodePositions[sector - 1].x + "px";
                car.style.top = nodePositions[sector - 1].y + "px";
            }

            function moveCar(carId, lap) {
                if (!carPositions[carId] && carPositions[carId] !== 0) return;
                
                // Mover al siguiente nodo (simulando progreso en la vuelta)
                carPositions[carId] = (carPositions[carId] + 1) % nodePositions.length;
                const car = document.getElementById("car-" + carId);
                
                if (car && nodePositions[carPositions[carId]]) {
                    car.style.transition = 'all 0.5s ease-in-out';
                    car.style.left = nodePositions[carPositions[carId]].x + "px";
                    car.style.top = nodePositions[carPositions[carId]].y + "px";
                }
            }

            function highlightSector(sector) {
                const sectorElements = document.getElementsByClassName('sector');
                const sectorElement = sectorElements[sector - 1];
                if (sectorElement) {
                    sectorElement.classList.add('pulse');
                    setTimeout(() => sectorElement.classList.remove('pulse'), 1500);
                }
            }

            function highlightCar(carId) {
                const car = document.getElementById("car-" + carId);
                if (car) {
                    car.classList.add('pulse');
                    setTimeout(() => car.classList.remove('pulse'), 1500);
                }
            }

            // Inicializar con circuito por defecto
            setupCircuit(5);
        });
    </script>
</body>
</html>`
