package rutas

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
	"github.com/jgutierrez746/clase_7_gin_bun/config"
	"github.com/jgutierrez746/clase_7_gin_bun/db"
	"github.com/jgutierrez746/clase_7_gin_bun/dto"
)

func ConsultarPeliculas(c *gin.Context) {
	// Definir contexto con tiempo de espera de solo 5 segundos
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var peliculas dto.PeliculasAllSelect
	if err := db.SelectAll(ctx, config.Tablas["pl"], &peliculas); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error consultando películas: " + err.Error(),
		})
		return
	}

	// Si no hay datos, responde vacío
	if len(peliculas) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"tematicas": []interface{}{},
			"mensaje":   "No hay películas registradas",
		})
		return
	}

	log.Printf("Se consultaron %d películas.", len(peliculas))
	c.JSON(http.StatusOK, gin.H{
		"peliculas": peliculas, // JSON con todos los campos (ID, Nombre, Slug)
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
	// Definir contexto con tiempo de espera de solo 5 segundos
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// tematicas := new(dto.AllTematicasSalida)a
	var pelicula dto.PeliculaSelectDTO
	if err := db.SelectOne(ctx, config.Tablas["pl"], &pelicula, "id = ?", id); err != nil {
		if err == sql.ErrNoRows { // Si no existe, 404
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Película no encontrada",
			})
			return
		}
		log.Printf("Error en SelectOne: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error consultando película: " + err.Error(),
		})
		return
	}

	log.Printf("Se consultó temática con ID: %s ", id)
	c.JSON(http.StatusOK, gin.H{
		"película": pelicula, // JSON con todos los campos (ID, Nombre, Slug...)
	})
}

/*
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
		fmt.Sprintf("%s.id, %s.anio, %s.titulo, %s.slug, %s.descripcion, %s.director, %s.created_at, %s.updated_at", pl, pl, pl, pl, pl, pl, pl, pl),
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
*/
/*
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
		fmt.Sprintf("%s.id, %s.anio, %s.titulo, %s.slug, %s.descripcion, %s.director, %s.created_at, %s.updated_at", pl, pl, pl, pl, pl, pl, pl, pl),
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
*/

func CrearPelicula(c *gin.Context) {
	var pelicula dto.PeliculaInsert
	if err := c.ShouldBindJSON(&pelicula); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error al procesar el JSON " + err.Error(),
		})
		return
	}

	nowChile := time.Now().In(config.Chilelocation)
	pelicula.Slug = slug.Make(pelicula.Titulo)
	pelicula.CreatedAt = nowChile
	pelicula.UpdatedAt = nowChile

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	if err := db.Insert(ctx, config.Tablas["pl"], &pelicula); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Se envía respuesta con el modelo actualizado
	c.JSON(http.StatusCreated, gin.H{
		"mensaje":  "Pelicula creada en Base de Datos",
		"pelicula": pelicula,
	})
}

func EditarPelicula(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Parámetro inválido.",
		})
		return
	}

	var input dto.PeliculaUpdate
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error al procesar el JSON " + err.Error(),
		})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	input.ID = int64(id)
	if input.Titulo != "" {
		input.Slug = slug.Make(input.Titulo)
	}
	input.UpdatedAt = time.Now().In(config.Chilelocation)

	// Ejecutamos el Update
	filasAfectadas, err := db.Update(ctx, config.Tablas["pl"], &input, "id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error en update: " + err.Error(),
		})
		return
	}

	if filasAfectadas == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Película no encontrada",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"mensaje":  "Película editada correctamente",
		"pelicula": input,
	})
}

func EliminarPelicula(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Parámetro inválido.",
		})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// Ejecutamos Delete
	filasAfectadas, err := db.Delete(ctx, config.Tablas["pl"], "id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error eliminando: " + err.Error(),
		})
		return
	}

	if filasAfectadas == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Película no encontrada",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"mensaje":    "Película eliminada correctamente",
		"eliminados": filasAfectadas,
	})
}
