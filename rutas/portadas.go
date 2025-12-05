package rutas

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jgutierrez746/clase_7_gin_bun/config"
	"github.com/jgutierrez746/clase_7_gin_bun/db"
	"github.com/jgutierrez746/clase_7_gin_bun/dto"
)

func ConsultarPortadasPelicula(c *gin.Context) {
	id := c.Param("id")
	if strings.TrimSpace(id) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "No se ingresó parámetro solicitado.",
		})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var portada dto.PortadaSelectDTO
	// Usamos SelectOne para traer una sola portada
	if err := db.SelectOne(ctx, config.Tablas["pp"], &portada, "p_id = ?", id); err != nil {
		fmt.Printf("Error buscando portada para p_id=%s: %v\n", id, err) // Debug log
		c.JSON(http.StatusOK, gin.H{
			"mensaje": "No hay portada registrada para esta película",
		})
		return
	}

	// Construir URL completa
	// Asumimos que el host es el mismo del request, o se puede configurar en .env
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	portada.Url = fmt.Sprintf("%s://%s/imagenes/%s", scheme, c.Request.Host, portada.NombreArchivo)

	c.JSON(http.StatusOK, gin.H{
		"portada": portada,
	})
}

func CrearPortada(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Parámetro inválido.",
		})
		return
	}

	// 1. Recibir archivo
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "No se subió ningún archivo: " + err.Error(),
		})
		return
	}

	// 2. Definir ruta y nombre
	ext := filepath.Ext(file.Filename)
	nuevoNombre := fmt.Sprintf("%d_%d%s", id, time.Now().UnixNano(), ext)
	// Cambio de directorio a "portadas"
	rutaDestino := filepath.Join("public", "upload", "portadas", nuevoNombre)

	// Asegurar que directorio existe
	if err := os.MkdirAll(filepath.Dir(rutaDestino), 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error creando directorio: " + err.Error(),
		})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// 3. Verificar si ya existe portada para esta película y borrarla (regla de única portada)
	var portadaExistente dto.PortadaSelectDTO
	if err := db.SelectOne(ctx, config.Tablas["pp"], &portadaExistente, "p_id = ?", id); err == nil {
		// Existe, borrar archivo viejo
		rutaVieja := filepath.Join("public", "upload", "portadas", portadaExistente.NombreArchivo)
		os.Remove(rutaVieja)
		// Borrar registro viejo
		db.Delete(ctx, config.Tablas["pp"], "id = ?", portadaExistente.ID)
	}

	// 4. Guardar archivo nuevo
	if err := c.SaveUploadedFile(file, rutaDestino); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error guardando archivo: " + err.Error(),
		})
		return
	}

	// 5. Guardar en BD
	portada := dto.PortadaInsertDTO{
		PID:           int64(id),
		NombreArchivo: nuevoNombre,
		CreatedAt:     time.Now().In(config.Chilelocation),
	}

	if err := db.Insert(ctx, config.Tablas["pp"], &portada); err != nil {
		os.Remove(rutaDestino)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error registrando en BD: " + err.Error(),
		})
		return
	}

	// Construir URL para respuesta
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	url := fmt.Sprintf("%s://%s/imagenes/%s", scheme, c.Request.Host, nuevoNombre)

	c.JSON(http.StatusCreated, gin.H{
		"mensaje": "Portada subida correctamente",
		"url":     url,
	})
}

func EliminarPortada(c *gin.Context) {
	// idf: id de la foto/portada
	idf := c.Param("idf")

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// 1. Obtener nombre archivo
	var portada dto.PortadaSelectDTO
	if err := db.SelectOne(ctx, config.Tablas["pp"], &portada, "id = ?", idf); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Portada no encontrada",
		})
		return
	}

	// 2. Borrar de BD
	filas, err := db.Delete(ctx, config.Tablas["pp"], "id = ?", idf)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error eliminando de BD: " + err.Error(),
		})
		return
	}

	if filas == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Portada no encontrada al intentar borrar",
		})
		return
	}

	// 3. Borrar archivo físico (ruta actualizada)
	rutaArchivo := filepath.Join("public", "upload", "portadas", portada.NombreArchivo)
	if err := os.Remove(rutaArchivo); err != nil {
		fmt.Printf("Error borrando archivo físico %s: %v\n", rutaArchivo, err)
	}

	c.JSON(http.StatusOK, gin.H{
		"mensaje": "Portada eliminada correctamente",
	})
}
