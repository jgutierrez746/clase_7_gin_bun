package modelos

import (
	"time"

	"github.com/uptrace/bun"
)

type TematicasModel struct {
	bun.BaseModel `bun:"table:tematicas"`

	ID        int64     `bun:",pk,autoincrement"`
	Nombre    string    `bun:",type:varchar(100),notnull"`
	Slug      string    `bun:",type:varchar(100),notnull,unique"`
	CreatedAt time.Time `bun:",type:timestamp,default:current_timestamp"`
	UpdatedAt time.Time `bun:",type:timestamp,default:current_timestamp,on_update:current_timestamp"`
}

type PeliculasModel struct {
	bun.BaseModel `bun:"table:peliculas"`

	ID          int64     `bun:",pk,autoincrement"`
	Anio        int       `bun:",type:smallint,nullzero"`
	Titulo      string    `bun:",type:varchar(255),notnull"`
	Slug        string    `bun:",type:varchar(255),notnull,unique"`
	Descripcion string    `bun:",type:text"`
	Director    string    `bun:",type:varchar(100),notnull"`
	CreatedAt   time.Time `bun:",type:timestamp,default:current_timestamp"`
	UpdatedAt   time.Time `bun:",type:timestamp,default:current_timestamp,on_update:current_timestamp"`
}

type PeliculaTematicaModel struct {
	bun.BaseModel `bun:"table:pelicula_tematicas"`

	PID        int64     `bun:"p_id,pk"`        // FK a Peliculas.ID
	TematicaID int64     `bun:"tematica_id,pk"` // FK a Tematicas.ID
	Orden      int       `bun:",nullzero"`
	CreatedAt  time.Time `bun:",type:timestamp,default:current_timestamp"`
}

type PortadaPeliculaModel struct {
	bun.BaseModel `bun:"table:portada_pelicula"`

	ID            int64     `bun:",pk,autoincrement"`
	PID           int64     `bun:"p_id"`
	NombreArchivo string    `bun:"nombre_archivo"`
	CreatedAt     time.Time `bun:",type:timestamp,default:current_timestamp"`
}
