package entity

import (
	"time"

	"github.com/google/uuid"
)

type Foto struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	EntradaID uuid.UUID `json:"entrada_id" gorm:"type:uuid;not null"`
	URLS3     string    `json:"url_s3" gorm:"type:varchar(500);not null"`
	Descricao string    `json:"descricao" gorm:"type:varchar(500)"`
	CriadoEm  time.Time `json:"criado_em" gorm:"autoCreateTime"`
}
