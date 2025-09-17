package models

import (
    "gorm.io/gorm"
)

type Piloto struct {
    gorm.Model
    Nombre        string  `json:"nombre" binding:"required"`
    Equipo        string  `json:"equipo" binding:"required"`
    Nacionalidad  string  `json:"nacionalidad"`
    Numero        int     `json:"numero"`
    Victorias     int     `json:"victorias"`
    Puntos        float64 `json:"puntos"`
    Podios        int     `json:"podios"`
    Poles         int     `json:"poles"`
    VueltasRapidas int    `json:"vueltas_rapidas"`
}
