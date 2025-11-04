package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jgutierrez746/clase_7_gin_bun/rutas"
	"github.com/joho/godotenv"
)

var prefijo = "/api/v1"

func adminMiddleware(c *gin.Context) {
	if c.GetHeader("Authorization") != "admin-token" { // ejmplo de token
		c.JSON(http.StatusUnauthorized, gin.H{
			"mensaje": "Acceso denegado",
		})
		c.Abort()
		return
	}
	c.Next()
}

func main() {
	// Cargar variables de entorno desde -env
	if err := godotenv.Load(); err != nil {
		log.Println("No se encontró .env, usando valores por defeto")
	}

	// Obtener puerto de .env o default 8085
	portStr := os.Getenv("PORT")
	port, err := strconv.Atoi(portStr)
	if err != nil || port == 0 {
		port = 8085
	}

	// Configurar Gin en modo release (sin logs verbose)
	gin.SetMode(gin.ReleaseMode)

	// Crear router
	router := gin.Default()

	// Ruta para archivos estaticos
	router.Static("/fotos", "./public/upload/fotos")

	// Ruta general sin grupo - varios ejemplos
	router.GET("/hola", rutas.Saludar)
	router.GET("/hola/:nombre", rutas.Saludar_con_nombre)
	router.GET("/query-string", rutas.Query_string)
	router.POST("/upload", rutas.Ejemplo_upload)

	// Grupo prefijo
	apiV1 := router.Group(prefijo)
	{
		// Grupo 1 users
		usersGroup := apiV1.Group("/users")
		{
			usersGroup.GET("", rutas.GetUsers)    // GET /api/v1/users
			usersGroup.POST("", rutas.CreateUser) // POST /api/v1/users
		}

		// Grupo 2 admin con middleware
		adminGroup := apiV1.Group("/admin", adminMiddleware)
		{
			adminGroup.GET("/dashboard", rutas.AdminOnly) // GET /admin/dashboard
		}
	}

	// Iniciar servidor
	fmt.Printf("servidor iniciado en http://localhost:%d\n", port)
	if err := router.Run(fmt.Sprintf(":%d", port)); err != nil {
		log.Fatal("Error al iniciar el servidor: ", err)
	}
}
