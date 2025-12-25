package jwt

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte(getSecret())

func getSecret() string {
	s := os.Getenv("JWT_SECRET")
	if s == "" {
		return "secreto_por_defecto_muy_seguro"
	}
	return s
}

func GenerateToken(userID int64, perfilID int64) (string, error) {
	claims := jwt.MapClaims{
		"user_id":   userID,
		"perfil_id": perfilID,
		"exp":       time.Now().Add(time.Hour * 24).Unix(), // 24 horas
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

func ValidateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de firma inesperado: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("token inválido")
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header requerido"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Formato de token inválido (Bearer <token>)"})
			c.Abort()
			return
		}

		claims, err := ValidateToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido o expirado"})
			c.Abort()
			return
		}

		// Setear variables en contexto
		c.Set("user_id", claims["user_id"])
		c.Set("perfil_id", claims["perfil_id"])
		c.Next()
	}
}

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		perfilID, exists := c.Get("perfil_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No autenticado"})
			c.Abort()
			return
		}

		// Convertir a float64 (JWT usa float para números JSON)
		pID, ok := perfilID.(float64)
		if !ok || int64(pID) != 1 { // Asumiendo ID 1 es Admin
			c.JSON(http.StatusForbidden, gin.H{"error": "Acceso denegado: Se requieren permisos de administrador"})
			c.Abort()
			return
		}
		c.Next()
	}
}
