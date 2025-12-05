package dto

import "time"

type PeliculaTematicasSelectOne struct {
	PID        int64     `json:"p_id"`
	TematicaID int64     `json:"tematica_id"`
	Orden      int       `json:"orden"`
	CreatedAt  time.Time `json:"created_at"`
}

type PeliculaTematicasAllSelect []PeliculaTematicasSelectOne

type PeliculaTematicasInsert struct {
	PID        int64     `json:"p_id" binding:"omitempty" bun:"p_id"`
	TematicaID int64     `json:"tematica_id" binding:"required"`
	Orden      int       `json:"orden" binding:"required"`
	CreatedAt  time.Time `json:"created_at" binding:"omitempty"`
}
