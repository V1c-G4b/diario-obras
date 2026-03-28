package entity

import "github.com/google/uuid"

type Responsavel struct {
	ID    uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Nome  string    `json:"nome" gorm:"type:varchar(255);not null"`
	Email string    `json:"email" gorm:"type:varchar(255);not null"`
	Cargo string    `json:"cargo" gorm:"type:varchar(255);not null"`
}
