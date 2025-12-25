package dto

type LoginDTO struct {
	Correo   string `json:"correo" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}
