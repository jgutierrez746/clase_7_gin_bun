package modelos

import (
	"time"

	"github.com/uptrace/bun"
)

type TematicasModel struct {
	bun.BaseModel `bun:"table:tematicas"`

	ID        int64     `bun:",pk,autoincrement" json:"id,omitempty"`
	Nombre    string    `bun:"nombre,notnull" json:"nombre" binding:"required"`
	Slug      string    `bun:"slug,notnull,unique" json:"slug"`
	CreatedAt time.Time `bun:",nullzero,type:timestamp,notnull,default:current_timestamp" json:"created_at,omitempty"`
	UpdatedAt time.Time `bun:",nullzero,type:timestamp,notnull,default:current_timestamp,on_update:current_timestamp" json:"updated_at,omitempty"`
}
