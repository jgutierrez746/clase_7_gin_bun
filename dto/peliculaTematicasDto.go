package dto

import "time"

type PeliculaTematicasSelectOne struct {
	PID        int64     `json:"p_id"`
	TematicaID int64     `json:"tematica_id"`
	Orden      int       `json:"orden"`
	CreatedAt  time.Time `json:"created_at"`
}

type PeliculaTematicasAllSelect []PeliculaTematicasSelectOne
