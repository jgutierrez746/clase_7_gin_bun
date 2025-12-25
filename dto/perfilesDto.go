package dto

type PerfilesSelectDTO struct {
	ID     int64  `json:"id" bun:"id"`
	Nombre string `json:"nombre" bun:"nombre"`
}

type PerfilesAllSelect []PerfilesSelectDTO

type PerfilesInsert struct {
	ID     int64  `json:"id,omitempty" bun:"id"`
	Nombre string `json:"nombre" binding:"required"`
}

type PerfilesUpdate struct {
	ID     int64  `json:"id,omitempty" bun:"id"`
	Nombre string `json:"nombre" binding:"required"`
}
