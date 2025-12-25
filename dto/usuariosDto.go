package dto

import "time"

type UsuarioPerfilDTO struct {
	ID           int64     `json:"id" bun:"id"`
	Nombre       string    `json:"nombre" bun:"nombre"`
	Correo       string    `json:"correo" bun:"correo"`
	Telefono     string    `json:"telefono" bun:"telefono"`
	PerfilID     int64     `json:"perfil_id" bun:"perfil_id"`
	PerfilNombre string    `json:"perfil" bun:"perfil_nombre"` // Join column
	CreatedAt    time.Time `json:"created_at" bun:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" bun:"updated_at"`
}

type UsuarioInsert struct {
	ID       int64  `json:"id,omitempty" bun:"id,pk,autoincrement"`
	Nombre   string `json:"nombre" binding:"required"`
	Correo   string `json:"correo" binding:"required,email"`
	Telefono string `json:"telefono" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
	PerfilID int64  `json:"perfil_id" binding:"required"`
}

type UsuarioUpdate struct {
	ID        int64     `json:"id,omitempty" bun:"id"`
	Nombre    string    `json:"nombre,omitempty"`
	Correo    string    `json:"correo,omitempty" binding:"omitempty,email"`
	Telefono  string    `json:"telefono,omitempty"`
	Password  string    `json:"password,omitempty" binding:"omitempty,min=6"`
	PerfilID  int64     `json:"perfil_id,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

type LoginDTO struct {
	Correo   string `json:"correo" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}
