package helpers

import (
	"github.com/jgutierrez746/clase_7_gin_bun/dto"
)

func AgruparPeliculas(rows []dto.PeliculaJoinRow) []dto.PeliculaSelectDTO {
	peliculasMap := make(map[int64]*dto.PeliculaSelectDTO)

	for _, r := range rows {
		p, exists := peliculasMap[r.ID]
		if !exists {
			p = &dto.PeliculaSelectDTO{
				ID:          r.ID,
				Anio:        r.Anio,
				Titulo:      r.Titulo,
				Slug:        r.Slug,
				Descripcion: r.Descripcion,
				CreatedAt:   r.CreatedAt,
				UpdatedAt:   r.UpdatedAt,
				Tematicas:   []dto.TematicasSelectOne{},
			}
			peliculasMap[r.ID] = p
		}

		if r.TematicaID != 0 {
			p.Tematicas = append(p.Tematicas, dto.TematicasSelectOne{
				ID:        r.TematicaID,
				Nombre:    r.Nombre,
				Slug:      r.SlugTem,
				CreatedAt: r.TCreatedAt,
				UpdatedAt: r.TUpdatedAt,
			})
		}
	}

	var peliculas []dto.PeliculaSelectDTO
	for _, p := range peliculasMap {
		peliculas = append(peliculas, *p)
	}
	return peliculas
}
