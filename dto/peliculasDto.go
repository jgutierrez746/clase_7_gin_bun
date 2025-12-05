package dto

import "time"

type PeliculaSelectDTO struct {
	ID          int64     `json:"id"`
	Anio        int       `json:"anio"`
	Titulo      string    `json:"titulo"`
	Slug        string    `json:"slug"`
	Descripcion string    `json:"descripcion"`
	Director    string    `json:"director"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type PeliculasAllSelect []PeliculaSelectDTO

type PeliculaInsert struct {
	ID          int64     `json:"id,omitempty" bun:",pk,autoincrement"`
	Anio        int       `json:"anio" binding:"required,min=1000,max=9999"`
	Titulo      string    `json:"titulo" binding:"required"`
	Slug        string    `json:"slug,omitempty"`
	Descripcion string    `json:"descripcion" binding:"required"`
	Director    string    `json:"director" binding:"required"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}

type PeliculaUpdate struct {
	ID          int64     `json:"id,omitempty" bun:",pk,autoincrement"`
	Anio        int       `json:"anio,omitempty" binding:"required,min=1000,max=9999"`
	Titulo      string    `json:"titulo,omitempty" binding:"required"`
	Slug        string    `json:"slug,omitempty"` // Este campo no se usa en el JSON, es solo para entregar la información en la respuesta
	Descripcion string    `json:"descripcion,omitempty" binding:"required"`
	Director    string    `json:"director,omitempty" binding:"required"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}
