package entity

import (
	"time"

	"github.com/google/uuid"
)

type Obra struct {
	ID           uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Nome         string    `json:"nome" gorm:"type:varchar(255);not null"`
	Endereco     string    `json:"endereco" gorm:"type:varchar(255);not null"`
	Progresso    float64   `json:"progresso" gorm:"type:float;default:0"`
	DataInicio   time.Time `json:"data_inicio" gorm:"type:date;not null"`
	DataEstimada time.Time `json:"data_estimada" gorm:"type:date;not null"`
	GastoTotal   float64   `json:"gasto_total" gorm:"type:float;default:0"`
	Entradas     []Entrada `json:"entradas,omitempty" gorm:"foreignKey:ObraID"`
}
