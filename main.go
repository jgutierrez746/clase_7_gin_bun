package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jgutierrez746/clase_7_gin_bun/db"
	"github.com/jgutierrez746/clase_7_gin_bun/modelos"
	"github.com/jgutierrez746/clase_7_gin_bun/rutas"
	"github.com/joho/godotenv"
)

/*
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
*/

var prefijo = "/api/v1"

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

	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbServer := os.Getenv("DB_SERVER")
	dbPort := os.Getenv("DB_PORT")
	dbParseTime := os.Getenv("DB_PARSE_TIME")

	variablesVacias := func(valores ...string) bool {
		for _, v := range valores {
			if v == "" {
				return true
			}
		}
		return false
	}(dbName, dbUser, dbPassword, dbServer, dbPort, dbParseTime)

	if variablesVacias {
		log.Fatal("Las variables de entorno no están bien definidas!")
	}

	// DSN de conexión a MySQL
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=%s", dbUser, dbPassword, dbServer, dbPort, dbName, dbParseTime) // ParseTime: triue para manejar timestamps como time.Time
	// log.Printf("DSN armado: %s", strings.ReplaceAll(dsn, dbPassword, "***"))

	if err := db.InitDB(dsn); err != nil {
		log.Fatal("Error initDB: ", err)
	}

	// Crear tablas en BDD
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := db.CreateTable(ctx, "", &modelos.TematicasModel{}); err != nil {
		log.Fatal(err)
	}

	// Configurar Gin en modo release (sin logs verbose)
	gin.SetMode(gin.ReleaseMode)

	// Crear router
	router := gin.Default()

	// Definición de Rutas HTTP
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
		tematicasGroup := apiV1.Group("/tematicas")
		{
			tematicasGroup.GET("", rutas.ConsultarTematicas)
			tematicasGroup.GET("/:id", rutas.ConsultarTematicasPorId)
			tematicasGroup.POST("", rutas.CrearTematica)
		}
		/*
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
		*/
	}

	// Iniciar servidor
	fmt.Printf("servidor iniciado en http://localhost:%d\n", port)
	if err := router.Run(fmt.Sprintf(":%d", port)); err != nil {
		log.Fatal("Error al iniciar el servidor: ", err)
	}
}
