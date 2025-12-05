package rutas

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jgutierrez746/clase_7_gin_bun/config"
	"github.com/jgutierrez746/clase_7_gin_bun/db"
	"github.com/jgutierrez746/clase_7_gin_bun/dto"
)

func ConsultarTematicasPelicula(c *gin.Context) {
	id := c.Param("id")
	if strings.TrimSpace(id) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "No se ingresó parámetro solicitado.",
		})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var modelo []dto.TematicasPeliculaJoinRow

	pt := config.Tablas["pt"]
	tm := config.Tablas["tm"]

	var tablasJoin = []string{
		fmt.Sprintf("LEFT JOIN %s ON (%s.tematica_id = %s.id)", tm, pt, tm),
	}

	var columnas = []string{
		fmt.Sprintf("%s.nombre, %s.orden", tm, pt),
	}

	where := fmt.Sprintf("%s.p_id = ?", pt)

	order := fmt.Sprintf("%s.orden ASC", pt)

	if err := db.SelectConJoin(ctx, config.Tablas["pt"], tablasJoin, columnas, &modelo, order, where, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error consultando temáticas asociadas: " + err.Error(),
		})
		return
	}

	if len(modelo) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"peliculas": []interface{}{},
			"mensaje":   "No hay temáticas registradas para esta película",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tematicas_asociadas": modelo,
		"total":               len(modelo),
	})
}

func CrearTematicasPelicula(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Parámetro inválido.",
		})
		return
	}

	var peliculaTematicas []dto.PeliculaTematicasInsert
	if err := c.ShouldBindJSON(&peliculaTematicas); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error al procesar el JSON " + err.Error(),
		})
		return
	}

	nowChile := time.Now().In(config.Chilelocation)

	for indice := range peliculaTematicas {
		peliculaTematicas[indice].PID = int64(id)
		peliculaTematicas[indice].CreatedAt = nowChile
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	insertados, err := db.InsertBatch(ctx, config.Tablas["pt"], peliculaTematicas)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Se envía respuesta con el modelo actualizado
	c.JSON(http.StatusCreated, gin.H{
		"mensaje":    "Temáticas asociadas se registraron correctamente",
		"insertados": insertados,
	})
}

func EliminarTematicaPelicula(c *gin.Context) {
	id := c.Param("id")
	idt := c.Param("idt")
	if strings.TrimSpace(id) == "" || strings.TrimSpace(idt) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "No se ingresó uno o más parámetros solicitados.",
		})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	where := "p_id = ? AND tematica_id = ?"

	// Ejecutamos Delete
	filasAfectadas, err := db.Delete(ctx, config.Tablas["pt"], where, id, idt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error eliminando temática asociada: " + err.Error(),
		})
		return
	}

	if filasAfectadas == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Temática asociada no encontrada",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"mensaje":    "Temática asociada eliminada correctamente",
		"eliminados": filasAfectadas,
	})
}
