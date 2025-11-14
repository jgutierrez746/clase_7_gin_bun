package dto

import "time"

type PeliculaSelectDTO struct {
	ID          int64                `json:"id"`
	Anio        int                  `json:"anio"`
	Titulo      string               `json:"titulo"`
	Slug        string               `json:"slug"`
	Descripcion string               `json:"descripcion"`
	CreatedAt   time.Time            `json:"created_at"`
	UpdatedAt   time.Time            `json:"updated_at"`
	Tematicas   []TematicasSelectOne `json:"tematicas"`
}
