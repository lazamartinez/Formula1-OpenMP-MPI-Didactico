package handlers

import (
	"formula1-crud-go/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ManejadorPilotos struct {
	DB *gorm.DB
}

func NuevoManejadorPilotos(db *gorm.DB) *ManejadorPilotos {
	return &ManejadorPilotos{DB: db}
}

// Obtener todos los pilotos
func (m *ManejadorPilotos) ObtenerPilotos(c *gin.Context) {
	var pilotos []models.Piloto
	resultado := m.DB.Find(&pilotos)
	if resultado.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": resultado.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, pilotos)
}

// Obtener un piloto por ID
func (m *ManejadorPilotos) ObtenerPiloto(c *gin.Context) {
	id := c.Param("id")
	var piloto models.Piloto
	resultado := m.DB.First(&piloto, id)
	if resultado.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Piloto no encontrado"})
		return
	}
	c.JSON(http.StatusOK, piloto)
}

// Crear nuevo piloto
func (m *ManejadorPilotos) CrearPiloto(c *gin.Context) {
	var piloto models.Piloto
	if err := c.ShouldBindJSON(&piloto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resultado := m.DB.Create(&piloto)
	if resultado.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": resultado.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, piloto)
}

// Actualizar piloto
func (m *ManejadorPilotos) ActualizarPiloto(c *gin.Context) {
	id := c.Param("id")
	var piloto models.Piloto

	if err := m.DB.First(&piloto, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Piloto no encontrado"})
		return
	}

	if err := c.ShouldBindJSON(&piloto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	m.DB.Save(&piloto)
	c.JSON(http.StatusOK, piloto)
}

// Eliminar piloto
func (m *ManejadorPilotos) EliminarPiloto(c *gin.Context) {
	id := c.Param("id")
	var piloto models.Piloto

	if err := m.DB.First(&piloto, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Piloto no encontrado"})
		return
	}

	m.DB.Delete(&piloto)
	c.JSON(http.StatusOK, gin.H{"mensaje": "Piloto eliminado correctamente"})
}

// Obtener estadísticas
func (m *ManejadorPilotos) ObtenerEstadisticas(c *gin.Context) {
	var totalPilotos int64
	var totalVictorias int64
	var totalPuntos float64
	var pilotoMasVictorias models.Piloto

	m.DB.Model(&models.Piloto{}).Count(&totalPilotos)
	m.DB.Model(&models.Piloto{}).Select("SUM(victorias)").Row().Scan(&totalVictorias)
	m.DB.Model(&models.Piloto{}).Select("SUM(puntos)").Row().Scan(&totalPuntos)
	m.DB.Order("victorias desc").First(&pilotoMasVictorias)

	estadisticas := gin.H{
		"total_pilotos":        totalPilotos,
		"total_victorias":      totalVictorias,
		"total_puntos":         totalPuntos,
		"piloto_mas_victorias": pilotoMasVictorias.Nombre,
		"victorias_piloto":     pilotoMasVictorias.Victorias,
	}

	c.JSON(http.StatusOK, estadisticas)
}

// Buscar pilotos por equipo
func (m *ManejadorPilotos) BuscarPorEquipo(c *gin.Context) {
	equipo := c.Query("equipo")
	var pilotos []models.Piloto

	if equipo == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parámetro 'equipo' requerido"})
		return
	}

	resultado := m.DB.Where("equipo LIKE ?", "%"+equipo+"%").Find(&pilotos)
	if resultado.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": resultado.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, pilotos)
}
