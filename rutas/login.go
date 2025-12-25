package rutas

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jgutierrez746/clase_7_gin_bun/config"
	"github.com/jgutierrez746/clase_7_gin_bun/db"
	"github.com/jgutierrez746/clase_7_gin_bun/dto"
	jwtPkg "github.com/jgutierrez746/clase_7_gin_bun/jwt"
	"golang.org/x/crypto/bcrypt"
)

func Login(c *gin.Context) {
	var input dto.LoginDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error de validación: " + err.Error(),
		})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// Nota: UsuarioInsert no tiene los tags bun:"table", pero db.SelectOne usa el nombre de tabla pasado
	// Sin embargo, UsuarioInsert tiene password, que es lo que necesitamos.
	// Debemos asegurar que el scan funcione. UsuarioInsert tiene tags json, bun lo inferirá si no hay tags struct
	// Para mayor seguridad usamos un struct ad-hoc o reutilizamos uno que tenga los campos db necesarios
	type UsuarioLogin struct {
		ID       int64  `bun:"id"`
		Password string `bun:"password"`
		PerfilID int64  `bun:"perfil_id"`
	}
	var userDB UsuarioLogin

	if err := db.SelectOne(ctx, config.Tablas["u"], &userDB, "correo = ?", input.Correo); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciales inválidas"}) // Usuario no encontrado
			return
		}
		log.Println("Error buscando usuario:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno"})
		return
	}

	// Verificar password
	if err := bcrypt.CompareHashAndPassword([]byte(userDB.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciales inválidas"}) // Password incorrecto
		return
	}

	// Generar Token
	token, err := jwtPkg.GenerateToken(userDB.ID, userDB.PerfilID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo generar el token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}
