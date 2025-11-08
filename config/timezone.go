package config

import (
	"log"
	"time"
)

var Chilelocation *time.Location

func Init() {
	// Carga la zona horaria de Chile al inicio del paquete (se ejecuta antes de main)
	var err error
	Chilelocation, err = time.LoadLocation("America/Santiago")
	if err != nil {
		log.Fatal("Error cargando zona horaria de Chile: ", err)
	}
}
