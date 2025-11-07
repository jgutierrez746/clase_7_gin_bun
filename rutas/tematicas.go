package rutas

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jgutierrez746/clase_7_gin_bun/db"
	"github.com/jgutierrez746/clase_7_gin_bun/modelos"
)

func ConsultarTematicas(c *gin.Context) {
	// Definir contexto con tiempo de espera de solo 5 segundos
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// tematicas := new(dto.AllTematicasSalida)a
	var tematicas []modelos.TematicasModel
	if err := db.SelectAll(ctx, "", &tematicas); err != nil {
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
	// Definir contexto con tiempo de espera de solo 5 segundos
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// tematicas := new(dto.AllTematicasSalida)a
	var tematica modelos.TematicasModel
	if err := db.SelectOne(ctx, "", &tematica, "id = ?", id); err != nil {
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
	var tematica modelos.TematicasModel
	if err := c.ShouldBindJSON(&tematica); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error al procesar el JSON " + err.Error(),
		})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	if err := db.Insert(ctx, &tematica); err != nil {
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
