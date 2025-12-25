package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jgutierrez746/clase_7_gin_bun/config"
	"github.com/jgutierrez746/clase_7_gin_bun/db"
	auth "github.com/jgutierrez746/clase_7_gin_bun/jwt"
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
	// Carga Zona horaria Chile
	config.Init()

	// Cargar variables de entorno desde -env
	if err := godotenv.Load(); err != nil {
		log.Println("No se encontró .env, usando valores por defecto")
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
	/*
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := db.CreateTable(ctx, &modelos.TematicasModel{}); err != nil {
			log.Fatal(err)
		}
		if err := db.CreateTable(ctx, &modelos.PeliculasModel{}); err != nil {
			log.Fatal(err)
		}
		if err := db.CreateTable(ctx, &modelos.PeliculaTematicaModel{}); err != nil {
			log.Fatal(err)
		}
		if err := db.CreateTable(ctx, &modelos.PortadaPeliculaModel{}); err != nil {
			log.Fatal(err)
		}
		if err := db.CreateTable(ctx, &modelos.PerfilesModel{}); err != nil {
			log.Fatal(err)
		}
		if err := db.CreateTable(ctx, &modelos.UsuariosModel{}); err != nil {
			log.Fatal(err)
		}
	*/

	// Agregar FKs individuales (genérico: llama tantas como necesites)
	/*
		if err := db.AgregarFK(ctx, config.Tablas["pt"], "p_id", config.Tablas["pl"], "id", "CASCADE"); err != nil {
			panic(err)
		}
		if err := db.AgregarFK(ctx, config.Tablas["pt"], "tematica_id", config.Tablas["tm"], "id", "CASCADE"); err != nil {
			panic(err)
		}
		if err := db.AgregarFK(ctx, config.Tablas["u"], "perfil_id", config.Tablas["p"], "id", ""); err != nil {
			panic(err)
		}
	*/

	// Configurar Gin en modo release (sin logs verbose)
	gin.SetMode(gin.ReleaseMode)

	// Crear router
	router := gin.Default()

	// Definición de Rutas HTTP
	// Ruta para archivos estaticos
	router.Static("/fotos", "./public/upload/fotos")
	router.Static("/imagenes", "./public/upload/portadas")

	// Grupo prefijo
	apiV1 := router.Group(prefijo)
	{
		apiV1.POST("/login", rutas.Login) // Ruta publica

		// Grupo protegido general
		protected := apiV1.Group("/")
		protected.Use(auth.AuthMiddleware()) // Middleware de autenticación global para estos grupos
		{

			tematicasGroup := protected.Group("/tematicas")
			{
				tematicasGroup.GET("", rutas.ConsultarTematicas)
				tematicasGroup.GET("/:id", rutas.ConsultarTematicasPorId)
				tematicasGroup.POST("", rutas.CrearTematica)
				tematicasGroup.PUT("/:id", rutas.EditarTematica)
				tematicasGroup.DELETE("/:id", rutas.EliminarTematica)
			}

			peliculasGroup := protected.Group("/peliculas")
			{
				peliculasGroup.GET("", rutas.ConsultarPeliculas)
				peliculasGroup.GET("/:id", rutas.ConsultarPeliculaPorId)
				peliculasGroup.POST("", rutas.CrearPelicula)
				peliculasGroup.PUT("/:id", rutas.EditarPelicula)
				peliculasGroup.DELETE("/:id", rutas.EliminarPelicula)

				tematicasPeliculaGroup := peliculasGroup.Group("/:id/tematicas")
				{
					tematicasPeliculaGroup.GET("", rutas.ConsultarTematicasPelicula)
					tematicasPeliculaGroup.POST("", rutas.CrearTematicasPelicula)
					tematicasPeliculaGroup.DELETE("/:idt", rutas.EliminarTematicaPelicula)
				}

				portadaPeliculaGroup := peliculasGroup.Group("/:id/portada")
				{
					portadaPeliculaGroup.GET("", rutas.ConsultarPortadasPelicula)
					portadaPeliculaGroup.POST("", rutas.CrearPortada)
					portadaPeliculaGroup.DELETE("/:idf", rutas.EliminarPortada)
				}
			}

			// Grupo Admin
			adminGroup := protected.Group("/")
			adminGroup.Use(auth.AdminMiddleware())
			{
				perfilesGroup := adminGroup.Group("/perfiles")
				{
					perfilesGroup.GET("", rutas.ConsultarPerfiles)
					perfilesGroup.GET("/:id", rutas.ConsultarPerfilPorId)
					perfilesGroup.POST("", rutas.CrearPerfil)
					perfilesGroup.PUT("/:id", rutas.EditarPerfil)
					perfilesGroup.DELETE("/:id", rutas.EliminarPerfil)
				}

				usuariosGroup := adminGroup.Group("/usuarios")
				{
					usuariosGroup.GET("", rutas.ConsultarUsuarios)
					usuariosGroup.GET("/:id", rutas.ConsultarUsuarioPorId)
					usuariosGroup.POST("", rutas.CrearUsuario)
					usuariosGroup.PUT("/:id", rutas.EditarUsuario)
					usuariosGroup.DELETE("/:id", rutas.EliminarUsuario)
				}
			}
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
