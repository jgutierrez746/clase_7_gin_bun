package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mysqldialect"
)

// DB es la instancia global de Bun para MySQL
var DB *bun.DB

// InitDB inicializa la conexión a MySQL usando Bun.
// ds: Cadena de conexión, e: "user:pass@tcp(localhost:3306)/dbname?parseTime=true".
func InitDB(dsn string) error {
	// Abre la conexión con el driver estándar de MySQL.
	sqldb, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("error abriendo conexión MySQL: %w", err)
	}

	// Crea la instancia de Bun con el dialecto MySQL.
	db := bun.NewDB(sqldb, mysqldialect.New())

	// Prueba de conexión.
	if err := db.Ping(); err != nil {
		return fmt.Errorf("error conectando a MySQL: %w", err)
	}

	log.Println("Conexión a MySQL exitosa con Bun.")
	DB = db
	return nil
}

// SelectAll realiza un SELECT de todas las filas de una tabla y las escanea en un slice de structs.
// Ej: var users []User; err := SelectAll(ctx, "users", &users)
func SelectAll(ctx context.Context, table string, dest interface{}) error {
	if DB == nil {
		return fmt.Errorf("db no inicializada") // Se debe llamar a InitDB primero si este error ocurre.
	}
	return DB.NewSelect().Table(table).Scan(ctx, dest)
}

// SelectOne realiza un SELECT de una sola fila de una tabla con clausula WHERE.
// Ej: var users User; err := SelectOne(ctx, "users", &user, "id = ?", 1)
func SelectOne(ctx context.Context, table string, dest interface{}, where string, args ...interface{}) error {
	if DB == nil {
		return fmt.Errorf("DB no inicializada") // Se debe llamar a InitDB primero si este error ocurre.
	}
	q := DB.NewSelect().Table(table).Where(where, args...)
	return q.Scan(ctx, dest)
}

// Update actualiza filas en una tabla usando un modelo (struct) y cláusula WHERE.
// Retorna el número de filas afectadas (int64).
// Ej: user.Name = "Nuevo"; affected, err := Update(ctx, "users", &user, "id = ?", user.ID)
func Update(ctx context.Context, table string, model interface{}, where string, args ...interface{}) (int64, error) {
	if DB == nil {
		return 0, fmt.Errorf("DB no inicializada") // Se debe llamar a InitDB primero si este error ocurre.
	}
	q := DB.NewUpdate().Table(table).Model(model).Where(where, args...)
	res, err := q.Exec(ctx)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// Delete borra filas de una tabla con cláusula WHERE.
// Retorna el número de filas afectadas (int64).
// Ej: affected, err := Delete(ctx, "users", "id = ?", 1)
func Delete(ctx context.Context, table string, where string, args ...interface{}) (int64, error) {
	if DB == nil {
		return 0, fmt.Errorf("DB no inicializada") // Se debe llamar a InitDB primero si este error ocurre.
	}
	q := DB.NewDelete().Table(table).Where(where, args...)
	res, err := q.Exec(ctx)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// SelectWithJoin realiza un SELECT con JOINs en una tabla principal, escaneando en un slice de structs.
// Ej:
//   - joins: []string{"LEFT JOIN orders ON orders.user_id = users.id"}
//   - where: "users.id > ?", 10
//
// Nota: El modelo (dest) debe mapear todos los campos de las tablas unidas.
func SelectConJoin(ctx context.Context, mainTable string, joins []string, dest interface{}, where string, args ...interface{}) error {
	if DB == nil {
		return fmt.Errorf("DB no inicializada") // Se debe llamar a InitDB primero si este error ocurre.
	}
	q := DB.NewSelect().Table(mainTable)

	// Agregar JOINs dinámicamente.
	for _, join := range joins {
		q = q.Join(join)
	}

	// WHERE opcional.
	if where != "" {
		q = q.Where(where, args...)
	}
	return q.Scan(ctx, dest)
}
