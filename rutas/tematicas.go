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

func ConsultarTematicas(c *gin.Context) {
	// Definir contexto con tiempo de espera de solo 5 segundos
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// tematicas := new(dto.AllTematicasSalida)a
	var tematicas dto.TemticasAllSelect
	if err := db.SelectAll(ctx, config.Tablas["tm"], &tematicas); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error consultando temáticas: " + err.Error(),
		})
		return
	}

	// Si no hay datos, responde vacío
	if len(tematicas) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"tematicas": []interface{}{},
			"mensaje":   "No hay temáticas registradas",
		})
		return
	}

	log.Printf("Se consultaron %d temáticas.", len(tematicas))
	c.JSON(http.StatusOK, gin.H{
		"tematicas": tematicas, // JSON con todos los campos (ID, Nombre, Slug)
		"total":     len(tematicas),
	})
}

func ConsultarTematicasPorId(c *gin.Context) {
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
	var tematica dto.TematicasSelectOne
	if err := db.SelectOne(ctx, config.Tablas["tm"], &tematica, "id = ?", id); err != nil {
		if err == sql.ErrNoRows { // Si no existe, 404
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Temática no encontrada",
			})
			return
		}
		log.Printf("Error en SelectOne: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error consultando temática: " + err.Error(),
		})
		return
	}

	log.Printf("Se consultó temática con ID: %s ", id)
	c.JSON(http.StatusOK, gin.H{
		"tematica": tematica, // JSON con todos los campos (ID, Nombre, Slug...)
	})
}

func CrearTematica(c *gin.Context) {
	var tematica dto.TematicasInsert
	if err := c.ShouldBindJSON(&tematica); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error al procesar el JSON " + err.Error(),
		})
		return
	}

	nowChile := time.Now().In(config.Chilelocation)
	tematica.Slug = slug.Make(tematica.Nombre)
	tematica.CreatedAt = nowChile
	tematica.UpdatedAt = nowChile

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	if err := db.Insert(ctx, config.Tablas["tm"], &tematica); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Se envía respuesta con el modelo actualizado
	c.JSON(http.StatusCreated, gin.H{
		"mensaje":  "Temática creada en Base de Datos",
		"tematica": tematica,
	})
}

func EditarTematica(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Parámetro inválido.",
		})
		return
	}

	var input dto.TematicasUpdate
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error al procesar el JSON " + err.Error(),
		})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	input.ID = int64(id)
	input.Slug = slug.Make(input.Nombre)
	input.UpdatedAt = time.Now().In(config.Chilelocation)

	// Ejecutamos el Update
	filasAfectadas, err := db.Update(ctx, config.Tablas["tm"], &input, "id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error en update: " + err.Error(),
		})
		return
	}

	if filasAfectadas == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Temática no encontrada",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"mensaje":  "Temática editada correctamente",
		"tematica": input,
	})
}

func EliminarTematica(c *gin.Context) {
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
	filasAfectadas, err := db.Delete(ctx, config.Tablas["tm"], "id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error eliminando: " + err.Error(),
		})
		return
	}

	if filasAfectadas == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Temática no encontrada",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"mensaje":    "Temática eliminada correctamente",
		"eliminados": filasAfectadas,
	})
}
