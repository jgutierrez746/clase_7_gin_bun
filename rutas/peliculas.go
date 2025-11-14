package rutas

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jgutierrez746/clase_7_gin_bun/config"
	"github.com/jgutierrez746/clase_7_gin_bun/db"
	"github.com/jgutierrez746/clase_7_gin_bun/dto"
	"github.com/jgutierrez746/clase_7_gin_bun/helpers"
)

func ConsultarPeliculas(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var modelo []dto.PeliculaJoinRow

	tm := config.Tablas["tm"]
	pl := config.Tablas["pl"]
	pt := config.Tablas["pt"]

	var tablasJoin = []string{
		fmt.Sprintf("LEFT JOIN %s ON %s.p_id = %s.id", pt, pt, pl),
		fmt.Sprintf("LEFT JOIN %s ON %s.tematica_id = %s.id", tm, pt, tm),
	}

	var columnas = []string{
		fmt.Sprintf("%s.id, %s.anio, %s.titulo, %s.slug, %s.descripcion, %s.created_at, %s.updated_at", pl, pl, pl, pl, pl, pl, pl),
		fmt.Sprintf("%s.id AS tematica_id, %s.nombre, %s.slug AS slug_tem, %s.created_at AS t_created_at, %s.updated_at AS t_updated_at", tm, tm, tm, tm, tm),
		fmt.Sprintf("%s.orden", pt),
	}

	order := fmt.Sprintf("%s.id DESC, %s.orden ASC", pl, pt)

	if err := db.SelectConJoin(ctx, config.Tablas["pl"], tablasJoin, columnas, &modelo, order, "", nil); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error consultando películas: " + err.Error(),
		})
		return
	}

	if len(modelo) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"peliculas": []interface{}{},
			"mensaje":   "No hay películas registradas",
		})
		return
	}

	peliculas := helpers.AgruparPeliculas(modelo)

	c.JSON(http.StatusOK, gin.H{
		"peliculas": peliculas,
		"total":     len(peliculas),
	})
}

func ConsultarPeliculaPorId(c *gin.Context) {
	id := c.Param("id")
	if strings.TrimSpace(id) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "No se ingresó parámetro solicitado.",
		})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var modelo []dto.PeliculaJoinRow

	tm := config.Tablas["tm"]
	pl := config.Tablas["pl"]
	pt := config.Tablas["pt"]

	var tablasJoin = []string{
		fmt.Sprintf("LEFT JOIN %s ON %s.p_id = %s.id", pt, pt, pl),
		fmt.Sprintf("LEFT JOIN %s ON %s.tematica_id = %s.id", tm, pt, tm),
	}

	var columnas = []string{
		fmt.Sprintf("%s.id, %s.anio, %s.titulo, %s.slug, %s.descripcion, %s.created_at, %s.updated_at", pl, pl, pl, pl, pl, pl, pl),
		fmt.Sprintf("%s.id AS tematica_id, %s.nombre, %s.slug AS slug_tem, %s.created_at AS t_created_at, %s.updated_at AS t_updated_at", tm, tm, tm, tm, tm),
		fmt.Sprintf("%s.orden", pt),
	}

	where := fmt.Sprintf("%s.id = ?", pl)

	order := fmt.Sprintf("%s.id DESC, %s.orden ASC", pl, pt)

	if err := db.SelectConJoin(ctx, config.Tablas["pl"], tablasJoin, columnas, &modelo, order, where, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error consultando películas: " + err.Error(),
		})
		return
	}

	peliculas := helpers.AgruparPeliculas(modelo)

	if len(peliculas) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Película no encontrada",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"pelicula": peliculas[0],
	})
}
