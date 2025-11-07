package modelos

import (
	"time"

	"github.com/uptrace/bun"
)

type TematicasModel struct {
	bun.BaseModel `bun:"table:tematicas"`

	ID        int64     `bun:",pk,autoincrement" json:"id"`
	Nombre    string    `bun:"nombre,notnull" json:"nombre" binding:"required"`
	Slug      string    `bun:"slug,notnull,unique" json:"slug" binding:"required"`
	CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
}
