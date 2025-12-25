package rutas

import (
	"context"
	"database/sql"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jgutierrez746/clase_7_gin_bun/config"
	"github.com/jgutierrez746/clase_7_gin_bun/db"
	"github.com/jgutierrez746/clase_7_gin_bun/dto"
)

func ConsultarPerfiles(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var perfiles dto.PerfilesAllSelect
	if err := db.SelectAll(ctx, config.Tablas["p"], &perfiles); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error consultando perfiles: " + err.Error(),
		})
		return
	}

	if len(perfiles) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"perfiles": []interface{}{},
			"mensaje":  "No hay perfiles registrados",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"perfiles": perfiles,
		"total":    len(perfiles),
	})
}

func ConsultarPerfilPorId(c *gin.Context) {
	id := c.Param("id")
	if strings.TrimSpace(id) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "No se ingresó parámetro solicitado.",
		})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var perfil dto.PerfilesSelectDTO
	if err := db.SelectOne(ctx, config.Tablas["p"], &perfil, "id = ?", id); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Perfil no encontrado",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error consultando perfil: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"perfil": perfil,
	})
}

func CrearPerfil(c *gin.Context) {
	var input dto.PerfilesInsert
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error al procesar el JSON: " + err.Error(),
		})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	if err := db.Insert(ctx, config.Tablas["p"], &input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"mensaje": "Perfil creado exitosamente",
		"perfil":  input,
	})
}

func EditarPerfil(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Parámetro inválido.",
		})
		return
	}

	var input dto.PerfilesUpdate
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error al procesar el JSON: " + err.Error(),
		})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	input.ID = int64(id)

	filasAfectadas, err := db.Update(ctx, config.Tablas["p"], &input, "id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error al actualizar perfil: " + err.Error(),
		})
		return
	}

	if filasAfectadas == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Perfil no encontrado",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"mensaje": "Perfil actualizado correctamente",
		"perfil":  input,
	})
}

func EliminarPerfil(c *gin.Context) {
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

	filasAfectadas, err := db.Delete(ctx, config.Tablas["p"], "id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error eliminando perfil: " + err.Error(),
		})
		return
	}

	if filasAfectadas == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Perfil no encontrado",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"mensaje":    "Perfil eliminado correctamente",
		"eliminados": filasAfectadas,
	})
}
