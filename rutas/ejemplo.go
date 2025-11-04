package rutas

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// Funciones de ejemplo, para ver funcionamiento de Gin
// Grupo usuarios
// Este será manejado como GET /users
func GetUsers(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"mensaje": "Lista de usuarios",
		"users":   []string{"Juan", "María", "Pedro"},
	})
}

// Este será manejado como POST /users
func CreateUser(c *gin.Context) {
	var jsonData struct {
		Nombre string `json:"nombre" binding:"required"`
		Email  string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&jsonData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"mensaje": "Usuario creado",
		"user":    jsonData,
	})
}

// Grupo admin
// AdminOnly se maneja con GET /admin/dashboard (ejemploc con middleware simulado)
func AdminOnly(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"mensaje": "Dashboard admin - acceso restringido",
	})
}

// Ruta general sin grupo
// Función Saludar
func Saludar(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"mensaje": "Hola desde Gin!",
	})
}

func Saludar_con_nombre(c *gin.Context) {
	nombre := c.Param("nombre")
	if strings.TrimSpace(nombre) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"mensaje": "Debes enviar un nombre como parametro!",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"mensaje": "Hola " + strings.TrimSpace(nombre),
	})
}

func Query_string(c *gin.Context) {
	id := c.Query("id")
	slug := c.Query("slug")
	c.JSON(http.StatusOK, gin.H{
		"mensaje": "Parámetros Query String | id = " + id + " | slug = " + slug,
	})
}

// Ejemplo de carga de archivos al servidor
func Ejemplo_upload(c *gin.Context) {
	file, err := c.FormFile("foto")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"mensaje": "Ocurrió un error inesperado!",
		})
		return
	}

	var extension = strings.Split(file.Filename, ".")[1]
	unixTime := time.Now().Unix()
	nombreArchivo := strconv.FormatInt(unixTime, 10) + "." + extension
	var archivo string = "public/upload/fotos/" + nombreArchivo

	c.SaveUploadedFile(file, archivo)

	c.JSON(http.StatusOK, gin.H{
		"mensaje":     "Archivo cargado exitosamente",
		"nombre_foto": nombreArchivo,
	})
}
