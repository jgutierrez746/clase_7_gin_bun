package dto

import "time"

type PortadaSelectDTO struct {
	ID            int64     `json:"id" bun:"id"`
	PID           int64     `json:"p_id" bun:"p_id"`
	NombreArchivo string    `json:"nombre_archivo" bun:"nombre_archivo"`
	Url           string    `json:"url" bun:"-"` // Campo calculado, no en BD
	CreatedAt     time.Time `json:"created_at" bun:"created_at"`
}

type PortadaInsertDTO struct {
	PID           int64     `json:"p_id" bun:"p_id"`
	NombreArchivo string    `json:"nombre_archivo" bun:"nombre_archivo"`
	CreatedAt     time.Time `json:"created_at" bun:",type:timestamp,default:current_timestamp"`
}
