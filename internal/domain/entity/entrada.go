package entity

import (
	"time"

	"github.com/google/uuid"
)

type Entrada struct {
	ID                      uuid.UUID         `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ObraID                  uuid.UUID         `json:"obra_id" gorm:"type:uuid;not null"`
	ResponsavelID           uuid.UUID         `json:"responsavel_id" gorm:"type:uuid;not null"`
	Observacoes             string            `json:"observacoes" gorm:"type:text"`
	Etapa                   string            `json:"etapa" gorm:"type:varchar(255);not null"`
	ProgressoEtapa          float64           `json:"progresso_etapa" gorm:"type:float;default:0"`
	CriadoEm                time.Time         `json:"criado_em" gorm:"autoCreateTime"`
	CustoDia                float64           `json:"custo_dia" gorm:"type:float;not null"`
	QuantidadeTrabalhadores int64             `json:"quantidade_trabalhadores" gorm:"type:int;not null"`
	CondicoesClimaticas     CondicaoClimatica `json:"condicoes_climaticas" gorm:"type:int;not null;default:0"`
	Paralisacao             bool              `json:"paralisacao" gorm:"type:bool;default:false"`
	Obra                    Obra              `json:"-" gorm:"foreignKey:ObraID"`
	Responsavel             Responsavel       `json:"responsavel,omitempty" gorm:"foreignKey:ResponsavelID"`
	Fotos                   []Foto            `json:"fotos,omitempty" gorm:"foreignKey:EntradaID"`
}

type CondicaoClimatica int

const (
	Ensolarado CondicaoClimatica = iota // 0
	Nublado                             // 1
	Chuvoso                             // 2
	Tempestade                          // 3
)
