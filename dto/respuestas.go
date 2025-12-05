package dto

import "time"

type PeliculaJoinRow struct {
	// Película
	ID          int64
	Anio        int
	Titulo      string
	Slug        string
	Descripcion string
	Director    string
	CreatedAt   time.Time
	UpdatedAt   time.Time

	// Temática
	TematicaID int64
	Nombre     string
	SlugTem    string
	TCreatedAt time.Time
	TUpdatedAt time.Time

	// Pelicula_Tematicas
	Orden int
}

type TematicasPeliculaJoinRow struct {
	// Temática
	Nombre string

	// Pelicula_Tematicas
	Orden int
}
