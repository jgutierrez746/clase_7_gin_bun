package dto

import (
	"time"
)

type TematicasSelectOne struct {
	ID        int64     `json:"id"`
	Nombre    string    `json:"nombre"`
	Slug      string    `json:"slug"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type TemticasAllSelect []TematicasSelectOne

type TematicasInsert struct {
	ID        int64     `json:"id,omitempty" bun:",pk,autoincrement"`
	Nombre    string    `json:"nombre" binding:"required"`
	Slug      string    `json:"slug,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

type TematicasBatchInsert []TematicasInsert

type TematicasUpdate struct {
	ID        int64     `json:"id,omitempty" bun:",pk,autoincrement"`
	Nombre    string    `json:"nombre" binding:"required"`
	Slug      string    `json:"slug,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
