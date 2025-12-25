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
	"github.com/jgutierrez746/clase_7_gin_bun/config"
	"github.com/jgutierrez746/clase_7_gin_bun/db"
	"github.com/jgutierrez746/clase_7_gin_bun/dto"
	"golang.org/x/crypto/bcrypt"
)

func ConsultarUsuarios(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var usuarios []dto.UsuarioPerfilDTO

	u := config.Tablas["u"]
	p := config.Tablas["p"]

	var tablasJoin = []string{
		"JOIN " + p + " ON " + u + ".perfil_id = " + p + ".id",
	}

	var columnas = []string{
		u + ".id", u + ".nombre", u + ".correo", u + ".telefono", u + ".perfil_id", u + ".created_at", u + ".updated_at",
		p + ".nombre AS perfil_nombre",
	}

	if err := db.SelectConJoin(ctx, config.Tablas["u"], tablasJoin, columnas, &usuarios, u+".id DESC", ""); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error consultando usuarios: " + err.Error(),
		})
		return
	}

	if len(usuarios) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"usuarios": []interface{}{},
			"mensaje":  "No hay usuarios registrados",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"usuarios": usuarios,
		"total":    len(usuarios),
	})
}

func ConsultarUsuarioPorId(c *gin.Context) {
	id := c.Param("id")
	if strings.TrimSpace(id) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "No se ingresó parámetro solicitado.",
		})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var usuarios []dto.UsuarioPerfilDTO

	u := config.Tablas["u"]
	p := config.Tablas["p"]

	var tablasJoin = []string{
		"JOIN " + p + " ON " + u + ".perfil_id = " + p + ".id",
	}

	var columnas = []string{
		u + ".id", u + ".nombre", u + ".correo", u + ".telefono", u + ".perfil_id", u + ".created_at", u + ".updated_at",
		p + ".nombre AS perfil_nombre",
	}

	where := u + ".id = ?"

	if err := db.SelectConJoin(ctx, config.Tablas["u"], tablasJoin, columnas, &usuarios, "", where, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error consultando usuario: " + err.Error(),
		})
		return
	}

	if len(usuarios) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Usuario no encontrado",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"usuario": usuarios[0],
	})
}

func CrearUsuario(c *gin.Context) {
	var input dto.UsuarioInsert
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error de validación: " + err.Error(),
		})
		return
	}

	// Hashear password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), 8)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error procesando contraseña",
		})
		return
	}
	input.Password = string(hashedPassword)

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// Validar que exista el perfil
	var perfilDummy dto.PerfilesSelectDTO
	// Usamos SelectOne para validar existencia (retornará error si no existe)
	if err := db.SelectOne(ctx, config.Tablas["p"], &perfilDummy, "id = ?", input.PerfilID); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusBadRequest, gin.H{"error": "El PerfilID no existe."})
			return
		}
		// Otro error de base de datos
		log.Println("Error verificando perfil:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error verificando perfil"})
		return
	}

	if err := db.Insert(ctx, config.Tablas["u"], &input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error creando usuario: " + err.Error(),
		})
		return
	}

	// Limpiar password para respuesta
	input.Password = ""

	c.JSON(http.StatusCreated, gin.H{
		"mensaje": "Usuario creado exitosamente",
		"usuario": input,
	})
}

func EditarUsuario(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Parámetro inválido.",
		})
		return
	}

	var input dto.UsuarioUpdate
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error de validación: " + err.Error(),
		})
		return
	}

	if input.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), 8)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error procesando contraseña",
			})
			return
		}
		input.Password = string(hashedPassword)
	}

	input.UpdatedAt = time.Now().In(config.Chilelocation)
	input.ID = int64(id)

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	filasAfectadas, err := db.Update(ctx, config.Tablas["u"], &input, "id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error actualizando usuario: " + err.Error(),
		})
		return
	}

	if filasAfectadas == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Usuario no encontrado",
		})
		return
	}

	input.Password = "" // No devolver hash
	c.JSON(http.StatusOK, gin.H{
		"mensaje": "Usuario actualizado correctamente",
		"usuario": input,
	})
}

func EliminarUsuario(c *gin.Context) {
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

	filasAfectadas, err := db.Delete(ctx, config.Tablas["u"], "id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error eliminando usuario: " + err.Error(),
		})
		return
	}

	if filasAfectadas == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Usuario no encontrado",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"mensaje":    "Usuario eliminado correctamente",
		"eliminados": filasAfectadas,
	})
}
